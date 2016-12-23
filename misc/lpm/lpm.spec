project.name = hooto-alphapress
project.version = 0.1.7.dev


%build
export PATH=$PATH:/usr/local/go/bin:/opt/gopath/bin
export GOPATH=/opt/gopath
mkdir -p {{.buildroot}}/bin
mkdir -p {{.buildroot}}/var/{cachedb,captchadb,log,storage}
go build -ldflags "-s -w" -o {{.buildroot}}/bin/hooto-alphapress main.go


%files
bin/hooto-alphapress
bin/keeper
etc/main.json.tpl
i18n/
modules/
websrv/mgr/tpls/
webui/htap/img/alpha2.png
webui/htap/img/search-16.png
webui/htap/img/ap.ico
webui/bs/3.3/fonts/


%js_compress
webui/bs/3.3/js/bootstrap.js
webui/cm/5/
webui/lessui/js/lessui.js
webui/lessui/js/browser-detect.js
webui/lessui/js/eventproxy.js
webui/lessui/js/sea.js
webui/htap/js/
websrv/mgr/tpls/js/
modules/


%css_compress
webui/bs/3.3/css/
webui/cm/5/
webui/purecss/pure.css
webui/lessui/css/lessui.css
webui/htap/css/
websrv/mgr/tpls/css/
modules/


%html_compress
modules/
websrv/mgr/tpls/
websrv/mgr/views/

%png_compress
webui/htap/img/
