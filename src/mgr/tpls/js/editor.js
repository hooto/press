var l5sEditor = {
    editors: {},
}

l5sEditor.Open = function(name)
{
    seajs.use([
        "~/codemirror/4/codemirror.min.css",
        "~/codemirror/4/codemirror.min.js",
    ],
    function() {
        seajs.use([
            "~/codemirror/4/mode/markdown/markdown.min.js",
            "~/codemirror/4/mode/xml/xml.min.js",
            "~/codemirror/4/addon/selection/active-line.min.js",
            // "~/codemirror/4/addon/scroll/annotatescrollbar.min.js",
        ],
        function() {

            var height = $("#field_"+ name +"_layout").height();

            l5sEditor.editors[name] = CodeMirror.fromTextArea(document.getElementById("field_"+ name), {
                mode            : "markdown",
                lineNumbers     : true,
                theme           : "default",
                lineWrapping    : true,
                styleActiveLine : true,
            });

            l5sEditor.editors[name].setSize("100%", height);

            l5sEditor.editors[name].on("change", function(cm) {
                // console.log("editor changed: "+ name);
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
    // l5sMgr.Ajax("/v1/text/markdown-render", {
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
    delete l5sEditor.editors[name];
}