var htapNode = {
    navPrefix: "node/index/",
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
    nodeOpToolsRefreshCurrent : null,
}

htapNode.Init = function(cb)
{
    htapNode.navRefresh(cb);
}

htapNode.navRefresh = function(cb)
{
    cb = cb || function(){};

    if (htapNode.speclsCurrent.length > 0) {

        // if (!l4iStorage.Get("htapm_spec_active")) {
        //     for (var i in htapNode.speclsCurrent) {
        //         l4iStorage.Set("htapm_spec_active", htapNode.speclsCurrent[i].meta.name);
        //         break;
        //     }
        // }

        // if (!l4iStorage.Get("htapm_spec_active")) {
        //     return cb();
        // }

        // console.log(htapNode.speclsCurrent);

        l4iTemplate.Render({
            dstid: "htapm-topbar-nav-node-specls",
            tplid: "htapm-topbar-nav-node-specls-tpl",
            data:  {
                active : l4iStorage.Get("htapm_spec_active"),
                items  : htapNode.speclsCurrent,
            },
        });

        return cb();
    }

    htapMgr.ApiCmd("mod-set/spec-list", {
        callback: function(err, data) {

            if (err || data.error || data.kind != "SpecList") {
                return cb();
            }

            //
            for (var i in data.items) {
                htapNode.speclsCurrent.push(data.items[i]);
                l4i.UrlEventRegister(
                    htapNode.navPrefix + data.items[i].meta.name,
                    htapNode.Index,
                    "htapm-topbar"
                );
            }

            //
            if (!l4iStorage.Get("htapm_spec_active")) {
                for (var i in htapNode.speclsCurrent) {
                    l4iStorage.Set("htapm_spec_active", htapNode.speclsCurrent[i].meta.name);
                    break;
                }
            }
            if (!l4iStorage.Get("htapm_spec_active")) {
                return cb();
            }

            l4iTemplate.Render({
                dstid: "htapm-topbar-nav-node-specls",
                tplid: "htapm-topbar-nav-node-specls-tpl",
                data:  {
                    // active : l4iStorage.Get("htapm_spec_active"),
                    items  : htapNode.speclsCurrent,
                },
            });

            cb();
        },
    });
}

htapNode.OpToolsRefresh = function(div_target)
{
    if (typeof div_target == "string" && div_target == htapNode.nodeOpToolsRefreshCurrent) {
        return;
    }

    $("#htapm-node-optools").empty();

    if (typeof div_target == "string") {

        var opt = $("#work-content").find(div_target);
        if (opt) {
            $("#htapm-node-optools").html(opt.html());
            htapNode.nodeOpToolsRefreshCurrent = div_target;
        }
    }
}

htapNode.Index = function(nav_href)
{
    if (!nav_href || nav_href.length <= htapNode.navPrefix.length) {
        return;
    }

    if (htapNode.speclsCurrent.length < 1) {
        return;
    }

    l4iStorage.Del("htapm_nodels_page");
    l4iStorage.Del("htapm_termls_page");

    htapNode.nodeOpToolsRefreshCurrent = null;
    l4iStorage.Set("htapm_nav_last_active", nav_href);
    l4iStorage.Set("htapm_spec_active", nav_href.substr(htapNode.navPrefix.length));

    var alertid = "#htapm-node-alert";

    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", function (tpl) {

            if (tpl) {
                $("#com-content").html(tpl);
            }

            var current = null;

            for (var i in htapNode.speclsCurrent) {

                if (htapNode.speclsCurrent[i].meta.name == l4iStorage.Get("htapm_spec_active")) {
                    current = htapNode.speclsCurrent[i];
                    break;
                }
            }

            if (!current) {
                return;
            }

            htapNode.specCurrent = current;

            if (!htapNode.specCurrent.nodeModels) {
                htapNode.specCurrent.nodeModels = [];
            }
            if (!htapNode.specCurrent.termModels) {
                htapNode.specCurrent.termModels = [];
            }

            var node_model_active = null;

            for (var i in htapNode.specCurrent.nodeModels) {

                if (!node_model_active) {
                    node_model_active = htapNode.specCurrent.nodeModels[i].meta.name;
                }

                if (l4iStorage.Get("htapm_nmodel_active") == htapNode.specCurrent.nodeModels[i].meta.name) {
                    node_model_active = htapNode.specCurrent.nodeModels[i].meta.name;
                    break;
                }
            }

            // console.log(l4iStorage.Get("htapm_nmodel_active"));

            if (!node_model_active) {
                return; // TODO
            }

                    //
            if (node_model_active != l4iStorage.Get("htapm_nmodel_active")) {
                l4iStorage.Set("htapm_nmodel_active", node_model_active);
            }

            //
            for (var i in htapNode.specCurrent.nodeModels) {
                if (node_model_active == htapNode.specCurrent.nodeModels[i].meta.name) {
                    htapNode.List(l4iStorage.Get("htapm_spec_active"), node_model_active);
                }
            }

            if (htapNode.specCurrent.nodeModels.length > 0) {
                l4iTemplate.Render({
                    dstid: "htapm-node-nmodels",
                    tplid: "htapm-node-nmodels-tpl",
                    data:  {
                        active: node_model_active,
                        items: htapNode.specCurrent.nodeModels,
                    },
                });
            } else {
                $("#htapm-node-nmodels").addClass("htapm-hide");
            }

            if (htapNode.specCurrent.termModels.length > 0) {
                l4iTemplate.Render({
                    dstid: "htapm-node-tmodels",
                    tplid: "htapm-node-tmodels-tpl",
                    data:  {
                        items: htapNode.specCurrent.termModels,
                    },
                });
            } else {
                $("#htapm-node-tmodels").addClass("htapm-hide");
            }
        });

        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        // template
        var el = document.getElementById("htapm-node-nmodels");
        if (!el || !el.length || el.length < 1) {
            htapMgr.TplCmd("node/index", {
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
    });
}

htapNode.List = function(modname, modelid)
{
    var alertid = "#htapm-node-alert",
        page = 0;

    if (!modname && l4iStorage.Get("htapm_spec_active")) {
        modname = l4iStorage.Get("htapm_spec_active");
    }

    if (!modelid && l4iStorage.Get("htapm_nmodel_active")) {
        modelid = l4iStorage.Get("htapm_nmodel_active");
    }

    if (l4iStorage.Get("htapm_nodels_page")) {
        page = l4iStorage.Get("htapm_nodels_page");
    }

    if (!modname || !modelid) {
        return;
    }

    var uri = "modname="+ modname +"&modelid="+ modelid +"&page="+ page;
    if (document.getElementById("qry_text")) {
        uri += "&qry_text="+ $("#qry_text").val();
    }

    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", "data", function (tpl, rsj) {

            if (tpl) {
                $("#work-content").html(tpl);
            }

            if (!rsj || rsj.kind != "NodeList"
                || !rsj.items || rsj.items.length < 1) {

                $("#htapm-nodels").empty();
                $("#htapm-termls").empty();
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
                dstid: "htapm-nodels",
                tplid: "htapm-nodels-tpl",
                data:  {
                    model   : rsj.model,
                    modname : modname,
                    modelid : modelid,
                    items   : rsj.items,
                    _status_def : htapNode.status_def,
                },
                success: function() {

                    rsj.meta.RangeLen = 20;

                    l4iTemplate.Render({
                        dstid : "htapm-nodels-pager",
                        tplid : "htapm-nodels-pager-tpl",
                        data  : l4i.Pager(rsj.meta),
                    });

                    htapNode.OpToolsRefresh("#htapm-node-list-opts");
                }
            });
        });

        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        // template
        var el = document.getElementById("htapm-nodels");
        if (!el || el.length < 1) {
            htapMgr.TplCmd("node/list", {
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

        htapMgr.ApiCmd("node/list?"+ uri, {
            callback: ep.done("data"),
        });
    });
}

htapNode.ListPage = function(page)
{
    l4iStorage.Set("htapm_nodels_page", parseInt(page));
    htapNode.List();
}

htapNode.Set = function(modname, modelid, nodeid)
{
    var alertid = "#htapm-node-alert";

    if (!modname && l4iStorage.Get("htapm_spec_active")) {
        modname = l4iStorage.Get("htapm_spec_active");
    }
    if (!modelid && l4iStorage.Get("htapm_nmodel_active")) {
        modelid = l4iStorage.Get("htapm_nmodel_active");
    }

    // console.log(modname +","+ modelid +","+ nodeid);

    if (!modname || !modelid) {
        return;
    }

    htapEditor.Clean();
    htapNode.nodeOpToolsRefreshCurrent = null;

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

            htapNode.setCurrent = data;
            data._status_def = htapNode.status_def;

            // console.log(data);

            l4iTemplate.Render({
                dstid: "htapm-nodeset-laymain",
                tplid: "htapm-nodeset-tpl",
                data:  data,
                success: function() {

                    var main_len = 0,
                        side_len = 0;
                    for (var i in data.model.fields) {

                        var field = data.model.fields[i];

                        switch (field.type) {

                        case "string":
                            main_len += 1;
                            break;

                        case "text":
                            main_len += 5;
                            break;

                        default:
                            side_len += 1;
                            break;
                        }
                    }
                    side_len += data.model.terms.length;

                    if (data.model.extensions.comment_perentry) {
                        side_len += 1;
                    }

                    if (data.model.extensions.permalink && data.model.extensions.permalink != "") {
                        main_len += 1;
                    }

                    var field_layout_target = "htapm-nodeset-fields";
                    if (side_len > 0 && main_len > side_len) {
                        field_layout_target = "htapm-nodeset-layside";
                    } else {
                        $("#htapm-nodeset-layside").addClass("htapm-hide");
                    }

                    //
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

                            tplid = "htapm-nodeset-tplstring";
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
                                htapEditor.Open(field.name, field.attr_format);
                            };

                            tplid = "htapm-nodeset-tpltext";
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

                            tplid = "htapm-nodeset-tplint";
                            break;

                        default:
                            continue;
                        }

                        l4iTemplate.Render({
                            dstid  : "htapm-nodeset-fields",
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

                            tplid = "htapm-nodeset-tplterm_tag";

                            l4iTemplate.Render({
                                dstid  : field_layout_target,
                                tplid  : tplid,
                                prepend: true,
                                data   : term,
                            });

                            break;

                        case "taxonomy":

                            if (!term.value) {
                                term.value = "0";
                            }

                            htapMgr.ApiCmd("term/list?modname="+ modname +"&modelid="+ term.meta.name, {
                                callback: function(err, data) {

                                    if (data.kind != "TermList") {
                                        return;
                                    }

                                    data.item = term;

                                    for (var i in data.items) {

                                        if (!data.items[i].pid) {
                                            data.items[i].pid = 0;
                                        }

                                        if (data.items[i].pid == 0) {
                                            data.items[i]._subs = htapTerm.ListSubRange(data.items, null, data.items[i].id, 0);
                                        }
                                    }

                                    tplid = "htapm-nodeset-tplterm_taxonomy";

                                    l4iTemplate.Render({
                                        dstid  : field_layout_target,
                                        tplid  : tplid,
                                        prepend: true,
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
                            dstid  : field_layout_target,
                            tplid  : "htapm-nodeset-tplext_comment_perentry",
                            append : true,
                            data   : {
                                _general_onoff: htapNode.general_onoff,
                                ext_comment_perentry: data.ext_comment_perentry,
                            },
                        });
                    }

                    if (data.model.extensions.permalink && data.model.extensions.permalink != "") {
                        l4iTemplate.Render({
                            dstid  : "htapm-nodeset-tops",
                            tplid  : "htapm-nodeset-tplext_permalink",
                            append : true,
                            data   : {
                                ext_permalink_name: data.ext_permalink_name,
                            },
                        });
                    }

                    l4iTemplate.Render({
                        dstid  : field_layout_target,
                        tplid  : "htapm-nodeset-tplstatus",
                        append : true,
                        data   : {
                            _status_def: htapNode.status_def,
                            status:      data.status,
                        },
                    });

                    htapNode.OpToolsRefresh("#htapm-node-set-opts");

                    if (data.create_new) {
                        $("#htapm-node-set-opts-label").text("Create new Content");
                    } else {
                        $("#htapm-node-set-opts-label").text("Editing");
                    }
                },
            });
        });

        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        htapMgr.TplCmd("node/set", {
            callback: function(err, tpl) {

                if (err) {
                    return ep.emit('error', err);
                }
                ep.emit("tpl", tpl);
            }
        });

        if (nodeid) {
            htapMgr.ApiCmd("node/entry?"+ uri +"&id="+ nodeid, {
                callback: ep.done("data"),
            });
        } else {
            htapMgr.ApiCmd("node-model/entry?"+ uri, {
                callback: function(err, data) {
                    ep.emit("data", {
                        kind  : "Node",
                        model : data,
                        id    : "",
                        title : "",
                        ext_comment_perentry: true,
                        create_new : true,
                    });
                },
            });
        }
    });
}


