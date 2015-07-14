var l5sComment = {

}

l5sComment.EmbedCommit = function()
{
    var form = $("#l5s-comment-embed-new-form");

    var req = {
        pid               : "",
        refer_id          : form.find("input[name=refer_id]").val(),
        refer_specid      : form.find("input[name=refer_specid]").val(),
        refer_datax_table : form.find("input[name=refer_datax_table]").val(),
        author            : form.find("input[name=author]").val(),
        content           : form.find("textarea[name=content]").val(),
    };

    // console.log(req);
    l5s.ApiCmd("+/comment/comment/set", {
        method : "POST",
        data   : JSON.stringify(req),
        callback : function(err, data) {

            // console.log(data);
            
            if (err) {
                return l4i.InnerAlert("#l5s-comment-embed-new-form-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#l5s-comment-embed-new-form-alert", 'alert-danger', data.error.message);
            }

            if (!data.kind || data.kind != "Comment") {
                return l4i.InnerAlert("#l5s-comment-embed-new-form-alert", 'alert-danger', "Network Exception");
            }

            req.meta = {
                id      : data.meta.id,
                created : l4i.TimeParseFormat(data.meta.created, "Y-m-d H:i"),
            };

            l4iTemplate.Render({
                dstid  : "l5s-comment-embed-list",
                tplid  : "l5s-comment-embed-tpl",
                data   : req,
                append : true,
                success : function() {
                    
                    $("#entry-"+ req.meta.id).css({
                        "outline": "#5cb85c solid 1px",
                    });

                    setTimeout(function() {
                        $("#entry-"+ req.meta.id).css({
                            "outline": "0px",
                        });
                    }, 2000);
                },
            });

            form.find("textarea[name=content]").val("");

            l4i.InnerAlert("#l5s-comment-embed-new-form-alert", 'alert-success', "Successfully commited");

            setTimeout(function() {
                $("#l5s-comment-embed-new-form-alert").hide(500);
            }, 2000);
        },
    });
}
