// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/hooto/hcaptcha/captcha4g"
	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/iam/iamapi"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/iomix/connect"
	"github.com/lynkdb/iomix/rdb"
	"github.com/lynkdb/iomix/rdb/modeler"
	"github.com/sysinner/incore/inapi"

	"github.com/hooto/hpress/api"
	"github.com/hooto/hpress/store"
)

var (
	Prefix         string
	Config         ConfigCommon
	AppName        = "hooto-press"
	Version        = "0.2"
	Release        = "1"
	SysVersionSign = ""
	CaptchaConfig  = captcha4g.DefaultConfig

	pod_inst_updated time.Time
	pod_inst         = "/home/action/.sysinner/pod_instance.json"

	User = &user.User{
		Uid:      "2048",
		Gid:      "2048",
		Username: "action",
		HomeDir:  "/home/action",
	}

	SysConfigList = api.SysConfigList{}
	inited        = false
)

type ConfigCommon struct {
	UrlBasePath           string                   `json:"url_base_path,omitempty"`
	ModuleDir             string                   `json:"module_dir,omitempty"`
	InstanceID            string                   `json:"instance_id"`
	AppInstance           iamapi.AppInstance       `json:"app_instance"`
	AppTitle              string                   `json:"app_title,omitempty"`
	HttpPort              uint16                   `json:"http_port"`
	HttpPortPprof         uint16                   `json:"http_port_pprof,omitempty"`
	IamServiceUrl         string                   `json:"iam_service_url"`
	IamServiceUrlFrontend string                   `json:"iam_service_url_frontend"`
	IoConnectors          connect.MultiConnOptions `json:"io_connectors"`
	RunMode               string                   `json:"run_mode,omitempty"`
}

func init() {

	SysConfigList.Kind = "SysConfigList"

	//
	SysConfigList.Insert(api.SysConfig{
		"frontend_header_site_name", "Site Name",
		"Site's Name", "",
	})
	SysConfigList.Insert(api.SysConfig{
		"frontend_header_site_logo_url", "",
		"", "",
	})
	SysConfigList.Insert(api.SysConfig{
		"frontend_footer_copyright", "2015~2017 hooto.com",
		"", "",
	})

	//
	SysConfigList.Insert(api.SysConfig{
		"frontend_html_head_subtitle", "HP",
		"Sub Title for HTML Head Title", "",
	})
	SysConfigList.Insert(api.SysConfig{
		"frontend_html_head_meta_keywords", "",
		"Meta Keywords in HTML Head for Search engine optimization", "",
	})
	SysConfigList.Insert(api.SysConfig{
		"frontend_html_head_meta_description", "",
		"Meta Description in HTML Head for Search engine optimization", "",
	})

	SysConfigList.Insert(api.SysConfig{
		"frontend_footer_analytics_scripts", "",
		"Embeded analytics scripts, ex. Google Analytics or Piwik ...", "text",
	})
	SysConfigList.Insert(api.SysConfig{
		"ls2_uri", "//127.0.0.1:9533/hpress/s2",
		"Storage Service URI", "",
	})

	//
	SysConfigList.Insert(api.SysConfig{
		"http_h_ac_allow_origin", "",
		"HTTP Access-Control-Allow-Origin", "",
	})

	go func() {
		for {
			time.Sleep(60e9)
			if err := sync_sysinner_config(); err != nil {
				hlog.Printf("error", "sync_sysinner_config err: %s", err.Error())
			}
		}
	}()
}

func Initialize(prefix string) error {

	if inited {
		return nil
	}

	var err error

	if prefix == "" {
		if prefix, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/.."); err != nil {
			prefix = "/home/action/apps/hooto-press"
		}
	}

	Prefix = filepath.Clean(prefix)
	Config.ModuleDir = Prefix + "/modules"

	file := Prefix + "/etc/main.json"
	if err := json.DecodeFile(file, &Config); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	{
		if Config.HttpPort == 0 {
			Config.HttpPort = 9533
		}
	}

	if Config.InstanceID != "" {
		SysVersionSign = idhash.HashToHexString([]byte(fmt.Sprintf("%s-%s-%s", Version, Release, Config.InstanceID)), 16)
	} else {
		SysVersionSign = "unreg"
	}

	if Config.AppTitle == "" {
		Config.AppTitle = "Hooto Press"
	}

	if err := sync_sysinner_config(); err != nil {
		return err
	}

	// Setting CAPTCHA
	CaptchaConfig.DataDir = Prefix + "/var/hcaptchadb"
	if err := captcha4g.Config(CaptchaConfig); err != nil {
		return err
	}

	// Default User
	if User, err = user.Current(); err != nil {
		return err
	}

	if err := store_init(); err != nil {
		return err
	}

	//
	{
		rs, err := store.Data.Query(rdb.NewQuerySet().From("sys_config").Limit(1000))
		if err != nil {
			hlog.Print("error", err.Error())
			return err
		}

		for _, v := range rs {

			SysConfigList.Insert(api.SysConfig{
				Key:   v.Field("key").String(),
				Value: v.Field("value").String(),
			})
		}
	}

	if err := module_init(); err != nil {
		return err
	}

	inited = true

	return nil
}

