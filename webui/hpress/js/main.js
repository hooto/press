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

var hpress = {
    base: "/hpress/",
    sys_version_sign: "1.0",
    debug: true,
}

hpress.urlver = function(debug_off) {
    var u = "?_v=" + hpress.sys_version_sign;
    if (!debug_off && hpress.debug) {
        u += "&_=" + Math.random();
    }
    return u;
}

hpress.Boot = function() {
    if (window._basepath && window._basepath.length > 1) {
        hpress.base = window._basepath;
        if (hpress.base.substring(hpress.base.length - 1) != "/") {
            hpress.base += "/";
        }
    }
    if (window._sys_version_sign && window._sys_version_sign.length > 1) {
        hpress.sys_version_sign = window._sys_version_sign;
    }

    if (!hpress.base || hpress.base == "") {
        hpress.base = "/";
    }

    seajs.config({
        base: hpress.base,
    });

    seajs.use([
        "~/hpress/js/jquery.js" + hpress.urlver(),
    ], function() {

        seajs.use([
            "~/lessui/js/lessui.js" + hpress.urlver(),
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
                modes.push("~/cm/5/mode/php/php.js" + hpress.urlver());
            case "htmlmixed":
                modes.push("~/cm/5/mode/xml/xml.js" + hpress.urlver());
                modes.push("~/cm/5/mode/javascript/javascript.js" + hpress.urlver());
                modes.push("~/cm/5/mode/css/css.js" + hpress.urlver());
                modes.push("~/cm/5/mode/htmlmixed/htmlmixed.js" + hpress.urlver());
                break;

            case "c":
            case "cpp":
            case "clike":
            case "java":
                lang = "clike";
                break;

            case "json":
                modes.push("~/cm/5/mode/javascript/javascript.js" + hpress.urlver());
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
                modes.push("~/cm/5/mode/" + lang + "/" + lang + ".js" + hpress.urlver());
                break;

            default:
                return;
        }

        seajs.use([
            "~/cm/5/lib/codemirror.css" + hpress.urlver(),
            "~/cm/5/lib/codemirror.js" + hpress.urlver(),
        ], function() {

            modes.push("~/cm/5/addon/runmode/runmode.js" + hpress.urlver());
            modes.push("~/cm/5/mode/clike/clike.js" + hpress.urlver());

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
        "~/hchart/hchart.js" + hpress.urlver(),
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
    if (!nav) {
        return;
    }

    var nav_path = window.location.pathname;
    if (!nav_path || nav_path == "") {
        nav_path = "/";
    }

    var found = false;
    while (true) {

        nav.find("a").each(function() {
            if (found) {
                return;
            }
            var href = $(this).attr("href");
            if (href && href == nav_path) {
                nav.find("a.active").removeClass("active");
                $(this).addClass("active");
                found = true;
            }
        });

        if (found) {
            break;
        }

        if (nav_path.lastIndexOf("/") > 0) {
            nav_path = nav_path.substr(0, nav_path.lastIndexOf("/"));
        } else {
            console.log("break 2");
            break;
        }
    }
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

            if (hpress.sys_version_sign == "unreg") {
                return window.location = "/hpress/mgr";
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

