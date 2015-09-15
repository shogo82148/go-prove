package main

import (
	"os"
	"runtime"

	prove "github.com/shogo82148/go-prove"
	formatter "github.com/shogo82148/go-prove/formatter"
	plugin "github.com/shogo82148/go-prove/plugin"
)

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	p := prove.NewProve()
	p.Formatter = &formatter.JUnitFormatter{}
	p.Plugins = append(p.Plugins, &plugin.TestMysqld{})
	p.ParseArgs(os.Args[1:])
	p.Run(nil)

	os.Exit(p.ExitCode)
}
