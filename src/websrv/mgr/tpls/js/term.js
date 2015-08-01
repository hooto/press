var l5sTerm = {

}

l5sTerm.List = function(modname, modelid)
{
    var alertid = "#l5smgr-node-alert";

    if (!modname && l4iStorage.Get("l5smgr_spec_active")) {
        modname = l4iStorage.Get("l5smgr_spec_active");
    }
    if (!modelid && l4iStorage.Get("l5smgr_tmodel_active")) {
        modelid = l4iStorage.Get("l5smgr_tmodel_active");
    }

    if (!modname || !modelid) {
        return;
    }

    var uri = "modname="+ modname +"&modelid="+ modelid;
    if (document.getElementById("qry_text")) {
        uri += "&qry_text="+ $("#qry_text").val();
    }

    // console.log(uri);
    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", "data", function (tpl, rsj) {
            
            if (tpl) {
                $("#work-content").html(tpl);
            }

            l4iStorage.Set("l5smgr_tmodel_active", modelid);

            if (!rsj || rsj.kind != "TermList" 
                || !rsj.items || rsj.items.length < 1) {
                
                $("#l5smgr-nodels").empty();
                $("#l5smgr-termls").empty();
                
                return l4i.InnerAlert(alertid, 'alert-danger', "Item Not Found");
            }

            $(alertid).hide();

            for (var i in rsj.items) {
                
                rsj.items[i].created = l4i.TimeParseFormat(rsj.items[i].created, "Y-m-d");
                rsj.items[i].updated = l4i.TimeParseFormat(rsj.items[i].updated, "Y-m-d H:i:s");

                if (!rsj.items[i].weight) {
                    rsj.items[i].weight = 0;
                }
            }

            l4iTemplate.Render({
                dstid: "l5smgr-termls",
                tplid: "l5smgr-termls-tpl",
                data:  {
                    model   : rsj.model,
                    modname : modname,
                    modelid : modelid,
                    items   : rsj.items,
                },
            });
        });
    
        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-termlist)");
        });

        // template
        var el = document.getElementById("l5smgr-termls");
        if (!el || el.length < 1) {
            l5sMgr.TplCmd("/term/list", {
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

        l5sMgr.ApiCmd("/term/list?"+ uri, {
            callback: ep.done("data"),           
        });
    });
}


l5sTerm.Set = function(modname, modelid, termid)
{
    var alertid = "#l5smgr-node-alert";

    if (!modname && l4iStorage.Get("l5smgr_spec_active")) {
        modname = l4iStorage.Get("l5smgr_spec_active");
    }
    if (!modelid && l4iStorage.Get("l5smgr_tmodel_active")) {
        modelid = l4iStorage.Get("l5smgr_tmodel_active");
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

            if (!data.state) {
                data.state = 1;
            }
            if (!data.weight) {
                data.weight = 0;
            }
            if (!data.pid) {
                data.pid = 0;
            }

            $(alertid).hide();

            l4iTemplate.Render({
                dstid: "l5smgr-termset",
                tplid: "l5smgr-termset-tpl",
                data:  data,
                success: function() {
                    
                },
            });
        });
    
        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        l5sMgr.TplCmd("/term/set", {
            callback: function(err, tpl) {
                    
                if (err) {
                    return ep.emit('error', err);
                }
                ep.emit("tpl", tpl);
            }
        });

        if (termid) {
            l5sMgr.ApiCmd("/term/entry?"+ uri +"&id="+ termid, {
                callback: ep.done("data"),           
            });
        } else {
            l5sMgr.ApiCmd("/term-model/entry?"+ uri, {
                callback: function(err, data) {
                    ep.emit("data", {
                        kind  : "Term",
                        model : data,
                        id    : "0",
                        pid   : "0",
                        title : "",
                        state : "1",
                        weight: "0",
                    });
                },           
            });
        }
    });
}

l5sTerm.SetCommit = function()
{
    var form = $("#l5smgr-termset"),
        alertid = "#l5smgr-node-alert";

    var req = {
        kind   : "Term",
        id     : parseInt(form.find("input[name=id]").val()),
        title  : form.find("input[name=title]").val(),
        state  : parseInt(form.find("input[name=state]").val()),
    }

    var model_type = form.find("input[name=model_type]").val();
    if (model_type = "taxonomy") {
        req.weight = parseInt(form.find("input[name=weight]").val());
        req.pid    = parseInt(form.find("input[name=pid]").val());
    } else if (model_type = "tag") {

    }

    console.log(JSON.stringify(req));

    //
    var uri = "modname="+ l4iStorage.Get("l5smgr_spec_active") +
        "&modelid="+ l4iStorage.Get("l5smgr_tmodel_active");

    l5sMgr.ApiCmd("term/set?"+ uri, {
        method : "POST",
        data   : JSON.stringify(req),
        callback : function(err, data) {

            if (!data || data.kind != "Term") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            form.find("input[name=id]").val(data.id);

            l4i.InnerAlert(alertid, 'alert-success', "Successful operation");
            setTimeout(l5sTerm.List, 500);
        }
    });
}
