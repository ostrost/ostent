define(function(require) {
	var React = require('react');
	return {
		mem_rows: function(Data, $mem) { return (<tr  key={"mem-rowby-kind-"+$mem.Kind}
  ><td
    >{$mem.Kind}</td
  ><td className="text-right"
    >{$mem.Free}</td
  ><td className="text-right bg-usepct"data-usepct={$mem.UsePct}
    >{$mem.UsePct}%</td
  ><td className="text-right"
    >{$mem.Used}</td
  ><td className="text-right"
    >{$mem.Total}</td
  ></tr
>); },
		panelmem: function(Data, rows) { return (<div  className={!Data.Params.Memn.Negative ? "" : "panel panel-default"}
  ><div className="h4 padding-left-like-panel-heading"
    ><a  href={Data.Params.Tlinks.Memn} onClick={this.handleClick}
      >Memory</a
    ></div
  ><ul   className={!Data.Params.Memn.Negative ? "collapse-hidden" : "list-group"}
    ><li className="list-group-item text-nowrap th"
      ><ul className="list-inline"
        ><li
          ><span
            ><b
              >Delay</b
            > <span className="badge"
              >{Data.Params.Memd}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Dlinks.Memd.Less.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Memd.Less.ExtraClass != null ? Data.Params.Dlinks.Memd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Dlinks.Memd.Less.Text}</a
><a href={Data.Params.Dlinks.Memd.More.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Memd.More.ExtraClass != null ? Data.Params.Dlinks.Memd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Dlinks.Memd.More.Text} +</a
></div
          ></li
        ><li
          ><span
            ><b
              >Rows</b
            > <span className="badge"
              >{Data.Params.Memn.Absolute}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Nlinks.Memn.Less.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Memn.Less.ExtraClass != null ? Data.Params.Nlinks.Memn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Nlinks.Memn.Less.Text}</a
><a href={Data.Params.Nlinks.Memn.More.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Memn.More.ExtraClass != null ? Data.Params.Nlinks.Memn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Nlinks.Memn.More.Text} +</a
></div
          ></li
        ></ul
      ></li
    ></ul
  ><table  className={Data.Params.Memn.Absolute != 0 ? "table table-hover" : "collapse-hidden"}
    ><thead
      ><tr
        ><th
          ></th
        ><th className="text-right"
          >Free</th
        ><th className="text-right"
          >Use%</th
        ><th className="text-right"
          >Used</th
        ><th className="text-right"
          >Total</th
        ></tr
      ></thead
    ><tbody
      >{rows}</tbody
    ></table
  ></div
>); },

		if_rows:  function(Data, $if)  { return (<tr  key={"if-rowby-name-"+$if.Name}
  ><td className="text-nowrap clip12" title={$if.Name}
    >{$if.Name}</td
  ><td className="text-right"
    >{$if.IP}</td
  ><td className="text-right text-nowrap"
    ><span className="mutext" title="Total BYTES modulo 4G"
      >{$if.BytesIn}/{$if.BytesOut}</span
    > <span title="BITS per second"
      >{$if.DeltaBitsIn}/{$if.DeltaBitsOut}</span
    ></td
  ><td className="text-right text-nowrap"
    ><span className="mutext" title="Total packets modulo 4G"
      >{$if.PacketsIn}/{$if.PacketsOut}</span
    > {$if.DeltaPacketsIn}/{$if.DeltaPacketsOut}</td
  ><td className="text-right text-nowrap"
    ><span className="mutext" title="Total drops modulo 4G"
      >{$if.DropsIn}/{$if.DropsOut}</span
    > {$if.DeltaDropsIn}/{$if.DeltaDropsOut}</td
  ><td className="text-right text-nowrap"
    ><span className="mutext" title="Total errors modulo 4G"
      >{$if.ErrorsIn}/{$if.ErrorsOut}</span
    > {$if.DeltaErrorsIn}/{$if.DeltaErrorsOut}</td
  ></tr
>); },
		panelif:  function(Data, rows) { return (<div  className={!Data.Params.Ifn.Negative ? "" : "panel panel-default"}
  ><div className="h4 padding-left-like-panel-heading"
    ><a  href={Data.Params.Tlinks.Ifn} onClick={this.handleClick}
      >Interfaces</a
    ></div
  ><ul   className={!Data.Params.Ifn.Negative ? "collapse-hidden" : "list-group"}
    ><li className="list-group-item text-nowrap th"
      ><ul className="list-inline"
        ><li
          ><span
            ><b
              >Delay</b
            > <span className="badge"
              >{Data.Params.Ifd}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Dlinks.Ifd.Less.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Ifd.Less.ExtraClass != null ? Data.Params.Dlinks.Ifd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Dlinks.Ifd.Less.Text}</a
><a href={Data.Params.Dlinks.Ifd.More.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Ifd.More.ExtraClass != null ? Data.Params.Dlinks.Ifd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Dlinks.Ifd.More.Text} +</a
></div
          ></li
        ><li
          ><span
            ><b
              >Rows</b
            > <span className="badge"
              >{Data.Params.Ifn.Absolute}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Nlinks.Ifn.Less.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Ifn.Less.ExtraClass != null ? Data.Params.Nlinks.Ifn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Nlinks.Ifn.Less.Text}</a
><a href={Data.Params.Nlinks.Ifn.More.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Ifn.More.ExtraClass != null ? Data.Params.Nlinks.Ifn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Nlinks.Ifn.More.Text} +</a
></div
          ></li
        ></ul
      ></li
    ></ul
  ><table  className={Data.Params.Ifn.Absolute != 0 ? "table table-hover" : "collapse-hidden"}
    ><thead
      ><tr
        ><th
          >Interface</th
        ><th className="text-right"
          >IP</th
        ><th className="text-right text-nowrap col-md-2" title="Bits In/Out per second"
          >IO <i
            >b</i
          >ps</th
        ><th className="text-right text-nowrap col-md-2" title="Packets In/Out per second"
          >Packets IO</th
        ><th className="text-right text-nowrap col-md-2" title="Drops In/Out per second"
          >Drops IO</th
        ><th className="text-right text-nowrap col-md-2" title="Errors In/Out per second"
          >Errors IO</th
        ></tr
      ></thead
    ><tbody
      >{rows}</tbody
    ></table
  ></div
>); },

		cpu_rows: function(Data, $cpu) { return (<tr  key={"cpu-rowby-N-"+$cpu.N}
  ><td className="text-right text-nowrap"
    >{$cpu.N}</td
  ><td className="text-right bg-usepct"
  data-usepct={$cpu.UserPct}
    >{$cpu.UserPct}%</td
  ><td className="text-right bg-usepct"
  data-usepct={$cpu.SysPct}
    >{$cpu.SysPct}%</td
  ><td className="text-right bg-usepct"
  data-usepct={$cpu.WaitPct}
    >{$cpu.WaitPct}%</td
  ><td className="text-right bg-usepct-inverse"
  data-usepct={$cpu.IdlePct}
    >{$cpu.IdlePct}%</td
  ></tr
>); },
		panelcpu: function(Data, rows) { return (<div  className={!Data.Params.CPUn.Negative ? "" : "panel panel-default"}
  ><div className="h4 padding-left-like-panel-heading"
    ><a  href={Data.Params.Tlinks.CPUn} onClick={this.handleClick}
      >CPU</a
    ></div
  ><ul   className={!Data.Params.CPUn.Negative ? "collapse-hidden" : "list-group"}
    ><li className="list-group-item text-nowrap th"
      ><ul className="list-inline"
        ><li
          ><span
            ><b
              >Delay</b
            > <span className="badge"
              >{Data.Params.CPUd}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Dlinks.CPUd.Less.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.CPUd.Less.ExtraClass != null ? Data.Params.Dlinks.CPUd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Dlinks.CPUd.Less.Text}</a
><a href={Data.Params.Dlinks.CPUd.More.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.CPUd.More.ExtraClass != null ? Data.Params.Dlinks.CPUd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Dlinks.CPUd.More.Text} +</a
></div
          ></li
        ><li
          ><span
            ><b
              >Rows</b
            > <span className="badge"
              >{Data.Params.CPUn.Absolute}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Nlinks.CPUn.Less.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.CPUn.Less.ExtraClass != null ? Data.Params.Nlinks.CPUn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Nlinks.CPUn.Less.Text}</a
