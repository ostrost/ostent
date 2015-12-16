let React      = require('react'),
    ReactDOM   = require('react-dom'),
    ReactPRM   = require('react-prm'),
    SparkLines = require('react-sparklines');
let ReactPureRenderMixin = ReactPRM;

var Sparkline = React.createClass({
  mixins: [ReactPureRenderMixin],
  getInitialState: function() { return {data: [], limit: 90, width: 180}; },
  componentDidUpdate: function(_, prevState) {
    var root = ReactDOM.findDOMNode(this.refs.root);
    if (root == null) {
      return;
    }
    var rootWidth = Math.floor(root.offsetWidth) - 10;
    if (prevState.width != rootWidth) {
      this.setState({width: rootWidth, limit: Math.round(rootWidth/2)});
    }
  },
  NewStateFrom: function(statentry) {
    var limit, data = [];
    if (this.state != null) {
      limit = this.state.limit;
      data  = this.state.data.slice(); // NB .slice https://github.com/borisyankov/react-sparklines/issues/27
    }
    if (this.props.col != null) {
      statentry = statentry[this.props.col];
    }
    data.push(+statentry);
    if (limit != null && data.length > limit) {
      data = data.slice(-limit);
    }
    this.setState({data: data});
  },
  render: function() {
    var spotsProps = {spotColors: {'-1': 'green', '1': 'red'}}; // reverse default
    if (this.props.defaultSpots) { delete spotsProps.spotColors; } // back to default
    return <div className="height-1rem" ref="root">
      <SparkLines.Sparklines
               data={this.state.data}
               limit={this.state.limit}
               width={this.state.width}
               height={35}>
        <SparkLines.SparklinesLine />
        <SparkLines.SparklinesSpots {...spotsProps} />
      </SparkLines.Sparklines>
    </div>;
  }
});

let jsdefines = {};
jsdefines.Sparkline = function(props) { return <Sparkline {...props} />; }

jsdefines.StateHandlingMixin = { // requires .Reduce method
  getInitialState: function() {
    return this.StateFrom(Data); // global Data
  },
  NewState: function(data) {
    let state = this.StateFrom(data);
    if (state != null) {
      this.setState(state);
    }
    var rkeys = Object.keys(this.refs);
    if (rkeys.length == 0) {
      return;
    }
    var statefrom;
    if (this.List != null) {
      statefrom = this.List(state);
    } else {
      var skeys = Object.keys(state);
      if (skeys.length != 1) {
        return;
      }
      statefrom = state[skeys[0]];
    }
    rkeys.forEach(function(rk) {
      var statentry;
      if (this.refs[rk] == null || (statentry = statefrom[rk]) == null) {
        return;
      }
      this.refs[rk].NewStateFrom(statentry);
    }, this);
  },
  StateFrom: function(data) {
    let state = this.Reduce(data);
    if (state != null) {
      for (let key in state) {
        if (state[key] == null) {
          delete state[key];
        }
      }
    }
    return state;
  }
};
jsdefines.HandlerMixin = {
  handleClick: function(e) {
    let href = e.target.getAttribute('href');
    if (href == null) {
      href = e.target.parentNode.getAttribute('href');
    }
    history.pushState({}, '', href);
    window.updates.sendSearch(href);
    e.stopPropagation();
    e.preventDefault();
    return void 0;
  }
};

// transformed from define_* templates:

jsdefines.define_hostname = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  Reduce: function(data) {
    return {
      hostname: data.hostname
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return <a href="/"  title={"hostname " + Data.hostname}
  >{Data.hostname}</a
>;
  }
});

jsdefines.define_panelcpu = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  List: function(data) {
    let list;
    if (data == null || data["cpu"] == null || (list = data["cpu"].List) == null) {
      return [];
    }
    return list;
  },
  Reduce: function(data) {
    return {
      params: data.params,
      cpu: data.cpu
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return <div className="grid-block hr-top"
  ><div className="col-lr large-1 text-right"
    ><div  className={!Data.params.CPUn.Negative ? "hide-showhide" : "show-showhide"}
      ><h1 className="h4 text-overflow"
        ><a  href={Data.params.Tlinks.CPUn} onClick={this.handleClick} title="CPU display options"
          ><span className="showhide-hide whitespace-pre float-left"
            >... </span
          >CPU</a
        ></h1
      ></div
    ></div
  ><div className="col-lr large-11"
    ><div  className={!Data.params.CPUn.Negative ? "hide-showhide" : "show-showhide"}
      ><div className="grid-block"
        ><ul className="menu showhide-show"
          ><li
            ><div className="input-group"
              ><div className="input-group-label small text-nowrap"
                >delay</div
              ><div className="input-group-button"
                ><a className="button small secondary disabled"
                  >{Data.params.CPUd}</a
                ></div
              ><div className="input-group-button"
                ><a href={Data.params.Dlinks.CPUd.Less.Href} className={"button small text-nowrap" + " " + (Data.params.Dlinks.CPUd.Less.ExtraClass != null ? Data.params.Dlinks.CPUd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Dlinks.CPUd.Less.Text}</a
></div
              ><div className="input-group-button"
                ><a href={Data.params.Dlinks.CPUd.More.Href} className={"button small text-nowrap" + " " + (Data.params.Dlinks.CPUd.More.ExtraClass != null ? Data.params.Dlinks.CPUd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Dlinks.CPUd.More.Text} +</a
></div
              ></div
            ></li
          ><li
            ><div className="input-group"
              ><div className="input-group-label small text-nowrap"
                >rows</div
              ><div className="input-group-button"
                ><a className="button small secondary disabled"
                  >{Data.params.CPUn.Absolute}</a
                ></div
              ><div className="input-group-button"
                ><a href={Data.params.Nlinks.CPUn.Less.Href} className={"button small success text-nowrap" + " " + (Data.params.Nlinks.CPUn.Less.ExtraClass != null ? Data.params.Nlinks.CPUn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Nlinks.CPUn.Less.Text}</a
></div
              ><div className="input-group-button"
                ><a href={Data.params.Nlinks.CPUn.More.Href} className={"button small success text-nowrap" + " " + (Data.params.Nlinks.CPUn.More.ExtraClass != null ? Data.params.Nlinks.CPUn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Nlinks.CPUn.More.Text} +</a
></div
              ></div
            ></li
          ></ul
        ></div
      ></div
    ><div  className={Data.params.CPUn.Absolute != 0 ? "stripe" : "hide"}
      ><div className="grid-block thead"
        ><span className="col small-1"
          >Core</span
        ><div className="grid-block wrap noscroll"
          ><span className="col small-1 text-right"
            > User</span
          ><span className="col small-1 text-right"
            > Sys</span
          ><span className="col small-1 text-right"
            > Wait</span
          ><span className="col small-1 text-right"
            > Idle</span
          ><span className="col"
            ></span
          ></div
        ></div
      >{this.List(Data).map(function($cpu, i) { return<div className="grid-block"
        ><span className="col small-1 text-nowrap"
          >{$cpu.N}</span
        ><div  key={"cpu-rowby-N-"+$cpu.N} className="grid-block wrap noscroll"
          ><span className="small-1 text-right"
            ><span className="col display-block bg-usepct"
data-usepct={$cpu.UserPct}
              > {$cpu.UserPct}%</span
            ></span
          ><span className="small-1 text-right"
            ><span className="col display-block bg-usepct"
data-usepct={$cpu.SysPct}
              > {$cpu.SysPct}%</span
            ></span
          ><span className="small-1 text-right"
            ><span className="col display-block bg-usepct"
data-usepct={$cpu.WaitPct}
              > {$cpu.WaitPct}%</span
            ></span
          ><span className="small-1 text-right"
            ><span className="col display-block bg-usepct-inverse"
data-usepct={$cpu.IdlePct}
              > {$cpu.IdlePct}%</span
            ></span
          ><span className="col-lr expand"
            >{jsdefines.Sparkline({ref: i, col: 'IdlePct', defaultSpots: true})}</span
          ></div
        ></div
      >})}</div
    ></div
  ></div
>;
  }
});

jsdefines.define_paneldf = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  List: function(data) {
    let list;
    if (data == null || data["diskUsage"] == null || (list = data["diskUsage"].List) == null) {
      return [];
    }
    return list;
  },
  Reduce: function(data) {
    return {
      params: data.params,
      diskUsage: data.diskUsage
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return <div className="grid-block hr-top"
  ><div className="col-lr large-1 text-right"
    ><div  className={!Data.params.Dfn.Negative ? "hide-showhide" : "show-showhide"}
      ><h1 className="h4 text-overflow"
        ><a  href={Data.params.Tlinks.Dfn} onClick={this.handleClick} title="Disk usage display options"
          ><span className="showhide-hide whitespace-pre float-left"
            >... </span
          >Disk usage</a
        ></h1
      ></div
    ></div
  ><div className="col-lr large-11"
    ><div  className={!Data.params.Dfn.Negative ? "hide-showhide" : "show-showhide"}
      ><div className="grid-block"
        ><ul className="menu showhide-show"
          ><li
            ><div className="input-group"
              ><div className="input-group-label small text-nowrap"
                >delay</div
              ><div className="input-group-button"
                ><a className="button small secondary disabled"
                  >{Data.params.Dfd}</a
                ></div
              ><div className="input-group-button"
                ><a href={Data.params.Dlinks.Dfd.Less.Href} className={"button small text-nowrap" + " " + (Data.params.Dlinks.Dfd.Less.ExtraClass != null ? Data.params.Dlinks.Dfd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Dlinks.Dfd.Less.Text}</a
></div
              ><div className="input-group-button"
                ><a href={Data.params.Dlinks.Dfd.More.Href} className={"button small text-nowrap" + " " + (Data.params.Dlinks.Dfd.More.ExtraClass != null ? Data.params.Dlinks.Dfd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Dlinks.Dfd.More.Text} +</a
></div
              ></div
            ></li
          ><li
            ><div className="input-group"
              ><div className="input-group-label small text-nowrap"
                >rows</div
              ><div className="input-group-button"
                ><a className="button small secondary disabled"
                  >{Data.params.Dfn.Absolute}</a
                ></div
              ><div className="input-group-button"
                ><a href={Data.params.Nlinks.Dfn.Less.Href} className={"button small success text-nowrap" + " " + (Data.params.Nlinks.Dfn.Less.ExtraClass != null ? Data.params.Nlinks.Dfn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Nlinks.Dfn.Less.Text}</a
></div
              ><div className="input-group-button"
                ><a href={Data.params.Nlinks.Dfn.More.Href} className={"button small success text-nowrap" + " " + (Data.params.Nlinks.Dfn.More.ExtraClass != null ? Data.params.Nlinks.Dfn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Nlinks.Dfn.More.Text} +</a
></div
              ></div
            ></li
          ></ul
        ></div
      ></div
    ><div  className={Data.params.Dfn.Absolute != 0 ? "stripe" : "hide"}
      ><div className="grid-block thead"
        ><span className="col small-1"
          ><a href={Data.params.Vlinks.Dfk[1-1].LinkHref} className={Data.params.Vlinks.Dfk[1-1].LinkClass} onClick={this.handleClick}  
  >Device<span className={Data.params.Vlinks.Dfk[1-1].CaretClass}
    ></span
  ></a
></span
        ><div className="grid-block wrap noscroll"
          ><span className="col small-1"
            > <a href={Data.params.Vlinks.Dfk[2-1].LinkHref} className={Data.params.Vlinks.Dfk[2-1].LinkClass} onClick={this.handleClick}  
  >Mounted<span className={Data.params.Vlinks.Dfk[2-1].CaretClass}
    ></span
  ></a
></span
          ><span className="col small-1 text-right"
            > <a href={Data.params.Vlinks.Dfk[6-1].LinkHref} className={Data.params.Vlinks.Dfk[6-1].LinkClass} onClick={this.handleClick}  
  >Total<span className={Data.params.Vlinks.Dfk[6-1].CaretClass}
    ></span
  ></a
></span
          ><span className="col small-1 text-right"
            > <a href={Data.params.Vlinks.Dfk[5-1].LinkHref} className={Data.params.Vlinks.Dfk[5-1].LinkClass} onClick={this.handleClick}  
  >Used<span className={Data.params.Vlinks.Dfk[5-1].CaretClass}
    ></span
  ></a
></span
          ><span className="col small-1 text-right"
            > <a href={Data.params.Vlinks.Dfk[3-1].LinkHref} className={Data.params.Vlinks.Dfk[3-1].LinkClass} onClick={this.handleClick}  
  >Avail<span className={Data.params.Vlinks.Dfk[3-1].CaretClass}
    ></span
  ></a
></span
          ><span className="col small-1 text-right"
            > <a href={Data.params.Vlinks.Dfk[4-1].LinkHref} className={Data.params.Vlinks.Dfk[4-1].LinkClass} onClick={this.handleClick}  
  >Use%<span className={Data.params.Vlinks.Dfk[4-1].CaretClass}
    ></span
  ></a
></span
          ><span className="col"
            ></span
          ></div
        ></div
      >{this.List(Data).map(function($df, i) { return<div className="grid-block"
        ><span className="col small-1 text-overflow"
          >{$df.DevName}</span
        ><div  key={"df-rowby-dirname-"+$df.DirName} className="grid-block wrap noscroll"
          ><span className="col small-1 text-overflow"
            > {$df.DirName}</span
          ><span className="col small-1 text-nowrap text-right"
            ><span className="mutext" title="Inodes total"
              > {$df.Inodes}</span
            > {$df.Total}</span
          ><span className="col small-1 text-nowrap text-right"
            ><span className="mutext" title="Inodes used"
              > {$df.Iused}</span
            > {$df.Used}</span
          ><span className="col small-1 text-nowrap text-right"
            ><span className="mutext" title="Inodes free"
              > {$df.Ifree}</span
            > {$df.Avail}</span
          ><span className="small-1 text-right text-nowrap"
            ><span className="col display-block bg-usepct" data-usepct={$df.UsePct}
              ><span className="mutext" title="Inodes use%"
                > {$df.IusePct}%</span
              > {$df.UsePct}%</span
            ></span
          ><span className="col-lr expand"
            >{jsdefines.Sparkline({ref: i, col: 'UsePct'})}</span
          ></div
        ></div
      >})}</div
    ></div
  ></div
>;
  }
});

jsdefines.define_panelif = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  List: function(data) {
    let list;
    if (data == null || data["ifaddrs"] == null || (list = data["ifaddrs"].List) == null) {
      return [];
    }
    return list;
  },
  Reduce: function(data) {
    return {
      params: data.params,
      ifaddrs: data.ifaddrs
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return <div className="grid-block hr-top"
  ><div className="col-lr large-1 text-right"
    ><div  className={!Data.params.Ifn.Negative ? "hide-showhide" : "show-showhide"}
      ><h1 className="h4 text-overflow"
        ><a  href={Data.params.Tlinks.Ifn} onClick={this.handleClick} title="Interfaces display options"
          ><span className="showhide-hide whitespace-pre float-left"
            >... </span
          >Interfaces</a
        ></h1
      ></div
    ></div
  ><div className="col-lr large-11"
    ><div  className={!Data.params.Ifn.Negative ? "hide-showhide" : "show-showhide"}
      ><div className="grid-block"
        ><ul className="menu showhide-show"
          ><li
            ><div className="input-group"
              ><div className="input-group-label small text-nowrap"
                >delay</div
              ><div className="input-group-button"
                ><a className="button small secondary disabled"
                  >{Data.params.Ifd}</a
                ></div
              ><div className="input-group-button"
                ><a href={Data.params.Dlinks.Ifd.Less.Href} className={"button small text-nowrap" + " " + (Data.params.Dlinks.Ifd.Less.ExtraClass != null ? Data.params.Dlinks.Ifd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Dlinks.Ifd.Less.Text}</a
></div
              ><div className="input-group-button"
                ><a href={Data.params.Dlinks.Ifd.More.Href} className={"button small text-nowrap" + " " + (Data.params.Dlinks.Ifd.More.ExtraClass != null ? Data.params.Dlinks.Ifd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Dlinks.Ifd.More.Text} +</a
></div
              ></div
            ></li
          ><li
            ><div className="input-group"
              ><div className="input-group-label small text-nowrap"
                >rows</div
              ><div className="input-group-button"
                ><a className="button small secondary disabled"
                  >{Data.params.Ifn.Absolute}</a
                ></div
              ><div className="input-group-button"
                ><a href={Data.params.Nlinks.Ifn.Less.Href} className={"button small success text-nowrap" + " " + (Data.params.Nlinks.Ifn.Less.ExtraClass != null ? Data.params.Nlinks.Ifn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Nlinks.Ifn.Less.Text}</a
></div
              ><div className="input-group-button"
                ><a href={Data.params.Nlinks.Ifn.More.Href} className={"button small success text-nowrap" + " " + (Data.params.Nlinks.Ifn.More.ExtraClass != null ? Data.params.Nlinks.Ifn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Nlinks.Ifn.More.Text} +</a
></div
              ></div
            ></li
          ></ul
        ></div
      ></div
    ><div  className={Data.params.Ifn.Absolute != 0 ? "stripe" : "hide"}
      ><div className="grid-block thead"
        ><span className="col small-1"
          >Interface</span
        ><div className="grid-block wrap noscroll"
          ><span className="col small-1 text-right"
            > IP</span
          ><span className="col small-2 text-right text-nowrap" title="Drops,Errors In/Out per second"
            > Loss IO ps</span
          ><span className="col small-2 text-right text-nowrap" title="Packets In/Out per second"
            > Packets IO ps</span
          ><span className="col small-2 text-right text-nowrap" title="Bits In/Out per second"
            > IO <i
              >b</i
            >ps</span
          ><span className="col"
            ></span
          ></div
        ></div
      >{this.List(Data).map(function($if, i) { return<div className="grid-block"
        ><span className="col small-1 text-overflow"
          >{$if.Name}</span
        ><div  key={"if-rowby-name-"+$if.Name} className="grid-block wrap noscroll"
          ><span className="col small-1 text-right text-overflow"
            >{$if.IP}</span
          ><span className="col small-2 text-right text-nowrap"
            > <span className="mutext" title="Total drops,errors modulo 4G"
              ><span title="Total drops In modulo 4G"
                >{$if.DropsIn}</span
              ><span  className={$if.DropsOut != null ? "" : "hide"}
                >/</span
              ><span  className={$if.DropsOut != null ? "" : "hide"} title="Total drops Out modulo 4G"
                >{$if.DropsOut}</span
              >,<span title="Total errors In modulo 4G"
                >{$if.ErrorsIn}</span
              >/<span title="Total errors Out modulo 4G"
                >{$if.ErrorsOut}</span
              ></span
            > <span  className={(($if.DeltaDropsIn == null || $if.DeltaDropsIn == "0") && ($if.DeltaDropsOut == null || $if.DeltaDropsOut == "0") && ($if.DeltaErrorsIn == null || $if.DeltaErrorsIn == "0") && ($if.DeltaErrorsOut == null || $if.DeltaErrorsOut == "0")) ? "mutext" : ""}
              ><span title="Drops In per second"
                >{$if.DeltaDropsIn}</span
              ><span  className={$if.DeltaDropsOut != null ? "" : "hide"}
                >/</span
              ><span  className={$if.DeltaDropsOut != null ? "" : "hide"} title="Drops Out per second"
                >{$if.DeltaDropsOut}</span
              >,<span title="Errors In per second"
                >{$if.DeltaErrorsIn}</span
              >/<span title="Errors Out per second"
                >{$if.DeltaErrorsOut}</span
              ></span
            ></span
          ><span className="col small-2 text-right text-nowrap"
            > <span className="mutext"
              ><span title="Total packets In modulo 4G"
                >{$if.PacketsIn}</span
              >/<span title="Total packets Out modulo 4G"
                >{$if.PacketsOut}</span
              ></span
            > <span title="Packets In per second"
              >{$if.DeltaPacketsIn}</span
            >/<span title="Packets Out per second"
              >{$if.DeltaPacketsOut}</span
            ></span
          ><span className="col small-2 text-right text-nowrap"
            > <span className="mutext"
              ><span title="Total BYTES In modulo 4G"
                >{$if.BytesIn}</span
              >/<span title="Total BYTES Out modulo 4G"
                >{$if.BytesOut}</span
              ></span
            > <span title="BITS In per second"
              >{$if.DeltaBitsIn}</span
            >/<span title="BITS Out per second"
              >{$if.DeltaBitsOut}</span
            ></span
          ><span className="col-lr expand"
            >{jsdefines.Sparkline({ref: i, col: 'DeltaBytesOutNum'})}</span
          ></div
        ></div
      >})}</div
    ></div
  ></div
>;
  }
});

jsdefines.define_panella = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  List: function(data) {
    let list;
    if (data == null || data["loadavg"] == null || (list = data["loadavg"].List) == null) {
      return [];
    }
    return list;
  },
  Reduce: function(data) {
    return {
      params: data.params,
      loadavg: data.loadavg
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return <div className="grid-block hr-top"
  ><div className="col-lr large-1 text-right"
    ><div  className={!Data.params.Lan.Negative ? "hide-showhide" : "show-showhide"}
      ><h1 className="h4 text-overflow"
        ><a  href={Data.params.Tlinks.Lan} onClick={this.handleClick} title="Load avg display options"
          ><span className="showhide-hide whitespace-pre float-left"
            >... </span
          >Load avg</a
        ></h1
      ></div
    ></div
  ><div className="col-lr large-11"
    ><div  className={!Data.params.Lan.Negative ? "hide-showhide" : "show-showhide"}
      ><div className="grid-block"
        ><ul className="menu showhide-show"
          ><li
            ><div className="input-group"
              ><div className="input-group-label small text-nowrap"
                >delay</div
              ><div className="input-group-button"
                ><a className="button small secondary disabled"
                  >{Data.params.Lad}</a
                ></div
              ><div className="input-group-button"
                ><a href={Data.params.Dlinks.Lad.Less.Href} className={"button small text-nowrap" + " " + (Data.params.Dlinks.Lad.Less.ExtraClass != null ? Data.params.Dlinks.Lad.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Dlinks.Lad.Less.Text}</a
></div
              ><div className="input-group-button"
                ><a href={Data.params.Dlinks.Lad.More.Href} className={"button small text-nowrap" + " " + (Data.params.Dlinks.Lad.More.ExtraClass != null ? Data.params.Dlinks.Lad.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Dlinks.Lad.More.Text} +</a
></div
              ></div
            ></li
          ><li
            ><div className="input-group"
              ><div className="input-group-label small text-nowrap"
                >rows</div
              ><div className="input-group-button"
                ><a className="button small secondary disabled"
                  >{Data.params.Lan.Absolute}</a
                ></div
              ><div className="input-group-button"
                ><a href={Data.params.Nlinks.Lan.Less.Href} className={"button small success text-nowrap" + " " + (Data.params.Nlinks.Lan.Less.ExtraClass != null ? Data.params.Nlinks.Lan.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Nlinks.Lan.Less.Text}</a
></div
              ><div className="input-group-button"
                ><a href={Data.params.Nlinks.Lan.More.Href} className={"button small success text-nowrap" + " " + (Data.params.Nlinks.Lan.More.ExtraClass != null ? Data.params.Nlinks.Lan.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Nlinks.Lan.More.Text} +</a
></div
              ></div
            ></li
          ></ul
        ></div
      ></div
    ><div  className={Data.params.Lan.Absolute != 0 ? "stripe" : "hide"}
      ><div className="grid-block thead"
        ><span className="col small-1"
          >Period</span
        ><div className="grid-block wrap noscroll"
          ><span className="col small-1 text-right"
            > Value</span
          ><span className="col"
            ></span
          ></div
        ></div
      >{this.List(Data).map(function($la, i) { return<div className="grid-block"
        ><span className="col small-1"
          >{$la.Period}m</span
        ><div  key={"la-rowby-period-"+$la.Period} className="grid-block wrap noscroll"
          ><span className="col small-1 text-right"
            > {$la.Value}</span
          ><span className="col-lr expand"
            >{jsdefines.Sparkline({ref: i, col: 'Value'})}</span
          ></div
        ></div
      >})}</div
    ></div
  ></div
>;
  }
});

jsdefines.define_panelmem = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  List: function(data) {
    let list;
    if (data == null || data["memory"] == null || (list = data["memory"].List) == null) {
      return [];
    }
    return list;
  },
  Reduce: function(data) {
    return {
      params: data.params,
      memory: data.memory
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return <div className="grid-block hr-top"
  ><div className="col-lr large-1 text-right"
    ><div  className={!Data.params.Memn.Negative ? "hide-showhide" : "show-showhide"}
      ><h1 className="h4 text-overflow"
        ><a  href={Data.params.Tlinks.Memn} onClick={this.handleClick} title="Memory display options"
          ><span className="showhide-hide whitespace-pre float-left"
            >... </span
          >Memory</a
        ></h1
      ></div
    ></div
  ><div className="col-lr large-11"
    ><div  className={!Data.params.Memn.Negative ? "hide-showhide" : "show-showhide"}
      ><div className="grid-block"
        ><ul className="menu showhide-show"
          ><li
            ><div className="input-group"
              ><div className="input-group-label small text-nowrap"
                >delay</div
              ><div className="input-group-button"
                ><a className="button small secondary disabled"
                  >{Data.params.Memd}</a
                ></div
              ><div className="input-group-button"
                ><a href={Data.params.Dlinks.Memd.Less.Href} className={"button small text-nowrap" + " " + (Data.params.Dlinks.Memd.Less.ExtraClass != null ? Data.params.Dlinks.Memd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Dlinks.Memd.Less.Text}</a
></div
              ><div className="input-group-button"
                ><a href={Data.params.Dlinks.Memd.More.Href} className={"button small text-nowrap" + " " + (Data.params.Dlinks.Memd.More.ExtraClass != null ? Data.params.Dlinks.Memd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Dlinks.Memd.More.Text} +</a
></div
              ></div
            ></li
          ><li
            ><div className="input-group"
              ><div className="input-group-label small text-nowrap"
                >rows</div
              ><div className="input-group-button"
                ><a className="button small secondary disabled"
                  >{Data.params.Memn.Absolute}</a
                ></div
              ><div className="input-group-button"
                ><a href={Data.params.Nlinks.Memn.Less.Href} className={"button small success text-nowrap" + " " + (Data.params.Nlinks.Memn.Less.ExtraClass != null ? Data.params.Nlinks.Memn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Nlinks.Memn.Less.Text}</a
></div
              ><div className="input-group-button"
                ><a href={Data.params.Nlinks.Memn.More.Href} className={"button small success text-nowrap" + " " + (Data.params.Nlinks.Memn.More.ExtraClass != null ? Data.params.Nlinks.Memn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Nlinks.Memn.More.Text} +</a
></div
              ></div
            ></li
          ></ul
        ></div
      ></div
    ><div  className={Data.params.Memn.Absolute != 0 ? "stripe" : "hide"}
      ><div className="grid-block thead"
        ><span className="col small-1"
          >Memory</span
        ><div className="grid-block wrap noscroll"
          ><span className="col small-1 text-right"
            > Total</span
          ><span className="col small-1 text-right"
            > Used</span
          ><span className="col small-1 text-right"
            > Free</span
          ><span className="col small-1 text-right"
            > Use%</span
          ><span className="col"
            ></span
          ></div
        ></div
      >{this.List(Data).map(function($mem, i) { return<div className="grid-block"
        ><span className="col small-1"
          >{$mem.Kind}</span
        ><div  key={"mem-rowby-kind-"+$mem.Kind} className="grid-block wrap noscroll"
          ><span className="col small-1 text-right"
            > {$mem.Total}</span
          ><span className="col small-1 text-right"
            > {$mem.Used}</span
          ><span className="col small-1 text-right"
            > {$mem.Free}</span
          ><span className="small-1 text-right"
            ><span className="col display-block bg-usepct" data-usepct={$mem.UsePct}
              > {$mem.UsePct}%</span
            ></span
          ><span className="col-lr expand"
            >{jsdefines.Sparkline({ref: i, col: 'UsePct'})}</span
          ></div
        ></div
      >})}</div
    ></div
  ></div
>;
  }
});

jsdefines.define_panelps = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  List: function(data) {
    let list;
    if (data == null || data["procs"] == null || (list = data["procs"].List) == null) {
      return [];
    }
    return list;
  },
  Reduce: function(data) {
    return {
      params: data.params,
      procs: data.procs
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return <div className="grid-block hr-top"
  ><div className="col-lr large-1 text-right"
    ><div  className={!Data.params.Psn.Negative ? "hide-showhide" : "show-showhide"}
      ><h1 className="h4 text-overflow"
        ><a  href={Data.params.Tlinks.Psn} onClick={this.handleClick} title="Processes display options"
          ><span className="showhide-hide whitespace-pre float-left"
            >... </span
          >Processes</a
        ></h1
      ></div
    ></div
  ><div className="col-lr large-11"
    ><div  className={!Data.params.Psn.Negative ? "hide-showhide" : "show-showhide"}
      ><div className="grid-block"
        ><ul className="menu showhide-show"
          ><li
            ><div className="input-group"
              ><div className="input-group-label small text-nowrap"
                >delay</div
              ><div className="input-group-button"
                ><a className="button small secondary disabled"
                  >{Data.params.Psd}</a
                ></div
              ><div className="input-group-button"
                ><a href={Data.params.Dlinks.Psd.Less.Href} className={"button small text-nowrap" + " " + (Data.params.Dlinks.Psd.Less.ExtraClass != null ? Data.params.Dlinks.Psd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Dlinks.Psd.Less.Text}</a
></div
              ><div className="input-group-button"
                ><a href={Data.params.Dlinks.Psd.More.Href} className={"button small text-nowrap" + " " + (Data.params.Dlinks.Psd.More.ExtraClass != null ? Data.params.Dlinks.Psd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Dlinks.Psd.More.Text} +</a
></div
              ></div
            ></li
          ><li
            ><div className="input-group"
              ><div className="input-group-label small text-nowrap"
                >rows</div
              ><div className="input-group-button"
                ><a className="button small secondary disabled"
                  >{Data.params.Psn.Absolute}</a
                ></div
              ><div className="input-group-button"
                ><a href={Data.params.Nlinks.Psn.Less.Href} className={"button small success text-nowrap" + " " + (Data.params.Nlinks.Psn.Less.ExtraClass != null ? Data.params.Nlinks.Psn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Nlinks.Psn.Less.Text}</a
></div
              ><div className="input-group-button"
                ><a href={Data.params.Nlinks.Psn.More.Href} className={"button small success text-nowrap" + " " + (Data.params.Nlinks.Psn.More.ExtraClass != null ? Data.params.Nlinks.Psn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Nlinks.Psn.More.Text} +</a
></div
              ></div
            ></li
          ></ul
        ></div
      ></div
    ><div  className={Data.params.Psn.Absolute != 0 ? "stripe" : "hide"}
      ><div className="grid-block thead"
        ><span className="col small-1 text-right"
          ><a href={Data.params.Vlinks.Psk[1-1].LinkHref} className={Data.params.Vlinks.Psk[1-1].LinkClass} onClick={this.handleClick}  
  >PID<span className={Data.params.Vlinks.Psk[1-1].CaretClass}
    ></span
  ></a
></span
        ><div className="grid-block wrap noscroll text-nowrap"
          ><span className="col small-1 text-right"
            > <a href={Data.params.Vlinks.Psk[2-1].LinkHref} className={Data.params.Vlinks.Psk[2-1].LinkClass} onClick={this.handleClick}  
  >UID<span className={Data.params.Vlinks.Psk[2-1].CaretClass}
    ></span
  ></a
></span
          ><span className="col small-1"
            > <a href={Data.params.Vlinks.Psk[3-1].LinkHref} className={Data.params.Vlinks.Psk[3-1].LinkClass} onClick={this.handleClick}  
  >USER<span className={Data.params.Vlinks.Psk[3-1].CaretClass}
    ></span
  ></a
></span
          ><span className="col small-1 text-right"
            > <a href={Data.params.Vlinks.Psk[4-1].LinkHref} className={Data.params.Vlinks.Psk[4-1].LinkClass} onClick={this.handleClick}  
  >PR<span className={Data.params.Vlinks.Psk[4-1].CaretClass}
    ></span
  ></a
></span
          ><span className="col small-1 text-right"
            > <a href={Data.params.Vlinks.Psk[5-1].LinkHref} className={Data.params.Vlinks.Psk[5-1].LinkClass} onClick={this.handleClick}  
  >NI<span className={Data.params.Vlinks.Psk[5-1].CaretClass}
    ></span
  ></a
></span
          ><span className="col small-1 text-right"
            > <a href={Data.params.Vlinks.Psk[6-1].LinkHref} className={Data.params.Vlinks.Psk[6-1].LinkClass} onClick={this.handleClick}  
  >VIRT<span className={Data.params.Vlinks.Psk[6-1].CaretClass}
    ></span
  ></a
></span
          ><span className="col small-1 text-right"
            > <a href={Data.params.Vlinks.Psk[7-1].LinkHref} className={Data.params.Vlinks.Psk[7-1].LinkClass} onClick={this.handleClick}  
  >RES<span className={Data.params.Vlinks.Psk[7-1].CaretClass}
    ></span
  ></a
></span
          ><span className="col small-1 text-center"
            > <a href={Data.params.Vlinks.Psk[8-1].LinkHref} className={Data.params.Vlinks.Psk[8-1].LinkClass} onClick={this.handleClick}  
  >TIME<span className={Data.params.Vlinks.Psk[8-1].CaretClass}
    ></span
  ></a
></span
          ><span className="col small-1"
            > <a href={Data.params.Vlinks.Psk[9-1].LinkHref} className={Data.params.Vlinks.Psk[9-1].LinkClass} onClick={this.handleClick}  
  >COMMAND<span className={Data.params.Vlinks.Psk[9-1].CaretClass}
    ></span
  ></a
></span
          ></div
        ></div
      >{this.List(Data).map(function($ps) { return<div className="grid-block"
        ><span className="col small-1 text-right"
          >{$ps.PID}</span
        ><div  key={"ps-rowby-pid-"+$ps.PID} className="grid-block wrap noscroll"
          ><span className="col small-1 text-right"
            > {$ps.UID}</span
          ><span className="col small-1"
            > {$ps.User}</span
          ><span className="col small-1 text-right"
            > {$ps.Priority}</span
          ><span className="col small-1 text-right"
            > {$ps.Nice}</span
          ><span className="col small-1 text-right"
            > {$ps.Size}</span
          ><span className="col small-1 text-right"
            > {$ps.Resident}</span
          ><span className="col small-1 text-center"
            > {$ps.Time}</span
          ><span className="col expand"
            > {$ps.Name}</span
          ></div
        ></div
      >})}</div
    ></div
  ></div
>;
  }
});

jsdefines.define_uptime = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  Reduce: function(data) {
    return {
      uptime: data.uptime
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return <span
  >{Data.uptime}</span
>;
  }
});


module.exports = jsdefines;

// Local variables:
// js-indent-level: 2
// js2-basic-offset: 2
// End:
