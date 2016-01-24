package prove

import (
	"io"
	"os"
	"os/exec"

	"github.com/mattn/go-shellwords"
	"github.com/shogo82148/go-tap"
)

type Test struct {
	Path  string
	Env   []string
	Exec  string
	Suite *tap.Testsuite
}

func (t *Test) Run() *tap.Testsuite {
	execParam, _ := shellwords.Parse(t.Exec)
	execParam = append(execParam, t.Path)
	cmd := exec.Command(execParam[0], execParam[1:]...)
	cmd.Env = t.Env
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Suite = errorTestsuite(err)
		return t.Suite
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Suite = errorTestsuite(err)
		return t.Suite
	}

	err = cmd.Start()
	if err != nil {
		t.Suite = errorTestsuite(err)
		return t.Suite
	}
	go io.Copy(os.Stderr, stderr)

	var suite *tap.Testsuite
	parser, err := tap.NewParser(stdout)
	if err != nil {
		suite = errorTestsuite(err)
	} else {
		suite, err = parser.Suite()
		if err != nil {
			suite = errorTestsuite(err)
		}
	}

	cmd.Wait()

	t.Suite = suite
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
