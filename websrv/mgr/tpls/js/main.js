var htpMgr = {
    frtbase : "/",
    base    : "/mgr/",
    api     : "/v1/",
    basetpl : "/mgr/-/",
    debug   : true,
}

htpMgr.debug_uri = function()
{
    if (!htpMgr.debug) {
        return "";
    }
    return "?_="+ Math.random();
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
            "~/htp/css/main.css"+ htpMgr.debug_uri(),
            "~/htp/js/marked.js",
            "-/css/main.css"+ htpMgr.debug_uri(),
            "-/css/defx.css",
            "-/js/spec.js"+ htpMgr.debug_uri(),
            "-/js/spec-editor.js"+ htpMgr.debug_uri(),
            "-/js/tablet.js",
            "-/js/lc-editor.js",
            "-/js/model.js"+ htpMgr.debug_uri(),
            "-/js/term.js"+ htpMgr.debug_uri(),
            "-/js/node.js"+ htpMgr.debug_uri(),
            "-/js/sys.js"+ htpMgr.debug_uri(),
            "-/js/s2.js"+ htpMgr.debug_uri(),
            "-/js/editor.js"+ htpMgr.debug_uri(),
        ], function() {

            setTimeout(htpMgr.BootInit, 300);
        });
    });
}

htpMgr.BootInit = function()
{
    l4i.debug = htpMgr.debug;

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

    l4i.Ajax(url, options)
}

htpMgr.AlertUserLogin = function()
{
    l4iAlert.Open("warn", "You are not logged in, or your login session has expired. Please sign in again", {
        close: false,
        buttons: [{
            title: "SIGN IN",
            href: htpMgr.frtbase +"auth/login",
        }],
    });
}

htpMgr.ApiCmd = function(url, options)
{
    if (options.nocache === undefined) {
        options.nocache = true;
    }

    var appcb = null;
    if (options.callback) {
        appcb = options.callback;
    }
    options.callback = function(err, data) {
        if (err == "Unauthorized") {
            return htpMgr.AlertUserLogin();
        }
        if (appcb) {
            appcb(err, data);
        }
    }

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
