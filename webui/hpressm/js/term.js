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

var hpressTerm = {
    taxonomy_ls_cache: null,
}

hpressTerm.List = function(modname, modelid) {
    var alertid = "#hpressm-node-alert",
        page = 0;

    if (!modname && l4iStorage.Get("hpressm_spec_active")) {
        modname = l4iStorage.Get("hpressm_spec_active");
    }
    if (!modelid && l4iStorage.Get("hpressm_tmodel_active")) {
        modelid = l4iStorage.Get("hpressm_tmodel_active");
    }
    if (l4iStorage.Get("hpressm_termls_page")) {
        page = l4iStorage.Get("hpressm_termls_page");
    }

    if (!modname || !modelid) {
        return;
    }

    var uri = "modname=" + modname + "&modelid=" + modelid + "&page=" + page;
    if (document.getElementById("qry_text")) {
        uri += "&qry_text=" + $("#qry_text").val();
    }

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create("tpl", "data", function(tpl, rsj) {

            if (tpl) {
                $("#work-content").html(tpl);
            }

            l4iStorage.Set("hpressm_tmodel_active", modelid);

            if (!rsj || rsj.kind != "TermList" ||
                !rsj.items || rsj.items.length < 1) {

                $("#hpressm-nodels").empty();
                $("#hpressm-termls").empty();

                l4i.InnerAlert(alertid, 'alert-info', "Item Not Found");
            } else {
                $(alertid).hide();
            }
            $("#hpressm-term-list-new-title").text("New " + rsj.model.title);

            if (!rsj.items) {
                rsj.items = [];
            }

            for (var i in rsj.items) {

                rsj.items[i].created = l4i.TimeParseFormat(rsj.items[i].created, "Y-m-d");
                rsj.items[i].updated = l4i.TimeParseFormat(rsj.items[i].updated, "Y-m-d H:i:s");

                if (!rsj.items[i].weight) {
                    rsj.items[i].weight = 0;
                }

                if (!rsj.items[i].pid) {
                    rsj.items[i].pid = 0;
                }

                if (rsj.model.type == "taxonomy" && rsj.items[i].pid == 0) {
                    rsj.items[i]._subs = hpressTerm.ListSubRange(rsj.items, null, rsj.items[i].id, 0);
                }
            }

            hpressTerm.taxonomy_ls_cache = rsj;

            l4iTemplate.Render({
                dstid: "hpressm-termls",
                tplid: "hpressm-termls-tpl",
                data: {
                    model: rsj.model,
                    modname: modname,
                    modelid: modelid,
                    items: rsj.items,
                },
                success: function() {

                    if (rsj.model.type != "taxonomy") {
                        rsj.meta.RangeLen = 20;

                        l4iTemplate.Render({
                            dstid: "hpressm-termls-pager",
                            tplid: "hpressm-termls-pager-tpl",
                            data: l4i.Pager(rsj.meta),
                        });
                    } else {
                        $("#hpressm-termls-pager").empty();
                    }

                    hpressNode.OpToolsRefresh("#hpressm-node-term-opts");
                }
            });
        });

        ep.fail(function(err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-termlist)");
        });

        // template
        var el = document.getElementById("hpressm-termls");
        if (!el || el.length < 1) {
            hpressMgr.TplCmd("term/list", {
                callback: function(err, tpl) {

                    if (err) {
                        return ep.emit('error', err);
                    }

                    ep.emit("tpl", tpl);
                }
            });
        } else {
            ep.emit("tpl", null);
        }

        hpressMgr.ApiCmd("term/list?" + uri, {
            callback: ep.done("data"),
        });
    });
}

hpressTerm.Sprint = function(num) {
    var s = "";
    for (i = 0; i < num; i++) {
        s += "&nbsp;&nbsp;&nbsp;&nbsp;";
    }

    return s;
}

hpressTerm.ListSubRange = function(ls, rs, pid, dpnum) {
    if (!rs) {
        rs = [];
    }

    dpnum++;

    for (var i in rs) {
        if (rs[i].id == pid) {
            // return rs;
        }
    }

    for (var i in ls) {

        if (ls[i].pid == pid) {
            ls[i]._dp = dpnum;
            rs.push(ls[i]);
            rs = hpressTerm.ListSubRange(ls, rs, ls[i].id, dpnum);
        }
    }

    return rs;
}

hpressTerm.ListPage = function(page) {
    l4iStorage.Set("hpressm_termls_page", parseInt(page));
    hpressTerm.List();
}

hpressTerm.Set = function(modname, modelid, termid) {
    var alertid = "#hpressm-node-alert";

    if (!modname && l4iStorage.Get("hpressm_spec_active")) {
        modname = l4iStorage.Get("hpressm_spec_active");
    }
    if (!modelid && l4iStorage.Get("hpressm_tmodel_active")) {
        modelid = l4iStorage.Get("hpressm_tmodel_active");
    }

    if (!modname || !modelid) {
        return;
    }

    var uri = "modname=" + modname + "&modelid=" + modelid;

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create("tpl", "data", function(tpl, data) {

            if (!tpl) {
                return; // TODO
            }

            $("#work-content").html(tpl);

            if (!data || data.kind != "Term") {
                return l4i.InnerAlert(alertid, 'alert-info', "Item Not Found");
            }

            if (!data.status) {
                data.status = 1;
            }
            if (!data.weight) {
                data.weight = 0;
            }
            if (!data.pid) {
                data.pid = 0;
            }

            data._taxonomy_ls = hpressTerm.taxonomy_ls_cache;

            $(alertid).hide();
            hpressNode.OpToolsRefresh("clean");

            l4iTemplate.Render({
                dstid: "hpressm-termset",
                tplid: "hpressm-termset-tpl",
                data: data,
                success: function() {},
            });
        });

        ep.fail(function(err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        hpressMgr.TplCmd("term/set", {
            callback: function(err, tpl) {

                if (err) {
                    return ep.emit('error', err);
                }
                ep.emit("tpl", tpl);
            }
        });

        if (termid) {
            hpressMgr.ApiCmd("term/entry?" + uri + "&id=" + termid, {
                callback: ep.done("data"),
            });
        } else {
            hpressMgr.ApiCmd("term-model/entry?" + uri, {
                callback: function(err, data) {
                    ep.emit("data", {
                        kind: "Term",
                        model: data,
                        id: "0",
                        pid: "0",
                        title: "",
                        status: "1",
                        weight: "0",
                    });
                },
            });
        }
    });
}

hpressTerm.SetCommit = function() {
    var form = $("#hpressm-termset"),
        alertid = "#hpressm-node-alert";

    var req = {
        kind: "Term",
        id: parseInt(form.find("input[name=id]").val()),
        title: form.find("input[name=title]").val(),
        status: parseInt(form.find("input[name=status]").val()),
    }

    var model_type = form.find("input[name=model_type]").val();
    if (model_type = "taxonomy") {
        req.weight = parseInt(form.find("input[name=weight]").val());
        req.pid = parseInt(form.find("select[name=pid]").val());
    } else if (model_type = "tag") {
        //
    }

    //
    var uri = "modname=" + l4iStorage.Get("hpressm_spec_active") +
    "&modelid=" + l4iStorage.Get("hpressm_tmodel_active");

    hpressMgr.ApiCmd("term/set?" + uri, {
        method: "POST",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (!data || data.kind != "Term") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            form.find("input[name=id]").val(data.id);

            l4i.InnerAlert(alertid, 'alert-success', "Successful operation");
            setTimeout(hpressTerm.List, 500);
        }
    });
}
