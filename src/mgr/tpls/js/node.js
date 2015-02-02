var l5sNode = {
    speclsCurrent: null,
    specCurrent: null,
    setCurrent: null,
}

l5sNode.Index = function()
{
    l5sMgr.Ajax(l5sMgr.base +"-/node/index.tpl", {
        callback: function(err, data) {

            $("#com-content").html(data);

            l5sMgr.Ajax("/v1/spec/list", {
                callback: function(err, data) {
                    
                    if (data === undefined || data.kind != "SpecList" 
                        || data.items === undefined || data.items.length < 1) {
                        return l4i.InnerAlert("#l5smgr-node-alert", 'alert-danger', "Content Type Not Found");
                    }

                    if (!l4iStorage.Get("l5smgr_node_active")) {
                        for (var i in data.items) {
                            l4iStorage.Set("l5smgr_node_active", data.items[i].metadata.id);
                            break;
                        }
                    }
                    if (!l4iStorage.Get("l5smgr_node_active")) {
                        // TODO
                        return;
                    }

                    for (var i in data.items) {
                        if (data.items[i].metadata.id == l4iStorage.Get("l5smgr_node_active")) {
                            l5sNode.specCurrent = data.items[i];
                            break;
                        }
                    }

                    if (!l5sNode.specCurrent) {
                        // TODO
                    }

                    // console.log(l5sNode.specCurrent);

                    if (l5sNode.specCurrent) {

                        // console.log(l5sNode.specCurrent);
                        
                        if (!l5sNode.specCurrent.nodeModels) {
                            l5sNode.specCurrent.nodeModels = [];
                        }
                        if (!l5sNode.specCurrent.termModels) {
                            l5sNode.specCurrent.termModels = [];
                        }

                        var node_model_active = null;
                       
                        for (var i in l5sNode.specCurrent.nodeModels) {
                                
                            if (!node_model_active) {
                                node_model_active = l5sNode.specCurrent.nodeModels[i].metadata.name;
                            }

                            if (l4iStorage.Get("l5smgr_node_model_active") == l5sNode.specCurrent.nodeModels[i].metadata.name) {
                                node_model_active = l5sNode.specCurrent.nodeModels[i].metadata.name;
                                break;
                            }
                        }

                        if (!node_model_active) {
                            return; // TODO
                        }

                        //
                        if (node_model_active != l4iStorage.Get("l5smgr_node_model_active")) {
                            l4iStorage.Set("l5smgr_node_model_active", node_model_active);
                        }

                        //
                        for (var i in l5sNode.specCurrent.nodeModels) {
                            // console.log(data.active +"."+ l5sNode.specCurrent.nodeModels[i].metadata.name)
                            if (node_model_active == l5sNode.specCurrent.nodeModels[i].metadata.name) {
                                l5sNode.List(l4iStorage.Get("l5smgr_node_active"), node_model_active);
                            }
                        }
                        // console.log(node_model_active);

                        l4iTemplate.Render({
                            dstid: "l5smgr-node-nmodels",
                            tplid: "l5smgr-node-nmodels-tpl",
                            data:  {
                                active: node_model_active,
                                items: l5sNode.specCurrent.nodeModels,
                            },
                        });
                    }

                    l5sNode.speclsCurrent = data.items;
                    // console.log(data);
                    l4iTemplate.Render({
                        dstid: "l5smgr-node-specls",
                        tplid: "l5smgr-node-specls-tpl",
                        data:  {
                            active : l4iStorage.Get("l5smgr_node_active"),
                            items  : data.items,
                        },
                    });
                },
            });
                
            // l5sNode.List();
        }
    });
}

l5sNode.List = function(specid, model)
{
    if (!specid && l4iStorage.Get("l5smgr_node_active")) {
        specid = l4iStorage.Get("l5smgr_node_active");
    }
    if (!model && l4iStorage.Get("l5smgr_node_model_active")) {
        model = l4iStorage.Get("l5smgr_node_model_active");
    }

    if (!specid || !model) {
        return;
    }

    var uri = "specid="+ specid +"&model="+ model;
    if (document.getElementById("qry_text")) {
        uri = "&qry_text="+ $("#qry_text").val();
    }

    // console.log(uri);
    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", "data", function (tpl, rsj) {
            
            if (tpl) {
                $("#work-content").html(tpl);
            }

            // console.log(rsj);
            // if (data typeof object)
            // var rsj = JSON.parse(data);

            if (rsj === undefined || rsj.kind != "NodeList" 
                || rsj.items === undefined || rsj.items.length < 1) {
                return l4i.InnerAlert("#l5smgr-node-alert", 'alert-danger', "Item Not Found");
            }

            $("#l5smgr-node-alert").hide();

            for (var i in rsj.items) {
                rsj.items[i].created = l4i.TimeParseFormat(rsj.items[i].created, "Y-m-d");
                rsj.items[i].updated = l4i.TimeParseFormat(rsj.items[i].updated, "Y-m-d");
            }

            l4iTemplate.Render({
                dstid: "l5smgr-nodels",
                tplid: "l5smgr-nodels-tpl",
                data:  {
                    specid : specid,
                    model  : model,
                    items  : rsj.items,
                },
            });
        });
    
        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        // template
        var el = document.getElementById("l5smgr-nodels");
        if (!el || el.length < 1) {
            l5sMgr.Ajax("-/node/list.tpl", {
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

        l5sMgr.Ajax("/v1/node/list?"+ uri, {
            callback: ep.done("data"),           
        });
    });
}

