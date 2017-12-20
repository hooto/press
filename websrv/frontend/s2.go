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
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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
	rezlocker = sync.NewPermitPool(1)
)

type S2 struct {
	*httpsrv.Controller
}

// resize
func (c S2) IndexAction() {

	c.AutoRender = false

	var (
		ipn      = c.Params.Get("ipn")
		obj_path = strings.TrimPrefix(filepath.Clean(c.Request.RequestPath), "hp/s2/")
		abs_path = config.Prefix + "/var/storage/" + obj_path
	)

	if ipn == "" {

		if fp, err := os.Open(abs_path); err == nil {

			c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
			http.ServeContent(c.Response.Out, c.Request.Request, obj_path, time.Now(), fp)
			fp.Close()

		} else {
			c.RenderError(404, "Object Not Found")
		}

		return
	}

	var (
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

	width, height, crop := 200, 200, true

	switch ipn {
	case "i6040":
		width, height = 60, 40

	case "s800":
		width, height, crop = 800, 800, false

	case "s800x":
		width, height, crop = 800, 8000, false

	case "thumb":
		width, height = 150, 150

	case "medium":
		width, height = 300, 168

	case "large":
		width, height = 800, 450

	default:
		c.RenderError(400, "Bad Request (ipn)")
		return
	}

	hid := "s2." + idhash.HashToHexString([]byte(obj_path+"."+ipn), 12)

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
	if width < src_bounds.Dx() && height < src_bounds.Dy() {

		if crop {
			if im, err := cutter.Crop(src_img, cutter.Config{
				Width:   width,
				Height:  height,
				Mode:    cutter.Centered,
				Options: cutter.Ratio,
			}); err == nil {
				dst_img = resize.Thumbnail(uint(width), uint(height), im, resize.Lanczos3)
			}
		}

		if dst_img == nil {
			dst_img = resize.Thumbnail(uint(width), uint(height), src_img, resize.Lanczos3)
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
