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

package v1

import (
	"archive/tar"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/types"
	"github.com/ulikunitz/xz"

	"github.com/hooto/hpress/api"
	"github.com/hooto/hpress/config"
	"github.com/hooto/hpress/modset"
)

var (
	spec_upload_size_max int64 = 8 * 1024 * 1024
)

func (c ModSet) SpecUploadCommitAction() {

	var set api.SpecUploadCommit

	defer c.RenderJson(&set)

	err := c.Request.JsonDecode(&set)
	if err != nil {
		set.Error = types.NewErrorMeta(api.ErrCodeBadArgument, "Bad Argument "+err.Error())
		return
	}

	if set.Size > spec_upload_size_max {
		set.Error = types.NewErrorMeta("400",
			fmt.Sprintf("the max size of Package can not more than %d", spec_upload_size_max))
		return
	}

	if len(set.Name) < 10 || !strings.HasSuffix(set.Name, ".txz") {
		set.Error = types.NewErrorMeta(api.ErrCodeBadArgument, "Invalid Name")
		return
	}

	if !iamclient.SessionAccessAllowed(c.Session, "editor.write", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	body64 := strings.SplitAfter(set.Data, ";base64,")
	if len(body64) != 2 {
		return
	}
	filedata, err := base64.StdEncoding.DecodeString(body64[1])
	if err != nil {
		set.Error = types.NewErrorMeta("400", "Package Not Found")
		return
	}

	if int64(len(filedata)) != set.Size {
		set.Error = types.NewErrorMeta("400", "Invalid Package Size")
		return
	}

	xzr, err := xz.NewReader(bytes.NewReader(filedata))
	if err != nil {
		set.Error = types.NewErrorMeta(api.ErrCodeBadArgument, err.Error())
		return
	}

	tr := tar.NewReader(xzr)
	if tr == nil {
		set.Error = types.NewErrorMeta("400", "Invalid Encoded Data")
		return
	}

	pkg_name := set.Name[:len(set.Name)-4]

	tmpdir := config.Prefix + "/var/tmp/" + pkg_name

	// if _, err := os.Stat(tmpdir); err == nil {
	// 	return
	// }

	files := map[string]int64{}

	for {

		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			set.Error = types.NewErrorMeta(api.ErrCodeBadArgument, err.Error())
			return
		}
		// fmt.Printf("Contents of %s\n", hdr.Name)

		if hdr.Name[len(hdr.Name)-1] == '/' {
			os.MkdirAll(tmpdir+"/"+hdr.Name, 0755)
			continue
		}

		// if strings.Index(hdr.Name, "/") > 0 {
		// 	os.MkdirAll(tmpdir+"/"+filepath.Dir(hdr.Name), 0755)
		// }

		fpo, err := os.OpenFile(tmpdir+"/"+hdr.Name, os.O_RDWR|os.O_CREATE, os.FileMode(hdr.Mode))
		if err != nil {
			set.Error = types.NewErrorMeta(api.ErrCodeBadArgument, err.Error())
			return
		}
		fpo.Seek(0, 0)
		fpo.Truncate(0)

		if _, err := io.Copy(fpo, tr); err != nil {
			set.Error = types.NewErrorMeta(api.ErrCodeBadArgument, err.Error())
			return
		}

		fpo.Close()

		files[hdr.Name] = hdr.Mode
	}

	var spec api.Spec
	if err := json.DecodeFile(tmpdir+"/spec.json", &spec); err != nil {
		set.Error = types.NewErrorMeta(api.ErrCodeBadArgument, err.Error())
		return
	}

	//
	spec.Meta.Name, err = modset.ModNameFilter(spec.Meta.Name)
	if err != nil {
		set.Error = types.NewErrorMeta(api.ErrCodeBadArgument, err.Error())
		return
	}

	spec.SrvName, err = api.SrvNameFilter(spec.SrvName)
	if err != nil {
		set.Error = types.NewErrorMeta(api.ErrCodeBadArgument, err.Error())
		return
	}

	specVersion, _ := strconv.Atoi(spec.Meta.ResourceVersion)

	if prev, err := modset.SpecFetch(spec.Meta.Name); err == nil {
		prevVersion, _ := strconv.Atoi(prev.Meta.ResourceVersion)
		if prevVersion > specVersion {
			set.Error = types.NewErrorMeta(api.ErrCodeBadArgument, "Invalid Version")
			return
		}
	}

	spec_dir := config.Prefix + "/modules/" + spec.Meta.Name

	for path, fmode := range files {
		if err := spec_file_sync(tmpdir+"/"+path, spec_dir+"/"+path, os.FileMode(fmode)); err != nil {
			set.Error = types.NewErrorMeta(api.ErrCodeBadArgument, err.Error())
			return
		}
	}

	// timenow := rdb.TimeNow("datetime")
	// jsb, _ := json.Encode(spec, "  ")

	// q := rdb.NewQuerySet().From("modules").Limit(1)
	// q.Where.And("name", spec.Meta.Name)

	// fields := map[string]interface{}{
	// 	"title":   spec.Title,
	// 	"version": spec.Meta.ResourceVersion,
	// 	"updated": timenow,
	// 	"body":    string(jsb),
	// }

	// if _, err := store.Data.Fetch(q); err == nil {
	// 	fr := rdb.NewFilter()
	// 	fr.And("name", spec.Meta.Name)
	// 	store.Data.Update("modules", fields, fr)
	// } else {
	// 	fields["name"] = spec.Meta.Name
	// 	fields["srvname"] = spec.SrvName
	// 	fields["created"] = timenow
	// 	fields["status"] = 1

	// 	store.Data.Insert("modules", fields)
	// }

	// config.SpecSet(&spec)

	// config.SpecSrvRefresh(spec.SrvName)

	modset.SpecSchemaSync(spec)

	// } else {

	// 	if err = modset.SpecInfoSet(set); err != nil {
	// 		set.Error = types.NewErrorMeta(api.ErrCodeInternalError, err.Error())
	// 		return
	// 	}
	// }

	// seted, err := modset.SpecFetch(set.Meta.Name)
	// if err != nil {
	// 	set.Error = types.NewErrorMeta(api.ErrCodeInternalError, err.Error())
	// 	return
	// }

	// modset.SpecSchemaSync(seted)

	set.Kind = "Spec"
}

func spec_file_sync(src, dst string, mod os.FileMode) error {

	fp_src, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fp_src.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	fp_dst, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, mod)
	if err != nil {
		return err
	}
	defer fp_dst.Close()

	fp_dst.Seek(0, 0)
	fp_dst.Truncate(0)

	if _, err := io.Copy(fp_dst, fp_src); err != nil {
		return err
	}

	return nil
}
