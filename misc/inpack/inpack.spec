[project]
name = hooto-press
version = 0.3.2.alpha
vendor = hooto.com
homepage = https://github.com/hooto/hpress
groups = app/other

%build
export PATH=$PATH:/usr/local/go/bin:/opt/gopath/bin
export GOPATH=/opt/gopath
mkdir -p {{.buildroot}}/bin
mkdir -p {{.buildroot}}/var/{hcaptchadb,log,storage,hpress_local_cache}
go build -ldflags "-s -w -X main.version={{.project__version}} -X main.release={{.project__release}}" -o {{.buildroot}}/bin/hooto-press main.go

sed -i 's/debug:\!0/debug:\!1/g' {{.buildroot}}/webui/hpress/js/main.js
sed -i 's/debug:\!0/debug:\!1/g' {{.buildroot}}/webui/hpressm/js/main.js

%files
bin/hooto-press
bin/keeper
etc/main.json.tpl
i18n/
modules/
webui/hpressm/
webui/hpress/img/alpha2.png
webui/hpress/img/search-16.png
webui/hpress/img/ap.ico
webui/bs/3.3/fonts/


%js_compress
webui/bs/3.3/js/bootstrap.js
webui/cm/5/
webui/lessui/js/lessui.js
webui/lessui/js/browser-detect.js
webui/lessui/js/eventproxy.js
webui/lessui/js/sea.js
webui/hpress/js/
webui/hpressm/js/
modules/
vendor/github.com/hooto/hchart/webui/

%css_compress
webui/bs/3.3/css/
webui/cm/5/
webui/purecss/pure.css
webui/lessui/css/lessui.css
webui/hpress/css/
webui/hpressm/css/
modules/


%html_compress
modules/
webui/hpressm/
websrv/mgr/views/

%png_compress
webui/hpress/img/
webui/hpressm/img/

