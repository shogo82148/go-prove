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
	c    *exec.Cmd
	env  []string
}

func init() {
	prove.AppendPluginLoader("harriet", prove.PluginLoaderFunc(harrietLoader))
}

func harrietLoader(name, args string) prove.Plugin {
	cmd := "harriet"
	cmdArgs := []string{"t/harriet"}

	a, _ := shellwords.Parse(args)
	if len(a) > 0 {
		cmd = a[0]
		cmdArgs = a[1:]
	}

	h := &Harriet{
		cmd:  cmd,
		args: cmdArgs,
	}
	if err := h.start(); err != nil {
		panic(err)
	}

	return h
}

func (p *Harriet) start() error {
	log.Printf("run harriet cmd: %s %s", p.cmd, p.args)
	cmd := exec.Command(p.cmd, p.args...)
	cmd.Env = os.Environ()

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}
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
			p.env = append(p.env, exportCmd[1])
			foundExport = true
		}
		if foundExport && t == "" {
			break
		}
	}
	p.c = cmd

	return s.Err()
}

func (p *Harriet) Run(w *prove.Worker) {
	w.Env = append(w.Env, p.env...)
}

func (p *Harriet) Close() error {
	if err := p.c.Process.Signal(os.Interrupt); err != nil {
		return err
	}
	return p.c.Wait()
}
