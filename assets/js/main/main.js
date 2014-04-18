var stateClass = React.createClass({
    getInitialState: function() { return {V: this.props.initialValue}; },
    render: function() {
	return (
	    React.DOM.span(null, this.state.V, this.props.append)
	);
    }
});
function xlabel_colorPercent(p) {
    return "label label-"+ _xcolorPercent(p);
}
function _xcolorPercent(p) {
    if (p > 90) { return "danger";  }
    if (p > 80) { return "warning"; }
    if (p > 20) { return "info";    }
    return "success";
}
var percentClass = React.createClass({
    getInitialState: function() { return {V: this.props.initialValue}; },
    render: function() {
	return (
	    React.DOM.span({className:xlabel_colorPercent(this.state.V)}, this.state.V, '%')
	);
    }
});

function newState(ID) {
    var node = document.getElementById(ID);
    n = new stateClass({
	elementID:    node,
	initialValue: node.innerHTML
    });
    React.renderComponent(n, n.props.elementID);
    return n;
}
function newPercent(ID) {
    var node = document.getElementById(ID);
    n = new percentClass({
	elementID:    node,
	initialValue: '',
    });
    // React.renderComponent(n, n.props.elementID);
    return n;
}

function default_button(btn) {
    btn.addClass('btn-default');
    btn.removeClass('btn-primary');
}
function primary_button(btn) {
    btn.addClass('btn-primary');
    btn.removeClass('btn-default');
}
function switch_class(label) {
    if (label.hasClass('btn-default')) {
	primary_button(label);
    } else {
	default_button(label);
    }
}

function update(Data) { // assert window["WebSocket"]
    var params = location.search.substr(1).split("&");
    for (var i in params) {
	if (params[i].split("=")[0] === "still") {
	    return;
	}
    }

    var   procTable =  procTableClass(null); // from gen/jscript.js
    var   diskTable =  diskTableClass(null); // from gen/jscript.js
    var inodesTable = inodeTableClass(null); // from gen/jscript.js
    var    cpuTable =   cpuTableClass(null); // from gen/jscript.js
    var    ifsTable =   ifsTableClass(null); // from gen/jscript.js
    var ifsPacketsTable = ifsPacketsTableClass(null); // from gen/jscript.js
    var ifsErrorsTable  = ifsErrorsTableClass(null);  // from gen/jscript.js

    React.renderComponent(  procTable, document.getElementById('ps-table'));
    React.renderComponent(  diskTable, document.getElementById('df-table'));
    React.renderComponent(inodesTable, document.getElementById('dfi-table'));
    React.renderComponent(   cpuTable, document.getElementById('cpu-table'));
    React.renderComponent(   ifsTable, document.getElementById('ifs-table'));
    React.renderComponent(ifsPacketsTable, document.getElementById('ifs-packets-table'));
    React.renderComponent(ifsErrorsTable,  document.getElementById('ifs-errors-table'));

    var onmessage = onmessage = function(event) {
	var data = JSON.parse(event.data);

	  procTable.setState(data.ProcTable);
	  diskTable.setState(data.DiskTable);
	inodesTable.setState(data.DiskTable);
	   cpuTable.setState(data.CPU);
	   ifsTable.setState(data.Interfaces);
	ifsPacketsTable.setState(data.Interfaces);
	ifsErrorsTable .setState(data.Interfaces);

	Data.About.Hostname  .setState({V: data.About.Hostname  });
	Data.About.IP        .setState({V: data.About.IP        });
	Data.System.Uptime   .setState({V: data.System.Uptime   });
	Data.System.LA       .setState({V: data.System.LA       });

	Data.RAM.Free        .setState({V: data.RAM.Free        });
	Data.RAM.Used        .setState({V: data.RAM.Used        });
	Data.RAM.Total       .setState({V: data.RAM.Total       });

	React.renderComponent(Data.RAM.UsePercent, Data.RAM.UsePercent.props.elementID);
	Data.RAM.UsePercent  .setState({V: data.RAM.UsePercent  });

	Data.Swap.Free       .setState({V: data.Swap.Free       });
	Data.Swap.Used       .setState({V: data.Swap.Used       });
	Data.Swap.Total      .setState({V: data.Swap.Total      });

    	React.renderComponent(Data.Swap.UsePercent, Data.Swap.UsePercent.props.elementID);
	Data.Swap.UsePercent .setState({V: data.Swap.UsePercent });

        $('span .tooltipable').tooltip(); // update the tooltips
    };

    var news = function() {
	var conn = new WebSocket("ws://" + HTTP_HOST + "/ws");
	var again = function() {
	    $("a.state").unbind('click');
	    window.setTimeout(news, 5000);
	};
	conn.onclose = again;
	conn.onerror = again;
	conn.onmessage = onmessage;

	conn.onopen = function() {
	    conn.send(location.search);
	    $(window).bind('popstate', function() {
		conn.send(location.search);
	    });
	};

	$("a.state").click(function() {
	    history.pushState({path: this.path}, '', this.href)
	    conn.send(this.search);
	    return false;
	});
    };
    news();
}

