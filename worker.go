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
			test.Run()
			if !test.Suite.Ok {
				w.prove.MarkAsFail()
			}
			w.prove.chanSuites <- test
			log.Printf("finish %s", test.Path)
		}
	}

	for _, p := range w.prove.Plugins {
		pp := p
		g := f
		f = func() {
			pp.Run(w, g)
		}
	}

	f()
	w.prove.wgWorkers.Done()
}
