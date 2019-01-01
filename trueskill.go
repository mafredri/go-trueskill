package trueskill

import (
	"errors"
	"fmt"
	"math"

	"github.com/mafredri/go-trueskill/collection"
	"github.com/mafredri/go-trueskill/gaussian"
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
	mu              float64 // Mean
	sigma           float64 // Standard deviation
	beta            float64 // Skill class width (length of skill chain)
	tau             float64 // Additive dynamics factor
	drawProbability float64 // Probability of a draw, between zero and a one
}

func (ts Config) String() string {
	return fmt.Sprintf("TrueSkill(mu=%.3f sigma=%.3f beta=%.3f tau=%.3f draw=%.1f%%)", ts.mu, ts.sigma, ts.beta, ts.tau, ts.drawProbability*100)
}

var (
	errDrawProbabilityOutOfRange = errors.New("draw probability must be between 0 and 100")
)

// Option represents a configuration option.
type Option func(c *Config)

// Mu sets the mean.
func Mu(mu float64) Option {
	return func(c *Config) {
		c.mu = mu
	}
}

// Sigma sets the standard deviation.
func Sigma(sigma float64) Option {
	return func(c *Config) {
		c.sigma = sigma
	}
}

// Beta sets the skill class width (length of skill chain).
func Beta(beta float64) Option {
	return func(c *Config) {
		c.beta = beta
	}
}

// Tau sets the additive dynamics factor.
func Tau(tau float64) Option {
	return func(c *Config) {
		c.tau = tau
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
		c.drawProbability = prob
	}, nil
}

// DrawProbabilityZero returns an Option that sets the draw probability
// to zero. Provides a convenient way to set the draw probability to
// zero without checking for errors.
func DrawProbabilityZero() Option {
	return func(c *Config) {
		c.drawProbability = 0
	}
}

// New creates a new TrueSkill configuration with default configuration.
// The configuration can be changed by providing one or multiple Option.
func New(opts ...Option) Config {
	c := Config{
		mu:              DefaultMu,
		sigma:           DefaultSigma,
		beta:            DefaultBeta,
		tau:             DefaultTau,
		drawProbability: DefaultDrawProbability,
	}
	for _, o := range opts {
		o(&c)
	}

	// Always represent the draw probability as a decimal value.
	c.drawProbability /= 100

	return c
}

// AdjustSkillsWithDraws returns the new skill level distribution for all provided
// players based on game configuration and draw status.
// For a N-player game, the draws parameter should have length n-1, where draws[i]
// represents whether player[i] and player[i+1] are in draw.
func (ts Config) AdjustSkillsWithDraws(players []Player, draws []bool) (newSkills []Player, probability float64) {
	// panic if draws slice length is not as expected
	if len(draws) != len(players)-1 {
		panic(fmt.Sprintf(
			"draws slice should have length %d but have %d instead",
			len(players)-1, len(draws)))
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

// AdjustSkills returns the new skill level distribution for all provided
// players based on game configuration and draw status.
// This function can only accept draw as bool which means all players have the
// same ranking. If you need to accept individual player draw state, please call
// AdjustSkillWithDraws.
func (ts Config) AdjustSkills(players []Player, draw bool) (newSkills []Player, probability float64) {
	draws := make([]bool, len(players)-1)
	for i := range draws {
		draws[i] = draw
	}

	return ts.AdjustSkillsWithDraws(players, draws)
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

// NewPlayer returns a new player with the mu and sigma from the game
// configuration.
func (ts Config) NewPlayer() Player {
	return NewPlayer(ts.mu, ts.sigma)
}

// TrueSkill returns the conservative TrueSkill of a player. The maximum
// TrueSkill is two times mu, in the default configuration a value between
// zero and fifty is returned.
func (ts Config) TrueSkill(p Player) float64 {
	trueSkill := p.Mu() - (ts.mu/ts.sigma)*p.Sigma()

	return math.Min(ts.mu*2, math.Max(0, trueSkill))
}
