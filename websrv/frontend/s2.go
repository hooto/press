package frontend

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/lessos/lessgo/httpsrv"

	"../../config"
	"../../store"
)

type S2 struct {
	*httpsrv.Controller
}

func (c S2) IndexAction() {

	c.AutoRender = false

	var (
		ipn         = c.Params.Get("ipn")
		object_path = strings.Trim(filepath.Clean(c.Request.RequestPath), "/")[3:]
	)

	obj, err := os.Open(config.Config.Prefix + "/var/storage/" + object_path)
	if err != nil {
		c.RenderError(404, "Object Not Found")
		return
	}

	if ipn == "" {
		c.Response.Out.Header().Set("Cache-Control", "max-age=86400")

		http.ServeContent(c.Response.Out, c.Request.Request, object_path, time.Now(), obj)

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

	width, height := 2000, 2000

	switch ipn {
	case "thumb":
		width, height = 150, 150
	case "medium":
		width, height = 300, 168
	case "large":
		width, height = 800, 450
	}

	srcImage, err := imaging.Decode(obj)

	dstImage := imaging.Fit(srcImage, width, height, imaging.Box)
	if dstImage != nil {

		buf := new(bytes.Buffer)

		switch ext {

		case ".jpg", ".jpeg":
			jpeg.Encode(buf, dstImage, &jpeg.Options{90})

		case ".png":
			png.Encode(buf, dstImage)

		case ".gif":
			gif.Encode(buf, dstImage, &gif.Options{NumColors: 256})

		default:
			c.RenderError(400, "Bad Request")
			return
		}

		store.CacheSetBytes([]byte(hid), buf.Bytes(), 3600)
		c.Response.Out.Write(buf.Bytes())

		return
	}

	c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
	http.ServeContent(c.Response.Out, c.Request.Request, object_path, time.Now(), obj)
}
