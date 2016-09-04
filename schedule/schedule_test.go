package schedule

import (
	"math"
	"testing"
)

func TestScheduleStepRun(t *testing.T) {
	simpleStepFunc := func(i int) float64 {
		return float64(i) + 0.5
	}

	step := NewStep(simpleStepFunc, 1)

	result := step.Run(0, -1)
	want := 1.5

	if result != want {
		t.Errorf("Run(0, -1) == %f, want %f", result, want)
	}
}

func TestScheduleSequenceRun(t *testing.T) {
	simpleStepFunc := func(i int) float64 {
		return float64(i) + 0.5
	}
	seq := NewSequence(
		NewStep(simpleStepFunc, 1),
		NewStep(simpleStepFunc, 2),
		NewStep(simpleStepFunc, 0),
	)

	result := seq.Run(0, -1)
	want := 2.5

	if result != want {
		t.Errorf("Run(0, -1) == %f, want %f", result, want)
	}
}

func TestScheduleLoopRun(t *testing.T) {
	var iter int
	delta := 5.0
	simpleStepReduceFunc := func(i int) float64 {
		delta = math.Max(delta-0.5, 0)
		iter++
		return delta
	}
	seq := NewSequence(
		NewStep(simpleStepReduceFunc, 1),
		NewStep(simpleStepReduceFunc, 1),
		NewStep(simpleStepReduceFunc, 1),
	)

	loop := NewLoop(seq, 0)

	result := loop.Run(0, -1)
	wantResult := 0.0
	wantIter := 3 * 4 // 4 iter of 3*0.5 to reduce 5.0 -> 0.0

	if result != wantResult {
		t.Errorf("Run(0, -1) == %f, want %f", result, wantResult)
	}

	if iter != wantIter {
		t.Errorf("Run(0, -1) iter == %d, want %d", iter, wantIter)
	}
}

func TestRunSchedule(t *testing.T) {
	var iter int
	delta := 10.0
	simpleStepReduceFunc := func(i int) float64 {
		delta = math.Max(delta-float64(i), 0)
		iter++
		return delta
	}
	simpleStepFunc := func(i int) float64 {
		return float64(i)
	}
	sequence := NewSequence(
		NewSequence(NewStep(simpleStepFunc, -1)),
		NewSequence(NewStep(simpleStepFunc, -2)),
		NewLoop(NewStep(simpleStepReduceFunc, 1), 0),
	)

	result := Run(sequence, -1)
	wantResult := 0.0
	wantIter := 10

	if result != wantResult {
		t.Errorf("Run(sequence, -1) == %f, want %f", result, wantResult)
	}

	if iter != wantIter {
		t.Errorf("Run(sequence, -1) iter == %d, want %d", iter, wantIter)
	}
}
