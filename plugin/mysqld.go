package plugin

import (
	"fmt"

	"github.com/lestrrat/go-test-mysqld"
	"github.com/shogo82148/go-prove"
)

type TestMysqld struct{}

func (p *TestMysqld) Run(w *prove.Worker, f func()) {
	mysqld, err := mysqltest.NewMysqld(nil)
	if err != nil {
		fmt.Errorf("mysql error: %s\n", err)
	}
	defer mysqld.Stop()

	address := mysqld.ConnectString(0)
	w.Env = append(w.Env, fmt.Sprintf("GO_PROVE_MYSQLD=%s", address))

	f()
}
