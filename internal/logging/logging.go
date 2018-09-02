package logging

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	stdlog "log"
	"os"
)

//NewLogger provides a set of default fields on a JSON logger
func NewLogger(serviceName string, version string) (logger log.Logger) {
	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = log.With(logger, "svc", serviceName)
	logger = log.With(logger, "version", version)
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	return
}