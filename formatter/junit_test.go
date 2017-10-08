package formatter

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/shogo82148/go-prove/test"
)

func TestJUnit_success(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(`print "1..1\nok 1\n";`)

	test := &test.Test{
		Path: f.Name(),
		Env:  os.Environ(),
		Exec: "perl",
	}

	test.Run()

	fmtter := &JUnitFormatter{}
	fmtter.OpenTest(test)
	b, _ := xml.MarshalIndent(fmtter.Suites, "", "")
	re := `^<testsuites><testsuite tests="1" failures="0" errors="0" skipped="0" time="0.[0-9]+" name="[^"]+">` +
		`<properties></properties><testcase classname="[^"]+" name="" time="0.[0-9]+">` +
		`<system-out><!\[CDATA\[ok 1` + "\n" +
		`\]\]></system-out></testcase></testsuite></testsuites>$`
	ok, err := regexp.Match(re, b)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf("incorrect output\n%s", string(b))
	}
}

func TestJUnit_fail(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(`print "1..1\nnot ok 1\n";`)

	test := &test.Test{
		Path: f.Name(),
		Env:  os.Environ(),
		Exec: "perl",
	}

	test.Run()

	fmtter := &JUnitFormatter{}
	fmtter.OpenTest(test)
	b, _ := xml.MarshalIndent(fmtter.Suites, "", "")
	re := `^<testsuites><testsuite tests="1" failures="1" errors="0" skipped="0" time="0.[0-9]+" name="[^"]+">` +
		`<properties></properties><testcase classname="[^"]+" name="" time="0.[0-9]+">` +
		`<failure message="not ok 1" type="TestFailed"></failure>` +
		`<system-out><!\[CDATA\[not ok 1` + "\n" +
		`\]\]></system-out></testcase></testsuite></testsuites>$`
	ok, err := regexp.Match(re, b)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf("incorrect output\n%s", string(b))
	}
}

func TestJUnit_failplan(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(`print "1..2\nok 1\n";`)

	test := &test.Test{
		Path: f.Name(),
		Env:  os.Environ(),
		Exec: "perl",
	}

	test.Run()

	fmtter := &JUnitFormatter{}
	fmtter.OpenTest(test)
	b, _ := xml.MarshalIndent(fmtter.Suites, "", "")
	re := `^<testsuites><testsuite tests="1" failures="0" errors="1" skipped="0" time="0.[0-9]+" name="[^"]+">` +
		`<properties></properties>` +
		`<testcase classname="[^"]+" name="" time="0.[0-9]+"><system-out><!\[CDATA\[ok 1` + "\n" + `\]\]></system-out></testcase>` // +
	//`<testcase classname="[^"]+" name="Number of runned tests does not match plan." time="0.[0-9]+">` +
	//`<failure message="Some test were not executed, The test died prematurely." type="Plan"><!\[CDATA\[Bad plan\]\]</failure>` +
	//`</testcase></testsuite></testsuites>$`
	ok, err := regexp.Match(re, b)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf("incorrect output\n%s", string(b))
	}
}
