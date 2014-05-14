/** @jsx React.DOM */ function emptyK(obj, key) {
return (obj      === undefined ||
	obj      === null      ||
	obj[key] === undefined ||
	obj[key] === null);
}

var ifsTableClass = React.createClass({displayName: 'ifsTableClass',
	getInitialState: function() { return Data.InterfacesBytes; },

	render: function() {
		var Data = {InterfacesBytes: this.state};
		var ifs_rows = emptyK(Data.InterfacesBytes, 'List') ?'': Data.InterfacesBytes.List.map(function($if) {
			return (React.DOM.tr( {key:$if.NameKey}, React.DOM.td(null, React.DOM.span( {dangerouslySetInnerHTML:{__html: $if.NameHTML}} )),React.DOM.td( {className:"digital"}, $if.DeltaIn),React.DOM.td( {className:"digital"}, $if.DeltaOut),React.DOM.td( {className:"digital"}, $if.In),React.DOM.td( {className:"digital"}, $if.Out)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null, "Interface"),React.DOM.th( {className:"digital nobr"}, "In",React.DOM.span( {className:"unit", title:"bits per second"}, React.DOM.b(null, "b"),"ps")),React.DOM.th( {className:"digital nobr"}, "Out",React.DOM.span( {className:"unit", title:"bits per second"}, React.DOM.b(null, "b"),"ps")),React.DOM.th( {className:"digital nobr"}, "In",React.DOM.span( {className:"unit", title:"total modulo 4G"}, "%4G")),React.DOM.th( {className:"digital nobr"}, "Out",React.DOM.span( {className:"unit", title:"total modulo 4G"}, "%4G")))),React.DOM.tbody(null, ifs_rows)));
		
	}
});

var ifsErrorsTableClass = React.createClass({displayName: 'ifsErrorsTableClass',
	getInitialState: function() { return Data.InterfacesErrors; },

	render: function() {
		var Data = {InterfacesErrors: this.state};
		var ifs_rows = emptyK(Data.InterfacesErrors, 'List') ?'': Data.InterfacesErrors.List.map(function($if) {
			return (React.DOM.tr( {key:$if.NameKey}, React.DOM.td(null, React.DOM.span( {dangerouslySetInnerHTML:{__html: $if.NameHTML}} )),React.DOM.td( {className:"digital"}, $if.DeltaIn),React.DOM.td( {className:"digital"}, $if.DeltaOut),React.DOM.td( {className:"digital"}, $if.In),React.DOM.td( {className:"digital"}, $if.Out)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null, "Interface"),React.DOM.th( {className:"digital nobr"}, "In ",React.DOM.span( {className:"unit", title:"per second"}, "ps")),React.DOM.th( {className:"digital nobr"}, "Out ",React.DOM.span( {className:"unit", title:"per second"}, "ps")),React.DOM.th( {className:"digital nobr"}, "In",React.DOM.span( {className:"unit", title:"modulo 4G"}, "%4G")),React.DOM.th( {className:"digital nobr"}, "Out",React.DOM.span( {className:"unit", title:"modulo 4G"}, "%4G")))),React.DOM.tbody(null, ifs_rows)));
		
	}
});

var ifsPacketsTableClass = React.createClass({displayName: 'ifsPacketsTableClass',
	getInitialState: function() { return Data.InterfacesPackets; },

	render: function() {
		var Data = {InterfacesPackets: this.state};
		var ifs_rows = emptyK(Data.InterfacesPackets, 'List') ?'': Data.InterfacesPackets.List.map(function($if) {
			return (React.DOM.tr( {key:$if.NameKey}, React.DOM.td(null, React.DOM.span( {dangerouslySetInnerHTML:{__html: $if.NameHTML}} )),React.DOM.td( {className:"digital"}, $if.DeltaIn),React.DOM.td( {className:"digital"}, $if.DeltaOut),React.DOM.td( {className:"digital"}, $if.In),React.DOM.td( {className:"digital"}, $if.Out)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null, "Interface"),React.DOM.th( {className:"digital nobr"}, "In ",React.DOM.span( {className:"unit", title:"per second"}, "ps")),React.DOM.th( {className:"digital nobr"}, "Out ",React.DOM.span( {className:"unit", title:"per second"}, "ps")),React.DOM.th( {className:"digital nobr"}, "In",React.DOM.span( {className:"unit", title:"modulo 4G"}, "%4G")),React.DOM.th( {className:"digital nobr"}, "Out",React.DOM.span( {className:"unit", title:"modulo 4G"}, "%4G")))),React.DOM.tbody(null, ifs_rows)));
		
	}
});

