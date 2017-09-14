package plugin

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/shogo82148/go-prove"
)

func TestHarriet(t *testing.T) {
	_, err := exec.LookPath("harriet")
	if err != nil {
		t.Log("harriet is not installed. skip this test.")
		return
	}

	w := prove.NewWorker(nil, 0)
	h := harrietLoader("harriet", "harriet _harriet_test")

	testfile := ""
	h.Run(w, func() {
		for _, e := range w.Env {
			if strings.HasPrefix(e, "GO_PROVE_TEST_HARRIET=") {
				testfile = e[len("GO_PROVE_TEST_HARRIET="):]
			}
		}

		// testfile should exist during running tests
		if testfile == "" {
			t.Error("environment value is not set")
			return
		}
		_, err := os.Stat(testfile)
		if err != nil {
			t.Error("unexpected error: ", err)
		}
	})
	h.(io.Closer).Close()

	// testfile should be removed because harriet command has already finished.
	if testfile == "" {
		t.Error("environment value is not set")
		return
	}
	_, err = os.Stat(testfile)
	if !os.IsNotExist(err) {
		t.Error("testfile still exists")
	}
}
