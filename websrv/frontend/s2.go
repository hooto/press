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

package frontend

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	gd "github.com/eryx/go-gd"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/sync"

	"../../config"
	"../../store"
)

var (
	rezlocker = sync.NewPermitPool(-1)
)

type S2 struct {
	*httpsrv.Controller
}

func (c S2) IndexAction() {

	c.AutoRender = false

	var (
		ipn         = c.Params.Get("ipn")
		object_path = strings.Trim(filepath.Clean(c.Request.RequestPath), "/")[3:]
		abs_path    = config.Config.Prefix + "/var/storage/" + object_path
	)

	if ipn == "" {

		if objfp, err := os.Open(abs_path); err == nil {

			c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
			http.ServeContent(c.Response.Out, c.Request.Request, object_path, time.Now(), objfp)
			objfp.Close()

		} else {
			c.RenderError(404, "Object Not Found")
		}

		return
	}

	ext := strings.ToLower(filepath.Ext(object_path))

	switch ext {

	case ".jpg", ".jpeg", ".png", ".gif":
		//

	default:
		c.RenderError(400, "Bad Request #01")
		return
	}

	h1 := md5.New()
	io.WriteString(h1, object_path+"."+ipn)
	hid := fmt.Sprintf("%x", h1.Sum(nil))[0:12]
	if rs := store.CacheGet(hid); rs.Status == "OK" {

		c.Response.Out.Header().Set("Cache-Control", "max-age=86400")

		switch ext {
		case ".jpg", ".jpeg":
			c.Response.Out.Header().Set("Content-type", "image/jpeg")
		case ".png":
			c.Response.Out.Header().Set("Content-type", "image/png")
		case ".git":
			c.Response.Out.Header().Set("Content-type", "image/gif")
		}

		c.Response.Out.Write(rs.Bytes())
		return
	}

	width, height, crop := 2000, 2000, true

	switch ipn {
	case "i6040":
		width, height = 60, 40
	case "s800":
		width, height, crop = 800, 800, false
	case "thumb":
		width, height = 150, 150
	case "medium":
		width, height = 300, 168
	case "large":
		width, height = 800, 450
	}

	pn := rezlocker.Pull()
	defer rezlocker.Push(pn)

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(objfp)

	defer func() {
		if x := recover(); x != nil {
			c.RenderError(404, "Object Not Found")
		}
	}()

	var ims *gd.Image

	switch ext {

	case ".jpg", ".jpeg":

		// ims = gd.CreateFromJpegFp(objfp)
		ims = gd.CreateFromJpeg(abs_path)
		// ims = gd.CreateFromJpegPtr(buf.Bytes())

	case ".png":

		// ims = gd.CreateFromPngFp(objfp)
		ims = gd.CreateFromPng(abs_path)
		// ims = gd.CreateFromPngPtr(buf.Bytes())

	case ".gif":

		// ims = gd.CreateFromGifFp(objfp)
		ims = gd.CreateFromGif(abs_path)
		// ims = gd.CreateFromGifPtr(buf.Bytes())
	}

	if ims == nil {
		return
	}

	defer ims.Destroy()

	sx, sy, sw, sh, dstw, dsth, _ := _resize_dimensions(ims.Sx(), ims.Sy(), width, height, crop)

	imd := gd.CreateTrueColor(dstw, dsth)
	defer imd.Destroy()

	if ext == ".png" {
		imd.AlphaBlending(false)
		imd.SaveAlpha(true)
	}

	ims.CopyResampled(imd, 0, 0, sx, sy, dstw, dsth, sw, sh)

	var out []byte

	switch ext {

	case ".jpg", ".jpeg":

		c.Response.Out.Header().Set("Content-type", "image/jpeg")
		out = gd.ImageToJpegBuffer(imd, 85)

	case ".png":

		c.Response.Out.Header().Set("Content-type", "image/jpeg")
		out = gd.ImageToPngBuffer(imd)

	case ".gif":

		c.Response.Out.Header().Set("Content-type", "image/gif")
		out = gd.ImageToGifBuffer(imd)

	}

	if len(out) > 10 {
		store.CacheSetBytes([]byte(hid), out, uint32(36000+rand.Intn(36000)))
	}

	// fmt.Println("rez size", len(out))

	c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
	c.Response.Out.Write(out)
}

func _resize_dimensions(srcw, srch, dstw, dsth int, crop bool) (sx, sy, sw, sh, dw, dh int, do bool) {

	if dstw > srcw && dsth > srch {
		return 0, 0, srcw, srch, srcw, srch, true
	}

	if !crop {

		cropw_ratio := float32(dstw) / float32(srcw)
		croph_ratio := float32(dsth) / float32(srch)

		if cropw_ratio < croph_ratio {
			dsth = int(cropw_ratio * float32(srch))
		} else {
			dstw = int(croph_ratio * float32(srcw))
		}
	}

	dst_ratio := float32(dstw) / float32(dsth)
	sdw := int(_float_min(float32(srch)*dst_ratio, float32(srcw)))
	sdh := int(_float_min(float32(srcw)/dst_ratio, float32(srch)))

	if dstw > sdw {
		dstw = sdw
	}

	if dsth > sdh {
		dsth = sdh
	}

	sdx := int((srcw - sdw) / 2)
	sdy := int((srch - sdh) / 2)

	return sdx, sdy, sdw, sdh, dstw, dsth, true
}

func _float_min(a, b float32) float32 {

	if a < b {
		return a
	}

	return b
}
