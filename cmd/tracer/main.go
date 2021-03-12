package main

import (
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/ssst0n3/go_instrumentation"
	"github.com/ssst0n3/go_instrumentation/filter"
	"github.com/ssst0n3/go_instrumentation/stmt"
	"golang.org/x/xerrors"
	"os"
	"os/exec"
)

func main() {
	trace := stmt.NewTrace()

	args := os.Args[1:]
	log.Logger.Debugf("args: %v", args)

	cmd := go_instrumentation.ParseCommand(args)
	switch cmd {
	case "compile":
		//log.Logger.Infof("args: %v", args)
		newArgs, err := go_instrumentation.Compile(args, trace, filter.BypassTooManyDetails)
		awesome_error.CheckFatal(err)
		//log.Logger.Debug(args)
		log.Logger.Debug(newArgs)
		args = newArgs
	case "link":
		go_instrumentation.Link(args)
	}
	err := go_instrumentation.ForwardCommand(args)
	var exitErr *exec.ExitError
	if err != nil {
		awesome_error.CheckErr(err)
		if xerrors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		} else {
			log.Logger.Fatal(err)
		}
	}
	go_instrumentation.Finish()
	os.Exit(0)
}
