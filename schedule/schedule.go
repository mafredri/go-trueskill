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

// Step is a step function in the schedule.
type Step struct {
	input    int // Input for function
	function func(i int) float64
}

// NewStep returns a step with provided function and input.
func NewStep(function func(i int) float64, input int) Step {
	return Step{input, function}
}

// Run runs the step function with the input and returns the delta.
func (s Step) Run(depth, maxDepth int) float64 {
	delta := s.function(s.input)
	return delta
}

// Sequence is a sequence of runnable schedules.
type Sequence struct {
	sequences []RunnableSchedule
}

// NewSequence returns a new sequence of runnable schedules.
func NewSequence(sequences []RunnableSchedule) Sequence {
	return Sequence{sequences}
}

// Run runs all runnable schedules in a sequence. The largest delta from any of the runnable schedules is returned.
func (s Sequence) Run(depth, maxDepth int) float64 {
	var delta float64
	for _, s := range s.sequences {
		delta = math.Max(delta, s.Run(depth+1, maxDepth))
	}

	return delta
}

// Loop is a runnable schedule that runs itself.
type Loop struct {
	schedule RunnableSchedule
	maxDelta float64
}

// NewLoop returns a new loop.
func NewLoop(schedule RunnableSchedule, maxDelta float64) Loop {
	return Loop{schedule, maxDelta}
}

// Run reruns the loop until a desired delta (maxDelta) is reached.
func (sl Loop) Run(depth, maxDepth int) float64 {
	delta := math.MaxFloat64
	for delta > sl.maxDelta {
		delta = sl.schedule.Run(depth+1, maxDepth)
	}

	return delta
}
