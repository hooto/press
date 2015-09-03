var l5sSys = {
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

l5sSys.Init = function()
{
    l4i.UrlEventRegister("sys/index", l5sSys.Index);
    l4i.UrlEventRegister("sys/status", l5sSys.Status);
    l4i.UrlEventRegister("sys/ids-status", l5sSys.IdentityStatus);
    l4i.UrlEventRegister("sys/config", l5sSys.Config);
}

l5sSys.Index = function()
{
    l5sMgr.TplCmd("sys/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            l5sSys.Status();
            // l5sSys.IdentityStatus();
            // l5sSys.Config();
        },
    });
}

l5sSys.Config = function()
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
    
        l5sMgr.ApiCmd("sys/config-list", {
            callback: ep.done('data'),
        });

        l5sMgr.TplCmd("sys/config", {
            callback: ep.done('tpl'),           
        });
    });
}

l5sSys.ConfigSetCommit = function()
{
    
    var form = $("#l5smgr-sys-configset"),
        alertid = "#l5smgr-sys-configset-alert",
        namereg = /^[a-z][a-z0-9_]+$/;

    var req = {
        items: [],
    }

    try {

        form.find(".l5smgr-sys-config-item").each(function() {
            
            req.items.push({
                key: $(this).attr("name"),
                value: $(this).val(),
            });
        });

    } catch (err) {
        l4i.InnerAlert(alertid, 'alert-danger', err);
        return;
    }

    l5sMgr.ApiCmd("sys/config-set", {
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


l5sSys.Status = function()
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
    
        l5sMgr.ApiCmd("sys/status", {
            callback: ep.done('data'),
        });

        l5sMgr.TplCmd("sys/status", {
            callback: ep.done('tpl'),           
        });
    });
}

l5sSys.IdentityStatus = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            if (!data) {
                return;
            }

            data._roles = l5sSys.roles;

            l4iTemplate.Render({
                dstid  : "work-content",
                tplsrc : tpl,
                data   : data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        l5sMgr.ApiCmd("sys/identity-status", {
            callback: ep.done('data'),
        });

        l5sMgr.TplCmd("sys/ids-status", {
            callback: ep.done('tpl'),           
        });
    });
}


l5sSys.IdentitySync = function()
{
    var form = $("#l5s-mgr-sys-ids");

    l5sMgr.Ajax("setup/app-register-put", {
        method : "POST",
        data   : form.serialize(),
        success: function(data) {
            
            if (data === undefined || data.kind != "AppInstanceRegister") {
                if (data.error) {
                    return l4i.InnerAlert("#l5s-mgr-sys-ids-alert", 'alert-danger', data.error.message);
                }

                return l4i.InnerAlert("#l5s-mgr-sys-ids-alert", 'alert-danger', "Network Connection Exception");
            }

            l4i.InnerAlert("#l5s-mgr-sys-ids-alert", 'alert-success', "Successful registered");
            
            window.setTimeout(function() {
                l5sSys.IdentityStatus();
            }, 1000);
        },
    });
}


l5sSys.UtilResourceSizeFormat = function(size)
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


l5sSys.UtilDurationFormat = function(timems, fix)
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

