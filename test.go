package prove

import (
	"io"
	"os"
	"os/exec"
	"syscall"

	"fmt"

	"github.com/mattn/go-shellwords"
	"github.com/shogo82148/go-tap"
)

type Test struct {
	Path string
	Env  []string
	Exec string

	// Merge test scripts' STDERR with their STDOUT.
	Merge bool

	Suite *tap.Testsuite
}

func (t *Test) Run() *tap.Testsuite {
	execParam, _ := shellwords.Parse(t.Exec)
	execParam = append(execParam, t.Path)
	cmd := exec.Command(execParam[0], execParam[1:]...)
	cmd.Env = t.Env

	r, w := io.Pipe()
	cmd.Stdout = w

	if t.Merge {
		cmd.Stderr = w
	} else {
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Start(); err != nil {
		t.Suite = errorTestsuite(err)
		return t.Suite
	}

	ch := make(chan *tap.Testsuite, 1)
	go func() {
		parser, err := tap.NewParser(r)
		if err != nil {
			ch <- errorTestsuite(err)
			return
		}
		suite, err := parser.Suite()
		if err != nil {
			ch <- errorTestsuite(err)
			return
		}
		ch <- suite
	}()

	err := cmd.Wait()
	w.Close()
	r.Close()

	suite := <-ch
	t.Suite = suite

	// check exit code
	if err == nil {
		return suite
	}
	// from http://qiita.com/hnakamur/items/5e6f22bda8334e190f63
	if e2, ok := err.(*exec.ExitError); ok {
		if s, ok := e2.Sys().(syscall.WaitStatus); ok {
			exitCode := s.ExitStatus()
			if exitCode != 0 {
				suite.Ok = false
				suite.Plan++
				suite.Tests = append(suite.Tests, &tap.Testline{
					Ok:          false,
					Num:         suite.Plan,
					Description: fmt.Sprintf("Test died with return code %d", exitCode),
				})
			}
			return suite
		}
	}

	suite.Ok = false
	suite.Plan++
	suite.Tests = append(suite.Tests, &tap.Testline{
		Ok:          false,
		Num:         suite.Plan,
		Description: "unexpected error",
		Diagnostic:  err.Error(),
	})
	return suite
}

func errorTestsuite(err error) *tap.Testsuite {
	return &tap.Testsuite{
		Ok: false,
		Tests: []*tap.Testline{
			&tap.Testline{
				Ok:          false,
				Num:         1,
				Description: "unexpected error",
				Diagnostic:  err.Error(),
			},
		},
		Plan:    1,
		Version: tap.DefaultTAPVersion,
	}
}
