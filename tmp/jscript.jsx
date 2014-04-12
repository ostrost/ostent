
/** @jsx React.DOM */ // -*- indent-tabs-mode: nil -*-// variable origin: jscript.amber
var ifsTableClass = React.createClass({
	getInitialState: function() { return Data.Interfaces; },
	render: function() {
		var Data = {Interfaces: this.state};
		var ifs_rows = Data.Interfaces.List.map(function($if) {
			return (<tr key={$if.Name}><td>{$if.Name}</td><td className="digital">{$if.DeltaIn}</td><td className="digital">{$if.DeltaOut}</td><td className="digital">{$if.In}</td><td className="digital">{$if.Out}</td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th>Interface</th><th className="digital nobr">In<span className="unit">ps</span></th><th className="digital nobr">Out<span className="unit">ps</span></th><th className="digital nobr">In<span className="unit">%4G</span></th><th className="digital nobr">Out<span className="unit">%4G</span></th></tr></thead><tbody>{ifs_rows}</tbody></table>);
		
	}
});

function diskname_function(disk) {
	if (disk.ShorDiskName === "") {
		return (<span>{disk.DiskName}</span>);
	}
	var span = (<span title={disk.DiskName} className="tooltipable" data-toggle="tooltip" data-placement="left">{disk.ShortDiskName}...</span>);
	$('span .tooltipable').tooltip();
	return span;
}

// variable origin: jscript.amber
var diskTableClass = React.createClass({
	getInitialState: function() { return Data.DiskTable; },
	render: function() {
		var Data = {DiskTable: this.state};
		
		var df_rows = Data.DiskTable.List.map(function($disk) {
			return (<tr key={$disk.DiskName}><td>{diskname_function($disk)}</td><td>{$disk.DirName}</td><td className="digital">{$disk.Avail}</td><td className="digital">{$disk.Used}&nbsp;<sup><span className={$disk.UsePercentClass}>{$disk.UsePercent}%</span></sup></td><td className="digital">{$disk.Total}</td></tr>);
		});
		

		
		return (<table className="table1 stripe-table"><thead><tr><th className="header">        <a href={Data.DiskTable.Links.DiskName.Href} className={Data.DiskTable.Links.DiskName.Class}>Device</a></th><th className="header">        <a href={Data.DiskTable.Links.DirName.Href}  className={Data.DiskTable.Links.DirName.Class}>Mounted</a></th><th className="header digital"><a href={Data.DiskTable.Links.Avail.Href}    className={Data.DiskTable.Links.Avail.Class}>Avail</a></th><th className="header digital"><a href={Data.DiskTable.Links.Used.Href}     className={Data.DiskTable.Links.Used.Class}>Used</a></th><th className="header digital"><a href={Data.DiskTable.Links.Total.Href}    className={Data.DiskTable.Links.Total.Class}>Total</a></th></tr></thead><tbody>{df_rows}</tbody></table>);
		
	}
});

// variable origin: jscript.amber
var inodeTableClass = React.createClass({
	getInitialState: function() { return Data.DiskTable; },
	render: function() {
		var Data = {DiskTable: this.state};
		
		var di_rows = Data.DiskTable.List.map(function($disk) {
			return (<tr key={$disk.DiskName}><td>{diskname_function($disk)}</td><td className="digital">{$disk.Ifree}</td><td className="digital">{$disk.Iused}&nbsp;<sup><span className={$disk.IusePercentClass}>{$disk.IusePercent}%</span></sup></td><td className="digital">{$disk.Inodes}</td></tr>);
		});
		

		
		return (<table className="table1 stripe-table"><thead><tr><th className="header">        Device</th><th className="header digital">Avail</th><th className="header digital">Used</th><th className="header digital">Total</th></tr></thead><tbody>{di_rows}</tbody></table>);
		
	}
});

// variable origin: jscript.amber
var cpuTableClass = React.createClass({
	getInitialState: function() { return Data.CPU; },
	render: function() {
		var Data = {CPU: this.state};
		var cpu_rows = Data.CPU.List.map(function($core) {
			return (<tr key={$core.N}><td className="digital nobr">{$core.N}</td><td className="digital"><span id="core0.User" className={$core.UserClass}>{$core.User}</span></td><td className="digital"><span id="core0.Sys"  className={$core.SysClass}>{$core.Sys}</span></td><td className="digital"><span id="core0.Idle" className={$core.IdleClass}>{$core.Idle}</span></td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th></th><th className="digital">User%</th><th className="digital">Sys%</th><th className="digital">Idle%</th></tr></thead><tbody>{cpu_rows}</tbody></table>);
		
	}
});

// origin: jscript.amber
var procTableClass  = React.createClass({
	getInitialState: function() { return Data.ProcTable; },
	render: function() {
		var Data = {ProcTable: this.state};
		var ps_rows = Data.ProcTable.List.map(function($proc) {
			return (<tr key={$proc.PID}><td className="digital">{$proc.PID}</td><td className="digital">{$proc.User}</td><td className="digital">{$proc.Priority}</td><td className="digital">{$proc.Nice}</td><td className="digital">{$proc.Size}</td><td className="digital">{$proc.Resident}</td><td className="center">{$proc.Time}</td><td>{$proc.Name}</td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th className="header digital"><a href={Data.ProcTable.Links.PID.Href}      className={Data.ProcTable.Links.PID.Class}>PID</a></th><th className="header digital"><a href={Data.ProcTable.Links.User.Href}     className={Data.ProcTable.Links.User.Class}>USER</a></th><th className="header digital"><a href={Data.ProcTable.Links.Priority.Href} className={Data.ProcTable.Links.Priority.Class}>PR</a></th><th className="header digital"><a href={Data.ProcTable.Links.Nice.Href}     className={Data.ProcTable.Links.Nice.Class}>NI</a></th><th className="header digital"><a href={Data.ProcTable.Links.Size.Href}     className={Data.ProcTable.Links.Size.Class}>VIRT</a></th><th className="header digital"><a href={Data.ProcTable.Links.Resident.Href} className={Data.ProcTable.Links.Resident.Class}>RES</a></th><th className="header center"> <a href={Data.ProcTable.Links.Time.Href}     className={Data.ProcTable.Links.Time.Class}>TIME</a></th><th className="header">        <a href={Data.ProcTable.Links.Name.Href}     className={Data.ProcTable.Links.Name.Class}>COMMAND</a></th></tr></thead><tbody id="procrows">{ps_rows}</tbody></table>);
		
	}
});
