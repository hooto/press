var hpressMgr = {
    frtbase: "/hpress/",
    base: "/hpress/mgr/",
    api: "/hpress/v1/",
    basetpl: "/hpress/~/hpressm/",
    debug: true,
}

hpressMgr.debug_uri = function() {
    if (!hpressMgr.debug) {
        return "";
    }
    return "?_=" + Math.random();
}

hpressMgr.Boot = function() {
    if (window._basepath && window._basepath.length > 1) {
        hpressMgr.frtbase = window._basepath;
        if (hpressMgr.frtbase.substring(hpressMgr.frtbase.length - 1) != "/") {
            hpressMgr.frtbase += "/";
        }
        if (!hpressMgr.frtbase || hpressMgr.frtbase == "") {
            hpressMgr.frtbase = "/";
        }
        hpressMgr.base = hpressMgr.frtbase + "/mgr/";
        hpressMgr.api = hpressMgr.frtbase + "v1/";
        hpressMgr.basetpl = hpressMgr.frtbase + "~/hpressm/";
    }

    seajs.config({
        base: hpressMgr.frtbase,
        alias: {
            ep: '~/lessui/js/eventproxy.js'
        },
    });

    seajs.use([
        "~/lessui/js/browser-detect.js",
        "~/hpress/js/jquery.js",
        "~/lessui/js/eventproxy.js",
    ], function() {

        var browser = BrowserDetect.browser;
        var version = BrowserDetect.version;
        var OS = BrowserDetect.OS;

        if (!((browser == 'Chrome' && version >= 22) ||
            (browser == 'Firefox' && version >= 31.0))) {
            $('body').load(window._basepath + "/error/browser");
            return;
        }

        seajs.use([
            "~/bs/3.3/css/bootstrap.css",
            "~/purecss/pure.css",
            "~/lessui/js/lessui.js",
            "~/lessui/css/lessui.css",
            "~/hpress/css/main.css" + hpressMgr.debug_uri(),
            "~/hpress/js/marked.js",
            "~/hpressm/css/main.css" + hpressMgr.debug_uri(),
            "~/hpressm/css/defx.css",
            "~/hpressm/js/spec.js" + hpressMgr.debug_uri(),
            "~/hpressm/js/spec-editor.js" + hpressMgr.debug_uri(),
            "~/hpressm/js/tablet.js",
            "~/hpressm/js/lc-editor.js",
            "~/hpressm/js/model.js" + hpressMgr.debug_uri(),
            "~/hpressm/js/term.js" + hpressMgr.debug_uri(),
            "~/hpressm/js/node.js" + hpressMgr.debug_uri(),
            "~/hpressm/js/sys.js" + hpressMgr.debug_uri(),
            "~/hpressm/js/s2.js" + hpressMgr.debug_uri(),
            "~/hpressm/js/editor.js" + hpressMgr.debug_uri(),
        ], function() {

            setTimeout(hpressMgr.BootInit, 300);
        });
    });
}

hpressMgr.BootInit = function() {
    l4i.debug = hpressMgr.debug;

    $("#hpressm-topbar").css({
        "display": "block"
    });

    hpressSys.Init();

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
        hpressEditor.sizeRefresh();
    });


    hpressSys.Init();
    hpressSpec.Init();
    hpressS2.Init();

    hpressNode.Init(function() {

        var navlast = l4iStorage.Get("hpressm_nav_last_active");
        if (!navlast) {
            navlast = "sys/index"
        }
        l4i.UrlEventHandler(navlast);
    });
}

hpressMgr.HttpSrvBasePath = function(url) {
    if (hpressMgr.base == "") {
        return url;
    }

    if (url.substr(0, 1) == "/") {
        return url;
    }

    return hpressMgr.base + url;
}

hpressMgr.Ajax = function(url, options) {
    options = options || {};

    //
    if (url.substr(0, 1) != "/" && url.substr(0, 4) != "http") {
        url = hpressMgr.HttpSrvBasePath(url);
    }

    l4i.Ajax(url, options)
}

hpressMgr.AlertUserLogin = function() {
    l4iAlert.Open("warn", "You are not logged in, or your login session has expired. Please sign in again", {
        close: false,
        buttons: [{
            title: "SIGN IN",
            href: hpressMgr.frtbase + "auth/login",
        }],
    });
}

hpressMgr.ApiCmd = function(url, options) {
    if (options.nocache === undefined) {
        options.nocache = true;
    }

    var appcb = null;
    if (options.callback) {
        appcb = options.callback;
    }
    options.callback = function(err, data) {
        if (err == "Unauthorized") {
            return hpressMgr.AlertUserLogin();
        }
        if (appcb) {
            appcb(err, data);
        }
    }

    hpressMgr.Ajax(hpressMgr.api + url, options);
}


hpressMgr.TplCmd = function(url, options) {
    hpressMgr.Ajax(hpressMgr.basetpl + url + ".tpl", options);
}


hpressMgr.Loader = function(target, uri) {
    hpressMgr.Ajax(hpressMgr.basetpl + uri + ".tpl", {
        callback: function(err, data) {
            $(target).html(data);
        }
    });
}


hpressMgr.BodyLoader = function(uri) {
    hpressMgr.Loader("#body-content", uri);
}


hpressMgr.ComLoader = function(uri) {
    hpressMgr.Loader("#com-content", uri);
}


hpressMgr.WorkLoader = function(uri) {
    hpressMgr.Loader("#work-content", uri);
}