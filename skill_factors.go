package trueskill

import (
	"github.com/mafredri/go-trueskill/collection"
	"github.com/mafredri/go-trueskill/factor"
	sched "github.com/mafredri/go-trueskill/schedule"
)

type skillFactors struct {
	skillPriorFactors                        []factor.Factor
	playerPerformances                       []int
	skillToPerformanceFactors                []factor.Factor
	playerPerformanceDifferences             []int
	performanceToPerformanceDifferencFactors []factor.Factor
	greatherThanOrWithinFactors              []factor.Factor
}

func buildSkillFactors(ts Config, players Players, draws []bool, varBag collection.DistributionBag) (skillFactors, []int, factor.List) {
	sf := skillFactors{}
	gf := factor.NewGaussianFactors()
	factorList := factor.NewList()

	numPlayers := players.Len()

	skillIndex := []int{}
	for i := 0; i < numPlayers; i++ {
		skillIndex = append(skillIndex, varBag.NextIndex())
	}

	for i := 0; i < numPlayers; i++ {
		priorSkill := players[i]
		gpf := gf.GaussianPrior(priorSkill.Mean(), priorSkill.Variance()+(ts.Tau*ts.Tau), skillIndex[i], varBag)
		sf.skillPriorFactors = append(sf.skillPriorFactors, gpf)
		factorList.Add(gpf)
	}

	for i := 0; i < numPlayers; i++ {
		sf.playerPerformances = append(sf.playerPerformances, varBag.NextIndex())
	}

	for i := 0; i < numPlayers; i++ {
		glf := gf.GaussianLikeliehood(ts.Beta*ts.Beta, sf.playerPerformances[i], skillIndex[i], varBag, varBag)
		sf.skillToPerformanceFactors = append(sf.skillToPerformanceFactors, glf)
		factorList.Add(glf)
	}

	for i := 0; i < numPlayers-1; i++ {
		sf.playerPerformanceDifferences = append(sf.playerPerformanceDifferences, varBag.NextIndex())
	}

	for i := 0; i < numPlayers-1; i++ {
		gws := gf.GaussianWeightedSum(1.0, -1.0, sf.playerPerformanceDifferences[i], sf.playerPerformances[i],
			sf.playerPerformances[i+1], varBag, varBag, varBag)
		sf.performanceToPerformanceDifferencFactors = append(sf.performanceToPerformanceDifferencFactors, gws)
		factorList.Add(gws)
	}

	// TODO: Calculate e separately for each
	epsilon := drawMargin(ts.Beta, ts.DrawProb)
	for i, draw := range draws {
		var f factor.Factor
		if draw {
			f = gf.GaussianWithin(epsilon, sf.playerPerformanceDifferences[i], varBag)
		} else {
			f = gf.GaussianGreaterThan(epsilon, sf.playerPerformanceDifferences[i], varBag)
		}
		sf.greatherThanOrWithinFactors = append(sf.greatherThanOrWithinFactors, f)
		factorList.Add(f)
	}

	return sf, skillIndex, factorList
}

func skillFactorListToScheduleStep(facs []factor.Factor, idx int) []sched.RunnableSchedule {
	steps := []sched.RunnableSchedule{}
	for _, f := range facs {
		step := sched.NewStep(f.UpdateMessage, idx)
		steps = append(steps, step)
	}

	return steps
}

// buildSkillFactorSchedule builds a full schedule that represents all the steps in a factor graph.
func buildSkillFactorSchedule(numPlayers int, sf skillFactors, loopMaxDelta float64) sched.RunnableSchedule {
	// Prior schedule initializes the skill priors for all players and updates the performance
	priorSchedule := sched.NewSequence([]sched.RunnableSchedule{
		sched.NewSequence(skillFactorListToScheduleStep(sf.skillPriorFactors, 0)),
		sched.NewSequence(skillFactorListToScheduleStep(sf.skillToPerformanceFactors, 0)),
	})

	// Loop schedule iterates until desired accuracy is reached
	var loopSchedule sched.RunnableSchedule

	if numPlayers == 2 {
		// In two player mode there is no to loop, just send the performance difference and the greater-than
		loopSchedule = sched.NewSequence([]sched.RunnableSchedule{
			sched.NewStep(sf.performanceToPerformanceDifferencFactors[0].UpdateMessage, 0),
			sched.NewStep(sf.greatherThanOrWithinFactors[0].UpdateMessage, 0),
		})
	} else {
		// Forward schedule updates the factor graph in one direction
		forwardSchedule := []sched.RunnableSchedule{}
		// ... and the backward schedule in the other direction
		backwardSchedule := []sched.RunnableSchedule{}

		for i := 0; i < numPlayers-2; i++ {
			forwardSteps := []sched.RunnableSchedule{
				sched.NewStep(sf.performanceToPerformanceDifferencFactors[i].UpdateMessage, 0),
				sched.NewStep(sf.greatherThanOrWithinFactors[i].UpdateMessage, 0),
				sched.NewStep(sf.performanceToPerformanceDifferencFactors[i].UpdateMessage, 2),
			}
			forwardSchedule = append(forwardSchedule, forwardSteps...)

			backwardSteps := []sched.RunnableSchedule{
				sched.NewStep(sf.performanceToPerformanceDifferencFactors[numPlayers-2-i].UpdateMessage, 0),
				sched.NewStep(sf.greatherThanOrWithinFactors[numPlayers-2-i].UpdateMessage, 0),
				sched.NewStep(sf.performanceToPerformanceDifferencFactors[numPlayers-2-i].UpdateMessage, 1),
			}
			backwardSchedule = append(backwardSchedule, backwardSteps...)
		}

		// Combine the backward and forward schedule so that they are run in said order
		combinedForwardBackwardSchedule := sched.NewSequence([]sched.RunnableSchedule{
			sched.NewSequence(forwardSchedule),
			sched.NewSequence(backwardSchedule),
		})

		// Loop through the forward and backward schedule until the delta stops changing by more than loopMaxDelta
		loopSchedule = sched.NewLoop(combinedForwardBackwardSchedule, loopMaxDelta)
	}

	innerSchedule := sched.NewSequence([]sched.RunnableSchedule{
		loopSchedule,
		sched.NewStep(sf.performanceToPerformanceDifferencFactors[0].UpdateMessage, 1),
		sched.NewStep(sf.performanceToPerformanceDifferencFactors[numPlayers-2].UpdateMessage, 2),
	})

	// Finally send the skill performances of all players
	posteriorSchedule := sched.NewSequence(skillFactorListToScheduleStep(sf.skillToPerformanceFactors, 1))

	// Combine all schedules into one runnable sequence
	fullSchedule := sched.NewSequence([]sched.RunnableSchedule{
		priorSchedule, innerSchedule, posteriorSchedule,
	})

	return fullSchedule
}
