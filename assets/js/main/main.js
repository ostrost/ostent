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

function Updatables(initial) {
    var models = [];
    function make(modelClass, opt, viewClass, elopt) {
        var el;
        if (elopt !== undefined) {
            el  =  elopt.el_added;
            delete elopt.el_added;
            _.extend(opt, elopt); // !

            if (el === undefined) {
                el = $();
            }
            _.map(elopt, function(a) { // add the values
                el = el.add(a);
            });
        }
        if (opt === undefined) {
            opt = {};
        }
        var modinit = modelClass.modelAttributes(initial);
        opt = _.extend(opt, modinit);            // !
        opt = _.extend(opt, {initial: modinit}); // !
        var model = new modelClass(opt);
	models.push(model);

        if (viewClass !== undefined) {
            new viewClass({el: el, model: model});
        }
	return model;
    }
    function set(data) {
	for (var i = 0; i < models.length; i++) {
            var ma = models[i].modelAttributes(data);
            for (var k in ma) {
                if (ma[k] === null) {
                    delete ma[k];
                }
            }
            models[i].set(ma);
	}
    }
    return {
	make: make,
	set:  set
    };
}

Updatables.declareModel = function(opt) {
    if (typeof opt == 'function') {
        opt = opt();
    }
    var modelClass = Backbone.Model.extend(opt);
    modelClass.modelAttributes = opt.modelAttributes;
    return modelClass;
};

Updatables.declareCollapseModel = function(ATTRIBUTE_HIDE) {
    return Updatables.declareModel(function () {
        return {
            Attribute_Hide: ATTRIBUTE_HIDE,
            modelAttributes: function(data) { return {
                Hide: data.ClientState[ATTRIBUTE_HIDE]
            }; },
            toggleHidden: function(s) { return _.object([ATTRIBUTE_HIDE], [!s.Hide]); }
        };
    });
};

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

function update(currentState, updatables) {
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

        setState(memtable, data.MEM);
        setState(cputable, data.CPU);
        setState(ifbytes,   data.IFbytes);
        setState(iferrors,  data.IFerrors);
	setState(ifpackets, data.IFpackets);
	setState(vagrant, {VagrantMachines: data.VagrantMachines,
                           VagrantError:  data.VagrantError,
                           VagrantErrord: data.VagrantErrord
                          });

        if (data.ClientState !== undefined) {
            console.log(JSON.stringify(data.ClientState), 'recvState');
        }
        currentState = _.extend(currentState, data.ClientState);
        data.ClientState = currentState;
        updatables.set(data);

        // update the tooltips
        // $('span .tooltipable').tooltip();
        $('span .tooltipable').popover({trigger: 'hover focus'});
        $('span .tooltipabledots').popover(); // the clickable dots
    };
    websocket = newwebsocket(onmessage);
}

function empty(obj) {
    return obj === undefined || obj === null;
}

var CollapseView = Backbone.View.extend({
    hidden: function() { return this.model.attributes.Hide; },

    initialize: function() {
	this.listenTo(this.model, 'change:Hide', this.redisplay_panel);
        this.init_target();
    },
    init_target: function() {
	_.map(this.model.attributes.target, function(t) {
            $(t).collapse({toggle: false}); // init collapsable objects
	}, this);
    },
    redisplay_panel: function() {
	_.map(this.model.attributes.target, function(t) {
            $(t).collapse(this.hidden() ? 'hide' : 'show');
	}, this);
    },
    toggleHidden: function() {
	return this.model.toggleHidden(this.model.attributes);
    },

    events: {'click': 'collapse_click'},
    collapse_click: function(e) {
	websocket.sendState(this.toggleHidden());
        e.preventDefault();
    }
});

var ExpandView = CollapseView.extend({
    expanded: function() { return this.model.attributes.Expand; },

    initialize_fromswitch: function() {
	this.listenTo(this.model, 'change:ExpandText', this.redisplay_expandtext);
	this.listenTo(this.model, 'change:Expandable', this.redisplay_expandable);
    },
    initialize: function() {
        this.initialize_fromswitch();
	this.listenTo(this.model, 'change:Expand',     this.change_expand);
        CollapseView.prototype.initialize.call(this); // does this.init_target();
    },
    change_expand: function() {
        this.redisplay_panel();
        this.redisplay_expand();
    },
    redisplay_expandtext: function() {
        var $expand_el = this.model.attributes.expand_el;
        $expand_el.text(this.model.attributes.ExpandText);
    },
    redisplay_expand: function() {
        var $expand_el = this.model.attributes.expand_el;
        if (this.expanded()) {
            primary_button($expand_el);
        } else {
            default_button($expand_el);
        }
    },
    redisplay_expandable: function() {
        var $expand_el = this.model.attributes.expand_el;
        var expandable = this.model.attributes.Expandable;
        if (empty(expandable)) {
            $expand_el.addClass('disabled');
        } else {
            $expand_el.removeClass('disabled');
        }
    },
    toggleExpandedState: function() {
        var te = this.model.toggleExpanded; //this.toggleExpanded !== undefined ? this.toggleExpanded : this.model.toggleExpanded;
	return te(this.model.attributes);
    },

    events: {'click': 'expand_click'},
    expand_click: function(e) {
	var clicked = $(e.target);
        if (clicked.is(this.model.attributes.header_el)) { // header clicked
            this.collapse_click(e);
            return;
        }

	var newState = this.toggleExpandedState();
	if (this.hidden()) { // if the panel was hidden by the header link
            newState = _.extend(newState, this.toggleHidden());
	}
	websocket.sendState(newState);
    }
});

