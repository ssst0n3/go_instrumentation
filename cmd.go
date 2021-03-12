package go_instrumentation

import (
	"fmt"
	"github.com/ssst0n3/awesome_libs/log"
	"os"
	"os/exec"
	"strings"
)

func ParseCommand(args []string) (cmd string) {
	binary := args[0]
	cmd = binary[strings.LastIndex(binary, "/")+1:]
	return
}

func ForwardCommand(args []string) error {
	path := args[0]
	args = args[1:]
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	quotedArgs := fmt.Sprintf("%+q", args)
	log.Logger.Debugf("forwarding command `%s %s`", path, quotedArgs[1:len(quotedArgs)-1])
	return cmd.Run()
}
