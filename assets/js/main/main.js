// requires underscore.js
// requires backbone.js
// requires jquery.js

function newwebsocket(onmessage) {
    var conn, connected = false;
    function sendJSON(obj) {
	if (conn && ~[2,3].indexOf(conn.readyState)) {
            connected = false;
            init();
	}
	if (!connected) {
            console.log('Cannot send this (re-connection failed):', obj);
            return;
	}
	conn.send(JSON.stringify(obj));
    }
    function sendClient(dict) {
	console.log(JSON.stringify(dict), 'sendClient');
	return sendJSON({Client: dict});
    }
    function sendSearch(search) {
	return sendJSON({Search: search});
    }

    function init() {
	var hostport;
	// hostport = HTTP_HOST; // global value from page context
	hostport = window.location.hostname + (location.port ? ':' + location.port : '');
	conn = new WebSocket('ws://' + hostport + '/ws'); // assert window["WebSocket"]

	conn.onopen = function() {
            connected = true;
            sendSearch(location.search);
            $(window).bind('popstate', function() {
                sendSearch(location.search);
            });
	};

        var statesel = 'table thead tr .header a.state';
	var again = function(_e) {
            $(statesel).unbind('click');
            window.setTimeout(init, 5000);
	};
	conn.onclose = again;
	conn.onerror = again;
	conn.onmessage = onmessage;

        $(statesel).click(function() {
            history.pushState({path: this.path}, '', this.href);
            sendSearch(this.search);
            return false;
	});
    }
    init();

    return {
        sendClient: sendClient,
        sendSearch: sendSearch
    };
}

// function empty(obj) {
//     return obj === undefined || obj === null;
// }
function emptyK(obj, key) {
    return (obj      === undefined ||
            obj      === null      ||
            obj[key] === undefined ||
            obj[key] === null);
}

var IFbytesCLASS = React.createClass({
  getInitialState: function() { return Data.IFbytes; },

  render: function() {
    var Data = {IFbytes: this.state};
    var rows = emptyK(Data.IFbytes, 'List') ?'': Data.IFbytes.List.map(function($if) { return ifbytes_rows(Data, $if); });
    return ifbytes_table(Data, rows);
  }
});

var IFerrorsCLASS = React.createClass({
  getInitialState: function() { return Data.IFerrors; },

  render: function() {
    var Data = {IFerrors: this.state};
    var rows = emptyK(Data.IFerrors, 'List') ?'': Data.IFerrors.List.map(function($if) { return iferrors_rows(Data, $if); });
    return iferrors_table(Data, rows);
  }
});

var IFpacketsCLASS = React.createClass({
  getInitialState: function() { return Data.IFpackets; },

  render: function() {
    var Data = {IFpackets: this.state};
    var rows = emptyK(Data.IFpackets, 'List') ?'': Data.IFpackets.List.map(function($if) { return ifpackets_rows(Data, $if); });
    return ifpackets_table(Data, rows);
  }
});

var DFbytesCLASS = React.createClass({
  getInitialState: function() { return {DFlinks: Data.DFlinks, DFbytes: Data.DFbytes}; },

  render: function() {
    var Data = this.state;
    var rows = emptyK(Data.DFbytes, 'List') ?'': Data.DFbytes.List.map(function($disk) { return dfbytes_rows(Data, $disk); });
    return dfbytes_table(Data, rows);
  }
});

var DFinodesCLASS = React.createClass({
  getInitialState: function() { return {DFlinks: Data.DFlinks, DFinodes: Data.DFinodes}; },

  render: function() {
    var Data = this.state;
    var rows = emptyK(Data.DFinodes, 'List') ?'': Data.DFinodes.List.map(function($disk) { return dfinodes_rows(Data, $disk); });
    return dfinodes_table(Data, rows);
  }
});

var MEMtableCLASS = React.createClass({
  getInitialState: function() { return Data.MEM; },

  render: function() {
    var Data = {MEM: this.state};
    var rows = emptyK(Data.MEM, 'List') ?'': Data.MEM.List.map(function($mem) { return mem_rows(Data, $mem); });
    return mem_table(Data, rows);
  }
});

var CPUtableCLASS = React.createClass({
  getInitialState: function() { return Data.CPU; },

  render: function() {
    var Data = {CPU: this.state};
    var rows = emptyK(Data.CPU, 'List') ?'': Data.CPU.List.map(function($core) { return cpu_rows(Data, $core); });
    return cpu_table(Data, rows);
  }
});

