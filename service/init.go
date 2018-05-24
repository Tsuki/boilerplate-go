package service

import (
	. "boilerplate-go/config"
	. "boilerplate-go/service/web"
)

func init() {
	Log.Info("Init Service")
	go StartGinServer()
}
