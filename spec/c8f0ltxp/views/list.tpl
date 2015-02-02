{{pagelet . "general" "html-header.tpl"}}
<body>
{{pagelet . "general" "nav-header.tpl" "topnav"}}
<div class="container">
  {{range $v := .list}}
  <h3>{{Field $v "title"}}</h3>
  <p>{{Field $v "content"}}</p>
  {{end}}
</div>
{{pagelet . "general" "footer.tpl"}}
</body>
</html>
