package trueskill

import (
	"testing"

	"github.com/mafredri/go-trueskill/mathextra"
)

const defaultEpsilon = 1e-5 // Precision for floating point comparison

func testPlayerSkillsWithErrorMargin(
	t *testing.T, playerSkills []Player, wantSkill []float64, epsilon float64) {
	for i, p := range playerSkills {
		if !mathextra.Float64AlmostEq(p.Mu(), wantSkill[i*2], epsilon) {
			t.Errorf("p%d.Mu() == %.5f, want %.5f", i, p.Mu(), wantSkill[i*2])
		}
		if !mathextra.Float64AlmostEq(p.Sigma(), wantSkill[i*2+1], epsilon) {
			t.Errorf("p%d.Sigma() == %.5f, want %.5f", i, p.Sigma(), wantSkill[i*2+1])
		}
	}
}

func testPlayerSkills(t *testing.T, playerSkills []Player, wantSkills []float64) {
	testPlayerSkillsWithErrorMargin(t, playerSkills, wantSkills, defaultEpsilon)
}

func testProbability(t *testing.T, probability, wantProbability float64) {
	probability = probability * 100
	if !mathextra.Float64AlmostEq(probability, wantProbability, defaultEpsilon) {
		t.Errorf("Probability == %.5f, want %.5f", probability, wantProbability)
	}
}

func TestTrueSkill_HeadToHead(t *testing.T) {
	wantSkill := []float64{
		29.3958320199992000, 7.1714755873261900,
		20.6041679800008000, 7.1714755873261900,
	}
	wantProbability := 47.7593111421005000

	draw := false
	ts := New()

	players := []Player{ts.NewPlayer(), ts.NewPlayer()}

	newPlayerSkills, probability := ts.AdjustSkills(players, draw)

	testPlayerSkills(t, newPlayerSkills, wantSkill)
	testProbability(t, probability, wantProbability)
}

func TestTrueSkill_HeadToHead_Draw(t *testing.T) {
	wantSkill := []float64{
		25.0000000000000000, 6.4575196623173100,
		25.0000000000000000, 6.4575196623173100,
	}
	wantProbability := 4.4813777157991800

	draw := true
	ts := New()

	players := []Player{ts.NewPlayer(), ts.NewPlayer()}

	newPlayerSkills, probability := ts.AdjustSkills(players, draw)

	testPlayerSkills(t, newPlayerSkills, wantSkill)
	testProbability(t, probability, wantProbability)
}

func TestTrueSkill_HeadToHead_NoDrawProbability(t *testing.T) {
	wantSkill := []float64{
		29.2054731068140000, 7.1948165514807000,
		20.7945268931860000, 7.1948165514807000,
	}
	wantProbability := 50.0000008292022000

	drawProbability, err := DrawProbability(0.0)
	if err != nil {
		t.Error(err)
	}
	draw := false

	ts := New(drawProbability)

	players := []Player{ts.NewPlayer(), ts.NewPlayer()}

	newPlayerSkills, probability := ts.AdjustSkills(players, draw)

	testPlayerSkills(t, newPlayerSkills, wantSkill)
	testProbability(t, probability, wantProbability)
}

func TestTrueSkill_HeadToHead_BetterPlayerWins(t *testing.T) {
	wantSkill := []float64{
		31.2295178021319000, 6.5230309127748200,
		18.7704821978681000, 6.5230309127748200,
	}
	wantProbability := 75.3779324752223000

	drawProbability, err := DrawProbability(10.0)
	if err != nil {
		t.Error(err)
	}
	draw := false

	ts := New(drawProbability)

	players := []Player{NewPlayer(29.396, 7.171), NewPlayer(20.604, 7.171)}

	newPlayerSkills, probability := ts.AdjustSkills(players, draw)

	testPlayerSkills(t, newPlayerSkills, wantSkill)
	testProbability(t, probability, wantProbability)
}

func TestTrueSkill_HeadToHead_BetterPlayerLoses(t *testing.T) {
	wantSkill := []float64{
		26.6428086088974000, 6.0399862622030400,
		23.3571913911026000, 6.0399862622030400,
	}
	wantProbability := 20.8198643622167000

	drawProbability, err := DrawProbability(10.0)
	if err != nil {
		t.Error(err)
	}
	draw := false

	ts := New(drawProbability)

	players := []Player{NewPlayer(20.604, 7.171), NewPlayer(29.396, 7.171)}

	newPlayerSkills, probability := ts.AdjustSkills(players, draw)

	testPlayerSkills(t, newPlayerSkills, wantSkill)
	testProbability(t, probability, wantProbability)
}

