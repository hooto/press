var l5s = {
    base : "/",
}

l5s.Boot = function()
{
    seajs.config({
        base: l5s.base,
    });

    seajs.use([
        "~/jquery/1.11/jquery.min.js",
    ],
    function() {

        seajs.use([
            // "~/bootstrap/3.3/js/bootstrap.min.js",
            "~/lessui/js/lessui.js",
        ],
        function() {

            for (var i in window.onload_hooks) {
                window.onload_hooks[i]();
            }            
        });
    });
}


l5s.CodeRender = function()
{
    $("code[class^='language-']").each(function(i, el) {

        var lang = el.className.substr("language-".length);

        var modes = [];

        if (lang == "html") {
            lang = "htmlmixed";
        }

        switch (lang) {

        case "php":
            modes.push("~/codemirror/5/mode/php/php.js");
        case "htmlmixed":
            modes.push("~/codemirror/5/mode/xml/xml.js");
            modes.push("~/codemirror/5/mode/javascript/javascript.js");
            modes.push("~/codemirror/5/mode/css/css.js");
            modes.push("~/codemirror/5/mode/htmlmixed/htmlmixed.js");
            break;

        case "c":
        case "cpp":
        case "clike":
        case "java":
            lang = "clike";
            break;

        case "json":
            modes.push("~/codemirror/5/mode/javascript/javascript.js");
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
            modes.push("~/codemirror/5/mode/"+ lang +"/"+ lang +".js");
            break;

        default:
            return;
        }

        seajs.use([
            "~/codemirror/5/lib/codemirror.css",
            "~/codemirror/5/lib/codemirror.js",
        ],
        function() {

            modes.push("~/codemirror/5/addon/runmode/runmode.js");
            modes.push("~/codemirror/5/mode/clike/clike.js");

            seajs.use(modes, function() {

                $(el).addClass('cm-s-default'); // apply a theme class
                CodeMirror.runMode($(el).text(), lang, $(el)[0]);
            });
        });
    });
}

l5s.NavActive = function(tplid, path)
{
    if (!tplid || !path) {
        return;
    }

    var nav = $("#"+ tplid);
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

l5s.Ajax = function(url, options)
{
    options = options || {};

    //
    if (url.substr(0, 1) != "/" && url.substr(0, 4) != "http") {
        url = l5s.base + url;
    }

    //
    if (/\?/.test(url)) {
        url += "&_=";
    } else {
        url += "?_=";
    }
    url += Math.random();

    //
    // url += "&access_token="+ l4iCookie.Get("access_token");

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


l5s.ActionLoader = function(target, uri)
{
    l5s.Ajax(l5s.base + uri, {
        callback: function(err, data) {
            $("#"+ target).html(data);
        }
    });
}


l5s.ApiCmd = function(url, options)
{
    l5s.Ajax(l5s.base + url, options);
}

