<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "v2/html-header.tpl"}}
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/-/static/gdoc/css/main.css"}}?v={{.sys_version_sign}}" type="text/css">
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/~/open-iconic/font/css/open-iconic-bootstrap.css"}}?v={{.sys_version_sign}}" type="text/css">
<body id="hp-body">
{{pagelet . "core/general" "v2/nav-header.tpl" "topnav"}}

<div class="hp-gdoc-index-frame-dark hp-gdoc-node-content hp-gdoc-bgimg-hexagons">
<div class="container" style="padding: 20px 10px; text-align: center;">
  <div class="hp-gdoc-index-frame-title">Explore Documents</div>
</div>
</div>

<div class="container" style="margin-top:10px">
<div class="hp-gdoc-nodels columns is-multiline" style="">
  {{range $v := .doc_list.Items}}
  <div class="hp-gdoc-nodels-item column is-6">
  <div class="card">
    <div class="card-content">
      <div class="media">
        <div class="media-left">
          <figure class="image is-48x48">
		    <img src="/hp/~/octicons/lib/svg/repo.svg" width="48" height="48">
          </figure>
		</div>
        <div class="media-content">
          <p class="title is-4">
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
		  <p style="text-align:right">
            <a class="button is-dark" href="{{$.baseuri}}/view/{{$v.ExtPermalinkName}}/">Read</a>
		  </p>
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


{{pagelet . "core/general" "v2/footer.tpl"}}

{{pagelet . "core/general" "html-footer.tpl"}}
</body>
</html>

