var htapMgr = {
    frtbase : "/",
    base    : "/mgr/",
    api     : "/v1/",
    basetpl : "/mgr/-/",
}

htapMgr.Boot = function()
{
    if (window._basepath) {
        htapMgr.frtbase = window._basepath;
        if (!htapMgr.frtbase || htapMgr.frtbase == "") {
            htapMgr.frtbase = "/";
        }
        htapMgr.base    = htapMgr.frtbase +"mgr/";
        htapMgr.api     = htapMgr.frtbase +"v1/";
        htapMgr.basetpl = htapMgr.frtbase +"mgr/-/";
    }

    seajs.config({
        base: htapMgr.base,
        alias: {
            ep: '~/lessui/js/eventproxy.js'
        },
    });

    seajs.use([
        "~/lessui/js/browser-detect.js",
        "~/htap/js/jquery.js",
        "~/lessui/js/eventproxy.js",
    ], function() {

        var browser = BrowserDetect.browser;
        var version = BrowserDetect.version;
        var OS      = BrowserDetect.OS;

        if (!((browser == 'Chrome' && version >= 22)
            || (browser == 'Firefox' && version >= 31.0))) {
            $('body').load(window._basepath +"/error/browser");
            return;
        }

        seajs.use([
            "~/bs/3.3/css/bootstrap.css",
            "~/purecss/pure.css",
            "~/lessui/js/lessui.js",
            "~/lessui/css/lessui.css",
            "~/htap/css/main.css?_="+ Math.random(),
            "~/htap/js/marked.js",
            "-/css/main.css?_="+ Math.random(),
            "-/css/defx.css",
            "-/js/spec.js?_="+ Math.random(),
            "-/js/spec-editor.js?_="+ Math.random(),
            "-/js/tablet.js",
            "-/js/lc-editor.js",
            "-/js/model.js?_="+ Math.random(),
            "-/js/term.js?_="+ Math.random(),
            "-/js/node.js?_="+ Math.random(),
            "-/js/sys.js?_="+ Math.random(),
            "-/js/s2.js?_="+ Math.random(),
            "-/js/editor.js?_="+ Math.random(),
        ], function() {

            setTimeout(htapMgr.BootInit, 300);
        });
    });
}

htapMgr.BootInit = function()
{
    $("#htapm-topbar").css({"display": "block"});

    htapSys.Init();

    marked.setOptions({
        renderer: new marked.Renderer(),
        gfm: true,
        tables: true,
        breaks: false,
        pedantic: false,
        sanitize: true,
        smartLists: true,
        smartypants: true
    });

    $(window).resize(function() {
        htapEditor.sizeRefresh();
    });


    htapSys.Init();
    htapSpec.Init();
    htapS2.Init();

    htapNode.Init(function() {

        var navlast = l4iStorage.Get("htapm_nav_last_active");
        if (!navlast) {
            navlast = "sys/index"
        }
        l4i.UrlEventHandler(navlast);
    });
}

htapMgr.HttpSrvBasePath = function(url)
{
    if (htapMgr.base == "") {
        return url;
    }

    if (url.substr(0, 1) == "/") {
        return url;
    }

    return htapMgr.base + url;
}

htapMgr.Ajax = function(url, options)
{
    options = options || {};

    //
    if (url.substr(0, 1) != "/" && url.substr(0, 4) != "http") {
        url = htapMgr.HttpSrvBasePath(url);
    }

    //
    if (/\?/.test(url)) {
        url += "&_=";
    } else {
        url += "?_=";
    }
    url += Math.random();

    //
    if (!options.method) {
        options.method = "GET";
    }

    //
    if (!options.timeout) {
        options.timeout = 10000;
    }

    //
    $.ajax({
        url     : url,
        type    : options.method,
        data    : options.data,
        timeout : options.timeout,
        success : function(rsp) {
            if (typeof options.callback === "function") {
                options.callback(null, rsp);
            }
            if (typeof options.success === "function") {
                options.success(rsp);
            }
        },
        error: function(xhr, textStatus, error) {
            if (typeof options.callback === "function") {
                options.callback(xhr.responseText, null);
            }
            if (typeof options.error === "function") {
                options.error(xhr, textStatus, error);
            }
        }
    });
}


htapMgr.ApiCmd = function(url, options)
{
    htapMgr.Ajax(htapMgr.api + url, options);
}


htapMgr.TplCmd = function(url, options)
{
    htapMgr.Ajax(htapMgr.basetpl + url +".tpl", options);
}


htapMgr.Loader = function(target, uri)
{
    htapMgr.Ajax(htapMgr.basetpl + uri +".tpl", {
        callback: function(err, data) {
            $(target).html(data);
        }
    });
}


htapMgr.BodyLoader = function(uri)
{
    htapMgr.Loader("#body-content", uri);
}


htapMgr.ComLoader = function(uri)
{
    htapMgr.Loader("#com-content", uri);
}


htapMgr.WorkLoader = function(uri)
{
    htapMgr.Loader("#work-content", uri);
}
