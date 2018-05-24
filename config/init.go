package config

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var Log *logrus.Logger

func init() {
	configLog()
	Log.Info("Init Config")
}

func configLog() {
	Log = &logrus.Logger{
		Out: os.Stdout,
		Formatter: new(PkgTextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  exPath + "/log/info.log",
		logrus.ErrorLevel: exPath + "/log/error.log",
	}
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
}
