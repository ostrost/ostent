/* eslint indent:0 no-console:0 */
/* eslint-env node */
/* global Data:false, options:false, pug:false */
/*
var React       = require('react'),
    ReactDOM    = require('react-dom'),
    ReconnectWS = require('reconnectingwebsocket'),
    jsxdefines  = require('./jsxdefines.js');
*/

var format = require('util').format;

var React      = require('react'),
    ReactDOM   = require('react-dom'),
    SparkLines = require('react-sparklines');

if (typeof(options) === 'undefined' || !options.ssr) {
  var ReconnectWS = require('reconnectingwebsocket');
}

class Sparkline extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {data: [], limit: 90, width: 180};
  }
  componentDidUpdate(_, prevState) {
    var root = ReactDOM.findDOMNode(this.refs.root);
    if (root == null) {
      return;
    }
    var rootWidth = Math.floor(root.offsetWidth) - 10;
    if (prevState.width != rootWidth) {
      this.setState({width: rootWidth, limit: Math.round(rootWidth/2)});
    }
  }
  NewState(statentry) {
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
  }
  render() {
    const e = React.createElement;
    let rootdiv = (c) => e('div', {ref: 'root'}, c);
    if (process.env.NODE_ENV !== 'production' && this.props.ssr) {
      return rootdiv();
    }
    var curveProps = {style: {strokeWidth: 1}};
    var spotsProps = {size: 2, spotColors: {'-1': 'green', '1': 'red'}}; // reverse default colors
    if (this.props.defaultSpots) { delete spotsProps.spotColors; } // back to default colors
    curveProps.key = 'curve';
    spotsProps.key = 'spots';
    return rootdiv(
      e(SparkLines.Sparklines, {
        data: this.state.data,
        limit: this.state.limit,
        width: this.state.width,
        svgWidth: this.state.width,
        height: 33,
        svgHeight: 33}, [
          e(SparkLines.SparklinesCurve, curveProps, null),
          e(SparkLines.SparklinesSpots, spotsProps, null),
        ]));
  }
}

