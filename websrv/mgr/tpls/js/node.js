var l5sNode = {
    speclsCurrent: [],
    specCurrent: null,
    setCurrent: null,
    cmEditor: null,
    cmEditors: {},
    general_onoff : [{
        type: true,
        name: "ON",
    },{
        type: false,
        name: "OFF",
    }],
    status_def : [{
        type: 1,
        name: "Publish",
    },{
        type: 2,
        name: "Draft",
    },{
        type: 3,
        name: "Private",
    }],
}

l5sNode.Init = function()
{
    l4i.UrlEventRegister("node/index", l5sNode.Index);
}

l5sNode.Index = function()
{
    l4iStorage.Set("l5smgr_nav_last_active", "node/index");
    
    var alertid = "#l5smgr-node-alert";

    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", "data", function (tpl, data) {
            
            if (tpl) {
                $("#com-content").html(tpl);
            }

            // console.log(data);

            //
            if (l5sNode.speclsCurrent.length < 1) {

                if (!data || data.kind != "SpecList" 
                    || !data.items || data.items.length < 1) {
                    return l4i.InnerAlert(alertid, 'alert-danger', "Content Type Not Found");
                }

                for (var i in data.items) {

                    if (!data.items[i].nodeModels || data.items[i].nodeModels.length < 1) {
                        continue;
                    }

                    l5sNode.speclsCurrent.push(data.items[i]);
                }
            }

            //
            if (!l4iStorage.Get("l5smgr_spec_active")) {
                for (var i in l5sNode.speclsCurrent) {
                    l4iStorage.Set("l5smgr_spec_active", l5sNode.speclsCurrent[i].meta.name);
                    break;
                }
            }
            if (!l4iStorage.Get("l5smgr_spec_active")) {
                // TODO
                return;
            }


            for (var i in l5sNode.speclsCurrent) {


                if (l5sNode.speclsCurrent[i].meta.name == l4iStorage.Get("l5smgr_spec_active")) {
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
                        node_model_active = l5sNode.specCurrent.nodeModels[i].meta.name;
                    }

                    if (l4iStorage.Get("l5smgr_nmodel_active") == l5sNode.specCurrent.nodeModels[i].meta.name) {
                        node_model_active = l5sNode.specCurrent.nodeModels[i].meta.name;
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
                    if (node_model_active == l5sNode.specCurrent.nodeModels[i].meta.name) {
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
            l5sMgr.TplCmd("node/index", {
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

        if (l5sNode.speclsCurrent.length > 0) {
            ep.emit("data", null);
        } else {

            l5sMgr.ApiCmd("mod-set/spec-list", {
                callback: ep.done("data"),           
            });
        }
    });
}

l5sNode.List = function(modname, modelid)
{
    var alertid = "#l5smgr-node-alert",
        page = 0;

    if (!modname && l4iStorage.Get("l5smgr_spec_active")) {
        modname = l4iStorage.Get("l5smgr_spec_active");
    }

    if (!modelid && l4iStorage.Get("l5smgr_nmodel_active")) {
        modelid = l4iStorage.Get("l5smgr_nmodel_active");
    }

    if (l4iStorage.Get("l5smgr_nodels_page")) {
        page = l4iStorage.Get("l5smgr_nodels_page");
    }

    if (!modname || !modelid) {
        return;
    }

    var uri = "modname="+ modname +"&modelid="+ modelid +"&page="+ page;
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

            if (!rsj || rsj.kind != "NodeList" 
                || !rsj.items || rsj.items.length < 1) {
                $("#l5smgr-nodels").empty();
                $("#l5smgr-termls").empty();
                return l4i.InnerAlert(alertid, 'alert-danger', "Item Not Found");
            }

            $(alertid).hide();

            for (var i in rsj.items) {
                rsj.items[i].created = l4i.TimeParseFormat(rsj.items[i].created, "Y-m-d");
                rsj.items[i].updated = l4i.TimeParseFormat(rsj.items[i].updated, "Y-m-d");

                if (!rsj.items[i].ext_access_counter) {
                    rsj.items[i].ext_access_counter = 0;
                }

                if (!rsj.items[i].ext_permalink_name) {
                    rsj.items[i].ext_permalink_name = "";
                }
            }

            l4iTemplate.Render({
                dstid: "l5smgr-nodels",
                tplid: "l5smgr-nodels-tpl",
                data:  {
                    model : rsj.model,
                    modname : modname,
                    modelid  : modelid,
                    items  : rsj.items,
                    _status_def : l5sNode.status_def,
                },
                success: function() {

                    rsj.meta.RangeLen = 20;

                    l4iTemplate.Render({
                        dstid  : "l5smgr-nodels-pager",
                        tplid  : "l5smgr-nodels-pager-tpl",
                        data   : l4i.Pager(rsj.meta),
                    });
                }
            });
        });
    
        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        // template
        var el = document.getElementById("l5smgr-nodels");
        if (!el || el.length < 1) {
            l5sMgr.TplCmd("node/list", {
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

        l5sMgr.ApiCmd("node/list?"+ uri, {
            callback: ep.done("data"),           
        });
    });
}

l5sNode.ListPage = function(page)
{
    l4iStorage.Set("l5smgr_nodels_page", parseInt(page));
    l5sNode.List();
}

l5sNode.Set = function(modname, modelid, nodeid)
{
    var alertid = "#l5smgr-node-alert";

    if (!modname && l4iStorage.Get("l5smgr_spec_active")) {
        modname = l4iStorage.Get("l5smgr_spec_active");
    }
    if (!modelid && l4iStorage.Get("l5smgr_nmodel_active")) {
        modelid = l4iStorage.Get("l5smgr_nmodel_active");
    }

    // console.log(modname +","+ modelid +","+ nodeid);

    if (!modname || !modelid) {
        return;
    }

    l5sEditor.Clean();

    var uri = "modname="+ modname +"&modelid="+ modelid;

    // console.log(uri);
    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", "data", function (tpl, data) {

            if (!tpl) {
                return; // TODO
            }

            $("#work-content").html(tpl);

            if (!data || data.kind != "Node") {
                return l4i.InnerAlert(alertid, 'alert-danger', "Item Not Found");
            }

            if (!data.status) {
                data.status = 1;
            }

            if (!data.model.terms) {
                data.model.terms = [];
            }

            if (!data.ext_comment_perentry) {
                data.ext_comment_perentry = false;
            }

            if (!data.ext_permalink_name) {
                data.ext_permalink_name = "";
            }

            $(alertid).hide();

            l5sNode.setCurrent = data;
            data._status_def = l5sNode.status_def;

            // console.log(data);

            l4iTemplate.Render({
                dstid: "l5smgr-nodeset",
                tplid: "l5smgr-nodeset-tpl",
                data:  data,
                success: function() {
                    
                    for (var i in data.model.fields) {

                        var field = data.model.fields[i];

                        var field_entry = {};

                        for (var j in data.fields) {
                            if (data.fields[i].name == field.name) {
                                field_entry = data.fields[i];
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

                            if (field_entry.attrs) {
                                for (var j in field_entry.attrs) {
                                    field["attr_"+ field_entry.attrs[j].key] = field_entry.attrs[j].value;
                                }
                            }

                            if (!field.value) {
                                field.value = "";
                            }

                            if (!field.attr_format) {
                                field.attr_format = "text";
                            }

                            cb = function() {
                                l5sEditor.Open(field.name, field.attr_format);
                            };                            

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
                            dstid  : "l5smgr-nodeset-fields",
                            tplid  : tplid,
                            append : true,
                            data   : field,
                            success: cb,
                        });
                    }

                    for (var i in data.model.terms) {

                        var term = data.model.terms[i];

                        for (var j in data.terms) {
                            if (data.terms[i].name == term.meta.name) {
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
                                dstid  : "l5smgr-nodeset-fields",
                                tplid  : tplid,
                                append : true,
                                data   : term,
                            });

                            break;

                        case "taxonomy":

                            if (!term.value) {
                                term.value = "0";
                            }


                            l5sMgr.ApiCmd("term/list?modname="+ modname +"&modelid="+ term.meta.name, {
                                callback: function(err, data) {
                                    
                                    if (data.kind != "TermList") {
                                        return;
                                    }

                                    data.item = term;

                                    tplid = "l5smgr-nodeset-tplterm_taxonomy";
                        
                                    l4iTemplate.Render({
                                        dstid  : "l5smgr-nodeset-fields",
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


                    if (data.model.extensions.comment_perentry) {
                        l4iTemplate.Render({
                            dstid  : "l5smgr-nodeset-fields",
                            tplid  : "l5smgr-nodeset-tplext_comment_perentry",
                            append : true,
                            data   : {
                                _general_onoff: l5sNode.general_onoff,
                                ext_comment_perentry: data.ext_comment_perentry,
                            },
                        });
                    }

                    if (data.model.extensions.permalink && data.model.extensions.permalink != "") {
                        l4iTemplate.Render({
                            dstid  : "l5smgr-nodeset-tops",
                            tplid  : "l5smgr-nodeset-tplext_permalink",
                            append : true,
                            data   : {
                                ext_permalink_name: data.ext_permalink_name,
                            },
                        });
                    }
                },
            });
        });

        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        l5sMgr.TplCmd("node/set", {
            callback: function(err, tpl) {
                    
                if (err) {
                    return ep.emit('error', err);
                }
                ep.emit("tpl", tpl);
            }
        });

        if (nodeid) {
            l5sMgr.ApiCmd("node/entry?"+ uri +"&id="+ nodeid, {
                callback: ep.done("data"),           
            });
        } else {
            l5sMgr.ApiCmd("node-model/entry?"+ uri, {
                callback: function(err, data) {
                    ep.emit("data", {
                        kind  : "Node",
                        model : data,
                        id    : "",
                        title : "",
                        ext_comment_perentry: true,
                    });
                },           
            });
        }
    });
}


l5sNode.SetCommit = function()
{
    var form = $("#l5smgr-nodeset"),
        alertid = "#l5smgr-node-alert";

    if (!l5sNode.setCurrent) {
        return;
    }

    l5sNode.setCurrent.title = form.find("input[name=title]").val();

    var req = {
        id     : form.find("input[name=id]").val(),
        title  : form.find("input[name=title]").val(),
        status : parseInt(form.find("select[name=status]").val()),
        fields : [],
        terms  : [],
        ext_comment_perentry: form.find("select[name=ext_comment_perentry]").val(),
        ext_permalink_name: form.find("input[name=ext_permalink_name]").val(),
    }

    if (req.ext_comment_perentry && req.ext_comment_perentry == "false") {
        req.ext_comment_perentry = false;
    } else {
        req.ext_comment_perentry = true;
    }

    // console.log("DDD");
    // console.log(req);
    for (var i in l5sNode.setCurrent.model.fields) {

        var field = l5sNode.setCurrent.model.fields[i];

        var field_set = {
            name: field.name,
            value: null,
            attrs: [],
        };

        switch (field.type) {

        case "text":

            var format = form.find("input[name=field_"+ field.name +"_attr_format]").val();
            if (!format) {
                format = "text";
            }

            field_set.attrs.push({key: "format", value: format});
            field_set.value = l5sEditor.Content(field.name);

            // console.log(format);

            // if (field.attrs) {
                
            //     for (var j in field.attrs) {

            //         if (field.attrs[j].key == "format" && field.attrs[j].value == "md") {
                        
            //             field_set.value = l5sEditor.Content(field.name);
                        
            //             field_set.attrs.push({key: "format", value: "md"});

            //             break;
            //         }
            //     }
            // }

            // if (!field_set.value) {
            //     field_set.value = form.find("textarea[name=field_"+ field.name +"]").val();
            // }

            break;
        
        case "string":
            field_set.value = form.find("input[name=field_"+ field.name +"]").val();
            break;

        case "int8":
        case "int16":
        case "int32":
        case "int64":
        case "uint8":
        case "uint16":
        case "uint32":
        case "uint64":
            field_set.value = form.find("input[name=field_"+ field.name +"]").val();
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
            val = form.find("input[name=term_"+ term.meta.name +"]").val();
            break;
        case "taxonomy":
            val = form.find("select[name=term_"+ term.meta.name +"]").val();
            break;
        }

        if (val) {
            req.terms.push({name: term.meta.name, value: val});
        }
    }

    // console.log(l5sNode.setCurrent.model.terms);
    // console.log(JSON.stringify(req));

    //
    var uri = "modname="+ l4iStorage.Get("l5smgr_spec_active");
    uri += "&modelid="+ l4iStorage.Get("l5smgr_nmodel_active");

    l5sMgr.ApiCmd("node/set?"+ uri, {
        method : "POST",
        data   : JSON.stringify(req),
        callback : function(err, data) {

            if (!data || data.kind != "Node") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            // console.log(data.id);
            form.find("input[name=id]").val(data.id);

            l4i.InnerAlert(alertid, 'alert-success', "Successful operation");
            setTimeout(function() {
                l5sNode.List();
                l5sEditor.Clean();
            }, 500);
        }
    });
}


l5sNode.Del = function(modname, modelid, id)
{
    console.log(modname, modelid, id);
    l4iModal.Open({
        title  : "Delete",
        tplsrc : '<div id="l5smgr-node-del" class="alert alert-danger">Are you sure to delete this?</div>',
        height : "200px",
        buttons: [{
            title: "Confirm to delete",
            onclick : 'l5sNode.DelCommit("'+modname+'","'+modelid+'","'+id+'")',
            style: "btn-danger",
        },{
            title: "Cancel",
            onclick : "l4iModal.Close()",
        }],
    });
}

l5sNode.DelCommit = function(modname, modelid, id)
{
    var alertid = "#l5smgr-node-del";
    var uri = "modname="+ modname + "&modelid="+ modelid +"&id="+ id;

    l5sMgr.ApiCmd("node/del?"+ uri, {
        callback : function(err, data) {

            if (!data || data.kind != "Node") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful deleted");
            setTimeout(function() {
                l5sNode.List();
                l4iModal.Close();
            }, 500);
        }
    });
}
