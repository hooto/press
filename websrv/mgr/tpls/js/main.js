var l5sMgr = {
    base    : "/mgr/",
    api     : "/v1/",
    basetpl : "/mgr/-/",
}

l5sMgr.Boot = function()
{
    seajs.config({
        base: l5sMgr.base,
        alias: {
            ep: '~/lessui/js/eventproxy.js'
        },
    });

    seajs.use([
        "~/lessui/js/BrowserDetect.js",
        "~/jquery/1.11/jquery.min.js",
        "~/lessui/js/eventproxy.js",
    ], function() {

        var browser = BrowserDetect.browser;
        var version = BrowserDetect.version;
        var OS      = BrowserDetect.OS;

        if (!((browser == 'Chrome' && version >= 22)
            || (browser == 'Firefox' && version >= 31.0))) { 
            $('body').load("/error/browser");
            return;
        }

        seajs.use([
            "~/bootstrap/3.3/css/bootstrap.min.css",
            "~/bootstrap/3.3/js/bootstrap.min.js",
            "~/lessui/js/lessui.js",
            "~/lessui/css/lessui.css",
            "~/l5s/css/main.css?_="+ Math.random(),
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
            "-/js/editor.js?_="+ Math.random(),
            "~/l5s/js/marked.min.js",
        ], function() {

            l5sSys.Init();

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
                l5sEditor.sizeRefresh();
            });

            l5sMgr.TplCmd("body", {
                callback: function(err, data) {
                
                    l5sSys.Init();
                    l5sNode.Init();
                    l5sSpec.Init();

                    $("#body-content").html(data);
                
                    l5sNode.Index();
                    // l5sSys.Index();
                    // l5sSpec.Index();
                }
            });
        });
    });
}

l5sMgr.SignOut = function()
{
    // l4iCookie.Del("access_token");

    // console.log("access_token");
    // return;

    l5sMgr.ApiCmd("ids/clean-cookie", {
        success: function() {

        },
    });
}

l5sMgr.Ajax = function(url, options)
{
    options = options || {};

    //
    if (url.substr(0, 1) != "/" && url.substr(0, 4) != "http") {
        url = l5sMgr.base + url;
    }

    //
    if (/\?/.test(url)) {
        url += "&_=";
    } else {
        url += "?_=";
    }
    url += Math.random();

    //
    if (l4iCookie.Get("access_token")) {
        url += "&access_token="+ l4iCookie.Get("access_token");
    }

    //
    if (options.method === undefined) {
        options.method = "GET";
    }

    //
    if (options.timeout === undefined) {
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
            // console.log(xhr.responseText);
            if (typeof options.callback === "function") {
                options.callback(xhr.responseText, null);
            }
            if (typeof options.error === "function") {
                options.error(xhr, textStatus, error);
            }
        }
    });
}



l5sMgr.ApiCmd = function(url, options)
{
    l5sMgr.Ajax(l5sMgr.api + url, options);
}

l5sMgr.TplCmd = function(url, options)
{
    l5sMgr.Ajax(l5sMgr.basetpl + url +".tpl", options);
}

l5sMgr.Loader = function(target, uri)
{
    l5sMgr.Ajax(l5sMgr.basetpl + uri +".tpl", {
        callback: function(err, data) {
            $(target).html(data);
        }
    });
}

l5sMgr.BodyLoader = function(uri)
{
    l5sMgr.Loader("#body-content", uri);
}

l5sMgr.ComLoader = function(uri)
{
    l5sMgr.Loader("#com-content", uri);
}

l5sMgr.WorkLoader = function(uri)
{
    l5sMgr.Loader("#work-content", uri);
}


