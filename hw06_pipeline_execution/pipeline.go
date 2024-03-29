package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	wrapper := func(in In, stage Stage) Out {
		out := make(Bi)
		go func() {
			defer func() {
				close(out)
				<-in
			}()
			for {
				select {
				case <-done:
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-done:
						return
					case out <- v:
					}
				}
			}
		}()
		return stage(out)
	}

	for _, stage := range stages {
		in = wrapper(in, stage)
	}

	return in
}
