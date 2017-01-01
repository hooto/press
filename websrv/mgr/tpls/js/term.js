var htpTerm = {
    taxonomy_ls_cache : null,
}

htpTerm.List = function(modname, modelid)
{
    var alertid = "#htpm-node-alert",
        page = 0;

    if (!modname && l4iStorage.Get("htpm_spec_active")) {
        modname = l4iStorage.Get("htpm_spec_active");
    }
    if (!modelid && l4iStorage.Get("htpm_tmodel_active")) {
        modelid = l4iStorage.Get("htpm_tmodel_active");
    }
    if (l4iStorage.Get("htpm_termls_page")) {
        page = l4iStorage.Get("htpm_termls_page");
    }

    if (!modname || !modelid) {
        return;
    }

    var uri = "modname="+ modname +"&modelid="+ modelid +"&page="+ page;
    if (document.getElementById("qry_text")) {
        uri += "&qry_text="+ $("#qry_text").val();
    }

    // console.log(uri);
    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", "data", function (tpl, rsj) {

            if (tpl) {
                $("#work-content").html(tpl);
            }

            l4iStorage.Set("htpm_tmodel_active", modelid);

            if (!rsj || rsj.kind != "TermList"
                || !rsj.items || rsj.items.length < 1) {

                $("#htpm-nodels").empty();
                $("#htpm-termls").empty();

                return l4i.InnerAlert(alertid, 'alert-danger', "Item Not Found");
            }

            $(alertid).hide();

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
                    rsj.items[i]._subs = htpTerm.ListSubRange(rsj.items, null, rsj.items[i].id, 0);
                }
            }

            htpTerm.taxonomy_ls_cache = rsj;

            l4iTemplate.Render({
                dstid: "htpm-termls",
                tplid: "htpm-termls-tpl",
                data:  {
                    model   : rsj.model,
                    modname : modname,
                    modelid : modelid,
                    items   : rsj.items,
                },
                success: function() {

                    if (rsj.model.type != "taxonomy") {
                        rsj.meta.RangeLen = 20;

                        l4iTemplate.Render({
                            dstid : "htpm-termls-pager",
                            tplid : "htpm-termls-pager-tpl",
                            data  : l4i.Pager(rsj.meta),
                        });
                    } else {
                        $("#htpm-termls-pager").empty();
                    }

                    htpNode.OpToolsRefresh("#htpm-node-term-opts");
                }
            });
        });

        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-termlist)");
        });

        // template
        var el = document.getElementById("htpm-termls");
        if (!el || el.length < 1) {
            htpMgr.TplCmd("term/list", {
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

        htpMgr.ApiCmd("term/list?"+ uri, {
            callback: ep.done("data"),
        });
    });
}

htpTerm.Sprint = function(num)
{
    var s = "";
    for (i = 0; i < num; i++) {
        s += "&nbsp;&nbsp;&nbsp;&nbsp;";
    }

    return s;
}

htpTerm.ListSubRange = function(ls, rs, pid, dpnum)
{
    if (!rs) {
        rs = [];
    }

    dpnum++;

    // console.log(rs, pid);
    // console.log(ls.length);

    for (var i in rs) {
        if (rs[i].id == pid) {
            // return rs;
        }
    }

    for (var i in ls) {

        if (ls[i].pid == pid) {
            ls[i]._dp = dpnum;
            rs.push(ls[i]);
            rs = htpTerm.ListSubRange(ls, rs, ls[i].id, dpnum);
        }
    }

    // console.log(rs);

    return rs;
}

htpTerm.ListPage = function(page)
{
    l4iStorage.Set("htpm_termls_page", parseInt(page));
    htpTerm.List();
}

htpTerm.Set = function(modname, modelid, termid)
{
    var alertid = "#htpm-node-alert";

    if (!modname && l4iStorage.Get("htpm_spec_active")) {
        modname = l4iStorage.Get("htpm_spec_active");
    }
    if (!modelid && l4iStorage.Get("htpm_tmodel_active")) {
        modelid = l4iStorage.Get("htpm_tmodel_active");
    }

    if (!modname || !modelid) {
        return;
    }

    var uri = "modname="+ modname +"&modelid="+ modelid;

    // console.log(uri);
    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", "data", function (tpl, data) {

            if (!tpl) {
                return; // TODO
            }

            $("#work-content").html(tpl);

            if (!data || data.kind != "Term") {
                return l4i.InnerAlert(alertid, 'alert-danger', "Item Not Found");
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

            data._taxonomy_ls = htpTerm.taxonomy_ls_cache;

            $(alertid).hide();
            htpNode.OpToolsRefresh();

            l4iTemplate.Render({
                dstid: "htpm-termset",
                tplid: "htpm-termset-tpl",
                data:  data,
                success: function() {

                },
            });
        });

        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        htpMgr.TplCmd("term/set", {
            callback: function(err, tpl) {

                if (err) {
                    return ep.emit('error', err);
                }
                ep.emit("tpl", tpl);
            }
        });

        if (termid) {
            htpMgr.ApiCmd("term/entry?"+ uri +"&id="+ termid, {
                callback: ep.done("data"),
            });
        } else {
            htpMgr.ApiCmd("term-model/entry?"+ uri, {
                callback: function(err, data) {
                    ep.emit("data", {
                        kind  : "Term",
                        model : data,
                        id    : "0",
                        pid   : "0",
                        title : "",
                        status : "1",
                        weight: "0",
                    });
                },
            });
        }
    });
}

htpTerm.SetCommit = function()
{
    var form = $("#htpm-termset"),
        alertid = "#htpm-node-alert";

    var req = {
        kind   : "Term",
        id     : parseInt(form.find("input[name=id]").val()),
        title  : form.find("input[name=title]").val(),
        status  : parseInt(form.find("input[name=status]").val()),
    }

    var model_type = form.find("input[name=model_type]").val();
    if (model_type = "taxonomy") {
        req.weight = parseInt(form.find("input[name=weight]").val());
        req.pid    = parseInt(form.find("select[name=pid]").val());
    } else if (model_type = "tag") {

    }

    // console.log(JSON.stringify(req));

    //
    var uri = "modname="+ l4iStorage.Get("htpm_spec_active") +
        "&modelid="+ l4iStorage.Get("htpm_tmodel_active");

    htpMgr.ApiCmd("term/set?"+ uri, {
        method : "POST",
        data   : JSON.stringify(req),
        callback : function(err, data) {

            if (!data || data.kind != "Term") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            form.find("input[name=id]").val(data.id);

            l4i.InnerAlert(alertid, 'alert-success', "Successful operation");
            setTimeout(htpTerm.List, 500);
        }
    });
}
