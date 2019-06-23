<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "bs4/html-header.tpl"}}
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/-/static/gdoc/css/main.css"}}?v={{.sys_version_sign}}" type="text/css">
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/~/open-iconic/font/css/open-iconic-bootstrap.css"}}?v={{.sys_version_sign}}" type="text/css">
<body>
{{pagelet . "core/general" "bs4/nav-header.tpl" "topnav"}}

<div class="hp-gdoc-index-frame-blue hp-gdoc-node-content hp-gdoc-bgimg-hexagons">
<div class="container" style="padding: 20px 10px; text-align: center;">
  <div class="hp-gdoc-index-frame-title">Explore Documents</div>
</div>
</div>

<div class="container">
<div class="hp-gdoc-nodels row" style="padding: 20px 0;">
  {{range $v := .doc_list.Items}}
  <div class="hp-gdoc-nodels-item col-sm-6">
  <div class="card">
    <div class="card-body">
      <h5 class="card-title">
        <a href="{{$.baseuri}}/view/{{$v.ExtPermalinkName}}/">{{FieldStringPrint $v "title" $.LANG}}</a>
      </h5>
      <div class="_hp-gdoc-nodels-info card-text">
      <span>
        <span class="oi oi-timer" title="icon name" aria-hidden="true"></span>
        {{UnixtimeFormat $v.Updated "2006-01-02"}}
      </span>
      {{range $term := $v.Terms}}
      {{if eq $term.Name "tags"}}
        {{if $term.Items}}
          <span>
            <span class="oi oi-tags" title="icon name" aria-hidden="true"></span>
            {{range $term_item := $term.Items}}
            <a href="{{$.baseuri}}/list?term_tags={{$term_item.Title}}" class="tag-item">{{$term_item.Title}}</a>
            {{end}}
          </span>
        {{end}}
      {{end}}
      {{end}}
      </div>
    </div>
    <div class="card-footer bg-transparent">
      <a class="btn btn-primary" href="{{$.baseuri}}/view/{{$v.ExtPermalinkName}}/">Read</a>
    </div>
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

