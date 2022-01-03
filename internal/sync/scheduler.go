package sync

type Scheduler interface {
	// Starts a new co-routine and tracks it in this scheduler
	NewCoroutine(ctx Context, fn func(Context))

	// Execute executes all coroutines until they are all blocked
	Execute(ctx Context)

	RunningCoroutines() int

	Exit(ctx Context)
}

type scheduler struct {
	coroutines []Coroutine
}

func NewScheduler() Scheduler {
	return &scheduler{
		coroutines: make([]Coroutine, 0),
	}
}

func (s *scheduler) NewCoroutine(ctx Context, fn func(Context)) {
	c := NewCoroutine(ctx, fn)
	s.coroutines = append(s.coroutines, c)
}

func (s *scheduler) Execute(ctx Context) {
	allBlocked := false
	for !allBlocked {
		allBlocked = true
		for i := 0; i < len(s.coroutines); i++ {
			c := s.coroutines[i]

			c.Execute()

			if c.Finished() {
				// Coroutine is finished, remove from list
				s.coroutines[i] = nil
				s.coroutines = append(s.coroutines[:i], s.coroutines[i+1:]...)
				i--
			} else {
				// Determine if coroutine made any progress or if it stayed blocked
				allBlocked = allBlocked && !c.Progress()
			}
		}
	}
}

func (s *scheduler) RunningCoroutines() int {
	return len(s.coroutines)
}

func (s *scheduler) Exit(_ Context) {
	for _, c := range s.coroutines {
		c.Exit()
	}
}