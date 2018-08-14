[project]
name = hooto-press
version = 0.4.4.alpha
vendor = hooto.com
homepage = https://github.com/hooto/hpress
groups = app/other

%build
export PATH=$PATH:/usr/local/go/bin:/opt/gopath/bin
export GOPATH=/opt/gopath
mkdir -p {{.buildroot}}/bin
mkdir -p {{.buildroot}}/var/{hcaptchadb,log,storage,hpress_local_cache,tmp}
go build -ldflags "-s -w -X main.version={{.project__version}} -X main.release={{.project__release}}" -o {{.buildroot}}/bin/hooto-press main.go
go build -ldflags "-s -w -X main.version={{.project__version}} -X main.release={{.project__release}}" -o {{.buildroot}}/bin/hp-s2-upgrade misc/s2-upgrade.go

sed -i 's/debug:\!0/debug:\!1/g' {{.buildroot}}/webui/hp/js/main.js
sed -i 's/debug:\!0/debug:\!1/g' {{.buildroot}}/webui/hpm/js/main.js
sed -i 's/debug:true/debug:false/g' {{.buildroot}}/webui/hp/js/main.js
sed -i 's/debug:true/debug:false/g' {{.buildroot}}/webui/hpm/js/main.js

%files
bin/hooto-press
bin/keeper
etc/config.json.tpl
i18n/
modules/core/
webui/hpm/
webui/hp/img/alpha2.png
webui/hp/img/search-16.png
webui/hp/img/ap.ico
webui/bs/3.3/fonts/
webui/octicons/
webui/open-iconic/
webui/katex/

%js_compress
webui/bs/3.3/js/bootstrap.js
webui/bs/4/js/bootstrap.js
webui/cm/5/
webui/lessui/js/lessui.js
webui/lessui/js/browser-detect.js
webui/lessui/js/eventproxy.js
webui/lessui/js/sea.js
webui/hp/js/
webui/hpm/js/
webui/katex/
modules/core/
vendor/github.com/hooto/hchart/webui/

%css_compress
webui/bs/3.3/css/
webui/bs/4/css/
webui/cm/5/
webui/purecss/pure.css
webui/lessui/css/lessui.css
webui/lessui/css/base.css
webui/hp/css/
webui/hpm/css/
webui/katex/
modules/core/


%html_compress
modules/core/
webui/hpm/
websrv/mgr/views/

%png_compress
webui/hp/img/
webui/hpm/img/

