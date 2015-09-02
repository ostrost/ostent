define(function(require) {
	var React = require('react');
	return {
		mem_rows:        function(Data, $mem)  { return (<tr  key={"mem-rowby-kind-"+$mem.Kind}
  ><td
    >{$mem.Kind}</td
  ><td className="text-right"
    >{$mem.Free}</td
  ><td className="text-right"
    ><span className="label" data-usepercent={$mem.UsePercent}
  >{$mem.UsePercent}%</span
>&nbsp;{$mem.Used}</td
  ><td className="text-right"
    >{$mem.Total}</td
  ></tr
>); },
		panelmem:        function(Data, rows)  { return (<div      className={!Data.Params.Memn.Negative ? "" : "panel panel-default"}
  >  <div    className={!Data.Params.Memn.Negative ? "" : "panel-heading"}
    >    <a    href={Data.Params.Tlinks.Memn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={!Data.Params.Memn.Negative ? "h4" : "h4 bg-info"}
        >Memory</b
      >    </a
    >  </div
  >  <table  className={!Data.Params.Memn.Negative ? "table collapse-hidden" : "table"}
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="nowrap text-right"
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
        ><div className="nowrap text-right"
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

		ifbytes_rows:    function(Data, $if)   { return (<tr  key={"ifbytes-rowby-name-"+$if.Name}
  ><td
    >{$if.Name}</td
  ><td className="text-right"
    >{$if.DeltaIn}</td
  ><td className="text-right"
    >{$if.DeltaOut}</td
  ><td className="text-right"
    >{$if.In}</td
  ><td className="text-right"
    >{$if.Out}</td
  ></tr
>); },
		iferrors_rows:   function(Data, $if)   { return (<tr  key={"iferrors-rowby-name-"+$if.Name}
  ><td
    >{$if.Name}</td
  ><td className="text-right"
    >{$if.DeltaIn}</td
  ><td className="text-right"
    >{$if.DeltaOut}</td
  ><td className="text-right"
    >{$if.In}</td
  ><td className="text-right"
    >{$if.Out}</td
  ></tr
>); },
		ifpackets_rows:  function(Data, $if)   { return (<tr  key={"ifpackets-rowby-name-"+$if.Name}
  ><td
    >{$if.Name}</td
  ><td className="text-right"
    >{$if.DeltaIn}</td
  ><td className="text-right"
    >{$if.DeltaOut}</td
  ><td className="text-right"
    >{$if.In}</td
  ><td className="text-right"
    >{$if.Out}</td
  ></tr
>); },
		panelif:         function(Data,r1,r2,r3){ return (<div      className={!Data.Params.Ifn.Negative ? "" : "panel panel-default"}
  >  <div    className={!Data.Params.Ifn.Negative ? "" : "panel-heading"}
    >    <a    href={Data.Params.Tlinks.Ifn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={!Data.Params.Ifn.Negative ? "h4" : "h4 bg-info"}
        >Interfaces</b
      >    </a
    >  </div
  >  <table  className={!Data.Params.Ifn.Negative ? "table collapse-hidden" : "table"}
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="nowrap text-right"
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
        ><div className="nowrap text-right"
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
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right"
          >Select</div
        ></td
      ><td colSpan="2"
        ><ul className="nav nav-pills"
          ><li  className={Data.Params.Ifn.Absolute != 0 && Data.Params.Ift.Absolute == 1 ? "active" : ""}
            ><a href={Data.Params.Vlinks.Ift[1-1].LinkHref} onClick={this.handleClick}
  >Packets</a
></li
          ><li  className={Data.Params.Ifn.Absolute != 0 && Data.Params.Ift.Absolute == 2 ? "active" : ""}
            ><a href={Data.Params.Vlinks.Ift[2-1].LinkHref} onClick={this.handleClick}
  >Errors</a
></li
          ><li  className={Data.Params.Ifn.Absolute != 0 && Data.Params.Ift.Absolute == 3 ? "active" : ""}
            ><a href={Data.Params.Vlinks.Ift[3-1].LinkHref} onClick={this.handleClick}
  >Bytes</a
></li
          ></ul
        ></td
      ></tr
    >  </table
  >  <table  className={Data.Params.Ifn.Absolute != 0 && Data.Params.Ift.Absolute == 1 ? "table table-hover" : "collapse-hidden"}
  ><thead
    ><tr
      ><th
        >Interface</th
      ><th className="text-right nowrap" title="per second"
        >In&nbsp;<span className="unit"
          >ps</span
        ></th
      ><th className="text-right nowrap" title="per second"
        >Out&nbsp;<span className="unit"
          >ps</span
        ></th
      ><th className="text-right nowrap" title="total modulo 4G"
        >In&nbsp;<span className="unit"
          >%4G</span
        ></th
      ><th className="text-right nowrap" title="total modulo 4G"
        >Out&nbsp;<span className="unit"
          >%4G</span
        ></th
      ></tr
    ></thead
  ><tbody
    >{r1}</tbody
  ></table
>
  <table  className={Data.Params.Ifn.Absolute != 0 && Data.Params.Ift.Absolute == 2 ? "table table-hover" : "collapse-hidden"}
  ><thead
    ><tr
      ><th
        >Interface</th
      ><th className="text-right nowrap" title="per second"
        >In&nbsp;<span className="unit"
          >ps</span
        ></th
      ><th className="text-right nowrap" title="per second"
        >Out&nbsp;<span className="unit"
          >ps</span
        ></th
      ><th className="text-right nowrap" title="modulo 4G"
        >In&nbsp;<span className="unit"
          >%4G</span
        ></th
      ><th className="text-right nowrap" title="modulo 4G"
        >Out&nbsp;<span className="unit"
          >%4G</span
        ></th
      ></tr
    ></thead
  ><tbody
    >{r2}</tbody
  ></table
>
  <table  className={Data.Params.Ifn.Absolute != 0 && Data.Params.Ift.Absolute == 3 ? "table table-hover" : "collapse-hidden"}
  ><thead
    ><tr
      ><th
        >Interface</th
      ><th className="text-right nowrap" title="BITS per second"
        >In<span className="unit"
          ><i
            >b</i
          >ps</span
        ></th
      ><th className="text-right nowrap" title="BITS per second"
        >Out<span className="unit"
          ><i
            >b</i
          >ps</span
        ></th
      ><th className="text-right nowrap" title="total BYTES modulo 4G"
        >In<span className="unit"
          ><i
            >B</i
          >%4G</span
        ></th
      ><th className="text-right nowrap" title="total BYTES modulo 4G"
        >Out<span className="unit"
          ><i
            >B</i
          >%4G</span
        ></th
      ></tr
    ></thead
  ><tbody
    >{r3}</tbody
  ></table
></div
>); },

		cpu_rows:        function(Data, $core) { return (<tr  key={"cpu-rowby-N-"+$core.N}
  ><td className="text-right nowrap"
    >{$core.N}</td
  ><td className="text-right"
    ><span className="usepercent-text" data-usepercent={$core.User}
      >{$core.User}</span
    ></td
  ><td className="text-right"
    ><span className="usepercent-text" data-usepercent={$core.Sys}
      >{$core.Sys}</span
    ></td
  ><td className="text-right"
    ><span className="usepercent-text" data-usepercent={$core.Wait}
      >{$core.Wait}</span
    ></td
  ><td className="text-right"
    ><span className="usepercent-text-inverse" data-usepercent={$core.Idle}
      >{$core.Idle}</span
    ></td
  ></tr
>); },
		panelcpu:        function(Data, rows)  { return (<div      className={!Data.Params.CPUn.Negative ? "" : "panel panel-default"}
  >  <div    className={!Data.Params.CPUn.Negative ? "" : "panel-heading"}
    >    <a    href={Data.Params.Tlinks.CPUn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={!Data.Params.CPUn.Negative ? "h4" : "h4 bg-info"}
        >CPU</b
      >    </a
    >  </div
  >  <table  className={!Data.Params.CPUn.Negative ? "table collapse-hidden" : "table"}
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="nowrap text-right"
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
        ><div className="nowrap text-right"
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
      ><th className="text-right nowrap"
        >User<span className="unit"
          >%</span
        ></th
      ><th className="text-right nowrap"
        >Sys<span className="unit"
          >%</span
        ></th
      ><th className="text-right nowrap"
        >Wait<span className="unit"
          >%</span
        ></th
      ><th className="text-right nowrap"
        >Idle<span className="unit"
          >%</span
        ></th
      ></tr
    ></thead
  ><tbody
    >{rows}</tbody
  ></table
></div
>); },

		dfbytes_rows:    function(Data, $disk) { return (<tr  key={"dfbytes-rowby-dirname-"+$disk.DirName}
  ><td
    >{$disk.DevName}</td
  ><td
    >{$disk.DirName}</td
  ><td className="text-right"
    >{$disk.Avail}</td
  ><td className="text-right"
    ><span className="label" data-usepercent={$disk.UsePercent}
  >{$disk.UsePercent}%</span
>&nbsp;{$disk.Used}</td
  ><td className="text-right"
    >{$disk.Total}</td
  ></tr
>); },
		dfinodes_rows:   function(Data, $disk) { return (<tr  key={"dfinodes-rowby-dirname-"+$disk.DirName}
  ><td
    >{$disk.DevName}</td
  ><td
    >{$disk.DirName}</td
  ><td className="text-right"
    >{$disk.Ifree}</td
  ><td className="text-right"
    ><span className="label" data-usepercent={$disk.IusePercent}
  >{$disk.IusePercent}%</span
>&nbsp;{$disk.Iused}</td
  ><td className="text-right"
    >{$disk.Inodes}</td
  ></tr
>); },
		paneldf:         function(Data,r1,r2)  { return (<div      className={!Data.Params.Dfn.Negative ? "" : "panel panel-default"}
  >  <div    className={!Data.Params.Dfn.Negative ? "" : "panel-heading"}
    >    <a    href={Data.Params.Tlinks.Dfn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={!Data.Params.Dfn.Negative ? "h4" : "h4 bg-info"}
        >Disk usage</b
      >    </a
    >  </div
  >  <table  className={!Data.Params.Dfn.Negative ? "table collapse-hidden" : "table"}
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="nowrap text-right"
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
        ><div className="nowrap text-right"
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
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="text-right"
          >Select</div
        ></td
      ><td colSpan="2"
        ><ul className="nav nav-pills"
          ><li  className={Data.Params.Dfn.Absolute != 0 && Data.Params.Dft.Absolute == 1 ? "active" : ""}
            ><a href={Data.Params.Vlinks.Dft[1-1].LinkHref} onClick={this.handleClick}
  >Inodes</a
></li
          ><li  className={Data.Params.Dfn.Absolute != 0 && Data.Params.Dft.Absolute == 2 ? "active" : ""}
            ><a href={Data.Params.Vlinks.Dft[2-1].LinkHref} onClick={this.handleClick}
  >Bytes</a
></li
          ></ul
        ></td
      ></tr
    >  </table
  >  <table  className={Data.Params.Dfn.Absolute != 0 && Data.Params.Dft.Absolute == 1 ? "table table-hover" : "collapse-hidden"}
  ><thead
    ><tr
      ><th className="header"
        >Device</th
      ><th className="header"
        >Mounted</th
      ><th className="header text-right"
        >Avail</th
      ><th className="header text-right"
        >Used</th
      ><th className="header text-right"
        >Total</th
      ></tr
    ></thead
  ><tbody
    >{r1}</tbody
  ></table
>
  <table  className={Data.Params.Dfn.Absolute != 0 && Data.Params.Dft.Absolute == 2 ? "table table-hover" : "collapse-hidden"}
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
    >  Used<span className={Data.Params.Vlinks.Dfk[4-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[5-1].LinkHref} className={Data.Params.Vlinks.Dfk[5-1].LinkClass}
    >  Total<span className={Data.Params.Vlinks.Dfk[5-1].CaretClass}
      ></span
    ></a
  ></th
></tr
    ></thead
  ><tbody
    >{r2}</tbody
  ></table
></div
>); },

		ps_rows:         function(Data, $proc) { return (<tr  key={"ps-rowby-pid-"+$proc.PID}
  ><td className="text-right"
    > {$proc.PID}</td
  ><td className="text-right"
    > {$proc.UID}</td
  ><td
    >{$proc.User}</td
  ><td className="text-right"
    > {$proc.Priority}</td
  ><td className="text-right"
    > {$proc.Nice}</td
  ><td className="text-right"
    > {$proc.Size}</td
  ><td className="text-right"
    > {$proc.Resident}</td
  ><td className="text-center"
    >{$proc.Time}</td
  ><td
    >{$proc.Name}</td
  ></tr
>); },
		panelps:         function(Data, rows)  { return (<div      className={!Data.Params.Psn.Negative ? "" : "panel panel-default"}
  >  <div    className={!Data.Params.Psn.Negative ? "" : "panel-heading"}
    >    <a    href={Data.Params.Tlinks.Psn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={!Data.Params.Psn.Negative ? "h4" : "h4 bg-info"}
        >Processes</b
      >    </a
    >  </div
  >  <table  className={!Data.Params.Psn.Negative ? "table collapse-hidden" : "table"}
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="nowrap text-right"
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
        ><div className="nowrap text-right"
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

		vagrant_rows:    function(Data, $mach) { return (<tr  key={"vagrant-rowby-uuid-"+$mach.UUID}
  ><td
    >{$mach.UUID}</td
  ><td
    >{$mach.Name}</td
  ><td
    >{$mach.Provider}</td
  ><td
    >{$mach.State}</td
  ><td
    >{$mach.Vagrantfile_path}</td
  ></tr
>); },
		vagrant_error:   function(Data)        { return (<tr key="vgerror"
  ><td colSpan="5"
    >{Data.VagrantError}</td
  ></tr
>); },
		panelvg:         function(Data, rows)  { return (<div      className={!Data.Params.Vgn.Negative ? "" : "panel panel-default"}
  >  <div    className={!Data.Params.Vgn.Negative ? "" : "panel-heading"}
    >    <a    href={Data.Params.Tlinks.Vgn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={!Data.Params.Vgn.Negative ? "h4" : "h4 bg-info"}
        >Vagrant</b
      >    </a
    >  </div
  >  <table  className={!Data.Params.Vgn.Negative ? "table collapse-hidden" : "table"}
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="nowrap text-right"
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
        ><div className="nowrap text-right"
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
