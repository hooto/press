var htapEditor = {
    editors: {},
}

htapEditor.Open = function(name, format)
{
    seajs.use([
        "~/cm/5/lib/codemirror.css",
        "~/cm/5/lib/codemirror.js",
    ],
    function() {

        seajs.use([
            "~/cm/5/mode/markdown/markdown.js",
            "~/cm/5/mode/xml/xml.js",
            "~/cm/5/addon/selection/active-line.js",
        ],
        function() {

            var lineNumbers = false;
            if (format == "md" || format == "html") {
                lineNumbers = true;
            }

            if (format == "md") {
                $("#field_"+ name +"_editor_mdr").show();
            } else {
                $("#field_"+ name +"_editor_mdr").hide();
            }

            $("#field_"+ name +"_editor_nav").find("a.active").removeClass("active");
            $("#field_"+ name +"_editor_nav").find("a.editor-nav-"+ format).addClass("active");

            var editor = htapEditor.editors[name];
            if (editor) {

                $("#field_"+ name +"_attr_format").val(format);

                htapEditor.editors[name].setOption("lineNumbers", lineNumbers);

                if (format != "md") {
                    htapEditor.PreviewClose(name);
                }

                return;
            }

            var height = $("#field_"+ name +"_layout").height(),
                width = $("#field_"+ name +"_layout").width();

            htapEditor.editors[name] = CodeMirror.fromTextArea(document.getElementById("field_"+ name), {
                mode            : "markdown",
                lineNumbers     : lineNumbers,
                theme           : "default",
                lineWrapping    : true,
                styleActiveLine : true,
                // viewportMargin  : Infinity,
            });

            htapEditor.editors[name].setSize(width, height);

            htapEditor.editors[name].on("change", function(cm) {
                htapEditor.previewChanged(name);
            });

            // console.log("post width: "+ $("#field_"+ name +"_layout").width());
            $("#field_"+ name +"_layout").addClass("htapm-editor-cm");
            $("#field_"+ name +"_tools").find(".preview_open").show();
        });
    });
}

htapEditor.sizeRefresh = function()
{
    // console.log("htapEditor.sizeRefresh");

    var dels = [];
    for (var name in htapEditor.editors) {

        var ok = document.getElementById("field_"+ name +"_layout");
        if (!ok) {
            dels.push(name);
            continue;
        }

        if (!$("#field_"+ name +"_colpreview").is(":visible")) {
            // console.log("preview skip");
            continue;
        }

        htapEditor.PreviewOpen(name);
    }

    for (var i in dels) {
        delete htapEditor.editors[dels[i]];
    }
}

htapEditor.PreviewOpen = function(name)
{
    // console.log($("#field_"+ name +"_preview"));

    var width = $("#field_"+ name +"_layout").width(),
        height = $("#field_"+ name +"_layout").height();

    var width5 = (width - 20) / 2;

    $("#field_"+ name +"_editor").css({width: width5 +"px"});
    $("#field_"+ name +"_editor").find(".CodeMirror").css({width: width5 +"px"});
    $("#field_"+ name +"_colpreview").css({width: width5 +"px", height: height +"px"});
    $("#field_"+ name +"_preview").css({width: width5 +"px", height: height +"px"});

    $("#field_"+ name +"_colspace").show();
    $("#field_"+ name +"_colpreview").show();

    $("#field_"+ name +"_editor").find(".CodeMirror").hover(function() {
        htapEditor.editorBindScroll(name);
    }, function() {
        htapEditor.editorUnBindScroll(name);
    });

    $("#field_"+ name +"_preview").hover(function() {
        htapEditor.previewBindScroll(name);
    }, function() {
        htapEditor.previewUnBindScroll(name);
    });

    htapEditor.previewChanged(name);

    $("#field_"+ name +"_tools").find(".preview_close").show();
    $("#field_"+ name +"_tools").find(".preview_open").hide();
}

htapEditor.PreviewClose = function(name)
{
    htapEditor.editorUnBindScroll(name);
    htapEditor.previewUnBindScroll(name);

    $("#field_"+ name +"_preview").empty();

    $("#field_"+ name +"_colspace").hide();
    $("#field_"+ name +"_colpreview").hide();

    $("#field_"+ name +"_editor").css({width: "100%"});

    $("#field_"+ name +"_tools").find(".preview_close").hide();
    $("#field_"+ name +"_tools").find(".preview_open").show();
}

htapEditor.previewChanged = function(name)
{
    if (!$("#field_"+ name +"_colpreview").is(":visible")) {
        return;
    }

    var editor = htapEditor.editors[name];
    if (!editor) {
        return;
    }

    var text = editor.getValue();

    // frountend markdown render
    $("#field_"+ name +"_preview").html(marked(text));

    // // backend markdown render
    // htapMgr.ApiCmd("/text/markdown-render", {
    //     method : "POST",
    //     data   : text,
    //     callback : function(err, data) {
    //         $("#field_"+ name +"_preview").html(data);
    //     }
    // });
}

htapEditor.editorBindScroll = function(name)
{
    htapEditor.previewUnBindScroll(name);

    $("#field_"+ name +"_editor").find(".CodeMirror-scroll").on("scroll", function() {

        var height    = $(this).outerHeight();
        var scrollTop = $(this).scrollTop();
        var percent   = (scrollTop / $(this)[0].scrollHeight);
        var preview   = $("#field_"+ name +"_preview");

        if (scrollTop === 0) {
            preview.scrollTop(0);
        } else if (scrollTop + height >= $(this)[0].scrollHeight) {
            preview.scrollTop(preview[0].scrollHeight);
        } else {
            preview.scrollTop(preview[0].scrollHeight * percent);
        }
    });
}

htapEditor.editorUnBindScroll = function(name)
{
    $("#field_"+ name +"_editor").find(".CodeMirror-scroll").unbind("scroll");
}

htapEditor.previewBindScroll = function(name)
{
    htapEditor.editorUnBindScroll(name);

    $("#field_"+ name +"_preview").on("scroll", function() {

        var height     = $(this).outerHeight();
        var scrollTop  = $(this).scrollTop();
        var percent    = (scrollTop / $(this)[0].scrollHeight);
        var editorView = $("#field_"+ name +"_editor").find(".CodeMirror-scroll");

        if (scrollTop === 0) {
            editorView.scrollTop(0);
        } else if (scrollTop + height >= $(this)[0].scrollHeight) {
            editorView.scrollTop(editorView[0].scrollHeight);
        } else {
            editorView.scrollTop(editorView[0].scrollHeight * percent);
        }
    });
}

htapEditor.previewUnBindScroll = function(name)
{
    $("#field_"+ name +"_preview").unbind("scroll");
}

htapEditor.Content = function(name)
{
    var editor = htapEditor.editors[name];
    if (editor) {
        return editor.getValue();
    }

    return null;
}

htapEditor.Close = function(name)
{
    var edr = htapEditor.editors[name];
    if (edr) {
        htapEditor.editors[name] = null;
        delete htapEditor.editors[name];
    }
}

htapEditor.Clean = function()
{
    for (var i in htapEditor.editors) {
        htapEditor.editors[i] = null;
        delete htapEditor.editors[i];
    }
}
