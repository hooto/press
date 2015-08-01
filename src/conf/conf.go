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

package conf

import (
	"errors"
	"fmt"

	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/eryx/hcaptcha/captcha"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/utils"
)

var (
	Config        ConfigCommon
	Version       = "1.0.0.dev00"
	captchaConfig captcha.Options
)

type ConfigCommon struct {
	Prefix             string      `json:"prefix,omitempty"`
	InstanceID         string      `json:"instance_id"`
	AppTitle           string      `json:"app_title,omitempty"`
	HttpPort           uint16      `json:"http_port"`
	IdentityServiceUrl string      `json:"identity_service_url"`
	Database           base.Config `json:"database"`
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

	fmt.Println(Config.Prefix)

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

	if _, err = Config.DatabaseInstance(); err != nil {
		return err
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