><a href={Data.Params.Nlinks.CPUn.More.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.CPUn.More.ExtraClass != null ? Data.Params.Nlinks.CPUn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Nlinks.CPUn.More.Text} +</a
></div
          ></li
        ></ul
      ></li
    ></ul
  ><table  className={Data.Params.CPUn.Absolute != 0 ? "table table-hover" : "collapse-hidden"}
    ><thead
      ><tr
        ><th
          ></th
        ><th className="text-right"
          >User</th
        ><th className="text-right"
          >Sys</th
        ><th className="text-right"
          >Wait</th
        ><th className="text-right"
          >Idle</th
        ></tr
      ></thead
    ><tbody
      >{rows}</tbody
    ></table
  ></div
>); },

		df_rows:  function(Data, $df)  { return (<tr  key={"df-rowby-dirname-"+$df.DirName}
  >  <td className="text-nowrap clip12" title={$df.DevName}
    >{$df.DevName}</td
  >  <td className="text-nowrap clip12" title={$df.DirName}
    >{$df.DirName}</td
  ><td className="text-right text-nowrap"
    ><span className="mutext" title="Inodes free"
      >{$df.Ifree}</span
    > {$df.Avail}</td
  ><td className="text-right bg-usepct text-nowrap"data-usepct={$df.UsePct}
    ><span className="mutext" title="Inodes use%"
      >{$df.IusePct}%</span
    > {$df.UsePct}%</td
  ><td className="text-right text-nowrap"
    ><span className="mutext" title="Inodes used"
      >{$df.Iused}</span
    > {$df.Used}</td
  ><td className="text-right text-nowrap"
    ><span className="mutext" title="Inodes total"
      >{$df.Inodes}</span
    > {$df.Total}</td
  ></tr
>); },
		paneldf:  function(Data,rows)  { return (<div  className={!Data.Params.Dfn.Negative ? "" : "panel panel-default"}
  ><div className="h4 padding-left-like-panel-heading"
    ><a  href={Data.Params.Tlinks.Dfn} onClick={this.handleClick}
      >Disk usage</a
    ></div
  ><ul   className={!Data.Params.Dfn.Negative ? "collapse-hidden" : "list-group"}
    ><li className="list-group-item text-nowrap th"
      ><ul className="list-inline"
        ><li
          ><span
            ><b
              >Delay</b
            > <span className="badge"
              >{Data.Params.Dfd}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Dlinks.Dfd.Less.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Dfd.Less.ExtraClass != null ? Data.Params.Dlinks.Dfd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Dlinks.Dfd.Less.Text}</a
><a href={Data.Params.Dlinks.Dfd.More.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Dfd.More.ExtraClass != null ? Data.Params.Dlinks.Dfd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Dlinks.Dfd.More.Text} +</a
></div
          ></li
        ><li
          ><span
            ><b
              >Rows</b
            > <span className="badge"
              >{Data.Params.Dfn.Absolute}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Nlinks.Dfn.Less.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Dfn.Less.ExtraClass != null ? Data.Params.Nlinks.Dfn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Nlinks.Dfn.Less.Text}</a
