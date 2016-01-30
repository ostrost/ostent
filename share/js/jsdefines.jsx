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
    return <div ref="root">
      <SparkLines.Sparklines
               data={this.state.data}
               limit={this.state.limit}
               width={this.state.width}
               height={33}>
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

jsdefines.define_cpu = React.createClass({
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
    return (
<div className="grid-block hr-top">
  <div className="col-lr large-1 text-right"><div className={Data.params.CPUn.Negative ? "show-showhide" : "hide-showhide"}>
    <h1 className="h4 text-overflow"><a title="CPU display options" href={Data.params.Tlinks.CPUn} onClick={this.handleClick}><span className="showhide-hide whitespace-pre float-left">... </span>CPU</a>
    </h1></div>
  </div>
  <div className="col-lr large-11"><div className={Data.params.CPUn.Negative ? "show-showhide" : "hide-showhide"}>
    <div className="grid-block">
      <ul className="menu showhide-show">
        <li>
          <div className="input-group">
            <div className="input-group-label small text-nowrap">delay</div>
            <div className="input-group-button"><a className="button small secondary disabled">{Data.params.CPUd}</a></div>
            <div className="input-group-button">
              <a href={Data.params.Dlinks.CPUd.Less.Href} onClick={this.handleClick}
                 className={"text-nowrap button small " + (Data.params.Dlinks.CPUd.Less.ExtraClass != null ? Data.params.Dlinks.CPUd.Less.ExtraClass : "")}
                >- {Data.params.Dlinks.CPUd.Less.Text}</a>
            </div>
            <div className="input-group-button">
              <a href={Data.params.Dlinks.CPUd.More.Href} onClick={this.handleClick}
                 className={"text-nowrap button small " + (Data.params.Dlinks.CPUd.More.ExtraClass != null ? Data.params.Dlinks.CPUd.More.ExtraClass : "")}
                >{Data.params.Dlinks.CPUd.More.Text} +</a>
            </div>
          </div>
        </li>
        <li>
          <div className="input-group">
            <div className="input-group-label small text-nowrap">rows</div>
            <div className="input-group-button"><a className="button small secondary disabled">{Data.params.CPUn.Absolute}</a></div>
            <div className="input-group-button">
              <a href={Data.params.Nlinks.CPUn.Less.Href} onClick={this.handleClick}
                 className={"text-nowrap button small success " + (Data.params.Nlinks.CPUn.Less.ExtraClass != null ? Data.params.Nlinks.CPUn.Less.ExtraClass : "")}
                >- {Data.params.Nlinks.CPUn.Less.Text}</a>
            </div>
            <div className="input-group-button">
              <a href={Data.params.Nlinks.CPUn.More.Href} onClick={this.handleClick}
                 className={"text-nowrap button small success " + (Data.params.Nlinks.CPUn.More.ExtraClass != null ? Data.params.Nlinks.CPUn.More.ExtraClass : "")}
                >{Data.params.Nlinks.CPUn.More.Text} +</a>
            </div>
          </div>
        </li>
      </ul>
    </div><div className={Data.params.CPUn.Absolute == 0 ? "hide":""}>
    <div className="grid-block vertical stripe">
      <div className="grid-block thead"><span className="expand col small-1">Core</span><span className="expand col small-1 text-right"> User%</span><span className="expand col small-1 text-right"> Sys%</span><span className="expand col small-1 text-right"> Wait%</span><span className="expand col small-1 text-right"> Idle%</span><span className="expand col"></span></div>
      
      
      {this.List(Data).map(function($cpu, i) { return (
      <div className="grid-block" key={"cpu-rowby-n-"+$cpu.N}><span className="expand col small-1 text-nowrap">{$cpu.N}</span><span className="expand col small-1 text-right bg-usepct"
       data-usepct={$cpu.UserPct}> {$cpu.UserPct}%</span>
      <span className="expand col small-1 text-right bg-usepct"
       data-usepct={$cpu.SysPct}> {$cpu.SysPct}%</span>
      <span className="expand col small-1 text-right bg-usepct"
       data-usepct={$cpu.WaitPct}> {$cpu.WaitPct}%</span>
      <span className="expand col small-1 text-right bg-usepct-inverse"
       data-usepct={$cpu.IdlePct}> {$cpu.IdlePct}%</span><span className="expand col-lr">{jsdefines.Sparkline({ref: i, col: 'IdlePct', defaultSpots: true})}</span></div>);})}
      
    </div></div></div>
  </div>
</div>);
  }
});

jsdefines.define_df = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  List: function(data) {
    let list;
    if (data == null || data["df"] == null || (list = data["df"].List) == null) {
      return [];
    }
    return list;
  },
  Reduce: function(data) {
    return {
      params: data.params,
      df: data.df
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return (
<div className="grid-block hr-top">
  <div className="col-lr large-1 text-right"><div className={Data.params.Dfn.Negative ? "show-showhide" : "hide-showhide"}>
    <h1 className="h4 text-overflow"><a title="Disk usage display options" href={Data.params.Tlinks.Dfn} onClick={this.handleClick}><span className="showhide-hide whitespace-pre float-left">... </span>Disk usage</a>
    </h1></div>
  </div>
  <div className="col-lr large-11"><div className={Data.params.Dfn.Negative ? "show-showhide" : "hide-showhide"}>
    <div className="grid-block">
      <ul className="menu showhide-show">
        <li>
          <div className="input-group">
            <div className="input-group-label small text-nowrap">delay</div>
            <div className="input-group-button"><a className="button small secondary disabled">{Data.params.Dfd}</a></div>
            <div className="input-group-button">
              <a href={Data.params.Dlinks.Dfd.Less.Href} onClick={this.handleClick}
                 className={"text-nowrap button small " + (Data.params.Dlinks.Dfd.Less.ExtraClass != null ? Data.params.Dlinks.Dfd.Less.ExtraClass : "")}
                >- {Data.params.Dlinks.Dfd.Less.Text}</a>
            </div>
            <div className="input-group-button">
              <a href={Data.params.Dlinks.Dfd.More.Href} onClick={this.handleClick}
                 className={"text-nowrap button small " + (Data.params.Dlinks.Dfd.More.ExtraClass != null ? Data.params.Dlinks.Dfd.More.ExtraClass : "")}
                >{Data.params.Dlinks.Dfd.More.Text} +</a>
            </div>
          </div>
        </li>
        <li>
          <div className="input-group">
            <div className="input-group-label small text-nowrap">rows</div>
            <div className="input-group-button"><a className="button small secondary disabled">{Data.params.Dfn.Absolute}</a></div>
            <div className="input-group-button">
              <a href={Data.params.Nlinks.Dfn.Less.Href} onClick={this.handleClick}
                 className={"text-nowrap button small success " + (Data.params.Nlinks.Dfn.Less.ExtraClass != null ? Data.params.Nlinks.Dfn.Less.ExtraClass : "")}
                >- {Data.params.Nlinks.Dfn.Less.Text}</a>
            </div>
            <div className="input-group-button">
              <a href={Data.params.Nlinks.Dfn.More.Href} onClick={this.handleClick}
                 className={"text-nowrap button small success " + (Data.params.Nlinks.Dfn.More.ExtraClass != null ? Data.params.Nlinks.Dfn.More.ExtraClass : "")}
                >{Data.params.Nlinks.Dfn.More.Text} +</a>
            </div>
          </div>
        </li>
      </ul>
    </div><div className={Data.params.Dfn.Absolute == 0 ? "hide":""}>
    <div className="grid-block vertical stripe">
      <div className="grid-block thead"><span className="expand col small-1 text-nowrap"><a href={Data.params.Vlinks.Dfk[1-1].LinkHref} className={Data.params.Vlinks.Dfk[1-1].LinkClass} onClick={this.handleClick}
            >Device<span className={Data.params.Vlinks.Dfk[1-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap"><a href={Data.params.Vlinks.Dfk[2-1].LinkHref} className={Data.params.Vlinks.Dfk[2-1].LinkClass} onClick={this.handleClick}
            >Mounted<span className={Data.params.Vlinks.Dfk[2-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap text-right"><a href={Data.params.Vlinks.Dfk[6-1].LinkHref} className={Data.params.Vlinks.Dfk[6-1].LinkClass} onClick={this.handleClick}
            >Total<span className={Data.params.Vlinks.Dfk[6-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap text-right"><a href={Data.params.Vlinks.Dfk[5-1].LinkHref} className={Data.params.Vlinks.Dfk[5-1].LinkClass} onClick={this.handleClick}
            >Used<span className={Data.params.Vlinks.Dfk[5-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap text-right"><a href={Data.params.Vlinks.Dfk[3-1].LinkHref} className={Data.params.Vlinks.Dfk[3-1].LinkClass} onClick={this.handleClick}
            >Avail<span className={Data.params.Vlinks.Dfk[3-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap text-right"><a href={Data.params.Vlinks.Dfk[4-1].LinkHref} className={Data.params.Vlinks.Dfk[4-1].LinkClass} onClick={this.handleClick}
            >Use%<span className={Data.params.Vlinks.Dfk[4-1].CaretClass}></span></a></span><span className="expand col"></span></div>
      
      
      {this.List(Data).map(function($df, i) { return (
      <div className="grid-block" key={"df-rowby-dirname-"+$df.DirName}><span className="expand col small-1 text-overflow">{$df.DevName}</span><span className="expand col small-1 text-overflow"> {$df.DirName}</span><span className="expand col small-1 text-overflow text-right gray"><span className="float-right"> {$df.Total}</span><span title="Inodes total"> {$df.Inodes}</span></span><span className="expand col small-1 text-overflow text-right gray"><span className="float-right"> {$df.Used}</span><span title="Inodes used"> {$df.Iused}</span></span><span className="expand col small-1 text-overflow text-right gray"><span className="float-right"> {$df.Avail}</span><span title="Inodes free"> {$df.Ifree}</span></span><span className="expand col small-1 text-overflow text-right gray bg-usepct" data-usepct={$df.UsePct}><span className="float-right"> {$df.UsePct}%</span><span title="Inodes use%"> {$df.IusePct}%</span></span><span className="expand col-lr">{jsdefines.Sparkline({ref: i, col: 'UsePct'})}</span></div>);})}
      
    </div></div></div>
  </div>
</div>);
  }
});

jsdefines.define_hostname = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  Reduce: function(data) {
    return {
      hostname: data.hostname
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return (<a href="/" title={"hostname "+Data.hostname}>{Data.hostname}</a>);
  }
});

jsdefines.define_if = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  List: function(data) {
    let list;
    if (data == null || data["netio"] == null || (list = data["netio"].List) == null) {
      return [];
    }
    return list;
  },
  Reduce: function(data) {
    return {
      params: data.params,
      netio: data.netio
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return (
<div className="grid-block hr-top">
  <div className="col-lr large-1 text-right"><div className={Data.params.Ifn.Negative ? "show-showhide" : "hide-showhide"}>
    <h1 className="h4 text-overflow"><a title="Interfaces display options" href={Data.params.Tlinks.Ifn} onClick={this.handleClick}><span className="showhide-hide whitespace-pre float-left">... </span>Interfaces</a>
    </h1></div>
  </div>
  <div className="col-lr large-11"><div className={Data.params.Ifn.Negative ? "show-showhide" : "hide-showhide"}>
    <div className="grid-block">
      <ul className="menu showhide-show">
        <li>
          <div className="input-group">
            <div className="input-group-label small text-nowrap">delay</div>
            <div className="input-group-button"><a className="button small secondary disabled">{Data.params.Ifd}</a></div>
            <div className="input-group-button">
              <a href={Data.params.Dlinks.Ifd.Less.Href} onClick={this.handleClick}
                 className={"text-nowrap button small " + (Data.params.Dlinks.Ifd.Less.ExtraClass != null ? Data.params.Dlinks.Ifd.Less.ExtraClass : "")}
                >- {Data.params.Dlinks.Ifd.Less.Text}</a>
            </div>
            <div className="input-group-button">
              <a href={Data.params.Dlinks.Ifd.More.Href} onClick={this.handleClick}
                 className={"text-nowrap button small " + (Data.params.Dlinks.Ifd.More.ExtraClass != null ? Data.params.Dlinks.Ifd.More.ExtraClass : "")}
                >{Data.params.Dlinks.Ifd.More.Text} +</a>
            </div>
          </div>
        </li>
        <li>
          <div className="input-group">
            <div className="input-group-label small text-nowrap">rows</div>
            <div className="input-group-button"><a className="button small secondary disabled">{Data.params.Ifn.Absolute}</a></div>
            <div className="input-group-button">
              <a href={Data.params.Nlinks.Ifn.Less.Href} onClick={this.handleClick}
                 className={"text-nowrap button small success " + (Data.params.Nlinks.Ifn.Less.ExtraClass != null ? Data.params.Nlinks.Ifn.Less.ExtraClass : "")}
                >- {Data.params.Nlinks.Ifn.Less.Text}</a>
            </div>
            <div className="input-group-button">
              <a href={Data.params.Nlinks.Ifn.More.Href} onClick={this.handleClick}
                 className={"text-nowrap button small success " + (Data.params.Nlinks.Ifn.More.ExtraClass != null ? Data.params.Nlinks.Ifn.More.ExtraClass : "")}
                >{Data.params.Nlinks.Ifn.More.Text} +</a>
            </div>
          </div>
        </li>
      </ul>
    </div><div className={Data.params.Ifn.Absolute == 0 ? "hide":""}>
    <div className="grid-block vertical stripe">
      <div className="grid-block thead"><span className="expand col small-1">Interface</span><span className="expand col small-1 text-right"> IP</span><span title="Drops,Errors In/Out per second" className="expand col small-2 text-right text-nowrap"> Loss IO ps</span><span title="Packets In/Out per second" className="expand col small-2 text-right text-nowrap"> Packets IO ps</span><span title="Bits In/Out per second" className="expand col small-2 text-right text-nowrap"> IO <i>b</i>ps</span><span className="expand col"></span></div>
      
      
      {this.List(Data).map(function($if, i) { return (
      <div className="grid-block" key={"if-rowby-name-"+$if.Name}><span className="expand col small-1 text-overflow">{$if.Name}</span><span className="expand col small-1 text-overflow text-right">{$if.IP}</span><span className="expand col small-2 text-right text-nowrap">&nbsp;<span title="Total drops,errors modulo 4G" className="gray"><span title="Total drops In modulo 4G">{$if.DropsIn}</span><span className={$if.DropsOut == "-1" ? "hide":""}>/</span><span className={$if.DropsOut == "-1" ? "hide":""} title="Total drops Out modulo 4G">{$if.DropsOut}</span>,<span title="Total errors In modulo 4G">{$if.ErrorsIn}</span>/<span title="Total errors Out modulo 4G">{$if.ErrorsOut}</span></span>&nbsp;<span className={($if.DeltaDropsIn == "0" && ($if.DeltaDropsOut == "-1" || $if.DeltaDropsOut == "0") && $if.DeltaErrorsIn == "0" && $if.DeltaErrorsOut == "0") ? "gray":""}><span title="Drops In per second">{$if.DeltaDropsIn}</span><span className={$if.DeltaDropsOut == "-1" ? "hide":""}>/</span><span className={$if.DeltaDropsOut == "-1" ? "hide":""} title="Drops Out per second">{$if.DeltaDropsOut}</span>,<span title="Errors In per second">{$if.DeltaErrorsIn}</span>/<span title="Errors Out per second">{$if.DeltaErrorsOut}</span></span></span><span className="expand col small-2 text-right text-nowrap">&nbsp;<span className="gray"><span title="Total packets In modulo 4G">{$if.PacketsIn}</span>/<span title="Total packets Out modulo 4G">{$if.PacketsOut}</span></span>&nbsp;<span title="Packets In per second">{$if.DeltaPacketsIn}</span>/<span title="Packets Out per second">{$if.DeltaPacketsOut}</span></span><span className="expand col small-2 text-right text-nowrap">&nbsp;<span className="gray"><span title="Total BYTES In modulo 4G">{$if.BytesIn}</span>/<span title="Total BYTES Out modulo 4G">{$if.BytesOut}</span></span>&nbsp;<span title="BITS In per second">{$if.DeltaBitsIn}</span>/<span title="BITS Out per second">{$if.DeltaBitsOut}</span></span><span className="expand col-lr">{jsdefines.Sparkline({ref: i, col: 'DeltaBytesOutNum'})}</span></div>);})}
      
    </div></div></div>
  </div>
</div>);
  }
});

jsdefines.define_la = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  List: function(data) {
    let list;
    if (data == null || data["la"] == null || (list = data["la"].List) == null) {
      return [];
    }
    return list;
  },
  Reduce: function(data) {
    return {
      params: data.params,
      la: data.la
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return (
<div className="grid-block hr-top">
  <div className="col-lr large-1 text-right"><div className={Data.params.Lan.Negative ? "show-showhide" : "hide-showhide"}>
    <h1 className="h4 text-overflow"><a title="Load avg display options" href={Data.params.Tlinks.Lan} onClick={this.handleClick}><span className="showhide-hide whitespace-pre float-left">... </span>Load avg</a>
    </h1></div>
  </div>
  <div className="col-lr large-11"><div className={Data.params.Lan.Negative ? "show-showhide" : "hide-showhide"}>
    <div className="grid-block">
      <ul className="menu showhide-show">
        <li>
          <div className="input-group">
            <div className="input-group-label small text-nowrap">delay</div>
            <div className="input-group-button"><a className="button small secondary disabled">{Data.params.Lad}</a></div>
            <div className="input-group-button">
              <a href={Data.params.Dlinks.Lad.Less.Href} onClick={this.handleClick}
                 className={"text-nowrap button small " + (Data.params.Dlinks.Lad.Less.ExtraClass != null ? Data.params.Dlinks.Lad.Less.ExtraClass : "")}
                >- {Data.params.Dlinks.Lad.Less.Text}</a>
            </div>
            <div className="input-group-button">
              <a href={Data.params.Dlinks.Lad.More.Href} onClick={this.handleClick}
                 className={"text-nowrap button small " + (Data.params.Dlinks.Lad.More.ExtraClass != null ? Data.params.Dlinks.Lad.More.ExtraClass : "")}
                >{Data.params.Dlinks.Lad.More.Text} +</a>
            </div>
          </div>
        </li>
        <li>
          <div className="input-group">
            <div className="input-group-label small text-nowrap">rows</div>
            <div className="input-group-button"><a className="button small secondary disabled">{Data.params.Lan.Absolute}</a></div>
            <div className="input-group-button">
              <a href={Data.params.Nlinks.Lan.Less.Href} onClick={this.handleClick}
                 className={"text-nowrap button small success " + (Data.params.Nlinks.Lan.Less.ExtraClass != null ? Data.params.Nlinks.Lan.Less.ExtraClass : "")}
                >- {Data.params.Nlinks.Lan.Less.Text}</a>
            </div>
            <div className="input-group-button">
              <a href={Data.params.Nlinks.Lan.More.Href} onClick={this.handleClick}
                 className={"text-nowrap button small success " + (Data.params.Nlinks.Lan.More.ExtraClass != null ? Data.params.Nlinks.Lan.More.ExtraClass : "")}
                >{Data.params.Nlinks.Lan.More.Text} +</a>
            </div>
          </div>
        </li>
      </ul>
    </div><div className={Data.params.Lan.Absolute == 0 ? "hide":""}>
    <div className="grid-block vertical stripe">
      <div className="grid-block thead"><span className="expand col small-1">Period</span><span className="expand col small-1 text-right"> Value</span><span className="expand col"></span></div>
      
      
      {this.List(Data).map(function($la, i) { return (
      <div className="grid-block" key={"la-rowby-period-"+$la.Period}><span className="expand col small-1">
        <div className="text-right width-3rem">{$la.Period}m</div></span><span className="expand col small-1 text-right"> {$la.Value}</span><span className="expand col-lr">{jsdefines.Sparkline({ref: i, col: 'Value'})}</span></div>);})}
      
    </div></div></div>
  </div>
</div>);
  }
});

jsdefines.define_mem = React.createClass({
  mixins: [ReactPureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
  List: function(data) {
    let list;
    if (data == null || data["mem"] == null || (list = data["mem"].List) == null) {
      return [];
    }
    return list;
  },
  Reduce: function(data) {
    return {
      params: data.params,
      mem: data.mem
    };
  },
  render: function() {
    let Data = this.state; // shadow global Data
    return (
<div className="grid-block hr-top">
  <div className="col-lr large-1 text-right"><div className={Data.params.Memn.Negative ? "show-showhide" : "hide-showhide"}>
    <h1 className="h4 text-overflow"><a title="Memory display options" href={Data.params.Tlinks.Memn} onClick={this.handleClick}><span className="showhide-hide whitespace-pre float-left">... </span>Memory</a>
    </h1></div>
  </div>
  <div className="col-lr large-11"><div className={Data.params.Memn.Negative ? "show-showhide" : "hide-showhide"}>
    <div className="grid-block">
      <ul className="menu showhide-show">
        <li>
          <div className="input-group">
            <div className="input-group-label small text-nowrap">delay</div>
            <div className="input-group-button"><a className="button small secondary disabled">{Data.params.Memd}</a></div>
            <div className="input-group-button">
              <a href={Data.params.Dlinks.Memd.Less.Href} onClick={this.handleClick}
                 className={"text-nowrap button small " + (Data.params.Dlinks.Memd.Less.ExtraClass != null ? Data.params.Dlinks.Memd.Less.ExtraClass : "")}
                >- {Data.params.Dlinks.Memd.Less.Text}</a>
            </div>
            <div className="input-group-button">
              <a href={Data.params.Dlinks.Memd.More.Href} onClick={this.handleClick}
                 className={"text-nowrap button small " + (Data.params.Dlinks.Memd.More.ExtraClass != null ? Data.params.Dlinks.Memd.More.ExtraClass : "")}
                >{Data.params.Dlinks.Memd.More.Text} +</a>
            </div>
          </div>
        </li>
        <li>
          <div className="input-group">
            <div className="input-group-label small text-nowrap">rows</div>
            <div className="input-group-button"><a className="button small secondary disabled">{Data.params.Memn.Absolute}</a></div>
            <div className="input-group-button">
              <a href={Data.params.Nlinks.Memn.Less.Href} onClick={this.handleClick}
                 className={"text-nowrap button small success " + (Data.params.Nlinks.Memn.Less.ExtraClass != null ? Data.params.Nlinks.Memn.Less.ExtraClass : "")}
                >- {Data.params.Nlinks.Memn.Less.Text}</a>
            </div>
            <div className="input-group-button">
              <a href={Data.params.Nlinks.Memn.More.Href} onClick={this.handleClick}
                 className={"text-nowrap button small success " + (Data.params.Nlinks.Memn.More.ExtraClass != null ? Data.params.Nlinks.Memn.More.ExtraClass : "")}
                >{Data.params.Nlinks.Memn.More.Text} +</a>
            </div>
          </div>
        </li>
      </ul>
    </div><div className={Data.params.Memn.Absolute == 0 ? "hide":""}>
    <div className="grid-block vertical stripe">
      <div className="grid-block thead"><span className="expand col small-1">Memory</span><span className="expand col small-1 text-right"> Total</span><span className="expand col small-1 text-right"> Used</span><span className="expand col small-1 text-right"> Free</span><span className="expand col small-1 text-right"> Use%</span><span className="expand col"></span></div>
      
      
      {this.List(Data).map(function($mem, i) { return (
      <div className="grid-block" key={"mem-rowby-kind-"+$mem.Kind}><span className="expand col small-1">{$mem.Kind}</span><span className="expand col small-1 text-right"> {$mem.Total}</span><span className="expand col small-1 text-right"> {$mem.Used}</span><span className="expand col small-1 text-right"> {$mem.Free}</span><span className="expand col small-1 text-right bg-usepct" data-usepct={$mem.UsePct}> {$mem.UsePct}%</span><span className="expand col-lr">{jsdefines.Sparkline({ref: i, col: 'UsePct'})}</span></div>);})}
      
    </div></div></div>
  </div>
</div>);
  }
});

jsdefines.define_ps = React.createClass({
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
    return (
<div className="grid-block hr-top">
  <div className="col-lr large-1 text-right"><div className={Data.params.Psn.Negative ? "show-showhide" : "hide-showhide"}>
    <h1 className="h4 text-overflow"><a title="Processes display options" href={Data.params.Tlinks.Psn} onClick={this.handleClick}><span className="showhide-hide whitespace-pre float-left">... </span>Processes</a>
    </h1></div>
  </div>
  <div className="col-lr large-11"><div className={Data.params.Psn.Negative ? "show-showhide" : "hide-showhide"}>
    <div className="grid-block">
      <ul className="menu showhide-show">
        <li>
          <div className="input-group">
            <div className="input-group-label small text-nowrap">delay</div>
            <div className="input-group-button"><a className="button small secondary disabled">{Data.params.Psd}</a></div>
            <div className="input-group-button">
              <a href={Data.params.Dlinks.Psd.Less.Href} onClick={this.handleClick}
                 className={"text-nowrap button small " + (Data.params.Dlinks.Psd.Less.ExtraClass != null ? Data.params.Dlinks.Psd.Less.ExtraClass : "")}
                >- {Data.params.Dlinks.Psd.Less.Text}</a>
            </div>
            <div className="input-group-button">
              <a href={Data.params.Dlinks.Psd.More.Href} onClick={this.handleClick}
                 className={"text-nowrap button small " + (Data.params.Dlinks.Psd.More.ExtraClass != null ? Data.params.Dlinks.Psd.More.ExtraClass : "")}
                >{Data.params.Dlinks.Psd.More.Text} +</a>
            </div>
          </div>
        </li>
        <li>
          <div className="input-group">
            <div className="input-group-label small text-nowrap">rows</div>
            <div className="input-group-button"><a className="button small secondary disabled">{Data.params.Psn.Absolute}</a></div>
            <div className="input-group-button">
              <a href={Data.params.Nlinks.Psn.Less.Href} onClick={this.handleClick}
                 className={"text-nowrap button small success " + (Data.params.Nlinks.Psn.Less.ExtraClass != null ? Data.params.Nlinks.Psn.Less.ExtraClass : "")}
                >- {Data.params.Nlinks.Psn.Less.Text}</a>
            </div>
            <div className="input-group-button">
              <a href={Data.params.Nlinks.Psn.More.Href} onClick={this.handleClick}
                 className={"text-nowrap button small success " + (Data.params.Nlinks.Psn.More.ExtraClass != null ? Data.params.Nlinks.Psn.More.ExtraClass : "")}
                >{Data.params.Nlinks.Psn.More.Text} +</a>
            </div>
          </div>
        </li>
      </ul>
    </div><div className={Data.params.Psn.Absolute == 0 ? "hide":""}>
    <div className="grid-block vertical stripe">
      <div className="grid-block thead"><span className="expand col small-1 text-nowrap text-right"><a href={Data.params.Vlinks.Psk[1-1].LinkHref} className={Data.params.Vlinks.Psk[1-1].LinkClass} onClick={this.handleClick}
            >PID<span className={Data.params.Vlinks.Psk[1-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap text-right"><a href={Data.params.Vlinks.Psk[2-1].LinkHref} className={Data.params.Vlinks.Psk[2-1].LinkClass} onClick={this.handleClick}
            >UID<span className={Data.params.Vlinks.Psk[2-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap"><a href={Data.params.Vlinks.Psk[3-1].LinkHref} className={Data.params.Vlinks.Psk[3-1].LinkClass} onClick={this.handleClick}
            >USER<span className={Data.params.Vlinks.Psk[3-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap text-right"><a href={Data.params.Vlinks.Psk[4-1].LinkHref} className={Data.params.Vlinks.Psk[4-1].LinkClass} onClick={this.handleClick}
            >PR<span className={Data.params.Vlinks.Psk[4-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap text-right"><a href={Data.params.Vlinks.Psk[5-1].LinkHref} className={Data.params.Vlinks.Psk[5-1].LinkClass} onClick={this.handleClick}
            >NI<span className={Data.params.Vlinks.Psk[5-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap text-right"><a href={Data.params.Vlinks.Psk[6-1].LinkHref} className={Data.params.Vlinks.Psk[6-1].LinkClass} onClick={this.handleClick}
            >VIRT<span className={Data.params.Vlinks.Psk[6-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap text-right"><a href={Data.params.Vlinks.Psk[7-1].LinkHref} className={Data.params.Vlinks.Psk[7-1].LinkClass} onClick={this.handleClick}
            >RES<span className={Data.params.Vlinks.Psk[7-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap text-center"><a href={Data.params.Vlinks.Psk[8-1].LinkHref} className={Data.params.Vlinks.Psk[8-1].LinkClass} onClick={this.handleClick}
            >TIME<span className={Data.params.Vlinks.Psk[8-1].CaretClass}></span></a></span><span className="expand col small-1 text-nowrap"><a href={Data.params.Vlinks.Psk[9-1].LinkHref} className={Data.params.Vlinks.Psk[9-1].LinkClass} onClick={this.handleClick}
            >COMMAND<span className={Data.params.Vlinks.Psk[9-1].CaretClass}></span></a></span></div>
      
      
      {this.List(Data).map(function($ps, i) { return (
      <div className="grid-block" key={"ps-rowby-pid-"+$ps.PID}><span className="expand col small-1 text-right">{$ps.PID}</span><span className="expand col small-1 text-right"> {$ps.UID}</span><span className="expand col small-1"> {$ps.User}</span><span className="expand col small-1 text-right"> {$ps.Priority}</span><span className="expand col small-1 text-right"> {$ps.Nice}</span><span className="expand col small-1 text-right"> {$ps.Size}</span><span className="expand col small-1 text-right"> {$ps.Resident}</span><span className="expand col small-1 text-center"> {$ps.Time}</span><span className="expand col"> {$ps.Name}</span></div>);})}
      
    </div></div></div>
  </div>
</div>);
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
    return (<span>{Data.uptime}</span>);
  }
});


module.exports = jsdefines;

// Local variables:
// js-indent-level: 2
// js2-basic-offset: 2
// End:
