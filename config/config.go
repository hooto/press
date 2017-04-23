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
	"os"
	"os/user"
	"path/filepath"

	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lynkdb/iomix/connect"
	"github.com/eryx/hcaptcha/captcha"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/types"

	"code.hooto.com/hooto/hootopress/api"
)

var (
	Prefix        string
	Config        ConfigCommon
	AppName       = "hooto-hootopress"
	Version       = "0.1.7.dev"
	CaptchaConfig = captcha.DefaultConfig

	User = &user.User{
		Uid:      "2048",
		Gid:      "2048",
		Username: "action",
		HomeDir:  "/home/action",
	}

	SysConfigList = api.SysConfigList{}
)

type ConfigCommon struct {
	UrlBasePath   string                   `json:"url_base_path,omitempty"`
	ModuleDir     string                   `json:"module_dir,omitempty"`
	InstanceID    string                   `json:"instance_id"`
	AppInstance   iamapi.AppInstance       `json:"app_instance"`
	AppTitle      string                   `json:"app_title,omitempty"`
	HttpPort      uint16                   `json:"http_port"`
	IamServiceUrl string                   `json:"iam_service_url"`
	Database      base.Config              `json:"database"`
	IoConnectors  connect.MultiConnOptions `json:"io_connectors"`
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
		"ls2_uri", "//127.0.0.1/bucket-name",
		"Storage Service URI", "",
	})
}

func Initialize(prefix string) error {

	var err error

	if prefix == "" {
		if prefix, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/.."); err != nil {
			prefix = "/opt/hooto/hootopress"
		}
	}

	Prefix = filepath.Clean(prefix)
	Config.ModuleDir = Prefix + "/modules"

	file := Prefix + "/etc/main.json"
	if err := json.DecodeFile(file, &Config); err != nil {
		return err
	}

	//
	if Config.IamServiceUrl == "" {
		return errors.New("Error: `iam_service_url` can not be null")
	}

	dcn, err := Config.DatabaseInstance()
	if err != nil {
		return err
	}

	rs, err := dcn.Base.Query(base.NewQuerySet().From("sys_config").Limit(1000))
	if err == nil {

		for _, v := range rs {

			SysConfigList.Insert(api.SysConfig{
				Key:   v.Field("key").String(),
				Value: v.Field("value").String(),
			})
		}
	}

	if Config.AppTitle == "" {
		Config.AppTitle = "HootoPress"
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

	return json.EncodeToFile(ConfigCommon{
		InstanceID:    Config.InstanceID,
		AppInstance:   Config.AppInstance,
		AppTitle:      Config.AppTitle,
		HttpPort:      Config.HttpPort,
		UrlBasePath:   Config.UrlBasePath,
		IamServiceUrl: Config.IamServiceUrl,
		Database:      Config.Database,
		IoConnectors:  Config.IoConnectors,
	}, Prefix+"/etc/main.json", "  ")
}
