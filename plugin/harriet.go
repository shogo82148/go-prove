package plugin

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/shogo82148/go-prove"
)

type Harriet struct {
	cmd  string
	args []string
}

func init() {
	prove.AppendPluginLoader("harriet", prove.PluginLoaderFunc(func(name, args string) prove.Plugin {
		cmd := "harriet"
		cmdArgs := []string{"t/harriet"}

		a, _ := shellwords.Parse(args)
		if len(a) > 0 {
			cmd = a[0]
			cmdArgs = a[1:]
		}
		return &Harriet{
			cmd:  cmd,
			args: cmdArgs,
		}
	}))
}

func (p *Harriet) Run(w *prove.Worker, f func()) {
	log.Printf("run harriet cmd: %s %s", p.cmd, p.args)

	cmd := exec.Command(p.cmd, p.args...)
	cmd.Env = w.Env

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()
	go io.Copy(os.Stderr, stderr)

	s := bufio.NewScanner(stdout)
	foundExport := false
	for s.Scan() {
		t := s.Text()
		if strings.HasPrefix(t, "export") {
			exportCmd, _ := shellwords.Parse(t)
			if len(exportCmd) < 2 {
				continue
			}
			log.Printf("export %s", exportCmd[1])
			w.Env = append(w.Env, exportCmd[1])
			foundExport = true
		}
		if foundExport && t == "" {
			break
		}
	}

	defer func() {
		cmd.Process.Signal(os.Interrupt)
		cmd.Wait()
	}()

	f()
}
