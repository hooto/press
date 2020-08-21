var gdoc = {
}


gdoc.PageEntryRender = function() {

    var idre = /^H\d{1}$/;
    var el = document.getElementById("hp-gdoc-page-entry-content");
    if (!el) {
        return;
    }
    var elo = document.getElementById("hp-gdoc-page-entry-toc");
    if (!elo) {
        return;
    }

    var items = el.childNodes;
    var nav = "";
    var last = 1;
    var num = 0;
    for (var i in items) {
        var item = items[i];
        if (!item.previousSibling) {
            continue;
        }
        if (!idre.test(item.previousSibling.nodeName)) {
            continue;
        }
        var level = parseInt(item.previousSibling.nodeName.substring(1));
        if (level > 3) {
            continue;
        } else if (level < 2) {
            level = 2;
        }
        if (level > last) {
            nav += "<ul>\n";
        } else if (level < last) {
            nav += "</li>\n</ul>\n";
        } else {
            nav += "</li>\n";
        }
        var tocid = "hp-gdoc-toc-" + i;
        nav += "<li><a href=\"#" + tocid + "\">" + item.previousSibling.outerText + "</a>";
        last = level;
        item.previousSibling.id = tocid;
        num += 1;
    }

    if (num < 2) {
        elo.style.display = "none";
        el.classList.replace("is-8", "is-10");
        return;
    }

    while (last > 1) {
        nav += "</li></ul>\n";
        last -= 1;
    }

    elo.innerHTML = "<nav class=\"hp-gdoc-page-toc-menu\">\n<h1>Page Nav</h1>" + nav + "</nav>";
}
