package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/mariusvanderwijden/threadpool"
	"gopkg.in/urfave/cli.v1"
)

var (
	app *cli.App
	// Flags
	testFlag = cli.StringFlag{
		Name:  "test",
		Usage: "Test name to execute",
	}
	threadFlag = cli.IntFlag{
		Name:  "threads",
		Usage: "Number of threads to use",
		Value: runtime.NumCPU(),
	}
	countFlag = cli.IntFlag{
		Name:  "count",
		Usage: "Number of executions per call",
		Value: 1,
	}
)

func init() {
	app = cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = "Marius van der Wijden"
	app.Email = "m.vanderwijden@live.de"
	app.Version = "0.1"
	app.Usage = "tool to repro flaky tests"
	app.Flags = []cli.Flag{
		testFlag,
		threadFlag,
		countFlag,
	}
	app.Action = repro
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func repro(c *cli.Context) error {
	var (
		pool     = threadpool.NewThreadPool(c.GlobalInt(threadFlag.Name))
		testName = c.GlobalString(testFlag.Name)
		count    = c.GlobalInt(countFlag.Name)
	)
	if count < 1 {
		count = 1
	}
	for {
		pool.Get(1)
		go func() {
			defer pool.Put(1)
			if err := execute(testName, count); err != nil {
				fmt.Printf("ERROR occurred: %v\n", err)
				os.Exit(1)
			}
		}()
	}
}

func execute(test string, count int) error {
	var (
		cnt  = strconv.Itoa(count)
		args []string
	)
	if test != "" {
		args = []string{
			"test",
			"--run",
			test,
			"-count",
			cnt,
			"./..."}
	} else {
		args = []string{
			"test",
			"-count",
			cnt,
			"./...",
		}
	}
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}
