package prove

import "os"

type Worker struct {
	prove *Prove
	Env   []string
}

func NewWorker(p *Prove) *Worker {
	return &Worker{
		prove: p,
		Env:   os.Environ(),
	}
}

func (w *Worker) Start() {
	w.prove.wgWorkers.Add(1)
	go w.run()
}

func (w *Worker) run() {
	f := func() {
		for test := range w.prove.chanTests {
			test.Env = w.Env
			w.prove.chanSuites <- test.Run()
		}
		w.prove.wgWorkers.Done()
	}

	for _, p := range w.prove.Plugins {
		f = func() { p.Run(w, f) }
	}

	f()
}
