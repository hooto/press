<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "html-header.tpl"}}
<body>
{{pagelet . "core/general" "nav-header.tpl" "topnav"}}

<div class="container">
  
  <div class="hp-ctn-header">
    <h2>Content Explore</h2>
  </div>

  <div class="row">

    <div class="col-md-9">
    
    <ul class="hp-nodels">
      {{range $v := .list.Items}}
      <li class="hp-nodels-item">
        <h4 class="hp-nodels-heading"><a href="{{$.baseuri}}/view/{{$v.SelfLink}}">{{FieldStringPrint $v "title" $.LANG}}</a></h4>
        <span class="hp-nodels-info">
            
            <span class="section">
              <span class="glyphicon glyphicon-time" aria-hidden="true"></span>&nbsp;
              {{UnixtimeFormat $v.Created "Y-m-d"}}
            </span>
            
            {{range $term := $v.Terms}}
              {{if eq $term.Name "categories"}}
              {{if $term.Items}}
              <span class="section">
                <span class="glyphicon glyphicon-th-list" aria-hidden="true"></span>&nbsp;
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
              <span class="section">
                <span class="glyphicon glyphicon-tags" aria-hidden="true"></span>&nbsp;          
                {{range $term_item := $term.Items}}
                <a href="{{$.baseuri}}/list?term_tags={{$term_item.Title}}" class="tag-item">{{$term_item.Title}}</a>
                {{end}}
              </span>
              {{end}}
              {{end}}
            {{end}}

        </span>

        <div class="hp-nodels-text">{{FieldHtmlSubPrint $v "content" 200 $.LANG}}</div>
      </li>
      {{end}}
    </ul>

    {{if .list_pager}}
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
    {{end}}

    </div>

    <div class="col-md-3">
        
        {{pagelet . .modname "search.tpl"}}

        {{pagelet . .modname "term/categories.tpl"}}

    </div> 
    
  </div>
</div>

{{pagelet . "core/general" "footer.tpl"}}

{{pagelet . "core/general" "html-footer.tpl"}}
</body>
</html>
