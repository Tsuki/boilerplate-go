package config

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

//var Log *logrus.Entry
var Log *logrus.Logger
var Logger *logrus.Logger

func init() {
	configLog()
	Log.Info("Init Config")
}

func configLog() {
	log := &logrus.Logger{
		Out: os.Stdout,
		Formatter: new(PkgJSONFormatter),
		//Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  exPath + "/log/info.log",
		logrus.ErrorLevel: exPath + "/log/error.log",
	}
	log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
	Logger = log
	Log = log
}