var PStableCLASS = React.createClass({
  getInitialState: function() { return {PStable: Data.PStable, PSlinks: Data.PSlinks}; },

  render: function() {
    var Data = this.state;
    var rows = emptyK(Data.PStable, 'List') ?'': Data.PStable.List.map(function($proc) { return ps_rows(Data, $proc); });
    return ps_table(Data, rows);
  }
});

var VGtableCLASS = React.createClass({
  getInitialState: function() { return {VagrantMachines: Data.VagrantMachines,
                                        VagrantError:  Data.VagrantError,
                                        VagrantErrord: Data.VagrantErrord
                                       }; },

  render: function() {
    var Data = this.state;
    var rows;
    if (Data.VagrantErrord !== undefined && Data.VagrantErrord) {
        rows = [vagrant_error(Data)];
    } else {
        rows = emptyK(Data.VagrantMachines, 'List') ?'': Data.VagrantMachines.List.map(function($machine) { return vagrant_rows(Data, $machine); });
    }
    return vagrant_table(Data, rows);
  }
});

var setState = function(obj, data, filterundefined) {
    if (data === undefined) { // null
        return;
    }
    // filter out undefined values
    data = _.object(_.filter(_.pairs(data), function(a) { return a[1] !== undefined; }));
    obj.setState(data);
};

var websocket; // a global

function update(currentClient, model) {
    var params = location.search.substr(1).split("&");
    for (var i in params) {
	if (params[i].split("=")[0] === "still") {
            return;
	}
    }

    // all *CLASS defined in gen/jscript.js
    var memtable  = React.renderComponent(MEMtableCLASS (null), document.getElementById('mem-table'));
    var pstable   = React.renderComponent(PStableCLASS  (null), document.getElementById('ps-table'));
    var dfbytes   = React.renderComponent(DFbytesCLASS  (null), document.getElementById('dfbytes-table'));
    var dfinodes  = React.renderComponent(DFinodesCLASS (null), document.getElementById('dfinodes-table'));
    var cputable  = React.renderComponent(CPUtableCLASS (null), document.getElementById('cpu-table'));
    var ifbytes   = React.renderComponent(IFbytesCLASS  (null), document.getElementById('ifbytes-table'));
    var iferrors  = React.renderComponent(IFerrorsCLASS (null), document.getElementById('iferrors-table'));
    var ifpackets = React.renderComponent(IFpacketsCLASS(null), document.getElementById('ifpackets-table'));
    var vagrant   = React.renderComponent(VGtableCLASS(null),   document.getElementById('vagrant-table'));

    var onmessage = function(event) {
	var data = JSON.parse(event.data);
        if (data.Client !== undefined && data.Client.DebugError !== undefined) {
            console.log('DEBUG ERROR', data.Client.DebugError);
        }

        setState(pstable,  {PStable:  data.PStable,  PSlinks: data.PSlinks});
	setState(dfbytes,  {DFbytes:  data.DFbytes,  DFlinks: data.DFlinks});
	setState(dfinodes, {DFinodes: data.DFinodes, DFlinks: data.DFlinks});

        setState(memtable,  data.MEM);
        setState(cputable,  data.CPU);
        setState(ifbytes,   data.IFbytes);
        setState(iferrors,  data.IFerrors);
	setState(ifpackets, data.IFpackets);
	setState(vagrant, {
            VagrantMachines: data.VagrantMachines,
            VagrantError:    data.VagrantError,
            VagrantErrord:   data.VagrantErrord
        });

        if (data.Client !== undefined) {
            console.log(JSON.stringify(data.Client), 'recvClient');
        }
        currentClient = _.extend(currentClient, data.Client);
        data.Client = currentClient;
        model.set(Model.attributes(data));

        // update the tooltips
        // $('span .tooltipable') .tooltip();
        $('span .tooltipable')    .popover({trigger: 'hover focus'});
        $('span .tooltipabledots').popover(); // the clickable dots
    };
    websocket = newwebsocket(onmessage);
}

var Model = Backbone.Model.extend({
    initialize: function() {
    }
});
Model.attributes = function(data) {
    return _.extend(data.Generic, data.Client);
};

