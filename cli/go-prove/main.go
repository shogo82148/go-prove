package main

import (
	"os"

	prove "github.com/shogo82148/go-prove"
	_ "github.com/shogo82148/go-prove/plugin"
)

func main() {
	p := prove.NewProve()
	p.ParseArgs(os.Args[1:])
	p.Run(nil)

	os.Exit(p.ExitCode)
}
