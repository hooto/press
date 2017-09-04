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

var hpressEditor = {
    editors: {},
}

hpressEditor.Open = function(name, format) {
    seajs.use([
        "~/cm/5/lib/codemirror.css",
        "~/cm/5/lib/codemirror.js",
    ], function() {

        seajs.use([
            "~/cm/5/mode/markdown/markdown.js",
            "~/cm/5/mode/xml/xml.js",
            "~/cm/5/addon/selection/active-line.js",
        ], function() {

            var lineNumbers = false;
            if (format == "md" || format == "html") {
                lineNumbers = true;
            }

            if (format == "md") {
                $("#field_" + name + "_editor_mdr").show();
            } else {
                $("#field_" + name + "_editor_mdr").hide();
            }

            $("#field_" + name + "_editor_nav").find("a.active").removeClass("active");
            $("#field_" + name + "_editor_nav").find("a.editor-nav-" + format).addClass("active");

            var editor = hpressEditor.editors[name];
            if (editor) {

                $("#field_" + name + "_attr_format").val(format);

                hpressEditor.editors[name].setOption("lineNumbers", lineNumbers);

                if (format != "md") {
                    hpressEditor.PreviewClose(name);
                }

                return;
            }

            var height = $("#field_" + name + "_layout").height(),
                width = $("#field_" + name + "_layout").width();

            hpressEditor.editors[name] = CodeMirror.fromTextArea(document.getElementById("field_" + name), {
                mode: "markdown",
                lineNumbers: lineNumbers,
                theme: "default",
                lineWrapping: true,
                styleActiveLine: true,
            // viewportMargin  : Infinity,
            });

            hpressEditor.editors[name].setSize(width, height);

            hpressEditor.editors[name].on("change", function(cm) {
                hpressEditor.previewChanged(name);
            });

            // console.log("post width: "+ $("#field_"+ name +"_layout").width());
            $("#field_" + name + "_layout").addClass("hpressm-editor-cm");
            $("#field_" + name + "_tools").find(".preview_open").show();
        });
    });
}

hpressEditor.sizeRefresh = function() {
    // console.log("hpressEditor.sizeRefresh");

    var dels = [];
    for (var name in hpressEditor.editors) {

        var ok = document.getElementById("field_" + name + "_layout");
        if (!ok) {
            dels.push(name);
            continue;
        }

        if (!$("#field_" + name + "_colpreview").is(":visible")) {
            // console.log("preview skip");
            continue;
        }

        hpressEditor.PreviewOpen(name);
    }

    for (var i in dels) {
        delete hpressEditor.editors[dels[i]];
    }
}

hpressEditor.PreviewOpen = function(name) {
    // console.log($("#field_"+ name +"_preview"));

    var width = $("#field_" + name + "_layout").width(),
        height = $("#field_" + name + "_layout").height();

    var width5 = (width - 20) / 2;

    $("#field_" + name + "_editor").css({
        width: width5 + "px"
    });
    $("#field_" + name + "_editor").find(".CodeMirror").css({
        width: width5 + "px"
    });
    $("#field_" + name + "_colpreview").css({
        width: width5 + "px",
        height: height + "px"
    });
    $("#field_" + name + "_preview").css({
        width: width5 + "px",
        height: height + "px"
    });

    $("#field_" + name + "_colspace").show();
    $("#field_" + name + "_colpreview").show();

    $("#field_" + name + "_editor").find(".CodeMirror").hover(function() {
        hpressEditor.editorBindScroll(name);
    }, function() {
        hpressEditor.editorUnBindScroll(name);
    });

    $("#field_" + name + "_preview").hover(function() {
        hpressEditor.previewBindScroll(name);
    }, function() {
        hpressEditor.previewUnBindScroll(name);
    });

    hpressEditor.previewChanged(name);

    $("#field_" + name + "_tools").find(".preview_close").show();
    $("#field_" + name + "_tools").find(".preview_open").hide();
}

hpressEditor.PreviewClose = function(name) {
    hpressEditor.editorUnBindScroll(name);
    hpressEditor.previewUnBindScroll(name);

    $("#field_" + name + "_preview").empty();

    $("#field_" + name + "_colspace").hide();
    $("#field_" + name + "_colpreview").hide();

    $("#field_" + name + "_editor").css({
        width: "100%"
    });

    $("#field_" + name + "_tools").find(".preview_close").hide();
    $("#field_" + name + "_tools").find(".preview_open").show();
}

hpressEditor.previewChanged = function(name) {
    if (!$("#field_" + name + "_colpreview").is(":visible")) {
        return;
    }

    var editor = hpressEditor.editors[name];
    if (!editor) {
        return;
    }

    var text = editor.getValue();

    // frountend markdown render
    $("#field_" + name + "_preview").html(marked(text));

// // backend markdown render
// hpressMgr.ApiCmd("/text/markdown-render", {
//     method : "POST",
//     data   : text,
//     callback : function(err, data) {
//         $("#field_"+ name +"_preview").html(data);
//     }
// });
}

hpressEditor.editorBindScroll = function(name) {
    hpressEditor.previewUnBindScroll(name);

    $("#field_" + name + "_editor").find(".CodeMirror-scroll").on("scroll", function() {

        var height = $(this).outerHeight();
        var scrollTop = $(this).scrollTop();
        var percent = (scrollTop / $(this)[0].scrollHeight);
        var preview = $("#field_" + name + "_preview");

        if (scrollTop === 0) {
            preview.scrollTop(0);
        } else if (scrollTop + height >= $(this)[0].scrollHeight) {
            preview.scrollTop(preview[0].scrollHeight);
        } else {
            preview.scrollTop(preview[0].scrollHeight * percent);
        }
    });
}

hpressEditor.editorUnBindScroll = function(name) {
    $("#field_" + name + "_editor").find(".CodeMirror-scroll").unbind("scroll");
}

hpressEditor.previewBindScroll = function(name) {
    hpressEditor.editorUnBindScroll(name);

    $("#field_" + name + "_preview").on("scroll", function() {

        var height = $(this).outerHeight();
        var scrollTop = $(this).scrollTop();
        var percent = (scrollTop / $(this)[0].scrollHeight);
        var editorView = $("#field_" + name + "_editor").find(".CodeMirror-scroll");

        if (scrollTop === 0) {
            editorView.scrollTop(0);
        } else if (scrollTop + height >= $(this)[0].scrollHeight) {
            editorView.scrollTop(editorView[0].scrollHeight);
        } else {
            editorView.scrollTop(editorView[0].scrollHeight * percent);
        }
    });
}

hpressEditor.previewUnBindScroll = function(name) {
    $("#field_" + name + "_preview").unbind("scroll");
}

hpressEditor.Content = function(name) {
    var editor = hpressEditor.editors[name];
    if (editor) {
        return editor.getValue();
    }

    return null;
}

hpressEditor.Close = function(name) {
    var edr = hpressEditor.editors[name];
    if (edr) {
        hpressEditor.editors[name] = null;
        delete hpressEditor.editors[name];
    }
}

hpressEditor.Clean = function() {
    for (var i in hpressEditor.editors) {
        hpressEditor.editors[i] = null;
        delete hpressEditor.editors[i];
    }
}
