package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := withDone(in, done)
	for _, s := range stages {
		if s != nil {
			out = s(withDone(out, done))
		}
	}

	return out
}

func withDone(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			default:
			}

			select {
			case <-done:
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				out <- value
			}
		}
	}()

	return out
}
