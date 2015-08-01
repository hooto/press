var l5sComment = {

}

l5sComment.EmbedLoader = function(dstid, ref_modname, ref_table, ref_id)
{
    l5s.ActionLoader(dstid, "+/comment/comment/embed?refer_modname="+ ref_modname +
        "&refer_datax_table="+ ref_table +"&refer_id="+ ref_id );
}

l5sComment.EmbedCommit = function()
{
    var form = $("#l5s-comment-embed-new-form"),
        alertid = "#l5s-comment-embed-new-form-alert";

    var req = {
        pid               : "",
        refer_id          : form.find("input[name=refer_id]").val(),
        refer_modname     : form.find("input[name=refer_modname]").val(),
        refer_datax_table : form.find("input[name=refer_datax_table]").val(),
        author            : form.find("input[name=author]").val(),
        content           : form.find("textarea[name=content]").val(),
        captcha_token     : form.find("input[name=captcha_token]").val(),
        captcha_word      : form.find("input[name=captcha_word]").val(),
    };

    // console.log(req);
    l5s.ApiCmd("+/comment/comment/set", {
        method : "POST",
        data   : JSON.stringify(req),
        callback : function(err, data) {

            // console.log(data);
            
            if (err) {
                return l4i.InnerAlert(alertid, 'alert-danger', err);
            }

            if (!data || data.error) {

                if (data.error && data.error.code == "CaptchaNotMatch") {

                    var captcha_token = Math.random();
                    form.find("input[name=captcha_token]").val(captcha_token);
                    form.find("#l5s-comment-captcha-url").attr("src", 
                        "/+/hcaptcha/api/image?hcaptcha_token="+ captcha_token);
                }

                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            if (!data.kind || data.kind != "Comment") {
                return l4i.InnerAlert(alertid, 'alert-danger', "Network Exception");
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

                    $("#l5s-comment-embed-list-header").css({"display": "block"});
                    
                    $("#entry-"+ req.meta.id).css({
                        "outline": "#5cb85c solid 1px",
                    });

                    setTimeout(function() {
                        $("#entry-"+ req.meta.id).css({
                            "outline": "0px",
                        });
                    }, 1500);
                },
            });

            form.find("textarea[name=content]").val("");
            form.find("textarea[name=captcha_word]").val("");

            l4i.InnerAlert(alertid, 'alert-success', "Successfully commited");

            setTimeout(function() {
                $(alertid).hide(500);
                l5sComment.EmbedFormHidden();
            }, 1500);
        },
    });
}


l5sComment.EmbedFormActive = function()
{
    $("#l5s-comment-embed-new-form-ctrl").css({"display": "none"});

    var form = $("#l5s-comment-embed-new-form"),
        captcha_token = Math.random();

    form.css({"display": "block"});
   
    form.find("input[name=captcha_token]").val(captcha_token);
    form.find("#l5s-comment-captcha-url").attr("src", 
        "/+/hcaptcha/api/image?hcaptcha_token="+ captcha_token);

    form.find("textarea[name=content]").focus();
}


l5sComment.EmbedFormHidden = function()
{
    $("#l5s-comment-embed-new-form-ctrl").css({"display": "block"});
    $("#l5s-comment-embed-new-form").css({"display": "none"});
}
