var hpressMgrModel = {

}

hpressMgrModel.List = function(tplid) {
    if (!tplid) {
        tplid = "hpressm-ctnls";
    }

    var qry_grpdev = $("#" + tplid + "-grpdev-qryid").attr("href").substr(1);
    var qry_grpdev_dp = "Groups";
    if (qry_grpdev != "") {
        qry_grpdev_dp = $("#" + tplid + "-grpdev-qrydp").text();
    }

    var uri = "qry_text=" + $("#" + tplid + "-qry-text").val();
    uri += "&qry_grpdev=" + qry_grpdev;

    var req = {
        items: [{
            name: "groups",
            uri: hpressMgr.base + "ext/lps/group/list",
        },
            {
                name: "info",
                uri: hpressMgr.base + "ext/lps/pkg-info/list?" + uri,
            },
        ],
    }

    //
    $.ajax({
        type: "POST",
        url: hpressMgr.base + "v1/mixer",
        data: JSON.stringify(req),
        timeout: 3000,
        success: function(rsp) {

            var rsj = JSON.parse(rsp);
            if (rsj === undefined || rsj.kind != "Mixer" || rsj.items === undefined) {
                $("#" + tplid + "-empty-alert").show();
                return;
            }

            if (rsj.items.groups === undefined) {
                $("#" + tplid + "-empty-alert").show();
                return;
            }

            if (rsj.items.info === undefined || rsj.items.info.kind != "PackageInfoList" || rsj.items.info.items === undefined) {
                rsj.items.info.items = [];
            }

            if (rsj.items.info.items.length > 0) {
                $("#" + tplid + "-empty-alert").hide();
            } else {
                $("#" + tplid + "-empty-alert").show();
            }

            for (var i in rsj.items.info.items) {
                rsj.items.info.items[i].meta.updated = l4i.TimeParseFormat(rsj.items.info.items[i].meta.updated, "Y-m-d");
            }

            lessTemplate.Render({
                dstid: tplid,
                tplid: tplid + "-tpl",
                data: rsj.items.info.items,
            });

            if (rsj.items.groups.kind !== undefined &&
                rsj.items.groups.kind == "PackageGroupList" &&
                rsj.items.groups.dev !== undefined &&
                rsj.items.groups.dev.length > 0) {
                lessTemplate.Render({
                    dstid: tplid + "-grpdev",
                    tplid: tplid + "-grpdev-tpl",
                    data: {
                        grpdev: rsj.items.groups.dev,
                        qry_grpdev: qry_grpdev,
                        qry_grpdev_dp: qry_grpdev_dp,
                    },
                });
            }
        },
        error: function(xhr, textStatus, error) {
            //lessAlert("#azt02e", 'alert-danger', textStatus+' '+xhr.responseText);
        }
    });
}