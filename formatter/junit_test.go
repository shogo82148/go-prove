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
	ok, err := regexp.Match(`<testsuites><testsuite tests="1" failures="0" time="0.[0-9]+" name="[^"]*"><properties></properties><testcase classname="[^"]*" name="" time="0.[0-9]+"><failure message="not ok" type=""></failure></testcase></testsuite></testsuites>`, b)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf("incorrect output\n%s", string(b))
	}
}