class Page extends React.PureComponent {
  maj_keys(prefix, keys) {
    return keys.reduce((o, k) => Object.assign(o, ({[k]: '[[ '+prefix+k+' ]]'})), {});
  }
  constructor(props) {
    super(props);
    if (props.initialState) {
      this.state = props.initialState;
    }
    if (!(process.env.NODE_ENV !== 'production' && props.ssr)) {
      return;
    }
    let maj_keys = this.maj_keys;
    let data = maj_keys('.Data.', ['Distrib', 'OstentUpgrade', 'OstentVersion']);
    data.Exporting = [maj_keys('$e.', ['Header', 'Text'])];
    data.la = {List: [maj_keys('$la.', [
      'Period',
      'Value',
    ])]};
    data.mem = {List: [maj_keys('$mem.', [
      'Kind',
      'Total',
      'Used',
      'Free',
      'UsePct',
    ])]};
    data.df = {List: [maj_keys('$df.', [
      'DevName',
      'DirName',
      'Total',
      'Inodes',
      'Used',
      'Iused',
      'Avail',
      'Ifree',
      'UsePct',
      'IusePct',
    ])]};
    data.cpu = {List: [maj_keys('$cpu.', [
      'N',
      'UserPct',
      'SysPct',
      'WaitPct',
      'IdlePct',
    ])]};
    data.netio = {List: [maj_keys('$if.', [
      'DeltaDropsIn',
      'DeltaDropsOut',
      'DeltaErrorsIn',
      'DeltaErrorsOut',
      'Name',
      'IP',
      'DropsIn',
      'DropsOut',
      'ErrorsIn',
      'ErrorsOut',
      'DeltaDropsIn',
      'DeltaDropsOut',
      'DeltaErrorsIn',
      'DeltaErrorsOut',
      'PacketsIn',
      'PacketsOut',
      'DeltaPacketsIn',
      'DeltaPacketsOut',
      'BytesIn',
      'BytesOut',
      'DeltaBitsIn',
      'DeltaBitsOut',
    ])]};
    data.procs = {List: [maj_keys('$ps.', [
      'PID',
      'UID',
      'User',
      'Priority',
      'Nice',
      'Size',
      'Resident',
      'Time',
      'Name',
    ])]};
    data.system_ostent = maj_keys('.Data.system_ostent.', [
      'hostname_short',
      'uptime_format',
    ]);
    data.params = {
      Tlinks: maj_keys('.Data.params.Tlinks.', [
        'Lan',
        'Memn',
        'Dfn',
        'CPUn',
        'Ifn',
        'Psn',
      ]),
      xVlinks: {
        Dfk: maj_keys('.Data.params.Vlinks.Dfk.', [
          'CaretClass',
          'LinkClass',
          'LinkHref',
        ]),
        Psk: maj_keys('.Data.params.Vlinks.Psk.', [
          'CaretClass',
          'LinkClass',
          'LinkHref',
        ]),
      },
      Vlinks: maj_keys('.Data.params.Vlinks.', [
        'Dfk',
        'Psk',
      ]),
      Nlinks: {
        Lan: {
          Less: maj_keys('.Data.params.Nlinks.Lan.Less.', ['ExtraClass', 'Href', 'Text']),
          More: maj_keys('.Data.params.Nlinks.Lan.More.', ['ExtraClass', 'Href', 'Text']),
        },
        Memn: {
          Less: maj_keys('.Data.params.Nlinks.Memn.Less.', ['ExtraClass', 'Href', 'Text']),
          More: maj_keys('.Data.params.Nlinks.Memn.More.', ['ExtraClass', 'Href', 'Text']),
        },
        Dfn: {
          Less: maj_keys('.Data.params.Nlinks.Dfn.Less.', ['ExtraClass', 'Href', 'Text']),
          More: maj_keys('.Data.params.Nlinks.Dfn.More.', ['ExtraClass', 'Href', 'Text']),
        },
        CPUn: {
          Less: maj_keys('.Data.params.Nlinks.CPUn.Less.', ['ExtraClass', 'Href', 'Text']),
          More: maj_keys('.Data.params.Nlinks.CPUn.More.', ['ExtraClass', 'Href', 'Text']),
        },
        Ifn: {
          Less: maj_keys('.Data.params.Nlinks.Ifn.Less.', ['ExtraClass', 'Href', 'Text']),
          More: maj_keys('.Data.params.Nlinks.Ifn.More.', ['ExtraClass', 'Href', 'Text']),
        },
        Psn: {
          Less: maj_keys('.Data.params.Nlinks.Psn.Less.', ['ExtraClass', 'Href', 'Text']),
          More: maj_keys('.Data.params.Nlinks.Psn.More.', ['ExtraClass', 'Href', 'Text']),
        },
      },
      Lan: maj_keys('.Data.params.Lan.', [
        'Absolute',
        'Negative',
      ]),
      Memn: maj_keys('.Data.params.Memn.', [
        'Absolute',
        'Negative',
      ]),
      Dfn: maj_keys('.Data.params.Dfn.', [
        'Absolute',
        'Negative',
      ]),
      CPUn: maj_keys('.Data.params.CPUn.', [
        'Absolute',
        'Negative',
      ]),
      Ifn: maj_keys('.Data.params.Ifn.', [
        'Absolute',
        'Negative',
      ]),
      Psn: maj_keys('.Data.params.Psn.', [
        'Absolute',
        'Negative',
      ]),
    };
    this.state = data;
  }
  NewState(data) {
    this.setState(data);
    Object.keys(this.refs).forEach((ref) => {
      if (this.refs[ref].props.getter) {
        let statentry = data['cpu'].List.filter((x) => x.N == this.refs[ref].props.getter)[0];
        this.refs[ref].NewState(statentry);
        return;
      }
      let refindex = ref.lastIndexOf('.');
      if (refindex == -1) {
        return;
      }
      let subindex = +ref.substr(refindex+1),
          datakey  =  ref.substr(0, refindex);
      let dkindex = datakey.lastIndexOf('.');
      if (dkindex != -1) {
        datakey = datakey.substr(dkindex+1);
      }

      let subdata = data[datakey];
      if (!subdata || !subdata.List) {
        return;
      }
      this.refs[ref].NewState(subdata.List[+subindex]);
    });
  }

