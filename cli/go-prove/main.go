package main

import (
	"os"
	"runtime"

	prove "github.com/shogo82148/go-prove"
	formatter "github.com/shogo82148/go-prove/formatter"
)

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	p := prove.NewProve()
	p.Formatter = &formatter.JUnitFormatter{}
	p.ParseArgs(os.Args[1:])
	p.Run(nil)

	os.Exit(p.ExitCode)
}
