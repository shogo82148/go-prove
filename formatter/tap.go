package formatter

import (
	"fmt"

	"github.com/shogo82148/go-prove/test"
	tap "github.com/shogo82148/go-tap"
)

type TapFormatter struct {
	Suites []*tap.Testsuite
}

func (f *TapFormatter) OpenTest(test *test.Test) {
	f.Suites = append(f.Suites, test.Suite)
}

func (f *TapFormatter) Report() {
	for _, s := range f.Suites {
		for _, t := range s.Tests {
			fmt.Printf("%#v", t)
		}
	}
}