var View = Backbone.View.extend({
    initialize: function() {
	this.listentext('IP',       $('#generic-ip'));
	this.listentext('Hostname', $('#generic-hostname'));
	this.listentext('Uptime',   $('#uptime #generic-uptime'));
	this.listentext('LA',       $('#generic-la'));

        var $hswapb = $('label[href="#showswap"]');
        this.listenactivate('HideSWAP', $hswapb, true);

        var $section_mem = $('#mem');
        var $section_if  = $('#if');
        var $section_cpu = $('#cpu');
        var $section_df  = $('#df');
        var $section_ps  = $('#ps');
        var $section_vg  = $('#vagrant');

        var $config_mem = $('#memconfig');
        var $config_if  = $('#ifconfig');
        var $config_cpu = $('#cpuconfig');
        var $config_df  = $('#dfconfig');
        var $config_ps  = $('#psconfig');
        var $config_vg  = $('#vgconfig');

        var $hidden_mem = $config_mem.find('.hiding');
        var $hidden_if  = $config_if .find('.hiding');
        var $hidden_cpu = $config_cpu.find('.hiding');
        var $hidden_df  = $config_df .find('.hiding');
        var $hidden_ps  = $config_ps .find('.hiding');
        var $hidden_vg  = $config_vg .find('.hiding');

        this.listenhide('HideMEM', $section_mem, $hidden_mem);
        this.listenhide('HideCPU', $section_cpu, $hidden_cpu);
        this.listenhide('HidePS',  $section_ps,  $hidden_ps);
        this.listenhide('HideVG',  $section_vg,  $hidden_vg);

        var $header_mem = $('header a[href="'+ $section_mem.selector +'"]');
        var $header_if  = $('header a[href="'+ $section_if .selector +'"]');
        var $header_cpu = $('header a[href="'+ $section_cpu.selector +'"]');
        var $header_df  = $('header a[href="'+ $section_df .selector +'"]');
        var $header_ps  = $('header a[href="'+ $section_ps .selector +'"]');
        var $header_vg  = $('header a[href="'+ $section_vg .selector +'"]');

        this.listentext('TabTitleIF', $header_if);
        this.listentext('TabTitleDF', $header_df);

        this.listenhide('HideconfigMEM', $config_mem, $header_mem, true);
        this.listenhide('HideconfigIF',  $config_if,  $header_if,  true);
        this.listenhide('HideconfigCPU', $config_cpu, $header_cpu, true);
        this.listenhide('HideconfigDF',  $config_df,  $header_df,  true);
        this.listenhide('HideconfigPS',  $config_ps,  $header_ps,  true);
        this.listenhide('HideconfigVG',  $config_vg,  $header_vg,  true);

        // NB by class
        var $tab_if    = $('.if-switch');
        var $tab_df    = $('.df-switch');
        var $panels_if = $('.if-tab');
        var $panels_df = $('.df-tab');

        this.listenTo(this.model, 'change:HideIF', this.change_collapsetabfunc('HideIF', 'TabIF', $panels_if, $tab_if));
        this.listenTo(this.model, 'change:HideDF', this.change_collapsetabfunc('HideDF', 'TabDF', $panels_df, $tab_df));
        this.listenTo(this.model, 'change:TabIF',  this.change_collapsetabfunc('HideIF', 'TabIF', $panels_if, $tab_if));
        this.listenTo(this.model, 'change:TabDF',  this.change_collapsetabfunc('HideDF', 'TabDF', $panels_df, $tab_df));

        var $psmore = $('label.more[href="#psmore"]');
        var $psless = $('label.less[href="#psless"]');
        this.listentext  ('PSplusText', $psmore);
        this.listenenable('PSnotExpandable',     $psmore);
//      this.listenenable('PSnotDecreasable',    $psless);

        this.listenrefresherror('RefreshErrorMEM', $config_mem.find('.refresh-group'));
        this.listenrefresherror('RefreshErrorIF',  $config_if .find('.refresh-group'));
        this.listenrefresherror('RefreshErrorCPU', $config_cpu.find('.refresh-group'));
        this.listenrefresherror('RefreshErrorDF',  $config_df .find('.refresh-group'));
        this.listenrefresherror('RefreshErrorPS',  $config_ps .find('.refresh-group'));
        this.listenrefresherror('RefreshErrorVG',  $config_vg .find('.refresh-group'));

        this.listenrefreshvalue('RefreshMEM',      $config_mem.find('.refresh-input'));
        this.listenrefreshvalue('RefreshIF',       $config_if .find('.refresh-input'));
        this.listenrefreshvalue('RefreshCPU',      $config_cpu.find('.refresh-input'));
        this.listenrefreshvalue('RefreshDF',       $config_df .find('.refresh-input'));
        this.listenrefreshvalue('RefreshPS',       $config_ps .find('.refresh-input'));
        this.listenrefreshvalue('RefreshVG',       $config_vg .find('.refresh-input'));

        var B = _.bind(function(c) { return _.bind(c, this); }, this);
        var expandable_sections = [
            [$section_if,  'ExpandIF',  'HideIF' ],
            [$section_cpu, 'ExpandCPU', 'HideCPU'],
            [$section_df,  'ExpandDF',  'HideDF' ]
        ];
        for (var i = 0; i < expandable_sections.length; ++i) {
            var S = expandable_sections[i][0];
            var E = expandable_sections[i][1];
            var H = expandable_sections[i][2];
            var $b = $('label[href="'+ S.selector +'"]');

            this.listenactivate(E, $b);
            $b.click( B(this.click_expandfunc(E, H)) );
        }

        $hswapb    .click( B(this.click_expandfunc('HideSWAP', 'HideMEM')) );
        $tab_if    .click( B(this.click_tabfunc('TabIF', 'HideIF')) );
        $tab_df    .click( B(this.click_tabfunc('TabDF', 'HideDF')) );

        $header_mem.click( B(this.click_expandfunc('HideconfigMEM')) );
        $header_if .click( B(this.click_expandfunc('HideconfigIF' )) );
        $header_cpu.click( B(this.click_expandfunc('HideconfigCPU')) );
        $header_df .click( B(this.click_expandfunc('HideconfigDF' )) );
        $header_ps .click( B(this.click_expandfunc('HideconfigPS' )) );
        $header_vg .click( B(this.click_expandfunc('HideconfigVG' )) );

        $hidden_mem.click( B(this.click_expandfunc('HideMEM')) );
        $hidden_if .click( B(this.click_expandfunc('HideIF' )) );
        $hidden_cpu.click( B(this.click_expandfunc('HideCPU')) );
        $hidden_df .click( B(this.click_expandfunc('HideDF' )) );
        $hidden_ps .click( B(this.click_expandfunc('HidePS' )) );
        $hidden_vg .click( B(this.click_expandfunc('HideVG' )) );

        $psmore    .click( B(this.click_psignalfunc('HidePS', true )) );
        $psless    .click( B(this.click_psignalfunc('HidePS', false)) );

        $config_mem.find('.refresh-input').on('input', B(this.submit_rsignalfunc('RefreshSignalMEM')) );
        $config_if .find('.refresh-input').on('input', B(this.submit_rsignalfunc('RefreshSignalIF' )) );
        $config_cpu.find('.refresh-input').on('input', B(this.submit_rsignalfunc('RefreshSignalCPU')) );
        $config_df .find('.refresh-input').on('input', B(this.submit_rsignalfunc('RefreshSignalDF' )) );
        $config_ps .find('.refresh-input').on('input', B(this.submit_rsignalfunc('RefreshSignalPS' )) );
        $config_vg .find('.refresh-input').on('input', B(this.submit_rsignalfunc('RefreshSignalVG' )) );
    },
    submit_rsignalfunc: function(R) {
        return function(e) {
            var sendc = _.object([R], [$(e.target).val()]);
            websocket.sendClient(sendc);
        };
    },

    listentext: function(K, $el) { this.listenTo(this.model, 'change:'+ K, this._text(K, $el)); },
//  listenHTML: function(K, $el) { this.listenTo(this.model, 'change:'+ K, this._HTML(K, $el)); },
         _text: function(K, $el) { return function() { var A = this.model.attributes; $el.text(A[K]); }; },
//       _HTML: function(K, $el) { return function() { var A = this.model.attributes; $el.html(A[K]); }; },

    listenrefresherror: function(E, $el) {
        this.listenTo(this.model, 'change:'+ E, function() {
            var A = this.model.attributes;
            $el[A[E] ? 'addClass' : 'removeClass']('has-warning');
        });
    },
    listenrefreshvalue: function(E, $el) {
        this.listenTo(this.model, 'change:'+ E, function() {
            var A = this.model.attributes;
            $el.prop('value', A[E]);
        });
    },

    listenenable: function(K, $el) {
        this.listenTo(this.model, 'change:'+ K, function() {
            var A = this.model.attributes;
            var V = A[K];
            V = V !== undefined && V;
            $el.prop('disabled', V);
            $el[V ? 'addClass' : 'removeClass']('disabled');
        });
    },

    listenactivate: function(K, $el, reverse) {
        this.listenTo(this.model, 'change:'+ K, function() {
            var A = this.model.attributes;
            var V = reverse !== undefined && reverse ? !A[K] : A[K];
            $el[V ? 'addClass' : 'removeClass']('active');
        });
    },
    listenhide: function(H, $el, $button_el, reverse) {
        this.listenTo(this.model, 'change:'+ H, function() {
            var A = this.model.attributes;
            $el.collapse(A[H] ? 'hide' : 'show'); // do what change_collapsefunc does

            // do what listenactivate does
            var V = reverse !== undefined && reverse ? !A[H] : A[H];
            $button_el[V ? 'addClass' : 'removeClass']('active');
        });
    },

    change_collapsefunc: function(H, $el) {
        return function() {
            var A = this.model.attributes;
            $el.collapse(A[H] ? 'hide' : 'show');
        };
    },

    change_collapsetabfunc: function(H, T, $el, $tabel) {
        return function() {
            var A = this.model.attributes;
            if (A[H]) { // hiding all
                $el.collapse('hide'); // do what change_collapsefunc does
                return;
            }
            var curtabid = A[T];
            var nots = _.map($el.not('[data-tabid="'+ curtabid +'"]'),
                             function(el) {
                                 var $el = $(el);
                                 $el.collapse('hide');
                                 return el;
                             });
            $($el.not(nots)).collapse('show');

            _.map($tabel, function(el) {
                var $el = $(el);
                var tabid_attr = +$el.attr('data-tabid'); // an int
                $el[tabid_attr === curtabid ? 'addClass' : 'removeClass']('active');
            });
        };
    },

    click_expandfunc: function(H, H2) {
        return function(e) {
            var A = this.model.attributes;
            var sendc = _.object([H], [!A[H]]);
            if (H2 !== undefined && A[H2]) { // if was hidden
                sendc = _.extend(sendc, _.object([H2], [!A[H2]]));
            }
            websocket.sendClient(sendc);
            e.preventDefault();
            e.stopPropagation(); // don't change checkbox/radio state
        };
    },
    click_tabfunc: function(T, H) {
        return function(e) {
            var newtabid = +$( $(e.target).attr('href') ).attr('data-tabid'); // THIS. +string makes an int
            var sendc = _.object([T], [newtabid]);
            var A = this.model.attributes;
            if (A[H]) { // if was hidden
                sendc = _.extend(sendc, _.object([H], [!A[H]]));
            }
            websocket.sendClient(sendc);
            e.preventDefault();
            e.stopPropagation(); // don't change checkbox/radio state
        };
    },
    click_psignalfunc: function(H, v) {
        return function(e) {
            var sendc = {MorePsignal: v};
            var A = this.model.attributes;
            if (A[H]) { // if was hidden
                sendc = _.extend(sendc, _.object([H], [!A[H]]));
            }
            websocket.sendClient(sendc);
            e.preventDefault();
            e.stopPropagation(); // don't change checkbox/radio state
        };
    }
});

function ready() {

    (new Headroom(document.querySelector("nav"), {
        offset: 71 - 51
        // "relative"" padding-top of the toprow
        // 71 is the absolute padding-top of the toprow
        // 51 is the height of the nav (50 +1px bottom border)
    })).init();

    $('.collapse').collapse({toggle: false}); // init collapsable objects

    // $('span .tooltipable')   .tooltip();
    $('span .tooltipable')      .popover({trigger: 'hover focus'});
    $('span .tooltipabledots')  .popover(); // the clickable dots

    $('[data-toggle="popover"]').popover(); // should be just #generic-hostname
    $('#generic-la').popover({
        trigger: 'hover focus',
        placement: 'right', // not 'auto right' until #generic-la is the last element for it's parent
        html: true, content: function() {
            return $('#uptime').html();
        }
    });

    $('body').on('click', function (e) { // hide the popovers on click outside
        $('span .tooltipabledots').each(function () {
            //the 'is' for buttons that trigger popups
            //the 'has' for icons within a button that triggers a popup
            if (!$(this).is(e.target) && $(this).has(e.target).length === 0 && $('.popover').has(e.target).length === 0) {
                $(this).popover('hide');
            }
        });
    });

    var model = new Model(Model.attributes(Data));
    var view  = new View({model: model});

    update(Data.Client, model);
}

// Local Variables:
// indent-tabs-mode: nil
// End:
