var l5sSpec = {

    specdef : {
        kind : "Spec",
        meta : {
            id   : "",
            name : "",
        },
        title : "",
    },
}

l5sSpec.Init = function()
{
    l4i.UrlEventRegister("spec/index", l5sSpec.Index);
}

l5sSpec.Index = function()
{
    l5sMgr.TplCmd("spec/index", {
        callback: function(err, data) {
                
            $("#com-content").html(data);
                
            l5sSpec.List();
        }
    });
}

l5sSpec.List = function()
{

    var uri = "";
    if (document.getElementById("qry_text")) {
        uri = "qry_text="+ $("#qry_text").val();
    }

    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", "data", function (tpl, rsj) {
            
            if (tpl) {
                $("#work-content").html(tpl);
            }
            // console.log(tpl);
            // if (data typeof object)
            // var rsj = JSON.parse(data);

            if (rsj === undefined || rsj.kind != "SpecList" 
                || rsj.items === undefined || rsj.items.length < 1) {
                return l4i.InnerAlert("#l5smgr-specls-alert", 'alert-danger', "Item Not Found");
            }

            $("#l5smgr-specls-alert").hide();

            for (var i in rsj.items) {
                
                if (rsj.items[i].nodeModels) {
                    rsj.items[i]._nodeModelsNum = rsj.items[i].nodeModels.length;
                } else {
                    rsj.items[i]._nodeModelsNum = 0;
                }

                if (rsj.items[i].termModels) {
                    rsj.items[i]._termModelsNum = rsj.items[i].termModels.length;
                } else {
                    rsj.items[i]._termModelsNum = 0;
                }

                if (rsj.items[i].actions) {
                    rsj.items[i]._actionsNum = rsj.items[i].actions.length;
                } else {
                    rsj.items[i]._actionsNum = 0;
                }

                if (!rsj.items[i].meta.created) {
                    rsj.items[i].meta.created = rsj.items[i].meta.updated;
                }
            }

            l4iTemplate.Render({
                dstid: "l5smgr-specls",
                tplid: "l5smgr-specls-tpl",
                data:  {
                    items: rsj.items,
                },
            });
        });
    
        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-speclist)");
        });

        // template
        var el = document.getElementById("l5smgr-specls");
        if (!el || el.length < 1) {
            l5sMgr.TplCmd("spec/list", {
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

        // l5sMgr.Ajax("-/spec/list.tpl", {
        //     callback: ep.done("tpl"),
        // });
    
        l5sMgr.ApiCmd("mod-set/spec-list?"+ uri, {
            callback: ep.done("data"),           
        });
    });
}

l5sSpec.InfoSet = function(name)
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            if (!data || !data.kind || data.kind != "Spec") {
                return alert("Spec Not Found");
            }

            var ptitle = "Info Settings";
            if (!name) {
                ptitle = "New Module";
            }

            l4iModal.Open({
                tplsrc : tpl,
                width  : 600,
                height : 300,
                title  : ptitle,
                data   : data,
                success : function() {
                    
                },
                buttons : [{
                    onclick : "l4iModal.Close()",
                    title   : "Close",
                }, {
                    onclick : "l5sSpec.InfoSetCommit()",
                    title   : "Save",
                    style   : "btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error, Please try again later "+ err);
        });

        l5sMgr.TplCmd("spec/info-set", {
            callback: ep.done('tpl'),
        });

        if (name) {
            
            l5sMgr.ApiCmd("mod-set/spec-entry?name="+ name, {
                callback: ep.done('data'),
            });
        } else {

            ep.emit("data", l4i.Clone(l5sSpec.specdef));
        }
        
    });

}

l5sSpec.InfoSetCommit = function()
{
    var form = $("#l5smgr-specset");
    var alertid = "#l5smgr-specset-alert";

    var req = {
        meta : {
            name : form.find("input[name=name]").val(),
        },
        title : form.find("input[name=title]").val(),
    };

    // console.log(req);

    l5sMgr.ApiCmd("mod-set/spec-info-set", {
        method  : "PUT",
        data    : JSON.stringify(req),
        success : function(data) {

            // console.log(data);

            if (!data || data.error || data.kind != "Spec") {

                if (data.error) {
                    return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
                }

                return l4i.InnerAlert(alertid, 'alert-danger', 'Network Connection Exception');
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful updated");
            
            l5sSpec.List();

            window.setTimeout(function() {
                l4iModal.Close();
            }, 1000);
        },
    });
}