l5sNode.Set = function(specid, model, nodeid)
{
    if (!specid && l4iStorage.Get("l5smgr_node_active")) {
        specid = l4iStorage.Get("l5smgr_node_active");
    }
    if (!model && l4iStorage.Get("l5smgr_node_model_active")) {
        model = l4iStorage.Get("l5smgr_node_model_active");
    }

    if (!specid || !model) {
        return;
    }

    var uri = "specid="+ specid +"&model="+ model;

    // console.log(uri);
    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", "data", function (tpl, data) {
            
            if (!tpl) {
                return; // TODO
            }

            $("#work-content").html(tpl);

            if (data === undefined || data.kind != "Node") {
                return l4i.InnerAlert("#l5smgr-node-alert", 'alert-danger', "Item Not Found");
            }

            if (!data.state) {
                data.state = 1;
            }

            $("#l5smgr-node-alert").hide();

            l5sNode.setCurrent = data;

            l4iTemplate.Render({
                dstid: "l5smgr-nodeset",
                tplid: "l5smgr-nodeset-tpl",
                data:  data,
                success: function() {
                    for (var i in data.model.fields) {

                        var field = data.model.fields[i];

                        for (var j in data.fields) {
                            if (data.fields[i].name == field.name) {
                                field.value = data.fields[i].value;
                                break;
                            }
                        }

                        var tplid = null;

                        switch (field.type) {
                        case "text":
                            
                            if (field.attrs) {
                                for (var j in field.attrs) {
                                    field["attr_"+ field.attrs[i].key] = field.attrs[i].value;
                                }
                            }

                            if (!field.value) {
                                field.value = "";
                            }

                            tplid = "l5smgr-nodeset-tpltext";
                            break;
                        case "int8":
                        case "int16":
                        case "int32":
                        case "int64":
                        case "uint8":
                        case "uint16":
                        case "uint32":
                        case "uint64":

                            if (!field.value) {
                                field.value = "0";
                            }

                            tplid = "l5smgr-nodeset-tplint";
                            break;
                        default:
                            continue;
                        }

                        l4iTemplate.Render({
                            dstid  : "l5smgr-nodeset",
                            tplid  : tplid,
                            append : true,
                            data   : field,
                        });
                    }
                },
            });
        });
    
        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        l5sMgr.Ajax("-/node/set.tpl", {
            callback: function(err, tpl) {
                    
                if (err) {
                    return ep.emit('error', err);
                }
                ep.emit("tpl", tpl);
            }
        });

        if (nodeid) {
            l5sMgr.Ajax("/v1/node/entry?"+ uri +"&id="+ nodeid, {
                callback: ep.done("data"),           
            });
        } else {
            l5sMgr.Ajax("/v1/node-model/entry?"+ uri, {
                callback: function(err, data) {
                    ep.emit("data", {
                        kind  : "Node",
                        model : data,
                        id    : "",
                        title : "",
                    });
                },           
            });
        }
    });
}

l5sNode.SetCommit = function()
{
    if (!l5sNode.setCurrent) {
        return;
    }

    l5sNode.setCurrent.title = $("#l5smgr-nodeset").find("input[name=title]").val();

    var vals = [];

    // console.log()
    for (var i in l5sNode.setCurrent.model.fields) {

        var field = l5sNode.setCurrent.model.fields[i];

        var val = null;

        switch (field.type) {

        case "text":
            val = $("#l5smgr-nodeset").find("textarea[name="+ field.name +"]").val();
            break;
        case "int8":
        case "int16":
        case "int32":
        case "int64":
        case "uint8":
        case "uint16":
        case "uint32":
        case "uint64":
            val = $("#l5smgr-nodeset").find("input[name="+ field.name +"]").val();
            break;
        }
        
        if (val) {
            vals.push({name: field.name, value: val});
        }
    }   

    var req = {
        id     : $("#l5smgr-nodeset").find("input[name=id]").val(),
        title  : $("#l5smgr-nodeset").find("input[name=title]").val(),
        state  : parseInt($("#l5smgr-nodeset").find("input[name=state]").val()),
        fields : vals,
    }

    // console.log(JSON.stringify(req));

    //
    var uri = "specid="+ l4iStorage.Get("l5smgr_node_active");
    uri += "&model="+ l4iStorage.Get("l5smgr_node_model_active");

    l5sMgr.Ajax("/v1/node/set?"+ uri, {
        method: "POST",
        data : JSON.stringify(req),
        callback : function(err, data) {

            if (data === undefined || data.kind != "Node") {
                return l4i.InnerAlert("#l5smgr-node-alert", 'alert-danger', data.error.message);
            }

            // console.log(data.id);
            $("#l5smgr-nodeset").find("input[name=id]").val(data.id);

            l4i.InnerAlert("#l5smgr-node-alert", 'alert-success', "Successful operation");
            setTimeout(l5sNode.List, 1000);
        }
    });
    // console.log(vals);
}