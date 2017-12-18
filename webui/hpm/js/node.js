// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

var hpNode = {
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
    node_refer_back: null,
    text_formats: [
        {
            name: "text",
            value: "Text",
        },
        {
            name: "html",
            value: "Html",
        },
        {
            name: "shtml",
            value: "Script Html",
        },
        {
            name: "md",
            value: "Makedown",
        },
    ],
}

hpNode.Init = function(cb) {
    hpNode.navRefresh(cb);
}

hpNode.navRefreshForce = function(cb) {
    hpNode.speclsCurrent = [];
    hpNode.navRefresh(cb);
}

hpNode.navRefresh = function(cb) {
    cb = cb || function() {};

    if (hpNode.speclsCurrent.length > 0) {

        // if (!l4iStorage.Get("hpm_spec_active")) {
        //     for (var i in hpNode.speclsCurrent) {
        //         l4iStorage.Set("hpm_spec_active", hpNode.speclsCurrent[i].meta.name);
        //         break;
        //     }
        // }

        // if (!l4iStorage.Get("hpm_spec_active")) {
        //     return cb();
        // }

        // console.log(hpNode.speclsCurrent);

        l4iTemplate.Render({
            dstid: "hpm-topbar-nav-node-specls",
            tplid: "hpm-topbar-nav-node-specls-tpl",
            data: {
                active: l4iStorage.Get("hpm_spec_active"),
                items: hpNode.speclsCurrent,
            },
        });

        return cb();
    }

    hpMgr.ApiCmd("mod-set/spec-list", {
        callback: function(err, data) {

            if (err || data.error || data.kind != "SpecList") {
                return cb();
            }

            hpNode.speclsCurrent = [];

            //
            for (var i in data.items) {
                if (!data.items[i].status || data.items[i].status != 1) {
                    continue;
                }
                hpNode.speclsCurrent.push(data.items[i]);
                l4i.UrlEventRegister(
                    hpNode.navPrefix + data.items[i].meta.name,
                    hpNode.Index,
                    "hpm-topbar"
                );
            }

            //
            if (!l4iStorage.Get("hpm_spec_active")) {
                for (var i in hpNode.speclsCurrent) {
                    l4iStorage.Set("hpm_spec_active", hpNode.speclsCurrent[i].meta.name);
                    break;
                }
            }
            if (!l4iStorage.Get("hpm_spec_active")) {
                return cb();
            }

            l4iTemplate.Render({
                dstid: "hpm-topbar-nav-node-specls",
                tplid: "hpm-topbar-nav-node-specls-tpl",
                data: {
                    // active : l4iStorage.Get("hpm_spec_active"),
                    items: hpNode.speclsCurrent,
                },
            });

            cb();
        },
    });
}

hpNode.OpToolsRefresh = function(div_target, fn) {
    if (typeof div_target == "string" && div_target == hpNode.nodeOpToolsRefreshCurrent) {
        return;
    }

    if (div_target == "clean") {
        hpNode.nodeOpToolsRefreshCurrent = null;
        $("#hpm-node-optools").empty();
        return;
    }

    $("#hpm-node-optools").empty();

    if (typeof div_target == "string") {

        var opt = $("#work-content").find(div_target);
        if (opt) {
            $("#hpm-node-optools").html(opt.html());
            hpNode.nodeOpToolsRefreshCurrent = div_target;
        }
    }

    if (fn) {
        fn();
    }
}

