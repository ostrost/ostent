/** @jsx React.DOM */ // -*- indent-tabs-mode: nil -*-// variable origin: jscript.amber
var ifsTableClass = React.createClass({displayName: 'ifsTableClass',
	getInitialState: function() { return Data.Interfaces; },
	render: function() {
		var Data = {Interfaces: this.state};
		var ifs_rows = Data.Interfaces.List.map(function($if) {
			return (React.DOM.tr( {key:$if.Name}, React.DOM.td(null, $if.Name),React.DOM.td( {className:"digital"}, $if.DeltaIn),React.DOM.td( {className:"digital"}, $if.DeltaOut),React.DOM.td( {className:"digital"}, $if.In),React.DOM.td( {className:"digital"}, $if.Out)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null, "Interface"),React.DOM.th( {className:"digital nobr"}, "In",React.DOM.span( {className:"unit"}, "ps")),React.DOM.th( {className:"digital nobr"}, "Out",React.DOM.span( {className:"unit"}, "ps")),React.DOM.th( {className:"digital nobr"}, "In",React.DOM.span( {className:"unit"}, "%4G")),React.DOM.th( {className:"digital nobr"}, "Out",React.DOM.span( {className:"unit"}, "%4G")))),React.DOM.tbody(null, ifs_rows)));
		
	}
});

function diskname_function(disk) {
	if (disk.ShorDiskName === "") {
		return (React.DOM.span(null, disk.DiskName));
	}
	var span = (React.DOM.span( {title:disk.DiskName, className:"tooltipable", 'data-toggle':"tooltip", 'data-placement':"left"}, disk.ShortDiskName,"..."));
	$('span .tooltipable').tooltip();
	return span;
}

// variable origin: jscript.amber
var diskTableClass = React.createClass({displayName: 'diskTableClass',
	getInitialState: function() { return Data.DiskTable; },
	render: function() {
		var Data = {DiskTable: this.state};
		
		var df_rows = Data.DiskTable.List.map(function($disk) {
			return (React.DOM.tr( {key:$disk.DiskName}, React.DOM.td(null, diskname_function($disk)),React.DOM.td(null, $disk.DirName),React.DOM.td( {className:"digital"}, $disk.Avail),React.DOM.td( {className:"digital"}, $disk.Used," ",React.DOM.sup(null, React.DOM.span( {className:$disk.UsePercentClass}, $disk.UsePercent,"%"))),React.DOM.td( {className:"digital"}, $disk.Total)));
		});
		

		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th( {className:"header"},         "        ",        React.DOM.a( {href:Data.DiskTable.Links.DiskName.Href, className:Data.DiskTable.Links.DiskName.Class}, "Device")),React.DOM.th( {className:"header"},         "        ",        React.DOM.a( {href:Data.DiskTable.Links.DirName.Href,  className:Data.DiskTable.Links.DirName.Class}, "Mounted")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.DiskTable.Links.Avail.Href,    className:Data.DiskTable.Links.Avail.Class}, "Avail")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.DiskTable.Links.Used.Href,     className:Data.DiskTable.Links.Used.Class}, "Used")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.DiskTable.Links.Total.Href,    className:Data.DiskTable.Links.Total.Class}, "Total")))),React.DOM.tbody(null, df_rows)));
		
	}
});

// variable origin: jscript.amber
var inodeTableClass = React.createClass({displayName: 'inodeTableClass',
	getInitialState: function() { return Data.DiskTable; },
	render: function() {
		var Data = {DiskTable: this.state};
		
		var di_rows = Data.DiskTable.List.map(function($disk) {
			return (React.DOM.tr( {key:$disk.DiskName}, React.DOM.td(null, diskname_function($disk)),React.DOM.td( {className:"digital"}, $disk.Ifree),React.DOM.td( {className:"digital"}, $disk.Iused," ",React.DOM.sup(null, React.DOM.span( {className:$disk.IusePercentClass}, $disk.IusePercent,"%"))),React.DOM.td( {className:"digital"}, $disk.Inodes)));
		});
		

		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th( {className:"header"},         "        Device"),React.DOM.th( {className:"header digital"}, "Avail"),React.DOM.th( {className:"header digital"}, "Used"),React.DOM.th( {className:"header digital"}, "Total"))),React.DOM.tbody(null, di_rows)));
		
	}
});

// variable origin: jscript.amber
var cpuTableClass = React.createClass({displayName: 'cpuTableClass',
	getInitialState: function() { return Data.CPU; },
	render: function() {
		var Data = {CPU: this.state};
		var cpu_rows = Data.CPU.List.map(function($core) {
			return (React.DOM.tr( {key:$core.N}, React.DOM.td( {className:"digital nobr"}, $core.N),React.DOM.td( {className:"digital"}, React.DOM.span( {id:"core0.User", className:$core.UserClass}, $core.User)),React.DOM.td( {className:"digital"}, React.DOM.span( {id:"core0.Sys",  className:$core.SysClass}, $core.Sys)),React.DOM.td( {className:"digital"}, React.DOM.span( {id:"core0.Idle", className:$core.IdleClass}, $core.Idle))));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null),React.DOM.th( {className:"digital"}, "User%"),React.DOM.th( {className:"digital"}, "Sys%"),React.DOM.th( {className:"digital"}, "Idle%"))),React.DOM.tbody(null, cpu_rows)));
		
	}
});

// origin: jscript.amber
var procTableClass  = React.createClass({displayName: 'procTableClass',
	getInitialState: function() { return Data.ProcTable; },
	render: function() {
		var Data = {ProcTable: this.state};
		var ps_rows = Data.ProcTable.List.map(function($proc) {
			return (React.DOM.tr( {key:$proc.PID}, React.DOM.td( {className:"digital"}, $proc.PID),React.DOM.td( {className:"digital"}, $proc.User),React.DOM.td( {className:"digital"}, $proc.Priority),React.DOM.td( {className:"digital"}, $proc.Nice),React.DOM.td( {className:"digital"}, $proc.Size),React.DOM.td( {className:"digital"}, $proc.Resident),React.DOM.td( {className:"center"}, $proc.Time),React.DOM.td(null, $proc.Name)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.PID.Href,      className:Data.ProcTable.Links.PID.Class}, "PID")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.User.Href,     className:Data.ProcTable.Links.User.Class}, "USER")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.Priority.Href, className:Data.ProcTable.Links.Priority.Class}, "PR")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.Nice.Href,     className:Data.ProcTable.Links.Nice.Class}, "NI")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.Size.Href,     className:Data.ProcTable.Links.Size.Class}, "VIRT")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.Resident.Href, className:Data.ProcTable.Links.Resident.Class}, "RES")),React.DOM.th( {className:"header center"},  " ", React.DOM.a( {href:Data.ProcTable.Links.Time.Href,     className:Data.ProcTable.Links.Time.Class}, "TIME")),React.DOM.th( {className:"header"},         "        ",        React.DOM.a( {href:Data.ProcTable.Links.Name.Href,     className:Data.ProcTable.Links.Name.Class}, "COMMAND")))),React.DOM.tbody( {id:"procrows"}, ps_rows)));
		
	}
});
