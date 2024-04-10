package test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_success(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "foo.t")
	if err := os.WriteFile(filename, []byte(`print "1..1\nok 1\n";`), 0644); err != nil {
		t.Fatal(err)
	}

	test := &Test{
		Path: filename,
		Env:  os.Environ(),
		Exec: "perl",
	}

	suite := test.Run()
	if suite == nil {
		t.Error("want not nil\ngot nil")
	}
	if test.Suite != suite {
		t.Errorf("want that %v equals to %v", test.Suite, suite)
	}

	if !suite.Ok {
		t.Error("want success\ngot fail")
	}
	if len(suite.Tests) != 1 {
		t.Errorf("want 1\ngot %d", len(suite.Tests))
	}
}

func TestRun_fail(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "foo.t")
	if err := os.WriteFile(filename, []byte(`print "1..1\nnot ok 1\n";`), 0644); err != nil {
		t.Fatal(err)
	}

	test := &Test{
		Path: filename,
		Env:  os.Environ(),
		Exec: "perl",
	}

	suite := test.Run()
	if suite == nil {
		t.Error("want not nil\ngot nil")
	}
	if test.Suite != suite {
		t.Errorf("want that %v equals to %v", test.Suite, suite)
	}

	if suite.Ok {
		t.Error("want fail\ngot success")
	}
	if len(suite.Tests) != 1 {
		t.Errorf("want 1\ngot %d", len(suite.Tests))
	}
}

func TestRun_failplan(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "foo.t")
	if err := os.WriteFile(filename, []byte(`print "1..2\nok 1\n";`), 0644); err != nil {
		t.Fatal(err)
	}

	test := &Test{
		Path: filename,
		Env:  os.Environ(),
		Exec: "perl",
	}

	suite := test.Run()
	if suite == nil {
		t.Error("want not nil\ngot nil")
	}
	if test.Suite != suite {
		t.Errorf("want that %v equals to %v", test.Suite, suite)
	}

	if suite.Ok {
		t.Error("want fail\ngot success")
	}
	if len(suite.Tests) != 1 {
		t.Errorf("want 1\ngot %d", len(suite.Tests))
	}
}

func TestRun_empty(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "foo.t")
	if err := os.WriteFile(filename, []byte(`die "test failed!!!";`), 0644); err != nil {
		t.Fatal(err)
	}

	test := &Test{
		Path: filename,
		Env:  os.Environ(),
		Exec: "perl",
	}

	suite := test.Run()
	if suite == nil {
		t.Error("want not nil\ngot nil")
	}
	if test.Suite != suite {
		t.Errorf("want that %v equals to %v", test.Suite, suite)
	}

	if suite.Ok {
		t.Error("want fail\ngot success")
	}
	if len(suite.Tests) != 2 {
		t.Errorf("want 2\ngot %d", len(suite.Tests))
	}
}

func TestRun_exitNonZero(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "foo.t")
	if err := os.WriteFile(filename, []byte(`print "1..1\nok 1\n"; exit 1;`), 0644); err != nil {
		t.Fatal(err)
	}

	test := &Test{
		Path: filename,
		Env:  os.Environ(),
		Exec: "perl",
	}

	suite := test.Run()
	if suite == nil {
		t.Error("want not nil\ngot nil")
	}
	if test.Suite != suite {
		t.Errorf("want that %v equals to %v", test.Suite, suite)
	}

	if suite.Ok {
		t.Error("want fail\ngot success")
	}
	if len(suite.Tests) != 2 {
		t.Errorf("want 2\ngot %d", len(suite.Tests))
	}

	if !suite.Tests[0].Ok {
		t.Error("want success\ngot fail")
	}
	if suite.Tests[1].Ok {
		t.Error("want fail\ngot success")
	}
}
