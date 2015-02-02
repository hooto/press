var l5sMgr = {
    base : "/mgr/",
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
            "~/lessui/css/lessui.min.css",
            "-/css/main.css?_="+ Math.random(),
            "-/js/spec.js?_="+ Math.random(),
            "-/js/node.js?_="+ Math.random(),
            "-/js/model.js?_="+ Math.random(),
        ], function() {

            l5sMgr.Ajax(l5sMgr.base +"-/body.tpl", {
                callback: function(err, data) {
                
                    $("#body-content").html(data);
                
                    l5sNode.Index();
                }
            });
        });
    });
}

// l5sMgr.ComLoader = function(uri)
// {
//     l5sMgr.Ajax("#com-content", uri);
// }

// l5sMgr.WorkLoader = function(uri)
// {
//     l5sMgr.Ajax("#work-content", uri);
// }

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
    url += "&access_token="+ l4iCookie.Get("access_token");

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
