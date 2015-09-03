// CodeMirror, copyright (c) by Marijn Haverbeke and others
// Distributed under an MIT license: http://codemirror.net/LICENSE
/**
 * Tag-closer extension for CodeMirror.
 *
 * This extension adds an "autoCloseTags" option that can be set to
 * either true to get the default behavior, or an object to further
 * configure its behavior.
 *
 * These are supported options:
 *
 * `whenClosing` (default true)
 *   Whether to autoclose when the '/' of a closing tag is typed.
 * `whenOpening` (default true)
 *   Whether to autoclose the tag when the final '>' of an opening
 *   tag is typed.
 * `dontCloseTags` (default is empty tags for HTML, none for XML)
 *   An array of tag names that should not be autoclosed.
 * `indentTags` (default is block tags for HTML, none for XML)
 *   An array of tag names that should, when opened, cause a
 *   blank line to be added inside the tag, and the blank line and
 *   closing line to be indented.
 *
 * See demos/closetag.html for a usage example.
 */(function(e){typeof exports=="object"&&typeof module=="object"?e(require("../../lib/codemirror"),require("../fold/xml-fold")):typeof define=="function"&&define.amd?define(["../../lib/codemirror","../fold/xml-fold"],e):e(CodeMirror)})(function(e){function r(r){if(r.getOption("disableInput"))return e.Pass;var i=r.listSelections(),s=[];for(var a=0;a<i.length;a++){if(!i[a].empty())return e.Pass;var f=i[a].head,l=r.getTokenAt(f),c=e.innerMode(r.getMode(),l.state),h=c.state;if(c.mode.name!="xml"||!h.tagName)return e.Pass;var p=r.getOption("autoCloseTags"),d=c.mode.configuration=="html",v=typeof p=="object"&&p.dontCloseTags||d&&t,m=typeof p=="object"&&p.indentTags||d&&n,g=h.tagName;l.end>f.ch&&(g=g.slice(0,g.length-l.end+f.ch));var y=g.toLowerCase();if(!g||l.type=="string"&&(l.end!=f.ch||!/[\"\']/.test(l.string.charAt(l.string.length-1))||l.string.length==1)||l.type=="tag"&&h.type=="closeTag"||l.string.indexOf("/")==l.string.length-1||v&&o(v,y)>-1||u(r,g,f,h,!0))return e.Pass;var b=m&&o(m,y)>-1;s[a]={indent:b,text:">"+(b?"\n\n":"")+"</"+g+">",newPos:b?e.Pos(f.line+1,0):e.Pos(f.line,f.ch+1)}}for(var a=i.length-1;a>=0;a--){var w=s[a];r.replaceRange(w.text,i[a].head,i[a].anchor,"+insert");var E=r.listSelections().slice(0);E[a]={head:w.newPos,anchor:w.newPos},r.setSelections(E),w.indent&&(r.indentLine(w.newPos.line,null,!0),r.indentLine(w.newPos.line+1,null,!0))}}function i(t,n){var r=t.listSelections(),i=[],s=n?"/":"</";for(var o=0;o<r.length;o++){if(!r[o].empty())return e.Pass;var a=r[o].head,f=t.getTokenAt(a),l=e.innerMode(t.getMode(),f.state),c=l.state;if(!(!n||f.type!="string"&&f.string.charAt(0)=="<"&&f.start==a.ch-1))return e.Pass;var h;if(l.mode.name!="xml")if(t.getMode().name=="htmlmixed"&&l.mode.name=="javascript")h=s+"script";else{if(t.getMode().name!="htmlmixed"||l.mode.name!="css")return e.Pass;h=s+"style"}else{if(!c.context||!c.context.tagName||u(t,c.context.tagName,a,c))return e.Pass;h=s+c.context.tagName}t.getLine(a.line).charAt(f.end)!=">"&&(h+=">"),i[o]=h}t.replaceSelections(i),r=t.listSelections();for(var o=0;o<r.length;o++)(o==r.length-1||r[o].head.line<r[o+1].head.line)&&t.indentLine(r[o].head.line)}function s(t){return t.getOption("disableInput")?e.Pass:i(t,!0)}function o(e,t){if(e.indexOf)return e.indexOf(t);for(var n=0,r=e.length;n<r;++n)if(e[n]==t)return n;return-1}function u(t,n,r,i,s){if(!e.scanForClosingTag)return!1;var o=Math.min(t.lastLine()+1,r.line+500),u=e.scanForClosingTag(t,r,null,o);if(!u||u.tag!=n)return!1;var a=i.context;for(var f=s?1:0;a&&a.tagName==n;a=a.prev)++f;r=u.to;for(var l=1;l<f;l++){var c=e.scanForClosingTag(t,r,null,o);if(!c||c.tag!=n)return!1;r=c.to}return!0}e.defineOption("autoCloseTags",!1,function(t,n,i){i!=e.Init&&i&&t.removeKeyMap("autoCloseTags");if(!n)return;var o={name:"autoCloseTags"};if(typeof n!="object"||n.whenClosing)o["'/'"]=function(e){return s(e)};if(typeof n!="object"||n.whenOpening)o["'>'"]=function(e){return r(e)};t.addKeyMap(o)});var t=["area","base","br","col","command","embed","hr","img","input","keygen","link","meta","param","source","track","wbr"],n=["applet","blockquote","body","button","div","dl","fieldset","form","frameset","h1","h2","h3","h4","h5","h6","head","html","iframe","layer","legend","object","ol","p","select","table","ul"];e.commands.closeTag=function(e){return i(e)}});