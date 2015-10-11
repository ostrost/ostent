define(function(require) {
  var React = require('react');
  var jsdefines = {};
  jsdefines.HandlerMixin = {
    handleClick: function(e) {
      var href = e.target.getAttribute('href');
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
  // all the define_* templates transformed into jsdefines.define_* = ...;

  jsdefines.define_panelcpu = React.createClass({
    mixins: [React.addons.PureRenderMixin, jsdefines.HandlerMixin],
    List: function(data) { // static
      var list;
      if (data != null && data["CPU"] != null && (list = data["CPU"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function(data) { // static
      return {
        Params: data.Params,
        CPU: data.CPU
      };
    },
    getInitialState: function() {
      return this.Reduce(Data); // global Data
    },
    render: function() {
      var Data = this.state;
      return <div  className={!Data.Params.CPUn.Negative ? "" : "panel panel-default"}
  ><div className="h4 padding-left-like-panel-heading"
    ><a  href={Data.Params.Tlinks.CPUn} onClick={this.handleClick}
      >CPU</a
    ></div
  ><ul   className={!Data.Params.CPUn.Negative ? "hidden" : "list-group"}
    ><li className="list-group-item text-nowrap th"
      ><ul className="list-inline"
        ><li
          ><span
            ><b
              >Delay</b
            > <span className="badge"
              >{Data.Params.CPUd}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Dlinks.CPUd.Less.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.CPUd.Less.ExtraClass != null ? Data.Params.Dlinks.CPUd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Dlinks.CPUd.Less.Text}</a
><a href={Data.Params.Dlinks.CPUd.More.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.CPUd.More.ExtraClass != null ? Data.Params.Dlinks.CPUd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Dlinks.CPUd.More.Text} +</a
></div
          ></li
        ><li
          ><span
            ><b
              >Rows</b
            > <span className="badge"
              >{Data.Params.CPUn.Absolute}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Nlinks.CPUn.Less.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.CPUn.Less.ExtraClass != null ? Data.Params.Nlinks.CPUn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Nlinks.CPUn.Less.Text}</a
><a href={Data.Params.Nlinks.CPUn.More.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.CPUn.More.ExtraClass != null ? Data.Params.Nlinks.CPUn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Nlinks.CPUn.More.Text} +</a
></div
          ></li
        ></ul
      ></li
    ></ul
  ><table  className={Data.Params.CPUn.Absolute != 0 ? "table table-hover" : "hidden"}
    ><thead
      ><tr
        ><th
          ></th
        ><th className="text-right"
          >User</th
        ><th className="text-right"
          >Sys</th
        ><th className="text-right"
          >Wait</th
        ><th className="text-right"
          >Idle</th
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
    mixins: [React.addons.PureRenderMixin, jsdefines.HandlerMixin],
    List: function(data) { // static
      var list;
      if (data != null && data["DF"] != null && (list = data["DF"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function(data) { // static
      return {
        Params: data.Params,
        DF: data.DF
      };
    },
    getInitialState: function() {
      return this.Reduce(Data); // global Data
    },
    render: function() {
      var Data = this.state;
      return <div  className={!Data.Params.Dfn.Negative ? "" : "panel panel-default"}
  ><div className="h4 padding-left-like-panel-heading"
    ><a  href={Data.Params.Tlinks.Dfn} onClick={this.handleClick}
      >Disk usage</a
    ></div
  ><ul   className={!Data.Params.Dfn.Negative ? "hidden" : "list-group"}
    ><li className="list-group-item text-nowrap th"
      ><ul className="list-inline"
        ><li
          ><span
            ><b
              >Delay</b
            > <span className="badge"
              >{Data.Params.Dfd}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Dlinks.Dfd.Less.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Dfd.Less.ExtraClass != null ? Data.Params.Dlinks.Dfd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Dlinks.Dfd.Less.Text}</a
><a href={Data.Params.Dlinks.Dfd.More.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Dfd.More.ExtraClass != null ? Data.Params.Dlinks.Dfd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Dlinks.Dfd.More.Text} +</a
></div
          ></li
        ><li
          ><span
            ><b
              >Rows</b
            > <span className="badge"
              >{Data.Params.Dfn.Absolute}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Nlinks.Dfn.Less.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Dfn.Less.ExtraClass != null ? Data.Params.Nlinks.Dfn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Nlinks.Dfn.Less.Text}</a
><a href={Data.Params.Nlinks.Dfn.More.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Dfn.More.ExtraClass != null ? Data.Params.Nlinks.Dfn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Nlinks.Dfn.More.Text} +</a
></div
          ></li
        ></ul
      ></li
    ></ul
  ><table  className={Data.Params.Dfn.Absolute != 0 ? "table table-hover" : "hidden"}
    ><thead
      ><tr className="text-nowrap"
        ><th className="header "
  ><a href={Data.Params.Vlinks.Dfk[1-1].LinkHref} className={Data.Params.Vlinks.Dfk[1-1].LinkClass} onClick={this.handleClick}  
    >Device<span className={Data.Params.Vlinks.Dfk[1-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Params.Vlinks.Dfk[2-1].LinkHref} className={Data.Params.Vlinks.Dfk[2-1].LinkClass} onClick={this.handleClick}  
    >Mounted<span className={Data.Params.Vlinks.Dfk[2-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[3-1].LinkHref} className={Data.Params.Vlinks.Dfk[3-1].LinkClass} onClick={this.handleClick}  
    >Avail<span className={Data.Params.Vlinks.Dfk[3-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[4-1].LinkHref} className={Data.Params.Vlinks.Dfk[4-1].LinkClass} onClick={this.handleClick}  
    >Use%<span className={Data.Params.Vlinks.Dfk[4-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[5-1].LinkHref} className={Data.Params.Vlinks.Dfk[5-1].LinkClass} onClick={this.handleClick}  
    >Used<span className={Data.Params.Vlinks.Dfk[5-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Dfk[6-1].LinkHref} className={Data.Params.Vlinks.Dfk[6-1].LinkClass} onClick={this.handleClick}  
    >Total<span className={Data.Params.Vlinks.Dfk[6-1].CaretClass}
      ></span
    ></a
  ></th
></tr
      ></thead
    ><tbody
      >{this.List(Data).map(function($df) { return<tr  key={"df-rowby-dirname-"+$df.DirName}
        >  <td className="text-nowrap clip12" title={$df.DevName}
          >{$df.DevName}</td
        >  <td className="text-nowrap clip12" title={$df.DirName}
          >{$df.DirName}</td
        ><td className="text-right text-nowrap"
          ><span className="mutext" title="Inodes free"
            >{$df.Ifree}</span
          > {$df.Avail}</td
        ><td className="text-right bg-usepct text-nowrap"data-usepct={$df.UsePct}
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
    mixins: [React.addons.PureRenderMixin, jsdefines.HandlerMixin],
    List: function(data) { // static
      var list;
      if (data != null && data["IF"] != null && (list = data["IF"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function(data) { // static
      return {
        Params: data.Params,
        IF: data.IF
      };
    },
    getInitialState: function() {
      return this.Reduce(Data); // global Data
    },
    render: function() {
      var Data = this.state;
      return <div  className={!Data.Params.Ifn.Negative ? "" : "panel panel-default"}
  ><div className="h4 padding-left-like-panel-heading"
    ><a  href={Data.Params.Tlinks.Ifn} onClick={this.handleClick}
      >Interfaces</a
    ></div
  ><ul   className={!Data.Params.Ifn.Negative ? "hidden" : "list-group"}
    ><li className="list-group-item text-nowrap th"
      ><ul className="list-inline"
        ><li
          ><span
            ><b
              >Delay</b
            > <span className="badge"
              >{Data.Params.Ifd}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Dlinks.Ifd.Less.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Ifd.Less.ExtraClass != null ? Data.Params.Dlinks.Ifd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Dlinks.Ifd.Less.Text}</a
><a href={Data.Params.Dlinks.Ifd.More.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Ifd.More.ExtraClass != null ? Data.Params.Dlinks.Ifd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Dlinks.Ifd.More.Text} +</a
></div
          ></li
        ><li
          ><span
            ><b
              >Rows</b
            > <span className="badge"
              >{Data.Params.Ifn.Absolute}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Nlinks.Ifn.Less.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Ifn.Less.ExtraClass != null ? Data.Params.Nlinks.Ifn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Nlinks.Ifn.Less.Text}</a
><a href={Data.Params.Nlinks.Ifn.More.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Ifn.More.ExtraClass != null ? Data.Params.Nlinks.Ifn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Nlinks.Ifn.More.Text} +</a
></div
          ></li
        ></ul
      ></li
    ></ul
  ><table  className={Data.Params.Ifn.Absolute != 0 ? "table table-hover" : "hidden"}
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
        ><td className="text-nowrap clip12" title={$if.Name}
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
            ><span  className={$if.DropsOut != null ? "" : "hidden"}
              >/</span
            ><span  className={$if.DropsOut != null ? "" : "hidden"} title="Total drops Out modulo 4G"
              >{$if.DropsOut}</span
            >,<span title="Total errors In modulo 4G"
              >{$if.ErrorsIn}</span
            >/<span title="Total errors Out modulo 4G"
              >{$if.ErrorsOut}</span
            ></span
          > <span  className={(($if.DeltaDropsIn == null || $if.DeltaDropsIn == "0") && ($if.DeltaDropsOut == null || $if.DeltaDropsOut == "0") && ($if.DeltaErrorsIn == null || $if.DeltaErrorsIn == "0") && ($if.DeltaErrorsOut == null || $if.DeltaErrorsOut == "0")) ? "mutext" : ""}
            ><span title="Drops In per second"
              >{$if.DeltaDropsIn}</span
            ><span  className={$if.DeltaDropsOut != null ? "" : "hidden"}
              >/</span
            ><span  className={$if.DeltaDropsOut != null ? "" : "hidden"} title="Drops Out per second"
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
    mixins: [React.addons.PureRenderMixin, jsdefines.HandlerMixin],
    List: function(data) { // static
      var list;
      if (data != null && data["MEM"] != null && (list = data["MEM"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function(data) { // static
      return {
        Params: data.Params,
        MEM: data.MEM
      };
    },
    getInitialState: function() {
      return this.Reduce(Data); // global Data
    },
    render: function() {
      var Data = this.state;
      return <div  className={!Data.Params.Memn.Negative ? "" : "panel panel-default"}
  ><div className="h4 padding-left-like-panel-heading"
    ><a  href={Data.Params.Tlinks.Memn} onClick={this.handleClick}
      >Memory</a
    ></div
  ><ul   className={!Data.Params.Memn.Negative ? "hidden" : "list-group"}
    ><li className="list-group-item text-nowrap th"
      ><ul className="list-inline"
        ><li
          ><span
            ><b
              >Delay</b
            > <span className="badge"
              >{Data.Params.Memd}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Dlinks.Memd.Less.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Memd.Less.ExtraClass != null ? Data.Params.Dlinks.Memd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Dlinks.Memd.Less.Text}</a
><a href={Data.Params.Dlinks.Memd.More.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Memd.More.ExtraClass != null ? Data.Params.Dlinks.Memd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Dlinks.Memd.More.Text} +</a
></div
          ></li
        ><li
          ><span
            ><b
              >Rows</b
            > <span className="badge"
              >{Data.Params.Memn.Absolute}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Nlinks.Memn.Less.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Memn.Less.ExtraClass != null ? Data.Params.Nlinks.Memn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Nlinks.Memn.Less.Text}</a
><a href={Data.Params.Nlinks.Memn.More.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Memn.More.ExtraClass != null ? Data.Params.Nlinks.Memn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Nlinks.Memn.More.Text} +</a
></div
          ></li
        ></ul
      ></li
    ></ul
  ><table  className={Data.Params.Memn.Absolute != 0 ? "table table-hover" : "hidden"}
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
        ><td className="text-right bg-usepct"data-usepct={$mem.UsePct}
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
    mixins: [React.addons.PureRenderMixin, jsdefines.HandlerMixin],
    List: function(data) { // static
      var list;
      if (data != null && data["PS"] != null && (list = data["PS"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function(data) { // static
      return {
        Params: data.Params,
        PS: data.PS
      };
    },
    getInitialState: function() {
      return this.Reduce(Data); // global Data
    },
    render: function() {
      var Data = this.state;
      return <div  className={!Data.Params.Psn.Negative ? "" : "panel panel-default"}
  ><div className="h4 padding-left-like-panel-heading"
    ><a  href={Data.Params.Tlinks.Psn} onClick={this.handleClick}
      >Processes</a
    ></div
  ><ul   className={!Data.Params.Psn.Negative ? "hidden" : "list-group"}
    ><li className="list-group-item text-nowrap th"
      ><ul className="list-inline"
        ><li
          ><span
            ><b
              >Delay</b
            > <span className="badge"
              >{Data.Params.Psd}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Dlinks.Psd.Less.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Psd.Less.ExtraClass != null ? Data.Params.Dlinks.Psd.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Dlinks.Psd.Less.Text}</a
><a href={Data.Params.Dlinks.Psd.More.Href} className={"btn btn-default" + " " + (Data.Params.Dlinks.Psd.More.ExtraClass != null ? Data.Params.Dlinks.Psd.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Dlinks.Psd.More.Text} +</a
></div
          ></li
        ><li
          ><span
            ><b
              >Rows</b
            > <span className="badge"
              >{Data.Params.Psn.Absolute}</span
            ></span
          > <div className="btn-group"
            ><a href={Data.Params.Nlinks.Psn.Less.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Psn.Less.ExtraClass != null ? Data.Params.Nlinks.Psn.Less.ExtraClass : "")} onClick={this.handleClick}  
  >- {Data.Params.Nlinks.Psn.Less.Text}</a
><a href={Data.Params.Nlinks.Psn.More.Href} className={"btn btn-default" + " " + (Data.Params.Nlinks.Psn.More.ExtraClass != null ? Data.Params.Nlinks.Psn.More.ExtraClass : "")} onClick={this.handleClick}  
  >{Data.Params.Nlinks.Psn.More.Text} +</a
></div
          ></li
        ></ul
      ></li
    ></ul
  ><table  className={Data.Params.Psn.Absolute != 0 ? "table table-hover" : "hidden"}
    ><thead
      ><tr className="text-nowrap"
        ><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[1-1].LinkHref} className={Data.Params.Vlinks.Psk[1-1].LinkClass} onClick={this.handleClick}  
    >PID<span className={Data.Params.Vlinks.Psk[1-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[2-1].LinkHref} className={Data.Params.Vlinks.Psk[2-1].LinkClass} onClick={this.handleClick}  
    >UID<span className={Data.Params.Vlinks.Psk[2-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Params.Vlinks.Psk[3-1].LinkHref} className={Data.Params.Vlinks.Psk[3-1].LinkClass} onClick={this.handleClick}  
    >USER<span className={Data.Params.Vlinks.Psk[3-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[4-1].LinkHref} className={Data.Params.Vlinks.Psk[4-1].LinkClass} onClick={this.handleClick}  
    >PR<span className={Data.Params.Vlinks.Psk[4-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[5-1].LinkHref} className={Data.Params.Vlinks.Psk[5-1].LinkClass} onClick={this.handleClick}  
    >NI<span className={Data.Params.Vlinks.Psk[5-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[6-1].LinkHref} className={Data.Params.Vlinks.Psk[6-1].LinkClass} onClick={this.handleClick}  
    >VIRT<span className={Data.Params.Vlinks.Psk[6-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-right"
  ><a href={Data.Params.Vlinks.Psk[7-1].LinkHref} className={Data.Params.Vlinks.Psk[7-1].LinkClass} onClick={this.handleClick}  
    >RES<span className={Data.Params.Vlinks.Psk[7-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header text-center"
  ><a href={Data.Params.Vlinks.Psk[8-1].LinkHref} className={Data.Params.Vlinks.Psk[8-1].LinkClass} onClick={this.handleClick}  
    >TIME<span className={Data.Params.Vlinks.Psk[8-1].CaretClass}
      ></span
    ></a
  ></th
><th className="header "
  ><a href={Data.Params.Vlinks.Psk[9-1].LinkHref} className={Data.Params.Vlinks.Psk[9-1].LinkClass} onClick={this.handleClick}  
    >COMMAND<span className={Data.Params.Vlinks.Psk[9-1].CaretClass}
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
  return jsdefines;
});
