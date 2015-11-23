'use strict';

var _typeofReactElement = typeof Symbol === 'function' && Symbol['for'] && Symbol['for']('react.element') || 60103;

define(function (require) {
  var React = require('react');
  var $ = require('jquery');
  var jsdefines = {};
  jsdefines.StateHandlingMixin = { // requires .Reduce method
    getInitialState: function getInitialState() {
      return this.StateFrom(Data); // global Data
    },
    NewState: function NewState(data) {
      var state = this.StateFrom(data);
      if (state != null) {
        this.setState(state);
      }
    },
    StateFrom: function StateFrom(data) {
      var state = this.Reduce(data);
      if (state != null) {
        for (var key in state) {
          if (state[key] == null) {
            delete state[key];
          }
        }
      }
      return state;
    }
  };
  jsdefines.HandlerMixin = {
    handleClick: function handleClick(e) {
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

  // transformed from define_* templates:

  jsdefines.define_hostname = React.createClass({
    displayName: 'define_hostname',

    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    Reduce: function Reduce(data) {
      return {
        hostname: data.hostname
      };
    },
    render: function render() {
      var Data = this.state; // shadow global Data
      return {
        $$typeof: _typeofReactElement,
        type: 'a',
        key: null,
        ref: null,
        props: {
          children: {
            $$typeof: _typeofReactElement,
            type: 'pre',
            key: null,
            ref: null,
            props: {
              children: Data.hostname
            },
            _owner: null
          },
          className: 'h5',
          title: "hostname " + Data.hostname,
          href: '/'
        },
        _owner: null
      };
    }
  });

  jsdefines.define_panelcpu = React.createClass({
    displayName: 'define_panelcpu',

    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    List: function List(data) {
      var list = undefined;
      if (data != null && data["cpu"] != null && (list = data["cpu"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function Reduce(data) {
      return {
        params: data.params,
        cpu: data.cpu
      };
    },
    render: function render() {
      var Data = this.state; // shadow global Data
      return {
        $$typeof: _typeofReactElement,
        type: 'div',
        key: null,
        ref: null,
        props: {
          children: [{
            $$typeof: _typeofReactElement,
            type: 'div',
            key: null,
            ref: null,
            props: {
              children: [{
                $$typeof: _typeofReactElement,
                type: 'div',
                key: null,
                ref: null,
                props: {
                  children: {
                    $$typeof: _typeofReactElement,
                    type: 'a',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'h5',
                        key: null,
                        ref: null,
                        props: {
                          children: 'CPU',
                          className: 'margin-bottom-0'
                        },
                        _owner: null
                      },
                      href: Data.params.Tlinks.CPUn,
                      onClick: this.handleClick
                    },
                    _owner: null
                  },
                  className: 'tabs-title menu-tab-padding'
                },
                _owner: null
              }, {
                $$typeof: _typeofReactElement,
                type: 'ul',
                key: null,
                ref: null,
                props: {
                  children: [{
                    $$typeof: _typeofReactElement,
                    type: 'li',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'div',
                        key: null,
                        ref: null,
                        props: {
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: 'delay',
                              className: 'input-group-label'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: Data.params.CPUd,
                              className: 'input-group-label label secondary'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['- ', Data.params.Dlinks.CPUd.Less.Text],
                              href: Data.params.Dlinks.CPUd.Less.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.CPUd.Less.ExtraClass != null ? Data.params.Dlinks.CPUd.Less.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: [Data.params.Dlinks.CPUd.More.Text, ' +'],
                              href: Data.params.Dlinks.CPUd.More.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.CPUd.More.ExtraClass != null ? Data.params.Dlinks.CPUd.More.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }],
                          className: 'input-group margin-bottom-0'
                        },
                        _owner: null
                      },
                      className: 'menu-text'
                    },
                    _owner: null
                  }, {
                    $$typeof: _typeofReactElement,
                    type: 'li',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'div',
                        key: null,
                        ref: null,
                        props: {
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: 'rows',
                              className: 'input-group-label'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: Data.params.CPUn.Absolute,
                              className: 'input-group-label label secondary'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['- ', Data.params.Nlinks.CPUn.Less.Text],
                              href: Data.params.Nlinks.CPUn.Less.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.CPUn.Less.ExtraClass != null ? Data.params.Nlinks.CPUn.Less.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: [Data.params.Nlinks.CPUn.More.Text, ' +'],
                              href: Data.params.Nlinks.CPUn.More.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.CPUn.More.ExtraClass != null ? Data.params.Nlinks.CPUn.More.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }],
                          className: 'input-group margin-bottom-0'
                        },
                        _owner: null
                      },
                      className: 'menu-text'
                    },
                    _owner: null
                  }],
                  className: 'float-left bar menu'
                },
                _owner: null
              }],
              className: !Data.params.CPUn.Negative ? "tabs tabs-border bar-less" : "tabs tabs-border",
              'data-tabs': true
            },
            _owner: null
          }, {
            $$typeof: _typeofReactElement,
            type: 'table',
            key: null,
            ref: null,
            props: {
              children: [{
                $$typeof: _typeofReactElement,
                type: 'thead',
                key: null,
                ref: null,
                props: {
                  children: {
                    $$typeof: _typeofReactElement,
                    type: 'tr',
                    key: null,
                    ref: null,
                    props: {
                      children: [{
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {},
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'User%',
                          className: 'text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Sys%',
                          className: 'text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Wait%',
                          className: 'text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Idle%',
                          className: 'text-right'
                        },
                        _owner: null
                      }]
                    },
                    _owner: null
                  }
                },
                _owner: null
              }, {
                $$typeof: _typeofReactElement,
                type: 'tbody',
                key: null,
                ref: null,
                props: {
                  children: this.List(Data).map(function ($cpu) {
                    return {
                      $$typeof: _typeofReactElement,
                      type: 'tr',
                      key: "cpu-rowby-N-" + $cpu.N,
                      ref: null,
                      props: {
                        children: [{
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $cpu.N,
                            className: 'text-right text-nowrap'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [$cpu.UserPct, '%'],
                            className: 'text-right bg-usepct',
                            'data-usepct': $cpu.UserPct
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [$cpu.SysPct, '%'],
                            className: 'text-right bg-usepct',
                            'data-usepct': $cpu.SysPct
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [$cpu.WaitPct, '%'],
                            className: 'text-right bg-usepct',
                            'data-usepct': $cpu.WaitPct
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [$cpu.IdlePct, '%'],
                            className: 'text-right bg-usepct-inverse',
                            'data-usepct': $cpu.IdlePct
                          },
                          _owner: null
                        }]
                      },
                      _owner: null
                    };
                  })
                },
                _owner: null
              }],
              className: Data.params.CPUn.Absolute != 0 ? "hover scroll-x margin-bottom-0" : "hide"
            },
            _owner: null
          }]
        },
        _owner: null
      };
    }
  });

  jsdefines.define_paneldf = React.createClass({
    displayName: 'define_paneldf',

    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    List: function List(data) {
      var list = undefined;
      if (data != null && data["diskUsage"] != null && (list = data["diskUsage"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function Reduce(data) {
      return {
        params: data.params,
        diskUsage: data.diskUsage
      };
    },
    render: function render() {
      var Data = this.state; // shadow global Data
      return {
        $$typeof: _typeofReactElement,
        type: 'div',
        key: null,
        ref: null,
        props: {
          children: [{
            $$typeof: _typeofReactElement,
            type: 'div',
            key: null,
            ref: null,
            props: {
              children: [{
                $$typeof: _typeofReactElement,
                type: 'div',
                key: null,
                ref: null,
                props: {
                  children: {
                    $$typeof: _typeofReactElement,
                    type: 'a',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'h5',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Disk usage',
                          className: 'margin-bottom-0'
                        },
                        _owner: null
                      },
                      href: Data.params.Tlinks.Dfn,
                      onClick: this.handleClick
                    },
                    _owner: null
                  },
                  className: 'tabs-title menu-tab-padding'
                },
                _owner: null
              }, {
                $$typeof: _typeofReactElement,
                type: 'ul',
                key: null,
                ref: null,
                props: {
                  children: [{
                    $$typeof: _typeofReactElement,
                    type: 'li',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'div',
                        key: null,
                        ref: null,
                        props: {
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: 'delay',
                              className: 'input-group-label'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: Data.params.Dfd,
                              className: 'input-group-label label secondary'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['- ', Data.params.Dlinks.Dfd.Less.Text],
                              href: Data.params.Dlinks.Dfd.Less.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Dfd.Less.ExtraClass != null ? Data.params.Dlinks.Dfd.Less.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: [Data.params.Dlinks.Dfd.More.Text, ' +'],
                              href: Data.params.Dlinks.Dfd.More.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Dfd.More.ExtraClass != null ? Data.params.Dlinks.Dfd.More.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }],
                          className: 'input-group margin-bottom-0'
                        },
                        _owner: null
                      },
                      className: 'menu-text'
                    },
                    _owner: null
                  }, {
                    $$typeof: _typeofReactElement,
                    type: 'li',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'div',
                        key: null,
                        ref: null,
                        props: {
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: 'rows',
                              className: 'input-group-label'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: Data.params.Dfn.Absolute,
                              className: 'input-group-label label secondary'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['- ', Data.params.Nlinks.Dfn.Less.Text],
                              href: Data.params.Nlinks.Dfn.Less.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Dfn.Less.ExtraClass != null ? Data.params.Nlinks.Dfn.Less.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: [Data.params.Nlinks.Dfn.More.Text, ' +'],
                              href: Data.params.Nlinks.Dfn.More.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Dfn.More.ExtraClass != null ? Data.params.Nlinks.Dfn.More.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }],
                          className: 'input-group margin-bottom-0'
                        },
                        _owner: null
                      },
                      className: 'menu-text'
                    },
                    _owner: null
                  }],
                  className: 'float-left bar menu'
                },
                _owner: null
              }],
              className: !Data.params.Dfn.Negative ? "tabs tabs-border bar-less" : "tabs tabs-border",
              'data-tabs': true
            },
            _owner: null
          }, {
            $$typeof: _typeofReactElement,
            type: 'table',
            key: null,
            ref: null,
            props: {
              children: [{
                $$typeof: _typeofReactElement,
                type: 'thead',
                key: null,
                ref: null,
                props: {
                  children: {
                    $$typeof: _typeofReactElement,
                    type: 'tr',
                    key: null,
                    ref: null,
                    props: {
                      children: [{
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['Device', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Dfk[1 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Dfk[1 - 1].LinkHref,
                              className: Data.params.Vlinks.Dfk[1 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header '
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['Mounted', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Dfk[2 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Dfk[2 - 1].LinkHref,
                              className: Data.params.Vlinks.Dfk[2 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header '
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['Avail', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Dfk[3 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Dfk[3 - 1].LinkHref,
                              className: Data.params.Vlinks.Dfk[3 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['Use%', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Dfk[4 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Dfk[4 - 1].LinkHref,
                              className: Data.params.Vlinks.Dfk[4 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['Used', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Dfk[5 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Dfk[5 - 1].LinkHref,
                              className: Data.params.Vlinks.Dfk[5 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['Total', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Dfk[6 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Dfk[6 - 1].LinkHref,
                              className: Data.params.Vlinks.Dfk[6 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header text-right'
                        },
                        _owner: null
                      }],
                      className: 'text-nowrap'
                    },
                    _owner: null
                  }
                },
                _owner: null
              }, {
                $$typeof: _typeofReactElement,
                type: 'tbody',
                key: null,
                ref: null,
                props: {
                  children: this.List(Data).map(function ($df) {
                    return {
                      $$typeof: _typeofReactElement,
                      type: 'tr',
                      key: "df-rowby-dirname-" + $df.DirName,
                      ref: null,
                      props: {
                        children: ['  ', {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $df.DevName,
                            className: 'text-nowrap'
                          },
                          _owner: null
                        }, '  ', {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $df.DirName,
                            className: 'text-nowrap'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [{
                              $$typeof: _typeofReactElement,
                              type: 'span',
                              key: null,
                              ref: null,
                              props: {
                                children: $df.Ifree,
                                className: 'mutext',
                                title: 'Inodes free'
                              },
                              _owner: null
                            }, ' ', $df.Avail],
                            className: 'text-right text-nowrap'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [{
                              $$typeof: _typeofReactElement,
                              type: 'span',
                              key: null,
                              ref: null,
                              props: {
                                children: [$df.IusePct, '%'],
                                className: 'mutext',
                                title: 'Inodes use%'
                              },
                              _owner: null
                            }, ' ', $df.UsePct, '%'],
                            className: 'text-right bg-usepct text-nowrap',
                            'data-usepct': $df.UsePct
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [{
                              $$typeof: _typeofReactElement,
                              type: 'span',
                              key: null,
                              ref: null,
                              props: {
                                children: $df.Iused,
                                className: 'mutext',
                                title: 'Inodes used'
                              },
                              _owner: null
                            }, ' ', $df.Used],
                            className: 'text-right text-nowrap'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [{
                              $$typeof: _typeofReactElement,
                              type: 'span',
                              key: null,
                              ref: null,
                              props: {
                                children: $df.Inodes,
                                className: 'mutext',
                                title: 'Inodes total'
                              },
                              _owner: null
                            }, ' ', $df.Total],
                            className: 'text-right text-nowrap'
                          },
                          _owner: null
                        }]
                      },
                      _owner: null
                    };
                  })
                },
                _owner: null
              }],
              className: Data.params.Dfn.Absolute != 0 ? "hover scroll-x margin-bottom-0" : "hide"
            },
            _owner: null
          }]
        },
        _owner: null
      };
    }
  });

  jsdefines.define_panelif = React.createClass({
    displayName: 'define_panelif',

    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    List: function List(data) {
      var list = undefined;
      if (data != null && data["ifaddrs"] != null && (list = data["ifaddrs"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function Reduce(data) {
      return {
        params: data.params,
        ifaddrs: data.ifaddrs
      };
    },
    render: function render() {
      var Data = this.state; // shadow global Data
      return {
        $$typeof: _typeofReactElement,
        type: 'div',
        key: null,
        ref: null,
        props: {
          children: [{
            $$typeof: _typeofReactElement,
            type: 'div',
            key: null,
            ref: null,
            props: {
              children: [{
                $$typeof: _typeofReactElement,
                type: 'div',
                key: null,
                ref: null,
                props: {
                  children: {
                    $$typeof: _typeofReactElement,
                    type: 'a',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'h5',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Interfaces',
                          className: 'margin-bottom-0'
                        },
                        _owner: null
                      },
                      href: Data.params.Tlinks.Ifn,
                      onClick: this.handleClick
                    },
                    _owner: null
                  },
                  className: 'tabs-title menu-tab-padding'
                },
                _owner: null
              }, {
                $$typeof: _typeofReactElement,
                type: 'ul',
                key: null,
                ref: null,
                props: {
                  children: [{
                    $$typeof: _typeofReactElement,
                    type: 'li',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'div',
                        key: null,
                        ref: null,
                        props: {
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: 'delay',
                              className: 'input-group-label'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: Data.params.Ifd,
                              className: 'input-group-label label secondary'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['- ', Data.params.Dlinks.Ifd.Less.Text],
                              href: Data.params.Dlinks.Ifd.Less.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Ifd.Less.ExtraClass != null ? Data.params.Dlinks.Ifd.Less.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: [Data.params.Dlinks.Ifd.More.Text, ' +'],
                              href: Data.params.Dlinks.Ifd.More.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Ifd.More.ExtraClass != null ? Data.params.Dlinks.Ifd.More.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }],
                          className: 'input-group margin-bottom-0'
                        },
                        _owner: null
                      },
                      className: 'menu-text'
                    },
                    _owner: null
                  }, {
                    $$typeof: _typeofReactElement,
                    type: 'li',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'div',
                        key: null,
                        ref: null,
                        props: {
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: 'rows',
                              className: 'input-group-label'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: Data.params.Ifn.Absolute,
                              className: 'input-group-label label secondary'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['- ', Data.params.Nlinks.Ifn.Less.Text],
                              href: Data.params.Nlinks.Ifn.Less.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Ifn.Less.ExtraClass != null ? Data.params.Nlinks.Ifn.Less.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: [Data.params.Nlinks.Ifn.More.Text, ' +'],
                              href: Data.params.Nlinks.Ifn.More.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Ifn.More.ExtraClass != null ? Data.params.Nlinks.Ifn.More.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }],
                          className: 'input-group margin-bottom-0'
                        },
                        _owner: null
                      },
                      className: 'menu-text'
                    },
                    _owner: null
                  }],
                  className: 'float-left bar menu'
                },
                _owner: null
              }],
              className: !Data.params.Ifn.Negative ? "tabs tabs-border bar-less" : "tabs tabs-border",
              'data-tabs': true
            },
            _owner: null
          }, {
            $$typeof: _typeofReactElement,
            type: 'table',
            key: null,
            ref: null,
            props: {
              children: [{
                $$typeof: _typeofReactElement,
                type: 'thead',
                key: null,
                ref: null,
                props: {
                  children: {
                    $$typeof: _typeofReactElement,
                    type: 'tr',
                    key: null,
                    ref: null,
                    props: {
                      children: [{
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Interface'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'IP',
                          className: 'text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: ['IO ', {
                            $$typeof: _typeofReactElement,
                            type: 'i',
                            key: null,
                            ref: null,
                            props: {
                              children: 'b'
                            },
                            _owner: null
                          }, 'ps'],
                          className: 'text-right text-nowrap col-md-3',
                          title: 'Bits In/Out per second'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Packets IO ps',
                          className: 'text-right text-nowrap col-md-3',
                          title: 'Packets In/Out per second'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Loss IO ps',
                          className: 'text-right text-nowrap col-md-3',
                          title: 'Drops,Errors In/Out per second'
                        },
                        _owner: null
                      }]
                    },
                    _owner: null
                  }
                },
                _owner: null
              }, {
                $$typeof: _typeofReactElement,
                type: 'tbody',
                key: null,
                ref: null,
                props: {
                  children: this.List(Data).map(function ($if) {
                    return {
                      $$typeof: _typeofReactElement,
                      type: 'tr',
                      key: "if-rowby-name-" + $if.Name,
                      ref: null,
                      props: {
                        children: [{
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $if.Name,
                            className: 'text-nowrap'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $if.IP,
                            className: 'text-right'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [{
                              $$typeof: _typeofReactElement,
                              type: 'span',
                              key: null,
                              ref: null,
                              props: {
                                children: [{
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.BytesIn,
                                    title: 'Total BYTES In modulo 4G'
                                  },
                                  _owner: null
                                }, '/', {
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.BytesOut,
                                    title: 'Total BYTES Out modulo 4G'
                                  },
                                  _owner: null
                                }],
                                className: 'mutext'
                              },
                              _owner: null
                            }, ' ', {
                              $$typeof: _typeofReactElement,
                              type: 'span',
                              key: null,
                              ref: null,
                              props: {
                                children: $if.DeltaBitsIn,
                                title: 'BITS In per second'
                              },
                              _owner: null
                            }, '/', {
                              $$typeof: _typeofReactElement,
                              type: 'span',
                              key: null,
                              ref: null,
                              props: {
                                children: $if.DeltaBitsOut,
                                title: 'BITS Out per second'
                              },
                              _owner: null
                            }],
                            className: 'text-right text-nowrap'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [{
                              $$typeof: _typeofReactElement,
                              type: 'span',
                              key: null,
                              ref: null,
                              props: {
                                children: [{
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.PacketsIn,
                                    title: 'Total packets In modulo 4G'
                                  },
                                  _owner: null
                                }, '/', {
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.PacketsOut,
                                    title: 'Total packets Out modulo 4G'
                                  },
                                  _owner: null
                                }],
                                className: 'mutext'
                              },
                              _owner: null
                            }, ' ', {
                              $$typeof: _typeofReactElement,
                              type: 'span',
                              key: null,
                              ref: null,
                              props: {
                                children: $if.DeltaPacketsIn,
                                title: 'Packets In per second'
                              },
                              _owner: null
                            }, '/', {
                              $$typeof: _typeofReactElement,
                              type: 'span',
                              key: null,
                              ref: null,
                              props: {
                                children: $if.DeltaPacketsOut,
                                title: 'Packets Out per second'
                              },
                              _owner: null
                            }],
                            className: 'text-right text-nowrap'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [{
                              $$typeof: _typeofReactElement,
                              type: 'span',
                              key: null,
                              ref: null,
                              props: {
                                children: [{
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.DropsIn,
                                    title: 'Total drops In modulo 4G'
                                  },
                                  _owner: null
                                }, {
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: '/',
                                    className: $if.DropsOut != null ? "" : "hide"
                                  },
                                  _owner: null
                                }, {
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.DropsOut,
                                    className: $if.DropsOut != null ? "" : "hide",
                                    title: 'Total drops Out modulo 4G'
                                  },
                                  _owner: null
                                }, ',', {
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.ErrorsIn,
                                    title: 'Total errors In modulo 4G'
                                  },
                                  _owner: null
                                }, '/', {
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.ErrorsOut,
                                    title: 'Total errors Out modulo 4G'
                                  },
                                  _owner: null
                                }],
                                className: 'mutext',
                                title: 'Total drops,errors modulo 4G'
                              },
                              _owner: null
                            }, ' ', {
                              $$typeof: _typeofReactElement,
                              type: 'span',
                              key: null,
                              ref: null,
                              props: {
                                children: [{
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.DeltaDropsIn,
                                    title: 'Drops In per second'
                                  },
                                  _owner: null
                                }, {
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: '/',
                                    className: $if.DeltaDropsOut != null ? "" : "hide"
                                  },
                                  _owner: null
                                }, {
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.DeltaDropsOut,
                                    className: $if.DeltaDropsOut != null ? "" : "hide",
                                    title: 'Drops Out per second'
                                  },
                                  _owner: null
                                }, ',', {
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.DeltaErrorsIn,
                                    title: 'Errors In per second'
                                  },
                                  _owner: null
                                }, '/', {
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.DeltaErrorsOut,
                                    title: 'Errors Out per second'
                                  },
                                  _owner: null
                                }],
                                className: ($if.DeltaDropsIn == null || $if.DeltaDropsIn == "0") && ($if.DeltaDropsOut == null || $if.DeltaDropsOut == "0") && ($if.DeltaErrorsIn == null || $if.DeltaErrorsIn == "0") && ($if.DeltaErrorsOut == null || $if.DeltaErrorsOut == "0") ? "mutext" : ""
                              },
                              _owner: null
                            }],
                            className: 'text-right text-nowrap'
                          },
                          _owner: null
                        }]
                      },
                      _owner: null
                    };
                  })
                },
                _owner: null
              }],
              className: Data.params.Ifn.Absolute != 0 ? "hover scroll-x margin-bottom-0" : "hide"
            },
            _owner: null
          }]
        },
        _owner: null
      };
    }
  });

  jsdefines.define_panelmem = React.createClass({
    displayName: 'define_panelmem',

    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    List: function List(data) {
      var list = undefined;
      if (data != null && data["memory"] != null && (list = data["memory"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function Reduce(data) {
      return {
        params: data.params,
        memory: data.memory
      };
    },
    render: function render() {
      var Data = this.state; // shadow global Data
      return {
        $$typeof: _typeofReactElement,
        type: 'div',
        key: null,
        ref: null,
        props: {
          children: [{
            $$typeof: _typeofReactElement,
            type: 'div',
            key: null,
            ref: null,
            props: {
              children: [{
                $$typeof: _typeofReactElement,
                type: 'div',
                key: null,
                ref: null,
                props: {
                  children: {
                    $$typeof: _typeofReactElement,
                    type: 'a',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'h5',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Memory',
                          className: 'margin-bottom-0'
                        },
                        _owner: null
                      },
                      href: Data.params.Tlinks.Memn,
                      onClick: this.handleClick
                    },
                    _owner: null
                  },
                  className: 'tabs-title menu-tab-padding'
                },
                _owner: null
              }, {
                $$typeof: _typeofReactElement,
                type: 'ul',
                key: null,
                ref: null,
                props: {
                  children: [{
                    $$typeof: _typeofReactElement,
                    type: 'li',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'div',
                        key: null,
                        ref: null,
                        props: {
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: 'delay',
                              className: 'input-group-label'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: Data.params.Memd,
                              className: 'input-group-label label secondary'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['- ', Data.params.Dlinks.Memd.Less.Text],
                              href: Data.params.Dlinks.Memd.Less.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Memd.Less.ExtraClass != null ? Data.params.Dlinks.Memd.Less.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: [Data.params.Dlinks.Memd.More.Text, ' +'],
                              href: Data.params.Dlinks.Memd.More.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Memd.More.ExtraClass != null ? Data.params.Dlinks.Memd.More.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }],
                          className: 'input-group margin-bottom-0'
                        },
                        _owner: null
                      },
                      className: 'menu-text'
                    },
                    _owner: null
                  }, {
                    $$typeof: _typeofReactElement,
                    type: 'li',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'div',
                        key: null,
                        ref: null,
                        props: {
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: 'rows',
                              className: 'input-group-label'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: Data.params.Memn.Absolute,
                              className: 'input-group-label label secondary'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['- ', Data.params.Nlinks.Memn.Less.Text],
                              href: Data.params.Nlinks.Memn.Less.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Memn.Less.ExtraClass != null ? Data.params.Nlinks.Memn.Less.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: [Data.params.Nlinks.Memn.More.Text, ' +'],
                              href: Data.params.Nlinks.Memn.More.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Memn.More.ExtraClass != null ? Data.params.Nlinks.Memn.More.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }],
                          className: 'input-group margin-bottom-0'
                        },
                        _owner: null
                      },
                      className: 'menu-text'
                    },
                    _owner: null
                  }],
                  className: 'float-left bar menu'
                },
                _owner: null
              }],
              className: !Data.params.Memn.Negative ? "tabs tabs-border bar-less" : "tabs tabs-border",
              'data-tabs': true
            },
            _owner: null
          }, {
            $$typeof: _typeofReactElement,
            type: 'table',
            key: null,
            ref: null,
            props: {
              children: [{
                $$typeof: _typeofReactElement,
                type: 'thead',
                key: null,
                ref: null,
                props: {
                  children: {
                    $$typeof: _typeofReactElement,
                    type: 'tr',
                    key: null,
                    ref: null,
                    props: {
                      children: [{
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {},
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Free',
                          className: 'text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Use%',
                          className: 'text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Used',
                          className: 'text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Total',
                          className: 'text-right'
                        },
                        _owner: null
                      }]
                    },
                    _owner: null
                  }
                },
                _owner: null
              }, {
                $$typeof: _typeofReactElement,
                type: 'tbody',
                key: null,
                ref: null,
                props: {
                  children: this.List(Data).map(function ($mem) {
                    return {
                      $$typeof: _typeofReactElement,
                      type: 'tr',
                      key: "mem-rowby-kind-" + $mem.Kind,
                      ref: null,
                      props: {
                        children: [{
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $mem.Kind
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $mem.Free,
                            className: 'text-right'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [$mem.UsePct, '%'],
                            className: 'text-right bg-usepct',
                            'data-usepct': $mem.UsePct
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $mem.Used,
                            className: 'text-right'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $mem.Total,
                            className: 'text-right'
                          },
                          _owner: null
                        }]
                      },
                      _owner: null
                    };
                  })
                },
                _owner: null
              }],
              className: Data.params.Memn.Absolute != 0 ? "hover scroll-x margin-bottom-0" : "hide"
            },
            _owner: null
          }]
        },
        _owner: null
      };
    }
  });

  jsdefines.define_panelps = React.createClass({
    displayName: 'define_panelps',

    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    List: function List(data) {
      var list = undefined;
      if (data != null && data["procs"] != null && (list = data["procs"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function Reduce(data) {
      return {
        params: data.params,
        procs: data.procs
      };
    },
    render: function render() {
      var Data = this.state; // shadow global Data
      return {
        $$typeof: _typeofReactElement,
        type: 'div',
        key: null,
        ref: null,
        props: {
          children: [{
            $$typeof: _typeofReactElement,
            type: 'div',
            key: null,
            ref: null,
            props: {
              children: [{
                $$typeof: _typeofReactElement,
                type: 'div',
                key: null,
                ref: null,
                props: {
                  children: {
                    $$typeof: _typeofReactElement,
                    type: 'a',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'h5',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Processes',
                          className: 'margin-bottom-0'
                        },
                        _owner: null
                      },
                      href: Data.params.Tlinks.Psn,
                      onClick: this.handleClick
                    },
                    _owner: null
                  },
                  className: 'tabs-title menu-tab-padding'
                },
                _owner: null
              }, {
                $$typeof: _typeofReactElement,
                type: 'ul',
                key: null,
                ref: null,
                props: {
                  children: [{
                    $$typeof: _typeofReactElement,
                    type: 'li',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'div',
                        key: null,
                        ref: null,
                        props: {
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: 'delay',
                              className: 'input-group-label'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: Data.params.Psd,
                              className: 'input-group-label label secondary'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['- ', Data.params.Dlinks.Psd.Less.Text],
                              href: Data.params.Dlinks.Psd.Less.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Psd.Less.ExtraClass != null ? Data.params.Dlinks.Psd.Less.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: [Data.params.Dlinks.Psd.More.Text, ' +'],
                              href: Data.params.Dlinks.Psd.More.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Dlinks.Psd.More.ExtraClass != null ? Data.params.Dlinks.Psd.More.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }],
                          className: 'input-group margin-bottom-0'
                        },
                        _owner: null
                      },
                      className: 'menu-text'
                    },
                    _owner: null
                  }, {
                    $$typeof: _typeofReactElement,
                    type: 'li',
                    key: null,
                    ref: null,
                    props: {
                      children: {
                        $$typeof: _typeofReactElement,
                        type: 'div',
                        key: null,
                        ref: null,
                        props: {
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: 'rows',
                              className: 'input-group-label'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: Data.params.Psn.Absolute,
                              className: 'input-group-label label secondary'
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['- ', Data.params.Nlinks.Psn.Less.Text],
                              href: Data.params.Nlinks.Psn.Less.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Psn.Less.ExtraClass != null ? Data.params.Nlinks.Psn.Less.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }, {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: [Data.params.Nlinks.Psn.More.Text, ' +'],
                              href: Data.params.Nlinks.Psn.More.Href,
                              className: "button secondary hollow text-nowrap input-group-button" + " " + (Data.params.Nlinks.Psn.More.ExtraClass != null ? Data.params.Nlinks.Psn.More.ExtraClass : ""),
                              onClick: this.handleClick
                            },
                            _owner: null
                          }],
                          className: 'input-group margin-bottom-0'
                        },
                        _owner: null
                      },
                      className: 'menu-text'
                    },
                    _owner: null
                  }],
                  className: 'float-left bar menu'
                },
                _owner: null
              }],
              className: !Data.params.Psn.Negative ? "tabs tabs-border bar-less" : "tabs tabs-border",
              'data-tabs': true
            },
            _owner: null
          }, {
            $$typeof: _typeofReactElement,
            type: 'table',
            key: null,
            ref: null,
            props: {
              children: [{
                $$typeof: _typeofReactElement,
                type: 'thead',
                key: null,
                ref: null,
                props: {
                  children: {
                    $$typeof: _typeofReactElement,
                    type: 'tr',
                    key: null,
                    ref: null,
                    props: {
                      children: [{
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['PID', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Psk[1 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Psk[1 - 1].LinkHref,
                              className: Data.params.Vlinks.Psk[1 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['UID', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Psk[2 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Psk[2 - 1].LinkHref,
                              className: Data.params.Vlinks.Psk[2 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['USER', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Psk[3 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Psk[3 - 1].LinkHref,
                              className: Data.params.Vlinks.Psk[3 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header '
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['PR', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Psk[4 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Psk[4 - 1].LinkHref,
                              className: Data.params.Vlinks.Psk[4 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['NI', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Psk[5 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Psk[5 - 1].LinkHref,
                              className: Data.params.Vlinks.Psk[5 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['VIRT', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Psk[6 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Psk[6 - 1].LinkHref,
                              className: Data.params.Vlinks.Psk[6 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['RES', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Psk[7 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Psk[7 - 1].LinkHref,
                              className: Data.params.Vlinks.Psk[7 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['TIME', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Psk[8 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Psk[8 - 1].LinkHref,
                              className: Data.params.Vlinks.Psk[8 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header text-center'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: {
                            $$typeof: _typeofReactElement,
                            type: 'a',
                            key: null,
                            ref: null,
                            props: {
                              children: ['COMMAND', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  className: Data.params.Vlinks.Psk[9 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.params.Vlinks.Psk[9 - 1].LinkHref,
                              className: Data.params.Vlinks.Psk[9 - 1].LinkClass,
                              onClick: this.handleClick
                            },
                            _owner: null
                          },
                          className: 'header '
                        },
                        _owner: null
                      }],
                      className: 'text-nowrap'
                    },
                    _owner: null
                  }
                },
                _owner: null
              }, {
                $$typeof: _typeofReactElement,
                type: 'tbody',
                key: null,
                ref: null,
                props: {
                  children: this.List(Data).map(function ($ps) {
                    return {
                      $$typeof: _typeofReactElement,
                      type: 'tr',
                      key: "ps-rowby-pid-" + $ps.PID,
                      ref: null,
                      props: {
                        children: [{
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [' ', $ps.PID],
                            className: 'text-right'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [' ', $ps.UID],
                            className: 'text-right'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $ps.User
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [' ', $ps.Priority],
                            className: 'text-right'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [' ', $ps.Nice],
                            className: 'text-right'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [' ', $ps.Size],
                            className: 'text-right'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: [' ', $ps.Resident],
                            className: 'text-right'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $ps.Time,
                            className: 'text-center'
                          },
                          _owner: null
                        }, {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $ps.Name
                          },
                          _owner: null
                        }]
                      },
                      _owner: null
                    };
                  })
                },
                _owner: null
              }],
              className: Data.params.Psn.Absolute != 0 ? "hover scroll-x margin-bottom-0" : "hide"
            },
            _owner: null
          }]
        },
        _owner: null
      };
    }
  });

  jsdefines.define_loadavg = React.createClass({
    displayName: 'define_loadavg',

    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    Reduce: function Reduce(data) {
      return {
        loadavg: data.loadavg
      };
    },
    render: function render() {
      var Data = this.state; // shadow global Data
      return {
        $$typeof: _typeofReactElement,
        type: 'span',
        key: null,
        ref: null,
        props: {
          children: Data.loadavg
        },
        _owner: null
      };
    }
  });

  jsdefines.define_uptime = React.createClass({
    displayName: 'define_uptime',

    mixins: [React.addons.PureRenderMixin, jsdefines.StateHandlingMixin, jsdefines.HandlerMixin],
    Reduce: function Reduce(data) {
      return {
        uptime: data.uptime
      };
    },
    render: function render() {
      var Data = this.state; // shadow global Data
      return {
        $$typeof: _typeofReactElement,
        type: 'span',
        key: null,
        ref: null,
        props: {
          children: Data.uptime
        },
        _owner: null
      };
    }
  });

  return jsdefines;
});
