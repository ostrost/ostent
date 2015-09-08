define(function(require) {
	var React = require('react');
	return {
		mem_rows: function(Data, $mem) { return (<tr  key={"mem-rowby-kind-"+$mem.Kind}
  ><td
    >{$mem.Kind}</td
  ><td className="text-right"
    >{$mem.Free}</td
  ><td className="text-right bg-usepct" data-usepct={$mem.UsePct}
    >{$mem.UsePct}%</td
  ><td className="text-right"
    >{$mem.Used}</td
  ><td className="text-right"
    >{$mem.Total}</td
  ></tr
>); },
		panelmem: function(Data, rows) { return (<div      className={!Data.Params.Memn.Negative ? "" : "panel panel-default"}
  >  <div    className={!Data.Params.Memn.Negative ? "" : "panel-heading"}
    >    <a    href={Data.Params.Tlinks.Memn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={!Data.Params.Memn.Negative ? "h4" : "h4 bg-info"}
        >Memory</b
      >    </a
    >  </div
  >  <table  className={!Data.Params.Memn.Negative ? "table collapse-hidden" : "table"}
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right text-nowrap"
          >Delay&nbsp;<span className="badge"
            >{Data.Params.Memd}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Dlinks.Memd.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Dlinks.Memd.Less.ExtraClass != null ? Data.Params.Dlinks.Memd.Less.ExtraClass : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Dlinks.Memd.Less.Text}</a
><a href={Data.Params.Dlinks.Memd.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Dlinks.Memd.More.ExtraClass != null ? Data.Params.Dlinks.Memd.More.ExtraClass : "")}
  
  >{Data.Params.Dlinks.Memd.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right text-nowrap"
          >Rows&nbsp;<span className="badge"
            >{Data.Params.Memn.Absolute}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Nlinks.Memn.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Nlinks.Memn.Less.ExtraClass != null ? Data.Params.Nlinks.Memn.Less.ExtraClass : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Nlinks.Memn.Less.Text}</a
><a href={Data.Params.Nlinks.Memn.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Nlinks.Memn.More.ExtraClass != null ? Data.Params.Nlinks.Memn.More.ExtraClass : "")}
  
  >{Data.Params.Nlinks.Memn.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    >  </table
  >  <table  className={Data.Params.Memn.Absolute != 0 ? "table table-hover" : "collapse-hidden"}
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
  >  <td className="text-nowrap clip12" title={$if.Name}
    >{$if.Name}</td
  ><td className="text-right text-nowrap"
    ><span className="text-graylighter" title="Total BYTES IN modulo 4G"
      >{$if.BytesIn}</span
    > <span title="Bits IN per second"
      >{$if.DeltaBitsIn}</span
    ></td
  ><td className="text-right text-nowrap"
    ><span className="text-graylighter" title="Total BYTES OUT modulo 4G"
      >{$if.BytesOut}</span
    > <span title="Bits OUT per second"
      >{$if.DeltaBitsOut}</span
    ></td
  ><td className="text-right text-nowrap"
    ><span className="text-graylighter" title="Total packets IN modulo 4G"
      >{$if.PacketsIn}</span
    > <span title="Packets IN per second"
      >{$if.DeltaPacketsIn}</span
    ></td
  ><td className="text-right text-nowrap"
    ><span className="text-graylighter" title="Total packets OUT modulo 4G"
      >{$if.PacketsOut}</span
    > <span title="Packets OUT per second"
      >{$if.DeltaPacketsOut}</span
    ></td
  ><td className="text-right text-nowrap"
    ><span className="text-graylighter" title="Total errors IN modulo 4G"
      >{$if.ErrorsIn}</span
    > <span title="Errors IN per second"
      >{$if.DeltaErrorsIn}</span
    ></td
  ><td className="text-right text-nowrap"
    ><span className="text-graylighter" title="Total errors OUT modulo 4G"
      >{$if.ErrorsOut}</span
    > <span title="Errors OUT per second"
      >{$if.DeltaErrorsOut}</span
    ></td
  ></tr
>); },
		panelif:  function(Data, rows) { return (<div      className={!Data.Params.Ifn.Negative ? "" : "panel panel-default"}
  >  <div    className={!Data.Params.Ifn.Negative ? "" : "panel-heading"}
    >    <a    href={Data.Params.Tlinks.Ifn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={!Data.Params.Ifn.Negative ? "h4" : "h4 bg-info"}
        >Interfaces</b
      >    </a
    >  </div
  >  <table  className={!Data.Params.Ifn.Negative ? "table collapse-hidden" : "table"}
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right text-nowrap"
          >Delay&nbsp;<span className="badge"
            >{Data.Params.Ifd}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Dlinks.Ifd.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Dlinks.Ifd.Less.ExtraClass != null ? Data.Params.Dlinks.Ifd.Less.ExtraClass : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Dlinks.Ifd.Less.Text}</a
><a href={Data.Params.Dlinks.Ifd.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Dlinks.Ifd.More.ExtraClass != null ? Data.Params.Dlinks.Ifd.More.ExtraClass : "")}
  
  >{Data.Params.Dlinks.Ifd.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right text-nowrap"
          >Rows&nbsp;<span className="badge"
            >{Data.Params.Ifn.Absolute}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Nlinks.Ifn.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Nlinks.Ifn.Less.ExtraClass != null ? Data.Params.Nlinks.Ifn.Less.ExtraClass : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Nlinks.Ifn.Less.Text}</a
><a href={Data.Params.Nlinks.Ifn.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Nlinks.Ifn.More.ExtraClass != null ? Data.Params.Nlinks.Ifn.More.ExtraClass : "")}
  
  >{Data.Params.Nlinks.Ifn.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    >  </table
  >  <table  className={Data.Params.Ifn.Absolute != 0 ? "table table-hover" : "collapse-hidden"}
  ><thead
    ><tr
      ><th
        >Interface</th
      ><th className="text-right normal" colSpan="6"
        >Bits <b title="Bits IN per second"
          >In</b
        >, <b title="Bits OUT per second"
          >Out</b
        > | Packets <b title="Packets IN per second"
          >In</b
        >, <b title="Packets OUT per second"
          >Out</b
        > | Errors <b title="Errors IN per second"
          >In</b
        >, <b title="Errors OUT per second"
          >Out</b
        ></th
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
  ><td className="text-right bg-usepct" data-usepct={$cpu.UserPct}
    >{$cpu.UserPct}%</td
  ><td className="text-right bg-usepct" data-usepct={$cpu.SysPct}
    >{$cpu.SysPct}%</td
  ><td className="text-right bg-usepct" data-usepct={$cpu.WaitPct}
    >{$cpu.WaitPct}%</td
  ><td className="text-right bg-usepct-inverse" data-usepct={$cpu.IdlePct}
    >{$cpu.IdlePct}%</td
  ></tr
>); },
		panelcpu: function(Data, rows) { return (<div      className={!Data.Params.CPUn.Negative ? "" : "panel panel-default"}
  >  <div    className={!Data.Params.CPUn.Negative ? "" : "panel-heading"}
    >    <a    href={Data.Params.Tlinks.CPUn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={!Data.Params.CPUn.Negative ? "h4" : "h4 bg-info"}
        >CPU</b
      >    </a
    >  </div
  >  <table  className={!Data.Params.CPUn.Negative ? "table collapse-hidden" : "table"}
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right text-nowrap"
          >Delay&nbsp;<span className="badge"
            >{Data.Params.CPUd}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Dlinks.CPUd.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Dlinks.CPUd.Less.ExtraClass != null ? Data.Params.Dlinks.CPUd.Less.ExtraClass : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Dlinks.CPUd.Less.Text}</a
><a href={Data.Params.Dlinks.CPUd.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Dlinks.CPUd.More.ExtraClass != null ? Data.Params.Dlinks.CPUd.More.ExtraClass : "")}
  
  >{Data.Params.Dlinks.CPUd.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right text-nowrap"
          >Rows&nbsp;<span className="badge"
            >{Data.Params.CPUn.Absolute}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Nlinks.CPUn.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Nlinks.CPUn.Less.ExtraClass != null ? Data.Params.Nlinks.CPUn.Less.ExtraClass : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Nlinks.CPUn.Less.Text}</a
><a href={Data.Params.Nlinks.CPUn.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Nlinks.CPUn.More.ExtraClass != null ? Data.Params.Nlinks.CPUn.More.ExtraClass : "")}
  
  >{Data.Params.Nlinks.CPUn.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    >  </table
  >  <table  className={Data.Params.CPUn.Absolute != 0 ? "table table-hover" : "collapse-hidden"}
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
    ><span className="text-graylighter" title="Inodes free"
      >{$df.Ifree}</span
    > {$df.Avail}</td
  ><td className="text-right text-nowrap bg-usepct" data-usepct={$df.UsePct}
    ><span className="text-graylighter" title="Inodes use%"
      >{$df.IusePct}%</span
    > {$df.UsePct}%</td
  ><td className="text-right text-nowrap"
    ><span className="text-graylighter" title="Inodes used"
      >{$df.Iused}</span
    > {$df.Used}</td
  ><td className="text-right text-nowrap"
    ><span className="text-graylighter" title="Inodes total"
      >{$df.Inodes}</span
    > {$df.Total}</td
  ></tr
>); },
		paneldf:  function(Data,rows)  { return (<div      className={!Data.Params.Dfn.Negative ? "" : "panel panel-default"}
  >  <div    className={!Data.Params.Dfn.Negative ? "" : "panel-heading"}
    >    <a    href={Data.Params.Tlinks.Dfn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={!Data.Params.Dfn.Negative ? "h4" : "h4 bg-info"}
        >Disk usage</b
      >    </a
    >  </div
  >  <table  className={!Data.Params.Dfn.Negative ? "table collapse-hidden" : "table"}
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right text-nowrap"
          >Delay&nbsp;<span className="badge"
            >{Data.Params.Dfd}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Dlinks.Dfd.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Dlinks.Dfd.Less.ExtraClass != null ? Data.Params.Dlinks.Dfd.Less.ExtraClass : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Dlinks.Dfd.Less.Text}</a
><a href={Data.Params.Dlinks.Dfd.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Dlinks.Dfd.More.ExtraClass != null ? Data.Params.Dlinks.Dfd.More.ExtraClass : "")}
  
  >{Data.Params.Dlinks.Dfd.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right text-nowrap"
          >Rows&nbsp;<span className="badge"
            >{Data.Params.Dfn.Absolute}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Nlinks.Dfn.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Nlinks.Dfn.Less.ExtraClass != null ? Data.Params.Nlinks.Dfn.Less.ExtraClass : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Nlinks.Dfn.Less.Text}</a
><a href={Data.Params.Nlinks.Dfn.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Nlinks.Dfn.More.ExtraClass != null ? Data.Params.Nlinks.Dfn.More.ExtraClass : "")}
  
  >{Data.Params.Nlinks.Dfn.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    >  </table
  >  <table  className={Data.Params.Dfn.Absolute != 0 ? "table table-hover" : "collapse-hidden"}
  ><thead
    ><tr
      ><th className="header "
  ><a href={Data.Params.Vlinks.Dfk[1-1].LinkHref} className={Data.Params.Vlinks.Dfk[1-1].LinkClass}
    >  Device<span className={Data.Params.Vlinks.Dfk[1-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Params.Vlinks.Dfk[2-1].LinkHref} className={Data.Params.Vlinks.Dfk[2-1].LinkClass}
    >  Mounted<span className={Data.Params.Vlinks.Dfk[2-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[3-1].LinkHref} className={Data.Params.Vlinks.Dfk[3-1].LinkClass}
    >  Avail<span className={Data.Params.Vlinks.Dfk[3-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[4-1].LinkHref} className={Data.Params.Vlinks.Dfk[4-1].LinkClass}
    >  Use%<span className={Data.Params.Vlinks.Dfk[4-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[5-1].LinkHref} className={Data.Params.Vlinks.Dfk[5-1].LinkClass}
    >  Used<span className={Data.Params.Vlinks.Dfk[5-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[6-1].LinkHref} className={Data.Params.Vlinks.Dfk[6-1].LinkClass}
    >  Total<span className={Data.Params.Vlinks.Dfk[6-1].CaretClass}
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
		panelps:  function(Data, rows) { return (<div      className={!Data.Params.Psn.Negative ? "" : "panel panel-default"}
  >  <div    className={!Data.Params.Psn.Negative ? "" : "panel-heading"}
    >    <a    href={Data.Params.Tlinks.Psn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={!Data.Params.Psn.Negative ? "h4" : "h4 bg-info"}
        >Processes</b
      >    </a
    >  </div
  >  <table  className={!Data.Params.Psn.Negative ? "table collapse-hidden" : "table"}
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right text-nowrap"
          >Delay&nbsp;<span className="badge"
            >{Data.Params.Psd}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Dlinks.Psd.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Dlinks.Psd.Less.ExtraClass != null ? Data.Params.Dlinks.Psd.Less.ExtraClass : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Dlinks.Psd.Less.Text}</a
><a href={Data.Params.Dlinks.Psd.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Dlinks.Psd.More.ExtraClass != null ? Data.Params.Dlinks.Psd.More.ExtraClass : "")}
  
  >{Data.Params.Dlinks.Psd.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right text-nowrap"
          >Rows&nbsp;<span className="badge"
            >{Data.Params.Psn.Absolute}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Nlinks.Psn.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Nlinks.Psn.Less.ExtraClass != null ? Data.Params.Nlinks.Psn.Less.ExtraClass : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Nlinks.Psn.Less.Text}</a
><a href={Data.Params.Nlinks.Psn.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Nlinks.Psn.More.ExtraClass != null ? Data.Params.Nlinks.Psn.More.ExtraClass : "")}
  
  >{Data.Params.Nlinks.Psn.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    >  </table
  >  <table  className={Data.Params.Psn.Absolute != 0 ? "table table-hover" : "collapse-hidden"}
  ><thead
    ><tr
      ><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[1-1].LinkHref} className={Data.Params.Vlinks.Psk[1-1].LinkClass}
    >  PID<span className={Data.Params.Vlinks.Psk[1-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[2-1].LinkHref} className={Data.Params.Vlinks.Psk[2-1].LinkClass}
    >  UID<span className={Data.Params.Vlinks.Psk[2-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Params.Vlinks.Psk[3-1].LinkHref} className={Data.Params.Vlinks.Psk[3-1].LinkClass}
    >  USER<span className={Data.Params.Vlinks.Psk[3-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[4-1].LinkHref} className={Data.Params.Vlinks.Psk[4-1].LinkClass}
    >  PR<span className={Data.Params.Vlinks.Psk[4-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[5-1].LinkHref} className={Data.Params.Vlinks.Psk[5-1].LinkClass}
    >  NI<span className={Data.Params.Vlinks.Psk[5-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[6-1].LinkHref} className={Data.Params.Vlinks.Psk[6-1].LinkClass}
    >  VIRT<span className={Data.Params.Vlinks.Psk[6-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[7-1].LinkHref} className={Data.Params.Vlinks.Psk[7-1].LinkClass}
    >  RES<span className={Data.Params.Vlinks.Psk[7-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-center"
  ><a href={Data.Params.Vlinks.Psk[8-1].LinkHref} className={Data.Params.Vlinks.Psk[8-1].LinkClass}
    >  TIME<span className={Data.Params.Vlinks.Psk[8-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Params.Vlinks.Psk[9-1].LinkHref} className={Data.Params.Vlinks.Psk[9-1].LinkClass}
    >  COMMAND<span className={Data.Params.Vlinks.Psk[9-1].CaretClass}
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
		panelvg:  function(Data, rows) { return (<div      className={!Data.Params.Vgn.Negative ? "" : "panel panel-default"}
  >  <div    className={!Data.Params.Vgn.Negative ? "" : "panel-heading"}
    >    <a    href={Data.Params.Tlinks.Vgn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={!Data.Params.Vgn.Negative ? "h4" : "h4 bg-info"}
        >Vagrant</b
      >    </a
    >  </div
  >  <table  className={!Data.Params.Vgn.Negative ? "table collapse-hidden" : "table"}
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right text-nowrap"
          >Delay&nbsp;<span className="badge"
            >{Data.Params.Vgd}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Dlinks.Vgd.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Dlinks.Vgd.Less.ExtraClass != null ? Data.Params.Dlinks.Vgd.Less.ExtraClass : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Dlinks.Vgd.Less.Text}</a
><a href={Data.Params.Dlinks.Vgd.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Dlinks.Vgd.More.ExtraClass != null ? Data.Params.Dlinks.Vgd.More.ExtraClass : "")}
  
  >{Data.Params.Dlinks.Vgd.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right text-nowrap"
          >Rows&nbsp;<span className="badge"
            >{Data.Params.Vgn.Absolute}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Nlinks.Vgn.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Nlinks.Vgn.Less.ExtraClass != null ? Data.Params.Nlinks.Vgn.Less.ExtraClass : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Nlinks.Vgn.Less.Text}</a
><a href={Data.Params.Nlinks.Vgn.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Nlinks.Vgn.More.ExtraClass != null ? Data.Params.Nlinks.Vgn.More.ExtraClass : "")}
  
  >{Data.Params.Nlinks.Vgn.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    >  </table
  >  <table  className={Data.Params.Vgn.Absolute != 0 ? "table table-hover" : "collapse-hidden"}
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
