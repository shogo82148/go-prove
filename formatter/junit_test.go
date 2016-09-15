package Formatter

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/shogo82148/go-prove"
)

func TestJUnit_success(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(`print "1..1\nok 1\n";`)

	test := &prove.Test{
		Path: f.Name(),
		Env:  os.Environ(),
		Exec: "perl",
	}

	test.Run()

	formatter := &JUnitFormatter{}
	formatter.OpenTest(test)
	b, _ := xml.MarshalIndent(formatter.Suites, "", "")
	ok, err := regexp.Match(`<testsuites><testsuite tests="1" failures="0" time="0.[0-9]+" name="[^"]*"><properties></properties><testcase classname="[^"]*" name="" time="0.[0-9]+"></testcase></testsuite></testsuites>`, b)
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

	test := &prove.Test{
		Path: f.Name(),
		Env:  os.Environ(),
		Exec: "perl",
	}

	test.Run()

	formatter := &JUnitFormatter{}
	formatter.OpenTest(test)
	b, _ := xml.MarshalIndent(formatter.Suites, "", "")
	ok, err := regexp.Match(`<testsuites><testsuite tests="1" failures="1" time="0.[0-9]+" name="[^"]*"><properties></properties><testcase classname="[^"]*" name="" time="0.[0-9]+"><failure message="not ok" type=""></failure></testcase></testsuite></testsuites>`, b)
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

	test := &prove.Test{
		Path: f.Name(),
		Env:  os.Environ(),
		Exec: "perl",
	}

	test.Run()

	formatter := &JUnitFormatter{}
	formatter.OpenTest(test)
	b, _ := xml.MarshalIndent(formatter.Suites, "", "")
	ok, err := regexp.Match(`<testsuites><testsuite tests="1" failures="0" errors="1" time="0.[0-9]+" name="[^"]*"><properties></properties><testcase classname="[^"]*" name="" time="0.[0-9]+"></testcase><testcase classname="[^"]*" name="Number of runned tests does not match plan." time="0.[0-9]+"><failure message="Some test were not executed, The test died prematurely." type="Plan"><!\[CDATA\[Bad plan\]\]></failure></testcase></testsuite></testsuites>`, b)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf("incorrect output\n%s", string(b))
	}
}

func TestJUnit_fail_with_diagnostic(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(`print "not ok 1\n#   Failed test at -e line 1.\n1..1";`)

	test := &prove.Test{
		Path: f.Name(),
		Env:  os.Environ(),
		Exec: "perl",
	}

	test.Run()

	formatter := &JUnitFormatter{}
	formatter.OpenTest(test)
	b, _ := xml.MarshalIndent(formatter.Suites, "", "")
	ok, err := regexp.Match(`<testsuites><testsuite tests="1" failures="1" time="0.[0-9]+" name="[^"]*"><properties></properties><testcase classname="[^"]*" name="" time="0.[0-9]+"><failure message="not ok" type=""><!\[CDATA\[#   Failed test at -e line 1.
\]\]></failure></testcase></testsuite></testsuites>`, b)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf("incorrect output\n%s", string(b))
	}
}
