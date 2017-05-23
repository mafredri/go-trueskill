package trueskill

import (
	"fmt"

	"github.com/mafredri/go-gaussian"
)

// Player is a player with a certain skill (mu, sigma).
type Player struct {
	gaussian.Gaussian
}

// NewPlayer returns a player from the provided mu (mean) and sigma
// (standard deviation).
func NewPlayer(mu, sigma float64) Player {
	return Player{
		Gaussian: gaussian.NewFromMeanAndStdDev(mu, sigma),
	}
}

// Mu returns the player mean.
func (p Player) Mu() float64 {
	return p.Mean()
}

// Sigma returns the player standard deviation.
func (p Player) Sigma() float64 {
	return p.StdDev()
}

func (p Player) String() string {
	return fmt.Sprintf("Player(mu=%.3f sigma=%.3f)", p.Mu(), p.Sigma())
}
