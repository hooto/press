// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

var hpressMgr = {
    frtbase: "/hpress/",
    base: "/hpress/mgr/",
    api: "/hpress/v1/",
    basetpl: "/hpress/~/hpressm/",
    sys_version_sign: "1.0",
    debug: true,
}

hpressMgr.urlver = function(debug_off) {
    var u = "?_v="+ hpressMgr.sys_version_sign;
    if (!debug_off && hpressMgr.debug) {
        u += "&_=" + Math.random();
    }
    return u;
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

    if (window._sys_version_sign && window._sys_version_sign.length > 1) {
        hpressMgr.sys_version_sign = window._sys_version_sign;
    }

    seajs.config({
        base: hpressMgr.frtbase,
        alias: {
            ep: '~/lessui/js/eventproxy.js' + hpressMgr.urlver(),
        },
    });

    seajs.use([
        "~/lessui/js/browser-detect.js" + hpressMgr.urlver(),
        "~/hpress/js/jquery.js" + hpressMgr.urlver(),
        "~/lessui/js/eventproxy.js" + hpressMgr.urlver(),
    ], function() {

        var browser = BrowserDetect.browser;
        var version = BrowserDetect.version;
        var OS = BrowserDetect.OS;

        if (!((browser == 'Chrome' && version >= 22) ||
            (browser == 'Firefox' && version >= 31.0))) {
            $('body').load(window._basepath + "/error/browser?" + hpressMgr.urlver());
            return;
        }

        seajs.use([
            "~/bs/3.3/css/bootstrap.css" + hpressMgr.urlver(),
            "~/purecss/pure.css" + hpressMgr.urlver(),
            "~/lessui/js/lessui.js" + hpressMgr.urlver(),
            "~/lessui/css/lessui.css" + hpressMgr.urlver(),
            "~/hpress/css/main.css" + hpressMgr.urlver(),
            "~/hpress/js/marked.js" + hpressMgr.urlver(),
            "~/hpressm/css/main.css" + hpressMgr.urlver(),
            "~/hpressm/css/defx.css" + hpressMgr.urlver(),
            "~/hpressm/js/spec.js" + hpressMgr.urlver(),
            "~/hpressm/js/spec-editor.js" + hpressMgr.urlver(),
            "~/hpressm/js/tablet.js" + hpressMgr.urlver(),
            "~/hpressm/js/lc-editor.js" + hpressMgr.urlver(),
            "~/hpressm/js/model.js" + hpressMgr.urlver(),
            "~/hpressm/js/term.js" + hpressMgr.urlver(),
            "~/hpressm/js/node.js" + hpressMgr.urlver(),
            "~/hpressm/js/sys.js" + hpressMgr.urlver(),
            "~/hpressm/js/s2.js" + hpressMgr.urlver(),
            "~/hpressm/js/editor.js" + hpressMgr.urlver(),
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
    hpressMgr.Ajax(hpressMgr.basetpl + url + ".tpl" + hpressMgr.urlver(true), options);
}


hpressMgr.Loader = function(target, uri) {
    hpressMgr.Ajax(hpressMgr.basetpl + uri + ".tpl" + hpressMgr.urlver(true), {
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
