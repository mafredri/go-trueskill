// Package trueskill implements the TrueSkill™ ranking system (by Microsoft) in
// Go.
package trueskill

import (
	"errors"
	"fmt"
	"math"

	"github.com/mafredri/go-gaussian"
	"github.com/mafredri/go-trueskill/collection"
	"github.com/mafredri/go-trueskill/schedule"
)

// Constants for the TrueSkill ranking system.
const (
	DefaultMu              = 25.0
	DefaultSigma           = DefaultMu / 3.0
	DefaultBeta            = DefaultSigma * 0.5
	DefaultTau             = DefaultSigma * 0.01
	DefaultDrawProbability = 10.0 // Percentage, between 0 and 100.

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

func (ts Config) String() string {
	return fmt.Sprintf("TrueSkill(mu=%.3f sigma=%.3f beta=%.3f tau=%.3f draw=%.1f%%)", ts.Mu, ts.Sigma, ts.Beta, ts.Tau, ts.DrawProb*100)
}

var (
	errDrawProbabilityOutOfRange = errors.New("draw probability must be between 0 and 100")
)

// Option represents a configuration option for TrueSkill.
type Option func(c *Config)

// Mu sets the mean for TrueSkill.
func Mu(mu float64) Option {
	return func(c *Config) {
		c.Mu = mu
	}
}

// Sigma sets the standard deviation for TrueSkill.
func Sigma(sigma float64) Option {
	return func(c *Config) {
		c.Sigma = sigma
	}
}

// Beta sets the skill class width for TrueSkill.
func Beta(beta float64) Option {
	return func(c *Config) {
		c.Beta = beta
	}
}

// Tau sets the additive dynamics factor for TrueSkill.
func Tau(tau float64) Option {
	return func(c *Config) {
		c.Tau = tau
	}
}

// DrawProbability takes a value between 0 and 100 and returns an Option that
// sets the probability of a draw. An error is returned if the input value is
// out of range.
func DrawProbability(prob float64) (Option, error) {
	if prob < 0.0 || prob > 100.0 {
		return nil, errDrawProbabilityOutOfRange
	}
	return func(c *Config) {
		c.DrawProb = prob / 100
	}, nil
}

// DrawProbabilityZero is an Option that sets the draw probability to
// zero. It exists as a convenience function.
func DrawProbabilityZero(c *Config) {
	c.DrawProb = 0
}

// New creates a new TrueSkill configuration with default values, unless
func New(opts ...Option) Config {
	c := Config{
		Mu:       DefaultMu,
		Sigma:    DefaultSigma,
		Beta:     DefaultBeta,
		Tau:      DefaultTau,
		DrawProb: DefaultDrawProbability / 100, // Percentage, between 0 and 100.
	}
	for _, o := range opts {
		o(&c)
	}
	return c
}

// AdjustSkills returns the new skill level distribution for all provided
// players based on game configuration and draw status.
func (ts Config) AdjustSkills(players []Player, draw bool) (newSkills []Player, probability float64) {
	draws := []bool{}
	for i := 0; i < len(players)-1; i++ {
		draws = append(draws, draw)
	}

	// TODO: Rewrite the distribution bag and simplify the factor list as well
	prior := gaussian.NewFromPrecision(0, 0)
	varBag := collection.NewDistributionBag(prior)

	skillFactors, skillIndex, factorList := buildSkillFactors(ts, players, draws, varBag)

	sched := buildSkillFactorSchedule(len(players), skillFactors, loopMaxDelta)

	// delta
	_ = schedule.Run(sched, -1)

	logZ := factorList.LogNormalization()
	probability = math.Exp(logZ)

	for _, id := range skillIndex {
		newSkills = append(newSkills, Player{Gaussian: varBag.Get(id)})
	}

	return newSkills, probability
}

// MatchQuality returns a float representing the quality of the match-up
// between players.
//
// Only two player match quality is supported at this time. Minus one is
// returned if the match-up is unsupported.
func (ts Config) MatchQuality(players []Player) float64 {
	if len(players) > 2 {
		return -1
	}

	return calculate2PlayerMatchQuality(ts, players[0], players[1])
}

// NewDefaultPlayer returns a new player with the mu and sigma from the game
// configuration.
func (ts Config) NewDefaultPlayer() Player {
	return NewPlayer(ts.Mu, ts.Sigma)
}

// TrueSkill returns the conservative TrueSkill of a player. The maximum
// TrueSkill is two times mu, in the default configuration a value between
// zero and fifty is returned.
func (ts Config) TrueSkill(p Player) float64 {
	trueSkill := p.Mu() - p.Sigma()*3

	return math.Min(ts.Mu*2, math.Max(0, trueSkill))
}
