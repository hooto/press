var htpEditor = {
    editors: {},
}

htpEditor.Open = function(name, format) {
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

            var editor = htpEditor.editors[name];
            if (editor) {

                $("#field_" + name + "_attr_format").val(format);

                htpEditor.editors[name].setOption("lineNumbers", lineNumbers);

                if (format != "md") {
                    htpEditor.PreviewClose(name);
                }

                return;
            }

            var height = $("#field_" + name + "_layout").height(),
                width = $("#field_" + name + "_layout").width();

            htpEditor.editors[name] = CodeMirror.fromTextArea(document.getElementById("field_" + name), {
                mode: "markdown",
                lineNumbers: lineNumbers,
                theme: "default",
                lineWrapping: true,
                styleActiveLine: true,
            // viewportMargin  : Infinity,
            });

            htpEditor.editors[name].setSize(width, height);

            htpEditor.editors[name].on("change", function(cm) {
                htpEditor.previewChanged(name);
            });

            // console.log("post width: "+ $("#field_"+ name +"_layout").width());
            $("#field_" + name + "_layout").addClass("htpm-editor-cm");
            $("#field_" + name + "_tools").find(".preview_open").show();
        });
    });
}

htpEditor.sizeRefresh = function() {
    // console.log("htpEditor.sizeRefresh");

    var dels = [];
    for (var name in htpEditor.editors) {

        var ok = document.getElementById("field_" + name + "_layout");
        if (!ok) {
            dels.push(name);
            continue;
        }

        if (!$("#field_" + name + "_colpreview").is(":visible")) {
            // console.log("preview skip");
            continue;
        }

        htpEditor.PreviewOpen(name);
    }

    for (var i in dels) {
        delete htpEditor.editors[dels[i]];
    }
}

htpEditor.PreviewOpen = function(name) {
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
        htpEditor.editorBindScroll(name);
    }, function() {
        htpEditor.editorUnBindScroll(name);
    });

    $("#field_" + name + "_preview").hover(function() {
        htpEditor.previewBindScroll(name);
    }, function() {
        htpEditor.previewUnBindScroll(name);
    });

    htpEditor.previewChanged(name);

    $("#field_" + name + "_tools").find(".preview_close").show();
    $("#field_" + name + "_tools").find(".preview_open").hide();
}

htpEditor.PreviewClose = function(name) {
    htpEditor.editorUnBindScroll(name);
    htpEditor.previewUnBindScroll(name);

    $("#field_" + name + "_preview").empty();

    $("#field_" + name + "_colspace").hide();
    $("#field_" + name + "_colpreview").hide();

    $("#field_" + name + "_editor").css({
        width: "100%"
    });

    $("#field_" + name + "_tools").find(".preview_close").hide();
    $("#field_" + name + "_tools").find(".preview_open").show();
}

htpEditor.previewChanged = function(name) {
    if (!$("#field_" + name + "_colpreview").is(":visible")) {
        return;
    }

    var editor = htpEditor.editors[name];
    if (!editor) {
        return;
    }

    var text = editor.getValue();

    // frountend markdown render
    $("#field_" + name + "_preview").html(marked(text));

// // backend markdown render
// htpMgr.ApiCmd("/text/markdown-render", {
//     method : "POST",
//     data   : text,
//     callback : function(err, data) {
//         $("#field_"+ name +"_preview").html(data);
//     }
// });
}

htpEditor.editorBindScroll = function(name) {
    htpEditor.previewUnBindScroll(name);

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

htpEditor.editorUnBindScroll = function(name) {
    $("#field_" + name + "_editor").find(".CodeMirror-scroll").unbind("scroll");
}

htpEditor.previewBindScroll = function(name) {
    htpEditor.editorUnBindScroll(name);

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

htpEditor.previewUnBindScroll = function(name) {
    $("#field_" + name + "_preview").unbind("scroll");
}

htpEditor.Content = function(name) {
    var editor = htpEditor.editors[name];
    if (editor) {
        return editor.getValue();
    }

    return null;
}

htpEditor.Close = function(name) {
    var edr = htpEditor.editors[name];
    if (edr) {
        htpEditor.editors[name] = null;
        delete htpEditor.editors[name];
    }
}

htpEditor.Clean = function() {
    for (var i in htpEditor.editors) {
        htpEditor.editors[i] = null;
        delete htpEditor.editors[i];
    }
}
