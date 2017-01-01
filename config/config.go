// Copyright 2015 lessOS.com, All rights reserved.
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

	"github.com/eryx/hcaptcha/captcha"
	"github.com/lessos/lessdb/skv"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/encoding/json"

	"code.hooto.com/hooto/hootopress/api"
)

var (
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
	Prefix        string      `json:"prefix,omitempty"`
	UrlBasePath   string      `json:"url_base_path,omitempty"`
	ModuleDir     string      `json:"module_dir,omitempty"`
	InstanceID    string      `json:"instance_id"`
	AppTitle      string      `json:"app_title,omitempty"`
	HttpPort      uint16      `json:"http_port"`
	IamServiceUrl string      `json:"iam_service_url"`
	Database      base.Config `json:"database"`
	CacheDB       skv.Config  `json:"cache_db,omitempty"`
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
		"frontend_footer_copyright", "2016 Demo",
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

	Config.Prefix = filepath.Clean(prefix)
	Config.ModuleDir = Config.Prefix + "/modules"

	file := Config.Prefix + "/etc/main.json"
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("Error: config file is not exists")
	}

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
		Config.AppTitle = "hooto AlphaPress"
	}

	// Setting CAPTCHA
	CaptchaConfig.DataDir = Config.Prefix + "/var/captchadb"

	if err := captcha.Config(CaptchaConfig); err != nil {
		return err
	}

	// Default User
	if User, err = user.Current(); err != nil {
		return err
	}

	//
	Config.CacheDB = skv.Config{
		DataDir: Config.Prefix + "/var/cachedb",
	}

	return module_init()
}

func Save() error {

	file := Config.Prefix + "/etc/main.json"
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return errors.New("Error: config file is not exists")
	}

	cfgExport := ConfigCommon{
		InstanceID:    Config.InstanceID,
		AppTitle:      Config.AppTitle,
		HttpPort:      Config.HttpPort,
		UrlBasePath:   Config.UrlBasePath,
		IamServiceUrl: Config.IamServiceUrl,
		Database:      Config.Database,
	}

	return json.EncodeToFile(cfgExport, file, "  ")
}