func sync_sysinner_config() error {

	if Config.RunMode == "local-dev" {
		return nil
	}

	var inst inapi.Pod
	if err := json.DecodeFile(pod_inst, &inst); err != nil {
		return err
	}

	if inst.Spec == nil ||
		len(inst.Spec.Boxes) == 0 ||
		inst.Spec.Boxes[0].Resources == nil {
		return errors.New("Not Pod Instance Setup")
	}

	var (
		opt       *inapi.AppOption
		optref    *inapi.AppOption
		data_opts connect.ConnOptions
		sync      = false
	)

	for _, app := range inst.Apps {

		if app.Spec.Meta.Name != "hooto-press" &&
			app.Spec.Meta.Name != "hooto-press-x1" {
			continue
		}

		//
		if opt == nil {
			opt = app.Operate.Options.Get("cfg/hooto-press") // TODO
		}

		if optref == nil {
			optref = app.Operate.Options.Get("cfg/sysinner-mysql") // TODO
		}
	}

	if opt == nil {
		return errors.New("No Configure Found")
	}

	if v, ok := opt.Items.Get("iam_service_url"); ok {
		if v.String() != Config.IamServiceUrl {
			Config.IamServiceUrl, sync = v.String(), true
		}
	} else {
		return errors.New("No Config.IamServiceUrl Found")
	}

	if v, ok := opt.Items.Get("iam_service_url_frontend"); ok {
		if v.String() != Config.IamServiceUrlFrontend {
			Config.IamServiceUrlFrontend, sync = v.String(), true
		}
	} else {
		Config.IamServiceUrlFrontend = Config.IamServiceUrl
	}

	if v, ok := opt.Items.Get("http_pprof_enable"); ok && v.String() == "1" {
		if p := Config.HttpPort + 1; p != Config.HttpPortPprof {
			Config.HttpPortPprof, sync = Config.HttpPort+1, true
		}
	} else {
		Config.HttpPortPprof = 0
	}

	if optref == nil {
		return errors.New("No Database Connection Configure Found")
	}

	ref_pod_id := inst.Meta.ID
	if optref.Ref != nil && optref.Ref.PodId != "" {
		ref_pod_id = optref.Ref.PodId
	}

	data_opts.Name = types.NameIdentifier("hpress_database")
	data_opts.Driver = "lynkdb/mysqlgo"

	if v, ok := optref.Items.Get("db_name"); ok {
		if data_opts.Value("dbname") != v.String() {
			data_opts.SetValue("dbname", v.String())
			sync = true
		}
	}

	if v, ok := optref.Items.Get("db_user"); ok {
		if data_opts.Value("user") != v.String() {
			data_opts.SetValue("user", v.String())
			sync = true
		}
	}

	if v, ok := optref.Items.Get("db_auth"); ok {
		if data_opts.Value("pass") != v.String() {
			data_opts.SetValue("pass", v.String())
			sync = true
		}
	}

	var nsz inapi.NsPodServiceMap
	if err := json.DecodeFile("/dev/shm/sysinner/nsz/"+ref_pod_id, &nsz); err != nil {
		return err
	}

	// TODO
	if srv := nsz.Get(3306); srv == nil || len(srv.Items) == 0 {
		return errors.New("No Pod ServicePort Found")
	} else {
		if data_opts.Value("host") != srv.Items[0].Ip {
			data_opts.SetValue("host", srv.Items[0].Ip)
			sync = true
		}
		if p := fmt.Sprintf("%d", srv.Items[0].Port); p != data_opts.Value("port") {
			data_opts.SetValue("port", p)
			sync = true
		}
	}

	Config.IoConnectors.SetOptions(data_opts)

	if sync {
		Save()
		hlog.Print("warn", "sysinner configs synced")
	}

	return nil
}

//
func store_init() error {

	{
		io_name := types.NewNameIdentifier("hpress_local_cache")
		opts := Config.IoConnectors.Options(io_name)

		if opts == nil {
			opts = &connect.ConnOptions{
				Name:      io_name,
				Connector: "iomix/skv/Connector",
				Driver:    types.NewNameIdentifier("lynkdb/kvgo"),
			}
		}

		if opts.Value("data_dir") == "" {
			opts.SetValue("data_dir", Prefix+"/var/"+string(io_name))
			Save()
		}

		Config.IoConnectors.SetOptions(*opts)
	}

	var dbname = "dbaction"
	{
		io_name := types.NewNameIdentifier("hpress_database")
		opts := Config.IoConnectors.Options(io_name)

		if opts == nil {
			opts = &connect.ConnOptions{
				Name:      io_name,
				Connector: "iomix/rdb/Connector",
				Driver:    types.NewNameIdentifier("lynkdb/mysqlgo"),
			}
		}

		if opts.Value("host") == "" {
			opts.SetValue("host", "localhost")
		}

		if opts.Value("port") == "" {
			opts.SetValue("port", "3306")
		}

		dbname = opts.Value("dbname")

		Config.IoConnectors.SetOptions(*opts)
	}

	if err := store.Init(Config.IoConnectors); err != nil {
		hlog.Printf("error", "store_init %s", err.Error())
		return err
	}

	ds, err := modeler.NewDatabaseEntryFromJson(dsBase)
	if err != nil {
		return err
	}

	mor, err := store.Data.Modeler()
	if err != nil {
		hlog.Printf("error", "store_init %s", err.Error())
		return err
	}

	err = mor.Sync(dbname, ds)
	if err != nil {
		hlog.Printf("error", "store_init %s", err.Error())
	}

	return err
}

func Save() error {
	return json.EncodeToFile(Config, Prefix+"/etc/main.json", "  ")
}
