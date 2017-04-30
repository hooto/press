var htpS2 = {

}

htpS2.Init = function()
{
    l4i.UrlEventRegister("s2/index", htpS2.Index, "htpm-topbar");
}

htpS2.Index = function()
{
    l4iStorage.Set("htpm_nav_last_active", "s2/index");

    htpMgr.TplCmd("s2/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            htpS2.ObjList();
        },
    });
}

htpS2.ObjList = function(path)
{
    if (!path) {
        path = l4iStorage.Get("htpm_s2_obj_path_active");
    }

    if (!path) {
        path = "/";
    }

    l4iStorage.Set("htpm_s2_obj_path_active", path);

    htpMgr.ApiCmd("s2-obj/list?path="+ path, {
        callback: function(err, data) {

            if (err || !data || !data.kind) {
                return;
            }

            if (!data.items) {
                data.items = [];
            }

            data._path = path;

            for (var i in data.items) {

                var name = data.items[i].name;

                data.items[i]._id = l4iString.CryptoMd5(path +"/"+ name);

                data.items[i]._abspath = path +"/"+ name;

                if (name.toLowerCase().substr(-4) == ".jpg" || name.toLowerCase().substr(-5) == ".jpeg"
                    || name.toLowerCase().substr(-4) == ".png" || name.toLowerCase().substr(-4) == ".gif") {
                    data.items[i]._isimg = true;
                } else {
                    data.items[i]._isimg = false;
                }
            }

            l4iTemplate.Render({
                dstid: "htpm-s2-objls",
                tplid: "htpm-s2-objls-tpl",
                data:  data,
            });


            var dirnav = [];

            dirnav.push({
                path   : "/",
                name   : "Home",
            });

            //
            path = l4i.StringTrim(path.replace(/\/+/g, "/"), "/");
            if (path.length > 0) {

                var prs = path.split("/");
                var ppath = "";
                for (var i in prs) {
                    ppath += "/"+ prs[i];
                    dirnav.push({
                        path   : ppath,
                        name   : prs[i],
                    });
                }
            }

            l4iTemplate.Render({
                dstid: "htpm-s2-objls-dirnav",
                tplid: "htpm-s2-objls-dirnav-tpl",
                data:  {
                    items: dirnav,
                },
            });
        },
    });
}

htpS2.ObjNew = function(type, path, file)
{
    if (!path) {
        path = l4iStorage.Get("htpm_s2_obj_path_active");

        if (!path) {
            path = "/";
        }
    }

    var formid = Math.random().toString(36).slice(2);

    var req = {
        title        : (type == "dir") ? "New Folder" : "New File",
        width        : 700,
        height       : 350,
        tplid        : "htpm-s2-objnew-tpl",
        data         : {
            formid   : formid,
            file     : file,
            path     : path,
            type     : type
        },
        buttons      : [
            {
                onclick : "htpS2.ObjNewSave(\""+ formid +"\")",
                title   : "Upload",
                style   : "btn-primary"
            },
            {
                onclick : "l4iModal.Close()",
                title   : "Close"
            }
        ]
    }

    // req.success = function() {
    //     $("#"+ formid +" :input[name=name]").focus();
    // }

    l4iModal.Open(req);
}

htpS2.ObjNewSave = function(formid)
{
    var elem = document.getElementById("htpm-s2-objnew-files");

    for (var i = 0; i < elem.files.length; i++) {

        htpS2._objNewUpload(formid, elem.files[i]);
    }
}

htpS2._objNewUpload = function(formid, file)
{
    var reader = new FileReader();

    reader.onload = (function(file) {

        return function(e) {

            if (e.target.readyState != FileReader.DONE) {
                return;
            }

            var ppath = $("#"+ formid +" :input[name=path]").val();

            htpMgr.ApiCmd("s2-obj/put", {
                method  : "POST",
                data    : JSON.stringify({
                    path    : ppath +"/"+ file.name,
                    size    : file.size,
                    body    : e.target.result,
                    encode  : "base64",
                }),
                callback : function(err, rsp) {

                    if (rsp && rsp.kind && rsp.kind == "FsFile") {

                        $("#"+ formid +"-alert").show().append("<div>"+ file.name +" OK</div>");
                        htpS2.ObjList(ppath);

                        setTimeout(function() {
                            l4iModal.Close();
                        }, 1000);
                    } else {

                        if (rsp.error) {
                            $("#"+ formid +"-alert").show().append("<div>"+ file.name +" Failed: "+ rsp.error.message +"</div>");
                        } else {
                            $("#"+ formid +"-alert").show().append("<div>"+ file.name +" Failed</div>");
                        }
                    }
                },
                error: function(status, message) {
                    $("#"+ formid +"-alert").show().append("<div>"+ file.name +" Failed</div>");
                }
            });
        };

    })(file);

    reader.readAsDataURL(file);
}

htpS2.ObjDel = function(path)
{
    //
    htpMgr.ApiCmd("s2-obj/del?path="+ path, {

        callback: function(err, data) {
            if (data.kind && data.kind == "FsFile") {
                $("#obj"+ l4iString.CryptoMd5(path)).remove();
            } else if (data.error) {
                alert(data.error.message);
            }
        },
    });
}

htpS2.UtilResourceSizeFormat = function(size)
{
    if (!size) {
        size = 0;
    }

    var ms = [
        [6, "EB"],
        [5, "PB"],
        [4, "TB"],
        [3, "GB"],
        [2, "MB"],
        [1, "KB"],
    ];
    for (var i in ms) {
        if (size > Math.pow(1024, ms[i][0])) {
            return (size / Math.pow(1024, ms[i][0])).toFixed(0) +" <span>"+ ms[i][1] +"</span>";
        }
    }

    if (size == 0) {
        return size;
    }

    return size + " <span>B</span>";
}
