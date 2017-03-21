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

	"code.hooto.com/lynkdb/iomix/skv"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/sync"

	// "github.com/shanemhansen/gogd"
	// "github.com/bamiaux/rez"
	// "github.com/disintegration/imaging"
	// gd "github.com/eryx/go-gd"
	"github.com/nfnt/resize"

	"code.hooto.com/hooto/hootopress/config"
	"code.hooto.com/hooto/hootopress/store"
)

var (
	rezlocker = sync.NewPermitPool(1)
)

type S2 struct {
	*httpsrv.Controller
}

// // rez
// func (c S2) Index_rezAction() {

// 	c.AutoRender = false

// 	var (
// 		ipn         = c.Params.Get("ipn")
// 		object_path = strings.Trim(filepath.Clean(c.Request.RequestPath), "/")[3:]
// 		abs_path    = config.Prefix + "/var/storage/" + object_path
// 	)

// 	if ipn == "" {

// 		if objfp, err := os.Open(abs_path); err == nil {

// 			c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
// 			http.ServeContent(c.Response.Out, c.Request.Request, object_path, time.Now(), objfp)
// 			objfp.Close()

// 		} else {
// 			c.RenderError(404, "Object Not Found")
// 		}

// 		return
// 	}

// 	ext := strings.ToLower(filepath.Ext(object_path))

// 	switch ext {

// 	case ".jpg", ".jpeg", ".png", ".gif":
// 		//

// 	default:
// 		c.RenderError(400, "Bad Request #01")
// 		return
// 	}

// 	h1 := md5.New()
// 	io.WriteString(h1, object_path+"."+ipn)
// 	hid := fmt.Sprintf("%x", h1.Sum(nil))[0:12]
// 	if rs := store.CacheGet(hid); rs.Status == "OK" {

// 		c.Response.Out.Header().Set("Cache-Control", "max-age=86400")

// 		switch ext {
// 		case ".jpg", ".jpeg":
// 			c.Response.Out.Header().Set("Content-type", "image/jpeg")
// 		case ".png":
// 			c.Response.Out.Header().Set("Content-type", "image/png")
// 		case ".git":
// 			c.Response.Out.Header().Set("Content-type", "image/gif")
// 		}

// 		c.Response.Out.Write(rs.Bytes())
// 		return
// 	}

// 	width, height, crop := 200, 200, true

// 	switch ipn {
// 	case "i6040":
// 		width, height = 60, 40
// 	case "s800":
// 		width, height, crop = 800, 800, false
// 	case "s800x":
// 		width, height, crop = 800, 8000, false
// 	case "thumb":
// 		width, height = 150, 150
// 	case "medium":
// 		width, height = 300, 168
// 	case "large":
// 		width, height = 800, 450
// 	}

// 	pn := rezlocker.Pull()
// 	defer rezlocker.Push(pn)

// 	defer func() {
// 		if x := recover(); x != nil {
// 			fmt.Println(x)
// 			c.RenderError(404, "Object Not Found")
// 		}
// 	}()

// 	fp, err := os.Open(abs_path)
// 	if err != nil {
// 		return
// 	}
// 	defer fp.Close()

// 	//
// 	var imgs image.Image

// 	switch ext {

// 	case ".jpg", ".jpeg":
// 		imgs, err = jpeg.Decode(fp)
// 	case ".png":
// 		imgs, err = png.Decode(fp)

// 	default:
// 		fmt.Println("error format")
// 		return
// 	}

// 	if err != nil {
// 		return
// 	}

// 	//
// 	// var imgd image.Image

// 	if crop {
// 		// imgd = resize.Thumbnail(uint(width), uint(height), imgs, resize.Lanczos3)
// 	} else {
// 		// imgd = resize.Resize(uint(width), uint(height), imgs, resize.Lanczos3)
// 	}

// 	// src, ok := imgs.(*image.YCbCr)
// 	// if !ok {
// 	// 	fmt.Println("001")
// 	// 	return
// 	// }

// 	p := imgs.Bounds().Size()
// 	w, h := width, getSize(width, p.Y, p.X)
// 	if p.X < p.Y {
// 		w, h = getSize(height, p.X, p.Y), height
// 	}

// 	imgd := image.NewYCbCr(image.Rect(0, 0, w, h), image.YCbCrSubsampleRatio444)

