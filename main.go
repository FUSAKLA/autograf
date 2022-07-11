package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
)

var (
	Version = "development"
)

type Context struct {
	logger logrus.FieldLogger
}

type VersionCmd struct{}

func (r *VersionCmd) Run(_ *Context) error {
	fmt.Println(Version)
	return nil
}

var CLI struct {
	Debug    bool        `help:"Enable debug logging."`
	Version  VersionCmd  `cmd:"" help:"Print version."`
	Generate GenerateCmd `cmd:"" default:"1" help:"Generates Grafana dashboard based on a given Prometheus metrics and prints it to stdout if not specified otherwise."`
}

func main() {
	ctx := kong.Parse(&CLI)
	rootLogger := logrus.New()
	rootLogger.SetOutput(os.Stderr)
	rootLogger.SetLevel(logrus.WarnLevel)
	if CLI.Debug {
		rootLogger.SetLevel(logrus.DebugLevel)
	}

	err := ctx.Run(&Context{
		logger: rootLogger.WithField("command", ctx.Command()),
	})
	ctx.FatalIfErrorf(err)
}
