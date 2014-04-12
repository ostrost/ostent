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

function update(Data) { // assert window["WebSocket"]
    $('span .tooltipable').tooltip();

    // https://stackoverflow.com/q/15725717
    $('[data-parent="#disk-accordion"]').on('click',function(e) {
	if ($($(this).attr('href')).hasClass('in')){
            e.stopPropagation();
	}
    });

    var params = location.search.substr(1).split("&");
    for (var i in params) {
	if (params[i].split("=")[0] === "nojs") {
	    return;
	}
    }

    var   procTable =  procTableClass(null); // from gen/jscript.js
    var   diskTable =  diskTableClass(null); // from gen/jscript.js
    var inodesTable = inodeTableClass(null); // from gen/jscript.js
    var    cpuTable =   cpuTableClass(null); // from gen/jscript.js
    var    ifsTable =   ifsTableClass(null); // from gen/jscript.js

    React.renderComponent(  procTable, document.getElementById('ps-table'));
    React.renderComponent(  diskTable, document.getElementById('df-table'));
    React.renderComponent(inodesTable, document.getElementById('dfi-table'));
    React.renderComponent(   cpuTable, document.getElementById('cpu-table'));
    React.renderComponent(   ifsTable, document.getElementById('ifs-table'));

    var onmessage = onmessage = function(event) {
	var data = JSON.parse(event.data);

	  procTable.setState(data.ProcTable);
	  diskTable.setState(data.DiskTable);
	inodesTable.setState(data.DiskTable);
	   cpuTable.setState(data.CPU);
	   ifsTable.setState(data.Interfaces);

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
