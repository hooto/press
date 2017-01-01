var htpMgr = {
    frtbase : "/",
    base    : "/mgr/",
    api     : "/v1/",
    basetpl : "/mgr/-/",
}

htpMgr.Boot = function()
{
    if (window._basepath) {
        htpMgr.frtbase = window._basepath;
        if (!htpMgr.frtbase || htpMgr.frtbase == "") {
            htpMgr.frtbase = "/";
        }
        htpMgr.base    = htpMgr.frtbase +"mgr/";
        htpMgr.api     = htpMgr.frtbase +"v1/";
        htpMgr.basetpl = htpMgr.frtbase +"mgr/-/";
    }

    seajs.config({
        base: htpMgr.base,
        alias: {
            ep: '~/lessui/js/eventproxy.js'
        },
    });

    seajs.use([
        "~/lessui/js/browser-detect.js",
        "~/htp/js/jquery.js",
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
            "~/htp/css/main.css?_="+ Math.random(),
            "~/htp/js/marked.js",
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

            setTimeout(htpMgr.BootInit, 300);
        });
    });
}

htpMgr.BootInit = function()
{
    $("#htpm-topbar").css({"display": "block"});

    htpSys.Init();

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
        htpEditor.sizeRefresh();
    });


    htpSys.Init();
    htpSpec.Init();
    htpS2.Init();

    htpNode.Init(function() {

        var navlast = l4iStorage.Get("htpm_nav_last_active");
        if (!navlast) {
            navlast = "sys/index"
        }
        l4i.UrlEventHandler(navlast);
    });
}

htpMgr.HttpSrvBasePath = function(url)
{
    if (htpMgr.base == "") {
        return url;
    }

    if (url.substr(0, 1) == "/") {
        return url;
    }

    return htpMgr.base + url;
}

htpMgr.Ajax = function(url, options)
{
    options = options || {};

    //
    if (url.substr(0, 1) != "/" && url.substr(0, 4) != "http") {
        url = htpMgr.HttpSrvBasePath(url);
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


htpMgr.ApiCmd = function(url, options)
{
    htpMgr.Ajax(htpMgr.api + url, options);
}


htpMgr.TplCmd = function(url, options)
{
    htpMgr.Ajax(htpMgr.basetpl + url +".tpl", options);
}


htpMgr.Loader = function(target, uri)
{
    htpMgr.Ajax(htpMgr.basetpl + uri +".tpl", {
        callback: function(err, data) {
            $(target).html(data);
        }
    });
}


htpMgr.BodyLoader = function(uri)
{
    htpMgr.Loader("#body-content", uri);
}


htpMgr.ComLoader = function(uri)
{
    htpMgr.Loader("#com-content", uri);
}


htpMgr.WorkLoader = function(uri)
{
    htpMgr.Loader("#work-content", uri);
}
