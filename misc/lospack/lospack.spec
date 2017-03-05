project.name = hootopress
project.version = 0.1.7.dev
project.vendor = hooto.com
project.homepage = https://code.hooto.com/hooto/hootopress
project.groups = app/other

%build
export PATH=$PATH:/usr/local/go/bin:/opt/gopath/bin
export GOPATH=/opt/gopath
mkdir -p {{.buildroot}}/bin
mkdir -p {{.buildroot}}/var/{cachedb,captchadb,log,storage}
go build -ldflags "-s -w" -o {{.buildroot}}/bin/hootopress main.go


%files
bin/hootopress
bin/keeper
etc/main.json.tpl
i18n/
modules/
websrv/mgr/tpls/
webui/htp/img/alpha2.png
webui/htp/img/search-16.png
webui/htp/img/ap.ico
webui/bs/3.3/fonts/


%js_compress
webui/bs/3.3/js/bootstrap.js
webui/cm/5/
webui/lessui/js/lessui.js
webui/lessui/js/browser-detect.js
webui/lessui/js/eventproxy.js
webui/lessui/js/sea.js
webui/htp/js/
websrv/mgr/tpls/js/
modules/


%css_compress
webui/bs/3.3/css/
webui/cm/5/
webui/purecss/pure.css
webui/lessui/css/lessui.css
webui/htp/css/
websrv/mgr/tpls/css/
modules/


%html_compress
modules/
websrv/mgr/tpls/
websrv/mgr/views/

%png_compress
webui/htp/img/
