/*

Package trueskill implements the TrueSkill™ ranking system (by Microsoft) in Go.

Create a new instance of trueskill with the default configuration:

	ts := trueskill.New()

Create new players with based on the trueskill configuration:

	p1 := ts.NewPlayer() // Same as trueskill.NewPlayer(25.0, 8.333)

The TrueSkill algorithm can be tweaked with configuration options:

	ts := trueskill.New(
		trueskill.Mu(200),
		trueskill.Sigma(66.666),
		trueskill.Beta(33.333),
		trueskill.Tau(0.666),
		trueskill.DrawProbabilityZero())

Adjust player skills:

	ts := trueskill.New()
	p1 := ts.NewPlayer()
	p2 := ts.NewPlayer()
	draw := false
	newSkills, probability := ts.AdjustSkills([]Player{p1, p2}, draw)
	p1 = newSkills[0]
	p2 = newSkills[1]

Check the conservative TrueSkill of a player:

	ts := trueskill.New()
	p1 := trueskill.NewPlayer(30, 1)
	fmt.Println(ts.TrueSkill(p1)) // 27

*/
package trueskill
