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

	var again = function(_e) {
            $("a.state").unbind('click');
            window.setTimeout(init, 5000);
	};
	conn.onclose = again;
	conn.onerror = again;
	conn.onmessage = onmessage;

        $("a.state").click(function() {
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

var websocket; // a global
function update(currentState, updatables) {
    var params = location.search.substr(1).split("&");
    for (var i in params) {
	if (params[i].split("=")[0] === "still") {
            return;
	}
    }

    // all *Class defined in gen/jscript.js
    var   procTable     = React.renderComponent(procTableClass      (null), document.getElementById('ps-table'));
    var disksinBytes    = React.renderComponent(disksinBytesClass   (null), document.getElementById('df-table'));
    var disksinInodes   = React.renderComponent(disksinInodesClass  (null), document.getElementById('dfi-table'));
    var    cpuTable     = React.renderComponent(cpuTableClass       (null), document.getElementById('cpu-table'));
    var    ifsTable     = React.renderComponent(ifsTableClass       (null), document.getElementById('ifs-table'));
    var ifsPacketsTable = React.renderComponent(ifsPacketsTableClass(null), document.getElementById('ifs-packets-table'));
    var ifsErrorsTable  = React.renderComponent(ifsErrorsTableClass (null), document.getElementById('ifs-errors-table'));

    var onmessage = function(event) {
	var data = JSON.parse(event.data);

        var setState = function(obj, data) {
            if (data !== undefined) { // null
                obj.setState(data);
            }
        };

        setState(procTable, data.ProcTable);

	var bytestate = {DisksinBytes: data.DisksinBytes};
	if (data.DiskLinks !== undefined) { bytestate.DiskLinks = data.DiskLinks; }
	setState(disksinBytes, bytestate);

	var inodestate = {DisksinInodes: data.DisksinInodes};
	if (data.DiskLinks !== undefined) { inodestate.DiskLinks = data.DiskLinks; }
	setState(disksinInodes, inodestate);

        setState(cpuTable, data.CPU);
        setState(ifsTable, data.InterfacesBytes);
	setState(ifsErrorsTable,  data.InterfacesErrors);
	setState(ifsPacketsTable, data.InterfacesPackets);

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

function empty(obj) {
    return obj === undefined || obj === null;
}

var ExpandButtonView = CollapseView.extend({
    expanded: function() { return this.model.attributes.Expand; },

    initialize_fromswitch: function() {
	this.listenTo(this.model, 'change:More',       this.redisplay_more);
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
    redisplay_more: function() {
        var $expand_el = this.model.attributes.expand_el;
        $expand_el.text(this.model.attributes.More);
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

var SwitchView = CollapseView.extend(ExpandButtonView.prototype).extend({
    initialize: function() {
        // CollapseView.prototype.initialize.call(this); // DO NOT CollapseView.initialize
	this.listenTo(this.model, 'change:Hide',       this.change_switch); // <- as in CollapseView.initialize
	this.listenTo(this.model, 'change:Expand',     this.change_switch); // <- as in ExpandButton.initialize
	this.listenTo(this.model, 'change:CurrentTab', this.change_switch);
        ExpandButtonView.prototype.initialize_fromswitch.call(this);
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
            More:       data[self.Attribute_Data].More        // immutable
        };
    };

    self.toggleHidden   = function(s) { return _.object([self.Attribute_Hide],   [!s.Hide]);   }; //, [!s[self.Attribute_Hide]]
    self.toggleExpanded = function(s) { return _.object([self.Attribute_Expand], [!s.Expand]); }; //, [!s[self.Attribute_Expand]]
    self.setTabState    = function(n) { return _.object([self.Attribute_CurrentTab], [n]); };

    return Updatables.declareModel(self);
}

var DisksSwitchModel = SwitchModel({
    Attribute_CurrentTab: 'CurrentDisksTab',
    Attribute_Expand:     'ExpandDisks',
    Attribute_Hide:       'HideDisks',
    Attribute_Data:       'Disks'
});

var NetworkSwitchModel = SwitchModel({
    Attribute_CurrentTab: 'CurrentNetworkTab',
    Attribute_Expand:     'ExpandNetwork',
    Attribute_Hide:       'HideNetwork',
    Attribute_Data:       'Network'
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
                More:       data.CPU.More        // immutable
            });
        }
        return r;
    };

    self.toggleHidden   = function(s) { return _.object([self.Attribute_Hide],   [!s.Hide]);   }; //, [!s[self.Attribute_Hide]]
    self.toggleExpanded = function(s) { return _.object([self.Attribute_Expand], [!s.Expand]); }; //, [!s[self.Attribute_Expand]]

    return self;
});

var ProcessesModel = Updatables.declareModel(function() {
    var self = {
        Attribute_Hide: 'HideProcesses'
    };
    self.modelAttributes = function(data) {
        var r = {
            Hide: data.ClientState[self.Attribute_Hide]
        };
        if (!empty(data.ProcTable)) {
            r = _.extend(r, {
                NotExpandable: data.ProcTable.NotExpandable,
                MoreText:      data.ProcTable.MoreText
            });
        }
        return r;
    };

    self.toggleHidden = function(s) { return _.object([self.Attribute_Hide], [!s.Hide]); }; //, [!s[self.Attribute_Hide]]

    self.more         = function()  { return {MoreProcessesSignal: true}; };
    self.less         = function()  { return {MoreProcessesSignal: false}; };
    return self;
});

var ProcessesView = CollapseView.extend({
    initialize: function() {
        CollapseView.prototype.initialize.call(this);
	this.listenTo(this.model, 'change:MoreText',      this.change_moretext);
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
    change_moretext: function() {
        var $more_el = this.model.attributes.more_el;
        $more_el.text(this.model.attributes.MoreText);
    },

    events: {'click': 'proc_click'},
    proc_click: function(e) {
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

    new CollapseView({ // MEMORY
	el: $('header a[href="#memory"]'),
	model: updatables.make(Updatables.declareModel(function () {
            var self = {
                Attribute_Hide: 'HideMemory'
            };
            self.modelAttributes = function(data) {
                return {
                    Hide: data.ClientState[self.Attribute_Hide]
                };
            };

            self.toggleHidden = function(s) { return _.object([self.Attribute_Hide], [!s.Hide]); }; //, [!s[self.Attribute_Hide]]
            return self;
        }), {target: $('#memory')})
    });

    updatables.make( // NETWORK
        NetworkSwitchModel, {target: $('.network-tab')}, SwitchView, {
            switch_el: $('label.network-switch'), // el_added:
            header_el: $('header  a[href="#network"]'),
            expand_el: $('label.all[href="#network"]')
        });

    updatables.make( // CPU
        ExpandCPUModel, {target: $('#cpu')}, ExpandButtonView, {
            header_el: $('header  a[href="#cpu"]'),
            expand_el: $('label.all[href="#cpu"]')
        });

    updatables.make( // DISKS
        DisksSwitchModel, {target: $('.disk-tab')}, SwitchView, {
            switch_el: $('label.disk-switch'), // el_added:
            header_el: $('header  a[href="#disks"]'),
            expand_el: $('label.all[href="#disks"]')
        });

    updatables.make( // PROCESSES
        ProcessesModel, {target: $('#processes')}, ProcessesView, {
            header_el: $('header a[href="#processes"]'),
            more_el: $('label.more[href="#psmore"]'),
            less_el: $('label.less[href="#psless"]')
        });

    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {IP: data.About.IP}; }}),
                    {}, UpdateView, {el: $('#About-IP')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {HostnameHTML: data.About.HostnameHTML}; }}),
                    {}, UpdateView, {el: $('#About-Hostname')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Uptime: data.System.Uptime}; }}),
                    {}, UpdateView, {el: $('#System-Uptime')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {LA: data.System.LA}; }}),
                    {}, UpdateView, {el: $('#System-LA')});

    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Free: data.RAM.Free}; }}),
                    {}, UpdateView, {el: $('#Data.RAM.Free')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Used: data.RAM.Used}; }}),
                    {}, UpdateView, {el: $('#Data.RAM.Used')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Total: data.RAM.Total}; }}),
                    {}, UpdateView, {el: $('#Data.RAM.Total')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {UsePercentHTML: data.RAM.UsePercentHTML}; }}),
                    {}, UpdateView, {el: $('#Data.RAM.UsePercent')});

    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Free: data.Swap.Free}; }}),
                    {}, UpdateView, {el: $('#Data.Swap.Free')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Used: data.Swap.Used}; }}),
                    {}, UpdateView, {el: $('#Data.Swap.Used')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {Total: data.Swap.Total}; }}),
                    {}, UpdateView, {el: $('#Data.Swap.Used')});
    updatables.make(Updatables.declareModel({modelAttributes: function(data) { return {UsePercentHTML: data.Swap.UsePercentHTML}; }}),
                    {}, UpdateView, {el: $('#Data.Swap.UsePercent')});

    update(Data.ClientState, updatables);
}

// Local Variables:
// indent-tabs-mode: nil
// End:
