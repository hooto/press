package state

import (
	"../conf"
	"sync"

	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/service/lessids"
)

const (
	ApiServiceUrl           = "http://127.0.0.1:1779/lessfly/v1"  // TODO DNS
	ExtServiceUrl           = "http://127.0.0.1:1779/lessfly/ext" // TODO DNS
	LessIdsOk           int = 1200
	LessIdsUnavailable  int = 1503
	LessIdsUnRegistered int = 1501
)

var (
	ZoneId       = ""
	LessIdsState int
	Locker       sync.Mutex
)

func init() {
}

func Refresh() {

	Locker.Lock()
	defer Locker.Unlock()

	// Check if lessIDS Service Available
	hc := httpclient.Get(lessids.ServiceUrl + "/status/info")
	var rsjson struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			ServiceStatus string `json:"serviceStatus"`
		} `json:"data"`
	}

	err := hc.ReplyJson(&rsjson)
	if err != nil || rsjson.Status != 200 || rsjson.Data.ServiceStatus != "ok" {
		LessIdsState = LessIdsUnavailable
	} else { // Check if lessFly Registered to lessIDS Service

		hc = httpclient.Get(lessids.ServiceUrl +
			"/app-auth/info?instanceid=" +
			conf.Config.InstanceID)
		var rsreg struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
			Data    struct {
				InstanceID string `json:"instance_id"`
				AppId      string `json:"app_id"`
				Version    string `json:"version"`
			} `json:"data"`
		}
		err := hc.ReplyJson(&rsreg)
		if err == nil && rsreg.Status == 404 {
			LessIdsState = LessIdsUnRegistered
		} else {
			LessIdsState = LessIdsOk
		}
	}
}
