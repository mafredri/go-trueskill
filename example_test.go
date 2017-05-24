package trueskill_test

import (
	"fmt"

	"github.com/mafredri/go-trueskill"
)

func Example() {
	ts := trueskill.New(trueskill.DrawProbabilityZero)
	p1 := ts.NewPlayer()
	p2 := ts.NewPlayer()
	draw := false
	skills := []trueskill.Player{p1, p2}
	newSkills, probability := ts.AdjustSkills(skills, draw)
	p1 = newSkills[0]
	p2 = newSkills[1]

	fmt.Println("Player 1:", p1)
	fmt.Println("Player 2:", p2)
	fmt.Printf("Probability of player 1 winning: %.1f\n", probability*100)
	fmt.Printf("Match quality before: %.3f, after: %.3f\n", ts.MatchQuality(skills), ts.MatchQuality(newSkills))

	// Output:
	// Player 1: Player(mu=29.205 sigma=7.195)
	// Player 2: Player(mu=20.795 sigma=7.195)
	// Probability of player 1 winning: 50.0
	// Match quality before: 0.447, after: 0.388
}
