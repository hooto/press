var htapSys = {
    roles : {
        items: [{
            idxid: 100,
            meta: {
                name : "Member",
            },
        },{
            idxid: 1000,
            meta: {
                name : "Guest",
            },
        }],
    },
}

htapSys.Init = function()
{
    l4i.UrlEventRegister("sys/index", htapSys.Index, "htapm-topbar");

    l4i.UrlEventRegister("sys/status", htapSys.Status, "htapm-sys-nav");
    l4i.UrlEventRegister("sys/iam-status", htapSys.IamStatus, "htapm-sys-nav");
    l4i.UrlEventRegister("sys/config", htapSys.Config, "htapm-sys-nav");
}

htapSys.Index = function()
{
    l4iStorage.Set("htapm_nav_last_active", "sys/index");

    htapMgr.TplCmd("sys/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            htapSys.Status();
            // htapSys.IamStatus();
            // htapSys.Config();
        },
    });
}

htapSys.Config = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            if (!data) {
                return;
            }

            for (var i in data.items) {
                if (!data.items[i].comment) {
                    data.items[i].comment = "";
                }
            }

            l4iTemplate.Render({
                dstid  : "work-content",
                tplsrc : tpl,
                data   : data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        htapMgr.ApiCmd("sys/config-list", {
            callback: ep.done('data'),
        });

        htapMgr.TplCmd("sys/config", {
            callback: ep.done('tpl'),
        });
    });
}

htapSys.ConfigSetCommit = function()
{

    var form = $("#htapm-sys-configset"),
        alertid = "#htapm-sys-configset-alert",
        namereg = /^[a-z][a-z0-9_]+$/;

    var req = {
        items: [],
    }

    try {

        form.find(".htapm-sys-config-item").each(function() {

            req.items.push({
                key: $(this).attr("name"),
                value: $(this).val(),
            });
        });

    } catch (err) {
        l4i.InnerAlert(alertid, 'alert-danger', err);
        return;
    }

    htapMgr.ApiCmd("sys/config-set", {
        method  : "PUT",
        data    : JSON.stringify(req),
        success : function(data) {

            if (!data || !data.kind || data.kind != "SysConfigList") {

                if (data.error) {
                    return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
                }

                return l4i.InnerAlert(alertid, 'alert-danger', "Network Connection Exception");
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful updated");
        },
    });
}


htapSys.Status = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            if (!data) {
                return;
            }

            // data._items = {};
            // for (var i in data.items) {
            //     data._items[data.items[i]["key"]] = data.items[i]["val"];
            // }

            l4iTemplate.Render({
                dstid  : "work-content",
                tplsrc : tpl,
                data   : data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        htapMgr.ApiCmd("sys/status", {
            callback: ep.done('data'),
        });

        htapMgr.TplCmd("sys/status", {
            callback: ep.done('tpl'),
        });
    });
}

htapSys.IamStatus = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            if (!data) {
                return;
            }

            data._roles = htapSys.roles;

            l4iTemplate.Render({
                dstid  : "work-content",
                tplsrc : tpl,
                data   : data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        htapMgr.ApiCmd("sys/iam-status", {
            callback: ep.done('data'),
        });

        htapMgr.TplCmd("sys/iam-status", {
            callback: ep.done('tpl'),
        });
    });
}


htapSys.IamSync = function()
{
    var form = $("#htap-mgr-sys-iam");

    htapMgr.Ajax("setup/app-register-put", {
        method : "POST",
        data   : form.serialize(),
        success: function(data) {

            if (data === undefined || data.kind != "AppInstanceRegister") {
                if (data.error) {
                    return l4i.InnerAlert("#htap-mgr-sys-iam-alert", 'alert-danger', data.error.message);
                }

                return l4i.InnerAlert("#htap-mgr-sys-iam-alert", 'alert-danger', "Network Connection Exception");
            }

            l4i.InnerAlert("#htap-mgr-sys-iam-alert", 'alert-success', "Successful registered");

            window.setTimeout(function() {
                htapSys.IamStatus();
            }, 1000);
        },
    });
}


htapSys.UtilResourceSizeFormat = function(size)
{
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


htapSys.UtilDurationFormat = function(timems, fix)
{
    var ms = [
        [86400000, "day"],
        [3600000, "hour"],
        [60000, "minute"],
        [1000, "second"],
    ];

    if (fix) {
        timems = parseInt(timems / fix);
    } else {
        timems = parseInt(timems);
    }

    var ts = "";

    for (var i in ms) {

        if (timems >= ms[i][0]) {

            var t = parseInt(timems / ms[i][0]);

            if (t > 0) {

                ts += t + " "+ ms[i][1];

                if (t > 1) {
                    ts += "s";
                }

                ts += ", ";

                timems = parseInt(timems % ms[i][0]);
            }
        }
    }

    if (ts.length > 2) {
        ts = ts.substr(0, ts.length - 2);
    } else if (timems > 0) {
        ts = timems +" microseconds";
    } else {
        ts = "0";
    }

    return ts;
}

