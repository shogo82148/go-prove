package prove

import (
	"log"
	"os"
)

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
			log.Printf("start %s", test.Path)
			w.prove.chanSuites <- test.Run()
			log.Printf("finish %s", test.Path)
		}
		w.prove.wgWorkers.Done()
	}

	for _, p := range w.prove.Plugins {
		f = func(g func()) func() {
			return func() { p.Run(w, g) }
		}(f)
	}

	f()
}
