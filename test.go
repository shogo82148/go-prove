package prove

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/shogo82148/go-tap"
)

type Test struct {
	Path  string
	Env   []string
	Exec  string
	Suite *tap.Testsuite
}

func (t *Test) Run() *tap.Testsuite {
	execParam := strings.Split(t.Exec, " ")
	execParam = append(execParam, t.Path)
	cmd := exec.Command(execParam[0], execParam[1:]...)
	cmd.Env = t.Env
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		log.Println(err)
	}
	go io.Copy(os.Stderr, stderr)

	parser, _ := tap.NewParser(stdout)
	suite, _ := parser.Suite()

	cmd.Wait()

	t.Suite = suite
	return suite
}
