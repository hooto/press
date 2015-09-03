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
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/eryx/hcaptcha/captcha"
	"github.com/lessos/lessdb/skv"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/utils"

	"../api"
)

var (
	Config        ConfigCommon
	Version       = "0.0.2.dev"
	captchaConfig captcha.Options

	User = &user.User{
		Uid:      "2048",
		Gid:      "2048",
		Username: "action",
		HomeDir:  "/home/action",
	}

	SysConfigList = api.SysConfigList{}
)

type ConfigCommon struct {
	Prefix             string      `json:"prefix,omitempty"`
	ModuleDir          string      `json:"module_dir,omitempty"`
	InstanceID         string      `json:"instance_id"`
	AppTitle           string      `json:"app_title,omitempty"`
	HttpPort           uint16      `json:"http_port"`
	IdentityServiceUrl string      `json:"identity_service_url"`
	Database           base.Config `json:"database"`
	CacheDB            skv.Config  `json:"cache_db,omitempty"`
}

func init() {

	SysConfigList.Kind = "SysConfigList"

	SysConfigList.Insert(api.SysConfig{"frontend_header_site_name", "CMS", "Site's Name", ""})
	SysConfigList.Insert(api.SysConfig{"frontend_header_site_logo_url", "", "", ""})
	SysConfigList.Insert(api.SysConfig{"frontend_footer_copyright", "2015 Demo", "", ""})

	SysConfigList.Insert(api.SysConfig{"frontend_html_head_subtitle", "CMS", "Sub Title for HTML Head Title", ""})
	SysConfigList.Insert(api.SysConfig{"frontend_html_head_meta_keywords", "", "Meta Keywords in HTML Head for Search engine optimization", ""})
	SysConfigList.Insert(api.SysConfig{"frontend_html_head_meta_description", "", "Meta Description in HTML Head for Search engine optimization", ""})

	SysConfigList.Insert(api.SysConfig{"frontend_footer_analytics_scripts", "",
		"Embeded analytics scripts, ex. Google Analytics or Piwik ...", "text"})
	SysConfigList.Insert(api.SysConfig{"ls2_uri", "//127.0.0.1/bucket-name", "Storage Service URI", ""})
}

func Initialize(prefix string) error {

	var err error

	if prefix == "" {
		prefix, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/..")
		if err != nil {
			prefix = "/opt/lesscms"
		}
	}
	reg, _ := regexp.Compile("/+")
	Config.Prefix = "/" + strings.Trim(reg.ReplaceAllString(prefix, "/"), "/")
	Config.ModuleDir = Config.Prefix + "/modules"

	file := Config.Prefix + "/etc/main.json"
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return errors.New("Error: config file is not exists")
	}

	fp, err := os.Open(file)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: Can not open (%s)", file))
	}
	defer fp.Close()

	cfgstr, err := ioutil.ReadAll(fp)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: Can not read (%s)", file))
	}

	if err = utils.JsonDecode(cfgstr, &Config); err != nil {
		return errors.New(fmt.Sprintf("Error: "+
			"config file invalid. (%s)", err.Error()))
	}

	//
	if Config.IdentityServiceUrl == "" {
		return errors.New("Error: `identity_service_url` can not be null")
	}

	httpsrv.GlobalService.Config.InstanceID = Config.InstanceID

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
		Config.AppTitle = "less CMS"
	}

	// Setting CAPTCHA
	captchaConfig = captcha.DefaultConfig

	captchaConfig.FontPath = Config.Prefix + "/vendor/github.com/eryx/hcaptcha/var/fonts/cmr10.ttf"
	captchaConfig.DataDir = Config.Prefix + "/var/captchadb"

	if err := captcha.Config(captchaConfig); err != nil {
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

	cfgExport := ConfigCommon{
		InstanceID:         Config.InstanceID,
		AppTitle:           Config.AppTitle,
		HttpPort:           Config.HttpPort,
		IdentityServiceUrl: Config.IdentityServiceUrl,
		Database:           Config.Database,
	}

	jsb, _ := utils.JsonEncodeIndent(cfgExport, "  ")

	file := Config.Prefix + "/etc/main.json"
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return errors.New("Error: config file is not exists")
	}

	fp, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: Can not open (%s)", file))
	}
	defer fp.Close()

	fp.Seek(0, 0)
	fp.Truncate(int64(len(jsb)))

	_, err = fp.Write(jsb)

	httpsrv.GlobalService.Config.InstanceID = Config.InstanceID

	return err
}
