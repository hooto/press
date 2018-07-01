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

package frontend

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/sync"
	"github.com/lynkdb/iomix/skv"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"

	"github.com/hooto/hpress/config"
	"github.com/hooto/hpress/store"
)

var (
	rezlocker      = sync.NewPermitPool(1)
	s2_path_reg    = regexp.MustCompile("^[0-9a-zA-Z_\\-\\.\\/]{1,100}$")
	s2_bucket_deft = "/deft/"
	s2_url_prefix  = "hp/s2"
)

func path_filter(path string) (string, error) {

	path = filepath.Clean(strings.Replace(strings.TrimSpace(path), " ", "-", -1))
	if !s2_path_reg.MatchString(path) {
		return path, fmt.Errorf("Invalid File Name")
	}

	if !strings.HasPrefix(path, s2_bucket_deft) {
		return "", errors.New("Invalid Bucket Name")
	}

	return path, nil
}

type S2 struct {
	*httpsrv.Controller
}

// resize
func (c S2) IndexAction() {

	c.AutoRender = false

	obj_path, err := path_filter(strings.TrimPrefix(c.Request.RequestPath, s2_url_prefix))
	if err != nil {
		c.RenderError(404, "Object Not Found")
		return
	}

	var (
		ipn       = c.Params.Get("ipn")
		ipl       = c.Params.Get("ipl")
		ipls      = [2]int{0, 0} // width, height
		iplc      = false
		ipl_step  = 64
		ipl_smin  = 64
		ipl_smax  = 2048
		abs_path  = config.Prefix + "/var/storage/" + obj_path
		ext       = strings.ToLower(filepath.Ext(obj_path))
		meta_type = ""
	)

	switch ext {

	case ".jpg", ".jpeg":
		meta_type = "image/jpeg"

	case ".png":
		meta_type = "image/png"

	case ".gif":
		meta_type = "image/gif"

	case ".svg":
		meta_type = "image/svg+xml"

	default:
		c.RenderError(400, "Bad Request #01")
		return
	}

	if (ipn == "" && ipl == "") ||
		meta_type == "image/svg+xml" {

		if fp, err := os.Open(abs_path); err == nil {

			c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
			http.ServeContent(c.Response.Out, c.Request.Request, obj_path, time.Now(), fp)
			fp.Close()

		} else {
			c.RenderError(404, "Object Not Found")
		}

		return
	}

	if ipl != "" {

		ar := strings.Split(ipl, ",")

		for _, k := range ar {

			if k[0] < 'a' || k[0] > 'z' {
				continue
			}

			switch k[0] {
			case 'w':
				if len(k) > 1 {
					ipls[0], _ = strconv.Atoi(k[1:])
				}

			case 'h':
				if len(k) > 1 {
					ipls[1], _ = strconv.Atoi(k[1:])
				}

			case 'c':
				iplc = true
			}
		}

	} else if ipn != "" { // TORM

		switch ipn {
		case "i6040":
			ipls[0], ipls[1] = 60, 40

		case "s800":
			ipls[0], ipls[1], iplc = 800, 800, false

		case "s800x":
			ipls[0], ipls[1], iplc = 800, 8000, false

		case "thumb":
			ipls[0], ipls[1] = 150, 150

		case "medium":
			ipls[0], ipls[1] = 300, 168

		case "large":
			ipls[0], ipls[1] = 800, 450

		default:
			c.RenderError(400, "Bad Request (ipn)")
			return
		}
	}

	if ipls[0] > 0 && ipls[1] < 1 {
		ipls[1] = ipls[0]
	} else if ipls[1] > 0 && ipls[0] < 1 {
		ipls[0] = ipls[1]
	}

	for i := range ipls {
		ipls[i] = ipls[i] - (ipls[i] % ipl_step)

		if ipls[i] < ipl_smin {
			ipls[i] = ipl_smin
		} else if ipls[i] > ipl_smax {
			ipls[i] = ipl_smax
		}
	}

	var (
		key = fmt.Sprintf("%s.%d.%d.%t", obj_path, ipls[0], ipls[1], iplc)
		hid = "s2." + idhash.HashToHexString([]byte(key), 12)
	)

	if rs := store.LocalCache.KvGet([]byte(hid)); rs.OK() {
		c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
		c.Response.Out.Header().Set("Content-type", meta_type)
		c.Response.Out.Write(rs.Bytex().Bytes())
		return
	}

	pn := rezlocker.Pull()
	defer rezlocker.Push(pn)

	defer func() {
		if x := recover(); x != nil {
			c.RenderError(404, "Object Not Found")
		}
	}()

	fp, err := os.Open(abs_path)
	if err != nil {
		c.RenderError(400, "Bad Request (invalid object path)")
		return
	}
	defer fp.Close()

	//
	var src_img image.Image

	switch ext {

	case ".jpg", ".jpeg":
		src_img, err = jpeg.Decode(fp)

	case ".png":
		src_img, err = png.Decode(fp)

	case ".gif":
		src_img, err = gif.Decode(fp)
	}

	if err != nil {
		c.RenderError(400, "Bad Request (invalid object format)")
		return
	}

	//
	var (
		dst_img    image.Image
		dst_buf    = new(bytes.Buffer)
		src_bounds = src_img.Bounds()
	)

	if ipls[0] < src_bounds.Dx() || ipls[1] < src_bounds.Dy() {

		if iplc {
			if im, err := cutter.Crop(src_img, cutter.Config{
				Width:   ipls[0],
				Height:  ipls[1],
				Mode:    cutter.Centered,
				Options: cutter.Ratio,
			}); err == nil {
				dst_img = resize.Thumbnail(uint(ipls[0]), uint(ipls[1]), im, resize.Lanczos3)
			}
		} else {

			rate := float32(1.0)

			if hrate := float32(ipls[1]) / float32(src_bounds.Dy()); hrate < rate {
				rate = hrate
			}

			if wrate := float32(ipls[0]) / float32(src_bounds.Dx()); wrate < rate {
				rate = wrate
			}

			if rate < 1 {
				ipls[0] = int(float32(src_bounds.Dx()) * rate)
				ipls[1] = int(float32(src_bounds.Dy()) * rate)
			}
		}

		if dst_img == nil {
			dst_img = resize.Thumbnail(uint(ipls[0]), uint(ipls[1]), src_img, resize.Lanczos3)
		}

	} else {
		dst_img = src_img
	}

	//
	switch ext {

	case ".jpg", ".jpeg":
		err = jpeg.Encode(dst_buf, dst_img, &jpeg.Options{90})

	case ".png":
		err = png.Encode(dst_buf, dst_img)

	case ".gif":
		err = gif.Encode(dst_buf, dst_img, &gif.Options{NumColors: 256})
	}
	if err != nil {
		c.RenderError(400, "Bad Request : "+err.Error())
		return
	}

	//
	if dst_buf.Len() > 10 {
		store.LocalCache.KvPut([]byte(hid), dst_buf.Bytes(), &skv.KvWriteOptions{
			Ttl: int64(36000+rand.Intn(36000)) * 1000,
		})
	}

	//
	c.Response.Out.Header().Set("Content-type", meta_type)
	c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
	c.Response.Out.Write(dst_buf.Bytes())
}