var SwitchView = CollapseView.extend(ExpandView.prototype).extend({
    initialize: function() {
        // CollapseView.prototype.initialize.call(this); // DO NOT CollapseView.initialize
	this.listenTo(this.model, 'change:Hide',       this.change_switch); // <- as in CollapseView.initialize
	this.listenTo(this.model, 'change:Expand',     this.change_switch); // <- as in ExpandView.initialize
	this.listenTo(this.model, 'change:CurrentTab', this.change_switch);
        ExpandView.prototype.initialize_fromswitch.call(this);
        this.init_target();
    },
    redisplay_tabs: function() {
	var target = this.model.attributes.target;
        var tabid  = this.model.tabid();
        var nots = _.map(target.not('[data-tabid="'+ tabid +'"]'),
                         function(el) {
                             var $el = $(el);
                             $el.collapse('hide');
                             return el;
                         });
        var el = target.not(nots);
        $(el).collapse('show');
    },
    redisplay_buttons: function() {
        var tabid = this.model.tabid();
        // _.map(this.$el, function(el) {
        _.map(this.model.attributes.switch_el, function(el) {
            var $el = $(el);

            if (!$el.hasClass('nondefault')) {
                return;
            }
            var tabid_attr = +$el.attr('data-tabid'); // an int
            if (tabid_attr === tabid) {
                primary_button($el);
            } else {
                default_button($el);
            }
        }, this);
    },
    change_switch: function() {
        if (this.hidden()) {
            this.redisplay_panel();
        }
        this.redisplay_tabs();
        this.redisplay_buttons();
        this.redisplay_expand(this.model.attributes.expand_el);
    },

    events: {'click': 'switch_click'},
    switch_click: function(e) {
	var clicked = $(e.target);
        if (clicked.is(this.model.attributes.header_el)) { // header clicked
            this.collapse_click(e);
            return;
        } else if (clicked.is(this.model.attributes.expand_el)) {
            this.expand_click(e);
            return;
        }

        var newtab_id = +$( clicked.attr('href') ).attr('data-tabid'); // THIS. +string makes an int
        var newState = this.model.setTabState(newtab_id);

	if (this.hidden()) { // if the panel was hidden by the header link
            newState = _.extend(newState, this.toggleHidden());
	}
	websocket.sendState(newState);
    }
});

function SwitchModel(self) {
    var easy = {};
    easy.tabid = function() { return this.attributes.CurrentTab; }; // return this.attributes[self.Attribute_CurrentTab];

    self = _.extend(self, easy);
    self.modelAttributes = function(data) {
        return {
            Hide:       data.ClientState[self.Attribute_Hide],
            Expand:     data.ClientState[self.Attribute_Expand],
            CurrentTab: data.ClientState[self.Attribute_CurrentTab],
            Expandable: data[self.Attribute_Data].Expandable, // immutable
            ExpandText: data[self.Attribute_Data].ExpandText  // immutable
        };
    };

    self.toggleHidden   = function(s) { return _.object([self.Attribute_Hide],   [!s.Hide]);   }; //, [!s[self.Attribute_Hide]]
    self.toggleExpanded = function(s) { return _.object([self.Attribute_Expand], [!s.Expand]); }; //, [!s[self.Attribute_Expand]]
    self.setTabState    = function(n) { return _.object([self.Attribute_CurrentTab], [n]); };

    return Updatables.declareModel(self);
}

var DFswitchmodel = SwitchModel({
    Attribute_CurrentTab: 'TabDF',
    Attribute_Expand:     'ExpandDF',
    Attribute_Hide:       'HideDF',
    Attribute_Data:       'DF'
});

var IFswitchmodel = SwitchModel({
    Attribute_CurrentTab: 'TabIF',
    Attribute_Expand:     'ExpandIF',
    Attribute_Hide:       'HideIF',
    Attribute_Data:       'IF'
});

