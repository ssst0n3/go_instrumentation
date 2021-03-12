package stmt

import (
	"github.com/sirupsen/logrus"
	"os"
)

func __stmt() {
	logger := logrus.New()
	file, _ := os.OpenFile("/tmp/instrumentation", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	logger.SetOutput(file)
	logger.Info()
}