><a href={Data.Params.Nlinks.Dfn.More.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Dfn.More.ExtraClass != null ? Data.Params.Nlinks.Dfn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Nlinks.Dfn.More.Text} +</a
></div
          ></li
        ></ul
      ></li
    ></ul
  ><table  className={Data.Params.Dfn.Absolute != 0 ? "table table-hover" : "collapse-hidden"}
    ><thead
      ><tr className="text-nowrap"
        ><th className="header "
  ><a href={Data.Params.Vlinks.Dfk[1-1].LinkHref} className={Data.Params.Vlinks.Dfk[1-1].LinkClass} onClick={this.handleClick}  
    >Device<span className={Data.Params.Vlinks.Dfk[1-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Params.Vlinks.Dfk[2-1].LinkHref} className={Data.Params.Vlinks.Dfk[2-1].LinkClass} onClick={this.handleClick}  
    >Mounted<span className={Data.Params.Vlinks.Dfk[2-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[3-1].LinkHref} className={Data.Params.Vlinks.Dfk[3-1].LinkClass} onClick={this.handleClick}  
    >Avail<span className={Data.Params.Vlinks.Dfk[3-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[4-1].LinkHref} className={Data.Params.Vlinks.Dfk[4-1].LinkClass} onClick={this.handleClick}  
    >Use%<span className={Data.Params.Vlinks.Dfk[4-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[5-1].LinkHref} className={Data.Params.Vlinks.Dfk[5-1].LinkClass} onClick={this.handleClick}  
    >Used<span className={Data.Params.Vlinks.Dfk[5-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[6-1].LinkHref} className={Data.Params.Vlinks.Dfk[6-1].LinkClass} onClick={this.handleClick}  
    >Total<span className={Data.Params.Vlinks.Dfk[6-1].CaretClass}
      ></span
    ></a
  ></th
></tr
      ></thead
    ><tbody
      >{rows}</tbody
    ></table
  ></div
>); },

		ps_rows:  function(Data, $ps)  { return (<tr  key={"ps-rowby-pid-"+$ps.PID}
  ><td className="text-right"
    > {$ps.PID}</td
  ><td className="text-right"
    > {$ps.UID}</td
  ><td
    >{$ps.User}</td
  ><td className="text-right"
    > {$ps.Priority}</td
  ><td className="text-right"
    > {$ps.Nice}</td
  ><td className="text-right"
    > {$ps.Size}</td
  ><td className="text-right"
    > {$ps.Resident}</td
  ><td className="text-center"
    >{$ps.Time}</td
  ><td
    >{$ps.Name}</td
  ></tr
>); },
		panelps:  function(Data, rows) { return (<div  className={!Data.Params.Psn.Negative ? "" : "panel panel-default"}
  ><div className="h4 padding-left-like-panel-heading"
    ><a  href={Data.Params.Tlinks.Psn} onClick={this.handleClick}
      >Processes</a
    ></div
  ><ul   className={!Data.Params.Psn.Negative ? "collapse-hidden" : "list-group"}
    ><li className="list-group-item text-nowrap th"
      ><ul className="list-inline"
        ><li
          ><span
            ><b
              >Delay</b
            > <span className="badge"
              >{Data.Params.Psd}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Dlinks.Psd.Less.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Psd.Less.ExtraClass != null ? Data.Params.Dlinks.Psd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Dlinks.Psd.Less.Text}</a
><a href={Data.Params.Dlinks.Psd.More.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Psd.More.ExtraClass != null ? Data.Params.Dlinks.Psd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Dlinks.Psd.More.Text} +</a
></div
          ></li
        ><li
          ><span
            ><b
              >Rows</b
            > <span className="badge"
              >{Data.Params.Psn.Absolute}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Nlinks.Psn.Less.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Psn.Less.ExtraClass != null ? Data.Params.Nlinks.Psn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Nlinks.Psn.Less.Text}</a
><a href={Data.Params.Nlinks.Psn.More.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Psn.More.ExtraClass != null ? Data.Params.Nlinks.Psn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Nlinks.Psn.More.Text} +</a
></div
          ></li
        ></ul
      ></li
    ></ul
  ><table  className={Data.Params.Psn.Absolute != 0 ? "table table-hover" : "collapse-hidden"}
    ><thead
      ><tr className="text-nowrap"
        ><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[1-1].LinkHref} className={Data.Params.Vlinks.Psk[1-1].LinkClass} onClick={this.handleClick}  
    >PID<span className={Data.Params.Vlinks.Psk[1-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[2-1].LinkHref} className={Data.Params.Vlinks.Psk[2-1].LinkClass} onClick={this.handleClick}  
    >UID<span className={Data.Params.Vlinks.Psk[2-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Params.Vlinks.Psk[3-1].LinkHref} className={Data.Params.Vlinks.Psk[3-1].LinkClass} onClick={this.handleClick}  
    >USER<span className={Data.Params.Vlinks.Psk[3-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[4-1].LinkHref} className={Data.Params.Vlinks.Psk[4-1].LinkClass} onClick={this.handleClick}  
    >PR<span className={Data.Params.Vlinks.Psk[4-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[5-1].LinkHref} className={Data.Params.Vlinks.Psk[5-1].LinkClass} onClick={this.handleClick}  
    >NI<span className={Data.Params.Vlinks.Psk[5-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[6-1].LinkHref} className={Data.Params.Vlinks.Psk[6-1].LinkClass} onClick={this.handleClick}  
    >VIRT<span className={Data.Params.Vlinks.Psk[6-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[7-1].LinkHref} className={Data.Params.Vlinks.Psk[7-1].LinkClass} onClick={this.handleClick}  
    >RES<span className={Data.Params.Vlinks.Psk[7-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-center"
  ><a href={Data.Params.Vlinks.Psk[8-1].LinkHref} className={Data.Params.Vlinks.Psk[8-1].LinkClass} onClick={this.handleClick}  
    >TIME<span className={Data.Params.Vlinks.Psk[8-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Params.Vlinks.Psk[9-1].LinkHref} className={Data.Params.Vlinks.Psk[9-1].LinkClass} onClick={this.handleClick}  
    >COMMAND<span className={Data.Params.Vlinks.Psk[9-1].CaretClass}
      ></span
    ></a
  ></th
></tr
      ></thead
    ><tbody
      >{rows}</tbody
    ></table
  ></div
>); },

		vg_rows:  function(Data, $vgm) { return (<tr  key={"vagrant-rowby-uuid-"+$vgm.UUID}
  ><td
    >{$vgm.UUID}</td
  ><td
    >{$vgm.Name}</td
  ><td
    >{$vgm.Provider}</td
  ><td
    >{$vgm.State}</td
  ><td
    >{$vgm.Vagrantfile_path}</td
  ></tr
>); },
		vg_error: function(Data)       { return (<tr key="vgerror"
  ><td colSpan="5"
    >{Data.VagrantError}</td
  ></tr
>); },
		panelvg:  function(Data, rows) { return (<div  className={!Data.Params.Vgn.Negative ? "" : "panel panel-default"}
  ><div className="h4 padding-left-like-panel-heading"
    ><a  href={Data.Params.Tlinks.Vgn} onClick={this.handleClick}
      >Vagrant</a
    ></div
  ><ul   className={!Data.Params.Vgn.Negative ? "collapse-hidden" : "list-group"}
    ><li className="list-group-item text-nowrap th"
      ><ul className="list-inline"
        ><li
          ><span
            ><b
              >Delay</b
            > <span className="badge"
              >{Data.Params.Vgd}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Dlinks.Vgd.Less.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Vgd.Less.ExtraClass != null ? Data.Params.Dlinks.Vgd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Dlinks.Vgd.Less.Text}</a
><a href={Data.Params.Dlinks.Vgd.More.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Vgd.More.ExtraClass != null ? Data.Params.Dlinks.Vgd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Dlinks.Vgd.More.Text} +</a
></div
          ></li
        ><li
          ><span
            ><b
              >Rows</b
            > <span className="badge"
              >{Data.Params.Vgn.Absolute}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Nlinks.Vgn.Less.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Vgn.Less.ExtraClass != null ? Data.Params.Nlinks.Vgn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Nlinks.Vgn.Less.Text}</a
><a href={Data.Params.Nlinks.Vgn.More.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Vgn.More.ExtraClass != null ? Data.Params.Nlinks.Vgn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Nlinks.Vgn.More.Text} +</a
></div
          ></li
        ></ul
      ></li
    ></ul
  ><table  className={Data.Params.Vgn.Absolute != 0 ? "table table-hover" : "collapse-hidden"}
    ><thead
      ><tr
        ><th
          >ID</th
        ><th
          >Name</th
        ><th
          >Provider</th
        ><th
          >State</th
        ><th
          >Directory</th
        ></tr
      ></thead
    ><tbody
      >{rows}</tbody
    ></table
  ></div
>); }
	};
});
