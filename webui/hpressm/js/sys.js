// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

var hpressSys = {
    roles: {
        items: [{
            idxid: 100,
            meta: {
                name: "Member",
            },
        }, {
            idxid: 1000,
            meta: {
                name: "Guest",
            },
        }],
    },
}

hpressSys.Init = function() {
    l4i.UrlEventRegister("sys/index", hpressSys.Index, "hpressm-topbar");

    l4i.UrlEventRegister("sys/status", hpressSys.Status, "hpressm-sys-nav");
    l4i.UrlEventRegister("sys/iam-status", hpressSys.IamStatus, "hpressm-sys-nav");
    l4i.UrlEventRegister("sys/config", hpressSys.Config, "hpressm-sys-nav");
}

hpressSys.Index = function() {
    l4iStorage.Set("hpressm_nav_last_active", "sys/index");

    hpressMgr.TplCmd("sys/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            hpressSys.Status();
        // hpressSys.IamStatus();
        // hpressSys.Config();
        },
    });
}

hpressSys.Config = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data) {
                return;
            }

            for (var i in data.items) {
                if (!data.items[i].comment) {
                    data.items[i].comment = "";
                }
            }

            l4iTemplate.Render({
                dstid: "work-content",
                tplsrc: tpl,
                data: data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        hpressMgr.ApiCmd("sys/config-list", {
            callback: ep.done('data'),
        });

        hpressMgr.TplCmd("sys/config", {
            callback: ep.done('tpl'),
        });
    });
}

hpressSys.ConfigSetCommit = function() {

    var form = $("#hpressm-sys-configset"),
        alertid = "#hpressm-sys-configset-alert",
        namereg = /^[a-z][a-z0-9_]+$/;

    var req = {
        items: [],
    }

    try {

        form.find(".hpressm-sys-config-item").each(function() {

            req.items.push({
                key: $(this).attr("name"),
                value: $(this).val(),
            });
        });

    } catch (err) {
        l4i.InnerAlert(alertid, 'alert-danger', err);
        return;
    }

    hpressMgr.ApiCmd("sys/config-set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

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


hpressSys.Status = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data) {
                return;
            }

            // data._items = {};
            // for (var i in data.items) {
            //     data._items[data.items[i]["key"]] = data.items[i]["val"];
            // }

            l4iTemplate.Render({
                dstid: "work-content",
                tplsrc: tpl,
                data: data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        hpressMgr.ApiCmd("sys/status", {
            callback: ep.done('data'),
        });

        hpressMgr.TplCmd("sys/status", {
            callback: ep.done('tpl'),
        });
    });
}

hpressSys.IamStatus = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

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

            data._roles = hpressSys.roles;

            l4iTemplate.Render({
                dstid: "work-content",
                tplsrc: tpl,
                data: data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        hpressMgr.ApiCmd("sys/iam-status", {
            callback: ep.done('data'),
        });

        hpressMgr.TplCmd("sys/iam-status", {
            callback: ep.done('tpl'),
        });
    });
}


hpressSys.IamSync = function() {
    var alert_id = "#hpress-mgr-sys-iam-alert",
        form = $("#hpress-mgr-sys-iam"),
        url = "";
    if (form) {
        var v = form.find("input[name=app_title]").val();
        if (v) {
            url += "&app_title=" + v;
        }
        v = form.find("input[name=instance_url]").val();
        if (v) {
            url += "&instance_url=" + v;
        }
    }

    hpressMgr.Ajax("setup/app-register-sync?" + url, {
        callback: function(err, data) {

            if (!data || data.kind != "AppInstanceRegister") {
                if (data.error) {
                    return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
                }

                return l4i.InnerAlert(alert_id, 'alert-danger', "Network Connection Exception");
            }

            l4i.InnerAlert(alert_id, 'alert-success', "Successful registered");

            window.setTimeout(function() {
                hpressSys.IamStatus();
            }, 1000);
        },
    });
}


hpressSys.UtilResourceSizeFormat = function(size) {
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
            return (size / Math.pow(1024, ms[i][0])).toFixed(0) + " <span>" + ms[i][1] + "</span>";
        }
    }

    if (size == 0) {
        return size;
    }

    return size + " <span>B</span>";
}


hpressSys.UtilDurationFormat = function(timems, fix) {
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

                ts += t + " " + ms[i][1];

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
        ts = timems + " microseconds";
    } else {
        ts = "0";
    }

    return ts;
}