// 	if err := rez.Convert(imgd, imgs, rez.NewBilinearFilter()); err != nil {
// 		fmt.Println("002", err)
// 		return
// 	}

// 	// buf := new(bytes.Buffer)

// 	switch ext {

// 	case ".jpg", ".jpeg":

// 		c.Response.Out.Header().Set("Content-type", "image/jpeg")
// 		// jpeg.Encode(buf, imgd, &jpeg.Options{90})
// 		jpeg.Encode(c.Response.Out, imgd, &jpeg.Options{90})

// 	case ".png":

// 		c.Response.Out.Header().Set("Content-type", "image/png")
// 		png.Encode(c.Response.Out, imgd)

// 	case ".gif":

// 		c.Response.Out.Header().Set("Content-type", "image/gif")
// 		// gif.Encode(buf, imgd, &gif.Options{NumColors: 256})
// 	}

// 	// if len(buf) > 10 {
// 	// 	// store.CacheSetBytes([]byte(hid), buf, int64(36000+rand.Intn(36000))*1000)
// 	// }

// 	c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
// 	// c.Response.Out.Write(buf.Bytes())
// }

// func getSize(a, b, c int) int {
// 	d := a * b / c
// 	return (d + 1) & -1
// }

// resize
func (c S2) IndexAction() {

	c.AutoRender = false

	var (
		ipn         = c.Params.Get("ipn")
		object_path = strings.Trim(filepath.Clean(c.Request.RequestPath), "/")[3:]
		abs_path    = config.Prefix + "/var/storage/" + object_path
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

	var (
		ext       = strings.ToLower(filepath.Ext(object_path))
		meta_type = ""
	)

	switch ext {

	case ".jpg", ".jpeg":
		meta_type = "image/jpeg"

	case ".png":
		meta_type = "image/png"

	case ".gif":
		meta_type = "image/gif"

	default:
		c.RenderError(400, "Bad Request #01")
		return
	}

	hid := "s2." + idhash.HashToHexString([]byte(object_path+"."+ipn), 12)

	if rs := store.LocalCache.KvGet([]byte(hid)); rs.OK() {

		c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
		c.Response.Out.Header().Set("Content-type", meta_type)
		c.Response.Out.Write(rs.Bytex().Bytes())

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
		dst_img image.Image
		dst_buf = new(bytes.Buffer)
	)

	if crop {
		dst_img = resize.Thumbnail(uint(width), uint(height), src_img, resize.Lanczos3)
	} else {
		dst_img = resize.Resize(uint(width), uint(height), src_img, resize.Lanczos3)
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
			TimeToLive: int64(36000+rand.Intn(36000)) * 1000,
		})
	}

	//
	c.Response.Out.Header().Set("Content-type", meta_type)
	c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
	c.Response.Out.Write(dst_buf.Bytes())
}

// // imaging
// func (c S2) Index_imagingAction() {

// 	c.AutoRender = false

// 	var (
// 		ipn         = c.Params.Get("ipn")
// 		object_path = strings.Trim(filepath.Clean(c.Request.RequestPath), "/")[3:]
// 		abs_path    = config.Prefix + "/var/storage/" + object_path
// 	)

// 	if ipn == "" {

// 		if objfp, err := os.Open(abs_path); err == nil {

// 			c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
// 			http.ServeContent(c.Response.Out, c.Request.Request, object_path, time.Now(), objfp)
// 			objfp.Close()

// 		} else {
// 			c.RenderError(404, "Object Not Found")
// 		}

// 		return
// 	}

// 	ext := strings.ToLower(filepath.Ext(object_path))

// 	switch ext {

// 	case ".jpg", ".jpeg", ".png", ".gif":
// 		//

// 	default:
// 		c.RenderError(400, "Bad Request #01")
// 		return
// 	}

// 	h1 := md5.New()
// 	io.WriteString(h1, object_path+"."+ipn)
// 	hid := fmt.Sprintf("%x", h1.Sum(nil))[0:12]
// 	if rs := store.CacheGet(hid); rs.Status == "OK" {

// 		c.Response.Out.Header().Set("Cache-Control", "max-age=86400")

// 		switch ext {
// 		case ".jpg", ".jpeg":
// 			c.Response.Out.Header().Set("Content-type", "image/jpeg")
// 		case ".png":
// 			c.Response.Out.Header().Set("Content-type", "image/png")
// 		case ".git":
// 			c.Response.Out.Header().Set("Content-type", "image/gif")
// 		}