  handleClick(e) {
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

  range(v, items, f, ...args) {
    let lines = [f(...args)];
    if (process.env.NODE_ENV !== 'production' && this.props.ssr) {
      lines.unshift(`[[range ${v} := ${items}]]`);
      lines.push('[[end]]');
    }
    return lines;
  }
  render() {
    let ostentUpgrade = (latestOstent) => {
      let lines = [pug`.top-bar-left: a(href=latestOstent) ${ this.state.OstentUpgrade }}`];
      if (process.env.NODE_ENV !== 'production' && this.props.ssr) {
        lines.unshift('[[if .Data.OstentUpgrade]]');
        lines.push('[[end]]');
      } else if (!this.state.OstentUpgrade) {
        lines = [];
      }
      return lines;
    };
    let exporting_line = (e, i) => pug`.row(key='exporting_line_' + i)
  span.expand.small-12: small: pre
  b ${ e.Header }
    | ${ e.Text }`;
    let exporting = () => {
      let elen = this.state.Exporting.length;
      let noutputs = `${ elen } output ${ elen !== 1 ? 's' : '' }`;
      if (process.env.NODE_ENV !== 'production' && this.props.ssr) {
        noutputs = '[[len .Data.Exporting]] output[[if ne 1 (len .Data.Exporting)]]s[[end]]';
      }
      let lines = [pug`.row.expanded.hr-bottom: .column.large-11.small-offset-1
  h5 Exporting
  .stripe
    .row.thead.nobold
      span.expand.small-12: small: pre
        | [outputs]
        |
        |
        i
          |     # [outputs.ostent] not counted nor shown. Counted 
          b ${ noutputs }
          | :
    ${ this.state.Exporting && this.state.Exporting.map((e, i) =>
      this.range('$e', '.Exporting', exporting_line, e, i)) }`];

      if (process.env.NODE_ENV !== 'production' && this.props.ssr) {
        lines.unshift('[[if ne (len .Data.Exporting) 0]]');
        lines.push('[[end]]');
      } else if (this.state.Exporting.length === 0) {
        lines = [];
      }
      return lines;
    };

    let htitle = 'hostname ' + this.state.system_ostent.hostname_short;
    let hostname = pug`a(href='/', title=htitle)
  = this.state.system_ostent.hostname_short`;
    let uptime = pug`span
  = this.state.system_ostent.uptime_format`;

    let vlink = (kparam, number, text) => {
      // console.log('vlink', kparam, number, text);
      let v = this.state.params.Vlinks[kparam];
      // console.log('v', v);
      if (process.env.NODE_ENV !== 'production' && this.props.ssr) {
        let y = v;
        v = format('(index %s %d)', v.replace(/(^\[\[\s?|\s?\]\]$)/g, ''), number-1);
        v = {
          CaretClass: v+'.CaretClass',
          LinkClass:  v+'.LinkClass',
          LinkHref:   v+'.LinkHref',
        };
        v = this.maj_keys(format('(index %s %d).', y.replace(/(^\[\[\s?|\s?\]\]$)/g, ''), number-1), [
          'CaretClass',
          'LinkClass',
          'LinkHref',
        ]);
      } else {
        v = v[number-1];
      }
      // console.log('v', v);
      let handleClick = (e) => this.handleClick(e);
      return pug`
a(className=v.LinkClass, href=v.LinkHref, onClick=handleClick)
  ${ text }
  span(className=v.CaretClass)`;
    };

    let table = (nparam, title, block) => {
      let handleClick = (e) => this.handleClick(e);
      let alink = (cls, fmt, link) => {
        cls = cls + ' ' + (link.ExtraClass || '');
        return pug`button(className=cls, href=link.Href, onClick=handleClick) ${ format(fmt, link.Text) }`;
      };
      let negative = this.state.params[nparam].Negative ? 'show-showhide' : 'hide-showhide';
      let absolute = this.state.params[nparam].Absolute == 0 ? 'hide' : '';
      if (process.env.NODE_ENV !== 'production' && this.props.ssr) {
        negative = `[[if .Data.params.${ nparam }.Negative]]show-showhide[[else]]hide-showhide[[end]]`;
        absolute = `[[if eq .Data.params.${ nparam }.Absolute 0]]hide[[end]]`;
      }
      let atitle = title + ' display options';
      return pug`.row.expanded.hr-bottom
  .column.text-right.small-1
    div(className=${ negative })
      h1.h5.text-overflow
        a(title=atitle, href=this.state.params.Tlinks[nparam], onClick=handleClick)
          span.showhide-hide.whitespace-pre.float-left ... 
          | ${ title }
  .column.large-11
    div(className=${ negative })
      ul.row.menu.showhide-show
        li: .input-group
          .input-group-label.text-nowrap rows
          .input-group-button: button.button.small.secondary.disabled: = this.state.params[nparam].Absolute
          - var cls = 'text-nowrap button small';
          .input-group-button: = alink(cls, '- %s', this.state.params.Nlinks[nparam].Less)
          .input-group-button: = alink(cls, '%s +', this.state.params.Nlinks[nparam].More)
      div(className=${ absolute })
        ${ block }`;
    };

    let la_line = (la, i) => pug`.row(key='la_line_' + i)
  span.expand.col.small-1: .text-right.width-3rem ${ la.Period }m
  span.expand.col.small-1.text-right  ${ la.Value }
  span.expand.col-lr: Sparkline(ref='la.'+i, col='Value', ssr=this.props.ssr)`;
    let la = table('Lan', 'Load avg', pug`.stripe
  .row.thead
    span.expand.col.small-1 Period
    span.expand.col.small-1.text-right  Value
    span.expand.col
  ${ this.state.la.List && this.state.la.List.map((la, i) =>
       this.range('$la', '.Data.la.List', la_line, la, i)) }`);

    let mem_line = (mem, i) => pug`.row(key='mem_line_' + i)
  span.expand.col.small-1 ${ mem.Kind }
  span.expand.col.small-1.text-right  ${ mem.Total }
  span.expand.col.small-1.text-right  ${ mem.Used }
  span.expand.col.small-1.text-right  ${ mem.Free }
  span.expand.col.small-1.text-right.bg-usepct(data-usepct=mem.UsePct) ${ mem.UsePct }%
  span.expand.col-lr: Sparkline(ref='mem.'+i, col='UsePct', ssr=this.props.ssr)`;
    let mem = table('Memn', 'Memory', pug`.stripe
  .row.thead
    span.expand.col.small-1 Memory
    span.expand.col.small-1.text-right  Total
    span.expand.col.small-1.text-right  Used
    span.expand.col.small-1.text-right  Free
    span.expand.col.small-1.text-right  Use%
    span.expand.col
  ${ this.state.mem.List && this.state.mem.List.map((mem, i) =>
       this.range('$mem', '.Data.mem.List', mem_line, mem, i)) }`);

    let df_line = (df, i) => pug`.row(key='df_line_'+i)
  span.expand.col.small-1.text-overflow ${ df.DevName }
  span.expand.col.small-1.text-overflow  ${ df.DirName }
  span.expand.col.small-1.text-overflow.text-right.gray
    span.float-right  ${ df.Total }
    span(title='Inodes total')  ${ df.Inodes }
  span.expand.col.small-1.text-overflow.text-right.gray
    span.float-right  ${ df.Used }
    span(title='Inodes used')  ${ df.Iused }
  span.expand.col.small-1.text-overflow.text-right.gray
    span.float-right  ${ df.Avail }
    span(title='Inodes free')  ${ df.Ifree }
  span.expand.col.small-1.text-overflow.text-right.gray.bg-usepct(data-usepct=df.UsePct)
    span.float-right  ${ df.UsePct }%
    span(title='Inodes use%')  ${ df.IusePct }%
  span.expand.col-lr
    Sparkline(ref='df.'+i, col='UsePct', ssr=this.props.ssr)`;
    let df = table('Dfn', 'Disk usage', pug`.stripe
  .row.thead
    //- FS:
    span.expand.col.small-1.text-nowrap:            = vlink('Dfk', 1, 'Device')
    //- MP:
    span.expand.col.small-1.text-nowrap:            = vlink('Dfk', 2, 'Mounted')
    //- TOTAL:
    span.expand.col.small-1.text-nowrap.text-right: = vlink('Dfk', 6, 'Total')
    //- USED:
    span.expand.col.small-1.text-nowrap.text-right: = vlink('Dfk', 5, 'Used')
    //- AVAIL:
    span.expand.col.small-1.text-nowrap.text-right: = vlink('Dfk', 3, 'Avail')
    //- USEPCT:
    span.expand.col.small-1.text-nowrap.text-right: = vlink('Dfk', 4, 'Use%')
    span.expand.col
  ${ this.state.df.List && this.state.df.List.map((df, i) =>
       this.range('$df', '.Data.df.List', df_line, df, i)) }`);

    //- let cpu_spark = (i, col) => (data) => +data.cpu.List[i][col];
    let cpu_spark = (N, col) => (data) => data.cpu.List.filter((x) => x.N == N)[0][col];
    let cpu_line = (cpu, i) => pug`.row(key='cpu_line_'+cpu.N)
  span.expand.col.small-1.text-nowrap ${ cpu.N }
  span.expand.col.small-1.text-right.bg-usepct(data-usepct=cpu.UserPct) ${ cpu.UserPct }%
  span.expand.col.small-1.text-right.bg-usepct(data-usepct=cpu.SysPct) ${ cpu.SysPct }%
  span.expand.col.small-1.text-right.bg-usepct(data-usepct=cpu.WaitPct) ${ cpu.WaitPct }%
  span.expand.col.small-1.text-right.bg-usepct-inverse(data-usepct=cpu.IdlePct) ${ cpu.IdlePct }%
  span.expand.col-lr
    //- ref=cpu.N+'.cpu.'+i, col='IdlePct'
    Sparkline(getter=cpu.N, col='IdlePct', ref='cpu.'+i, defaultSpots=true, ssr=this.props.ssr)`;
    let cpu = table('CPUn', 'CPU', pug`.stripe
  .row.thead
    span.expand.col.small-1 Core
    span.expand.col.small-1.text-right  User%
    span.expand.col.small-1.text-right  Sys%
    span.expand.col.small-1.text-right  Wait%
    span.expand.col.small-1.text-right  Idle%
    span.expand.col
  ${ this.state.cpu.List && this.state.cpu.List.map((cpu, i) =>
       this.range('$cpu', '.Data.cpu.List', cpu_line, cpu, i)) }`);

    let if_line = (if_, i) => {
      let deltacls = (if_.DeltaDropsIn == "0" &&
                      if_.DeltaDropsOut == "0" &&
                      if_.DeltaErrorsIn == "0" &&
                      if_.DeltaErrorsOut == "0") ? 'gray' : '';
      if (process.env.NODE_ENV !== 'production' && this.props.ssr) {
        deltacls = '[[if and (eq $if.DeltaDropsIn `0`) (eq $if.DeltaDropsOut `0`) (eq $if.DeltaErrorsIn `0`) (eq $if.DeltaErrorsOut `0`)]]gray[[end]]';
      }
      return pug`.row(key='if_line_'+i)
  span.expand.col.small-1.text-overflow ${ if_.Name }
  span.expand.col.small-1.text-overflow.text-right ${ if_.IP }
  span.expand.col.small-2.text-right.text-nowrap
    |  
    span.gray(title='Total drops,errors modulo 4G')
      span(title='Total drops In modulo 4G') ${ if_.DropsIn }
      | /
      span(title='Total drops Out modulo 4G') ${ if_.DropsOut }
      | ,
      span(title='Total errors In modulo 4G') ${ if_.ErrorsIn }
      | /
      span(title='Total errors Out modulo 4G') ${ if_.ErrorsOut }
    |  
    span(className=deltacls)
      span(title='Drops In per second') ${ if_.DeltaDropsIn }
      | /
      span(title='Drops Out per second') ${ if_.DeltaDropsOut }
      | ,
      span(title='Errors In per second') ${ if_.DeltaErrorsIn }
      | /
      span(title='Errors Out per second') ${ if_.DeltaErrorsOut }
  span.expand.col.small-2.text-right.text-nowrap
    |  
    span.gray
      span(title='Total packets In modulo 4G') ${ if_.PacketsIn }
      | /
      span(title='Total packets Out modulo 4G') ${ if_.PacketsOut }
    |  
    span(title='Packets In per second') ${ if_.DeltaPacketsIn }
    | /
    span(title='Packets Out per second') ${ if_.DeltaPacketsOut }
  span.expand.col.small-2.text-right.text-nowrap
    |  
    span.gray
      span(title='Total BYTES In modulo 4G') ${ if_.BytesIn }
      | /
      span(title='Total BYTES Out modulo 4G') ${ if_.BytesOut }
    |  
    span(title='BITS In per second') ${ if_.DeltaBitsIn }
    | /
    span(title='BITS Out per second') ${ if_.DeltaBitsOut }
  span.expand.col-lr
    Sparkline(ref='netio.'+i, col='DeltaBytesOutNum', ssr=this.props.ssr)`;
    };
    let ifs = table('Ifn', 'Interfaces', pug`.stripe
  .row.thead
    span.expand.col.small-1 Interface
    span.expand.col.small-1.text-right  IP
    span.expand.col.small-2.text-right.text-nowrap(title='Drops,Errors In/Out per second')  Loss IO ps
    span.expand.col.small-2.text-right.text-nowrap(title='Packets In/Out per second')  Packets IO ps
    span.expand.col.small-2.text-right.text-nowrap(title='Bits In/Out per second')
      |  IO 
      i b
      | ps
    span.expand.col
  ${ this.state.netio.List && this.state.netio.List.map((if_, i) =>
       this.range('$if', '.Data.netio.List', if_line, if_, i)) }`);

    let ps_line = (ps, i) => pug`.row(key='ps_line_'+i)
  span.expand.col.small-1.text-right ${ ps.PID }
  span.expand.col.small-1.text-right  ${ ps.UID }
  span.expand.col.small-1  ${ ps.User }
  span.expand.col.small-1.text-right  ${ ps.Priority }
  span.expand.col.small-1.text-right  ${ ps.Nice }
  span.expand.col.small-1.text-right  ${ ps.Size }
  span.expand.col.small-1.text-right  ${ ps.Resident }
  span.expand.col.small-1.text-center  ${ ps.Time }
  span.expand.col  ${ ps.Name }`;
    let ps = table('Psn', 'Processes', pug`.stripe
  .row.thead
    span.expand.col.small-1.text-nowrap.text-right:  = vlink('Psk', 1, 'PID')
    span.expand.col.small-1.text-nowrap.text-right:  = vlink('Psk', 2, 'UID')
    span.expand.col.small-1.text-nowrap:             = vlink('Psk', 3, 'USER')
    //- PRI:
    span.expand.col.small-1.text-nowrap.text-right:  = vlink('Psk', 4, 'PR')
    //- NICE:
    span.expand.col.small-1.text-nowrap.text-right:  = vlink('Psk', 5, 'NI')
    span.expand.col.small-1.text-nowrap.text-right:  = vlink('Psk', 6, 'VIRT')
    span.expand.col.small-1.text-nowrap.text-right:  = vlink('Psk', 7, 'RES')
    span.expand.col.small-1.text-nowrap.text-center: = vlink('Psk', 8, 'TIME')
    //- NAME:
    span.expand.col.small-1.text-nowrap:             = vlink('Psk', 9, 'COMMAND')
  ${ this.state.procs.List && this.state.procs.List.map((ps, i) =>
       this.range('$ps', '.Data.procs.List', ps_line, ps, i)) }`);

    let latestOstent = 'https://www.ostrost.com/ostent/releases/latest?cmp=' + this.state.OstentVersion;
    return pug`.top-bar
  .top-bar-title: h2.h5.margin-bottom-0
    = hostname
    |  
    a(href=latestOstent) ostent
  = ostentUpgrade(latestOstent)
  div: .top-bar-right: h2.h5.margin-bottom-0: small
    | ${ this.state.Distrib } up 
    span.whitespace-pre: =uptime

= la
= mem
= df
= cpu
= ifs
= ps
= exporting()
`;
  }
};

/*
function neweventsource(onmessage) {
  var init, conn = null;
  var sendSearch = function(search) {
    location.search = search;
    console.log('new search', search);
    conn.close();
    window.setTimeout(init, 100);
  };
  var sendLocationSearch = function() { sendSearch(location.search); };
  init = function() {
    conn = new EventSource('/index.sse' + location.search);
    conn.onopen = function() {
      console.log('sse opened');
      window.addEventListener('popstate', sendLocationSearch);
    };
    conn.onclose = function() {
      console.log('sse closed');
      window.removeEventListener('popstate', sendLocationSearch);
    };
    conn.onerror = function() {
      console.log('sse errord');
      window.removeEventListener('popstate', sendLocationSearch);
    };
    conn.onmessage = onmessage;
  };
  init();
  return {
    sendSearch: sendSearch,
    close: function() { return conn.close(); }
  };
}; // */

function main(data) {
  if (data.params.Still.Absolute != 0) {
    return;
  }

  let page = ReactDOM.hydrate(
    React.createElement(Page, {initialState: data}),
    document.getElementById('page'));

  var wscheme = location.protocol === 'https:' ? 'wss' : 'ws';
  var ws = new ReconnectWS(wscheme + '://' + location.host + '/index.ws');

  // sendSearch is also referenced as window.updates.sendSearch
  ws.sendSearch = function(search) {
    console.log('ws send', search);
    ws.send(JSON.stringify({Search: search}));
  };
  ws.sendLocationSearch = function() { ws.sendSearch(location.search); };

  ws.onclose = function() {
    console.log('ws closed');
    window.removeEventListener('popstate', ws.sendLocationSearch);
  };
  ws.onopen = function() {
    console.log('ws opened');
    ws.sendLocationSearch();
    window.addEventListener('popstate', ws.sendLocationSearch);
  };
  ws.onmessage = function(event) { // the onmessage
    var state = JSON.parse(event.data);
    if (state == null) {
      return;
    }
    if ((state.Reload != null) && state.Reload) {
      window.setTimeout((function() { location.reload(true); }), 5000);
      window.setTimeout(ws.close, 2000);
      console.log('in 5s: location.reload(true)');
      console.log('in 2s: ws.close()');
      return;
    }
    if (state.Error != null) {
      console.log('Error', state.Error);
      return;
    }
    page.NewState(state);
  };

  window.updates = ws; // neweventsource(onmessage);
}

if (typeof(Data) !== 'undefined') { main(Data); } // global Data

module.exports = Page;