function ready() {
    $('label.all').on('click', function(e) {
	var label = $(this);
	switch_class(label);
	var parents = label.attr('data-parent').split(',');
	for (var i in parents) {
	    var parent = parents[i];

	    if ($(parent +'.collapse').not('.in').length > 0) {
		var labels = $($('label[href="' + parent +'"]').not(this));
		if (labels.length > 0) { // show last active panel
		    $($('label[href="' + parent +'"].active').attr('href')).collapse('show');
		} else {
		    $(parent).collapse('show');
		}
		/*
		// console.log('label[href="' + parent +'"].active');
		// console.log($('label[href="' + parent +'"].active').attr('href'));
		var active = $('label[href="' + parent +'"].active');
		console.log('active', active.length, 'label[href="' + parent +'"].active');
		if (active.length > 0) {
		    $(active.attr('href')).collapse('show');
		} else {
		    console.log('here');
		    $($('label[href="' + parent +'"]').attr('href')).collapse('show');
		} // */
	    }

	    var collapsing;
	    collapsing = $(parent +' tr.collapse').not('.in').length > 0;
	    if (collapsing) {
		$(parent +' tr.collapse').collapse('show');
	    } else {
		$(parent +' tr.collapse').collapse('hide');
	    }
	    // collapsing = $(parent +' tr.collapse').not('.in').length > 0; // update
	    collapsing = !collapsing;
	    // setting defined in jscript.{amber,js} globals
	         if (parent == '#network-acc') { collapsing_interfaces = collapsing; }
	    else if (parent == '#disk-acc')    { collapsing_disks      = collapsing; }
	    else if (parent == '#cpu')         { collapsing_cpu        = collapsing; }
	}
    });

    /* $('label.all').tooltip({
	container: 'body',
	placement: 'left',
	trigger: 'click'
    });
    $('label.all').tooltip('show');
    $(window).resize(function() {
	$('label.all').tooltip('show');
    });
    $('.tooltip .tooltip-arrow').addClass('ii-tooltip-arrow');
    $('.tooltip .tooltip-inner').addClass('ii-tooltip-inner'); // */

    $('span .tooltipable').tooltip();

    var accs = {
	'#disk-acc': {
	    none: '#disks',
	    label: 'label.disk-switch'
	},
	'#network-acc': {
	    none: '#interfaces',
	    label: 'label.network-switch'
	}
    };
    for (acc in accs) {
	var prop = accs[acc];
	$(acc +' .panel-collapse').each(function(_i, panel) {
	    $(panel).collapse({toggle: false, parent: acc}); // init objects
	});
	accbind(acc, prop);
    }

    update({
	About: { Hostname:   newState('Data.About.Hostname')
		 , IP:       newState('Data.About.IP')
	       },
	System: { Uptime:    newState('Data.System.Uptime')
		  , LA:      newState('Data.System.LA')
		},
	RAM: { Free:         newState('Data.RAM.Free')
	       , Used:       newState('Data.RAM.Used')
	       , UsePercent: newPercent('Data.RAM.UsePercent')
	       , Total:      newState('Data.RAM.Total')
	     },
	Swap: { Free:         newState('Data.Swap.Free')
		, Used:       newState('Data.Swap.Used')
		, UsePercent: newPercent('Data.Swap.UsePercent')
		, Total:      newState('Data.Swap.Total')
	      }
    });
}

function accbind(acc, prop) {
    $('[data-parent="'+ acc +'"]').on('click', function(e) {
	if ($(this).hasClass('nondefault')) {
	    switch_class($(this));
	}
	var tagName = $(this).prop('tagName');
	if (tagName == 'LABEL') {
	    $('[data-parent="'+ acc +'"].nondefault').not(this).each(function (_i, label) {
		default_button($(label));
	    });
	}

	var href = $(this).attr('href');
	if (tagName == 'LABEL') { // acc in Object.keys(accs)
	    if ($(href).hasClass('in')) {
		// https://stackoverflow.com/q/15725717
		e.stopPropagation();
		return;
	    }
	    if ($(acc +' .panel-collapse.in').not(prop.none).length == 0) {
		$($('label[href="' + parent +'"].active').attr('href')).collapse('show');
		// $($(prop.label +'.active').attr('href')).collapse('show');
	    }
	} else if ($(href).hasClass('in')) {
	    // show LAST collapsed
	    $($(prop.label +'.active').attr('href')).collapse('show');
	    // primary_button($(prop.label +'.active.nondefault'));
	    e.stopPropagation();
	    e.preventDefault();
	}
    });
}