// 		c.Response.Out.Write(rs.Bytes())
// 		return
// 	}

// 	width, height, crop := 200, 200, true

// 	switch ipn {
// 	case "i6040":
// 		width, height = 60, 40
// 	case "s800":
// 		width, height, crop = 800, 800, false
// 	case "s800x":
// 		width, height, crop = 800, 8000, false
// 	case "thumb":
// 		width, height = 150, 150
// 	case "medium":
// 		width, height = 300, 168
// 	case "large":
// 		width, height = 800, 450
// 	}

// 	pn := rezlocker.Pull()
// 	defer rezlocker.Push(pn)

// 	defer func() {
// 		if x := recover(); x != nil {
// 			c.RenderError(404, "Object Not Found")
// 		}
// 	}()

// 	imgs, err := imaging.Open(abs_path)
// 	if err != nil {
// 		return
// 	}

// 	var imgd *image.NRGBA

// 	if crop {
// 		imgd = imaging.Thumbnail(imgs, width, height, imaging.Box)
// 	} else {
// 		imgd = imaging.Resize(imgs, width, height, imaging.Box)
// 	}

// 	if imgd == nil {
// 		return
// 	}

// 	// buf := new(bytes.Buffer)

// 	switch ext {

// 	case ".jpg", ".jpeg":

// 		c.Response.Out.Header().Set("Content-type", "image/jpeg")
// 		// jpeg.Encode(buf, imgd, &jpeg.Options{90})
// 		jpeg.Encode(c.Response.Out, imgd, &jpeg.Options{90})

// 	case ".png":

// 		c.Response.Out.Header().Set("Content-type", "image/png")
// 		png.Encode(c.Response.Out, imgd)

// 	case ".gif":

// 		c.Response.Out.Header().Set("Content-type", "image/gif")
// 		// gif.Encode(buf, imgd, &gif.Options{NumColors: 256})
// 	}

// 	// if len(buf) > 10 {
// 	// 	// store.CacheSetBytes([]byte(hid), buf, int64(36000+rand.Intn(36000))*1000)
// 	// }

// 	c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
// 	// c.Response.Out.Write(buf.Bytes())
// }

// // gd
// func (c S2) Index_gdAction() {

// 	c.AutoRender = false

// 	var (
// 		ipn         = c.Params.Get("ipn")
// 		object_path = strings.Trim(filepath.Clean(c.Request.RequestPath), "/")[3:]
// 		abs_path    = config.Prefix + "/var/storage/" + object_path
// 	)

// 	if ipn == "" {

// 		if objfp, err := os.Open(abs_path); err == nil {

// 			c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
// 			http.ServeContent(c.Response.Out, c.Request.Request, object_path, time.Now(), objfp)
// 			objfp.Close()

// 		} else {
// 			c.RenderError(404, "Object Not Found")
// 		}

// 		return
// 	}

// 	ext := strings.ToLower(filepath.Ext(object_path))

// 	switch ext {

// 	case ".jpg", ".jpeg", ".png", ".gif":
// 		//

// 	default:
// 		c.RenderError(400, "Bad Request #01")
// 		return
// 	}

// 	// h1 := md5.New()
// 	// io.WriteString(h1, object_path+"."+ipn)
// 	// hid := fmt.Sprintf("%x", h1.Sum(nil))[0:12]
// 	// if rs := store.CacheGet(hid); rs.Status == "OK" {

// 	// 	c.Response.Out.Header().Set("Cache-Control", "max-age=86400")

// 	// 	switch ext {
// 	// 	case ".jpg", ".jpeg":
// 	// 		c.Response.Out.Header().Set("Content-type", "image/jpeg")
// 	// 	case ".png":
// 	// 		c.Response.Out.Header().Set("Content-type", "image/png")
// 	// 	case ".git":
// 	// 		c.Response.Out.Header().Set("Content-type", "image/gif")
// 	// 	}

// 	// 	c.Response.Out.Write(rs.Bytes())
// 	// 	return
// 	// }

// 	width, height, crop := 2000, 2000, true

// 	switch ipn {
// 	case "i6040":
// 		width, height = 60, 40
// 	case "s800":
// 		width, height, crop = 800, 800, false
// 	case "s800x":
// 		width, height, crop = 800, 8000, false
// 	case "thumb":
// 		width, height = 150, 150
// 	case "medium":
// 		width, height = 300, 168
// 	case "large":
// 		width, height = 800, 450
// 	}

// 	pn := rezlocker.Pull()
// 	defer rezlocker.Push(pn)

// 	// buf := new(bytes.Buffer)
// 	// buf.ReadFrom(objfp)

// 	defer func() {
// 		if x := recover(); x != nil {
// 			c.RenderError(404, "Object Not Found")
// 		}
// 	}()

// 	out := _img_resize(ext, abs_path, width, height, crop)
// 	if len(out) < 10 {
// 		c.RenderError(404, "Object Not Found")
// 		return
// 	}

// 	// store.CacheSetBytes([]byte(hid), _bytes_clone(out), int64(864000+rand.Intn(864000))*1000)

// 	//
// 	switch ext {

// 	case ".jpg", ".jpeg":

// 		c.Response.Out.Header().Set("Content-type", "image/jpeg")

// 	case ".png":

// 		c.Response.Out.Header().Set("Content-type", "image/jpeg")

// 	case ".gif":

// 		c.Response.Out.Header().Set("Content-type", "image/gif")
// 	}

// 	c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
// 	c.Response.Out.Write(_bytes_clone(out))
// }

// func _bytes_clone(src []byte) []byte {

// 	dst := make([]byte, len(src))
// 	copy(dst, src)

// 	return dst
// }

// func _img_resize(ext, abs_path string, width, height int, crop bool) []byte {

// 	var out []byte

// 	ims := create_image(ext, abs_path)
// 	if ims == nil {
// 		return out
// 	}
// 	defer ims.Destroy()

// 	sx, sy, sw, sh, dstw, dsth, _ := _resize_dimensions(ims.Sx(), ims.Sy(), width, height, crop)

// 	imd := gd.CreateTrueColor(dstw, dsth)
// 	defer imd.Destroy()

// 	if ext == ".png" {
// 		imd.AlphaBlending(false)
// 		imd.SaveAlpha(true)
// 	}

// 	ims.CopyResampled(imd, 0, 0, sx, sy, dstw, dsth, sw, sh)

// 	switch ext {

// 	case ".jpg", ".jpeg":
// 		out = gd.ImageToJpegBuffer(imd, 85)

// 	case ".png":
// 		out = gd.ImageToPngBuffer(imd)

// 	case ".gif":
// 		out = gd.ImageToGifBuffer(imd)

// 	}

// 	return _bytes_clone(out)
// }

// func create_image(ext, abs_path string) *gd.Image {

// 	switch ext {

// 	case ".jpg", ".jpeg":

// 		return gd.CreateFromJpeg(abs_path)

// 	case ".png":

// 		return gd.CreateFromPng(abs_path)

// 	case ".gif":

// 		return gd.CreateFromGif(abs_path)
// 	}

// 	return nil
// }

// func _resize_dimensions(srcw, srch, dstw, dsth int, crop bool) (sx, sy, sw, sh, dw, dh int, do bool) {

// 	if dstw > srcw && dsth > srch {
// 		return 0, 0, srcw, srch, srcw, srch, true
// 	}

// 	if !crop {

// 		cropw_ratio := float32(dstw) / float32(srcw)
// 		croph_ratio := float32(dsth) / float32(srch)

// 		if cropw_ratio < croph_ratio {
// 			dsth = int(cropw_ratio * float32(srch))
// 		} else {
// 			dstw = int(croph_ratio * float32(srcw))
// 		}
// 	}

// 	dst_ratio := float32(dstw) / float32(dsth)
// 	sdw := int(_float_min(float32(srch)*dst_ratio, float32(srcw)))
// 	sdh := int(_float_min(float32(srcw)/dst_ratio, float32(srch)))

// 	if dstw > sdw {
// 		dstw = sdw
// 	}

// 	if dsth > sdh {
// 		dsth = sdh
// 	}

// 	sdx := int((srcw - sdw) / 2)
// 	sdy := int((srch - sdh) / 2)

// 	return sdx, sdy, sdw, sdh, dstw, dsth, true
// }

// func _float_min(a, b float32) float32 {

// 	if a < b {
// 		return a
// 	}

// 	return b
// }
