#!/bin/bash
set -x
go clean -cache
go build -work -p 1 -a -toolexec ~/pentest_project/go_instrumentation/cmd/tracer/tracer