var ExpandMEMModel = Updatables.declareModel(function() {
    var self = {
        Attribute_Expand: 'HideSWAP',
        Attribute_Hide:   'HideMEM'
    };
    self.modelAttributes = function(data) {
        return {
            Expand:    !data.ClientState[self.Attribute_Expand], // NB inverse
            Hide:       data.ClientState[self.Attribute_Hide]
        };
    };

    self.toggleHidden   = function(s) { return _.object([self.Attribute_Hide],   [!s.Hide]);   };
    self.toggleExpanded = function(s) { return _.object([self.Attribute_Expand], [ s.Expand]); }; // NB reverse, thus not "!"

    return self;
});

var ExpandVGModel = Updatables.declareModel(function() {
    var self = {
        Attribute_Hide:   'HideMEM'
    };
    self.modelAttributes = function(data) {
        return {
            Expand:    !data.ClientState[self.Attribute_Expand], // NB inverse
            Hide:       data.ClientState[self.Attribute_Hide]
        };
    };

    self.toggleHidden   = function(s) { return _.object([self.Attribute_Hide],   [!s.Hide]);   };
    self.toggleExpanded = function(s) { return _.object([self.Attribute_Expand], [ s.Expand]); }; // NB reverse, thus not "!"

    return self;
});

var ExpandCPUModel = Updatables.declareModel(function() {
    var self = {
        Attribute_Expand: 'ExpandCPU',
        Attribute_Hide:   'HideCPU'
    };
    self.modelAttributes = function(data) {
        var r = {
            Hide:       data.ClientState[self.Attribute_Hide],
            Expand:     data.ClientState[self.Attribute_Expand]
        };
        if (!empty(data.CPU)) {
            r = _.extend(r, {
                Expandable: data.CPU.Expandable, // immutable
                ExpandText: data.CPU.ExpandText  // immutable
            });
        }
        return r;
    };

    self.toggleHidden   = function(s) { return _.object([self.Attribute_Hide],   [!s.Hide]);   }; //, [!s[self.Attribute_Hide]]
    self.toggleExpanded = function(s) { return _.object([self.Attribute_Expand], [!s.Expand]); }; //, [!s[self.Attribute_Expand]]

    return self;
});

var PSmodel = Updatables.declareModel(function() {
    var self = {
        Attribute_Hide: 'HidePS'
    };
    self.modelAttributes = function(data) {
        var r = {
            Hide: data.ClientState[self.Attribute_Hide]
        };
        if (!empty(data.PStable)) {
            r = _.extend(r, {
                NotExpandable: data.PStable.NotExpandable,
                PlusText:      data.PStable.PlusText
            });
        }
        return r;
    };

    self.toggleHidden = function(s) { return _.object([self.Attribute_Hide], [!s.Hide]); }; //, [!s[self.Attribute_Hide]]

    self.more         = function()  { return {MorePsignal: true}; };
    self.less         = function()  { return {MorePsignal: false}; };
    return self;
});

var PSview = CollapseView.extend({
    initialize: function() {
        CollapseView.prototype.initialize.call(this);
	this.listenTo(this.model, 'change:PlusText',      this.change_plustext);
	this.listenTo(this.model, 'change:NotExpandable', this.change_notexpandable);
    },
    change_notexpandable: function() {
        var $more_el = this.model.attributes.more_el;
        if (this.model.attributes.NotExpandable) {
            $more_el.addClass('disabled');
        } else {
            $more_el.removeClass('disabled');
        }
    },
    change_plustext: function() {
        var $more_el = this.model.attributes.more_el;
        $more_el.text(this.model.attributes.PlusText);
    },

    events: {'click': 'ps_click'},
    ps_click: function(e) {
	var clicked = $(e.target);
        if (clicked.is(this.model.attributes.header_el)) { // header clicked
            this.collapse_click(e);
            return;
        }
        do {
            var func;
            if (clicked.is(this.model.attributes.more_el)) { // more clicked
                func = 'more';
            } else if (clicked.is(this.model.attributes.less_el)) { // less clicked
                func = 'less';
            } else {
                break;
            }
            var newState = this.model[func]();
            if (this.hidden()) { // if the panel was hidden by the header link
                newState = _.extend(newState, this.toggleHidden());
            }
            websocket.sendState(newState);
        } while (0);
        e.stopPropagation(); // otherwise input checkbox gets checked/unchecked
        e.preventDefault();
    }
});

