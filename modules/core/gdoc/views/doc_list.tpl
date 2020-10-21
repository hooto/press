<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "v3/html-header.tpl"}}
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/-/static/gdoc/css/main.css"}}?v={{.sys_version_sign}}" type="text/css">
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/~/open-iconic/font/css/open-iconic-bootstrap.css"}}?v={{.sys_version_sign}}" type="text/css">
<body id="hp-body">
{{pagelet . "core/general" "v3/nav-header.tpl" "topnav" "topbar_class=navbar-light"}}

<div class="hp-gdoc-index-frame-dark hp-gdoc-node-content hp-gdoc-bgimg-hexagons">
<div class="container" style="padding: 40px 20px; text-align: center;">
  <div class="hp-gdoc-index-frame-title">Explore Documents</div>
</div>
</div>

<div class="container" style="margin-top:10px">
<div class="hp-gdoc-nodels row" style="">
  {{range $v := .doc_list.Items}}
  <div class="hp-gdoc-nodels-item col-md-6">
  <div class="card mb-4">
    <div class="row">
      <div class="col-auto" style="width:70px; padding:20px 20px;">
        <div class="media-left">
          <figure class="image is-48x48">
		    <img src="/hp/~/octicons/lib/svg/repo.svg" width="48" height="48">
          </figure>
		</div>
      </div>
      <div class="col">
        <div class="card-body">
          <p class="card-title" style="font-size:1.2rem;font-weight:bold;">
		    <a href="{{$.baseuri}}/view/{{$v.ExtPermalinkName}}/">{{FieldStringPrint $v "title" $.LANG}}</a>
		  </p>
		  <p class="subtitle is-6">
            <span>
              <span class="oi oi-timer" title="icon name" aria-hidden="true"></span>
              {{UnixtimeFormat $v.Updated "2006-01-02"}}
            </span>
            <span>
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
            </span>
		  </p>
        </div>
        <div class="card-footer text-right bg-transparent border-0">
           <a class="btn btn-dark" href="{{$.baseuri}}/view/{{$v.ExtPermalinkName}}/">Read</a>
        </div>
      </div>
    </div>
  </div>
  </div>
  {{end}}

    {{if .list_pager}}
    <nav class="pagination is-centered hp-pagination">
      {{if .list_pager.FirstPageNumber}}
      <a class="pagination-previous" href="{{$.baseuri}}/list?page={{.list_pager.FirstPageNumber}}">First</a>
      {{end}}
      <ul class="pagination-list">
      {{range $index, $page := .list_pager.RangePages}}
      <li>
        <a class="pagination-link {{if eq $page $.list_pager.CurrentPageNumber}}is-current{{end}}" href="{{$.baseuri}}/list?{{FilterUri $ "page" $page}}">{{$page}}</a>
      </li>
      {{end}}
      </ul>

      {{if .list_pager.LastPageNumber}}
      <a class="pagination-next" href="{{$.baseuri}}/list?page={{.list_pager.LastPageNumber}}">Last</a>
      {{end}}
    </nav>
    {{end}}

</div>
</div>


{{pagelet . "core/general" "v3/footer.tpl"}}

{{pagelet . "core/general" "html-footer.tpl"}}
</body>
</html>

