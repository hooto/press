var l5sTerm = {

}

l5sTerm.List = function(specid, modelid)
{
    if (!specid && l4iStorage.Get("l5smgr_spec_active")) {
        specid = l4iStorage.Get("l5smgr_spec_active");
    }
    if (!modelid && l4iStorage.Get("l5smgr_tmodel_active")) {
        modelid = l4iStorage.Get("l5smgr_tmodel_active");
    }

    if (!specid || !modelid) {
        return;
    }

    var uri = "specid="+ specid +"&modelid="+ modelid;
    if (document.getElementById("qry_text")) {
        uri = "&qry_text="+ $("#qry_text").val();
    }

    // console.log(uri);
    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", "data", function (tpl, rsj) {
            
            if (tpl) {
                $("#work-content").html(tpl);
            }

            l4iStorage.Set("l5smgr_tmodel_active", modelid);

            if (rsj === undefined || rsj.kind != "TermList" 
                || rsj.items === undefined || rsj.items.length < 1) {
                $("#l5smgr-nodels").empty();
                $("#l5smgr-termls").empty();
                return l4i.InnerAlert("#l5smgr-node-alert", 'alert-danger', "Item Not Found");
            }

            $("#l5smgr-node-alert").hide();

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
                    specid  : specid,
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
            l5sMgr.Ajax("-/term/list.tpl", {
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

        l5sMgr.Ajax("/v1/term/list?"+ uri, {
            callback: ep.done("data"),           
        });
    });
}


l5sTerm.Set = function(specid, modelid, termid)
{
    if (!specid && l4iStorage.Get("l5smgr_spec_active")) {
        specid = l4iStorage.Get("l5smgr_spec_active");
    }
    if (!modelid && l4iStorage.Get("l5smgr_tmodel_active")) {
        modelid = l4iStorage.Get("l5smgr_tmodel_active");
    }

    if (!specid || !modelid) {
        return;
    }

    var uri = "specid="+ specid +"&modelid="+ modelid;

    // console.log(uri);
    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", "data", function (tpl, data) {
            
            if (!tpl) {
                return; // TODO
            }

            $("#work-content").html(tpl);

            if (data === undefined || data.kind != "Term") {
                return l4i.InnerAlert("#l5smgr-node-alert", 'alert-danger', "Item Not Found");
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

            $("#l5smgr-node-alert").hide();

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

        l5sMgr.Ajax("-/term/set.tpl", {
            callback: function(err, tpl) {
                    
                if (err) {
                    return ep.emit('error', err);
                }
                ep.emit("tpl", tpl);
            }
        });

        if (termid) {
            l5sMgr.Ajax("/v1/term/entry?"+ uri +"&id="+ termid, {
                callback: ep.done("data"),           
            });
        } else {
            l5sMgr.Ajax("/v1/term-model/entry?"+ uri, {
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
    var req = {
        kind   : "Term",
        id     : parseInt($("#l5smgr-termset").find("input[name=id]").val()),
        title  : $("#l5smgr-termset").find("input[name=title]").val(),
        state  : parseInt($("#l5smgr-termset").find("input[name=state]").val()),
    }

    var model_type = $("#l5smgr-termset").find("input[name=model_type]").val();
    if (model_type = "taxonomy") {
        req.weight = parseInt($("#l5smgr-termset").find("input[name=weight]").val());
        req.pid = parseInt($("#l5smgr-termset").find("input[name=pid]").val());
    } else if (model_type = "tag") {

    }

    // console.log(req);

    //
    var uri = "specid="+ l4iStorage.Get("l5smgr_spec_active");
    uri += "&modelid="+ l4iStorage.Get("l5smgr_tmodel_active");

    l5sMgr.Ajax("/v1/term/set?"+ uri, {
        method : "POST",
        data   : JSON.stringify(req),
        callback : function(err, data) {

            if (data === undefined || data.kind != "Term") {
                return l4i.InnerAlert("#l5smgr-node-alert", 'alert-danger', data.error.message);
            }

            $("#l5smgr-termset").find("input[name=id]").val(data.id);

            l4i.InnerAlert("#l5smgr-node-alert", 'alert-success', "Successful operation");
            setTimeout(l5sTerm.List, 500);
        }
    });
}
