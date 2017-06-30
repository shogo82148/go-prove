package main

import (
	"os"

	prove "github.com/shogo82148/go-prove"
	formatter "github.com/shogo82148/go-prove/formatter"
	_ "github.com/shogo82148/go-prove/plugin"
)

func main() {
	p := prove.NewProve()
	p.Formatter = &formatter.JUnitFormatter{}
	p.ParseArgs(os.Args[1:])
	p.Run(nil)

	os.Exit(p.ExitCode)
}
