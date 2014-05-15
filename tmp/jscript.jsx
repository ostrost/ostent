/** @jsx React.DOM */ function emptyK(obj, key) {
return (obj      === undefined ||
	obj      === null      ||
	obj[key] === undefined ||
	obj[key] === null);
}

var IFbytesCLASS = React.createClass({
	getInitialState: function() { return Data.IFbytes; },

	render: function() {
		var Data = {IFbytes: this.state};
		var rows = emptyK(Data.IFbytes, 'List') ?'': Data.IFbytes.List.map(function($if) {
			return (<tr key={$if.NameKey}><td><span dangerouslySetInnerHTML={{__html: $if.NameHTML}} /></td><td className="digital">{$if.DeltaIn}</td><td className="digital">{$if.DeltaOut}</td><td className="digital">{$if.In}</td><td className="digital">{$if.Out}</td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th>Interface</th><th className="digital nobr">In<span className="unit" title="bits per second"><b>b</b>ps</span></th><th className="digital nobr">Out<span className="unit" title="bits per second"><b>b</b>ps</span></th><th className="digital nobr">In<span className="unit" title="total modulo 4G">%4G</span></th><th className="digital nobr">Out<span className="unit" title="total modulo 4G">%4G</span></th></tr></thead><tbody>{rows}</tbody></table>);
		
	}
});

var IFerrorsCLASS = React.createClass({
	getInitialState: function() { return Data.IFerrors; },

	render: function() {
		var Data = {IFerrors: this.state};
		var rows = emptyK(Data.IFerrors, 'List') ?'': Data.IFerrors.List.map(function($if) {
			return (<tr key={$if.NameKey}><td><span dangerouslySetInnerHTML={{__html: $if.NameHTML}} /></td><td className="digital">{$if.DeltaIn}</td><td className="digital">{$if.DeltaOut}</td><td className="digital">{$if.In}</td><td className="digital">{$if.Out}</td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th>Interface</th><th className="digital nobr">In&nbsp;<span className="unit" title="per second">ps</span></th><th className="digital nobr">Out&nbsp;<span className="unit" title="per second">ps</span></th><th className="digital nobr">In<span className="unit" title="modulo 4G">%4G</span></th><th className="digital nobr">Out<span className="unit" title="modulo 4G">%4G</span></th></tr></thead><tbody>{rows}</tbody></table>);
		
	}
});

var IFpacketsCLASS = React.createClass({
	getInitialState: function() { return Data.IFpackets; },

	render: function() {
		var Data = {IFpackets: this.state};
		var rows = emptyK(Data.IFpackets, 'List') ?'': Data.IFpackets.List.map(function($if) {
			return (<tr key={$if.NameKey}><td><span dangerouslySetInnerHTML={{__html: $if.NameHTML}} /></td><td className="digital">{$if.DeltaIn}</td><td className="digital">{$if.DeltaOut}</td><td className="digital">{$if.In}</td><td className="digital">{$if.Out}</td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th>Interface</th><th className="digital nobr">In&nbsp;<span className="unit" title="per second">ps</span></th><th className="digital nobr">Out&nbsp;<span className="unit" title="per second">ps</span></th><th className="digital nobr">In<span className="unit" title="modulo 4G">%4G</span></th><th className="digital nobr">Out<span className="unit" title="modulo 4G">%4G</span></th></tr></thead><tbody>{rows}</tbody></table>);
		
	}
});

var DFbytesCLASS = React.createClass({
	getInitialState: function() { return {DFlinks: Data.DFlinks, DFbytes: Data.DFbytes}; },

	render: function() {
		var Data = this.state;
		var rows = emptyK(Data.DFbytes, 'List') ?'': Data.DFbytes.List.map(function($disk) {
			return (<tr key={$disk.DirNameKey}><td className="nobr"><span dangerouslySetInnerHTML={{__html: $disk.DiskNameHTML}} /></td><td className="nobr"><span dangerouslySetInnerHTML={{__html: $disk.DirNameHTML}} /></td><td className="digital">{$disk.Avail}</td><td className="digital">{$disk.Used}&nbsp;<sup><span className={$disk.UsePercentClass}>{$disk.UsePercent}%</span></sup></td><td className="digital">{$disk.Total}</td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th className="header">        <a href={Data.DFlinks.DiskName.Href} className={Data.DFlinks.DiskName.Class}>Device<span  className={Data.DFlinks.DiskName.CaretClass} /></a></th><th className="header">        <a href={Data.DFlinks.DirName.Href}  className={Data.DFlinks.DirName.Class} >Mounted<span className={Data.DFlinks.DirName.CaretClass}  /></a></th><th className="header digital"><a href={Data.DFlinks.Avail.Href}    className={Data.DFlinks.Avail.Class}   >Avail<span   className={Data.DFlinks.Avail.CaretClass}    /></a></th><th className="header digital"><a href={Data.DFlinks.Used.Href}     className={Data.DFlinks.Used.Class}    >Used<span    className={Data.DFlinks.Used.CaretClass}     /></a></th><th className="header digital"><a href={Data.DFlinks.Total.Href}    className={Data.DFlinks.Total.Class}   >Total<span   className={Data.DFlinks.Total.CaretClass}    /></a></th></tr></thead><tbody>{rows}</tbody></table>);
		
	}
});

var DFinodesCLASS = React.createClass({
	getInitialState: function() { return {DFlinks: Data.DFlinks, DFinodes: Data.DFinodes}; },

	render: function() {
		var Data = this.state;
		var rows = emptyK(Data.DFinodes, 'List') ?'': Data.DFinodes.List.map(function($disk) {
			return (<tr key={$disk.DirNameKey}><td className="nobr"><span dangerouslySetInnerHTML={{__html: $disk.DiskNameHTML}} /></td><td className="nobr"><span dangerouslySetInnerHTML={{__html: $disk.DirNameHTML}} /></td><td className="digital">{$disk.Ifree}</td><td className="digital">{$disk.Iused}&nbsp;<sup><span className={$disk.IusePercentClass}>{$disk.IusePercent}%</span></sup></td><td className="digital">{$disk.Inodes}</td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th className="header">Device</th><th className="header">Mounted</th><th className="header digital">Avail</th><th className="header digital">Used</th><th className="header digital">Total</th></tr></thead><tbody>{rows}</tbody></table>);
		
	}
});

var CPUtableCLASS = React.createClass({
	getInitialState: function() { return Data.CPU; },

	render: function() {
		var Data = {CPU: this.state};
		var rows = Data.CPU.List.map(function($core) {
			return (<tr key={$core.N}><td className="digital nobr">{$core.N}</td><td className="digital"><span id="core0.User" className={$core.UserClass}>{$core.User}</span></td><td className="digital"><span id="core0.Sys"  className={$core.SysClass}>{$core.Sys}</span></td><td className="digital"><span id="core0.Idle" className={$core.IdleClass}>{$core.Idle}</span></td></tr>);
		});
		
		return (<table className="table1 stripe-table"><thead><tr><th></th><th className="digital nobr">User<span className="unit">%</span></th><th className="digital nobr">Sys<span className="unit">%</span></th><th className="digital nobr">Idle<span className="unit">%</span></th></tr></thead><tbody>{rows}</tbody></table>);
		
	}
});

var PStableCLASS = React.createClass({
	getInitialState: function() { return Data.PStable; },

	render: function() {
		var Data = {PStable: this.state};
		var rows = Data.PStable.List.map(function($proc) {
			return (<tr key={$proc.PID}><td className="digital">{$proc.PID}</td><td className="digital"><span dangerouslySetInnerHTML={{__html: $proc.UserHTML}} /></td><td className="digital">{$proc.Priority}</td><td className="digital">{$proc.Nice}</td><td className="digital">{$proc.Size}</td><td className="digital">{$proc.Resident}</td><td className="center">{$proc.Time}</td><td className="nobr"><span dangerouslySetInnerHTML={{__html: $proc.NameHTML}} /></td></tr>);
		});
		
		return (<table className="table2 stripe-table"><thead><tr><th className="header digital"><a href={Data.PStable.Links.PID.Href}      className={Data.PStable.Links.PID.Class}     >PID<span     className={Data.PStable.Links.PID.CaretClass}      /></a></th><th className="header digital"><a href={Data.PStable.Links.User.Href}     className={Data.PStable.Links.User.Class}    >USER<span    className={Data.PStable.Links.User.CaretClass}     /></a></th><th className="header digital"><a href={Data.PStable.Links.Priority.Href} className={Data.PStable.Links.Priority.Class}>PR<span      className={Data.PStable.Links.Priority.CaretClass} /></a></th><th className="header digital"><a href={Data.PStable.Links.Nice.Href}     className={Data.PStable.Links.Nice.Class}    >NI<span      className={Data.PStable.Links.Nice.CaretClass}     /></a></th><th className="header digital"><a href={Data.PStable.Links.Size.Href}     className={Data.PStable.Links.Size.Class}    >VIRT<span    className={Data.PStable.Links.Size.CaretClass}     /></a></th><th className="header digital"><a href={Data.PStable.Links.Resident.Href} className={Data.PStable.Links.Resident.Class}>RES<span     className={Data.PStable.Links.Resident.CaretClass} /></a></th><th className="header center"> <a href={Data.PStable.Links.Time.Href}     className={Data.PStable.Links.Time.Class}    >TIME<span    className={Data.PStable.Links.Time.CaretClass}     /></a></th><th className="header">        <a href={Data.PStable.Links.Name.Href}     className={Data.PStable.Links.Name.Class}    >COMMAND<span className={Data.PStable.Links.Name.CaretClass}     /></a></th></tr></thead><tbody>{rows}</tbody></table>);
		
	}
});
