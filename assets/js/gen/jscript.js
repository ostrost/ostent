/** @jsx React.DOM */ var collapsing_interfaces = true; // default is true
var ifsTableClass = React.createClass({displayName: 'ifsTableClass',
	getInitialState: function() { return Data.Interfaces; },
	render: function() {
		var Data = {Interfaces: this.state};
		var ifs_rows = Data.Interfaces.List.map(function($if) {
			if (!collapsing_interfaces) {
				$if.CollapseClass = "collapse in";
			}
			return (React.DOM.tr( {key:$if.NameKey, className:$if.CollapseClass}, React.DOM.td(null, React.DOM.span( {dangerouslySetInnerHTML:{__html: $if.NameHTML}} )),React.DOM.td( {className:"digital"}, $if.DeltaInBytes),React.DOM.td( {className:"digital"}, $if.DeltaOutBytes),React.DOM.td( {className:"digital"}, $if.InBytes),React.DOM.td( {className:"digital"}, $if.OutBytes)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null, "Interface"),React.DOM.th( {className:"digital nobr"}, "In",React.DOM.span( {className:"unit"}, React.DOM.b(null, "b"),"ps")),React.DOM.th( {className:"digital nobr"}, "Out",React.DOM.span( {className:"unit"}, React.DOM.b(null, "b"),"ps")),React.DOM.th( {className:"digital nobr"}, "In",React.DOM.span( {className:"unit"}, "%4G")),React.DOM.th( {className:"digital nobr"}, "Out",React.DOM.span( {className:"unit"}, "%4G")))),React.DOM.tbody(null, ifs_rows)));
		
	}
});

var ifsPacketsTableClass = React.createClass({displayName: 'ifsPacketsTableClass',
	getInitialState: function() { return Data.Interfaces; },
	render: function() {
		var Data = {Interfaces: this.state};
		var ifs_rows = Data.Interfaces.List.map(function($if) {
			if (!collapsing_interfaces) {
				$if.CollapseClass = "collapse in";
			}
			return (React.DOM.tr( {key:$if.NameKey, className:$if.CollapseClass}, React.DOM.td(null, React.DOM.span( {dangerouslySetInnerHTML:{__html: $if.NameHTML}} )),React.DOM.td( {className:"digital"}, $if.DeltaInPackets),React.DOM.td( {className:"digital"}, $if.DeltaOutPackets),React.DOM.td( {className:"digital"}, $if.InPackets),React.DOM.td( {className:"digital"}, $if.OutPackets)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null, "Interface"),React.DOM.th( {className:"digital nobr"}, "In ",React.DOM.span( {className:"unit"}, "ps")),React.DOM.th( {className:"digital nobr"}, "Out ",React.DOM.span( {className:"unit"}, "ps")),React.DOM.th( {className:"digital nobr"}, "In",React.DOM.span( {className:"unit"}, "%4G")),React.DOM.th( {className:"digital nobr"}, "Out",React.DOM.span( {className:"unit"}, "%4G")))),React.DOM.tbody(null, ifs_rows)));
		
	}
});

var ifsErrorsTableClass = React.createClass({displayName: 'ifsErrorsTableClass',
	getInitialState: function() { return Data.Interfaces; },
	render: function() {
		var Data = {Interfaces: this.state};
		var ifs_rows = Data.Interfaces.List.map(function($if) {
			if (!collapsing_interfaces) {
				$if.CollapseClass = "collapse in";
			}
			return (React.DOM.tr( {key:$if.NameKey, className:$if.CollapseClass}, React.DOM.td(null, React.DOM.span( {dangerouslySetInnerHTML:{__html: $if.NameHTML}} )),React.DOM.td( {className:"digital"}, $if.DeltaInErrors),React.DOM.td( {className:"digital"}, $if.DeltaOutErrors),React.DOM.td( {className:"digital"}, $if.InErrors),React.DOM.td( {className:"digital"}, $if.OutErrors)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null, "Interface"),React.DOM.th( {className:"digital nobr"}, "In ",React.DOM.span( {className:"unit"}, "ps")),React.DOM.th( {className:"digital nobr"}, "Out ",React.DOM.span( {className:"unit"}, "ps")),React.DOM.th( {className:"digital nobr"}, "In",React.DOM.span( {className:"unit"}, "%4G")),React.DOM.th( {className:"digital nobr"}, "Out",React.DOM.span( {className:"unit"}, "%4G")))),React.DOM.tbody(null, ifs_rows)));
		
	}
});

var collapsing_disks = true; // default is true
var diskTableClass = React.createClass({displayName: 'diskTableClass',
	getInitialState: function() { return Data.DiskTable; },
	render: function() {
		var Data = {DiskTable: this.state};
		var df_rows = Data.DiskTable.List.map(function($disk) {
			if (!collapsing_disks) {
				$disk.CollapseClass = "collapse in";
			}
			return (React.DOM.tr( {key:$disk.DiskNameKey, className:$disk.CollapseClass}, React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $disk.DiskNameHTML}} )),React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $disk.DirNameHTML}} )),React.DOM.td( {className:"digital"}, $disk.Avail),React.DOM.td( {className:"digital"}, $disk.Used," ",React.DOM.sup(null, React.DOM.span( {className:$disk.UsePercentClass}, $disk.UsePercent,"%"))),React.DOM.td( {className:"digital"}, $disk.Total)));
		});

		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th( {className:"header"},         "        ",        React.DOM.a( {href:Data.DiskTable.Links.DiskName.Href, className:Data.DiskTable.Links.DiskName.Class}, "Device")),React.DOM.th( {className:"header"},         "        ",        React.DOM.a( {href:Data.DiskTable.Links.DirName.Href,  className:Data.DiskTable.Links.DirName.Class}, "Mounted")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.DiskTable.Links.Avail.Href,    className:Data.DiskTable.Links.Avail.Class}, "Avail")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.DiskTable.Links.Used.Href,     className:Data.DiskTable.Links.Used.Class}, "Used")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.DiskTable.Links.Total.Href,    className:Data.DiskTable.Links.Total.Class}, "Total")))),React.DOM.tbody(null, df_rows)));
		
	}
});

var inodeTableClass = React.createClass({displayName: 'inodeTableClass',
	getInitialState: function() { return Data.DiskTable; },
	render: function() {
		var Data = {DiskTable: this.state};
		var di_rows = Data.DiskTable.List.map(function($disk) {
			if (!collapsing_disks) {
				$disk.CollapseClass = "collapse in";
			}
			return (React.DOM.tr( {key:$disk.DiskNameKey, className:$disk.CollapseClass}, React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $disk.DiskNameHTML}} )),React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $disk.DirNameHTML}} )),React.DOM.td( {className:"digital"}, $disk.Ifree),React.DOM.td( {className:"digital"}, $disk.Iused," ",React.DOM.sup(null, React.DOM.span( {className:$disk.IusePercentClass}, $disk.IusePercent,"%"))),React.DOM.td( {className:"digital"}, $disk.Inodes)));
		});

		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th( {className:"header"}, "Device"),React.DOM.th( {className:"header"}, "Mounted"),React.DOM.th( {className:"header digital"}, "Avail"),React.DOM.th( {className:"header digital"}, "Used"),React.DOM.th( {className:"header digital"}, "Total"))),React.DOM.tbody(null, di_rows)));
		
	}
});

var collapsing_cpu = true; // default is true
var cpuTableClass = React.createClass({displayName: 'cpuTableClass',
	getInitialState: function() { return Data.CPU; },
	render: function() {
		var Data = {CPU: this.state};
		var cpu_rows = Data.CPU.List.map(function($core) {
			if (!collapsing_cpu) {
				$core.CollapseClass = "collapse in";
			}
			return (React.DOM.tr( {key:$core.N, className:$core.CollapseClass}, React.DOM.td( {className:"digital nobr"}, $core.N),React.DOM.td( {className:"digital"}, React.DOM.span( {id:"core0.User", className:$core.UserClass}, $core.User)),React.DOM.td( {className:"digital"}, React.DOM.span( {id:"core0.Sys",  className:$core.SysClass}, $core.Sys)),React.DOM.td( {className:"digital"}, React.DOM.span( {id:"core0.Idle", className:$core.IdleClass}, $core.Idle))));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null),React.DOM.th( {className:"digital nobr"}, "User",React.DOM.span( {className:"unit"}, "%")),React.DOM.th( {className:"digital nobr"}, "Sys",React.DOM.span( {className:"unit"}, "%")),React.DOM.th( {className:"digital nobr"}, "Idle",React.DOM.span( {className:"unit"}, "%")))),React.DOM.tbody(null, cpu_rows)));
		
	}
});

var procTableClass  = React.createClass({displayName: 'procTableClass',
	getInitialState: function() { return Data.ProcTable; },
	render: function() {
		var Data = {ProcTable: this.state};
		var ps_rows = Data.ProcTable.List.map(function($proc) {
			return (React.DOM.tr( {key:$proc.PID}, React.DOM.td( {className:"digital"}, $proc.PID),React.DOM.td( {className:"digital"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $proc.UserHTML}} )),React.DOM.td( {className:"digital"}, $proc.Priority),React.DOM.td( {className:"digital"}, $proc.Nice),React.DOM.td( {className:"digital"}, $proc.Size),React.DOM.td( {className:"digital"}, $proc.Resident),React.DOM.td( {className:"center"}, $proc.Time),React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $proc.NameHTML}} ))));
		});
		
		return (React.DOM.table( {className:"table2 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.PID.Href,      className:Data.ProcTable.Links.PID.Class}, "PID")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.User.Href,     className:Data.ProcTable.Links.User.Class}, "USER")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.Priority.Href, className:Data.ProcTable.Links.Priority.Class}, "PR")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.Nice.Href,     className:Data.ProcTable.Links.Nice.Class}, "NI")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.Size.Href,     className:Data.ProcTable.Links.Size.Class}, "VIRT")),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.Resident.Href, className:Data.ProcTable.Links.Resident.Class}, "RES")),React.DOM.th( {className:"header center"},  " ", React.DOM.a( {href:Data.ProcTable.Links.Time.Href,     className:Data.ProcTable.Links.Time.Class}, "TIME")),React.DOM.th( {className:"header"},         "        ",        React.DOM.a( {href:Data.ProcTable.Links.Name.Href,     className:Data.ProcTable.Links.Name.Class}, "COMMAND")))),React.DOM.tbody( {id:"procrows"}, ps_rows)));
		
	}
});
