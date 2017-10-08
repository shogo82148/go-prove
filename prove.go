package prove

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/shogo82148/go-prove/formatter"
	"github.com/shogo82148/go-prove/test"
	"github.com/soh335/sliceflag"
)

type Prove struct {
	FlagSet *flag.FlagSet

	// Run N test jobs in parallel
	Jobs int

	// Execution command for running test
	Exec string

	Formatter Formatter
	formatter string

	Merge   bool
	Plugins []Plugin

	chanTests  chan *test.Test
	chanSuites chan *test.Test
	wgWorkers  *sync.WaitGroup
	pluginArgs []string
	version    bool
	help       bool

	Mutex    *sync.Mutex
	ExitCode int
}

type Formatter interface {
	// Called to create a new test
	OpenTest(test *test.Test)

	// Prints the report after all tests are run
	Report()
}

type Plugin interface {
	Run(w *Worker, f func())
}

type PluginLoader interface {
	Load(name, args string) Plugin
}

type PluginLoaderFunc func(name, args string) Plugin

func (f PluginLoaderFunc) Load(name, args string) Plugin {
	return f(name, args)
}

var pluginLoaders map[string]PluginLoader = map[string]PluginLoader{}

func AppendPluginLoader(name string, loader PluginLoader) {
	pluginLoaders[name] = loader
}

func NewProve() *Prove {
	p := &Prove{
		FlagSet:    flag.NewFlagSet("prove", flag.ExitOnError),
		Mutex:      &sync.Mutex{},
		ExitCode:   0,
		Plugins:    []Plugin{},
		chanTests:  make(chan *test.Test),
		chanSuites: make(chan *test.Test),
		wgWorkers:  &sync.WaitGroup{},
	}
	p.FlagSet.IntVar(&p.Jobs, "j", 1, "shorthand of -jobs option")
	p.FlagSet.IntVar(&p.Jobs, "jobs", 1, "Run N test jobs in parallel")
	p.FlagSet.BoolVar(&p.version, "V", false, "shorthand of -version option")
	p.FlagSet.BoolVar(&p.version, "version", false, "Show version of go-prove")
	p.FlagSet.StringVar(&p.Exec, "e", "perl", "shorthand of -exec option")
	p.FlagSet.StringVar(&p.Exec, "exec", "perl", "Interpreter to run the tests")
	p.FlagSet.StringVar(&p.formatter, "formatter", "", "Result formatter to use")
	p.FlagSet.BoolVar(&p.Merge, "merge", false, "Merge test scripts' STDERR with their STDOUT")
	p.FlagSet.BoolVar(&p.help, "h", false, "show this help")
	p.FlagSet.BoolVar(&p.help, "help", false, "show this help")
	p.FlagSet.BoolVar(&p.help, "?", false, "show this help")
	sliceflag.StringVar(p.FlagSet, &p.pluginArgs, "plugin", []string{}, "plugins")
	sliceflag.StringVar(p.FlagSet, &p.pluginArgs, "P", []string{}, "plugins")
	return p
}

func (p *Prove) ParseArgs(args []string) {
	p.FlagSet.Parse(args)

	for _, plugin := range p.pluginArgs {
		a := strings.SplitN(plugin, "=", 2)
		name := a[0]
		pluginArgs := ""
		if len(a) >= 2 {
			pluginArgs = a[1]
		}

		loader, ok := pluginLoaders[name]
		if !ok {
			panic("plugin " + name + " not found")
		}
		p.Plugins = append(p.Plugins, loader.Load(name, pluginArgs))
	}
}

func (p *Prove) Run(args []string) {
	if args != nil {
		p.ParseArgs(args)
	}

	if p.version {
		fmt.Printf("go-prove %s, %s built for %s/%s\n", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return
	}

	if p.help {
		p.FlagSet.PrintDefaults()
		return
	}

	files := p.FindTestFiles()

	if p.Jobs < 1 {
		p.Jobs = 1
	}
	for i := 0; i < p.Jobs; i++ {
		w := NewWorker(p, i)
		w.Start()
	}

	switch p.formatter {
	case "junit":
		p.Formatter = &formatter.JUnitFormatter{}
	case "prove":
		fallthrough
	case "":
		p.Formatter = &formatter.TapFormatter{}
	default:
		panic(fmt.Sprintf("unknown formatter: %s", p.formatter))
	}

	go func() {
		for _, path := range files {
			p.chanTests <- &test.Test{
				Path:  path,
				Env:   []string{},
				Exec:  p.Exec,
				Merge: p.Merge,
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

	// clean up plugins
	for i := range p.Plugins {
		if c, ok := p.Plugins[len(p.Plugins)-1-i].(io.Closer); ok {
			c.Close()
		}
	}
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
