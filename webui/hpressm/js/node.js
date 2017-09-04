var hpressNode = {
    navPrefix: "node/index/",
    speclsCurrent: [],
    specCurrent: null,
    setCurrent: null,
    cmEditor: null,
    cmEditors: {},
    general_onoff: [{
        type: true,
        name: "ON",
    }, {
        type: false,
        name: "OFF",
    }],
    status_def: [{
        type: 1,
        name: "Publish",
    }, {
        type: 2,
        name: "Draft",
    }, {
        type: 3,
        name: "Private",
    }],
    nodeOpToolsRefreshCurrent: null,
}

hpressNode.Init = function(cb) {
    hpressNode.navRefresh(cb);
}

hpressNode.navRefresh = function(cb) {
    cb = cb || function() {};

    if (hpressNode.speclsCurrent.length > 0) {

        // if (!l4iStorage.Get("hpressm_spec_active")) {
        //     for (var i in hpressNode.speclsCurrent) {
        //         l4iStorage.Set("hpressm_spec_active", hpressNode.speclsCurrent[i].meta.name);
        //         break;
        //     }
        // }

        // if (!l4iStorage.Get("hpressm_spec_active")) {
        //     return cb();
        // }

        // console.log(hpressNode.speclsCurrent);

        l4iTemplate.Render({
            dstid: "hpressm-topbar-nav-node-specls",
            tplid: "hpressm-topbar-nav-node-specls-tpl",
            data: {
                active: l4iStorage.Get("hpressm_spec_active"),
                items: hpressNode.speclsCurrent,
            },
        });

        return cb();
    }

    hpressMgr.ApiCmd("mod-set/spec-list", {
        callback: function(err, data) {

            if (err || data.error || data.kind != "SpecList") {
                return cb();
            }

            //
            for (var i in data.items) {
                hpressNode.speclsCurrent.push(data.items[i]);
                l4i.UrlEventRegister(
                    hpressNode.navPrefix + data.items[i].meta.name,
                    hpressNode.Index,
                    "hpressm-topbar"
                );
            }

            //
            if (!l4iStorage.Get("hpressm_spec_active")) {
                for (var i in hpressNode.speclsCurrent) {
                    l4iStorage.Set("hpressm_spec_active", hpressNode.speclsCurrent[i].meta.name);
                    break;
                }
            }
            if (!l4iStorage.Get("hpressm_spec_active")) {
                return cb();
            }

            l4iTemplate.Render({
                dstid: "hpressm-topbar-nav-node-specls",
                tplid: "hpressm-topbar-nav-node-specls-tpl",
                data: {
                    // active : l4iStorage.Get("hpressm_spec_active"),
                    items: hpressNode.speclsCurrent,
                },
            });

            cb();
        },
    });
}

hpressNode.OpToolsRefresh = function(div_target) {
    if (typeof div_target == "string" && div_target == hpressNode.nodeOpToolsRefreshCurrent) {
        return;
    }

    if (div_target == "clean") {
        hpressNode.nodeOpToolsRefreshCurrent = null;
        $("#hpressm-node-optools").empty();
        return;
    }

    $("#hpressm-node-optools").empty();

    if (typeof div_target == "string") {

        var opt = $("#work-content").find(div_target);
        if (opt) {
            $("#hpressm-node-optools").html(opt.html());
            hpressNode.nodeOpToolsRefreshCurrent = div_target;
        }
    }
}

hpressNode.Index = function(nav_href) {
    if (!nav_href || nav_href.length <= hpressNode.navPrefix.length) {
        return;
    }

    if (hpressNode.speclsCurrent.length < 1) {
        return;
    }

    l4iStorage.Del("hpressm_nodels_page");
    l4iStorage.Del("hpressm_termls_page");

    hpressNode.nodeOpToolsRefreshCurrent = null;
    l4iStorage.Set("hpressm_nav_last_active", nav_href);
    l4iStorage.Set("hpressm_spec_active", nav_href.substr(hpressNode.navPrefix.length));

    var alertid = "#hpressm-node-alert";

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create("tpl", function(tpl) {

            if (tpl) {
                $("#com-content").html(tpl);
            }

            var current = null;

            for (var i in hpressNode.speclsCurrent) {

                if (hpressNode.speclsCurrent[i].meta.name == l4iStorage.Get("hpressm_spec_active")) {
                    current = hpressNode.speclsCurrent[i];
                    break;
                }
            }

            if (!current) {
                return;
            }

            hpressNode.specCurrent = current;

            if (!hpressNode.specCurrent.nodeModels) {
                hpressNode.specCurrent.nodeModels = [];
            }
            if (!hpressNode.specCurrent.termModels) {
                hpressNode.specCurrent.termModels = [];
            }

            var node_model_active = null;

            for (var i in hpressNode.specCurrent.nodeModels) {

                if (!node_model_active) {
                    node_model_active = hpressNode.specCurrent.nodeModels[i].meta.name;
                }

                if (l4iStorage.Get("hpressm_nmodel_active") == hpressNode.specCurrent.nodeModels[i].meta.name) {
                    node_model_active = hpressNode.specCurrent.nodeModels[i].meta.name;
                    break;
                }
            }

            // console.log(l4iStorage.Get("hpressm_nmodel_active"));

            if (!node_model_active) {
                return; // TODO
            }

            //
            if (node_model_active != l4iStorage.Get("hpressm_nmodel_active")) {
                l4iStorage.Set("hpressm_nmodel_active", node_model_active);
            }

            //
            for (var i in hpressNode.specCurrent.nodeModels) {
                if (node_model_active == hpressNode.specCurrent.nodeModels[i].meta.name) {
                    hpressNode.List(l4iStorage.Get("hpressm_spec_active"), node_model_active);
                }
            }

            if (hpressNode.specCurrent.nodeModels.length > 0) {
                l4iTemplate.Render({
                    dstid: "hpressm-node-nmodels",
                    tplid: "hpressm-node-nmodels-tpl",
                    data: {
                        active: node_model_active,
                        items: hpressNode.specCurrent.nodeModels,
                    },
                });
            } else {
                $("#hpressm-node-nmodels").addClass("hpressm-hide");
            }

            if (hpressNode.specCurrent.termModels.length > 0) {
                l4iTemplate.Render({
                    dstid: "hpressm-node-tmodels",
                    tplid: "hpressm-node-tmodels-tpl",
                    data: {
                        items: hpressNode.specCurrent.termModels,
                    },
                });
            } else {
                $("#hpressm-node-tmodels").addClass("hpressm-hide");
            }
        });

        ep.fail(function(err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        // template
        var el = document.getElementById("hpressm-node-nmodels");
        if (!el || !el.length || el.length < 1) {
            hpressMgr.TplCmd("node/index", {
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

hpressNode.List = function(modname, modelid) {
    var alertid = "#hpressm-node-alert",
        page = 0;

    if (!modname && l4iStorage.Get("hpressm_spec_active")) {
        modname = l4iStorage.Get("hpressm_spec_active");
    }

    if (!modelid && l4iStorage.Get("hpressm_nmodel_active")) {
        modelid = l4iStorage.Get("hpressm_nmodel_active");
    }

    if (l4iStorage.Get("hpressm_nodels_page")) {
        page = l4iStorage.Get("hpressm_nodels_page");
    }

    if (!modname || !modelid) {
        return;
    }

    var uri = "modname=" + modname + "&modelid=" + modelid + "&page=" + page;
    uri += "&fields=no_fields&terms=no_terms";
    if (document.getElementById("qry_text")) {
        uri += "&qry_text=" + $("#qry_text").val();
    }

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create("tpl", "data", function(tpl, rsj) {

            if (tpl) {
                $("#work-content").html(tpl);
            }

            if (!rsj || rsj.kind != "NodeList" ||
                !rsj.items || rsj.items.length < 1) {

                $("#hpressm-nodels").empty();
                $("#hpressm-termls").empty();
                l4i.InnerAlert(alertid, 'alert-info', "Item Not Found");
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
                dstid: "hpressm-nodels",
                tplid: "hpressm-nodels-tpl",
                data: {
                    model: rsj.model,
                    modname: modname,
                    modelid: modelid,
                    items: rsj.items,
                    _status_def: hpressNode.status_def,
                },
                success: function() {

                    rsj.meta.RangeLen = 20;

                    l4iTemplate.Render({
                        dstid: "hpressm-nodels-pager",
                        tplid: "hpressm-nodels-pager-tpl",
                        data: l4i.Pager(rsj.meta),
                    });

                    hpressNode.OpToolsRefresh("#hpressm-node-list-opts");
                }
            });
        });

        ep.fail(function(err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        // template
        var el = document.getElementById("hpressm-nodels");
        if (!el || el.length < 1) {
            hpressMgr.TplCmd("node/list", {
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

        hpressMgr.ApiCmd("node/list?" + uri, {
            callback: ep.done("data"),
        });
    });
}

hpressNode.ListPage = function(page) {
    l4iStorage.Set("hpressm_nodels_page", parseInt(page));
    hpressNode.List();
}

hpressNode.ListBatchSelectAll = function() {
    var form = $("#hpressm-nodels");
    if (!form) {
        return;
    }

    var checked = false;
    if (form.find(".hpressm-nodels-chk-all").is(':checked')) {
        checked = true;
    }

    form.find(".hpressm-nodels-chk-item").each(function() {
        if (checked) {
            $(this).prop("checked", true);
        } else {
            $(this).prop("checked", false);
        }
    });

    hpressNode.ListBatchSelectTodoBtnRefresh(checked);
}

hpressNode.ListBatchSelectTodoBtnRefresh = function(onoff) {
    if (onoff !== true && onoff !== false) {

        onoff = false;

        $("#hpressm-nodels").find(".hpressm-nodels-chk-item").each(function() {

            if ($(this).is(":checked")) {
                onoff = true;
                return (false);
            }
        });
    }

    if (onoff === true) {
        $("#hpressm-nodels-batch-select-todo-btn").css({
            "display": "block"
        });
    } else {
        $("#hpressm-nodels-batch-select-todo-btn").css({
            "display": "none"
        });
    }
}

hpressNode.ListBatchSelectTodo = function() {
    var form = $("#hpressm-nodels");
    if (!form) {
        return;
    }

    var select_num = 0;

    form.find(".hpressm-nodels-chk-item").each(function() {

        if ($(this).is(":checked")) {
            select_num++;
        }
    });

    var params = {
        select_num: select_num,
    };

    hpressMgr.TplCmd("node/list-batch-set", {
        callback: function(err, data) {

            if (err) {
                return;
            }

            l4iModal.Open({
                title: "Batch operation",
                tplsrc: data,
                data: params,
                width: 800,
                height: 300,
                buttons: [{
                    title: "Confirm to delete",
                    onclick: "hpressNode.ListBatchSelectTodoDelete()",
                    style: "btn-danger",
                }, {
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }],
            });
        },
    })
}

hpressNode.ListBatchSelectTodoDelete = function(modname, modelid) {
    if (!modname && l4iStorage.Get("hpressm_spec_active")) {
        modname = l4iStorage.Get("hpressm_spec_active");
    }
    if (!modelid && l4iStorage.Get("hpressm_nmodel_active")) {
        modelid = l4iStorage.Get("hpressm_nmodel_active");
    }

    if (!modname || !modelid) {
        return;
    }

    var ids = [];

    $("#hpressm-nodels").find(".hpressm-nodels-chk-item").each(function() {

        if ($(this).is(":checked")) {
            ids.push($(this).val());
        }
    });

    var alertid = "#hpressm-nodels-batch-set-alert";

    hpressNode.DelBatch(modname, modelid, ids, function(err, data) {

        if (err) {
            l4i.InnerAlert(alertid, 'alert-danger', err);
        } else if (data && data.error) {
            l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
        } else if (data && data.kind == "Node") {
            l4i.InnerAlert(alertid, 'alert-success', "Successful operation");
            hpressNode.List();
            setTimeout(function() {
                l4iModal.Close();
            }, 500);
        } else {
            l4i.InnerAlert(alertid, 'alert-danger', "unknown error");
        }
    });
}

hpressNode.Set = function(modname, modelid, nodeid) {
    var alertid = "#hpressm-node-alert";

    if (!modname && l4iStorage.Get("hpressm_spec_active")) {
        modname = l4iStorage.Get("hpressm_spec_active");
    }
    if (!modelid && l4iStorage.Get("hpressm_nmodel_active")) {
        modelid = l4iStorage.Get("hpressm_nmodel_active");
    }

    // console.log(modname +","+ modelid +","+ nodeid);

    if (!modname || !modelid) {
        return;
    }

    hpressEditor.Clean();
    hpressNode.nodeOpToolsRefreshCurrent = null;

    var uri = "modname=" + modname + "&modelid=" + modelid;

    // console.log(uri);
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create("tpl", "data", function(tpl, data) {

            if (!tpl) {
                return; // TODO
            }

            $("#work-content").html(tpl);

            if (!data || data.kind != "Node") {
                return l4i.InnerAlert(alertid, 'alert-info', "Item Not Found");
            }

            if (!data.status) {
                data.status = 1;
            }

            if (!data.model.terms) {
                data.model.terms = [];
            }

            if (!data.ext_comment_enable) {
                data.ext_comment_enable = false;
            }

            if (!data.ext_comment_perentry) {
                data.ext_comment_perentry = false;
            }

            if (!data.ext_permalink_name) {
                data.ext_permalink_name = "";
            }

            $(alertid).hide();

            hpressNode.setCurrent = data;
            data._status_def = hpressNode.status_def;

            // console.log(data);

            l4iTemplate.Render({
                dstid: "hpressm-nodeset-laymain",
                tplid: "hpressm-nodeset-tpl",
                data: data,
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

                    var field_layout_target = "hpressm-nodeset-fields";
                    if (side_len > 0 && main_len > side_len) {
                        field_layout_target = "hpressm-nodeset-layside";
                    } else {
                        $("#hpressm-nodeset-layside").addClass("hpressm-hide");
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

                                tplid = "hpressm-nodeset-tplstring";
                                break;

                            case "text":

                                if (field.attrs) {
                                    for (var j in field.attrs) {
                                        field["attr_" + field.attrs[j].key] = field.attrs[j].value;
                                    }
                                }

                                if (field_entry.attrs) {
                                    for (var j in field_entry.attrs) {
                                        field["attr_" + field_entry.attrs[j].key] = field_entry.attrs[j].value;
                                    }
                                }

                                if (!field.value) {
                                    field.value = "";
                                }

                                if (!field.attr_format) {
                                    field.attr_format = "text";
                                }

                                cb = function() {
                                    hpressEditor.Open(field.name, field.attr_format);
                                };

                                tplid = "hpressm-nodeset-tpltext";
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

                                tplid = "hpressm-nodeset-tplint";
                                break;

                            default:
                                continue;
                        }

                        l4iTemplate.Render({
                            dstid: "hpressm-nodeset-fields",
                            tplid: tplid,
                            append: true,
                            data: field,
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

                                tplid = "hpressm-nodeset-tplterm_tag";

                                l4iTemplate.Render({
                                    dstid: field_layout_target,
                                    tplid: tplid,
                                    prepend: true,
                                    data: term,
                                });

                                break;

                            case "taxonomy":

                                if (!term.value) {
                                    term.value = "0";
                                }

                                hpressMgr.ApiCmd("term/list?modname=" + modname + "&modelid=" + term.meta.name, {
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
                                                data.items[i]._subs = hpressTerm.ListSubRange(data.items, null, data.items[i].id, 0);
                                            }
                                        }

                                        tplid = "hpressm-nodeset-tplterm_taxonomy";

                                        l4iTemplate.Render({
                                            dstid: field_layout_target,
                                            tplid: tplid,
                                            prepend: true,
                                            data: data,
                                        });
                                    },
                                });

                                break;

                            default:
                                continue;
                        }
                    }


                    if (data.model.extensions.comment_enable && data.model.extensions.comment_perentry) {
                        l4iTemplate.Render({
                            dstid: field_layout_target,
                            tplid: "hpressm-nodeset-tplext_comment_perentry",
                            append: true,
                            data: {
                                _general_onoff: hpressNode.general_onoff,
                                ext_comment_perentry: data.ext_comment_perentry,
                            },
                        });
                    }

                    if (data.model.extensions.permalink && data.model.extensions.permalink != "") {
                        l4iTemplate.Render({
                            dstid: "hpressm-nodeset-tops",
                            tplid: "hpressm-nodeset-tplext_permalink",
                            append: true,
                            data: {
                                ext_permalink_name: data.ext_permalink_name,
                            },
                        });
                    }

                    l4iTemplate.Render({
                        dstid: field_layout_target,
                        tplid: "hpressm-nodeset-tplstatus",
                        append: true,
                        data: {
                            _status_def: hpressNode.status_def,
                            status: data.status,
                        },
                    });

                    hpressNode.OpToolsRefresh("#hpressm-node-set-opts");

                    if (data.create_new) {
                        $("#hpressm-node-set-opts-label").text("Create new Content");
                    } else {
                        $("#hpressm-node-set-opts-label").text("Editing");
                    }
                },
            });
        });

        ep.fail(function(err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        hpressMgr.TplCmd("node/set", {
            callback: function(err, tpl) {

                if (err) {
                    return ep.emit('error', err);
                }
                ep.emit("tpl", tpl);
            }
        });

        if (nodeid) {
            hpressMgr.ApiCmd("node/entry?" + uri + "&id=" + nodeid, {
                callback: ep.done("data"),
            });
        } else {
            hpressMgr.ApiCmd("node-model/entry?" + uri, {
                callback: function(err, data) {
                    ep.emit("data", {
                        kind: "Node",
                        model: data,
                        id: "",
                        title: "",
                        ext_comment_perentry: true,
                        create_new: true,
                    });
                },
            });
        }
    });
}


hpressNode.SetCommit = function() {
    var form = $("#hpressm-nodeset-layout"),
        alertid = "#hpressm-node-alert";

    if (!hpressNode.setCurrent) {
        return;
    }

    hpressNode.setCurrent.title = form.find("input[name=title]").val();

    var req = {
        id: form.find("input[name=id]").val(),
        title: form.find("input[name=title]").val(),
        status: parseInt(form.find("select[name=status]").val()),
        fields: [],
        terms: [],
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
    for (var i in hpressNode.setCurrent.model.fields) {

        var field = hpressNode.setCurrent.model.fields[i];

        var field_set = {
            name: field.name,
            value: null,
            attrs: [],
        };

        switch (field.type) {

            case "text":

                var format = form.find("input[name=field_" + field.name + "_attr_format]").val();
                if (!format) {
                    format = "text";
                }

                field_set.attrs.push({
                    key: "format",
                    value: format
                });
                field_set.value = hpressEditor.Content(field.name);

                // console.log(format);

                // if (field.attrs) {

                //     for (var j in field.attrs) {

                //         if (field.attrs[j].key == "format" && field.attrs[j].value == "md") {

                //             field_set.value = hpressEditor.Content(field.name);

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
                field_set.value = form.find("input[name=field_" + field.name + "]").val();
                break;

            case "int8":
            case "int16":
            case "int32":
            case "int64":
            case "uint8":
            case "uint16":
            case "uint32":
            case "uint64":
                field_set.value = form.find("input[name=field_" + field.name + "]").val();
                break;

        }

        if (field_set.value) {
            req.fields.push(field_set);
        }
    }

    for (var i in hpressNode.setCurrent.model.terms) {

        var term = hpressNode.setCurrent.model.terms[i];

        var val = null;

        switch (term.type) {

            case "tag":
                val = form.find("input[name=term_" + term.meta.name + "]").val();
                break;
            case "taxonomy":
                val = form.find("select[name=term_" + term.meta.name + "]").val();
                break;
        }

        if (val) {
            req.terms.push({
                name: term.meta.name,
                value: val
            });
        }
    }

    // console.log(hpressNode.setCurrent.model.terms);
    // console.log(JSON.stringify(req));

    //
    var uri = "modname=" + l4iStorage.Get("hpressm_spec_active");
    uri += "&modelid=" + l4iStorage.Get("hpressm_nmodel_active");

    hpressMgr.ApiCmd("node/set?" + uri, {
        method: "POST",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (!data || data.kind != "Node") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            // console.log(data.id);
            form.find("input[name=id]").val(data.id);

            l4i.InnerAlert(alertid, 'alert-success', "Successful operation");
            setTimeout(function() {
                hpressNode.List();
                hpressEditor.Clean();
            }, 500);
        }
    });
}


hpressNode.Del = function(modname, modelid, id) {
    l4iModal.Open({
        title: "Delete",
        tplsrc: '<div id="hpressm-node-del" class="alert alert-danger">Are you sure to delete this?</div>',
        height: "200px",
        buttons: [{
            title: "Confirm to delete",
            onclick: 'hpressNode.DelCommit("' + modname + '","' + modelid + '","' + id + '")',
            style: "btn-danger",
        }, {
            title: "Cancel",
            onclick: "l4iModal.Close()",
        }],
    });
}

hpressNode.DelCommit = function(modname, modelid, id) {
    var alertid = "#hpressm-node-del";
    var uri = "modname=" + modname + "&modelid=" + modelid + "&id=" + id;

    hpressMgr.ApiCmd("node/del?" + uri, {
        callback: function(err, data) {

            if (!data || data.kind != "Node") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful deleted");
            setTimeout(function() {
                hpressNode.List();
                l4iModal.Close();
            }, 500);
        }
    });
}

hpressNode.DelBatch = function(modname, modelid, ids, cb) {
    var uri = "modname=" + modname + "&modelid=" + modelid + "&id=" + ids.join(",");

    hpressMgr.ApiCmd("node/del?" + uri, {
        callback: cb,
    });
}