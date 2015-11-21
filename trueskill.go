// Package trueskill implements the TrueSkill™ algorithm created by Microsoft in Go.
package trueskill

import (
	"errors"
	"math"
	"sort"

	"github.com/mafredri/go-gaussian"
	"github.com/mafredri/go-trueskill/collection"
	"github.com/mafredri/go-trueskill/schedule"
)

// Constants for the TrueSkill algorithm
const (
	defaultMu    = 25.0
	defaultSigma = defaultMu / 3.0
	defaultBeta  = defaultSigma * 0.5
	defaultTau   = defaultSigma * 0.01
	loopMaxDelta = 1e-4 // Desired accuracy for factor graph loop schedule
)

// TSConf represents a configuration for a game
type TSConf struct {
	Mu       float64 // Mean
	Sigma    float64 // Standard deviation
	Beta     float64 // Skill class width (length of skill chain)
	Tau      float64 // Additive dynamics factor
	DrawProb float64 // Probability of a draw, 0.0 - 100.0
}

// New creates a new TrueSkill configuration from the provided values
func New(mu, sigma, beta, tau, drawProbPercentage float64) (TSConf, error) {
	if drawProbPercentage < 0.0 || drawProbPercentage > 100.0 {
		return TSConf{}, errors.New("Draw probability must be between 0 and 100.")
	}
	return TSConf{mu, sigma, beta, tau, drawProbPercentage / 100}, nil
}

// NewWithDefaults creates a new game configuration from the default TrueSkill configuration
func NewWithDefaults(drawProbPercentage float64) (TSConf, error) {
	return New(defaultMu, defaultSigma, defaultBeta, defaultTau, drawProbPercentage)
}

// TrueSkill returns the conservative true skill for a player
func (ts TSConf) TrueSkill(p Player) int64 {
	trueSkill := p.Mu() - p.Sigma()*3
	trueSkill = math.Ceil(trueSkill)

	return int64(math.Min(ts.Mu*2, math.Max(0, trueSkill)))
}

// NewPlayerWithDefaults creates a new player with default mu / sigma from the game configuration
func (ts TSConf) NewPlayerWithDefaults() Player {
	return NewPlayer(ts.Mu, ts.Sigma)
}

// MatchQuality calculates and returns the match quality
func (ts TSConf) MatchQuality(players Players) (float64, error) {
	if players.Len() > 2 {
		return 0, errors.New("Match quality currently supports only 2 players.")
	}

	return calculate2PlayerMatchQuality(ts, players[0], players[1]), nil
}

// AdjustSkills calculates the new skill level distribution for all provided players based on game configuration and
// draw status
func (ts TSConf) AdjustSkills(players Players, draw bool) (Players, float64) {
	// Sort
	sort.Sort(players)
	draws := []bool{}
	for i := 0; i < players.Len()-1; i++ {
		draws = append(draws, draw)
	}

	// TODO: Rewrite distribution bag
	prior := gaussian.NewFromPrecision(0, 0)
	varBag := collection.NewDistributionBag(prior)

	skillFactors, skillIndex, factorList := buildSkillFactors(ts, players, draws, varBag)

	sched := buildSkillFactorSchedule(players.Len(), skillFactors, loopMaxDelta)

	// delta
	_ = schedule.Run(sched, -1)

	logZ := factorList.LogNormalization()
	probability := math.Exp(logZ)

	newPlayerSkills := Players{}
	for _, id := range skillIndex {
		newPlayerSkills = append(newPlayerSkills, Player{Gaussian: varBag.Get(id)})
	}

	return newPlayerSkills, probability
}
