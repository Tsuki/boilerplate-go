package config

import (
	"runtime"
	"strings"
	"path"
	"github.com/sirupsen/logrus"
	"fmt"
)

var gopath = func() string {
	_, file, _, _ := runtime.Caller(0)
	parts := strings.Split(file, "/src/")
	if len(parts) == 2 {
		return parts[0] + "/src/"
	}
	return ""
}()

type filePos struct {
	Pkg  string `json:"pkg"`
	File string `json:"file"`
	Func string `json:"func"`
	Line int    `json:"line"`
}

type PkgTextFormatter logrus.TextFormatter

func (f *PkgTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	formatter := &(logrus.TextFormatter{})
	stacktrace := getFilePos(5)
	entry.Data["pos"] = fmt.Sprintf("%s.%s(%s:%d)", stacktrace.Pkg, stacktrace.Func, stacktrace.File, stacktrace.Line)
	return formatter.Format(entry)
}

func getFilePos(skip int) filePos {
	pc, fullPath, line, _ := runtime.Caller(skip + 1)
	pkg, file := pkgFile(fullPath)
	return filePos{
		Pkg:  pkg,
		File: file,
		Func: funcName(pc, pkg),
		Line: line,
	}
}

func pkgFile(fullPath string) (pkg, file string) {
	relativePath := strings.TrimPrefix(fullPath, gopath)
	pkg, file = path.Split(relativePath)
	pkg = strings.TrimSuffix(pkg, "/")
	return
}

func funcName(pc uintptr, pkg string) string {
	fn := runtime.FuncForPC(pc)
	s := strings.TrimPrefix(fn.Name(), pkg)
	i := strings.Index(s, ".")
	if i >= 0 {
		s = s[i+1:]
	}
	return s
}
