<style>
.page-header {
	margin: 10px 0;
	font-height: 100%;
}
.htp-sys-table {
  font-size: 12px;
}
.htp-sys-table td {
  padding: 5px !important;
}
.htp-sys-table tr.line {
  border-top: 1px solid #ccc;
}
</style>

<!-- <div class="page-header">
  <h2>System Monitor Status <small></small></h2>
</div> -->

<div class="panel panel-default">
  <div class="panel-heading">System Monitor Status</div>
  <div class="panel-body">

    <table width="100%" class="htp-sys-table">
      
      <tr>
        <td width="30%">App Instance ID</td>
        <td>{[=it.instance_id]}</td>
      </tr>
      <tr>
        <td>App Version</td>
        <td>{[=it.app_version]}</td>
      </tr>
      <tr>
        <td>Runtime Version</td>
        <td>{[=it.runtime_version]}</td>
      </tr>
      <tr>
        <td>Uptime</td>
        <td>{[=l4i.TimeParseFormat(it.uptime, "Y-m-d H:i:s")]}</td>
      </tr>

      <tr class="line">
        <td>Current Coroutine Number</td>
        <td>{[=it.coroutine_number]}</td>
      </tr>
      <tr>
        <td>Current Memory Allocated</td>
        <td>{[=htpSys.UtilResourceSizeFormat(it.memstats.alloc)]}</td>
      </tr>
      <tr>
        <td>Total Memory Allocated</td>
        <td>{[=htpSys.UtilResourceSizeFormat(it.memstats.total_alloc)]}</td>
      </tr>
      <tr>
        <td>Memory obtained from system</td>
        <td>{[=htpSys.UtilResourceSizeFormat(it.memstats.sys)]}</td>
      </tr>

      <tr class="line">
        <td>Next GC Recycle</td>
        <td>{[=htpSys.UtilResourceSizeFormat(it.memstats.next_gc)]}</td>
      </tr>
      <tr>
        <td>Since Last GC Time</td>
        <td>{[=htpSys.UtilDurationFormat((new Date().getTime()) - (it.memstats.last_gc / 1000000))]}</td>
      </tr>
      <tr>
        <td>Total GC Pause</td>
        <td>{[=htpSys.UtilDurationFormat(it.memstats.pause_total_ns, 1000000)]}</td>
      </tr>
      <tr>
        <td>Total GC Times</td>
        <td>{[=it.memstats.num_gc]}</td>
      </tr>
      <tr>
        <td>Average GC Pause</td>
        <td>{[=htpSys.UtilDurationFormat((it.memstats.pause_total_ns / it.memstats.num_gc), 1000000)]}</td>
      </tr>


      <!-- <tr class="line">
        <td>CpuNum</td>
        <td>{[=it.info.cpu_num]}</td>
      </tr>
      <tr>
        <td>Uptime</td>
        <td>{[=htpSys.UtilDurationFormat(it.info.uptime * 1000)]}</td>
      </tr>
      <tr>
        <td>Loads</td>
        <td>{[=it.info.loads[0]]}</td>
      </tr>
      <tr>
        <td>MemTotal</td>
        <td>{[=htpSys.UtilResourceSizeFormat(it.info.mem_total)]}</td>
      </tr>
      <tr>
        <td>MemFree</td>
        <td>{[=htpSys.UtilResourceSizeFormat(it.info.mem_free)]}</td>
      </tr>
      <tr>
        <td>MemShared</td>
        <td>{[=htpSys.UtilResourceSizeFormat(it.info.mem_shared)]}</td>
      </tr>
      <tr>
        <td>MemBuffer</td>
        <td>{[=htpSys.UtilResourceSizeFormat(it.info.mem_buffer)]}</td>
      </tr>
      <tr>
        <td>MemUsed</td>
        <td>{[=htpSys.UtilResourceSizeFormat(it.info.mem_used)]}</td>
      </tr>
      <tr>
        <td>SwapTotal</td>
        <td>{[=htpSys.UtilResourceSizeFormat(it.info.swap_total)]}</td>
      </tr>
      <tr>
        <td>SwapFree</td>
        <td>{[=htpSys.UtilResourceSizeFormat(it.info.swap_free)]}</td>
      </tr>
      <tr>
        <td>Procs</td>
        <td>{[=it.info.procs]}</td>
      </tr>   -->    

      
    </table>

  </div>
</div>