var UpdateView = Backbone.View.extend({
    events: {}, // sanity check
    initialize: function() {
	this.listenTo(this.model, 'change:'+ this.update_key(), this.change);
    },
    update_key: function() {
        return _.keys(this.model.attributes.initial)[0];
    },
    change: function() {
        var key = this.update_key();
        var func = /HTML$/.test(key) ? 'html' : 'text';
        this.$el[func](this.model.attributes[key]);
    }
});

function ready() {
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

    var updatables = Updatables(Data);

    new CollapseView({ el: $('[href="#memconfig"]'), // MEM CONFIG
                       model: updatables.make(Updatables.declareCollapseModel('ConfigMEM'),
                                              {target: $('#memconfig')}) });

    new CollapseView({ el: $('[href="#ifconfig"]'), // IF CONFIG
                       model: updatables.make(Updatables.declareCollapseModel('ConfigIF'),
                                              {target: $('#ifconfig')}) });

    new CollapseView({ el: $('[href="#cpuconfig"]'), // CPU CONFIG
                       model: updatables.make(Updatables.declareCollapseModel('ConfigCPU'),
                                              {target: $('#cpuconfig')}) });

    new CollapseView({ el: $('[href="#dfconfig"]'), // DF CONFIG
                       model: updatables.make(Updatables.declareCollapseModel('ConfigDF'),
                                              {target: $('#dfconfig')}) });

    new CollapseView({ el: $('[href="#psconfig"]'), // PS CONFIG
                       model: updatables.make(Updatables.declareCollapseModel('ConfigPS'),
                                              {target: $('#psconfig')}) });

    new CollapseView({ el: $('[href="#vgconfig"]'), // VG CONFIG
                       model: updatables.make(Updatables.declareCollapseModel('ConfigVG'),
                                              {target: $('#vgconfig')}) });

    // updatables.make(Updatables.declareCollapseModel('HideMEM'), // MEMORY
    //                 {target: $('#mem')}, CollapseView, { el: $('header a[href="#mem"]') });

    updatables.make( // MEM
        ExpandMEMModel, {target: $('#mem')}, ExpandView, {
            header_el: $('header a[href="#mem"]'),
            expand_el: $('label[href="#showswap"]')
        });

    updatables.make( // IF
        IFswitchmodel, {target: $('.network-tab')}, SwitchView, {
            switch_el: $('label.network-switch'),
            header_el: $('header a[href="#if"]'),
            expand_el: $('label.all[href="#if"]')
        });

    updatables.make( // CPU
        ExpandCPUModel, {target: $('#cpu')}, ExpandView, {
            header_el: $('header a[href="#cpu"]'),
            expand_el: $('label.all[href="#cpu"]')
        });

    updatables.make( // DF
        DFswitchmodel, {target: $('.disk-tab')}, SwitchView, {
            switch_el: $('label.disk-switch'),
            header_el: $('header a[href="#df"]'),
            expand_el: $('label.all[href="#df"]')
        });

    updatables.make( // PS
        PSmodel, {target: $('#ps')}, PSview, {
            header_el: $('header a[href="#ps"]'),
            more_el: $('label.more[href="#psmore"]'),
            less_el: $('label.less[href="#psless"]')
        });

    updatables.make( // VG
        Updatables.declareCollapseModel('HideVG'), {target: $('#vagrant')},
        CollapseView, {header_el: $('header a[href="#vagrant"]')});

    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {IP: data.Generic.IP}; }}),
                    {}, UpdateView, {el: $('#generic-ip')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {HostnameHTML: data.Generic.HostnameHTML}; }}),
                    {}, UpdateView, {el: $('#generic-hostname')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Uptime: data.Generic.Uptime}; }}),
                    {}, UpdateView, {el: $('#generic-uptime')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {LA: data.Generic.LA}; }}),
                    {}, UpdateView, {el: $('#generic-la')});

    /*
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Free: data.RAM.Free}; }}),
                    {}, UpdateView, {el: $('#ram-free')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Used: data.RAM.Used}; }}),
                    {}, UpdateView, {el: $('#ram-used')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Total: data.RAM.Total}; }}),
                    {}, UpdateView, {el: $('#ram-total')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {UsePercentHTML: data.RAM.UsePercentHTML}; }}),
                    {}, UpdateView, {el: $('#ram-usepercent')});

    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Free: data.Swap.Free}; }}),
                    {}, UpdateView, {el: $('#swap-free')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Used: data.Swap.Used}; }}),
                    {}, UpdateView, {el: $('#swap-used')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Total: data.Swap.Total}; }}),
                    {}, UpdateView, {el: $('#swap-total')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {UsePercentHTML: data.Swap.UsePercentHTML}; }}),
                    {}, UpdateView, {el: $('#swap-usepercent')});
    */

    update(Data.ClientState, updatables);
}

// Local Variables:
// indent-tabs-mode: nil
// End:
