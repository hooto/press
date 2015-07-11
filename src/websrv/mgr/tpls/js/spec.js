var l5sSpec = {

}

l5sSpec.Index = function()
{
    l5sMgr.Ajax(l5sMgr.base +"-/spec/index.tpl", {
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
                rsj.items[i].metadata.updated = l4i.TimeParseFormat(rsj.items[i].metadata.updated, "Y-m-d");
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
            l5sMgr.Ajax("-/spec/list.tpl", {
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
    
        l5sMgr.Ajax("/v1/spec/list?"+ uri, {
            callback: ep.done("data"),           
        });
    });
}

