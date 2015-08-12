<style>
.page-header {
	margin: 10px 0;
	font-height: 100%;
}
.l5s-sys-table {
  font-size: 14px;
}
.l5s-sys-table td {
  padding: 5px !important;
}
.l5s-sys-table tr.line {
  border-top: 1px solid #ccc;
}
</style>

<div class="panel panel-default">
  <div class="panel-heading">Identity Service Status</div>
  <div class="panel-body">

    <table width="100%" style="margin-bottom:20px;">
      <tr>
        <td width="20%">Identity Service Url</td>
        <td><input type="text" name="ids_url" class="form-control" placeholder="Enter the Identity Service URL" value="{[=it.service_url]}" readonly></td>
      </tr>
    </table>

    <div id="l5s-mgr-sys-ids-alert"></div>

    <form id="l5s-mgr-sys-ids" action="#">
    <table width="100%" class="table table-hover l5s-sys-table">
      <thead><tr>
        <th width="20%"></th>
        <th width="40%">Local App Info</th>
        <th>Registered in Identity Service</th>
      </tr></thead>
      
      <tr>
        <td>App Version</td>
        <td>{[=it.instance_self.version]}</td>
        <td>{[=it.instance_registered.version]}</td>
      </tr>

      <tr>
        <td>App ID</td>
        <td>{[=it.instance_self.app_id]}</td>
        <td>{[=it.instance_registered.app_id]}</td>
      </tr>

      <tr>
        <td>App Name</td>
        <td>
          <input type="text" name="app_title" class="form-control input-sm" 
            placeholder="Enter the App Name" value="{[=it.instance_self.app_title]}">
        </td>
        <td>{[=it.instance_registered.app_title]}</td>
      </tr>

      <tr>
        <td>Entry URL</td>
        <td>
          <input type="text" name="instance_url" class="form-control input-sm" 
            placeholder="Enter the Entry URL of App Instance" value="{[=it.instance_self.url]}">
        </td>
        <td>{[=it.instance_registered.url]}</td>
      </tr>

      <tr>
        <td>Privileges</td>
        <td>
          <table class="table">
          <thead>
            <tr>
              <th>Privilege</th>
              <th>Roles</th>
            </tr>
          </thead>
          <tbody>
            {[~it.instance_self.privileges :v]}
            <tr>
              <td>
                <p><strong>{[=v.desc]}</strong></p>
                <p>{[=v.privilege]}</p>
              </td>
              <td>
              {[ if (v.roles) { ]}
              {[~v.roles :rv]}
                {[~it._roles.items :drv]}
                {[ if (rv == drv.idxid) { ]}
                  <p>{[=drv.meta.name]}</p>
                {[ } ]}
                {[~]}
              {[~]}
              {[ } ]}
              </td>
            </tr>
            {[~]}
          </tbody>
          </table>
        </td>
        <td>
          <table class="table">
          <thead>
            <tr>
              <th>Privilege</th>
              <th>Roles</th>
            </tr>
          </thead>
          <tbody>
            {[~it.instance_registered.privileges :v]}
            <tr>
              <td>
                <p><strong>{[=v.desc]}</strong></p>
                <p>{[=v.privilege]}</p>
              </td>
              <td>
              {[ if (v.roles) { ]}
              {[~v.roles :rv]}
                {[~it._roles.items :drv]}
                {[ if (rv == drv.idxid) { ]}
                  <p>{[=drv.meta.name]}</p>
                {[ } ]}
                {[~]}
              {[~]}
              {[ } ]}
              </td>
            </tr>
            {[~]}
          </tbody>
          </table>
        </td>
      </tr>
      
    </table>
    </form>

    <div class="text-center">
      <button type="submit" class="btn btn-success" onclick="l5sSys.IdentitySync()">Sync to Identity Service</button>
    </div>

  </div>
</div>