func TestTrueSkill_4PFreeForAll(t *testing.T) {
	wantSkill := []float64{
		33.2066809656313000, 6.3481091698077100,
		27.4014546938433000, 5.7871629348447600,
		22.5985453061884000, 5.7871629348413500,
		16.7933190343613000, 6.3481091698145000,
	}
	wantProbability := 3.1576468103466900

	drawProbability, err := DrawProbability(10.0)
	if err != nil {
		t.Error(err)
	}
	draw := false

	ts := New(drawProbability)

	players := []Player{ts.NewPlayer(),
		ts.NewPlayer(),
		ts.NewPlayer(),
		ts.NewPlayer()}

	newPlayerSkills, probability := ts.AdjustSkills(players, draw)

	testPlayerSkills(t, newPlayerSkills, wantSkill)
	testProbability(t, probability, wantProbability)
}

func TestTrueSkill_4PFreeForAll_WithDraws(t *testing.T) {
	wantSkill := []float64{
		28.162, 5.712,
		28.162, 5.712,
		21.836, 5.712,
		21.836, 5.712,
	}
	wantProbability := 0.09406

	drawProbability, err := DrawProbability(10.0)
	if err != nil {
		t.Error(err)
	}

	draws := []bool{
		true,
		false,
		true,
	}

	ts := New(drawProbability)

	players := []Player{
		ts.NewPlayer(),
		ts.NewPlayer(),
		ts.NewPlayer(),
		ts.NewPlayer()}

	newPlayerSkills, probability := ts.AdjustSkillsWithDraws(players, draws)

	// In draw, the skill calculation for two drawed players are not guaranteed
	// to be the same, but should be close. We gives a larger epsilon here as
	// the error margin
	testPlayerSkillsWithErrorMargin(t, newPlayerSkills, wantSkill, 0.01)
	testProbability(t, probability, wantProbability)
}

func TestTrueSkill_8PFreeForAll(t *testing.T) {
	wantSkill := []float64{
		36.7710964365458000, 5.7492832446709400,
		32.2423455786852000, 5.1329106221072500,
		29.0739837977166000, 4.9427131496384800,
		26.3221792432643000, 4.8745473884181200,
		23.6778207568207000, 4.8745473884091200,
		20.9260162023759000, 4.9427131496114100,
		17.7576544214189000, 5.1329106220642000,
		13.2289035635009000, 5.7492832447392400,
	}
	wantProbability := 0.0006565500448901

	drawProbability, err := DrawProbability(10.0)
	if err != nil {
		t.Error(err)
	}
	draw := false

	ts := New(drawProbability)

	players := []Player{ts.NewPlayer(),
		ts.NewPlayer(),
		ts.NewPlayer(),
		ts.NewPlayer(),
		ts.NewPlayer(),
		ts.NewPlayer(),
		ts.NewPlayer(),
		ts.NewPlayer()}

	newPlayerSkills, probability := ts.AdjustSkills(players, draw)

	testPlayerSkills(t, newPlayerSkills, wantSkill)
	testProbability(t, probability, wantProbability)
}

func TestTrueSkill_MatchQuality_HeadToHead(t *testing.T) {
	wantMatchQuality := 44.7

	ts := New()

	players := []Player{ts.NewPlayer(), ts.NewPlayer()}

	matchQuality := ts.MatchQuality(players)
	if matchQuality == -1 {
		t.Error("Match quality was -1")
	}

	matchQuality = matchQuality * 100
	if !mathextra.Float64AlmostEq(matchQuality, wantMatchQuality, 1e-1) {
		t.Errorf("Probability == %.1f, want %.1f", matchQuality, wantMatchQuality)
	}

	players = append(players, ts.NewPlayer())
	matchQuality = ts.MatchQuality(players)
	if matchQuality != -1 {
		t.Errorf("bad match quality for >2 players; got %v, want %v", matchQuality, -1)
	}
}

func TestTrueSkillForPlayer(t *testing.T) {
	ts := New(DrawProbabilityZero())

	player := ts.NewPlayer()
	skill := ts.TrueSkill(player)
	if skill != 0 {
		t.Errorf("wrong trueskill for new player; got %v, want %v", skill, 0)
	}

	player = NewPlayer(30, 1)
	skill = ts.TrueSkill(player)
	if skill != 27 {
		t.Errorf("wrong trueskill for new player; got %v, want %v", skill, 27)
	}
}
