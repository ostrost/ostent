
/** @jsx React.DOM */ function emptyK(obj, key) {
return (obj      === undefined ||
	obj      === null      ||
	obj[key] === undefined ||
	obj[key] === null);
}

var ifsTableClass = React.createClass({
	getInitialState: function() { return Data.InterfacesBytes; },

	render: function() {
		var Data = {InterfacesBytes: this.state};
		var ifs_rows = emptyK(Data.InterfacesBytes, 'List') ?'': Data.InterfacesBytes.List.map(function($if) {
			return (<tr key={$if.NameKey}><td><span dangerouslySetInnerHTML={{__html: $if.NameHTML}} /></td><td className="digital">{$if.DeltaIn}</td><td className="digital">{$if.DeltaOut}</td><td className="digital">{$if.In}</td><td className="digital">{$if.Out}</td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th>Interface</th><th className="digital nobr">In<span className="unit" title="bits per second"><b>b</b>ps</span></th><th className="digital nobr">Out<span className="unit" title="bits per second"><b>b</b>ps</span></th><th className="digital nobr">In<span className="unit" title="total modulo 4G">%4G</span></th><th className="digital nobr">Out<span className="unit" title="total modulo 4G">%4G</span></th></tr></thead><tbody>{ifs_rows}</tbody></table>);
		
	}
});

var ifsErrorsTableClass = React.createClass({
	getInitialState: function() { return Data.InterfacesErrors; },

	render: function() {
		var Data = {InterfacesErrors: this.state};
		var ifs_rows = emptyK(Data.InterfacesErrors, 'List') ?'': Data.InterfacesErrors.List.map(function($if) {
			return (<tr key={$if.NameKey}><td><span dangerouslySetInnerHTML={{__html: $if.NameHTML}} /></td><td className="digital">{$if.DeltaIn}</td><td className="digital">{$if.DeltaOut}</td><td className="digital">{$if.In}</td><td className="digital">{$if.Out}</td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th>Interface</th><th className="digital nobr">In&nbsp;<span className="unit" title="per second">ps</span></th><th className="digital nobr">Out&nbsp;<span className="unit" title="per second">ps</span></th><th className="digital nobr">In<span className="unit" title="modulo 4G">%4G</span></th><th className="digital nobr">Out<span className="unit" title="modulo 4G">%4G</span></th></tr></thead><tbody>{ifs_rows}</tbody></table>);
		
	}
});

var ifsPacketsTableClass = React.createClass({
	getInitialState: function() { return Data.InterfacesPackets; },

	render: function() {
		var Data = {InterfacesPackets: this.state};
		var ifs_rows = emptyK(Data.InterfacesPackets, 'List') ?'': Data.InterfacesPackets.List.map(function($if) {
			return (<tr key={$if.NameKey}><td><span dangerouslySetInnerHTML={{__html: $if.NameHTML}} /></td><td className="digital">{$if.DeltaIn}</td><td className="digital">{$if.DeltaOut}</td><td className="digital">{$if.In}</td><td className="digital">{$if.Out}</td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th>Interface</th><th className="digital nobr">In&nbsp;<span className="unit" title="per second">ps</span></th><th className="digital nobr">Out&nbsp;<span className="unit" title="per second">ps</span></th><th className="digital nobr">In<span className="unit" title="modulo 4G">%4G</span></th><th className="digital nobr">Out<span className="unit" title="modulo 4G">%4G</span></th></tr></thead><tbody>{ifs_rows}</tbody></table>);
		
	}
});

var disksinBytesClass = React.createClass({
	getInitialState: function() { return {DiskLinks: Data.DiskLinks, DisksinBytes: Data.DisksinBytes}; },

	render: function() {
		var Data = this.state;
		var rows = emptyK(Data.DisksinBytes, 'List') ?'': Data.DisksinBytes.List.map(function($disk) {
			return (<tr key={$disk.DiskNameKey}><td className="nobr"><span dangerouslySetInnerHTML={{__html: $disk.DiskNameHTML}} /></td><td className="nobr"><span dangerouslySetInnerHTML={{__html: $disk.DirNameHTML}} /></td><td className="digital">{$disk.Avail}</td><td className="digital">{$disk.Used}&nbsp;<sup><span className={$disk.UsePercentClass}>{$disk.UsePercent}%</span></sup></td><td className="digital">{$disk.Total}</td></tr>);
		});

		
		return (<table className="table1 stripe-table"><thead><tr><th className="header">        <a href={Data.DiskLinks.DiskName.Href} className={Data.DiskLinks.DiskName.Class}>Device</a></th><th className="header">        <a href={Data.DiskLinks.DirName.Href}  className={Data.DiskLinks.DirName.Class}>Mounted</a></th><th className="header digital"><a href={Data.DiskLinks.Avail.Href}    className={Data.DiskLinks.Avail.Class}>Avail</a></th><th className="header digital"><a href={Data.DiskLinks.Used.Href}     className={Data.DiskLinks.Used.Class}>Used</a></th><th className="header digital"><a href={Data.DiskLinks.Total.Href}    className={Data.DiskLinks.Total.Class}>Total</a></th></tr></thead><tbody>{rows}</tbody></table>);
		
	}
});

var disksinInodesClass = React.createClass({
	getInitialState: function() { return {DiskLinks: Data.DiskLinks, DisksinInodes: Data.DisksinInodes}; },

	render: function() {
		var Data = this.state;
		var rows = emptyK(Data.DisksinInodes, 'List') ?'': Data.DisksinInodes.List.map(function($disk) {
			return (<tr key={$disk.DiskNameKey}><td className="nobr"><span dangerouslySetInnerHTML={{__html: $disk.DiskNameHTML}} /></td><td className="nobr"><span dangerouslySetInnerHTML={{__html: $disk.DirNameHTML}} /></td><td className="digital">{$disk.Ifree}</td><td className="digital">{$disk.Iused}&nbsp;<sup><span className={$disk.IusePercentClass}>{$disk.IusePercent}%</span></sup></td><td className="digital">{$disk.Inodes}</td></tr>);
		});

		
		return (<table className="table1 stripe-table"><thead><tr><th className="header">Device</th><th className="header">Mounted</th><th className="header digital">Avail</th><th className="header digital">Used</th><th className="header digital">Total</th></tr></thead><tbody>{rows}</tbody></table>);
		
	}
});

var cpuTableClass = React.createClass({
	getInitialState: function() { return Data.CPU; },

	render: function() {
		var Data = {CPU: this.state};
		var cpu_rows = Data.CPU.List.map(function($core) {
			return (<tr key={$core.N}><td className="digital nobr">{$core.N}</td><td className="digital"><span id="core0.User" className={$core.UserClass}>{$core.User}</span></td><td className="digital"><span id="core0.Sys"  className={$core.SysClass}>{$core.Sys}</span></td><td className="digital"><span id="core0.Idle" className={$core.IdleClass}>{$core.Idle}</span></td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th></th><th className="digital nobr">User<span className="unit">%</span></th><th className="digital nobr">Sys<span className="unit">%</span></th><th className="digital nobr">Idle<span className="unit">%</span></th></tr></thead><tbody>{cpu_rows}</tbody></table>);
		
	}
});

var procTableClass  = React.createClass({
	getInitialState: function() { return Data.ProcTable; },

	render: function() {
		var Data = {ProcTable: this.state};
		var ps_rows = Data.ProcTable.List.map(function($proc) {
			return (<tr key={$proc.PID}><td className="digital">{$proc.PID}</td><td className="digital"><span dangerouslySetInnerHTML={{__html: $proc.UserHTML}} /></td><td className="digital">{$proc.Priority}</td><td className="digital">{$proc.Nice}</td><td className="digital">{$proc.Size}</td><td className="digital">{$proc.Resident}</td><td className="center">{$proc.Time}</td><td className="nobr"><span dangerouslySetInnerHTML={{__html: $proc.NameHTML}} /></td></tr>);
		});
		
		return (<table className="table2 stripe-table"><thead><tr><th className="header digital"><a href={Data.ProcTable.Links.PID.Href}      className={Data.ProcTable.Links.PID.Class}>PID</a></th><th className="header digital"><a href={Data.ProcTable.Links.User.Href}     className={Data.ProcTable.Links.User.Class}>USER</a></th><th className="header digital"><a href={Data.ProcTable.Links.Priority.Href} className={Data.ProcTable.Links.Priority.Class}>PR</a></th><th className="header digital"><a href={Data.ProcTable.Links.Nice.Href}     className={Data.ProcTable.Links.Nice.Class}>NI</a></th><th className="header digital"><a href={Data.ProcTable.Links.Size.Href}     className={Data.ProcTable.Links.Size.Class}>VIRT</a></th><th className="header digital"><a href={Data.ProcTable.Links.Resident.Href} className={Data.ProcTable.Links.Resident.Class}>RES</a></th><th className="header center"> <a href={Data.ProcTable.Links.Time.Href}     className={Data.ProcTable.Links.Time.Class}>TIME</a></th><th className="header">        <a href={Data.ProcTable.Links.Name.Href}     className={Data.ProcTable.Links.Name.Class}>COMMAND</a></th></tr></thead><tbody id="procrows">{ps_rows}</tbody></table>);
		
	}
});
