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
        "~/bootstrap/3.3/js/bootstrap.min.js"
    ], function() {    

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

            case "go":
            case "javascript":
            case "css":
            case "json":
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
            ], function() {

                modes.push("~/codemirror/5/addon/runmode/runmode.js");
                modes.push("~/codemirror/5/mode/clike/clike.js");

                seajs.use(modes, function() {

                    $(el).addClass('cm-s-default'); // apply a theme class
                    CodeMirror.runMode($(el).text(), lang, $(el)[0]);
                });
            });
        });    
    });
}
