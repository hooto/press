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

var hpressSpec = {

    specdef: {
        kind: "Spec",
        meta: {
            id: "",
            name: "",
        },
		srvname: "",
        title: "",
    },

    termdef: {
        kind: "TermModel",
        meta: {
            name: "",
        },
        type: "taxonomy",
        title: "",
    },

    nodedef: {
        kind: "NodeModel",
        meta: {
            name: "",
        },
        title: "",
        fields: [],
        extensions: {
            access_counter: false,
            comment_enable: false,
            comment_perentry: false,
        },
    },

    actiondef: {
        kind: "SpecAction",
        name: "",
        datax: [],
    },

    action_dataxdef: {
        name: "",
        type: "node.list",
        pager: false,
        query: {
            table: "",
            limit: 10,
            order: "",
        },
        cache_ttl: 0,
    },

    routedef: {
        kind: "SpecRoute",
        path: "",
        dataAction: "",
        template: "",
        params: {},
    },

    field_idx_typedef: [{
        type: 0,
        name: "No Index",
    }, {
        type: 1,
        name: "General Index",
    }, {
        type: 2,
        name: "Unique Index",
    // }, {
    //     type : 3,
    //     name : "Primary Key",
    }],

    term_typedef: [{
        type: "taxonomy",
        name: "Categories",
    }, {
        type: "tag",
        name: "Tags",
    }],

    // datax_typedef : [{
    //     type : "node.list",
    //     name : "Node List",
    // },{
    //     type : "node.entry",
    //     name : "Node Entry",
    // }],
    datax_typedef: [{
        type: "list",
        name: "List",
    }, {
        type: "entry",
        name: "Entry",
    }],

    field_typedef: [{
        type: "bool",
        name: "Bool",
    }, {
        type: "string",
        name: "Varchar",
    }, {
        type: "text",
        name: "Text",
    }, {
        type: "date",
        name: "Date",
    }, {
        type: "datetime",
        name: "Datetime",
    }, {
        type: "int8",
        name: "int8",
    }, {
        type: "uint8",
        name: "uint8",
    }, {
        type: "int16",
        name: "int16",
    }, {
        type: "uint16",
        name: "uint16",
    }, {
        type: "int32",
        name: "int32",
    }, {
        type: "uint32",
        name: "uint32",
    }, {
        type: "int64",
        name: "int64",
    }, {
        type: "uint64",
        name: "uint64",
    }, {
        type: "float",
        name: "Float",
    }, {
        type: "decimal",
        name: "Decimal Float",
    }],

    general_onoff: [{
        type: true,
        name: "ON",
    }, {
        type: false,
        name: "OFF",
    }],

    permalink_def: [{
        type: "",
        name: "Off",
    }, {
        type: "name",
        name: "Name",
    }, {
        type: "day-name",
        name: "Day and Name",
    }],
}


hpressSpec.Init = function() {
    l4i.UrlEventRegister("spec/index", hpressSpec.Index, "hpressm-topbar");
}

hpressSpec.Index = function() {
    l4iStorage.Set("hpressm_nav_last_active", "spec/index");

    hpressMgr.TplCmd("spec/index", {
        callback: function(err, data) {

            $("#com-content").html(data);

            hpressSpec.List();
        }
    });
}

hpressSpec.List = function() {

    var uri = "";
    if (document.getElementById("qry_text")) {
        uri = "qry_text=" + $("#qry_text").val();
    }

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create("tpl", "data", function(tpl, rsj) {

            if (tpl) {
                $("#work-content").html(tpl);
            }
            // console.log(tpl);
            // if (data typeof object)
            // var rsj = JSON.parse(data);

            if (rsj === undefined || rsj.kind != "SpecList" ||
                rsj.items === undefined || rsj.items.length < 1) {
                return l4i.InnerAlert("#hpressm-specls-alert", 'alert-info', "Item Not Found");
            }

            $("#hpressm-specls-alert").hide();

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

                if (rsj.items[i].views) {
                    rsj.items[i]._viewsNum = rsj.items[i].views.length;
                } else {
                    rsj.items[i]._viewsNum = 0;
                }

                if (rsj.items[i].router.routes) {
                    rsj.items[i]._routesNum = rsj.items[i].router.routes.length;
                } else {
                    rsj.items[i]._routesNum = 0;
                }

                if (!rsj.items[i].meta.created) {
                    rsj.items[i].meta.created = rsj.items[i].meta.updated;
                }
            }

            l4iTemplate.Render({
                dstid: "hpressm-specls",
                tplid: "hpressm-specls-tpl",
                data: {
                    items: rsj.items,
                },
            });
        });

        ep.fail(function(err) {
            // TODO
            alert("SpecListRefresh error, Please try again later (EC:app-speclist)");
        });

        // template
        var el = document.getElementById("hpressm-specls");
        if (!el || el.length < 1) {
            hpressMgr.TplCmd("spec/list", {
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

        // hpressMgr.Ajax("-/spec/list.tpl", {
        //     callback: ep.done("tpl"),
        // });

        hpressMgr.ApiCmd("mod-set/spec-list?" + uri, {
            callback: ep.done("data"),
        });
    });
}

hpressSpec.InfoSet = function(name) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind || data.kind != "Spec") {
                return alert("Spec Not Found");
            }

            var ptitle = "Info Settings";
            if (!name) {
                ptitle = "New Module";
            }

            l4iModal.Open({
                tplsrc: tpl,
                width: 600,
                height: 400,
                title: ptitle,
                data: data,
                success: function() {},
                buttons: [{
                    onclick: "l4iModal.Close()",
                    title: "Close",
                }, {
                    onclick: "hpressSpec.InfoSetCommit()",
                    title: "Save",
                    style: "btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error, Please try again later " + err);
        });

        hpressMgr.TplCmd("spec/info-set", {
            callback: ep.done('tpl'),
        });

        if (name) {

            hpressMgr.ApiCmd("mod-set/spec-entry?name=" + name, {
                callback: ep.done('data'),
            });
        } else {

            ep.emit("data", l4i.Clone(hpressSpec.specdef));
        }

    });

}

hpressSpec.InfoSetCommit = function() {
    var form = $("#hpressm-specset");
    var alertid = "#hpressm-specset-alert";

    var req = {
        meta: {
            name: form.find("input[name=name]").val(),
        },
        srvname: form.find("input[name=srvname]").val(),
        title: form.find("input[name=title]").val(),
    };

    // console.log(req);

    hpressMgr.ApiCmd("mod-set/spec-info-set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            // console.log(data);

            if (!data || data.error || data.kind != "Spec") {

                if (data.error) {
                    return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
                }

                return l4i.InnerAlert(alertid, 'alert-danger', 'Network Connection Exception');
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful updated");

            hpressSpec.List();

            window.setTimeout(function() {
                l4iModal.Close();
            }, 1000);
        },
    });
}


// Spec::Term
hpressSpec.TermList = function(modname) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind || data.kind != "Spec") {

                if (data.error) {
                    return alert(data.error.message);
                }

                return alert("SpecTermList Not Found");
            }

            if (!data.termModels) {
                data.termModels = [];
            }

            l4iModal.Open({
                id: "term-model-ls",
                tplsrc: tpl,
                width: 900,
                height: 500,
                title: "Term List",
                // data   : data,
                success: function() {

                    l4iTemplate.Render({
                        dstid: "hpressm-spec-termls",
                        tplid: "hpressm-spec-termls-tpl",
                        data: data,
                    });
                },
                buttons: [{
                    onclick: "hpressSpec.TermSet(\"" + modname + "\")",
                    title: "New Term",
                    style: "btn-primary",
                }, {
                    onclick: "l4iModal.Close()",
                    title: "Close",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error, Please try again later " + err);
        });

        hpressMgr.TplCmd("spec/term/list", {
            callback: ep.done('tpl'),
        });

        hpressMgr.ApiCmd("mod-set/spec-entry?name=" + modname, {
            callback: ep.done('data'),
        });
    });
}


hpressSpec.TermSet = function(modname, modelid) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind || data.kind != "TermModel") {

                if (data.error) {
                    return alert(data.error.message);
                }

                return alert("TermModel Not Found");
            }

            data._modname = modname;

            var ptitle = "Term Settings";
            if (!modelid) {
                ptitle = "New Term";
            }

            l4iModal.Open({
                id: "term-model-set",
                tplsrc: tpl,
                title: ptitle,
                data: data,
                success: function() {},
                buttons: [{
                    onclick: "l4iModal.Close()",
                    title: "Close",
                }, {
                    onclick: "hpressSpec.TermSetCommit()",
                    title: "Save",
                    style: "btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error, Please try again later " + err);
        });

        hpressMgr.TplCmd("spec/term/set", {
            callback: ep.done('tpl'),
        });

        if (modelid) {

            hpressMgr.ApiCmd("term-model/entry?modname=" + modname + "&modelid=" + modelid, {
                callback: ep.done('data'),
            });
        } else {

            ep.emit("data", l4i.Clone(hpressSpec.termdef));
        }
    });
}