hpNode.Index = function(nav_href) {
    if (!nav_href || nav_href.length <= hpNode.navPrefix.length) {
        return;
    }

    if (hpNode.speclsCurrent.length < 1) {
        return;
    }

    l4iStorage.Del("hpm_nodels_page");
    l4iStorage.Del("hpm_termls_page");

    hpNode.nodeOpToolsRefreshCurrent = null;
    l4iStorage.Set("hpm_nav_last_active", nav_href);
    l4iStorage.Set("hpm_spec_active", nav_href.substr(hpNode.navPrefix.length));

    var alertid = "#hpm-node-alert";

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create("tpl", function(tpl) {

            if (tpl) {
                $("#com-content").html(tpl);
            }

            var current = null;

            for (var i in hpNode.speclsCurrent) {

                if (hpNode.speclsCurrent[i].meta.name == l4iStorage.Get("hpm_spec_active")) {
                    current = hpNode.speclsCurrent[i];
                    break;
                }
            }

            if (!current) {
                return;
            }

            hpNode.specCurrent = current;

            if (!hpNode.specCurrent.nodeModels) {
                hpNode.specCurrent.nodeModels = [];
            }
            if (!hpNode.specCurrent.termModels) {
                hpNode.specCurrent.termModels = [];
            }

            var node_model_active = null;

            for (var i in hpNode.specCurrent.nodeModels) {

                if (!node_model_active) {
                    node_model_active = hpNode.specCurrent.nodeModels[i].meta.name;
                }

                if (l4iStorage.Get("hpm_nmodel_active") == hpNode.specCurrent.nodeModels[i].meta.name) {
                    node_model_active = hpNode.specCurrent.nodeModels[i].meta.name;
                    break;
                }
            }

            // console.log(l4iStorage.Get("hpm_nmodel_active"));

            if (!node_model_active) {
                return; // TODO
            }

            //
            if (node_model_active != l4iStorage.Get("hpm_nmodel_active")) {
                l4iStorage.Set("hpm_nmodel_active", node_model_active);
            }

            //
            for (var i in hpNode.specCurrent.nodeModels) {
                if (node_model_active == hpNode.specCurrent.nodeModels[i].meta.name) {
                    hpNode.List(l4iStorage.Get("hpm_spec_active"), node_model_active);
                }
            }

            if (hpNode.specCurrent.nodeModels.length > 0) {
                l4iTemplate.Render({
                    dstid: "hpm-node-nmodels",
                    tplid: "hpm-node-nmodels-tpl",
                    data: {
                        active: node_model_active,
                        items: hpNode.specCurrent.nodeModels,
                    },
                });
            } else {
                $("#hpm-node-nmodels").addClass("hpm-hide");
            }

            if (hpNode.specCurrent.termModels.length > 0) {
                l4iTemplate.Render({
                    dstid: "hpm-node-tmodels",
                    tplid: "hpm-node-tmodels-tpl",
                    data: {
                        items: hpNode.specCurrent.termModels,
                    },
                });
            } else {
                $("#hpm-node-tmodels").addClass("hpm-hide");
            }
        });

        ep.fail(function(err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        // template
        var el = document.getElementById("hpm-node-nmodels");
        if (!el || !el.length || el.length < 1) {
            hpMgr.TplCmd("node/index", {
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

hpNode.List = function(modname, modelid, referid) {
    var alertid = "#hpm-node-alert",
        page = 0;

    if (!modname && l4iStorage.Get("hpm_spec_active")) {
        modname = l4iStorage.Get("hpm_spec_active");
    }

    if (!modelid && l4iStorage.Get("hpm_nmodel_active")) {
        modelid = l4iStorage.Get("hpm_nmodel_active");
    }

    if (!referid && l4iStorage.Get("hpm_node_refer_active")) {
        referid = l4iStorage.Get("hpm_node_refer_active");
    }
    if (!referid) {
        referid = "";
    } else {
        l4iStorage.Set("hpm_node_refer_active", referid);
    }

    if (l4iStorage.Get("hpm_nodels_page")) {
        page = l4iStorage.Get("hpm_nodels_page");
    }

    if (!modname || !modelid) {
        return;
    }

    var uri = "modname=" + modname + "&modelid=" + modelid + "&ext_node_refer=" + referid + "&page=" + page;
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

                $("#hpm-nodels").empty();
                $("#hpm-termls").empty();
                l4i.InnerAlert(alertid, 'alert-info', "Item Not Found");
            } else {
                $(alertid).hide();
            }

            if (!rsj.model.extensions.node_refer) {
                hpNode.node_refer_back = null;
                $("#hpm-node-list-refer-back").css({
                    "display": "none"
                });
            } else {
                hpNode.node_refer_back = rsj.model.extensions.node_refer;
                $("#hpm-node-list-refer-back").css({
                    "display": "block"
                });
            }
            l4iStorage.Set("hpm_spec_active", modname);
            l4iStorage.Set("hpm_nmodel_active", modelid);
            $("#hpm-node-list-new-title").text("New " + rsj.model.title);

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

                if (!rsj.items[i].ext_node_refer) {
                    rsj.items[i].ext_node_refer = "";
                }
            }

            l4iTemplate.Render({
                dstid: "hpm-nodels",
                tplid: "hpm-nodels-tpl",
                data: {
                    model: rsj.model,
                    modname: modname,
                    modelid: modelid,
                    items: rsj.items,
                    _status_def: hpNode.status_def,
                },
                success: function() {

                    rsj.meta.RangeLen = 20;

                    l4iTemplate.Render({
                        dstid: "hpm-nodels-pager",
                        tplid: "hpm-nodels-pager-tpl",
                        data: l4i.Pager(rsj.meta),
                    });

                    hpNode.OpToolsRefresh("#hpm-node-list-opts");
                }
            });
        });

        ep.fail(function(err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        // template
        var el = document.getElementById("hpm-nodels");
        if (!el || el.length < 1) {
            hpMgr.TplCmd("node/list", {
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

        hpMgr.ApiCmd("node/list?" + uri, {
            callback: ep.done("data"),
        });
    });
}

hpNode.ListPage = function(page) {
    l4iStorage.Set("hpm_nodels_page", parseInt(page));
    hpNode.List();
}

hpNode.ListBatchSelectAll = function() {
    var form = $("#hpm-nodels");
    if (!form) {
        return;
    }

    var checked = false;
    if (form.find(".hpm-nodels-chk-all").is(':checked')) {
        checked = true;
    }

    form.find(".hpm-nodels-chk-item").each(function() {
        if (checked) {
            $(this).prop("checked", true);
        } else {
            $(this).prop("checked", false);
        }
    });

    hpNode.ListBatchSelectTodoBtnRefresh(checked);
}

hpNode.ListBatchSelectTodoBtnRefresh = function(onoff) {
    if (onoff !== true && onoff !== false) {

        onoff = false;

        $("#hpm-nodels").find(".hpm-nodels-chk-item").each(function() {

            if ($(this).is(":checked")) {
                onoff = true;
                return (false);
            }
        });
    }

    if (onoff === true) {
        $("#hpm-nodels-batch-select-todo-btn").css({
            "display": "block"
        });
    } else {
        $("#hpm-nodels-batch-select-todo-btn").css({
            "display": "none"
        });
    }
}

hpNode.ListBatchSelectTodo = function() {
    var form = $("#hpm-nodels");
    if (!form) {
        return;
    }

    var select_num = 0;

    form.find(".hpm-nodels-chk-item").each(function() {

        if ($(this).is(":checked")) {
            select_num++;
        }
    });

    var params = {
        select_num: select_num,
    };

    hpMgr.TplCmd("node/list-batch-set", {
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
                    onclick: "hpNode.ListBatchSelectTodoDelete()",
                    style: "btn-danger",
                }, {
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }],
            });
        },
    })
}

