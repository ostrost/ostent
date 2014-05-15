/** @jsx React.DOM */ function emptyK(obj, key) {
return (obj      === undefined ||
	obj      === null      ||
	obj[key] === undefined ||
	obj[key] === null);
}

var IFbytesCLASS = React.createClass({displayName: 'IFbytesCLASS',
	getInitialState: function() { return Data.IFbytes; },

	render: function() {
		var Data = {IFbytes: this.state};
		var rows = emptyK(Data.IFbytes, 'List') ?'': Data.IFbytes.List.map(function($if) {
			return (React.DOM.tr( {key:$if.NameKey}, React.DOM.td(null, React.DOM.span( {dangerouslySetInnerHTML:{__html: $if.NameHTML}} )),React.DOM.td( {className:"digital"}, $if.DeltaIn),React.DOM.td( {className:"digital"}, $if.DeltaOut),React.DOM.td( {className:"digital"}, $if.In),React.DOM.td( {className:"digital"}, $if.Out)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null, "Interface"),React.DOM.th( {className:"digital nobr", title:"BITS per second"}, "In",React.DOM.span( {className:"unit"}, React.DOM.b(null, "b"),"ps")),React.DOM.th( {className:"digital nobr", title:"BITS per second"}, "Out",React.DOM.span( {className:"unit"}, React.DOM.b(null, "b"),"ps")),React.DOM.th( {className:"digital nobr", title:"total BYTES modulo 4G"}, "In",React.DOM.span( {className:"unit"}, React.DOM.b(null, "B"),"%4G")),React.DOM.th( {className:"digital nobr", title:"total BYTES modulo 4G"}, "Out",React.DOM.span( {className:"unit"}, React.DOM.b(null, "B"),"%4G")))),React.DOM.tbody(null, rows)));
		
	}
});

var IFerrorsCLASS = React.createClass({displayName: 'IFerrorsCLASS',
	getInitialState: function() { return Data.IFerrors; },

	render: function() {
		var Data = {IFerrors: this.state};
		var rows = emptyK(Data.IFerrors, 'List') ?'': Data.IFerrors.List.map(function($if) {
			return (React.DOM.tr( {key:$if.NameKey}, React.DOM.td(null, React.DOM.span( {dangerouslySetInnerHTML:{__html: $if.NameHTML}} )),React.DOM.td( {className:"digital"}, $if.DeltaIn),React.DOM.td( {className:"digital"}, $if.DeltaOut),React.DOM.td( {className:"digital"}, $if.In),React.DOM.td( {className:"digital"}, $if.Out)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null, "Interface"),React.DOM.th( {className:"digital nobr", title:"per second"}, "In ",React.DOM.span( {className:"unit"}, "ps")),React.DOM.th( {className:"digital nobr", title:"per second"}, "Out ",React.DOM.span( {className:"unit"}, "ps")),React.DOM.th( {className:"digital nobr", title:"modulo 4G"}, "In ",React.DOM.span( {className:"unit"}, "%4G")),React.DOM.th( {className:"digital nobr", title:"modulo 4G"}, "Out ",React.DOM.span( {className:"unit"}, "%4G")))),React.DOM.tbody(null, rows)));
		
	}
});

var IFpacketsCLASS = React.createClass({displayName: 'IFpacketsCLASS',
	getInitialState: function() { return Data.IFpackets; },

	render: function() {
		var Data = {IFpackets: this.state};
		var rows = emptyK(Data.IFpackets, 'List') ?'': Data.IFpackets.List.map(function($if) {
			return (React.DOM.tr( {key:$if.NameKey}, React.DOM.td(null, React.DOM.span( {dangerouslySetInnerHTML:{__html: $if.NameHTML}} )),React.DOM.td( {className:"digital"}, $if.DeltaIn),React.DOM.td( {className:"digital"}, $if.DeltaOut),React.DOM.td( {className:"digital"}, $if.In),React.DOM.td( {className:"digital"}, $if.Out)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null, "Interface"),React.DOM.th( {className:"digital nobr", title:"per second"}, "In ",React.DOM.span( {className:"unit"}, "ps")),React.DOM.th( {className:"digital nobr", title:"per second"}, "Out ",React.DOM.span( {className:"unit"}, "ps")),React.DOM.th( {className:"digital nobr", title:"total modulo 4G"}, "In ",React.DOM.span( {className:"unit"}, "%4G")),React.DOM.th( {className:"digital nobr", title:"total modulo 4G"}, "Out ",React.DOM.span( {className:"unit"}, "%4G")))),React.DOM.tbody(null, rows)));
		
	}
});

