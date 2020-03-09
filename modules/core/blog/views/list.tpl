<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "v2/html-header.tpl"}}
<body id="hp-body">

{{pagelet . "core/general" "v2/nav-header.tpl" "topnav"}}

<div class="container" style="margin-top:10px">

  <div class="columns">
    <div class="column is-9">
      <div class="hp-ctn-title">
        Content Explore
      </div>
    </div>
    <div class="column is-3">
      <form action="{{.baseuri}}/list">
        <div class="field has-addons">
          <div class="control">
            <input type="text" class="input"
              placeholder=""
              name="qry_text"
              value="{{.qry_text}}">
          </div>
          <div class="control">
            <input class="button is-dark" type="submit" value="Search">
          </div>
        </div>
      </form>
    </div>
  </div>
</div>

<div class="container">
  <div class="columns">

    <div class="column is-9">

    <ul class="hp-node-list">
      {{range $v := .list.Items}}
      <li class="hp-node-list-item">
        <h4 class="hp-node-list-heading">
          <a href="{{$.baseuri}}/view/{{$v.ID}}.html">{{FieldStringPrint $v "title" $.LANG}}</a>
        </h4>
        <div class="hp-node-list-info">

            <span class="info-item">
              Published : {{UnixtimeFormat $v.Created "Y-m-d"}}
            </span>

            {{range $term := $v.Terms}}
              {{if eq $term.Name "categories"}}
              {{if $term.Items}}
              <span class="info-item">
                Categories :
                {{range $term_item := $term.Items}}
                <a href="{{$.baseuri}}/list?term_categories={{printf "%d" $term_item.ID}}">{{$term_item.Title}}</a>
                {{end}}
              </span>
              {{end}}
              {{end}}
            {{end}}

            {{range $term := $v.Terms}}
              {{if eq $term.Name "tags"}}
              {{if $term.Items}}
              <span class="info-item">
                Tags :
                {{range $term_item := $term.Items}}
                <a href="{{$.baseuri}}/list?term_tags={{$term_item.Title}}" class="info-tag-item">{{$term_item.Title}}</a>
                {{end}}
              </span>
              {{end}}
              {{end}}
            {{end}}

        </div>

        <div class="hp-node-list-text">{{FieldHtmlSubPrint $v "content" 200 $.LANG}}</div>
      </li>
      {{end}}
    </ul>

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

    <div class="column is-3">
        {{pagelet . .modname "term/categories.tpl"}}
    </div>

  </div>
</div>

{{pagelet . "core/general" "v2/footer.tpl"}}

{{pagelet . "core/general" "html-footer.tpl"}}
</body>
</html>