hpressSpec.TermSetCommit = function() {
    var form = $("#hpressm-spec-termset"),
        alertid = "#hpressm-spec-termset-alert";

    var req = {
        meta: {
            name: form.find("input[name=name]").val(),
        },
        type: form.find("select[name=type]").val(),
        title: form.find("input[name=title]").val(),
        modname: form.find("input[name=modname]").val(),
    };

    hpressMgr.ApiCmd("mod-set/spec-term-set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            // console.log(data);

            if (!data || !data.kind || data.kind != "TermModel") {

                if (data.error) {
                    return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
                }

                return l4i.InnerAlert(alertid, 'alert-danger', 'Network Connection Exception');
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful updated");

            window.setTimeout(function() {

                l4iModal.Prev(function() {

                    hpressSpec.List();

                    hpressMgr.ApiCmd("mod-set/spec-entry?name=" + req.modname, {
                        callback: function(err, data) {

                            if (err || !data || !data.kind || data.kind != "Spec") {
                                return;
                            }

                            l4iTemplate.Render({
                                dstid: "hpressm-spec-termls",
                                tplid: "hpressm-spec-termls-tpl",
                                data: data,
                            });
                        },
                    });
                });

            }, 1000);
        },
    });
}



// Spec::Node
hpressSpec.NodeList = function(modname) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind || data.kind != "Spec") {

                if (data.error) {
                    return alert(data.error.message);
                }

                return alert("SpecNodeList Not Found");
            }

            if (!data.nodeModels) {
                data.nodeModels = [];
            }

            for (var i in data.nodeModels) {

                if (!data.nodeModels[i].fields) {
                    data.nodeModels[i].fields = [];
                }

                data.nodeModels[i]._fieldsNum = data.nodeModels[i].fields.length;

                if (!data.nodeModels[i].terms) {
                    data.nodeModels[i].terms = [];
                }

                data.nodeModels[i]._termsNum = data.nodeModels[i].terms.length;
            }

            l4iModal.Open({
                id: "node-model-ls",
                tplsrc: tpl,
                width: 700,
                height: 400,
                title: "Node List",
                // data   : data,
                success: function() {

                    l4iTemplate.Render({
                        dstid: "hpressm-spec-nodels",
                        tplid: "hpressm-spec-nodels-tpl",
                        data: data,
                    });
                },
                buttons: [{
                    onclick: "hpressSpec.NodeSet(\"" + modname + "\")",
                    title: "New Node",
                    style: "btn-primary",
                }, {
                    onclick: "l4iModal.Close()",
                    title: "Close",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error, Please try again later " + err);
        });

        hpressMgr.TplCmd("spec/node/list", {
            callback: ep.done('tpl'),
        });

        hpressMgr.ApiCmd("mod-set/spec-entry?name=" + modname, {
            callback: ep.done('data'),
        });
    });
}


hpressSpec.NodeSet = function(modname, modelid) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind || data.kind != "NodeModel") {

                if (data.error) {
                    return alert(data.error.message);
                }

                return alert("NodeModel Not Found");
            }

            data._modname = modname;

            var ptitle = "Node Settings";
            if (!modelid) {
                ptitle = "New Node";
            }

            //
            for (var i in data.fields) {

                if (!data.fields[i].length) {
                    data.fields[i].length = 0;
                }

                if (!data.fields[i].length) {
                    data.fields[i].indexType = 0;
                }

                data.fields[i]._seqid = Math.random().toString(16).slice(2);
            }

            if (!data.extensions) {
                data.extensions = {};
            }
            if (!data.extensions.access_counter) {
                data.extensions.access_counter = false;
            }
            if (!data.extensions.comment_enable) {
                data.extensions.comment_enable = false;
            }
            if (!data.extensions.comment_perentry) {
                data.extensions.comment_perentry = false;
            }

            data._field_idx_typedef = hpressSpec.field_idx_typedef;
            data._field_typedef = hpressSpec.field_typedef;
            data._general_onoff = hpressSpec.general_onoff;
            data._permalink_def = hpressSpec.permalink_def;

            //
            if (!data.terms) {
                data.terms = [];
            }

            for (var i in data.terms) {
                data.terms[i]._seqid = Math.random().toString(16).slice(2);
            }

            data._term_typedef = hpressSpec.term_typedef;
            // data._field_termdef = hpressSpec.termdef;


            l4iModal.Open({
                id: "node-model-set",
                tplsrc: tpl,
                title: ptitle,
                data: data,
                width: 980,
                height: 600,
                success: function() {},
                buttons: [{
                    onclick: "hpressSpec.NodeSetFieldAppend()",
                    title: "New Field",
                    style: "btn-primary",
                }, {
                    onclick: "hpressSpec.NodeSetTermAppend()",
                    title: "New Term",
                    style: "btn-primary",
                }, {
                    onclick: "hpressSpec.NodeSetCommit()",
                    title: "Save",
                    style: "btn-primary",
                }, {
                    onclick: "l4iModal.Close()",
                    title: "Close",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error, Please try again later " + err);
        });

        hpressMgr.TplCmd("spec/node/set", {
            callback: ep.done('tpl'),
        });

        if (modelid) {

            hpressMgr.ApiCmd("node-model/entry?modname=" + modname + "&modelid=" + modelid, {
                callback: ep.done('data'),
            });
        } else {

            ep.emit("data", l4i.Clone(hpressSpec.nodedef));
        }
    });
}

