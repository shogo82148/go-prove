package Formatter

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/shogo82148/go-prove"
	tap "github.com/shogo82148/go-tap"
)

type JUnitFormatter struct {
	Suites JUnitTestSuites
}

// JUnitTestSuites is a collection of JUnit test suites.
type JUnitTestSuites struct {
	XMLName xml.Name `xml:"testsuites"`
	Suites  []JUnitTestSuite
}

// JUnitTestSuite is a single JUnit test suite which may contain many
// testcases.
type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Errors     int             `xml:"errors,attr"`
	Skipped    int             `xml:"skipped,attr"`
	Time       string          `xml:"time,attr"`
	Name       string          `xml:"name,attr"`
	Properties []JUnitProperty `xml:"properties>property,omitempty"`
	TestCases  []JUnitTestCase
}

// JUnitTestCase is a single test case with its result.
type JUnitTestCase struct {
	XMLName     xml.Name          `xml:"testcase"`
	Classname   string            `xml:"classname,attr"`
	Name        string            `xml:"name,attr"`
	Time        string            `xml:"time,attr"`
	SkipMessage *JUnitSkipMessage `xml:"skipped,omitempty"`
	Failure     *JUnitFailure     `xml:"failure,omitempty"`
	SystemOut   *JUnitSystemOut   `xml:"system-out,omitempty"`
	SystemErr   *JUnitSystemErr   `xml:"system-err,omitempty"`
}

// JUnitSkipMessage contains the reason why a testcase was skipped.
type JUnitSkipMessage struct {
	Message string `xml:"message,attr"`
}

// JUnitProperty represents a key/value pair used to define properties.
type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// JUnitFailure contains data related to a failed test.
type JUnitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",cdata"`
}

// JUnitSystemOut contains the standard out.
type JUnitSystemOut struct {
	Contents string `xml:",cdata"`
}

// JunitSystemErr contains the standard error.
type JUnitSystemErr struct {
	Contents string `xml:",cdata"`
}

func (f *JUnitFormatter) formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.3f", d.Seconds())
}

func (f *JUnitFormatter) OpenTest(test *prove.Test) {
	className := strings.Replace(test.Path, "/", "_", -1)
	className = strings.Replace(className, ".", "_", -1)

	suite := test.Suite

	ts := JUnitTestSuite{
		Time: f.formatDuration(suite.Time),
		Name: className,
	}

	for _, line := range suite.Tests {
		testCase := JUnitTestCase{
			Classname: className,
			Name:      line.Description,
			Time:      f.formatDuration(line.Time),
			SystemOut: &JUnitSystemOut{
				Contents: line.GoString(),
			},
		}
		if !line.Ok {
			ts.Failures++
			testCase.Failure = &JUnitFailure{
				Message:  line.String(),
				Type:     "TestFailed",
				Contents: line.Diagnostic,
			}
		}
		if line.Directive == tap.Skip {
			ts.Skipped++
			testCase.SkipMessage = &JUnitSkipMessage{
				Message: line.Diagnostic,
			}
		}
		ts.Tests++
		ts.TestCases = append(ts.TestCases, testCase)
	}

	if suite.Plan < 0 {
		ts.Errors++
		testCase := JUnitTestCase{
			Classname: className,
			Name:      "Test died too soon, even before plan.",
			Time:      "0.000",
			Failure: &JUnitFailure{
				Message:  "The test suite died before a plan was produced. You need to have a plan.",
				Type:     "Plan",
				Contents: "No plan",
			},
		}
		ts.TestCases = append(ts.TestCases, testCase)
	} else if len(suite.Tests) != suite.Plan {
		ts.Errors++
		testCase := JUnitTestCase{
			Classname: className,
			Name:      "Number of runned tests does not match plan.",
			Time:      "0.000",
			Failure: &JUnitFailure{
				Message:  "Some test were not executed, The test died prematurely.",
				Type:     "Plan",
				Contents: "Bad plan",
			},
		}
		ts.TestCases = append(ts.TestCases, testCase)
	}

	f.Suites.Suites = append(f.Suites.Suites, ts)
}

func (f *JUnitFormatter) Report() {
	out := os.Stdout
	io.WriteString(out, xml.Header)
	enc := xml.NewEncoder(out)
	enc.Indent("", "    ")
	enc.Encode(f.Suites)
	enc.Flush()
}
