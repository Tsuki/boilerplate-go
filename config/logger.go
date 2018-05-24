package config

import (
	"runtime"
	"strings"
	"path"
	"github.com/sirupsen/logrus"
	"io"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"fmt"
	"sync"
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

type PkgJSONFormatter struct {
	Level      string `json:"level"`
	Msg        string `json:"msg"`
	Time       string `json:"time"`
	Package    string `json:"pkg"`
	File       string `json:"file"`
	Func       string `json:"func"`
	Line       int    `json:"line"`
	Pos        string `json:"pos"`
	isTerminal bool
	sync.Once
}

func (f *PkgJSONFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}
}
func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

func (f *PkgJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	f.Do(func() { f.init(entry) })
	//var bytes []byte
	if f.isTerminal {
		formatter := &(logrus.TextFormatter{})
		stacktrace := getFilePos(5)
		entry.Data["pos"] = fmt.Sprintf("%s.%s(%s:%d)", stacktrace.Pkg, stacktrace.Func, stacktrace.File, stacktrace.Line)
		return formatter.Format(entry)
	} else {
		stacktrace := getFilePos(5)
		jsonFormatter := &(logrus.JSONFormatter{})
		entry.Data["pos"] = fmt.Sprintf("%s.%s(%s:%d)", stacktrace.Pkg, stacktrace.Func, stacktrace.File, stacktrace.Line)
		return jsonFormatter.Format(entry)
		//customFormatter := PkgJSONFormatter{}
		//json.Unmarshal(bytes, &customFormatter)
		//customFormatter.File = stacktrace.File
		//customFormatter.Line = stacktrace.Line
		//customFormatter.Package = stacktrace.Pkg
		//customFormatter.Func = stacktrace.Func
		//customFormatter.Pos = fmt.Sprintf("%s.%s(%s:%d)", stacktrace.Pkg, stacktrace.Func, stacktrace.File, stacktrace.Line)
		//serialized, err := json.Marshal(customFormatter)
		//if err != nil {
		//	return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
		//}
		//return bytes, nil
	}
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
