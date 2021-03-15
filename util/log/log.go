package log

import (
	"github.com/beego/beego/v2/core/logs"
)

var (
	TGLog  = logs.NewLogger(10000)
	OPQLog = logs.NewLogger(10000)
)

func LogInit() {

	err := TGLog.SetLogger("file", `{"filename":"./log/tg.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	if err != nil {
		panic(err)
	}

	err = TGLog.SetLogger("console")
	if err != nil {
		panic(err)
	}
	// TGLog.EnableFuncCallDepth(true)
	TGLog.Async()

	err = OPQLog.SetLogger("file", `{"filename":"./log/opq.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	if err != nil {
		panic(err)
	}

	err = OPQLog.SetLogger("console")
	if err != nil {
		panic(err)
	}

	//TGLog.EnableFuncCallDepth(true)
	OPQLog.Async()
}
