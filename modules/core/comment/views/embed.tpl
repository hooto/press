<div class="htp-comment-embed">

  <header id="htp-comment-embed-list-header" style="display:{{if len .list.Items}}block{{else}}none{{end}}">
    <nav class="nav-primary">
      <ul>
        <li>
          <span>Comments</span>
        </li>
      </ul>
    </nav>
  </header>

  <div id="htp-comment-embed-list" class="list">
    {{range $v := .list.Items}}
    <div class="entry">
      <div class="avatar">
        <img src="{{HttpSrvBasePath "htp/+/comment/~/img/user-default.png"}}">
      </div>  
      
      <div class="body">
        <div class="info">
          <strong>{{FieldSubString $v.Fields "author" 50}}</strong>
          <small>@{{TimeFormat $v.Created "atom" "Y-m-d H:i"}}</small>
        </div>
        <p>{{FieldSubHtml $v.Fields "content" 2000}}</p>
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

  <div class="list">
    <div class="entry">
      <div class="avatar">
        <img src="{{HttpSrvBasePath "htp/+/comment/~/img/user-default.png"}}">
      </div>  
      
      <div id="htp-comment-embed-new-form-ctrl" class="body">
        <div>
          <div class="info"><strong>Guest</strong></div>
          <div>
            <input type="text" class="form-control" name="author" placeholder="Leave a comment ..." onclick="htpComment.EmbedFormActive()">
          </div>
        </div>
      </div>

      <div id="htp-comment-embed-new-form" class="body new" style="display:none;">

        <div id="htp-comment-embed-new-form-alert"></div>

        <input type="hidden" name="refer_id" value="{{.new_form_refer_id}}">
        <input type="hidden" name="refer_modname" value="{{.new_form_refer_modname}}">
        <input type="hidden" name="refer_datax_table" value="{{.new_form_refer_datax_table}}">
        <input type="hidden" name="captcha_token" value="">

        <div class="form-group">
          <label>Your name</label>
          <input type="text" class="form-control" name="author" value="{{.new_form_author}}">
        </div>

        <div class="form-group">
          <label>Content</label>
          <textarea class="form-control" rows="3" name="content"></textarea>
        </div>

        <div class="form-group">
          <label>Verification</label>
          <div>
            <div class="row">          
              <div class="col-xs-6">
                <input type="text" class="form-control" name="captcha_word" value="">
                <span class="help-block">Type the characters you see in the right picture</span>
              </div>
              <div class="col-xs-6" style="background-color: #dce6ff;">
                <img id="htp-comment-captcha-url" src="">
              </div>
            </div>        
          </div>
        </div>

        <button class="btn btn-default" onclick="htpComment.EmbedCommit()">Commit</button>

      </div>
    </div>
  </div>
</div>

<script id="htp-comment-embed-tpl" type="text/html">
<div class="entry" id="entry-{[=it.meta.id]}">
  <div class="avatar">
    <img src="{[=htp.HttpSrvBasePath("+/comment/~/img/user-default.png")]}">
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

