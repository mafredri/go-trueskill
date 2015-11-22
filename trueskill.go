// Package trueskill implements the TrueSkill™ ranking system (by Microsoft) in Go.
package trueskill

import (
	"errors"
	"math"

	"github.com/mafredri/go-gaussian"
	"github.com/mafredri/go-trueskill/collection"
	"github.com/mafredri/go-trueskill/schedule"
)

// Constants for the TrueSkill ranking system.
const (
	DefaultMu                 = 25.0
	DefaultSigma              = DefaultMu / 3.0
	DefaultBeta               = DefaultSigma * 0.5
	DefaultTau                = DefaultSigma * 0.01
	DefaultDrawProbPercentage = 10.0

	loopMaxDelta = 1e-4 // Desired accuracy for factor graph loop schedule
)

// Config is the configuration for the TrueSkill ranking system
type Config struct {
	Mu       float64 // Mean
	Sigma    float64 // Standard deviation
	Beta     float64 // Skill class width (length of skill chain)
	Tau      float64 // Additive dynamics factor
	DrawProb float64 // Probability of a draw, between zero and a one
}

// New creates a new TrueSkill configuration from the provided values
func New(mu, sigma, beta, tau, drawProbPercentage float64) (Config, error) {
	if drawProbPercentage < 0.0 || drawProbPercentage > 100.0 {
		return Config{}, errors.New("Draw probability must be between 0 and 100.")
	}
	return Config{mu, sigma, beta, tau, drawProbPercentage / 100}, nil
}

// NewDefault returns a new game configuration with the default TrueSkill configuration.
func NewDefault(drawProbPercentage float64) (Config, error) {
	return New(DefaultMu, DefaultSigma, DefaultBeta, DefaultTau, drawProbPercentage)
}

// AdjustSkills returns the new skill level distribution for all provided players based on game configuration and draw
// status.
func (ts Config) AdjustSkills(players Players, draw bool) (Players, float64) {
	// Sort players
	// sort.Sort(players)

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

// MatchQuality returns a float representing the quality of the match-up between players.
//
// Only two player match quality is supported at this time. Minus one is returned if the match-up is unsupported.
func (ts Config) MatchQuality(players Players) float64 {
	if players.Len() > 2 {
		return -1
	}

	return calculate2PlayerMatchQuality(ts, players[0], players[1])
}

// NewDefaultPlayer returns a new player with the mu and sigma from the game configuration.
func (ts Config) NewDefaultPlayer() Player {
	return NewPlayer(ts.Mu, ts.Sigma)
}

// TrueSkill returns the conservative TrueSkill of a player. The maximum TrueSkill is two times mu, in the default
// configuration a value between zero and fifty is returned.
func (ts Config) TrueSkill(p Player) int64 {
	trueSkill := p.Mu() - p.Sigma()*3
	trueSkill = math.Ceil(trueSkill)

	return int64(math.Min(ts.Mu*2, math.Max(0, trueSkill)))
}
