// requires underscore.js
// requires backbone.js
// requires jquery.js

function default_button(btn) {
    btn.addClass('btn-default');
    btn.removeClass('btn-primary');
}
function primary_button(btn) {
    btn.addClass('btn-primary');
    btn.removeClass('btn-default');
}

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
    function sendState(dict) {
	console.log(JSON.stringify(dict), 'sendState');
	return sendJSON({State: dict});
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
        sendState: sendState,
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
  getInitialState: function() { return Data.PStable; },

  render: function() {
    var Data = {PStable: this.state};
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

var websocket; // a global

function update(currentState, model) {
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

        var setState = function(obj, data) {
            if (data !== undefined) { // null
                obj.setState(data);
            }
        };

        setState(pstable, data.PStable);

	var bytestate = {DFbytes: data.DFbytes};
	if (data.DFlinks !== undefined) { bytestate.DFlinks = data.DFlinks; }
	setState(dfbytes, bytestate);

	var inodestate = {DFinodes: data.DFinodes};
	if (data.DFlinks !== undefined) { inodestate.DFlinks = data.DFlinks; }
	setState(dfinodes, inodestate);

        setState(memtable,  data.MEM);
        setState(cputable,  data.CPU);
        setState(ifbytes,   data.IFbytes);
        setState(iferrors,  data.IFerrors);
	setState(ifpackets, data.IFpackets);
	setState(vagrant, {
            VagrantMachines: data.VagrantMachines,
            VagrantError:  data.VagrantError,
            VagrantErrord: data.VagrantErrord
        });

        if (data.ClientState !== undefined) {
            console.log(JSON.stringify(data.ClientState), 'recvState');
        }
        currentState = _.extend(currentState, data.ClientState);
        data.ClientState = currentState;
        model.set(Model.attributes(data));

        // update the tooltips
        // $('span .tooltipable').tooltip();
        $('span .tooltipable').popover({trigger: 'hover focus'});
        $('span .tooltipabledots').popover(); // the clickable dots
    };
    websocket = newwebsocket(onmessage);
}

var Model = Backbone.Model.extend({
    initialize: function() {
    }
});
Model.attributes = function(data) {
    var A = _.extend(data.Generic, data.ClientState);
    A = _.extend(A, {
        PlusText: data.PStable.PlusText
    });
    return A;
};

var View = Backbone.View.extend({
    initialize: function() {
	this.listenchange_Textfunc('IP',           $('#generic-ip'));
	this.listenchange_HTMLfunc('HostnameHTML', $('#generic-hostname'));
	this.listenchange_Textfunc('Uptime',       $('#generic-uptime'));
	this.listenchange_Textfunc('LA',           $('#generic-la'));

        var $hswapb = $('label[href="#showswap"]');
        this.listenchange_buttonfunc('HideSWAP', $hswapb, true);

        var $section_mem = $('#mem');
        var $section_if  = $('#if');
        var $section_cpu = $('#cpu');
        var $section_df  = $('#df');
        var $section_ps  = $('#ps');
        var $section_vg  = $('#vagrant');
        this.listenhide('HideMEM', $section_mem);
        this.listenhide('HideCPU', $section_cpu);
        this.listenhide('HidePS',  $section_ps);
        this.listenhide('HideVG',  $section_vg);

        var $config_mem = $('#memconfig');
        var $config_if  = $('#ifconfig');
        var $config_cpu = $('#cpuconfig');
        var $config_df  = $('#dfconfig');
        var $config_ps  = $('#psconfig');
        var $config_vg  = $('#vgconfig');

        this.listenhide('HideconfigMEM', $config_mem);
        this.listenhide('HideconfigIF',  $config_if);
        this.listenhide('HideconfigCPU', $config_cpu);
        this.listenhide('HideconfigDF',  $config_df);
        this.listenhide('HideconfigPS',  $config_ps);
        this.listenhide('HideconfigVG',  $config_vg);

        var $tab_if    = $('label.network-switch');
        var $tab_df    = $('label.disk-switch');
        var $panels_if = $('.network-tab'); // by class
        var $panels_df = $('.disk-tab');    // by class

        this.listenTo(this.model, 'change:HideIF', this.change_collapsetabfunc('HideIF', 'TabIF', $panels_if, $tab_if));
        this.listenTo(this.model, 'change:HideDF', this.change_collapsetabfunc('HideDF', 'TabDF', $panels_df, $tab_df));
        this.listenTo(this.model, 'change:TabIF',  this.change_collapsetabfunc('HideIF', 'TabIF', $panels_if, $tab_if));
        this.listenTo(this.model, 'change:TabDF',  this.change_collapsetabfunc('HideDF', 'TabDF', $panels_df, $tab_df));

        var $psmore = $('label.more[href="#psmore"]');
        var $psless = $('label.less[href="#psless"]');
        this.listenchange_Textfunc('PlusText', $psmore);

        var B = _.bind(function(c) { return _.bind(c, this); }, this);

        $hswapb.click( B(this.click_expandfunc('HideSWAP')) );
        $tab_if.click( B(this.click_tabfunc('TabIF', 'HideIF')) );
        $tab_df.click( B(this.click_tabfunc('TabDF', 'HideDF')) );

        var expandable_sections = [
            [$section_if,  'ExpandIF',  'HideIF'],
            [$section_cpu, 'ExpandCPU', 'HideCPU'],
            [$section_df,  'ExpandDF',  'HideDF' ]
        ];
        for (var i = 0; i < expandable_sections.length; ++i) {
            var S  = expandable_sections[i][0];
            var K  = expandable_sections[i][1];
            var KK = expandable_sections[i][2];
            var $b = $('label.all[href="'+ S.selector +'"]');

            this.listenchange_buttonfunc(K, $b);
            $b.click( B(this.click_expandfunc(K, KK)) );
        }

        $('[href="'+ $config_mem.selector +'"]').click( B(this.click_expandfunc('HideconfigMEM', 'HideMEM')) );
        $('[href="'+ $config_if .selector +'"]').click( B(this.click_expandfunc('HideconfigIF',  'HideIF' )) );
        $('[href="'+ $config_cpu.selector +'"]').click( B(this.click_expandfunc('HideconfigCPU', 'HideCPU')) );
        $('[href="'+ $config_df .selector +'"]').click( B(this.click_expandfunc('HideconfigDF',  'HideDF' )) );
        $('[href="'+ $config_ps .selector +'"]').click( B(this.click_expandfunc('HideconfigPS',  'HidePS' )) );
        $('[href="'+ $config_vg .selector +'"]').click( B(this.click_expandfunc('HideconfigVG',  'HideVG' )) );

        $('header a[href="'+ $section_mem.selector +'"]').click( B(this.click_expandfunc('HideMEM', 'HideconfigMEM', true)) );
        $('header a[href="'+ $section_if .selector +'"]').click( B(this.click_expandfunc('HideIF',  'HideconfigIF',  true)) );
        $('header a[href="'+ $section_cpu.selector +'"]').click( B(this.click_expandfunc('HideCPU', 'HideconfigCPU', true)) );
        $('header a[href="'+ $section_df .selector +'"]').click( B(this.click_expandfunc('HideDF',  'HideconfigDF',  true)) );
        $('header a[href="'+ $section_ps .selector +'"]').click( B(this.click_expandfunc('HidePS',  'HideconfigPS',  true)) );
        $('header a[href="'+ $section_vg .selector +'"]').click( B(this.click_expandfunc('HideVG',  'HideconfigVG',  true)) );

        $psmore.click( B(this.click_psignalfunc('HidePS', true )) );
        $psless.click( B(this.click_psignalfunc('HidePS', false)) );
    },

    listenchange_Textfunc: function(K, $el){ this.listenTo(this.model, 'change:'+ K, this.change_Textfunc(K, $el)); },
    listenchange_HTMLfunc: function(K, $el){ this.listenTo(this.model, 'change:'+ K, this.change_HTMLfunc(K, $el)); },
          change_Textfunc: function(K, $el) { return function() { var A = this.model.attributes; $el.text(A[K]); }; },
          change_HTMLfunc: function(K, $el) { return function() { var A = this.model.attributes; $el.html(A[K]); }; },

    listenchange_buttonfunc: function(K, $el, reverse) {
        this.listenTo(this.model, 'change:'+ K, function() {
            var A = this.model.attributes;
            var V = reverse !== undefined && reverse ? !A[K] : A[K];
            var c = V ? primary_button : default_button;
            c($el);
        });
    },
    listenhide: function(K, $el) {
        this.listenTo(this.model, 'change:'+ K, this.change_collapsefunc(K, $el));
    },

    change_collapsetabfunc: function(K, KK, $el, $tabel) {
        return function() {
            var A = this.model.attributes;
            if (A[K]) { // hiding all
                $el.collapse('hide'); // do what change_collapsefunc does
                return;
            }
            var curtabid = A[KK];
            var nots = _.map($el.not('[data-tabid="'+ curtabid +'"]'),
                             function(el) {
                                 var $el = $(el);
                                 $el.collapse('hide');
                                 return el;
                             });
            $($el.not(nots)).collapse('show');

            _.map($tabel, function(el) {
                var $el = $(el);

                if (!$el.hasClass('nondefault')) {
                    return;
                }
                var tabid_attr = +$el.attr('data-tabid'); // an int
                if (tabid_attr === curtabid) {
                    primary_button($el);
                } else {
                    default_button($el);
                }
            });
        };
    },
    change_collapsefunc: function(K, $el) {
        return function() {
            var A = this.model.attributes;
            $el.collapse(A[K] ? 'hide' : 'show');
        };
    },

    click_psignalfunc: function(KK, v) {
        return function(e) {
            var newstate = {MorePsignal: v};
            var A = this.model.attributes;
            if (A[KK]) { // if was hidden
                newstate = _.extend(newstate, _.object([KK], [!A[KK]]));
            }
            websocket.sendState(newstate);
            e.preventDefault();
            e.stopPropagation(); // don't check/uncheck the checkbox
        };
    },
    click_tabfunc: function(K, KK) {
        return function(e) {
            var newtabid = +$( $(e.target).attr('href') ).attr('data-tabid'); // THIS. +string makes an int
            var newstate = _.object([K], [newtabid]);
            var A = this.model.attributes;
            if (A[KK]) { // if was hidden
                newstate = _.extend(newstate, _.object([KK], [!A[KK]]));
            }
            websocket.sendState(newstate);
            e.preventDefault();
        };
    },
    click_expandfunc: function(K, KK, isheader) {
        isheader = isheader !== undefined && isheader;
        return function(e) {
            var A = this.model.attributes;
            var V = A[K];
            var newstate = _.object([K], [!V]);
            if (KK !== undefined) {
                do {
                    if (V) {
                        if (isheader || !A[KK]) {
                            break;
                        }
                    } else if (!isheader || A[KK]) {
                        break;
                    }
                    newstate = _.extend(newstate, _.object([KK], [!A[KK]]));
                } while (0);
            }
            websocket.sendState(newstate);
            e.preventDefault();
        };
    }
});

function ready() {

    // construct an instance of Headroom, passing the element
    (new Headroom(document.querySelector("nav"), {
        offset: 71 - 51 // "padding-top" of the toprow container - navbar height (50px by default) + bottom border (1px)
    })).init();

    $('.collapse').collapse({toggle: false}); // init collapsable objects

    // $('span .tooltipable').tooltip();
    $('span .tooltipable').popover({trigger: 'hover focus'});
    $('span .tooltipabledots').popover(); // the clickable dots

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

    update(Data.ClientState, model);
}

// Local Variables:
// indent-tabs-mode: nil
// End:
