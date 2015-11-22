define(function(require) {
  let React = require('react');
  let $     = require('jquery');
  let jsdefines = {};
  jsdefines.StateHandlingMixin = { // requires .Reduce method
    getInitialState: function() {
      return this.StateFrom(Data); // global Data
    },
    NewState: function(data) {
      let state = this.StateFrom(data);
      if (state != null) {
        this.setState(state);
      }
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
        href = $(e.target).parent().get(0).getAttribute('href');
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
    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    Reduce: function(data) {
      return {
        hostname: data.hostname
      };
    },
    render: function() {
      let Data = this.state; // shadow global Data
      return <a  title={"hostname " + Data.hostname} href="/"
  ><h4 className="clip12 margin-bottom-0"
    >{Data.hostname}</h4
  ></a
>;
    }
  });

  jsdefines.define_panelcpu = React.createClass({
    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    List: function(data) {
      let list;
      if (data != null && data["cpu"] != null && (list = data["cpu"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function(data) {
      return {
        params: data.params,
        cpu: data.cpu
      };
    },
    render: function() {
      let Data = this.state; // shadow global Data
      return <div
  ><div  className={!Data.params.CPUn.Negative ? "tabs tabs-border bar-less" : "tabs tabs-border"} data-tabs
    ><div className="tabs-title menu-tab-padding"
      ><a  href={Data.params.Tlinks.CPUn} onClick={this.handleClick}
        ><h5 className="margin-bottom-0"
          >CPU</h5
        ></a
      ></div
    ><ul className="float-left bar menu"
      ><li className="menu-text"
        ><div className="input-group margin-bottom-0"
          ><span className="input-group-label"
            >delay</span
          ><span className="input-group-label label secondary"
            >{Data.params.CPUd}</span
          ><a href={Data.params.Dlinks.CPUd.Less.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.CPUd.Less.ExtraClass != null ? Data.params.Dlinks.CPUd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Dlinks.CPUd.Less.Text}</a
><a href={Data.params.Dlinks.CPUd.More.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.CPUd.More.ExtraClass != null ? Data.params.Dlinks.CPUd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Dlinks.CPUd.More.Text} +</a
></div
        ></li
      ><li className="menu-text"
        ><div className="input-group margin-bottom-0"
          ><span className="input-group-label"
            >rows</span
          ><span className="input-group-label label secondary"
            >{Data.params.CPUn.Absolute}</span
          ><a href={Data.params.Nlinks.CPUn.Less.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.CPUn.Less.ExtraClass != null ? Data.params.Nlinks.CPUn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Nlinks.CPUn.Less.Text}</a
><a href={Data.params.Nlinks.CPUn.More.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.CPUn.More.ExtraClass != null ? Data.params.Nlinks.CPUn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Nlinks.CPUn.More.Text} +</a
></div
        ></li
      ></ul
    ></div
  ><table  className={Data.params.CPUn.Absolute != 0 ? "hover scroll-x margin-bottom-0" : "hide"}
    ><thead
      ><tr
        ><th
          ></th
        ><th className="text-right"
          >User%</th
        ><th className="text-right"
          >Sys%</th
        ><th className="text-right"
          >Wait%</th
        ><th className="text-right"
          >Idle%</th
        ></tr
      ></thead
    ><tbody
      >{this.List(Data).map(function($cpu) { return<tr  key={"cpu-rowby-N-"+$cpu.N}
        ><td className="text-right text-nowrap"
          >{$cpu.N}</td
        ><td className="text-right bg-usepct"
  data-usepct={$cpu.UserPct}
          >{$cpu.UserPct}%</td
        ><td className="text-right bg-usepct"
  data-usepct={$cpu.SysPct}
          >{$cpu.SysPct}%</td
        ><td className="text-right bg-usepct"
  data-usepct={$cpu.WaitPct}
          >{$cpu.WaitPct}%</td
        ><td className="text-right bg-usepct-inverse"
  data-usepct={$cpu.IdlePct}
          >{$cpu.IdlePct}%</td
        ></tr
      >})}</tbody
    ></table
  ></div
>;
    }
  });

  jsdefines.define_paneldf = React.createClass({
    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    List: function(data) {
      let list;
      if (data != null && data["diskUsage"] != null && (list = data["diskUsage"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function(data) {
      return {
        params: data.params,
        diskUsage: data.diskUsage
      };
    },
    render: function() {
      let Data = this.state; // shadow global Data
      return <div
  ><div  className={!Data.params.Dfn.Negative ? "tabs tabs-border bar-less" : "tabs tabs-border"} data-tabs
    ><div className="tabs-title menu-tab-padding"
      ><a  href={Data.params.Tlinks.Dfn} onClick={this.handleClick}
        ><h5 className="margin-bottom-0"
          >Disk usage</h5
        ></a
      ></div
    ><ul className="float-left bar menu"
      ><li className="menu-text"
        ><div className="input-group margin-bottom-0"
          ><span className="input-group-label"
            >delay</span
          ><span className="input-group-label label secondary"
            >{Data.params.Dfd}</span
          ><a href={Data.params.Dlinks.Dfd.Less.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Dfd.Less.ExtraClass != null ? Data.params.Dlinks.Dfd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Dlinks.Dfd.Less.Text}</a
><a href={Data.params.Dlinks.Dfd.More.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Dfd.More.ExtraClass != null ? Data.params.Dlinks.Dfd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Dlinks.Dfd.More.Text} +</a
></div
        ></li
      ><li className="menu-text"
        ><div className="input-group margin-bottom-0"
          ><span className="input-group-label"
            >rows</span
          ><span className="input-group-label label secondary"
            >{Data.params.Dfn.Absolute}</span
          ><a href={Data.params.Nlinks.Dfn.Less.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Dfn.Less.ExtraClass != null ? Data.params.Nlinks.Dfn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Nlinks.Dfn.Less.Text}</a
><a href={Data.params.Nlinks.Dfn.More.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Dfn.More.ExtraClass != null ? Data.params.Nlinks.Dfn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Nlinks.Dfn.More.Text} +</a
></div
        ></li
      ></ul
    ></div
  ><table  className={Data.params.Dfn.Absolute != 0 ? "hover scroll-x margin-bottom-0" : "hide"}
    ><thead
      ><tr className="text-nowrap"
        ><th className="header "
  ><a href={Data.params.Vlinks.Dfk[1-1].LinkHref} className={Data.params.Vlinks.Dfk[1-1].LinkClass} onClick={this.handleClick}  
    >Device<span className={Data.params.Vlinks.Dfk[1-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.params.Vlinks.Dfk[2-1].LinkHref} className={Data.params.Vlinks.Dfk[2-1].LinkClass} onClick={this.handleClick}  
    >Mounted<span className={Data.params.Vlinks.Dfk[2-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.params.Vlinks.Dfk[3-1].LinkHref} className={Data.params.Vlinks.Dfk[3-1].LinkClass} onClick={this.handleClick}  
    >Avail<span className={Data.params.Vlinks.Dfk[3-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.params.Vlinks.Dfk[4-1].LinkHref} className={Data.params.Vlinks.Dfk[4-1].LinkClass} onClick={this.handleClick}  
    >Use%<span className={Data.params.Vlinks.Dfk[4-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.params.Vlinks.Dfk[5-1].LinkHref} className={Data.params.Vlinks.Dfk[5-1].LinkClass} onClick={this.handleClick}  
    >Used<span className={Data.params.Vlinks.Dfk[5-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.params.Vlinks.Dfk[6-1].LinkHref} className={Data.params.Vlinks.Dfk[6-1].LinkClass} onClick={this.handleClick}  
    >Total<span className={Data.params.Vlinks.Dfk[6-1].CaretClass}
      ></span
    ></a
  ></th
></tr
      ></thead
    ><tbody
      >{this.List(Data).map(function($df) { return<tr  key={"df-rowby-dirname-"+$df.DirName}
        >  <td className="text-nowrap"
          >{$df.DevName}</td
        >  <td className="text-nowrap"
          >{$df.DirName}</td
        ><td className="text-right text-nowrap"
          ><span className="mutext" title="Inodes free"
            >{$df.Ifree}</span
          > {$df.Avail}</td
        ><td className="text-right bg-usepct text-nowrap" data-usepct={$df.UsePct}
          ><span className="mutext" title="Inodes use%"
            >{$df.IusePct}%</span
          > {$df.UsePct}%</td
        ><td className="text-right text-nowrap"
          ><span className="mutext" title="Inodes used"
            >{$df.Iused}</span
          > {$df.Used}</td
        ><td className="text-right text-nowrap"
          ><span className="mutext" title="Inodes total"
            >{$df.Inodes}</span
          > {$df.Total}</td
        ></tr
      >})}</tbody
    ></table
  ></div
>;
    }
  });

  jsdefines.define_panelif = React.createClass({
    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    List: function(data) {
      let list;
      if (data != null && data["ifaddrs"] != null && (list = data["ifaddrs"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function(data) {
      return {
        params: data.params,
        ifaddrs: data.ifaddrs
      };
    },
    render: function() {
      let Data = this.state; // shadow global Data
      return <div
  ><div  className={!Data.params.Ifn.Negative ? "tabs tabs-border bar-less" : "tabs tabs-border"} data-tabs
    ><div className="tabs-title menu-tab-padding"
      ><a  href={Data.params.Tlinks.Ifn} onClick={this.handleClick}
        ><h5 className="margin-bottom-0"
          >Interfaces</h5
        ></a
      ></div
    ><ul className="float-left bar menu"
      ><li className="menu-text"
        ><div className="input-group margin-bottom-0"
          ><span className="input-group-label"
            >delay</span
          ><span className="input-group-label label secondary"
            >{Data.params.Ifd}</span
          ><a href={Data.params.Dlinks.Ifd.Less.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Ifd.Less.ExtraClass != null ? Data.params.Dlinks.Ifd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Dlinks.Ifd.Less.Text}</a
><a href={Data.params.Dlinks.Ifd.More.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Ifd.More.ExtraClass != null ? Data.params.Dlinks.Ifd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Dlinks.Ifd.More.Text} +</a
></div
        ></li
      ><li className="menu-text"
        ><div className="input-group margin-bottom-0"
          ><span className="input-group-label"
            >rows</span
          ><span className="input-group-label label secondary"
            >{Data.params.Ifn.Absolute}</span
          ><a href={Data.params.Nlinks.Ifn.Less.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Ifn.Less.ExtraClass != null ? Data.params.Nlinks.Ifn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Nlinks.Ifn.Less.Text}</a
><a href={Data.params.Nlinks.Ifn.More.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Ifn.More.ExtraClass != null ? Data.params.Nlinks.Ifn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Nlinks.Ifn.More.Text} +</a
></div
        ></li
      ></ul
    ></div
  ><table  className={Data.params.Ifn.Absolute != 0 ? "hover scroll-x margin-bottom-0" : "hide"}
    ><thead
      ><tr
        ><th
          >Interface</th
        ><th className="text-right"
          >IP</th
        ><th className="text-right text-nowrap col-md-3" title="Bits In/Out per second"
          >IO <i
            >b</i
          >ps</th
        ><th className="text-right text-nowrap col-md-3" title="Packets In/Out per second"
          >Packets IO ps</th
        ><th className="text-right text-nowrap col-md-3" title="Drops,Errors In/Out per second"
          >Loss IO ps</th
        ></tr
      ></thead
    ><tbody
      >{this.List(Data).map(function($if) { return<tr  key={"if-rowby-name-"+$if.Name}
        ><td className="text-nowrap"
          >{$if.Name}</td
        ><td className="text-right"
          >{$if.IP}</td
        ><td className="text-right text-nowrap"
          ><span className="mutext"
            ><span title="Total BYTES In modulo 4G"
              >{$if.BytesIn}</span
            >/<span title="Total BYTES Out modulo 4G"
              >{$if.BytesOut}</span
            ></span
          > <span title="BITS In per second"
            >{$if.DeltaBitsIn}</span
          >/<span title="BITS Out per second"
            >{$if.DeltaBitsOut}</span
          ></td
        ><td className="text-right text-nowrap"
          ><span className="mutext"
            ><span title="Total packets In modulo 4G"
              >{$if.PacketsIn}</span
            >/<span title="Total packets Out modulo 4G"
              >{$if.PacketsOut}</span
            ></span
          > <span title="Packets In per second"
            >{$if.DeltaPacketsIn}</span
          >/<span title="Packets Out per second"
            >{$if.DeltaPacketsOut}</span
          ></td
        ><td className="text-right text-nowrap"
          ><span className="mutext" title="Total drops,errors modulo 4G"
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
          ></td
        ></tr
      >})}</tbody
    ></table
  ></div
>;
    }
  });

  jsdefines.define_panelmem = React.createClass({
    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    List: function(data) {
      let list;
      if (data != null && data["memory"] != null && (list = data["memory"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function(data) {
      return {
        params: data.params,
        memory: data.memory
      };
    },
    render: function() {
      let Data = this.state; // shadow global Data
      return <div
  ><div  className={!Data.params.Memn.Negative ? "tabs tabs-border bar-less" : "tabs tabs-border"} data-tabs
    ><div className="tabs-title menu-tab-padding"
      ><a  href={Data.params.Tlinks.Memn} onClick={this.handleClick}
        ><h5 className="margin-bottom-0"
          >Memory</h5
        ></a
      ></div
    ><ul className="float-left bar menu"
      ><li className="menu-text"
        ><div className="input-group margin-bottom-0"
          ><span className="input-group-label"
            >delay</span
          ><span className="input-group-label label secondary"
            >{Data.params.Memd}</span
          ><a href={Data.params.Dlinks.Memd.Less.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Memd.Less.ExtraClass != null ? Data.params.Dlinks.Memd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Dlinks.Memd.Less.Text}</a
><a href={Data.params.Dlinks.Memd.More.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Memd.More.ExtraClass != null ? Data.params.Dlinks.Memd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Dlinks.Memd.More.Text} +</a
></div
        ></li
      ><li className="menu-text"
        ><div className="input-group margin-bottom-0"
          ><span className="input-group-label"
            >rows</span
          ><span className="input-group-label label secondary"
            >{Data.params.Memn.Absolute}</span
          ><a href={Data.params.Nlinks.Memn.Less.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Memn.Less.ExtraClass != null ? Data.params.Nlinks.Memn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Nlinks.Memn.Less.Text}</a
><a href={Data.params.Nlinks.Memn.More.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Memn.More.ExtraClass != null ? Data.params.Nlinks.Memn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Nlinks.Memn.More.Text} +</a
></div
        ></li
      ></ul
    ></div
  ><table  className={Data.params.Memn.Absolute != 0 ? "hover scroll-x margin-bottom-0" : "hide"}
    ><thead
      ><tr
        ><th
          ></th
        ><th className="text-right"
          >Free</th
        ><th className="text-right"
          >Use%</th
        ><th className="text-right"
          >Used</th
        ><th className="text-right"
          >Total</th
        ></tr
      ></thead
    ><tbody
      >{this.List(Data).map(function($mem) { return<tr  key={"mem-rowby-kind-"+$mem.Kind}
        ><td
          >{$mem.Kind}</td
        ><td className="text-right"
          >{$mem.Free}</td
        ><td className="text-right bg-usepct" data-usepct={$mem.UsePct}
          >{$mem.UsePct}%</td
        ><td className="text-right"
          >{$mem.Used}</td
        ><td className="text-right"
          >{$mem.Total}</td
        ></tr
      >})}</tbody
    ></table
  ></div
>;
    }
  });

  jsdefines.define_panelps = React.createClass({
    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    List: function(data) {
      let list;
      if (data != null && data["procs"] != null && (list = data["procs"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function(data) {
      return {
        params: data.params,
        procs: data.procs
      };
    },
    render: function() {
      let Data = this.state; // shadow global Data
      return <div
  ><div  className={!Data.params.Psn.Negative ? "tabs tabs-border bar-less" : "tabs tabs-border"} data-tabs
    ><div className="tabs-title menu-tab-padding"
      ><a  href={Data.params.Tlinks.Psn} onClick={this.handleClick}
        ><h5 className="margin-bottom-0"
          >Processes</h5
        ></a
      ></div
    ><ul className="float-left bar menu"
      ><li className="menu-text"
        ><div className="input-group margin-bottom-0"
          ><span className="input-group-label"
            >delay</span
          ><span className="input-group-label label secondary"
            >{Data.params.Psd}</span
          ><a href={Data.params.Dlinks.Psd.Less.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Psd.Less.ExtraClass != null ? Data.params.Dlinks.Psd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Dlinks.Psd.Less.Text}</a
><a href={Data.params.Dlinks.Psd.More.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Psd.More.ExtraClass != null ? Data.params.Dlinks.Psd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Dlinks.Psd.More.Text} +</a
></div
        ></li
      ><li className="menu-text"
        ><div className="input-group margin-bottom-0"
          ><span className="input-group-label"
            >rows</span
          ><span className="input-group-label label secondary"
            >{Data.params.Psn.Absolute}</span
          ><a href={Data.params.Nlinks.Psn.Less.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Psn.Less.ExtraClass != null ? Data.params.Nlinks.Psn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.params.Nlinks.Psn.Less.Text}</a
><a href={Data.params.Nlinks.Psn.More.Href} className={"button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Psn.More.ExtraClass != null ? Data.params.Nlinks.Psn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.params.Nlinks.Psn.More.Text} +</a
></div
        ></li
      ></ul
    ></div
  ><table  className={Data.params.Psn.Absolute != 0 ? "hover scroll-x margin-bottom-0" : "hide"}
    ><thead
      ><tr className="text-nowrap"
        ><th className="header text-right"
  ><a href={Data.params.Vlinks.Psk[1-1].LinkHref} className={Data.params.Vlinks.Psk[1-1].LinkClass} onClick={this.handleClick}  
    >PID<span className={Data.params.Vlinks.Psk[1-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.params.Vlinks.Psk[2-1].LinkHref} className={Data.params.Vlinks.Psk[2-1].LinkClass} onClick={this.handleClick}  
    >UID<span className={Data.params.Vlinks.Psk[2-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.params.Vlinks.Psk[3-1].LinkHref} className={Data.params.Vlinks.Psk[3-1].LinkClass} onClick={this.handleClick}  
    >USER<span className={Data.params.Vlinks.Psk[3-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.params.Vlinks.Psk[4-1].LinkHref} className={Data.params.Vlinks.Psk[4-1].LinkClass} onClick={this.handleClick}  
    >PR<span className={Data.params.Vlinks.Psk[4-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.params.Vlinks.Psk[5-1].LinkHref} className={Data.params.Vlinks.Psk[5-1].LinkClass} onClick={this.handleClick}  
    >NI<span className={Data.params.Vlinks.Psk[5-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.params.Vlinks.Psk[6-1].LinkHref} className={Data.params.Vlinks.Psk[6-1].LinkClass} onClick={this.handleClick}  
    >VIRT<span className={Data.params.Vlinks.Psk[6-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.params.Vlinks.Psk[7-1].LinkHref} className={Data.params.Vlinks.Psk[7-1].LinkClass} onClick={this.handleClick}  
    >RES<span className={Data.params.Vlinks.Psk[7-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-center"
  ><a href={Data.params.Vlinks.Psk[8-1].LinkHref} className={Data.params.Vlinks.Psk[8-1].LinkClass} onClick={this.handleClick}  
    >TIME<span className={Data.params.Vlinks.Psk[8-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.params.Vlinks.Psk[9-1].LinkHref} className={Data.params.Vlinks.Psk[9-1].LinkClass} onClick={this.handleClick}  
    >COMMAND<span className={Data.params.Vlinks.Psk[9-1].CaretClass}
      ></span
    ></a
  ></th
></tr
      ></thead
    ><tbody
      >{this.List(Data).map(function($ps) { return<tr  key={"ps-rowby-pid-"+$ps.PID}
        ><td className="text-right"
          > {$ps.PID}</td
        ><td className="text-right"
          > {$ps.UID}</td
        ><td
          >{$ps.User}</td
        ><td className="text-right"
          > {$ps.Priority}</td
        ><td className="text-right"
          > {$ps.Nice}</td
        ><td className="text-right"
          > {$ps.Size}</td
        ><td className="text-right"
          > {$ps.Resident}</td
        ><td className="text-center"
          >{$ps.Time}</td
        ><td
          >{$ps.Name}</td
        ></tr
      >})}</tbody
    ></table
  ></div
>;
    }
  });

  jsdefines.define_loadavg = React.createClass({
    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    Reduce: function(data) {
      return {
        loadavg: data.loadavg
      };
    },
    render: function() {
      let Data = this.state; // shadow global Data
      return <span
  >{Data.loadavg}</span
>;
    }
  });

  jsdefines.define_uptime = React.createClass({
    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
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

  return jsdefines;
});
