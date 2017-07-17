// Copyright 2015~2017 hooto Author, All rights reserved.
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

	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/loscore/losapi"
	"code.hooto.com/lynkdb/iomix/connect"
	"github.com/eryx/hcaptcha/captcha"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/types"

	"code.hooto.com/hooto/hooto-press/api"
)

var (
	Prefix        string
	Config        ConfigCommon
	AppName       = "hooto-press"
	Version       = "0.2.0.dev"
	CaptchaConfig = captcha.DefaultConfig

	pod_inst_updated time.Time
	pod_inst         = "/home/action/.los/pod_instance.json"

	User = &user.User{
		Uid:      "2048",
		Gid:      "2048",
		Username: "action",
		HomeDir:  "/home/action",
	}

	SysConfigList = api.SysConfigList{}
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
	Database              base.Config              `json:"database"`
	IoConnectors          connect.MultiConnOptions `json:"io_connectors"`
	RunMode               string                   `json:"run_mode,omitempty"`
}

func init() {

	SysConfigList.Kind = "SysConfigList"

	SysConfigList.Insert(api.SysConfig{
		"frontend_header_site_name", "CMS",
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

	SysConfigList.Insert(api.SysConfig{
		"frontend_html_head_subtitle", "CMS",
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
		"ls2_uri", "//127.0.0.1/s2/bucket-name",
		"Storage Service URI", "",
	})
}

func Initialize(prefix string) error {

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

	if Config.RunMode != "local-dev" {

		var inst losapi.Pod
		if err := json.DecodeFile(pod_inst, &inst); err != nil {
			return err
		}

		if inst.Spec == nil ||
			len(inst.Spec.Boxes) == 0 ||
			inst.Spec.Boxes[0].Resources == nil {
			return errors.New("Not Pod Instance Setup")
		}

		var (
			opt    *losapi.AppOption
			optref *losapi.AppOption
		)

		for _, app := range inst.Apps {

			if app.Spec.Meta.Name != "hooto-press" {
				continue
			}

			//
			if opt == nil {
				opt = app.Operate.Options.Get("cfg/hooto-press") // TODO
			}

			if optref == nil {
				optref = app.Operate.Options.Get("cfg/los-mysql") // TODO
			}
		}

		if opt == nil {
			return errors.New("No Configure Found")
		}

		if v, ok := opt.Items.Get("iam_service_url"); ok {
			Config.IamServiceUrl = v.String()
		} else {
			return errors.New("No Config.IamServiceUrl Found")
		}

		if v, ok := opt.Items.Get("iam_service_url_frontend"); ok {
			Config.IamServiceUrlFrontend = v.String()
		} else {
			Config.IamServiceUrlFrontend = Config.IamServiceUrl
		}

		if v, ok := opt.Items.Get("http_pprof_enable"); ok && v.String() == "1" {
			Config.HttpPortPprof = Config.HttpPort + 1
		} else {
			Config.HttpPortPprof = 0
		}

		if optref == nil || optref.Ref == nil {
			return errors.New("No Database Connection Configure Found")
		}

		Config.Database.Driver = "mysql"

		if v, ok := optref.Items.Get("db_name"); ok {
			Config.Database.Dbname = v.String()
		}

		if v, ok := optref.Items.Get("db_user"); ok {
			Config.Database.User = v.String()
		}

		if v, ok := optref.Items.Get("db_auth"); ok {
			Config.Database.Pass = v.String()
		}

		if optref.Ref.PodId == "" {
			return errors.New("No UpStream Pod Found")
		}

		var nsz losapi.NsPodServiceMap
		if err := json.DecodeFile("/dev/shm/los/nsz/"+optref.Ref.PodId, &nsz); err != nil {
			return err
		}

		// TODO
		if srv := nsz.Get(3306); srv == nil || len(srv.Items) == 0 {
			return errors.New("No Pod ServicePort Found")
		} else {
			Config.Database.Host = srv.Items[0].Ip
			Config.Database.Port = fmt.Sprintf("%d", srv.Items[0].Port)
		}
	}

	Save()

	dcn, err := Config.DatabaseInstance()
	if err != nil {
		return err
	}

	{
		rs, err := dcn.Base.Query(base.NewQuerySet().From("sys_config").Limit(1000))
		if err != nil {
			return err
		}

		for _, v := range rs {

			SysConfigList.Insert(api.SysConfig{
				Key:   v.Field("key").String(),
				Value: v.Field("value").String(),
			})
		}
	}

	if Config.AppTitle == "" {
		Config.AppTitle = "Hooto Press"
	}

	// Setting CAPTCHA
	CaptchaConfig.DataDir = Prefix + "/var/captchadb"

	if err := captcha.Config(CaptchaConfig); err != nil {
		return err
	}

	// Default User
	if User, err = user.Current(); err != nil {
		return err
	}

	if err := init_data(); err != nil {
		return err
	}

	return module_init()
}

//
func init_data() error {

	io_name := types.NewNameIdentifier("htp_local_cache")
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
	}

	Config.IoConnectors.SetOptions(*opts)

	return nil
}

func Save() error {
	return json.EncodeToFile(Config, Prefix+"/etc/main.json", "  ")
}
