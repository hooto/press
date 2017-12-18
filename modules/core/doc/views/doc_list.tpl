<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "bs4/html-header.tpl"}}
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/-/static/doc/css/main.css"}}?v={{.sys_version_sign}}" type="text/css">
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/~/open-iconic/font/css/open-iconic-bootstrap.css"}}?v={{.sys_version_sign}}" type="text/css">
<body>
{{pagelet . "core/general" "bs4/nav-header.tpl" "topnav"}}

<div class="hpdoc_index_frame_blue hpdoc_node_content hpdoc_bgimg_hexagons">
<div class="container" style="padding: 20px 10px; text-align: center;">
  <div class="hpdoc_index_frame_title">Explore documents</div>
</div>
</div>

<div class="container">
<div class="hpdoc-nodels" style="padding: 20px 0;">
  {{range $v := .doc_list.Items}}
  <div class="hpdoc-nodels-item row">
  <div class="col-sm-10">
    <div class="hpdoc-nodels-title"><a href="{{$.baseuri}}/entry/{{$v.ExtPermalinkName}}/">{{$v.Title}}</a></div>
    <div class="hpdoc-nodels-info">
    {{range $term := $v.Terms}}
      {{if eq $term.Name "tags"}}
        {{if $term.Items}}
          <span>
            <img src="/hp/~/open-iconic/svg/tags.svg" width="16" height="16" class="hpdoc_icon">&nbsp;          
            {{range $term_item := $term.Items}}
            <a href="{{$.baseuri}}/list?term_tags={{$term_item.Title}}" class="tag-item">{{$term_item.Title}}</a>
            {{end}}
          </span>
        {{end}}
      {{end}}
    {{end}}
    </div>
  </div>
  <div class="col-sm-2" style="text-align:right;padding-right:0">
    <a class="btn btn-success" href="{{$.baseuri}}/entry/{{$v.ExtPermalinkName}}/">Read</a>
  </div>
  </div>
  {{end}}

  {{if .list_pager}}
  <div>
    <ul class="pagination pagination-sm">
      {{if .list_pager.FirstPageNumber}}
      <li>
        <a href="{{$.baseuri}}/list?page={{.list_pager.FirstPageNumber}}">First</a>
      </li>
      {{end}}

      {{range $index, $page := .list_pager.RangePages}}
      <li {{if eq $page $.list_pager.CurrentPageNumber}}class="active"{{end}}>
        <a href="{{$.baseuri}}/list?{{FilterUri $ "page" $page}}">{{$page}}</a>
      </li>
      {{end}}
      
      {{if .list_pager.LastPageNumber}}
      <li>
        <a href="{{$.baseuri}}/list?page={{.list_pager.LastPageNumber}}">Last</a>
      </li>
      {{end}}
    </ul>
  </div>
  {{end}}
</div>
</div>


{{pagelet . "core/general" "footer.tpl"}}

{{pagelet . "core/general" "html-footer.tpl"}}
</body>
</html>

