package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	// "time"

	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/logger"
	"github.com/lessos/lessgo/service/lessids"

	"./conf"
	"./datax"
	// "./state"
)

import (
	ext_comment "../spec/comment/websrv"
	cdef "./websrv/frontend"
	cmgr "./websrv/mgr"
	capi "./websrv/v1"
)

var (
	flagPrefix     = flag.String("prefix", "", "the prefix folder path")
	flagCpuprofile = flag.String("pprof", "", "the pprof path")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// render functions
	httpsrv.GlobalService.Config.TemplateFuncRegister("TimeFormat", datax.TimeFormat)
	httpsrv.GlobalService.Config.TemplateFuncRegister("FieldDebug", datax.FieldDebug)
	httpsrv.GlobalService.Config.TemplateFuncRegister("FieldString", datax.FieldString)
	httpsrv.GlobalService.Config.TemplateFuncRegister("FieldSubString", datax.FieldSubString)
	httpsrv.GlobalService.Config.TemplateFuncRegister("FieldHtml", datax.FieldHtml)
	httpsrv.GlobalService.Config.TemplateFuncRegister("FieldSubHtml", datax.FieldSubHtml)
	httpsrv.GlobalService.Config.TemplateFuncRegister("pagelet", datax.Pagelet)
	httpsrv.GlobalService.Config.TemplateFuncRegister("FilterUri", datax.FilterUri)
}

func main() {

	//
	flag.Parse()
	if err := conf.Initialize(*flagPrefix); err != nil {
		fmt.Println("Error on conf.Initialize", err)
		logger.Printf("error", "conf.Initialize error: %v", err)
		os.Exit(1)
	}

	lessids.ServiceUrl = conf.Config.LessIdsUrl

	// httpsrv.Config.UrlBasePath = "cmf"
	httpsrv.GlobalService.Config.HttpPort = conf.Config.HttpPort
	// httpsrv.Config.LessIdsServiceUrl = conf.Config.LessIdsUrl

	// state
	// for {

	// 	state.Refresh()

	// 	if state.LessIdsState == state.LessIdsUnRegistered ||
	// 		state.LessIdsState == state.LessIdsOk {
	// 		break
	// 	}

	// 	time.Sleep(3e9)
	// }

	// conf.SpecRefresh("c8f0ltxp")

	//
	// httpsrv.Config.I18n(conf.Config.Prefix + "/src/i18n/en.json")
	// httpsrv.Config.I18n(conf.Config.Prefix + "/src/i18n/zh_CN.json")

	httpsrv.GlobalService.ModuleRegister("/+/comment", ext_comment.NewModule())

	//
	httpsrv.GlobalService.ModuleRegister("/v1", capi.NewModule())
	httpsrv.GlobalService.ModuleRegister("/mgr", cmgr.NewModule())
	httpsrv.GlobalService.ModuleRegister("/", cdef.NewModule())

	//
	fmt.Println("Running")
	httpsrv.GlobalService.Start()
}
