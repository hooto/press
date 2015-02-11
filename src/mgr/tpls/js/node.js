var l5sNode = {
    speclsCurrent: null,
    specCurrent: null,
    setCurrent: null,
    cmEditor: null,
    cmEditors: {},
}

l5sNode.Index = function()
{
    // console.log(uri);
    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", "data", function (tpl, data) {
            
            if (tpl) {
                $("#com-content").html(tpl);
            }

            if (!l5sNode.speclsCurrent) {

                if (!data || data.kind != "SpecList" 
                    || data.items === undefined || data.items.length < 1) {
                    return l4i.InnerAlert("#l5smgr-node-alert", 'alert-danger', "Content Type Not Found");
                }

                l5sNode.speclsCurrent = data.items;
            }

            //
            if (!l4iStorage.Get("l5smgr_spec_active")) {
                for (var i in l5sNode.speclsCurrent) {
                    l4iStorage.Set("l5smgr_spec_active", l5sNode.speclsCurrent[i].metadata.id);
                    break;
                }
            }
            if (!l4iStorage.Get("l5smgr_spec_active")) {
                // TODO
                return;
            }


            for (var i in l5sNode.speclsCurrent) {
                if (l5sNode.speclsCurrent[i].metadata.id == l4iStorage.Get("l5smgr_spec_active")) {
                    l5sNode.specCurrent = l5sNode.speclsCurrent[i];
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

                    if (l4iStorage.Get("l5smgr_nmodel_active") == l5sNode.specCurrent.nodeModels[i].metadata.name) {
                        node_model_active = l5sNode.specCurrent.nodeModels[i].metadata.name;
                        break;
                    }
                }

                // console.log(l4iStorage.Get("l5smgr_nmodel_active"));

                if (!node_model_active) {
                    return; // TODO
                }

                        //
                if (node_model_active != l4iStorage.Get("l5smgr_nmodel_active")) {
                    l4iStorage.Set("l5smgr_nmodel_active", node_model_active);
                }

                //
                for (var i in l5sNode.specCurrent.nodeModels) {
                    if (node_model_active == l5sNode.specCurrent.nodeModels[i].metadata.name) {
                        l5sNode.List(l4iStorage.Get("l5smgr_spec_active"), node_model_active);
                    }
                }

                l4iTemplate.Render({
                    dstid: "l5smgr-node-nmodels",
                    tplid: "l5smgr-node-nmodels-tpl",
                    data:  {
                        active: node_model_active,
                        items: l5sNode.specCurrent.nodeModels,
                    },
                });

                l4iTemplate.Render({
                    dstid: "l5smgr-node-tmodels",
                    tplid: "l5smgr-node-tmodels-tpl",
                    data:  {
                        items: l5sNode.specCurrent.termModels,
                    },
                });
            }

            var el = document.getElementById("l5smgr-node-specls");
            if (!el || !el.length || el.length < 1) {

                l4iTemplate.Render({
                    dstid: "l5smgr-node-specls",
                    tplid: "l5smgr-node-specls-tpl",
                    data:  {
                        active : l4iStorage.Get("l5smgr_spec_active"),
                        items  : l5sNode.speclsCurrent,
                    },
                });
            }
        });
    
        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        // template
        var el = document.getElementById("l5smgr-node-specls");
        if (!el || !el.length || el.length < 1) {
            l5sMgr.Ajax("-/node/index.tpl", {
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

        if (l5sNode.speclsCurrent) {
            ep.emit("data", null);
        } else {

            l5sMgr.Ajax("/v1/spec/list", {
                callback: ep.done("data"),           
            });
        }
    });
}