var DFbytesCLASS = React.createClass({displayName: 'DFbytesCLASS',
	getInitialState: function() { return {DFlinks: Data.DFlinks, DFbytes: Data.DFbytes}; },

	render: function() {
		var Data = this.state;
		var rows = emptyK(Data.DFbytes, 'List') ?'': Data.DFbytes.List.map(function($disk) {
			return (React.DOM.tr( {key:$disk.DirNameKey}, React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $disk.DiskNameHTML}} )),React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $disk.DirNameHTML}} )),React.DOM.td( {className:"digital"}, $disk.Avail),React.DOM.td( {className:"digital"}, $disk.Used," ",React.DOM.sup(null, React.DOM.span( {className:$disk.UsePercentClass}, $disk.UsePercent,"%"))),React.DOM.td( {className:"digital"}, $disk.Total)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th( {className:"header"},         "        ",        React.DOM.a( {href:Data.DFlinks.DiskName.Href, className:Data.DFlinks.DiskName.Class}, "Device",React.DOM.span(  {className:Data.DFlinks.DiskName.CaretClass} ))),React.DOM.th( {className:"header"},         "        ",        React.DOM.a( {href:Data.DFlinks.DirName.Href,  className:Data.DFlinks.DirName.Class} , "Mounted",React.DOM.span( {className:Data.DFlinks.DirName.CaretClass}  ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.DFlinks.Avail.Href,    className:Data.DFlinks.Avail.Class}   , "Avail",React.DOM.span(   {className:Data.DFlinks.Avail.CaretClass}    ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.DFlinks.Used.Href,     className:Data.DFlinks.Used.Class}    , "Used",React.DOM.span(    {className:Data.DFlinks.Used.CaretClass}     ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.DFlinks.Total.Href,    className:Data.DFlinks.Total.Class}   , "Total",React.DOM.span(   {className:Data.DFlinks.Total.CaretClass}    ))))),React.DOM.tbody(null, rows)));
		
	}
});

var DFinodesCLASS = React.createClass({displayName: 'DFinodesCLASS',
	getInitialState: function() { return {DFlinks: Data.DFlinks, DFinodes: Data.DFinodes}; },

	render: function() {
		var Data = this.state;
		var rows = emptyK(Data.DFinodes, 'List') ?'': Data.DFinodes.List.map(function($disk) {
			return (React.DOM.tr( {key:$disk.DirNameKey}, React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $disk.DiskNameHTML}} )),React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $disk.DirNameHTML}} )),React.DOM.td( {className:"digital"}, $disk.Ifree),React.DOM.td( {className:"digital"}, $disk.Iused," ",React.DOM.sup(null, React.DOM.span( {className:$disk.IusePercentClass}, $disk.IusePercent,"%"))),React.DOM.td( {className:"digital"}, $disk.Inodes)));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th( {className:"header"}, "Device"),React.DOM.th( {className:"header"}, "Mounted"),React.DOM.th( {className:"header digital"}, "Avail"),React.DOM.th( {className:"header digital"}, "Used"),React.DOM.th( {className:"header digital"}, "Total"))),React.DOM.tbody(null, rows)));
		
	}
});

