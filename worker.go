package prove

type Worker struct {
	prove *Prove
}

func NewWorker(p *Prove) *Worker {
	return &Worker{
		prove: p,
	}
}

func (w *Worker) Start() {
	w.prove.wgWorkers.Add(1)
	go w.run()
}

func (w *Worker) run() {
	for test := range w.prove.chanTests {
		w.prove.chanSuites <- test.Run()
	}
	w.prove.wgWorkers.Done()
}