hpressSpec.NodeSetFieldAppend = function() {
    l4iTemplate.Render({
        dstid: "hpressm-spec-node-fields",
        tplid: "hpressm-spec-node-field-item-tpl",
        append: true,
        data: {
            _field_typedef: hpressSpec.field_typedef,
            _field_idx_typedef: hpressSpec.field_idx_typedef,
            _type: "bool",
            _indexType: 0,
            _seqid: Math.random().toString(16).slice(2),
        },
    });
}

hpressSpec.NodeSetFieldAttrAppend = function(seqid) {
    var dst = document.getElementById("hpressm-spec-node-field-attr-item-tpl");
    if (!dst) {
        return;
    }
    var source = dst.value || dst.innerHTML;

    $("#field-seq-" + seqid).find(".hpressm-spec-node-field-attrs").append(source);
}

hpressSpec.NodeSetTermAppend = function(seqid) {
    l4iTemplate.Render({
        dstid: "hpressm-spec-node-terms",
        tplid: "hpressm-spec-node-term-item-tpl",
        append: true,
        data: {
            _term_typedef: hpressSpec.term_typedef,
            _type: "taxonomy",
            _seqid: Math.random().toString(16).slice(2),
        },
    });
}

hpressSpec.NodeSetCommit = function() {
    var form = $("#hpressm-spec-nodeset"),
        alertid = "#hpressm-spec-nodeset-alert",
        namereg = /^[a-z][a-z0-9_]+$/;

    var req = {
        meta: {
            name: form.find("input[name=name]").val(),
        },
        title: form.find("input[name=title]").val(),
        modname: form.find("input[name=modname]").val(),
        fields: [],
        terms: [],
        extensions: {
            access_counter: false,
            comment_enable: false,
            comment_perentry: false,
        },
    };

    if (form.find("select[name=ext_access_counter]").val() == "true") {
        req.extensions.access_counter = true;
    }

    if (form.find("select[name=ext_comment_enable]").val() == "true") {
        req.extensions.comment_enable = true;
    }

    if (form.find("select[name=ext_comment_perentry]").val() == "true") {
        req.extensions.comment_perentry = true;
    }

    req.extensions.permalink = form.find("select[name=ext_permalink]").val();

    try {

        form.find(".hpressm-spec-node-field-item").each(function() {

            var field = {
                name: $(this).find("input[name=field_name]").val(),
                title: $(this).find("input[name=field_title]").val(),
                type: $(this).find("select[name=field_type]").val(),
                length: $(this).find("input[name=field_length]").val(),
                indexType: parseInt($(this).find("select[name=field_index_type]").val()),
                attrs: [],
            };

            if (!field.name || field.name == "") {
                return;
            }

            if (!namereg.test(field.name)) {
                throw "Invalid Field Name : " + field.name;
            }

            $(this).find(".hpressm-spec-node-field-attrs").each(function() {

                var attr_key = $(this).find("input[name=field_attr_key]").val();

                if (attr_key) {

                    if (!namereg.test(attr_key)) {
                        throw "Invalid Field Attribute Key : " + attr_key;
                    }

                    field.attrs.push({
                        key: attr_key,
                        value: $(this).find("input[name=field_attr_value]").val(),
                    });
                }
            });

            req.fields.push(field);
        });

        form.find(".hpressm-spec-node-term-item").each(function() {

            var term = {
                meta: {
                    name: $(this).find("input[name=term_name]").val(),
                },
                title: $(this).find("input[name=term_title]").val(),
                type: $(this).find("select[name=term_type]").val(),
            };

            if (!term.meta.name || term.meta.name == "") {
                return;
            }

            if (!namereg.test(term.meta.name)) {
                throw "Invalid Term Name : " + term.meta.name;
            }

            req.terms.push(term);
        });

    } catch (err) {
        l4i.InnerAlert(alertid, 'alert-danger', err);
        return;
    }

    // console.log(req);


    hpressMgr.ApiCmd("mod-set/spec-node-set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            // console.log(data);

            if (!data || !data.kind || data.kind != "NodeModel") {

                if (data.error) {
                    return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
                }

                return l4i.InnerAlert(alertid, 'alert-danger', 'Network Connection Exception');
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful updated");

            window.setTimeout(function() {

                l4iModal.Prev(function() {

                    hpressSpec.List();

                    hpressMgr.ApiCmd("mod-set/spec-entry?name=" + req.modname, {
                        callback: function(err, data) {

                            if (err || !data || !data.kind || data.kind != "Spec") {
                                return;
                            }

                            if (!data.nodeModels) {
                                data.nodeModels = [];
                            }

                            for (var i in data.nodeModels) {

                                if (!data.nodeModels[i].fields) {
                                    data.nodeModels[i].fields = [];
                                }

                                data.nodeModels[i]._fieldsNum = data.nodeModels[i].fields.length;

                                if (!data.nodeModels[i].terms) {
                                    data.nodeModels[i].terms = [];
                                }

                                data.nodeModels[i]._termsNum = data.nodeModels[i].terms.length;
                            }

                            l4iTemplate.Render({
                                dstid: "hpressm-spec-nodels",
                                tplid: "hpressm-spec-nodels-tpl",
                                data: data,
                            });
                        },
                    });
                });

            }, 1000);
        },
    });
}


