package go_instrumentation

import (
	"github.com/sirupsen/logrus"
	"github.com/ssst0n3/awesome_libs/log"
	"os"
)

var LogWriter *os.File

func init() {
	LogWriter, err := os.OpenFile("/tmp/out", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		os.Exit(-1)
	}
	log.Logger.Level = logrus.DebugLevel
	log.Logger.SetOutput(LogWriter)
	log.Logger.Infof("======START====")
}

func Finish() {
	defer LogWriter.Close()
	log.Logger.Infof("=====FINISHED====\n\n")
}