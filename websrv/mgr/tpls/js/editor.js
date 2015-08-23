var l5sEditor = {
    editors: {},
}

l5sEditor.Open = function(name, format)
{
    seajs.use([
        "~/codemirror/5/lib/codemirror.css",
        "~/codemirror/5/lib/codemirror.js",
    ],
    function() {

        seajs.use([
            "~/codemirror/5/mode/markdown/markdown.js",
            "~/codemirror/5/mode/xml/xml.js",
            "~/codemirror/5/addon/selection/active-line.js",
        ],
        function() {

            // console.log(format);

            var lineNumbers = false;
            if (format == "md") {
                lineNumbers = true;
                $("#field_"+ name +"_editor_mdr").show();
            } else {
                $("#field_"+ name +"_editor_mdr").hide();
            }
            
            $("#field_"+ name +"_editor_nav").find("a.active").removeClass("active");
            $("#field_"+ name +"_editor_nav").find("a.editor-nav-"+ format).addClass("active");

            var editor = l5sEditor.editors[name];
            if (editor) {

                $("#field_"+ name +"_attr_format").val(format);
   
                l5sEditor.editors[name].setOption("lineNumbers", lineNumbers);

                if (format != "md") {
                    l5sEditor.PreviewClose(name);                    
                }

                return;
            }

            var height = $("#field_"+ name +"_layout").height();

            l5sEditor.editors[name] = CodeMirror.fromTextArea(document.getElementById("field_"+ name), {
                mode            : "markdown",
                lineNumbers     : lineNumbers,
                theme           : "default",
                lineWrapping    : true,
                styleActiveLine : true,
            });

            l5sEditor.editors[name].setSize("100%", height);

            l5sEditor.editors[name].on("change", function(cm) {
                l5sEditor.previewChanged(name);
            });

            $("#field_"+ name +"_layout").addClass("l5smgr-editor-cm");
            $("#field_"+ name +"_tools").find(".preview_open").show();
        });
    });
}

l5sEditor.sizeRefresh = function()
{
    // console.log("l5sEditor.sizeRefresh");

    var dels = [];
    for (var name in l5sEditor.editors) {

        var ok = document.getElementById("field_"+ name +"_layout");
        if (!ok) {
            dels.push(name);
            continue;
        }

        if (!$("#field_"+ name +"_colpreview").is(":visible")) {
            // console.log("preview skip");
            continue;
        }

        l5sEditor.PreviewOpen(name);
    }

    for (var i in dels) {
        delete l5sEditor.editors[dels[i]];
    }
}

l5sEditor.PreviewOpen = function(name)
{
    // console.log($("#field_"+ name +"_preview"));

    var width = $("#field_"+ name +"_layout").width(),
        height = $("#field_"+ name +"_layout").height();

    var width5 = (width - 20) / 2;

    $("#field_"+ name +"_editor").css({width: width5 +"px"});    
    $("#field_"+ name +"_colpreview").css({width: width5 +"px", height: height +"px"});
    $("#field_"+ name +"_preview").css({width: width5 +"px", height: height +"px"});

    $("#field_"+ name +"_colspace").show();
    $("#field_"+ name +"_colpreview").show();

    $("#field_"+ name +"_editor").find(".CodeMirror").hover(function() {
        l5sEditor.editorBindScroll(name);
    }, function() {
        l5sEditor.editorUnBindScroll(name);
    });

    $("#field_"+ name +"_preview").hover(function() {
        l5sEditor.previewBindScroll(name);
    }, function() {
        l5sEditor.previewUnBindScroll(name);
    });

    l5sEditor.previewChanged(name);
    
    $("#field_"+ name +"_tools").find(".preview_close").show();
    $("#field_"+ name +"_tools").find(".preview_open").hide();
}

l5sEditor.PreviewClose = function(name)
{
    l5sEditor.editorUnBindScroll(name);
    l5sEditor.previewUnBindScroll(name);

    $("#field_"+ name +"_preview").empty();

    $("#field_"+ name +"_colspace").hide();
    $("#field_"+ name +"_colpreview").hide();

    $("#field_"+ name +"_editor").css({width: "100%"});

    $("#field_"+ name +"_tools").find(".preview_close").hide();
    $("#field_"+ name +"_tools").find(".preview_open").show();
}

l5sEditor.previewChanged = function(name)
{
    if (!$("#field_"+ name +"_colpreview").is(":visible")) {
        return;
    }

    var editor = l5sEditor.editors[name];
    if (!editor) {
        return;
    }

    var text = editor.getValue();

    // frountend markdown render
    $("#field_"+ name +"_preview").html(marked(text));

    // // backend markdown render
    // l5sMgr.ApiCmd("/text/markdown-render", {
    //     method : "POST",
    //     data   : text,
    //     callback : function(err, data) {
    //         $("#field_"+ name +"_preview").html(data);
    //     }
    // });
}

l5sEditor.editorBindScroll = function(name)
{
    l5sEditor.previewUnBindScroll(name);

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

l5sEditor.editorUnBindScroll = function(name)
{
    $("#field_"+ name +"_editor").find(".CodeMirror-scroll").unbind("scroll");
}

l5sEditor.previewBindScroll = function(name)
{
    l5sEditor.editorUnBindScroll(name);

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

l5sEditor.previewUnBindScroll = function(name)
{
    $("#field_"+ name +"_preview").unbind("scroll");
}

l5sEditor.Content = function(name)
{
    var editor = l5sEditor.editors[name];
    if (editor) {
        return editor.getValue();
    }

    return null;
}

l5sEditor.Close = function(name)
{
    var edr = l5sEditor.editors[name];
    if (edr) {
        l5sEditor.editors[name] = null;
        delete l5sEditor.editors[name];
    }
}

l5sEditor.Clean = function()
{
    for (var i in l5sEditor.editors) {
        l5sEditor.editors[i] = null;
        delete l5sEditor.editors[i];
    }
}