// Spec::Action
hpressSpec.ActionList = function(modname) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind || data.kind != "Spec") {

                if (data.error) {
                    return alert(data.error.message);
                }

                return alert("SpecActionList Not Found");
            }

            if (!data.actions) {
                data.actions = [];
            }

            for (var i in data.actions) {

                if (!data.actions[i].datax) {
                    data.actions[i].datax = [];
                }

                data.actions[i]._dataxNum = data.actions[i].datax.length;
            }

            data._modname = modname;

            l4iModal.Open({
                id: "spec-action-ls",
                tplsrc: tpl,
                width: 700,
                height: 400,
                title: "Action List",
                // data   : data,
                success: function() {

                    l4iTemplate.Render({
                        dstid: "hpressm-spec-actionls",
                        tplid: "hpressm-spec-actionls-tpl",
                        data: data,
                    });
                },
                buttons: [{
                    onclick: "hpressSpec.ActionSet(\"" + modname + "\")",
                    title: "New Action",
                    style: "btn-primary",
                }, {
                    onclick: "l4iModal.Close()",
                    title: "Close",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error, Please try again later " + err);
        });

        hpressMgr.TplCmd("spec/action/list", {
            callback: ep.done('tpl'),
        });

        hpressMgr.ApiCmd("mod-set/spec-entry?name=" + modname, {
            callback: ep.done('data'),
        });
    });
}


