var htpNode = {
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

htpNode.Init = function(cb)
{
    htpNode.navRefresh(cb);
}

htpNode.navRefresh = function(cb)
{
    cb = cb || function(){};

    if (htpNode.speclsCurrent.length > 0) {

        // if (!l4iStorage.Get("htpm_spec_active")) {
        //     for (var i in htpNode.speclsCurrent) {
        //         l4iStorage.Set("htpm_spec_active", htpNode.speclsCurrent[i].meta.name);
        //         break;
        //     }
        // }

        // if (!l4iStorage.Get("htpm_spec_active")) {
        //     return cb();
        // }

        // console.log(htpNode.speclsCurrent);

        l4iTemplate.Render({
            dstid: "htpm-topbar-nav-node-specls",
            tplid: "htpm-topbar-nav-node-specls-tpl",
            data:  {
                active : l4iStorage.Get("htpm_spec_active"),
                items  : htpNode.speclsCurrent,
            },
        });

        return cb();
    }

    htpMgr.ApiCmd("mod-set/spec-list", {
        callback: function(err, data) {

            if (err || data.error || data.kind != "SpecList") {
                return cb();
            }

            //
            for (var i in data.items) {
                htpNode.speclsCurrent.push(data.items[i]);
                l4i.UrlEventRegister(
                    htpNode.navPrefix + data.items[i].meta.name,
                    htpNode.Index,
                    "htpm-topbar"
                );
            }

            //
            if (!l4iStorage.Get("htpm_spec_active")) {
                for (var i in htpNode.speclsCurrent) {
                    l4iStorage.Set("htpm_spec_active", htpNode.speclsCurrent[i].meta.name);
                    break;
                }
            }
            if (!l4iStorage.Get("htpm_spec_active")) {
                return cb();
            }

            l4iTemplate.Render({
                dstid: "htpm-topbar-nav-node-specls",
                tplid: "htpm-topbar-nav-node-specls-tpl",
                data:  {
                    // active : l4iStorage.Get("htpm_spec_active"),
                    items  : htpNode.speclsCurrent,
                },
            });

            cb();
        },
    });
}

htpNode.OpToolsRefresh = function(div_target)
{
    if (typeof div_target == "string" && div_target == htpNode.nodeOpToolsRefreshCurrent) {
        return;
    }

    if (div_target == "clean") {
        htpNode.nodeOpToolsRefreshCurrent = null;
        $("#htpm-node-optools").empty();
        return;
    }

    $("#htpm-node-optools").empty();

    if (typeof div_target == "string") {

        var opt = $("#work-content").find(div_target);
        if (opt) {
            $("#htpm-node-optools").html(opt.html());
            htpNode.nodeOpToolsRefreshCurrent = div_target;
        }
    }
}

htpNode.Index = function(nav_href)
{
    if (!nav_href || nav_href.length <= htpNode.navPrefix.length) {
        return;
    }

    if (htpNode.speclsCurrent.length < 1) {
        return;
    }

    l4iStorage.Del("htpm_nodels_page");
    l4iStorage.Del("htpm_termls_page");

    htpNode.nodeOpToolsRefreshCurrent = null;
    l4iStorage.Set("htpm_nav_last_active", nav_href);
    l4iStorage.Set("htpm_spec_active", nav_href.substr(htpNode.navPrefix.length));

    var alertid = "#htpm-node-alert";

    seajs.use(["ep"], function (EventProxy) {

        var ep = EventProxy.create("tpl", function (tpl) {

            if (tpl) {
                $("#com-content").html(tpl);
            }

            var current = null;

            for (var i in htpNode.speclsCurrent) {

                if (htpNode.speclsCurrent[i].meta.name == l4iStorage.Get("htpm_spec_active")) {
                    current = htpNode.speclsCurrent[i];
                    break;
                }
            }

            if (!current) {
                return;
            }

            htpNode.specCurrent = current;

            if (!htpNode.specCurrent.nodeModels) {
                htpNode.specCurrent.nodeModels = [];
            }
            if (!htpNode.specCurrent.termModels) {
                htpNode.specCurrent.termModels = [];
            }

            var node_model_active = null;

            for (var i in htpNode.specCurrent.nodeModels) {

                if (!node_model_active) {
                    node_model_active = htpNode.specCurrent.nodeModels[i].meta.name;
                }

                if (l4iStorage.Get("htpm_nmodel_active") == htpNode.specCurrent.nodeModels[i].meta.name) {
                    node_model_active = htpNode.specCurrent.nodeModels[i].meta.name;
                    break;
                }
            }

            // console.log(l4iStorage.Get("htpm_nmodel_active"));

            if (!node_model_active) {
                return; // TODO
            }

                    //
            if (node_model_active != l4iStorage.Get("htpm_nmodel_active")) {
                l4iStorage.Set("htpm_nmodel_active", node_model_active);
            }

            //
            for (var i in htpNode.specCurrent.nodeModels) {
                if (node_model_active == htpNode.specCurrent.nodeModels[i].meta.name) {
                    htpNode.List(l4iStorage.Get("htpm_spec_active"), node_model_active);
                }
            }

            if (htpNode.specCurrent.nodeModels.length > 0) {
                l4iTemplate.Render({
                    dstid: "htpm-node-nmodels",
                    tplid: "htpm-node-nmodels-tpl",
                    data:  {
                        active: node_model_active,
                        items: htpNode.specCurrent.nodeModels,
                    },
                });
            } else {
                $("#htpm-node-nmodels").addClass("htpm-hide");
            }

            if (htpNode.specCurrent.termModels.length > 0) {
                l4iTemplate.Render({
                    dstid: "htpm-node-tmodels",
                    tplid: "htpm-node-tmodels-tpl",
                    data:  {
                        items: htpNode.specCurrent.termModels,
                    },
                });
            } else {
                $("#htpm-node-tmodels").addClass("htpm-hide");
            }
        });

        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        // template
        var el = document.getElementById("htpm-node-nmodels");
        if (!el || !el.length || el.length < 1) {
            htpMgr.TplCmd("node/index", {
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

htpNode.List = function(modname, modelid)
{
    var alertid = "#htpm-node-alert",
        page = 0;

    if (!modname && l4iStorage.Get("htpm_spec_active")) {
        modname = l4iStorage.Get("htpm_spec_active");
    }

    if (!modelid && l4iStorage.Get("htpm_nmodel_active")) {
        modelid = l4iStorage.Get("htpm_nmodel_active");
    }

    if (l4iStorage.Get("htpm_nodels_page")) {
        page = l4iStorage.Get("htpm_nodels_page");
    }

    if (!modname || !modelid) {
        return;
    }

    var uri = "modname="+ modname +"&modelid="+ modelid +"&page="+ page;
    uri += "&fields=no_fields&terms=no_terms";
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

                $("#htpm-nodels").empty();
                $("#htpm-termls").empty();
                l4i.InnerAlert(alertid, 'alert-danger', "Item Not Found");
            } else {
                $(alertid).hide();
            }

            if (!rsj.items) {
                rsj.items = [];
            }

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
                dstid: "htpm-nodels",
                tplid: "htpm-nodels-tpl",
                data:  {
                    model   : rsj.model,
                    modname : modname,
                    modelid : modelid,
                    items   : rsj.items,
                    _status_def : htpNode.status_def,
                },
                success: function() {

                    rsj.meta.RangeLen = 20;

                    l4iTemplate.Render({
                        dstid : "htpm-nodels-pager",
                        tplid : "htpm-nodels-pager-tpl",
                        data  : l4i.Pager(rsj.meta),
                    });

                    htpNode.OpToolsRefresh("#htpm-node-list-opts");
                }
            });
        });

        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        // template
        var el = document.getElementById("htpm-nodels");
        if (!el || el.length < 1) {
            htpMgr.TplCmd("node/list", {
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

        htpMgr.ApiCmd("node/list?"+ uri, {
            callback: ep.done("data"),
        });
    });
}

htpNode.ListPage = function(page)
{
    l4iStorage.Set("htpm_nodels_page", parseInt(page));
    htpNode.List();
}

htpNode.ListBatchSelectAll = function()
{
    var form = $("#htpm-nodels");
    if (!form) {
        return;
    }

    var checked = false;
    if (form.find(".htpm-nodels-chk-all").is(':checked')) {
        checked = true;
    }

    form.find(".htpm-nodels-chk-item").each(function() {
        if (checked) {
            $(this).prop("checked", true);
        } else {
            $(this).prop("checked", false);
        }
    });

    htpNode.ListBatchSelectTodoBtnRefresh(checked);
}

htpNode.ListBatchSelectTodoBtnRefresh = function(onoff)
{
    if (onoff !== true && onoff !== false) {

        onoff = false;

        $("#htpm-nodels").find(".htpm-nodels-chk-item").each(function() {

            if ($(this).is(":checked")) {
                onoff = true;
                return(false);
            }
        });
    }

    if (onoff === true) {
        $("#htpm-nodels-batch-select-todo-btn").css({"display": "block"});
    } else {
        $("#htpm-nodels-batch-select-todo-btn").css({"display": "none"});
    }
}

htpNode.ListBatchSelectTodo = function()
{
    var form = $("#htpm-nodels");
    if (!form) {
        return;
    }

    var select_num = 0;

    form.find(".htpm-nodels-chk-item").each(function() {

        if ($(this).is(":checked")) {
            select_num++;
        }
    });

    var params = {
        select_num: select_num,
    };

    htpMgr.TplCmd("node/list-batch-set", {
        callback: function(err, data) {

            if (err) {
                return;
            }

            l4iModal.Open({
                title  : "Batch operation",
                tplsrc : data,
                data   : params,
                width  : 800,
                height : 300,
                buttons: [{
                    title: "Confirm to delete",
                    onclick : "htpNode.ListBatchSelectTodoDelete()",
                    style: "btn-danger",
                }, {
                    title: "Cancel",
                    onclick : "l4iModal.Close()",
                }],
            });
        },
    })
}

htpNode.ListBatchSelectTodoDelete = function(modname, modelid)
{
    if (!modname && l4iStorage.Get("htpm_spec_active")) {
        modname = l4iStorage.Get("htpm_spec_active");
    }
    if (!modelid && l4iStorage.Get("htpm_nmodel_active")) {
        modelid = l4iStorage.Get("htpm_nmodel_active");
    }

    if (!modname || !modelid) {
        return;
    }

    var ids = [];

    $("#htpm-nodels").find(".htpm-nodels-chk-item").each(function() {

        if ($(this).is(":checked")) {
            ids.push($(this).val());
        }
    });

    var alertid = "#htpm-nodels-batch-set-alert";

    htpNode.DelBatch(modname, modelid, ids, function(err, data) {

        if (err) {
            l4i.InnerAlert(alertid, 'alert-danger', err);
        } else if (data && data.error) {
            l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
        } else if (data && data.kind == "Node") {
            l4i.InnerAlert(alertid, 'alert-success', "Successful operation");
            htpNode.List();
            setTimeout(function() {
                l4iModal.Close();
            }, 500);
        } else {
            l4i.InnerAlert(alertid, 'alert-danger', "unknown error");
        }
    });
}

htpNode.Set = function(modname, modelid, nodeid)
{
    var alertid = "#htpm-node-alert";

    if (!modname && l4iStorage.Get("htpm_spec_active")) {
        modname = l4iStorage.Get("htpm_spec_active");
    }
    if (!modelid && l4iStorage.Get("htpm_nmodel_active")) {
        modelid = l4iStorage.Get("htpm_nmodel_active");
    }

    // console.log(modname +","+ modelid +","+ nodeid);

    if (!modname || !modelid) {
        return;
    }

    htpEditor.Clean();
    htpNode.nodeOpToolsRefreshCurrent = null;

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

            htpNode.setCurrent = data;
            data._status_def = htpNode.status_def;

            // console.log(data);

            l4iTemplate.Render({
                dstid: "htpm-nodeset-laymain",
                tplid: "htpm-nodeset-tpl",
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

                    var field_layout_target = "htpm-nodeset-fields";
                    if (side_len > 0 && main_len > side_len) {
                        field_layout_target = "htpm-nodeset-layside";
                    } else {
                        $("#htpm-nodeset-layside").addClass("htpm-hide");
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

                            tplid = "htpm-nodeset-tplstring";
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
                                htpEditor.Open(field.name, field.attr_format);
                            };

                            tplid = "htpm-nodeset-tpltext";
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

                            tplid = "htpm-nodeset-tplint";
                            break;

                        default:
                            continue;
                        }

                        l4iTemplate.Render({
                            dstid  : "htpm-nodeset-fields",
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

                            tplid = "htpm-nodeset-tplterm_tag";

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

                            htpMgr.ApiCmd("term/list?modname="+ modname +"&modelid="+ term.meta.name, {
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
                                            data.items[i]._subs = htpTerm.ListSubRange(data.items, null, data.items[i].id, 0);
                                        }
                                    }

                                    tplid = "htpm-nodeset-tplterm_taxonomy";

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
                            tplid  : "htpm-nodeset-tplext_comment_perentry",
                            append : true,
                            data   : {
                                _general_onoff: htpNode.general_onoff,
                                ext_comment_perentry: data.ext_comment_perentry,
                            },
                        });
                    }

                    if (data.model.extensions.permalink && data.model.extensions.permalink != "") {
                        l4iTemplate.Render({
                            dstid  : "htpm-nodeset-tops",
                            tplid  : "htpm-nodeset-tplext_permalink",
                            append : true,
                            data   : {
                                ext_permalink_name: data.ext_permalink_name,
                            },
                        });
                    }

                    l4iTemplate.Render({
                        dstid  : field_layout_target,
                        tplid  : "htpm-nodeset-tplstatus",
                        append : true,
                        data   : {
                            _status_def: htpNode.status_def,
                            status:      data.status,
                        },
                    });

                    htpNode.OpToolsRefresh("#htpm-node-set-opts");

                    if (data.create_new) {
                        $("#htpm-node-set-opts-label").text("Create new Content");
                    } else {
                        $("#htpm-node-set-opts-label").text("Editing");
                    }
                },
            });
        });

        ep.fail(function (err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        htpMgr.TplCmd("node/set", {
            callback: function(err, tpl) {

                if (err) {
                    return ep.emit('error', err);
                }
                ep.emit("tpl", tpl);
            }
        });

        if (nodeid) {
            htpMgr.ApiCmd("node/entry?"+ uri +"&id="+ nodeid, {
                callback: ep.done("data"),
            });
        } else {
            htpMgr.ApiCmd("node-model/entry?"+ uri, {
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


htpNode.SetCommit = function()
{
    var form = $("#htpm-nodeset-layout"),
        alertid = "#htpm-node-alert";

    if (!htpNode.setCurrent) {
        return;
    }

    htpNode.setCurrent.title = form.find("input[name=title]").val();

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
    for (var i in htpNode.setCurrent.model.fields) {

        var field = htpNode.setCurrent.model.fields[i];

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
            field_set.value = htpEditor.Content(field.name);

            // console.log(format);

            // if (field.attrs) {

            //     for (var j in field.attrs) {

            //         if (field.attrs[j].key == "format" && field.attrs[j].value == "md") {

            //             field_set.value = htpEditor.Content(field.name);

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

    for (var i in htpNode.setCurrent.model.terms) {

        var term = htpNode.setCurrent.model.terms[i];

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

    // console.log(htpNode.setCurrent.model.terms);
    // console.log(JSON.stringify(req));

    //
    var uri = "modname="+ l4iStorage.Get("htpm_spec_active");
    uri += "&modelid="+ l4iStorage.Get("htpm_nmodel_active");

    htpMgr.ApiCmd("node/set?"+ uri, {
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
                htpNode.List();
                htpEditor.Clean();
            }, 500);
        }
    });
}


htpNode.Del = function(modname, modelid, id)
{
    l4iModal.Open({
        title  : "Delete",
        tplsrc : '<div id="htpm-node-del" class="alert alert-danger">Are you sure to delete this?</div>',
        height : "200px",
        buttons: [{
            title: "Confirm to delete",
            onclick : 'htpNode.DelCommit("'+modname+'","'+modelid+'","'+id+'")',
            style: "btn-danger",
        },{
            title: "Cancel",
            onclick : "l4iModal.Close()",
        }],
    });
}

htpNode.DelCommit = function(modname, modelid, id)
{
    var alertid = "#htpm-node-del";
    var uri = "modname="+ modname + "&modelid="+ modelid +"&id="+ id;

    htpMgr.ApiCmd("node/del?"+ uri, {
        callback : function(err, data) {

            if (!data || data.kind != "Node") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful deleted");
            setTimeout(function() {
                htpNode.List();
                l4iModal.Close();
            }, 500);
        }
    });
}

htpNode.DelBatch = function(modname, modelid, ids, cb)
{
    var uri = "modname="+ modname + "&modelid="+ modelid +"&id="+ ids.join(",");

    htpMgr.ApiCmd("node/del?"+ uri, {
        callback : cb,
    });
}
