var htapComment = {

}

htapComment.EmbedLoader = function(dstid, ref_modname, ref_table, ref_id)
{
    htap.ActionLoader(dstid, "+/comment/comment/embed?refer_modname="+ ref_modname +
        "&refer_datax_table="+ ref_table +"&refer_id="+ ref_id );
}

htapComment.EmbedCommit = function()
{
    var form = $("#htap-comment-embed-new-form"),
        alertid = "#htap-comment-embed-new-form-alert";

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
    htap.ApiCmd("+/comment/comment/set", {
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
                    form.find("#htap-comment-captcha-url").attr("src", 
                        htap.HttpSrvBasePath("+/hcaptcha/api/image?hcaptcha_token="+ captcha_token));
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
                dstid  : "htap-comment-embed-list",
                tplid  : "htap-comment-embed-tpl",
                data   : req,
                append : true,
                success : function() {

                    $("#htap-comment-embed-list-header").css({"display": "block"});
                    
                    $("#entry-"+ req.meta.id).css({
                        "outline": "#5cb85c solid 2px",
                    });

                    setTimeout(function() {
                        $("#entry-"+ req.meta.id).css({
                            "outline": "0px",
                        });
                    }, 1500);
                },
            });

            form.find("textarea[name=content]").val("");
            form.find("input[name=captcha_word]").val("");

            l4i.InnerAlert(alertid, 'alert-success', "Successfully commited");

            setTimeout(function() {
                htapComment.EmbedFormHidden();
                $(alertid).hide(500);
            }, 1500);
        },
    });
}


htapComment.EmbedFormActive = function()
{
    $("#htap-comment-embed-new-form-ctrl").css({"display": "none"});

    var form = $("#htap-comment-embed-new-form"),
        captcha_token = Math.random();

    form.css({"display": "block"});

    form.find("input[name=captcha_token]").val(captcha_token);
    form.find("#htap-comment-captcha-url").attr("src", 
        htap.HttpSrvBasePath("+/hcaptcha/api/image?hcaptcha_token="+ captcha_token));

    form.find("textarea[name=content]").focus();

    $("body").scrollTop($(window).scrollTop() + 350);
}


htapComment.EmbedFormHidden = function()
{
    $("#htap-comment-embed-new-form-ctrl").css({"display": "block"});
    $("#htap-comment-embed-new-form").slideUp(500);//css({"display": "none"});
}