hpressSpec.ActionSet = function(modname, modelid) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'nodeModels', 'termModels', 'data', function(tpl, nodeModels, termModels, data) {

            // console.log(data);
            // return;

            if (!data || !data.kind || data.kind != "SpecAction") {

                if (data.error) {
                    return alert(data.error.message);
                }

                return alert("SpecAction Not Found");
            }

            data._modname = modname;

            var ptitle = "Action Settings";
            if (!modelid) {
                ptitle = "New Action";
            }

            //
            if (!data.datax) {
                data.datax = [];
            }

            //
            for (var i in data.datax) {

                data.datax[i]._seqid = Math.random().toString(16).slice(2);

                if (!data.datax[i].pager) {
                    data.datax[i].pager = false;
                }

                if (!data.datax[i].query.limit) {
                    data.datax[i].query.limit = 1;
                }

                if (!data.datax[i].query.order) {
                    data.datax[i].query.order = "";
                }

                if (!data.datax[i].cache_ttl) {
                    data.datax[i].cache_ttl = 0;
                }
            }

            data._nodeModels = nodeModels;
            data._termModels = termModels;

            data._datax_typedef = l4i.Clone(hpressSpec.datax_typedef);

            // console.log(data);
            // return;

            l4iModal.Open({
                id: "spec-action-set",
                tplsrc: tpl,
                title: ptitle,
                data: data,
                width: 980,
                height: 550,
                success: function() {

                    if (!modelid) {
                        hpressSpec.ActionSetDataxAppend(modname);
                    }
                },
                buttons: [{
                    onclick: "hpressSpec.ActionSetDataxAppend(\"" + modname + "\")",
                    title: "New Datax",
                    style: "btn-primary",
                }, {
                    onclick: "hpressSpec.ActionSetCommit()",
                    title: "Save",
                    style: "btn-primary",
                }, {
                    onclick: "l4iModal.Close()",
                    title: "Close",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error, Please try again later " + err);
        });

        hpressMgr.TplCmd("spec/action/set", {
            callback: ep.done('tpl'),
        });

        hpressMgr.ApiCmd("mod-set/spec-entry?name=" + modname, {
            callback: function(err, data) {

                if (err) {
                    ep.emit("error", err);
                    return;
                }

                // console.log(data);

                if (!data || !data.kind || data.kind != "Spec") {
                    ep.emit("error", "Spec Not Found");
                    return;
                }

                //
                if (!data.nodeModels) {
                    data.nodeModels = [];
                }
                ep.emit("nodeModels", data.nodeModels);

                //
                if (!data.termModels) {
                    data.termModels = [];
                }
                ep.emit("termModels", data.termModels);

                //
                if (modelid) {

                    for (var i in data.actions) {

                        if (data.actions[i].name == modelid) {
                            data.actions[i].kind = "SpecAction";
                            ep.emit("data", data.actions[i]);
                            return;
                        }
                    }

                    ep.emit("error", "Spec Not Found");

                } else {
                    ep.emit("data", l4i.Clone(hpressSpec.actiondef));
                }
            },
        });
    });
}

hpressSpec.ActionSetDataxAppend = function(modname) {
    hpressMgr.ApiCmd("mod-set/spec-entry?name=" + modname, {
        callback: function(err, data) {

            if (err) {
                return alert(err);
            }

            if (!data || !data.kind || data.kind != "Spec") {
                return alert("Spec Not Found");
            }

            //
            if (!data.nodeModels) {
                data.nodeModels = [];
            }

            //
            if (!data.termModels) {
                data.termModels = [];
            }

            var action = l4i.Clone(hpressSpec.actiondef);

            action._nodeModels = data.nodeModels;
            action._termModels = data.termModels;
            action._datax_typedef = hpressSpec.datax_typedef;
            action._seqid = Math.random().toString(16).slice(2);

            l4iTemplate.Render({
                dstid: "hpressm-spec-action-dataxs",
                tplid: "hpressm-spec-action-datax-item-tpl",
                append: true,
                data: action,
            });
        },
    });
}

hpressSpec.ActionSetCommit = function() {
    var form = $("#hpressm-spec-actionset"),
        alertid = "#hpressm-spec-actionset-alert",
        namereg = /^[a-z][a-z0-9_]+$/;

    var req = {
        name: form.find("input[name=name]").val(),
        modname: form.find("input[name=modname]").val(),
        datax: [],
    };

    if (!namereg.test(req.name)) {
        return l4i.InnerAlert(alertid, 'alert-danger', 'Invalid Action Name');
    }

    if (!req.modname || req.modname == "") {
        return l4i.InnerAlert(alertid, 'alert-danger', 'Invalid Module Name');
    }

    try {

        form.find(".hpressm-spec-action-datax-item").each(function() {

            var datax = {
                name: $(this).find("input[name=datax_name]").val(),
                type: $(this).find("select[name=datax_type]").val(),
                query: {
                    table: $(this).find("select[name=datax_query_table]").val(),
                    limit: parseInt($(this).find("input[name=datax_query_limit]").val()),
                    order: $(this).find("input[name=datax_query_order]").val(),
                },
                pager: $(this).find("select[name=datax_pager]").val(),
                cache_ttl: parseInt($(this).find("input[name=datax_cache_ttl]").val()),
            };

            if (!datax.name || datax.name == "") {
                return;
            }

            if (!namereg.test(datax.name)) {
                throw "Invalid Datax Name : " + datax.name;
            }

            if (datax.pager == "true") {
                datax.pager = true;
            } else {
                datax.pager = false;
            }

            if (datax.type != "list" && datax.type != "entry") {
                datax.type = "list";
            }

            if (datax.query.table.substr(0, 5) == "node.") {
                datax.type = "node." + datax.type;
            } else if (datax.query.table.substr(0, 5) == "term.") {
                datax.type = "term." + datax.type;
            } else {
                throw "Invalid Query Table Name : " + datax.query.table;
            }

            if (!namereg.test(datax.query.table.slice(5))) {
                throw "Invalid Query Table Name : " + datax.query.table.slice(5);
            }

            datax.query.table = datax.query.table.slice(5);

            req.datax.push(datax);
        });

    } catch (err) {
        l4i.InnerAlert(alertid, 'alert-danger', err);
        return;
    }

    hpressMgr.ApiCmd("mod-set/spec-action-set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            // console.log(data);

            if (!data || !data.kind || data.kind != "Action") {

                if (data.error) {
                    return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
                }

                return l4i.InnerAlert(alertid, 'alert-danger', 'Network Connection Exception');
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful updated");

            window.setTimeout(function() {

                l4iModal.Prev(function() {

                    hpressSpec.List();

                    hpressMgr.ApiCmd("mod-set/spec-entry?name=" + req.modname, {
                        callback: function(err, data) {

                            if (err || !data || !data.kind || data.kind != "Spec") {
                                return;
                            }

                            if (!data.actions) {
                                data.actions = [];
                            }

                            for (var i in data.actions) {

                                if (!data.actions[i].datax) {
                                    data.actions[i].datax = [];
                                }

                                data.actions[i]._dataxNum = data.actions[i].datax.length;
                            }

                            data._modname = req.modname;

                            l4iTemplate.Render({
                                dstid: "hpressm-spec-actionls",
                                tplid: "hpressm-spec-actionls-tpl",
                                data: data,
                            });
                        },
                    });
                });

            }, 1000);
        },
    });
}



// Spec::Router
hpressSpec.RouteList = function(modname) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind || data.kind != "Spec") {

                if (data.error) {
                    return alert(data.error.message);
                }

                return alert("SpecRouteList Not Found");
            }

            data._modname = modname;

            if (!data.actions) {
                data.actions = [];
            }

            if (!data.router.routes) {
                data.router.routes = [];
            }

            for (var i in data.router.routes) {

                if (!data.router.routes[i].params) {
                    data.router.routes[i].params = {};
                }

                data.router.routes[i]._paramsNum = 0;
                for (var j in data.router.routes[i].params) {
                    data.router.routes[i]._paramsNum++;
                }
            }

            l4iModal.Open({
                id: "spec-route-ls",
                tplsrc: tpl,
                width: 900,
                height: 500,
                title: "Route List",
                success: function() {

                    l4iTemplate.Render({
                        dstid: "hpressm-spec-routels",
                        tplid: "hpressm-spec-routels-tpl",
                        data: data,
                    });
                },
                buttons: [{
                    onclick: "hpressSpec.RouteSet(\"" + modname + "\")",
                    title: "New Route",
                    style: "btn-primary",
                }, {
                    onclick: "l4iModal.Close()",
                    title: "Close",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error, Please try again later " + err);
        });

        hpressMgr.TplCmd("spec/router/list", {
            callback: ep.done('tpl'),
        });

        hpressMgr.ApiCmd("mod-set/spec-entry?name=" + modname, {
            callback: ep.done('data'),
        });
    });
}


hpressSpec.RouteSet = function(modname, modelid) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'actions', 'data', function(tpl, actions, data) {

            if (!data || !data.kind || data.kind != "SpecRoute") {

                if (data.error) {
                    return alert(data.error.message);
                }

                return alert("SpecRoute Not Found");
            }

            data._modname = modname;
            data._actions = actions;

            var ptitle = "Route Settings";
            if (!modelid) {
                ptitle = "New Route";
            }

            if (!data.params) {
                data.params = [];
            }

            l4iModal.Open({
                id: "spec-route-set",
                tplsrc: tpl,
                title: ptitle,
                success: function() {

                    l4iTemplate.Render({
                        dstid: "hpressm-spec-routeset",
                        tplid: "hpressm-spec-routeset-tpl",
                        data: data,
                        success: function() {

                            for (var i in data.params) {
                                l4iTemplate.Render({
                                    dstid: "hpressm-spec-route-params",
                                    tplid: "hpressm-spec-route-param-item-tpl",
                                    append: true,
                                    data: {
                                        _key: i,
                                        _value: data.params[i],
                                    },
                                });
                            }
                        }
                    });
                },
                buttons: [{
                    onclick: "l4iModal.Close()",
                    title: "Close",
                }, {
                    onclick: "hpressSpec.RouteSetParamAppend()",
                    title: "New Param",
                    style: "btn-primary",
                }, {
                    onclick: "hpressSpec.RouteSetCommit()",
                    title: "Save",
                    style: "btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error, Please try again later " + err);
        });

        hpressMgr.TplCmd("spec/router/set", {
            callback: ep.done('tpl'),
        });

        hpressMgr.ApiCmd("mod-set/spec-entry?name=" + modname, {
            callback: function(err, data) {

                if (err) {
                    ep.emit("error", err);
                    return;
                }

                if (!data || !data.kind || data.kind != "Spec") {
                    ep.emit("error", "Spec Not Found");
                    return;
                }

                if (!data.actions) {
                    data.actions = [];
                }
                ep.emit("actions", data.actions);

                //
                if (modelid) {

                    for (var i in data.router.routes) {

                        if (data.router.routes[i].path == modelid) {
                            data.router.routes[i].kind = "SpecRoute";
                            ep.emit("data", data.router.routes[i]);
                            return;
                        }
                    }

                    ep.emit("error", "Spec Not Found");

                } else {
                    ep.emit("data", l4i.Clone(hpressSpec.routedef));
                }
            },
        });
    });
}

hpressSpec.RouteSetParamAppend = function() {
    l4iTemplate.Render({
        dstid: "hpressm-spec-route-params",
        tplid: "hpressm-spec-route-param-item-tpl",
        append: true,
        data: {
            _seqid: Math.random().toString(16).slice(2),
            _key: "",
            _value: "",
        },
    });
}

hpressSpec.RouteSetCommit = function() {
    var form = $("#hpressm-spec-routeset"),
        alertid = "#hpressm-spec-routeset-alert",
        namereg = /^[a-z][a-z0-9_]+$/;
    ;

    var req = {
        path: form.find("input[name=path]").val(),
        dataAction: form.find("select[name=data_action]").val(),
        template: form.find("input[name=template]").val(),
        modname: form.find("input[name=modname]").val(),
        params: {},
    };

    try {

        form.find(".hpressm-spec-route-param-item").each(function() {

            var param_key = $(this).find("input[name=param_key]").val(),
                param_value = $(this).find("input[name=param_value]").val();

            if (!param_key || !param_value) {
                return;
            }

            if (!namereg.test(param_key)) {
                throw "Invalid Param Name : " + param_key;
            }

            req.params[param_key] = param_value;
        });

    } catch (err) {
        l4i.InnerAlert(alertid, 'alert-danger', err);
        return;
    }

    hpressMgr.ApiCmd("mod-set/spec-route-set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (!data || !data.kind || data.kind != "SpecRoute") {

                if (data.error) {
                    return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
                }

                return l4i.InnerAlert(alertid, 'alert-danger', 'Network Connection Exception');
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful updated");

            window.setTimeout(function() {

                l4iModal.Prev(function() {

                    hpressSpec.List();

                    hpressMgr.ApiCmd("mod-set/spec-entry?name=" + req.modname, {
                        callback: function(err, data) {

                            if (err || !data || !data.kind || data.kind != "Spec") {
                                return;
                            }

                            data._modname = req.modname;

                            if (!data.router.routes) {
                                data.router.routes = [];
                            }

                            for (var i in data.router.routes) {

                                if (!data.router.routes[i].params) {
                                    data.router.routes[i].params = {};
                                }

                                data.router.routes[i]._paramsNum = 0;
                                for (var j in data.router.routes[i].params) {
                                    data.router.routes[i]._paramsNum++;
                                }
                            }

                            l4iTemplate.Render({
                                dstid: "hpressm-spec-routels",
                                tplid: "hpressm-spec-routels-tpl",
                                data: data,
                            });
                        },
                    });
                });

            }, 1000);
        },
    });
}


hpressSpec.RouteSetTemplateSelect = function(modname) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind || data.kind != "SpecTemplateList") {

                if (data.error) {
                    return alert(data.error.message);
                }

                return alert("SpecTemplateList Not Found");
            }

            data._modname = modname;

            l4iModal.Open({
                id: "spec-route-template-select",
                tplsrc: tpl,
                title: "Select a Template",
                data: data,
                success: function() {},
                buttons: [{
                    onclick: "l4iModal.Close()",
                    title: "Close",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error, Please try again later " + err);
        });

        hpressMgr.TplCmd("spec/view/list", {
            callback: ep.done('tpl'),
        });

        hpressMgr.ApiCmd("mod-set/fs-tpl-list?modname=" + modname, {
            callback: ep.done('data'),
        });
    });
}

hpressSpec.RouteSetTemplateSelectOne = function(path) {
    l4iModal.Prev(function() {
        $("#hpressm-spec-routeset-template").attr("value", path);
    });
}
