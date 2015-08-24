package Formatter

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/shogo82148/go-tap"
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
	Contents string `xml:",chardata"`
}

func (f *JUnitFormatter) formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.3f", d.Seconds())
}

func (f *JUnitFormatter) OpenTest(suite *tap.Testsuite) {
	ts := JUnitTestSuite{
		Tests:      0,
		Failures:   0,
		Time:       f.formatDuration(suite.Time),
		Name:       "hoge",
		Properties: []JUnitProperty{},
		TestCases:  []JUnitTestCase{},
	}

	for _, line := range suite.Tests {
		testCase := JUnitTestCase{
			Classname: "fuga",
			Name:      line.Description,
			Time:      f.formatDuration(line.Time),
			Failure:   nil,
		}
		if !line.Ok {
			ts.Failures++
			testCase.Failure = &JUnitFailure{
				Message:  "not ok",
				Type:     "",
				Contents: line.Diagnostic,
			}
		}
		ts.Tests++
		ts.TestCases = append(ts.TestCases, testCase)
	}

	f.Suites.Suites = append(f.Suites.Suites, ts)
}

func (f *JUnitFormatter) Report() {
	bytes, _ := xml.MarshalIndent(f.Suites, "", "\t")
	os.Stdout.Write(bytes)
}
