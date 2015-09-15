package plugin

import (
	"fmt"
	"log"

	"github.com/lestrrat/go-test-mysqld"
	"github.com/shogo82148/go-prove"
)

type TestMysqld struct{}

func (p *TestMysqld) Run(w *prove.Worker, f func()) {
	log.Printf("run mysqld")
	mysqld, err := mysqltest.NewMysqld(nil)
	if err != nil {
		log.Printf("mysql error: %s\n", err)
	}
	defer mysqld.Stop()

	address := mysqld.ConnectString(0)
	w.Env = append(w.Env, fmt.Sprintf("GO_PROVE_MYSQLD=%s", address))

	f()
}
