package prove

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/shogo82148/go-prove/test"
)

type testPlugin int

var pluginResChan chan int

func (p testPlugin) Run(w *Worker, f func()) {
	pluginResChan <- int(p)
	f()
	pluginResChan <- int(p)
}

func Test__run(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "foo.t")
	if err := os.WriteFile(filename, []byte(`print "1..1\nok 1\n";`), 0644); err != nil {
		t.Fatal(err)
	}

	test := &test.Test{
		Path: filename,
		Env:  os.Environ(),
		Exec: "perl",
	}

	pluginResChan = make(chan int, 4)

	p := NewProve()
	w := NewWorker(p, 0)
	p.Plugins = []Plugin{
		testPlugin(1),
		testPlugin(2),
	}

	w.Start()
	p.chanTests <- test
	go func() {
		for range p.chanSuites {
		}
	}()
	close(p.chanTests)
	p.wgWorkers.Wait()
	close(pluginResChan)

	pluginRes := make([]int, 0, 4)
	for res := range pluginResChan {
		pluginRes = append(pluginRes, res)
	}
	if !reflect.DeepEqual(pluginRes, []int{2, 1, 1, 2}) {
		t.Errorf(
			"plugin exec is not valid: got: %v, expect: [2 1 1 2]",
			pluginRes,
		)
	}
}
