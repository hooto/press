package main

import (
	"./conf"
	"./datax"
	"flag"
	"fmt"
	"github.com/lessos/lessgo/logger"
	"github.com/lessos/lessgo/pagelet"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
)

import (
	capi "./apiserver/v1"
	cdef "./controllers"
	cmgr "./mgr/controllers"
)

var (
	flagPrefix     = flag.String("prefix", "", "the prefix folder path")
	flagCpuprofile = flag.String("pprof", "", "the pprof path")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	go http.ListenAndServe("localhost:6060", nil)

	if *flagCpuprofile != "" {
		f, err := os.Create(*flagCpuprofile)
		if err != nil {
			fmt.Println(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	//
	flag.Parse()
	if err := conf.Initialize(*flagPrefix); err != nil {
		fmt.Println("Error on conf.Initialize", err)
		logger.Printf("error", "conf.Initialize error: %v", err)
		os.Exit(1)
	}

	// pagelet.Config.UrlBasePath = "cmf"
	pagelet.Config.HttpPort = conf.Config.HttpPort
	pagelet.Config.LessIdsServiceUrl = conf.Config.LessIdsUrl

	// pagelet.Config.ViewFuncRegistry("Query", datax.Query)
	// pagelet.Config.ViewFuncRegistry("NewQuery", datax.NewQuery)
	// pagelet.Config.ViewFuncRegistry("QueryEntry", datax.QueryEntry)
	pagelet.Config.ViewFuncRegistry("Field", datax.Field)
	pagelet.Config.ViewFuncRegistry("pagelet", datax.Pagelet)

	//
	pagelet.Config.I18n(conf.Config.Prefix + "/src/i18n/en.json")
	pagelet.Config.I18n(conf.Config.Prefix + "/src/i18n/zh_CN.json")

	//
	pagelet.Config.RouteAppend("v1", "/:controller/:action")
	pagelet.RegisterController("v1", (*capi.Node)(nil))
	pagelet.RegisterController("v1", (*capi.Term)(nil))
	pagelet.RegisterController("v1", (*capi.Spec)(nil))
	pagelet.RegisterController("v1", (*capi.NodeModel)(nil))
	pagelet.RegisterController("v1", (*capi.TermModel)(nil))

	//
	pagelet.Config.RouteStaticAppend("mgr", "/~", conf.Config.Prefix+"/static")
	pagelet.Config.RouteStaticAppend("mgr", "/-", conf.Config.Prefix+"/src/mgr/tpls")
	pagelet.Config.RouteAppend("mgr", "/:controller/:action")
	pagelet.RegisterController("mgr", (*cmgr.Index)(nil))

	//
	pagelet.Config.RouteStaticAppend("default", "/~", conf.Config.Prefix+"/static")
	pagelet.Config.ViewPath("default", conf.Config.Prefix+"/src/views")
	pagelet.Config.RouteAppend("default", "/:controller/:action")
	pagelet.RegisterController("default", (*cdef.Index)(nil))
	pagelet.RegisterController("default", (*cdef.Error)(nil))

	//
	fmt.Println("Running")
	pagelet.Run()
}
