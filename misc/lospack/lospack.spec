project.name = hooto-press
project.version = 0.2.2.dev
project.vendor = hooto.com
project.homepage = https://github.com/hooto/hooto-press
project.groups = app/other

%build
export PATH=$PATH:/usr/local/go/bin:/opt/gopath/bin
export GOPATH=/opt/gopath
mkdir -p {{.buildroot}}/bin
mkdir -p {{.buildroot}}/var/{captchadb,log,storage,htp_local_cache}
go build -ldflags "-s -w" -o {{.buildroot}}/bin/hooto-press main.go

sed -i 's/debug:\!0/debug:\!1/g' {{.buildroot}}/webui/htpm/js/main.js

%files
bin/hooto-press
bin/keeper
etc/main.json.tpl
i18n/
modules/
webui/htpm/
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
webui/htpm/js/
modules/
vendor/code.hooto.com/hooto/chart/webui/

%css_compress
webui/bs/3.3/css/
webui/cm/5/
webui/purecss/pure.css
webui/lessui/css/lessui.css
webui/htp/css/
webui/htpm/css/
modules/


%html_compress
modules/
webui/htpm/
websrv/mgr/views/

%png_compress
webui/htp/img/
webui/htpm/img/