htapNode.SetCommit = function()
{
    var form = $("#htapm-nodeset-laymain"),
        alertid = "#htapm-node-alert";

    if (!htapNode.setCurrent) {
        return;
    }

    htapNode.setCurrent.title = form.find("input[name=title]").val();

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
    for (var i in htapNode.setCurrent.model.fields) {

        var field = htapNode.setCurrent.model.fields[i];

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
            field_set.value = htapEditor.Content(field.name);

            // console.log(format);

            // if (field.attrs) {

            //     for (var j in field.attrs) {

            //         if (field.attrs[j].key == "format" && field.attrs[j].value == "md") {

            //             field_set.value = htapEditor.Content(field.name);

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

    for (var i in htapNode.setCurrent.model.terms) {

        var term = htapNode.setCurrent.model.terms[i];

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

    // console.log(htapNode.setCurrent.model.terms);
    // console.log(JSON.stringify(req));

    //
    var uri = "modname="+ l4iStorage.Get("htapm_spec_active");
    uri += "&modelid="+ l4iStorage.Get("htapm_nmodel_active");

    htapMgr.ApiCmd("node/set?"+ uri, {
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
                htapNode.List();
                htapEditor.Clean();
            }, 500);
        }
    });
}


htapNode.Del = function(modname, modelid, id)
{
    l4iModal.Open({
        title  : "Delete",
        tplsrc : '<div id="htapm-node-del" class="alert alert-danger">Are you sure to delete this?</div>',
        height : "200px",
        buttons: [{
            title: "Confirm to delete",
            onclick : 'htapNode.DelCommit("'+modname+'","'+modelid+'","'+id+'")',
            style: "btn-danger",
        },{
            title: "Cancel",
            onclick : "l4iModal.Close()",
        }],
    });
}

htapNode.DelCommit = function(modname, modelid, id)
{
    var alertid = "#htapm-node-del";
    var uri = "modname="+ modname + "&modelid="+ modelid +"&id="+ id;

    htapMgr.ApiCmd("node/del?"+ uri, {
        callback : function(err, data) {

            if (!data || data.kind != "Node") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful deleted");
            setTimeout(function() {
                htapNode.List();
                l4iModal.Close();
            }, 500);
        }
    });
}