l5sNode.List = function(specid, modelid)
{
    if (!specid && l4iStorage.Get("l5smgr_spec_active")) {
        specid = l4iStorage.Get("l5smgr_spec_active");
    }
    if (!modelid && l4iStorage.Get("l5smgr_nmodel_active")) {
        modelid = l4iStorage.Get("l5smgr_nmodel_active");
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

            // console.log(rsj);
            // if (data typeof object)
            // var rsj = JSON.parse(data);

            if (rsj === undefined || rsj.kind != "NodeList" 
                || rsj.items === undefined || rsj.items.length < 1) {
                $("#l5smgr-nodels").empty();
                $("#l5smgr-termls").empty();
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
                    modelid  : modelid,
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

l5sNode.Set = function(specid, modelid, nodeid)
{
    if (!specid && l4iStorage.Get("l5smgr_spec_active")) {
        specid = l4iStorage.Get("l5smgr_spec_active");
    }
    if (!modelid && l4iStorage.Get("l5smgr_nmodel_active")) {
        modelid = l4iStorage.Get("l5smgr_nmodel_active");
    }

    // console.log(specid +","+ modelid +","+ nodeid);

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

            if (data === undefined || data.kind != "Node") {
                return l4i.InnerAlert("#l5smgr-node-alert", 'alert-danger', "Item Not Found");
            }

            if (!data.state) {
                data.state = 1;
            }

            if (!data.model.terms) {
                data.model.terms = [];
            }

            $("#l5smgr-node-alert").hide();

            l5sNode.setCurrent = data;
            // console.log(data);

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
                        var cb = null;

                        switch (field.type) {

                        case "string":

                            if (!field.value) {
                                field.value = "";
                            }

                            tplid = "l5smgr-nodeset-tplstring";
                            break;

                        case "text":
                            
                            if (field.attrs) {
                                for (var j in field.attrs) {
                                    field["attr_"+ field.attrs[j].key] = field.attrs[j].value;
                                }
                            }

                            if (!field.value) {
                                field.value = "";
                            }

                            if (field.attr_format && field.attr_format == "md") {
                                cb = function() {
                                    l5sEditor.Open(field.name);
                                }
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
                            success: cb,
                        });
                    }

                    for (var i in data.model.terms) {

                        var term = data.model.terms[i];

                        for (var j in data.terms) {
                            if (data.terms[i].name == term.metadata.name) {
                                term.value = data.terms[i].value;
                                break;
                            }
                        }

                        var tplid = null;

                        switch (term.type) {

                        case "tag":

                            if (!term.value) {
                                term.value = "";
                            }

                            tplid = "l5smgr-nodeset-tplterm_tag";

                            l4iTemplate.Render({
                                dstid  : "l5smgr-nodeset",
                                tplid  : tplid,
                                append : true,
                                data   : term,
                            });

                            break;

                        case "taxonomy":

                            if (!term.value) {
                                term.value = "0";
                            }


                            l5sMgr.Ajax("/v1/term/list?specid="+ specid +"&modelid="+ term.metadata.name, {
                                callback: function(err, data) {
                                    
                                    if (data.kind != "TermList") {
                                        return;
                                    }

                                    data.item = term;

                                    tplid = "l5smgr-nodeset-tplterm_taxonomy";
                        
                                    l4iTemplate.Render({
                                        dstid  : "l5smgr-nodeset",
                                        tplid  : tplid,
                                        append : true,
                                        data   : data,
                                    });
                                },           
                            });                            

                            break;

                        default:
                            continue;
                        }
                        
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

    var req = {
        id     : $("#l5smgr-nodeset").find("input[name=id]").val(),
        title  : $("#l5smgr-nodeset").find("input[name=title]").val(),
        state  : parseInt($("#l5smgr-nodeset").find("input[name=state]").val()),
        fields : [],
        terms  : [],
    }

    // console.log()
    for (var i in l5sNode.setCurrent.model.fields) {

        var field = l5sNode.setCurrent.model.fields[i];

        var field_set = {
            name: field.name,
            value: null,
            attrs: [],
        };

        switch (field.type) {

        case "text":

            if (field.attrs) {
                for (var j in field.attrs) {

                    if (field.attrs[j].key == "format" && field.attrs[j].value == "md") {
                        
                        field_set.value = l5sEditor.Content(field.name);
                        
                        field_set.attrs.push({key: "format", value: "md"});

                        break;
                    }
                }
            }

            if (!field_set.value) {
                field_set.value = $("#l5smgr-nodeset").find("textarea[name=field_"+ field.name +"]").val();
            }

            break;
        
        case "string":
            field_set.value = $("#l5smgr-nodeset").find("input[name=field_"+ field.name +"]").val();
            break;

        case "int8":
        case "int16":
        case "int32":
        case "int64":
        case "uint8":
        case "uint16":
        case "uint32":
        case "uint64":
            field_set.value = $("#l5smgr-nodeset").find("input[name=field_"+ field.name +"]").val();
            break;

        }
        
        if (field_set.value) {
            req.fields.push(field_set);
        }
    }

    for (var i in l5sNode.setCurrent.model.terms) {
        
        var term = l5sNode.setCurrent.model.terms[i];

        var val = null;

        switch (term.type) {

        case "tag":
            val = $("#l5smgr-nodeset").find("input[name=term_"+ term.metadata.name +"]").val();
            break;
        case "taxonomy":
            val = $("#l5smgr-nodeset").find("select[name=term_"+ term.metadata.name +"]").val();
            break;
        }

        if (val) {
            req.terms.push({name: term.metadata.name, value: val});
        }
    }

    // console.log(l5sNode.setCurrent.model.terms);
    // console.log(JSON.stringify(req));
    // console.log(req);

    //
    var uri = "specid="+ l4iStorage.Get("l5smgr_spec_active");
    uri += "&modelid="+ l4iStorage.Get("l5smgr_nmodel_active");

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
            setTimeout(l5sNode.List, 500);
        }
    });
}