var CPUtableCLASS = React.createClass({displayName: 'CPUtableCLASS',
	getInitialState: function() { return Data.CPU; },

	render: function() {
		var Data = {CPU: this.state};
		var rows = Data.CPU.List.map(function($core) {
			return (React.DOM.tr( {key:$core.N}, React.DOM.td( {className:"digital nobr"}, $core.N),React.DOM.td( {className:"digital"}, React.DOM.span( {id:"core0.User", className:$core.UserClass}, $core.User)),React.DOM.td( {className:"digital"}, React.DOM.span( {id:"core0.Sys",  className:$core.SysClass}, $core.Sys)),React.DOM.td( {className:"digital"}, React.DOM.span( {id:"core0.Idle", className:$core.IdleClass}, $core.Idle))));
		});
		
		return (React.DOM.table( {className:"table1 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th(null),React.DOM.th( {className:"digital nobr"}, "User",React.DOM.span( {className:"unit"}, "%")),React.DOM.th( {className:"digital nobr"}, "Sys",React.DOM.span( {className:"unit"}, "%")),React.DOM.th( {className:"digital nobr"}, "Idle",React.DOM.span( {className:"unit"}, "%")))),React.DOM.tbody(null, rows)));
		
	}
});

var PStableCLASS = React.createClass({displayName: 'PStableCLASS',
	getInitialState: function() { return Data.PStable; },

	render: function() {
		var Data = {PStable: this.state};
		var rows = Data.PStable.List.map(function($proc) {
			return (React.DOM.tr( {key:$proc.PID}, React.DOM.td( {className:"digital"}, $proc.PID),React.DOM.td( {className:"digital"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $proc.UserHTML}} )),React.DOM.td( {className:"digital"}, $proc.Priority),React.DOM.td( {className:"digital"}, $proc.Nice),React.DOM.td( {className:"digital"}, $proc.Size),React.DOM.td( {className:"digital"}, $proc.Resident),React.DOM.td( {className:"center"}, $proc.Time),React.DOM.td( {className:"nobr"}, React.DOM.span( {dangerouslySetInnerHTML:{__html: $proc.NameHTML}} ))));
		});
		
		return (React.DOM.table( {className:"table2 stripe-table"}, React.DOM.thead(null, React.DOM.tr(null, React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.PStable.Links.PID.Href,      className:Data.PStable.Links.PID.Class}     , "PID",React.DOM.span(     {className:Data.PStable.Links.PID.CaretClass}      ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.PStable.Links.User.Href,     className:Data.PStable.Links.User.Class}    , "USER",React.DOM.span(    {className:Data.PStable.Links.User.CaretClass}     ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.PStable.Links.Priority.Href, className:Data.PStable.Links.Priority.Class}, "PR",React.DOM.span(      {className:Data.PStable.Links.Priority.CaretClass} ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.PStable.Links.Nice.Href,     className:Data.PStable.Links.Nice.Class}    , "NI",React.DOM.span(      {className:Data.PStable.Links.Nice.CaretClass}     ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.PStable.Links.Size.Href,     className:Data.PStable.Links.Size.Class}    , "VIRT",React.DOM.span(    {className:Data.PStable.Links.Size.CaretClass}     ))),React.DOM.th( {className:"header digital"}, React.DOM.a( {href:Data.PStable.Links.Resident.Href, className:Data.PStable.Links.Resident.Class}, "RES",React.DOM.span(     {className:Data.PStable.Links.Resident.CaretClass} ))),React.DOM.th( {className:"header center"},  " ", React.DOM.a( {href:Data.PStable.Links.Time.Href,     className:Data.PStable.Links.Time.Class}    , "TIME",React.DOM.span(    {className:Data.PStable.Links.Time.CaretClass}     ))),React.DOM.th( {className:"header"},         "        ",        React.DOM.a( {href:Data.PStable.Links.Name.Href,     className:Data.PStable.Links.Name.Class}    , "COMMAND",React.DOM.span( {className:Data.PStable.Links.Name.CaretClass}     ))))),React.DOM.tbody(null, rows)));
		
	}
});
