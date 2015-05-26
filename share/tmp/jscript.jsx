define(function(require) {
	var React = require('react');
	return {
		mem_rows:        function(Data, $mem)  { return (<tr key={"mem-rowby-kind-"+$mem.Kind}
  ><td
    >{$mem.Kind}</td
  ><td className="text-right"
    >{$mem.Free}</td
  ><td className="text-right"
    >{$mem.Used}&nbsp;<sup
      ><span  className={LabelClassColorPercent($mem.UsePercent)}
  >{$mem.UsePercent}%</span
></sup
    ></td
  ><td className="text-right"
    >{$mem.Total}</td
  ></tr
>); },
		mem_table:       function(Data, rows)  { return (<table className="table1 stripe-table"
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
>); },

		ifbytes_rows:    function(Data, $if)   { return (<tr key={"ifbytes-rowby-name-"+$if.Name}
  ><td
    ><input id={"if-bytes-name-"+$if.Name}  className="collapse-checkbox" type="checkbox" aria-hidden="true" hidden
  ></input
><label htmlFor={"if-bytes-name-"+$if.Name} className="clip" style={{maxWidth: '12ch'}}
  >{$if.Name}</label
></td
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
		ifbytes_table:   function(Data, rows)  { return (<table className="table1 stripe-table"
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
    >{rows}</tbody
  ></table
>); },
		iferrors_rows:   function(Data, $if)   { return (<tr key={"iferrors-rowby-name-"+$if.Name}
  ><td
    ><input id={"if-errors-name-"+$if.Name}  className="collapse-checkbox" type="checkbox" aria-hidden="true" hidden
  ></input
><label htmlFor={"if-errors-name-"+$if.Name} className="clip" style={{maxWidth: '12ch'}}
  >{$if.Name}</label
></td
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
		iferrors_table:  function(Data, rows)  { return (<table className="table1 stripe-table"
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
    >{rows}</tbody
  ></table
>); },
		ifpackets_rows:  function(Data, $if)   { return (<tr key={"ifpackets-rowby-name-"+$if.Name}
  ><td
    ><input id={"if-packets-name-"+$if.Name}  className="collapse-checkbox" type="checkbox" aria-hidden="true" hidden
  ></input
><label htmlFor={"if-packets-name-"+$if.Name} className="clip" style={{maxWidth: '12ch'}}
  >{$if.Name}</label
></td
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
		ifpackets_table: function(Data, rows)  { return (<table className="table1 stripe-table"
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
    >{rows}</tbody
  ></table
>); },

		cpu_rows:        function(Data, $core) { return (<tr key={"cpu-rowby-N-"+$core.N}
  ><td className="text-right nowrap"
    >{$core.N}</td
  ><td className="text-right"
    ><span className={$core.UserClass}
      >{$core.User}</span
    ></td
  ><td className="text-right"
    ><span className={$core.SysClass} 
      >{$core.Sys}</span
    ></td
  ><td className="text-right"
    ><span className={$core.IdleClass}
      >{$core.Idle}</span
    ></td
  ></tr
>); },
		cpu_table:       function(Data, rows)  { return (<table className="table1 stripe-table"
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
>); },

		dfbytes_rows:    function(Data, $disk) { return (<tr key={"dfbytes-rowby-dirname-"+$disk.DirName}
  ><td className="nowrap"
    ><input id={"df-bytes-devname-"+$disk.DevName}  className="collapse-checkbox" type="checkbox" aria-hidden="true" hidden
  ></input
><label htmlFor={"df-bytes-devname-"+$disk.DevName} className="clip" style={{maxWidth: '12ch'}}
  >{$disk.DevName}</label
></td
  ><td className="nowrap"
    ><input id={"df-bytes-dirname-"+$disk.DirName}  className="collapse-checkbox" type="checkbox" aria-hidden="true" hidden
  ></input
><label htmlFor={"df-bytes-dirname-"+$disk.DirName} className="clip" style={{maxWidth: '6ch'}}
  >{$disk.DirName}</label
></td
  ><td className="text-right"
    >{$disk.Avail}</td
  ><td className="text-right"
    >{$disk.Used}&nbsp;<sup
      ><span  className={LabelClassColorPercent($disk.UsePercent)}
  >{$disk.UsePercent}%</span
></sup
    ></td
  ><td className="text-right"
    >{$disk.Total}</td
  ></tr
>); },
		dfbytes_table:   function(Data, rows)  { return (<table className="table1 stripe-table"
  ><thead
    ><tr
      ><th className="header "
  ><a href={Data.Links.Params.df.FS.Href} className={Data.Links.Params.df.FS.Class}
    >{Data.Links.Params.df.FS.Text}<span className={Data.Links.Params.df.FS.CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Links.Params.df.MP.Href} className={Data.Links.Params.df.MP.Class}
    >{Data.Links.Params.df.MP.Text}<span className={Data.Links.Params.df.MP.CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Links.Params.df.AVAIL.Href} className={Data.Links.Params.df.AVAIL.Class}
    >{Data.Links.Params.df.AVAIL.Text}<span className={Data.Links.Params.df.AVAIL.CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Links.Params.df.USED.Href} className={Data.Links.Params.df.USED.Class}
    >{Data.Links.Params.df.USED.Text}<span className={Data.Links.Params.df.USED.CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Links.Params.df.TOTAL.Href} className={Data.Links.Params.df.TOTAL.Class}
    >{Data.Links.Params.df.TOTAL.Text}<span className={Data.Links.Params.df.TOTAL.CaretClass}
      ></span
    ></a
  ></th
></tr
    ></thead
  ><tbody
    >{rows}</tbody
  ></table
>); },
		dfinodes_rows:   function(Data, $disk) { return (<tr key={"dfinodes-rowby-dirname-"+$disk.DirName}
  ><td className="nowrap"
    ><input id={"df-inodes-devname-"+$disk.DevName}  className="collapse-checkbox" type="checkbox" aria-hidden="true" hidden
  ></input
><label htmlFor={"df-inodes-devname-"+$disk.DevName} className="clip" style={{maxWidth: '12ch'}}
  >{$disk.DevName}</label
></td
  ><td className="nowrap"
    ><input id={"df-inodes-devname-"+$disk.DirName}  className="collapse-checkbox" type="checkbox" aria-hidden="true" hidden
  ></input
><label htmlFor={"df-inodes-devname-"+$disk.DirName} className="clip" style={{maxWidth: '6ch'}}
  >{$disk.DirName}</label
></td
  ><td className="text-right"
    >{$disk.Ifree}</td
  ><td className="text-right"
    >{$disk.Iused}&nbsp;<sup
      ><span  className={LabelClassColorPercent($disk.IusePercent)}
  >{$disk.IusePercent}%</span
></sup
    ></td
  ><td className="text-right"
    >{$disk.Inodes}</td
  ></tr
>); },
		dfinodes_table:  function(Data, rows)  { return (<table className="table1 stripe-table"
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
    >{rows}</tbody
  ></table
>); },

		ps_rows:         function(Data, $proc) { return (<tr key={"ps-rowby-pid-"+$proc.PIDstring}
  ><td className="text-right"
    > {$proc.PID}</td
  ><td className="text-right"
    > {$proc.UID}</td
  ><td
    >            <input id={"psuser-pid-"+$proc.PIDstring}  className="collapse-checkbox" type="checkbox" aria-hidden="true" hidden
  ></input
><label htmlFor={"psuser-pid-"+$proc.PIDstring} className="clip" style={{maxWidth: '12ch'}}
  >{$proc.User}</label
></td
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
  ><td className="nowrap"
    >     <input id={"psname-pid-"+$proc.PIDstring}  className="collapse-checkbox" type="checkbox" aria-hidden="true" hidden
  ></input
><label htmlFor={"psname-pid-"+$proc.PIDstring} className="clip" style={{maxWidth: '42ch'}}
  >{$proc.Name}</label
></td
  ></tr
>); },
		ps_table:        function(Data, rows)  { return (<table className="table2 stripe-table"
  ><thead
    ><tr
      ><th className="header text-right"
  ><a href={Data.Links.Params.ps.PID.Href} className={Data.Links.Params.ps.PID.Class}
    >{Data.Links.Params.ps.PID.Text}<span className={Data.Links.Params.ps.PID.CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Links.Params.ps.UID.Href} className={Data.Links.Params.ps.UID.Class}
    >{Data.Links.Params.ps.UID.Text}<span className={Data.Links.Params.ps.UID.CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Links.Params.ps.USER.Href} className={Data.Links.Params.ps.USER.Class}
    >{Data.Links.Params.ps.USER.Text}<span className={Data.Links.Params.ps.USER.CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Links.Params.ps.PRI.Href} className={Data.Links.Params.ps.PRI.Class}
    >{Data.Links.Params.ps.PRI.Text}<span className={Data.Links.Params.ps.PRI.CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Links.Params.ps.NICE.Href} className={Data.Links.Params.ps.NICE.Class}
    >{Data.Links.Params.ps.NICE.Text}<span className={Data.Links.Params.ps.NICE.CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Links.Params.ps.VIRT.Href} className={Data.Links.Params.ps.VIRT.Class}
    >{Data.Links.Params.ps.VIRT.Text}<span className={Data.Links.Params.ps.VIRT.CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Links.Params.ps.RES.Href} className={Data.Links.Params.ps.RES.Class}
    >{Data.Links.Params.ps.RES.Text}<span className={Data.Links.Params.ps.RES.CaretClass}
      ></span
    ></a
  ></th
><th className="header text-center"
  ><a href={Data.Links.Params.ps.TIME.Href} className={Data.Links.Params.ps.TIME.Class}
    >{Data.Links.Params.ps.TIME.Text}<span className={Data.Links.Params.ps.TIME.CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Links.Params.ps.NAME.Href} className={Data.Links.Params.ps.NAME.Class}
    >{Data.Links.Params.ps.NAME.Text}<span className={Data.Links.Params.ps.NAME.CaretClass}
      ></span
    ></a
  ></th
></tr
    ></thead
  ><tbody
    >{rows}</tbody
  ></table
>); },

		vagrant_rows:    function(Data, $mach) { return (<tr key={"vagrant-rowby-uuid-"+$mach.UUID}
  ><td
    >       <input id={"vagrant-uuid-"+$mach.UUID}  className="collapse-checkbox" type="checkbox" aria-hidden="true" hidden
  ></input
><label htmlFor={"vagrant-uuid-"+$mach.UUID} className="clip" style={{maxWidth: '7ch'}}
  >{$mach.UUID}</label
></td
  ><td
    >       {$mach.Name}</td
  ><td
    >       {$mach.Provider}</td
  ><td
    >       <input id={"vagrant-state-"+$mach.UUID}  className="collapse-checkbox" type="checkbox" aria-hidden="true" hidden
  ></input
><label htmlFor={"vagrant-state-"+$mach.UUID} className="clip" style={{maxWidth: '8ch'}}
  >{$mach.State}</label
></td
  ><td
    >       <input id={"vagrant-filepath-"+$mach.UUID}  className="collapse-checkbox" type="checkbox" aria-hidden="true" hidden
  ></input
><label htmlFor={"vagrant-filepath-"+$mach.UUID} className="clip" style={{maxWidth: '50ch'}}
  >{$mach.Vagrantfile_path}</label
></td
  ></tr
>); },
		vagrant_error:   function(Data)        { return (<tr key="vgerror"
  ><td colspan="5"
    >{Data.VagrantError}</td
  ></tr
>); },
		vagrant_table:   function(Data, rows)  { return (<table id="vgtable" className="table1 stripe-table"
  ><thead
    ><tr
      ><th
        >id</th
      ><th
        >name</th
      ><th
        >provider</th
      ><th
        >state</th
      ><th
        >directory</th
      ></tr
    ></thead
  ><tbody
    >{rows}</tbody
  ></table
>); }
	};
});
