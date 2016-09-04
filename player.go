package trueskill

import "github.com/mafredri/go-gaussian"

// Player is a player with a certain skill (mu, sigma).
type Player struct {
	Rank int
	gaussian.Gaussian
}

// NewPlayer returns a player from the provided mu (mean) and sigma (standard
// deviation).
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

// Players is a list of players (skills).
type Players []Player

func (p Players) Len() int           { return len(p) }
func (p Players) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Players) Less(i, j int) bool { return p[i].Rank < p[j].Rank }
