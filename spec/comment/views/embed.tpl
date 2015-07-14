<div class="l5s-comment-embed">
  
  {{if len .list.Items}}
  <header>
    <nav class="nav-primary">
      <ul>
        <li>
          <span>Comments</span>
        </li>
      </ul>
    </nav>
  </header>
  {{end}}

  <div id="l5s-comment-embed-list" class="list">
    {{range $v := .list.Items}}
    <div class="entry">
      <div class="avatar">
        <img src="/+/comment/~/img/user-default.png">
      </div>  
      
      <div class="body">
        <div class="info">
          <strong>{{FieldSubString $v.Fields "author" 50}}</strong>
          <small>@{{TimeFormat $v.Created "atom" "Y-m-d H:i"}}</small>
        </div>
        <p>{{FieldSubHtml $v.Fields "content" 300}}</p>
      </div>
    </div>
    {{end}}
  </div>

  <header>
    <nav class="nav-primary">
      <ul>
        <li>
          <span>New Comment</span>
        </li>
      </ul>
    </nav>
  </header>

  <div id="l5s-comment-embed-new-form" class="new">

    <div id="l5s-comment-embed-new-form-alert"></div>

    <input type="hidden" name="refer_id" value="{{.new_form_refer_id}}">
    <input type="hidden" name="refer_specid" value="{{.new_form_refer_specid}}">
    <input type="hidden" name="refer_datax_table" value="{{.new_form_refer_datax_table}}">

    <div class="form-group">
      <label>Your name</label>
      <input type="text" class="form-control" name="author" value="{{.new_form_author}}">
    </div>

    <div class="form-group">
      <label>Content</label>
      <textarea class="form-control" rows="3" name="content"></textarea>
    </div>

    <button class="btn btn-default" onclick="l5sComment.EmbedCommit()">Commit</button>

  </div>
</div>

<script id="l5s-comment-embed-tpl" type="text/html">
<div class="entry" id="entry-{[=it.meta.id]}">
  <div class="avatar">
    <img src="/+/comment/~/img/user-default.png">
  </div>
  <div class="body">
    <div class="info">
      <strong>{[=it.author]}</strong>
      <small>@{[=it.meta.created]}</small>
    </div>
    <p>{[=it.content]}</p>
  </div>
</div>
</script>