hpNode.ListBatchSelectTodoDelete = function(modname, modelid) {
    if (!modname && l4iStorage.Get("hpm_spec_active")) {
        modname = l4iStorage.Get("hpm_spec_active");
    }
    if (!modelid && l4iStorage.Get("hpm_nmodel_active")) {
        modelid = l4iStorage.Get("hpm_nmodel_active");
    }

    if (!modname || !modelid) {
        return;
    }

    var ids = [];

    $("#hpm-nodels").find(".hpm-nodels-chk-item").each(function() {

        if ($(this).is(":checked")) {
            ids.push($(this).val());
        }
    });

    var alertid = "#hpm-nodels-batch-set-alert";

    hpNode.DelBatch(modname, modelid, ids, function(err, data) {

        if (err) {
            l4i.InnerAlert(alertid, 'alert-danger', err);
        } else if (data && data.error) {
            l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
        } else if (data && data.kind == "Node") {
            l4i.InnerAlert(alertid, 'alert-success', "Successful operation");
            hpNode.List();
            setTimeout(function() {
                l4iModal.Close();
            }, 500);
        } else {
            l4i.InnerAlert(alertid, 'alert-danger', "unknown error");
        }
    });
}

hpNode.ReferBack = function() {
    if (hpNode.node_refer_back) {
        hpNode.List(null, hpNode.node_refer_back);
    }
}

