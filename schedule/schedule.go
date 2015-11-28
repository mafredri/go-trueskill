// Package schedule provides a scheduler capable of scheduling functions to be run in sequences and loops.
package schedule

import "math"

// RunnableSchedule provides an interface for the run function.
type RunnableSchedule interface {
	Run(depth, maxDepth int) float64
}

// Run runs a schedule starting from zero depth.
func Run(schedule RunnableSchedule, maxDepth int) float64 {
	return schedule.Run(0, maxDepth)
}

type step struct {
	input    int // Input for function
	function func(i int) float64
}

// NewStep returns a step in the schedule with provided function and input.
func NewStep(function func(i int) float64, input int) RunnableSchedule {
	return step{input, function}
}

// Run runs the step function with the input and returns the delta.
func (s step) Run(depth, maxDepth int) float64 {
	delta := s.function(s.input)
	return delta
}

type sequence struct {
	sequences []RunnableSchedule
}

// NewSequence returns a new sequence of runnable schedules.
func NewSequence(sequences []RunnableSchedule) RunnableSchedule {
	return sequence{sequences}
}

// Run runs all runnable schedules in a sequence. The largest delta from any of the runnable schedules is returned.
func (s sequence) Run(depth, maxDepth int) float64 {
	var delta float64
	for _, s := range s.sequences {
		delta = math.Max(delta, s.Run(depth+1, maxDepth))
	}

	return delta
}

type loop struct {
	schedule RunnableSchedule
	maxDelta float64
}

// NewLoop returns a new runnable loop schedule.
func NewLoop(schedule RunnableSchedule, maxDelta float64) RunnableSchedule {
	return loop{schedule, maxDelta}
}

// Run reruns the loop until a desired delta (maxDelta) is reached.
func (l loop) Run(depth, maxDepth int) float64 {
	delta := math.MaxFloat64
	for delta > l.maxDelta {
		delta = l.schedule.Run(depth+1, maxDepth)
	}

	return delta
}
