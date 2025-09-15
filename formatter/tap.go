package formatter

import (
	"fmt"

	"github.com/shogo82148/go-prove/test"
	tap "github.com/shogo82148/go-tap"
)

// TapFormatter formats the test result into TAP https://testanything.org/
type TapFormatter struct {
	Suites []*tap.Testsuite
}

// OpenTest implements prove.Formatter
func (f *TapFormatter) OpenTest(test *test.Test) {
	f.Suites = append(f.Suites, test.Suite)
}

// Report implements prove.Formatter
func (f *TapFormatter) Report() error {
	for _, s := range f.Suites {
		for _, t := range s.Tests {
			if _, err := fmt.Printf("%#v", t); err != nil {
				return err
			}
		}
	}
	return nil
}
