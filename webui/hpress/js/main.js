var hpress = {
    base: "/hpress/",
}

hpress.Boot = function() {
    if (window._basepath && window._basepath.length > 1) {
        hpress.base = window._basepath;
        if (hpress.base.substring(hpress.base.length - 1) != "/") {
            hpress.base += "/";
        }
    }

    if (!hpress.base || hpress.base == "") {
        hpress.base = "/";
    }

    seajs.config({
        base: hpress.base,
    });

    seajs.use([
        "~/hpress/js/jquery.js",
    ], function() {

        seajs.use([
            "~/lessui/js/lessui.js",
        ], function() {

            setTimeout(function() {
                for (var i in window.onload_hooks) {
                    window.onload_hooks[i]();
                }
            }, 100);
        });
    });
}

hpress.HttpSrvBasePath = function(url) {
    if (hpress.base == "") {
        return url;
    }

    if (url.substr(0, 1) == "/") {
        return url;
    }

    return hpress.base + url;
}

hpress.CodeRender = function() {
    $("code[class^='language-']").each(function(i, el) {

        var lang = el.className.substr("language-".length);
        if (lang == "hchart" || lang == "hooto_chart") {
            return hpress.hchartRender(i, el);
        }

        var modes = [];

        if (lang == "html") {
            lang = "htmlmixed";
        }

        switch (lang) {

            case "php":
                modes.push("~/cm/5/mode/php/php.js");
            case "htmlmixed":
                modes.push("~/cm/5/mode/xml/xml.js");
                modes.push("~/cm/5/mode/javascript/javascript.js");
                modes.push("~/cm/5/mode/css/css.js");
                modes.push("~/cm/5/mode/htmlmixed/htmlmixed.js");
                break;

            case "c":
            case "cpp":
            case "clike":
            case "java":
                lang = "clike";
                break;

            case "json":
                modes.push("~/cm/5/mode/javascript/javascript.js");
                lang = "application/ld+json";
                break;

            case "go":
            case "javascript":
            case "css":
            case "xml":
            case "yaml":
            case "lua":
            case "markdown":
            case "r":
            case "shell":
            case "sql":
            case "swift":
            case "erlang":
            case "nginx":
                modes.push("~/cm/5/mode/" + lang + "/" + lang + ".js");
                break;

            default:
                return;
        }

        seajs.use([
            "~/cm/5/lib/codemirror.css",
            "~/cm/5/lib/codemirror.js",
        ], function() {

            modes.push("~/cm/5/addon/runmode/runmode.js");
            modes.push("~/cm/5/mode/clike/clike.js");

            seajs.use(modes, function() {

                $(el).addClass('cm-s-default'); // apply a theme class
                CodeMirror.runMode($(el).text(), lang, $(el)[0]);
            });
        });
    });
}

hpress.hchartRender = function(i, elem) {
    var elem_id = "hchart-id-" + i;
    elem.setAttribute("id", elem_id);
    seajs.use([
        "~/hchart/hchart.js",
    ], function() {
        hooto_chart.basepath = hpress.base + "/~/hchart";
        hooto_chart.opts_width = "600px";
        hooto_chart.opts_height = "400px";
        hooto_chart.JsonRenderElement(elem, elem_id);
    });
}

hpress.NavActive = function(tplid, path) {
    if (!tplid || !path) {
        return;
    }

    var nav = $("#" + tplid);
    nav.find("a").each(function() {

        var href = $(this).attr("href");

        if (href) {

            if (href.match(path)) {
                nav.find("a.active").removeClass("active");
                $(this).addClass("active");
            }
        }
    });
}

hpress.Ajax = function(url, options) {
    options = options || {};

    //
    if (url.substr(0, 1) != "/" && url.substr(0, 4) != "http") {
        url = hpress.HttpSrvBasePath(url);
    }

    l4i.Ajax(url, options)
}

hpress.ActionLoader = function(target, url) {
    hpress.Ajax(hpress.HttpSrvBasePath(url), {
        callback: function(err, data) {
            $("#" + target).html(data);
        }
    });
}

hpress.ApiCmd = function(url, options) {
    hpress.Ajax(hpress.HttpSrvBasePath(url), options);
}

hpress.AuthSessionRefresh = function() {
    hpress.Ajax(hpress.HttpSrvBasePath("auth/session"), {
        callback: function(err, data) {

            if (err || !data || data.kind != "AuthSession") {

                return l4iTemplate.Render({
                    dstid: "hpress-topbar-userbar",
                    tplid: "hpress-topbar-user-unsigned-tpl",
                });
            }

            l4iTemplate.Render({
                dstid: "hpress-topbar-userbar",
                tplid: "hpress-topbar-user-signed-tpl",
                data: data,
                success: function() {

                    $("#hpress-topbar-userbar").hover(function() {
                        $("#hpress-topbar-user-signed-modal").fadeIn(200);
                    }, function() {}
                    );
                    $("#hpress-topbar-user-signed-modal").hover(function() {}, function() {
                        $("#hpress-topbar-user-signed-modal").fadeOut(200);
                    }
                    );
                },
            });
        },
    });
}