var disksinBytesClass = React.createClass({displayName: 'disksinBytesClass',
	getInitialState: function() { return {DiskLinks: Data.DiskLinks, DisksinBytes: Data.DisksinBytes}; },

	render: function() {
		var Data = this.state;
		var rows = emptyK(Data.DisksinBytes, 'List') ?'': Data.DisksinBytes.List.map(function($disk) {
			return (React.DOM.tr( {key:$disk.DirNameKey}, React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $disk.DiskNameHTML}} )),React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $disk.DirNameHTML}} )),React.DOM.td( {className:"digital"}, $disk.Avail),React.DOM.td( {className:"digital"}, $disk.Used," ",React.DOM.sup(null, React.DOM.span( {className:$disk.UsePercentClass}, $disk.UsePercent,"%"))),React.DOM.td( {className:"digital"}, $disk.Total)));
		});

		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th( {className:"header"},         "        ",        React.DOM.a( {href:Data.DiskLinks.DiskName.Href, className:Data.DiskLinks.DiskName.Class}, "Device",React.DOM.span(  {className:Data.DiskLinks.DiskName.CaretClass} ))),React.DOM.th( {className:"header"},         "        ",        React.DOM.a( {href:Data.DiskLinks.DirName.Href,  className:Data.DiskLinks.DirName.Class} , "Mounted",React.DOM.span( {className:Data.DiskLinks.DirName.CaretClass}  ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.DiskLinks.Avail.Href,    className:Data.DiskLinks.Avail.Class}   , "Avail",React.DOM.span(   {className:Data.DiskLinks.Avail.CaretClass}    ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.DiskLinks.Used.Href,     className:Data.DiskLinks.Used.Class}    , "Used",React.DOM.span(    {className:Data.DiskLinks.Used.CaretClass}     ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.DiskLinks.Total.Href,    className:Data.DiskLinks.Total.Class}   , "Total",React.DOM.span(   {className:Data.DiskLinks.Total.CaretClass}    ))))),React.DOM.tbody(null, rows)));
		
	}
});

var disksinInodesClass = React.createClass({displayName: 'disksinInodesClass',
	getInitialState: function() { return {DiskLinks: Data.DiskLinks, DisksinInodes: Data.DisksinInodes}; },

	render: function() {
		var Data = this.state;
		var rows = emptyK(Data.DisksinInodes, 'List') ?'': Data.DisksinInodes.List.map(function($disk) {
			return (React.DOM.tr( {key:$disk.DirNameKey}, React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $disk.DiskNameHTML}} )),React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $disk.DirNameHTML}} )),React.DOM.td( {className:"digital"}, $disk.Ifree),React.DOM.td( {className:"digital"}, $disk.Iused," ",React.DOM.sup(null, React.DOM.span( {className:$disk.IusePercentClass}, $disk.IusePercent,"%"))),React.DOM.td( {className:"digital"}, $disk.Inodes)));
		});

		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th( {className:"header"}, "Device"),React.DOM.th( {className:"header"}, "Mounted"),React.DOM.th( {className:"header digital"}, "Avail"),React.DOM.th( {className:"header digital"}, "Used"),React.DOM.th( {className:"header digital"}, "Total"))),React.DOM.tbody(null, rows)));
		
	}
});

var cpuTableClass = React.createClass({displayName: 'cpuTableClass',
	getInitialState: function() { return Data.CPU; },

	render: function() {
		var Data = {CPU: this.state};
		var cpu_rows = Data.CPU.List.map(function($core) {
			return (React.DOM.tr( {key:$core.N}, React.DOM.td( {className:"digital nobr"}, $core.N),React.DOM.td( {className:"digital"}, React.DOM.span( {id:"core0.User", className:$core.UserClass}, $core.User)),React.DOM.td( {className:"digital"}, React.DOM.span( {id:"core0.Sys",  className:$core.SysClass}, $core.Sys)),React.DOM.td( {className:"digital"}, React.DOM.span( {id:"core0.Idle", className:$core.IdleClass}, $core.Idle))));
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
		
		return (React.DOM.table( {className:"table2 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.PID.Href,      className:Data.ProcTable.Links.PID.Class}     , "PID",React.DOM.span(     {className:Data.ProcTable.Links.PID.CaretClass}      ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.User.Href,     className:Data.ProcTable.Links.User.Class}    , "USER",React.DOM.span(    {className:Data.ProcTable.Links.User.CaretClass}     ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.Priority.Href, className:Data.ProcTable.Links.Priority.Class}, "PR",React.DOM.span(      {className:Data.ProcTable.Links.Priority.CaretClass} ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.Nice.Href,     className:Data.ProcTable.Links.Nice.Class}    , "NI",React.DOM.span(      {className:Data.ProcTable.Links.Nice.CaretClass}     ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.Size.Href,     className:Data.ProcTable.Links.Size.Class}    , "VIRT",React.DOM.span(    {className:Data.ProcTable.Links.Size.CaretClass}     ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.ProcTable.Links.Resident.Href, className:Data.ProcTable.Links.Resident.Class}, "RES",React.DOM.span(     {className:Data.ProcTable.Links.Resident.CaretClass} ))),React.DOM.th( {className:"header center"},  " ", React.DOM.a( {href:Data.ProcTable.Links.Time.Href,     className:Data.ProcTable.Links.Time.Class}    , "TIME",React.DOM.span(    {className:Data.ProcTable.Links.Time.CaretClass}     ))),React.DOM.th( {className:"header"},         "        ",        React.DOM.a( {href:Data.ProcTable.Links.Name.Href,     className:Data.ProcTable.Links.Name.Class}    , "COMMAND",React.DOM.span( {className:Data.ProcTable.Links.Name.CaretClass}     ))))),React.DOM.tbody( {id:"procrows"}, ps_rows)));
		
	}
});
