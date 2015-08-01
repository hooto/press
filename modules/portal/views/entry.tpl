{{pagelet . "general" "html-header.tpl"}}
<body>
{{pagelet . "general" "nav-header.tpl" "topnav"}}

<div class="container">
  
  <div class="l5s-ctn-header">
    <h2>Content Entry</h2>
  </div>

  <div class="row">
    <div class="col-md-12">
    
      <div class="l5s-nodev">
        
        <div class="header">
          <h2>{{.entry.Title}}</h2>
          <div class="hinfo">
            <span class="section">
              <span class="glyphicon glyphicon-time" aria-hidden="true"></span>&nbsp;
              {{TimeFormat .entry.Created "atom" "Y-m-d H:i:s"}}
            </span>
            
            {{range $term := .entry.Terms}}
              {{if eq $term.Name "categories"}}
              {{if $term.Items}}
              <span class="section">
                <span class="glyphicon glyphicon-folder-open" aria-hidden="true"></span>&nbsp;
                {{range $term_item := $term.Items}}
                <a href="{{$.baseuri}}/list?term_categories={{printf "%d" $term_item.ID}}">{{$term_item.Title}}</a>
                {{end}}
              </span>
              {{end}}
              {{end}}
            {{end}}
            
            {{range $term := .entry.Terms}}
              {{if eq $term.Name "tags"}}
              {{if $term.Items}}
              <span class="section">
                <span class="glyphicon glyphicon-tags" aria-hidden="true"></span>&nbsp;          
                {{range $term_item := $term.Items}}
                <a href="{{$.baseuri}}/list?term_tags={{$term_item.Title}}" class="">{{$term_item.Title}}</a>
                {{end}}
              </span>
              {{end}}
              {{end}}
            {{end}}

          </div>
        </div>

        <div class="content">{{FieldHtml .entry.Fields "content"}}</div>
      </div>

      <div id="{{.modname}}-comments">comments loading</div>
      
    </div>
    
  </div>  

</div>

{{pagelet . "general" "footer.tpl"}}

</body>
</html>
<script type="text/javascript">

window.onload_hooks.push(function() {
    seajs.use([
        "+/comment/~/index.js",
        "+/comment/~/index.css",
    ],
    function() {
        l5sComment.EmbedLoader("{{.modname}}-comments", "{{.modname}}", "{{.__datax_table__}}", "{{.entry.ID}}");
    });
});

window.onload_hooks.push(function() {
    l5s.CodeRender();
});

</script>
