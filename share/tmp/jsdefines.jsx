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
		panelmem:        function(Data, rows)  { return (<div
  ><div
    ><a       href={Data.Params.Toggle.Configmem} onClick={this.handleClick} className="btn-block"
      >  <span  className={Data.Params.Configmem ? "h4 bg-info" : "h4"}
        >Memory</span
      ></a
    ></div
  ><div
    ><div  className={Data.Params.Configmem ? "config-margintop" : "config-margintop collapse-hidden"} id="memconfig"
      ><form  action={"/form/"+Data.Params} className="form-inline"
        ><input className="hidden-submit" type="submit"
        ></input
      ><div className="btn-toolbar"
        ><div className="btn-group btn-group-sm" role="group"
          ><a  className={Data.Params.Hidemem ? "btn btn-default active" : "btn btn-default"}
    href={Data.Params.Toggle.Hidemem} onClick={this.handleClick} 
            >Hidden</a
          ><a  className={Data.Params.Hideswap ? "btn btn-default active" : "btn btn-default"}
    href={Data.Params.Toggle.Hideswap} onClick={this.handleClick}
            >Hide swap</a
          ></div
        ><div className="btn-group btn-group-sm" role="group"
          ><div  className={Data.Params.Errors && Data.Params.Errors.Refreshmem ? "input-group input-group-sm refresh-group has-warning" : "input-group input-group-sm refresh-group"}
  ><span className="input-group-addon"
    >Refresh</span
  >  <input className="form-control refresh-input width-fourem" type="text" placeholder={Data.MinRefresh}  name="refreshmem"  value={Data.Params.Refreshmem} onChange={this.handleChange}
  ></input></div
></div
        ></div
      ></form
    ></div
  ></div
><div
  ><div  className={Data.Params.Hidemem ? "collapse-hidden" : ""}
    ><table className="table table-hover"
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
  ></div
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
		panelif:         function(Data,r1,r2,r3){ return (<div
  ><div
    ><a       href={Data.Params.Toggle.Configif} onClick={this.handleClick} className="btn-block"
      >  <span  className={Data.Params.Configif ? "h4 bg-info" : "h4"}
        >Interfaces</span
      ></a
    ></div
  ><div
    ><div  className={Data.Params.Configif ? "config-margintop" : "config-margintop collapse-hidden"} id="ifconfig"
      ><form  action={"/form/"+Data.Params} className="form-inline"
        ><input className="hidden-submit" type="submit"
        ></input
      ><div className="btn-toolbar"
        ><div className="btn-group btn-group-sm" role="group"
          ><a  className={Data.Params.Hideif ? "btn btn-default active" : "btn btn-default"}
    href={Data.Params.Toggle.Hideif} onClick={this.handleClick}
            >Hidden</a
          ><a  className={Data.ExpandableIF ? "btn btn-default" : "btn btn-default disabled"}
    href={Data.Params.Toggle.Expandif} onClick={this.handleClick}
            >{Data.ExpandtextIF}</a
          ></div
        ><div className="btn-group btn-group-sm" role="group"
          ><div  className={Data.Params.Errors && Data.Params.Errors.Refreshif ? "input-group input-group-sm refresh-group has-warning" : "input-group input-group-sm refresh-group"}
  ><span className="input-group-addon"
    >Refresh</span
  >  <input className="form-control refresh-input width-fourem" type="text" placeholder={Data.MinRefresh}  name="refreshif"  value={Data.Params.Refreshif} onChange={this.handleChange}
  ></input></div
></div
        ></div
      ></form
    ><ul className="nav nav-tabs config-margintop"
      ><li  className={Data.Params.Ift == 1 ? "active" : ""}
        ><a href={Data.Params.Variations.Ift[1-1].LinkHref} onClick={this.handleClick}
  >Packets</a
></li
      ><li  className={Data.Params.Ift == 2 ? "active" : ""}
        ><a href={Data.Params.Variations.Ift[2-1].LinkHref} onClick={this.handleClick}
  >Errors</a
></li
      ><li  className={Data.Params.Ift == 3 ? "active" : ""}
        ><a href={Data.Params.Variations.Ift[3-1].LinkHref} onClick={this.handleClick}
  >Bytes</a
></li
      ></ul
    ></div
  ></div
><div
  ><div  className={Data.Params.Hideif ? "collapse-hidden" : ""}
    ><div  className={Data.Params.Ift == 1 ? "" : "collapse-hidden"}
      ><table className="table table-hover"
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
></div
    ><div  className={Data.Params.Ift == 2 ? "" : "collapse-hidden"}
      ><table className="table table-hover"
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
></div
    ><div  className={Data.Params.Ift == 3 ? "" : "collapse-hidden"}
      ><table className="table table-hover"
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
    ></div
  ></div
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
    ><span className="usepercent-text-inverse" data-usepercent={$core.Idle}
      >{$core.Idle}</span
    ></td
  ></tr
>); },
		panelcpu:        function(Data, rows)  { return (<div
  ><div
    ><a       href={Data.Params.Toggle.Configcpu} onClick={this.handleClick} className="btn-block"
      >  <span  className={Data.Params.Configcpu ? "h4 bg-info" : "h4"}
        >CPU</span
      ></a
    ></div
  ><div
    ><div  className={Data.Params.Configcpu ? "config-margintop" : "config-margintop collapse-hidden"} id="cpuconfig"
      ><form  action={"/form/"+Data.Params} className="form-inline"
        ><input className="hidden-submit" type="submit"
        ></input
      ><div className="btn-toolbar"
        ><div className="btn-group btn-group-sm" role="group"
          ><a  className={Data.Params.Hidecpu ? "btn btn-default active" : "btn btn-default"}
    href={Data.Params.Toggle.Hidecpu} onClick={this.handleClick}
            >Hidden</a
          ><a  className={Data.CPU.ExpandableCPU ? "btn btn-default" : "btn btn-default disabled"}
    href={Data.Params.Toggle.Expandcpu} onClick={this.handleClick}
            >{Data.CPU.ExpandtextCPU}</a
          ></div
        ><div className="btn-group btn-group-sm" role="group"
          ><div  className={Data.Params.Errors && Data.Params.Errors.Refreshcpu ? "input-group input-group-sm refresh-group has-warning" : "input-group input-group-sm refresh-group"}
  ><span className="input-group-addon"
    >Refresh</span
  >  <input className="form-control refresh-input width-fourem" type="text" placeholder={Data.MinRefresh}  name="refreshcpu"  value={Data.Params.Refreshcpu} onChange={this.handleChange}
  ></input></div
></div
        ></div
      ></form
    ></div
  ></div
><div
  ><div  className={Data.Params.Hidecpu ? "collapse-hidden" : ""}
    ><table className="table table-hover"
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
        >Idle<span className="unit"
          >%</span
        ></th
      ></tr
    ></thead
  ><tbody
    >{rows}</tbody
  ></table
></div
  ></div
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
		paneldf:         function(Data,r1,r2)  { return (<div
  ><div
    ><a     href={Data.Params.Toggle.Configdf} onClick={this.handleClick} className="btn-block"
      ><span  className={Data.Params.Configdf ? "h4 bg-info" : "h4"}
        >Disk usage</span
      ></a
    ></div
  ><div
    ><div  className={Data.Params.Configdf ? "config-margintop" : "config-margintop collapse-hidden"} id="dfconfig"
      ><form  action={"/form/"+Data.Params} className="form-inline"
        ><input className="hidden-submit" type="submit"
        ></input
      ><div className="btn-toolbar"
        ><div className="btn-group btn-group-sm" role="group"
          ><a  className={Data.Params.Hidedf ? "btn btn-default active" : "btn btn-default"}
    href={Data.Params.Toggle.Hidedf} onClick={this.handleClick}
            >Hidden</a
          ><a  className={Data.ExpandableDF ? "btn btn-default" : "btn btn-default disabled"}
    href={Data.Params.Toggle.Expanddf} onClick={this.handleClick}
            >{Data.ExpandtextDF}</a
          ></div
        ><div className="btn-group btn-group-sm" role="group"
          ><div  className={Data.Params.Errors && Data.Params.Errors.Refreshdf ? "input-group input-group-sm refresh-group has-warning" : "input-group input-group-sm refresh-group"}
  ><span className="input-group-addon"
    >Refresh</span
  >  <input className="form-control refresh-input width-fourem" type="text" placeholder={Data.MinRefresh}  name="refreshdf"  value={Data.Params.Refreshdf} onChange={this.handleChange}
  ></input></div
></div
        ></div
      ></form
    ><ul className="nav nav-tabs config-margintop"
      ><li  className={Data.Params.Dft == 1 ? "active" : ""}
        ><a href={Data.Params.Variations.Dft[1-1].LinkHref} onClick={this.handleClick}
  >Inodes</a
></li
      ><li  className={Data.Params.Dft == 2 ? "active" : ""}
        ><a href={Data.Params.Variations.Dft[2-1].LinkHref} onClick={this.handleClick}
  >Bytes</a
></li
      ></ul
    ></div
  ></div
><div
  ><div  className={Data.Params.Hidedf ? "collapse-hidden" : ""}
    ><div  className={Data.Params.Dft == 1 ? "" : "collapse-hidden"}
      ><table className="table table-hover"
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
></div
    ><div  className={Data.Params.Dft == 2 ? "" : "collapse-hidden"}
      ><table className="table table-hover"
  ><thead
    ><tr
      ><th className="header "
  ><a href={Data.Params.Variations.Dfk[1-1].LinkHref} className={Data.Params.Variations.Dfk[1-1].LinkClass}
    >Device
  <span className={Data.Params.Variations.Dfk[1-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Params.Variations.Dfk[2-1].LinkHref} className={Data.Params.Variations.Dfk[2-1].LinkClass}
    >Mounted
  <span className={Data.Params.Variations.Dfk[2-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Variations.Dfk[3-1].LinkHref} className={Data.Params.Variations.Dfk[3-1].LinkClass}
    >Avail
  <span className={Data.Params.Variations.Dfk[3-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Variations.Dfk[4-1].LinkHref} className={Data.Params.Variations.Dfk[4-1].LinkClass}
    >Used
  <span className={Data.Params.Variations.Dfk[4-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Variations.Dfk[5-1].LinkHref} className={Data.Params.Variations.Dfk[5-1].LinkClass}
    >Total
  <span className={Data.Params.Variations.Dfk[5-1].CaretClass}
      ></span
    ></a
  ></th
></tr
    ></thead
  ><tbody
    >{r2}</tbody
  ></table
></div
    ></div
  ></div
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
		panelps:         function(Data, rows)  { return (<div      className={(Data.Params.Psn != "!0" && Data.Params.Psn >= 0) ? "" : "panel panel-default"}
  >  <div    className={(Data.Params.Psn != "!0" && Data.Params.Psn >= 0) ? "" : "panel-heading"}
    >    <a    href={Data.Params.Toggle.Psn} onClick={this.handleClick} className="panel-title btn-block"
      >      <b  className={(Data.Params.Psn != "!0" && Data.Params.Psn >= 0) ? "h4" : "h4 bg-info"}
        >Processes</b
      >    </a
    >  </div
  >  <table  className={(Data.Params.Psn != "!0" && Data.Params.Psn >= 0) ? "table collapse-hidden" : "table"} id="psconfig"
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="nowrap text-right"
          >Delay&nbsp;<span className="badge"
            >{Data.Params.Psd}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Delayed.Psd.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Delayed.Psd.Less.Class != null ? Data.Params.Delayed.Psd.Less.Class : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Delayed.Psd.Less.Text}</a
><a href={Data.Params.Delayed.Psd.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Delayed.Psd.More.Class != null ? Data.Params.Delayed.Psd.More.Class : "")}
  
  >{Data.Params.Delayed.Psd.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td className="col-md-10" colSpan="2"
        ></td
      ></tr
    ><tr className="panel-config"
      ><td className="col-md-2"
        ><div className="nowrap text-right"
          >Rows&nbsp;<span className="badge"
            >{Data.Params.Psn == "!0" ? 0 : (Data.Params.Psn < 0 ? -Data.Params.Psn : Data.Params.Psn)}</span
          ></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Numbered.Psn.Less.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Numbered.Psn.Less.Class != null ? Data.Params.Numbered.Psn.Less.Class : "")}
  
  ><span className="xlabel xlabel-default"
    >-</span
  > {Data.Params.Numbered.Psn.Less.Text}</a
><a href={Data.Params.Numbered.Psn.More.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Numbered.Psn.More.Class != null ? Data.Params.Numbered.Psn.More.Class : "")}
  
  >{Data.Params.Numbered.Psn.More.Text} <span className="xlabel xlabel-default"
    >+</span
  ></a
></div
        ></td
      ><td
        ><div className="btn-group nowrap-group" role="group"
          ><a href={Data.Params.Numbered.Psn.Zero.Href} onClick={this.handleClick} className={"btn btn-default" + " " + (Data.Params.Numbered.Psn.Zero.Class != null ? Data.Params.Numbered.Psn.Zero.Class : "")}
  
  >{Data.Params.Numbered.Psn.Zero.Text} <span className="xlabel xlabel-default"
    ></span
  ></a
></div
        ></td
      ><td className="col-md-10"
        ></td
      ></tr
    >  </table
  >  <table  className={(Data.Params.Psn == "!0" || Data.Params.Psn == 0) ? "collapse-hidden" : "table table-hover"}
  ><thead
    ><tr
      ><th className="header text-right"
  ><a href={Data.Params.Variations.Psk[1-1].LinkHref} className={Data.Params.Variations.Psk[1-1].LinkClass}
    >PID
  <span className={Data.Params.Variations.Psk[1-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Variations.Psk[8-1].LinkHref} className={Data.Params.Variations.Psk[8-1].LinkClass}
    >UID
  <span className={Data.Params.Variations.Psk[8-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Params.Variations.Psk[9-1].LinkHref} className={Data.Params.Variations.Psk[9-1].LinkClass}
    >USER
  <span className={Data.Params.Variations.Psk[9-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Variations.Psk[2-1].LinkHref} className={Data.Params.Variations.Psk[2-1].LinkClass}
    >PR
  <span className={Data.Params.Variations.Psk[2-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Variations.Psk[3-1].LinkHref} className={Data.Params.Variations.Psk[3-1].LinkClass}
    >NI
  <span className={Data.Params.Variations.Psk[3-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Variations.Psk[4-1].LinkHref} className={Data.Params.Variations.Psk[4-1].LinkClass}
    >VIRT
  <span className={Data.Params.Variations.Psk[4-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Variations.Psk[5-1].LinkHref} className={Data.Params.Variations.Psk[5-1].LinkClass}
    >RES
  <span className={Data.Params.Variations.Psk[5-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-center"
  ><a href={Data.Params.Variations.Psk[6-1].LinkHref} className={Data.Params.Variations.Psk[6-1].LinkClass}
    >TIME
  <span className={Data.Params.Variations.Psk[6-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Params.Variations.Psk[7-1].LinkHref} className={Data.Params.Variations.Psk[7-1].LinkClass}
    >COMMAND
  <span className={Data.Params.Variations.Psk[7-1].CaretClass}
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
		panelvg:         function(Data, rows)  { return (<div
  ><div
    ><a       href={Data.Params.Toggle.Configvg} onClick={this.handleClick} className="btn-block"
      >  <span  className={Data.Params.Configvg ? "h4 bg-info" : "h4"}
        >Vagrant</span
      ></a
    ></div
  ><div
    ><div  className={Data.Params.Configvg ? "config-margintop" : "config-margintop collapse-hidden"} id="vgconfig"
      ><form  action={"/form/"+Data.Params} className="form-inline"
        ><input className="hidden-submit" type="submit"
        ></input
      ><div className="btn-toolbar"
        ><div className="btn-group btn-group-sm" role="group"
          ><a  className={Data.Params.Hidevg ? "btn btn-default active" : "btn btn-default"}
    href={Data.Params.Toggle.Hidevg} onClick={this.handleClick}
            >Hidden</a
          ></div
        ><div className="btn-group btn-group-sm" role="group"
          ><div  className={Data.Params.Errors && Data.Params.Errors.Refreshvg ? "input-group input-group-sm refresh-group has-warning" : "input-group input-group-sm refresh-group"}
  ><span className="input-group-addon"
    >Refresh</span
  >  <input className="form-control refresh-input width-fourem" type="text" placeholder={Data.MinRefresh}  name="refreshvg"  value={Data.Params.Refreshvg} onChange={this.handleChange}
  ></input></div
></div
        ></div
      ></form
    ></div
  ></div
><div
  ><div  className={Data.Params.Hidevg ? "collapse-hidden" : ""}
    ><table className="table table-hover"
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
  ></div
></div
>); }
	};
});
