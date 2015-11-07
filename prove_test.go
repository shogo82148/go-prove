package prove

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestParseArgs(t *testing.T) {
	p := NewProve()
	p.ParseArgs([]string{"-j", "100", "--exec", "foobar"})

	if p.Jobs != 100 {
		t.Errorf("want 100\ngot %d", p.Jobs)
	}
	if p.Exec != "foobar" {
		t.Errorf("want foobar\ngot %s", p.Exec)
	}
}

func TestFindTestFiles(t *testing.T) {
	// create dummy test files
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	//	defer os.RemoveAll(dir)
	ioutil.WriteFile(filepath.Join(dir, "foo.t"), []byte{}, 0644)
	ioutil.WriteFile(filepath.Join(dir, "foo.pl"), []byte{}, 0644)
	os.MkdirAll(filepath.Join(dir, "foo", "bar"), 0777)
	ioutil.WriteFile(filepath.Join(dir, "foo", "bar", "foo.t"), []byte{}, 0644)

	// pass directory name
	{
		p := NewProve()
		p.ParseArgs([]string{dir})
		testFiles := p.FindTestFiles()
		expected := []string{filepath.Join(dir, "foo.t"), filepath.Join(dir, "foo", "bar", "foo.t")}
		sort.Strings(testFiles)
		sort.Strings(expected)
		if !reflect.DeepEqual(testFiles, expected) {
			t.Errorf("want %v\ngot %v", expected, testFiles)
		}
	}

	// pass file name
	{
		p := NewProve()
		p.ParseArgs([]string{filepath.Join(dir, "foo.t")})
		testFiles := p.FindTestFiles()
		expected := []string{filepath.Join(dir, "foo.t")}
		if !reflect.DeepEqual(testFiles, expected) {
			t.Errorf("want %v\ngot %v", expected, testFiles)
		}
	}
}
