package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/lestrrat/go-test-mysqld"
	"github.com/shogo82148/go-tap"
)

var exitCode = 0
var mutex *sync.Mutex

type Job struct {
	path string
	env  []string
	exec string
}

// JUnitTestSuites is a collection of JUnit test suites.
type JUnitTestSuites struct {
	XMLName xml.Name `xml:"testsuites"`
	Suites  []JUnitTestSuite
}

// JUnitTestSuite is a single JUnit test suite which may contain many
// testcases.
type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Time       string          `xml:"time,attr"`
	Name       string          `xml:"name,attr"`
	Properties []JUnitProperty `xml:"properties>property,omitempty"`
	TestCases  []JUnitTestCase
}

// JUnitTestCase is a single test case with its result.
type JUnitTestCase struct {
	XMLName     xml.Name          `xml:"testcase"`
	Classname   string            `xml:"classname,attr"`
	Name        string            `xml:"name,attr"`
	Time        string            `xml:"time,attr"`
	SkipMessage *JUnitSkipMessage `xml:"skipped,omitempty"`
	Failure     *JUnitFailure     `xml:"failure,omitempty"`
}

// JUnitSkipMessage contains the reason why a testcase was skipped.
type JUnitSkipMessage struct {
	Message string `xml:"message,attr"`
}

// JUnitProperty represents a key/value pair used to define properties.
type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// JUnitFailure contains data related to a failed test.
type JUnitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

func main() {
	mutex = &sync.Mutex{}
	proveMain()
	os.Exit(exitCode)
}

func proveMain() {
	var jobs int
	var execParam string
	flag.IntVar(&jobs, "j", 1, "Run N test jobs in parallel")
	flag.IntVar(&jobs, "jobs", 1, "Run N test jobs in parallel")
	flag.StringVar(&execParam, "exec", "perl", "")
	flag.Parse()

	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	files := findTestFiles()

	wg := &sync.WaitGroup{}
	wg2 := &sync.WaitGroup{}
	chanPath := make(chan string, len(files))
	chanSuite := make(chan JUnitTestSuite)
	chanDone := make(chan bool)
	for i := 0; i < jobs; i++ {
		wg2.Add(1)
		go func() {
			mysqld, _ := mysqltest.NewMysqld(nil)
			defer mysqld.Stop()
			defer wg2.Done()
			for {
				address := mysqld.ConnectString(0)

				select {
				case path := <-chanPath:
					log.Println("start: " + path)
					job := &Job{
						path: path,
						env: []string{
							fmt.Sprintf("GO_PROVE_MYSQLD=%s", address),
						},
						exec: execParam,
					}
					chanSuite <- job.run()
					wg.Done()
					log.Println("done: " + path)
				case <-chanDone:
					return
				}
			}
		}()
	}

	for _, path := range files {
		wg.Add(1)
		chanPath <- path
	}

	suites := JUnitTestSuites{}
	for range files {
		suites.Suites = append(suites.Suites, <-chanSuite)
	}

	bytes, _ := xml.MarshalIndent(suites, "", "\t")
	os.Stdout.Write(bytes)

	wg.Wait()

	for i := 0; i < jobs; i++ {
		chanDone <- true
	}
	wg2.Wait()
}

func findTestFiles() []string {
	files := []string{}
	filepath.Walk(
		"t",
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(path, ".t") {
				files = append(files, path)
			}

			return nil
		})
	return files
}

func (j *Job) run() JUnitTestSuite {
	execParam := strings.Split(j.exec, " ")
	execParam = append(execParam, j.path)
	cmd := exec.Command(execParam[0], execParam[1:]...)
	cmd.Env = j.env
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		log.Println(err)
	}
	go io.Copy(os.Stderr, stderr)

	suite := j.newJUnitTestSuite(stdout)

	cmd.Wait()

	return suite
}

func (j *Job) newJUnitTestSuite(reader io.Reader) JUnitTestSuite {
	className := strings.Replace(j.path, "/", "_", -1)
	className = strings.Replace(className, ".", "_", -1)

	start := time.Now()
	lastTestEnd := start

	parser, _ := tap.NewParser(reader)

	ts := JUnitTestSuite{
		Tests:      0,
		Failures:   0,
		Time:       "0.000",
		Name:       className,
		Properties: []JUnitProperty{},
		TestCases:  []JUnitTestCase{},
	}

	for {
		line, err := parser.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}
		testEnd := time.Now()

		testCase := JUnitTestCase{
			Classname: className,
			Name:      line.Description,
			Time:      fmt.Sprintf("%.3f", (testEnd.Sub(lastTestEnd)).Seconds()),
			Failure:   nil,
		}
		if !line.Ok {
			ts.Failures++
			testCase.Failure = &JUnitFailure{
				Message:  "not ok",
				Type:     "",
				Contents: line.Diagnostic,
			}
		}
		ts.Tests++
		ts.TestCases = append(ts.TestCases, testCase)
		lastTestEnd = testEnd
	}

	suite, _ := parser.Suite()
	if !suite.Ok {
		markAsFail()
	}

	end := time.Now()
	ts.Time = fmt.Sprintf("%.3f", (end.Sub(start)).Seconds())
	return ts
}

func markAsFail() {
	mutex.Lock()
	defer mutex.Unlock()
	exitCode = 1
}
