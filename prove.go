package prove

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Prove struct {
	FlagSet *flag.FlagSet

	// Run N test jobs in parallel
	Jobs int

	// Execution command for running test
	Exec string

	Formatter Formatter

	Plugins []Plugin

	chanTests  chan *Test
	chanSuites chan *Test
	wgWorkers  *sync.WaitGroup

	Mutex    *sync.Mutex
	ExitCode int
}

type Formatter interface {
	// Called to create a new test
	OpenTest(test *Test)

	// Prints the report after all tests are run
	Report()
}

type Plugin interface {
	Run(w *Worker, f func())
}

func NewProve() *Prove {
	p := &Prove{
		FlagSet:    flag.NewFlagSet("prove", flag.ExitOnError),
		Mutex:      &sync.Mutex{},
		ExitCode:   0,
		Plugins:    []Plugin{},
		chanTests:  make(chan *Test),
		chanSuites: make(chan *Test),
		wgWorkers:  &sync.WaitGroup{},
	}
	p.FlagSet.IntVar(&p.Jobs, "j", 1, "Run N test jobs in parallel")
	p.FlagSet.IntVar(&p.Jobs, "job", 1, "Run N test jobs in parallel")
	p.FlagSet.StringVar(&p.Exec, "exec", "perl", "")
	return p
}

func (p *Prove) ParseArgs(args []string) {
	p.FlagSet.Parse(args)
}

func (p *Prove) Run(args []string) {
	if args != nil {
		p.ParseArgs(args)
	}

	files := p.FindTestFiles()

	if p.Jobs < 1 {
		p.Jobs = 1
	}
	for i := 0; i < p.Jobs; i++ {
		w := NewWorker(p)
		w.Start()
	}

	go func() {
		for _, path := range files {
			p.chanTests <- &Test{
				Path: path,
				Env:  []string{},
				Exec: p.Exec,
			}
		}
		close(p.chanTests)
		p.wgWorkers.Wait()
		close(p.chanSuites)
	}()

	for suite := range p.chanSuites {
		p.Formatter.OpenTest(suite)
	}
	p.Formatter.Report()
}

// Find Test Files
func (p *Prove) FindTestFiles() []string {
	files := []string{}
	if p.FlagSet.NArg() == 0 {
		files = p.findTestFiles(files, "t")
	} else {
		for _, parent := range p.FlagSet.Args() {
			files = p.findTestFiles(files, parent)
		}
	}
	return files
}

func (p *Prove) findTestFiles(files []string, parent string) []string {
	stat, _ := os.Stat(parent)
	if !stat.IsDir() {
		return append(files, parent)
	}

	filepath.Walk(
		parent,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(path, ".t") {
				files = append(files, path)
			}

			return nil
		})
	return files
}

func (p *Prove) MarkAsFail() {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.ExitCode = 1
}