hpNode.Set = function(modname, modelid, nodeid, referid) {
    var alertid = "#hpm-node-alert";

    if (!modname && l4iStorage.Get("hpm_spec_active")) {
        modname = l4iStorage.Get("hpm_spec_active");
    }
    if (!modelid && l4iStorage.Get("hpm_nmodel_active")) {
        modelid = l4iStorage.Get("hpm_nmodel_active");
    }
    if (!referid && l4iStorage.Get("hpm_node_refer_active")) {
        referid = l4iStorage.Get("hpm_node_refer_active");
    }

    // console.log(modname +","+ modelid +","+ nodeid);

    if (!modname || !modelid) {
        return;
    }

    hpEditor.Clean();
    hpNode.nodeOpToolsRefreshCurrent = null;

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

            if (!data.ext_node_refer) {
                data.ext_node_refer = "";
            }

            $(alertid).hide();

            hpNode.setCurrent = data;
            data._status_def = hpNode.status_def;

            // console.log(data);

            l4iTemplate.Render({
                dstid: "hpm-nodeset-laymain",
                tplid: "hpm-nodeset-tpl",
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

                    if (data.model.extensions.node_refer && data.model.extensions.node_refer != "") {
                        main_len += 1;
                    }

                    var field_layout_target = "hpm-nodeset-fields";
                    if (side_len > 0 && main_len > side_len) {
                        field_layout_target = "hpm-nodeset-layside";
                    } else {
                        $("#hpm-nodeset-layside").addClass("hpm-hide");
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

                                tplid = "hpm-nodeset-tplstring";
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
                                if (!field.attr_formats) {
                                    field.attr_formats = "text,html,md";
                                }

                                var formats = field.attr_formats.split(",");
                                var set_format = null;
                                field._formats = [];

                                for (var i in hpNode.text_formats) {
                                    if (formats.indexOf(hpNode.text_formats[i].name) > -1) {
                                        field._formats.push({
                                            name: hpNode.text_formats[i].name,
                                            value: hpNode.text_formats[i].value,
                                        });
                                        if (field.attr_format == hpNode.text_formats[i].name) {
                                            set_format = hpNode.text_formats[i].name;
                                        }
                                    }
                                }
                                if (field._formats.length < 1) {
                                    field._formats.push({
                                        name: hpNode.text_formats[0].name,
                                        value: hpNode.text_formats[0].value,
                                    });
                                }
                                if (!set_format) {
                                    field.attr_format = field._formats[0].name;
                                }

                                cb = function() {
                                    hpEditor.Open(field.name, field.attr_format);
                                };

                                tplid = "hpm-nodeset-tpltext";
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

                                tplid = "hpm-nodeset-tplint";
                                break;

                            default:
                                continue;
                        }

                        l4iTemplate.Render({
                            dstid: "hpm-nodeset-fields",
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

                                tplid = "hpm-nodeset-tplterm_tag";

                                l4iTemplate.Render({
                                    dstid: field_layout_target,
                                    tplid: tplid,
                                    prepend: true,
                                    data: data.model.terms[i],
                                });

                                break;

                            case "taxonomy":

                                if (!term.value) {
                                    term.value = "0";
                                }

                                hpMgr.ApiCmd("term/list?modname=" + modname + "&modelid=" + term.meta.name, {
                                    async: false,
                                    callback: function(err, term_list) {

                                        if (term_list.kind != "TermList") {
                                            return;
                                        }

                                        term_list.item = term;

                                        for (var i in term_list.items) {

                                            if (!term_list.items[i].pid) {
                                                term_list.items[i].pid = 0;
                                            }

                                            if (term_list.items[i].pid == 0) {
                                                term_list.items[i]._subs = hpTerm.ListSubRange(term_list.items, null, term_list.items[i].id, 0);
                                            }
                                        }

                                        tplid = "hpm-nodeset-tplterm_taxonomy";

                                        l4iTemplate.Render({
                                            dstid: field_layout_target,
                                            tplid: tplid,
                                            prepend: true,
                                            data: term_list,
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
                            tplid: "hpm-nodeset-tplext_comment_perentry",
                            append: true,
                            data: {
                                _general_onoff: hpNode.general_onoff,
                                ext_comment_perentry: data.ext_comment_perentry,
                            },
                        });
                    }

                    if (!data.ext_node_refer || data.ext_node_refer.length < 12) {
                        if (referid) {
                            data.ext_node_refer = referid;
                        }
                    }
                    if (data.model.extensions.node_refer && data.model.extensions.node_refer != "") {
                        l4iTemplate.Render({
                            dstid: "hpm-nodeset-tops",
                            tplid: "hpm-nodeset-tplext_node_refer",
                            append: true,
                            data: {
                                ext_node_refer: data.ext_node_refer,
                            },
                        });
                    }

                    if (data.model.extensions.permalink && data.model.extensions.permalink != "") {
                        l4iTemplate.Render({
                            dstid: "hpm-nodeset-tops",
                            tplid: "hpm-nodeset-tplext_permalink",
                            append: true,
                            data: {
                                ext_permalink_name: data.ext_permalink_name,
                            },
                        });
                    }

                    l4iTemplate.Render({
                        dstid: field_layout_target,
                        tplid: "hpm-nodeset-tplstatus",
                        append: true,
                        data: {
                            _status_def: hpNode.status_def,
                            status: data.status,
                        },
                    });

                    hpNode.OpToolsRefresh("#hpm-node-set-opts");

                    if (data.create_new) {
                        $("#hpm-node-set-opts-label").text("Create new Content");
                    } else {
                        $("#hpm-node-set-opts-label").text("Editing");
                    }

                    hpMgr.hotkey_ctrl_s = hpNode.SetSave;
                },
            });
        });

        ep.fail(function(err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-nodelist)");
        });

        hpMgr.TplCmd("node/set", {
            callback: function(err, tpl) {

                if (err) {
                    return ep.emit('error', err);
                }
                ep.emit("tpl", tpl);
            }
        });

        if (nodeid) {
            hpMgr.ApiCmd("node/entry?" + uri + "&id=" + nodeid, {
                callback: ep.done("data"),
            });
        } else {
            hpMgr.ApiCmd("node-model/entry?" + uri, {
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

hpNode.SetSave = function() {
    if (!hpNode.setCurrent) {
        hpMgr.hotkey_ctrl_s = null;
        return;
    }
    hpNode.SetCommit({
        save: true
    });
}

hpNode.SetCommit = function(options) {
    options = options || {};
    var form = $("#hpm-nodeset-layout"),
        alertid = "#hpm-node-alert";

    if (!hpNode.setCurrent) {
        return;
    }

    hpNode.setCurrent.title = form.find("input[name=title]").val();

    var req = {
        id: form.find("input[name=id]").val(),
        title: form.find("input[name=title]").val(),
        status: parseInt(form.find("select[name=status]").val()),
        fields: [],
        terms: [],
        ext_comment_perentry: form.find("select[name=ext_comment_perentry]").val(),
        ext_permalink_name: form.find("input[name=ext_permalink_name]").val(),
        ext_node_refer: form.find("input[name=ext_node_refer]").val(),
    }

    if (req.ext_comment_perentry && req.ext_comment_perentry == "false") {
        req.ext_comment_perentry = false;
    } else {
        req.ext_comment_perentry = true;
    }

    // return console.log(req);
    for (var i in hpNode.setCurrent.model.fields) {

        var field = hpNode.setCurrent.model.fields[i];

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
                field_set.value = hpEditor.Content(field.name);

                // console.log(format);

                // if (field.attrs) {

                //     for (var j in field.attrs) {

                //         if (field.attrs[j].key == "format" && field.attrs[j].value == "md") {

                //             field_set.value = hpEditor.Content(field.name);

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

    for (var i in hpNode.setCurrent.model.terms) {

        var term = hpNode.setCurrent.model.terms[i];

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

    // console.log(hpNode.setCurrent.model.terms);
    // console.log(JSON.stringify(req));

    //
    var uri = "modname=" + l4iStorage.Get("hpm_spec_active");
    uri += "&modelid=" + l4iStorage.Get("hpm_nmodel_active");

    hpMgr.ApiCmd("node/set?" + uri, {
        method: "POST",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (!data || data.kind != "Node") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            // console.log(data.id);
            form.find("input[name=id]").val(data.id);

            l4i.InnerAlert(alertid, 'alert-success', "Successful operation");
            if (options.save) {
                return;
            }
            hpNode.setCurrent = null;
            setTimeout(function() {
                hpNode.List();
                hpEditor.Clean();
            }, 500);
        }
    });
}


hpNode.Del = function(modname, modelid, id) {
    l4iModal.Open({
        title: "Delete",
        tplsrc: '<div id="hpm-node-del" class="alert alert-danger">Are you sure to delete this?</div>',
        height: "200px",
        buttons: [{
            title: "Confirm to delete",
            onclick: 'hpNode.DelCommit("' + modname + '","' + modelid + '","' + id + '")',
            style: "btn-danger",
        }, {
            title: "Cancel",
            onclick: "l4iModal.Close()",
        }],
    });
}

hpNode.DelCommit = function(modname, modelid, id) {
    var alertid = "#hpm-node-del";
    var uri = "modname=" + modname + "&modelid=" + modelid + "&id=" + id;

    hpMgr.ApiCmd("node/del?" + uri, {
        callback: function(err, data) {

            if (!data || data.kind != "Node") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful deleted");
            setTimeout(function() {
                hpNode.List();
                l4iModal.Close();
            }, 500);
        }
    });
}

hpNode.DelBatch = function(modname, modelid, ids, cb) {
    var uri = "modname=" + modname + "&modelid=" + modelid + "&id=" + ids.join(",");

    hpMgr.ApiCmd("node/del?" + uri, {
        callback: cb,
    });
}
