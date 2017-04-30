var htpSys = {
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

htpSys.Init = function()
{
    l4i.UrlEventRegister("sys/index", htpSys.Index, "htpm-topbar");

    l4i.UrlEventRegister("sys/status", htpSys.Status, "htpm-sys-nav");
    l4i.UrlEventRegister("sys/iam-status", htpSys.IamStatus, "htpm-sys-nav");
    l4i.UrlEventRegister("sys/config", htpSys.Config, "htpm-sys-nav");
}

htpSys.Index = function()
{
    l4iStorage.Set("htpm_nav_last_active", "sys/index");

    htpMgr.TplCmd("sys/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            htpSys.Status();
            // htpSys.IamStatus();
            // htpSys.Config();
        },
    });
}

htpSys.Config = function()
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

        htpMgr.ApiCmd("sys/config-list", {
            callback: ep.done('data'),
        });

        htpMgr.TplCmd("sys/config", {
            callback: ep.done('tpl'),
        });
    });
}

htpSys.ConfigSetCommit = function()
{

    var form = $("#htpm-sys-configset"),
        alertid = "#htpm-sys-configset-alert",
        namereg = /^[a-z][a-z0-9_]+$/;

    var req = {
        items: [],
    }

    try {

        form.find(".htpm-sys-config-item").each(function() {

            req.items.push({
                key: $(this).attr("name"),
                value: $(this).val(),
            });
        });

    } catch (err) {
        l4i.InnerAlert(alertid, 'alert-danger', err);
        return;
    }

    htpMgr.ApiCmd("sys/config-set", {
        method  : "PUT",
        data    : JSON.stringify(req),
        callback : function(err, data) {

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


htpSys.Status = function()
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

        htpMgr.ApiCmd("sys/status", {
            callback: ep.done('data'),
        });

        htpMgr.TplCmd("sys/status", {
            callback: ep.done('tpl'),
        });
    });
}

htpSys.IamStatus = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            if (!data) {
                return;
            }

            if (!data.instance_self.privileges) {
                data.instance_self.privileges = [];
            }
            for (var i in data.instance_self.privileges) {

                if (!data.instance_self.privileges[i].roles) {
                    data.instance_self.privileges[i].roles = [];
                }
            }

            if (!data.instance_registered.privileges) {
                data.instance_registered.privileges = [];
            }
            for (var i in data.instance_registered.privileges) {

                if (!data.instance_registered.privileges[i].roles) {
                    data.instance_registered.privileges[i].roles = [];
                }
            }

            data._roles = htpSys.roles;

            l4iTemplate.Render({
                dstid  : "work-content",
                tplsrc : tpl,
                data   : data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        htpMgr.ApiCmd("sys/iam-status", {
            callback: ep.done('data'),
        });

        htpMgr.TplCmd("sys/iam-status", {
            callback: ep.done('tpl'),
        });
    });
}


htpSys.IamSync = function()
{
    var form = $("#htp-mgr-sys-iam");
    var alert_id = "#htp-mgr-sys-iam-alert";

    htpMgr.Ajax("setup/app-register-sync", {
        method : "POST",
        data   : form.serialize(),
        callback: function(err, data) {

            console.log("ASDFASD");

            if (!data || data.kind != "AppInstanceRegister") {
                if (data.error) {
                    return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
                }

                return l4i.InnerAlert(alert_id, 'alert-danger', "Network Connection Exception");
            }

            l4i.InnerAlert(alert_id, 'alert-success', "Successful registered");

            window.setTimeout(function() {
                htpSys.IamStatus();
            }, 1000);
        },
    });
}


htpSys.UtilResourceSizeFormat = function(size)
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


htpSys.UtilDurationFormat = function(timems, fix)
{
    var ms = [
        [86400000, "day"],
        [3600000, "hour"],
        [60000, "minute"],
        [1000, "second"],
    ];

    if (!timems) {
        timems = 0;
    }

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

