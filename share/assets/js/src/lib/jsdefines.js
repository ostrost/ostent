'use strict';

var _typeofReactElement = typeof Symbol === 'function' && Symbol['for'] && Symbol['for']('react.element') || 60103;

define(function (require) {
  var React = require('react');
  var jsdefines = {};
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
  // all the define_* templates transformed into jsdefines.define_* = ...;

  jsdefines.define_panelcpu = React.createClass({
    displayName: 'define_panelcpu',

    mixins: [React.addons.PureRenderMixin, jsdefines.HandlerMixin],
    List: function List(data) {
      // static
      var list;
      if (data != null && data["CPU"] != null && (list = data["CPU"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function Reduce(data) {
      // static
      return {
        Params: data.Params,
        CPU: data.CPU
      };
    },
    getInitialState: function getInitialState() {
      return this.Reduce(Data); // global Data
    },
    render: function render() {
      var Data = this.state;
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
              children: {
                $$typeof: _typeofReactElement,
                type: 'a',
                key: null,
                ref: null,
                props: {
                  children: 'CPU',
                  href: Data.Params.Tlinks.CPUn,
                  onClick: this.handleClick
                },
                _owner: null
              },
              className: 'h4 padding-left-like-panel-heading'
            },
            _owner: null
          }, {
            $$typeof: _typeofReactElement,
            type: 'ul',
            key: null,
            ref: null,
            props: {
              children: {
                $$typeof: _typeofReactElement,
                type: 'li',
                key: null,
                ref: null,
                props: {
                  children: {
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
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'b',
                                key: null,
                                ref: null,
                                props: {
                                  children: 'Delay'
                                },
                                _owner: null
                              }, ' ', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  children: Data.Params.CPUd,
                                  className: 'badge'
                                },
                                _owner: null
                              }]
                            },
                            _owner: null
                          }, ' ', {
                            $$typeof: _typeofReactElement,
                            type: 'div',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: ['- ', Data.Params.Dlinks.CPUd.Less.Text],
                                  href: Data.Params.Dlinks.CPUd.Less.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Dlinks.CPUd.Less.ExtraClass != null ? Data.Params.Dlinks.CPUd.Less.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }, {
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: [Data.Params.Dlinks.CPUd.More.Text, ' +'],
                                  href: Data.Params.Dlinks.CPUd.More.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Dlinks.CPUd.More.ExtraClass != null ? Data.Params.Dlinks.CPUd.More.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }],
                              className: 'btn-group'
                            },
                            _owner: null
                          }]
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'li',
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
                                type: 'b',
                                key: null,
                                ref: null,
                                props: {
                                  children: 'Rows'
                                },
                                _owner: null
                              }, ' ', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  children: Data.Params.CPUn.Absolute,
                                  className: 'badge'
                                },
                                _owner: null
                              }]
                            },
                            _owner: null
                          }, ' ', {
                            $$typeof: _typeofReactElement,
                            type: 'div',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: ['- ', Data.Params.Nlinks.CPUn.Less.Text],
                                  href: Data.Params.Nlinks.CPUn.Less.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Nlinks.CPUn.Less.ExtraClass != null ? Data.Params.Nlinks.CPUn.Less.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }, {
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: [Data.Params.Nlinks.CPUn.More.Text, ' +'],
                                  href: Data.Params.Nlinks.CPUn.More.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Nlinks.CPUn.More.ExtraClass != null ? Data.Params.Nlinks.CPUn.More.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }],
                              className: 'btn-group'
                            },
                            _owner: null
                          }]
                        },
                        _owner: null
                      }],
                      className: 'list-inline'
                    },
                    _owner: null
                  },
                  className: 'list-group-item text-nowrap th'
                },
                _owner: null
              },
              className: !Data.Params.CPUn.Negative ? "hidden" : "list-group"
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
                          children: 'User',
                          className: 'text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Sys',
                          className: 'text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Wait',
                          className: 'text-right'
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'th',
                        key: null,
                        ref: null,
                        props: {
                          children: 'Idle',
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
              className: Data.Params.CPUn.Absolute != 0 ? "table table-hover" : "hidden"
            },
            _owner: null
          }],
          className: !Data.Params.CPUn.Negative ? "" : "panel panel-default"
        },
        _owner: null
      };
    }
  });

  jsdefines.define_paneldf = React.createClass({
    displayName: 'define_paneldf',

    mixins: [React.addons.PureRenderMixin, jsdefines.HandlerMixin],
    List: function List(data) {
      // static
      var list;
      if (data != null && data["DF"] != null && (list = data["DF"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function Reduce(data) {
      // static
      return {
        Params: data.Params,
        DF: data.DF
      };
    },
    getInitialState: function getInitialState() {
      return this.Reduce(Data); // global Data
    },
    render: function render() {
      var Data = this.state;
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
              children: {
                $$typeof: _typeofReactElement,
                type: 'a',
                key: null,
                ref: null,
                props: {
                  children: 'Disk usage',
                  href: Data.Params.Tlinks.Dfn,
                  onClick: this.handleClick
                },
                _owner: null
              },
              className: 'h4 padding-left-like-panel-heading'
            },
            _owner: null
          }, {
            $$typeof: _typeofReactElement,
            type: 'ul',
            key: null,
            ref: null,
            props: {
              children: {
                $$typeof: _typeofReactElement,
                type: 'li',
                key: null,
                ref: null,
                props: {
                  children: {
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
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'b',
                                key: null,
                                ref: null,
                                props: {
                                  children: 'Delay'
                                },
                                _owner: null
                              }, ' ', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  children: Data.Params.Dfd,
                                  className: 'badge'
                                },
                                _owner: null
                              }]
                            },
                            _owner: null
                          }, ' ', {
                            $$typeof: _typeofReactElement,
                            type: 'div',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: ['- ', Data.Params.Dlinks.Dfd.Less.Text],
                                  href: Data.Params.Dlinks.Dfd.Less.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Dlinks.Dfd.Less.ExtraClass != null ? Data.Params.Dlinks.Dfd.Less.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }, {
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: [Data.Params.Dlinks.Dfd.More.Text, ' +'],
                                  href: Data.Params.Dlinks.Dfd.More.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Dlinks.Dfd.More.ExtraClass != null ? Data.Params.Dlinks.Dfd.More.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }],
                              className: 'btn-group'
                            },
                            _owner: null
                          }]
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'li',
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
                                type: 'b',
                                key: null,
                                ref: null,
                                props: {
                                  children: 'Rows'
                                },
                                _owner: null
                              }, ' ', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  children: Data.Params.Dfn.Absolute,
                                  className: 'badge'
                                },
                                _owner: null
                              }]
                            },
                            _owner: null
                          }, ' ', {
                            $$typeof: _typeofReactElement,
                            type: 'div',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: ['- ', Data.Params.Nlinks.Dfn.Less.Text],
                                  href: Data.Params.Nlinks.Dfn.Less.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Nlinks.Dfn.Less.ExtraClass != null ? Data.Params.Nlinks.Dfn.Less.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }, {
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: [Data.Params.Nlinks.Dfn.More.Text, ' +'],
                                  href: Data.Params.Nlinks.Dfn.More.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Nlinks.Dfn.More.ExtraClass != null ? Data.Params.Nlinks.Dfn.More.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }],
                              className: 'btn-group'
                            },
                            _owner: null
                          }]
                        },
                        _owner: null
                      }],
                      className: 'list-inline'
                    },
                    _owner: null
                  },
                  className: 'list-group-item text-nowrap th'
                },
                _owner: null
              },
              className: !Data.Params.Dfn.Negative ? "hidden" : "list-group"
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
                                  className: Data.Params.Vlinks.Dfk[1 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Dfk[1 - 1].LinkHref,
                              className: Data.Params.Vlinks.Dfk[1 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Dfk[2 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Dfk[2 - 1].LinkHref,
                              className: Data.Params.Vlinks.Dfk[2 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Dfk[3 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Dfk[3 - 1].LinkHref,
                              className: Data.Params.Vlinks.Dfk[3 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Dfk[4 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Dfk[4 - 1].LinkHref,
                              className: Data.Params.Vlinks.Dfk[4 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Dfk[5 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Dfk[5 - 1].LinkHref,
                              className: Data.Params.Vlinks.Dfk[5 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Dfk[6 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Dfk[6 - 1].LinkHref,
                              className: Data.Params.Vlinks.Dfk[6 - 1].LinkClass,
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
                            className: 'text-nowrap clip12',
                            title: $df.DevName
                          },
                          _owner: null
                        }, '  ', {
                          $$typeof: _typeofReactElement,
                          type: 'td',
                          key: null,
                          ref: null,
                          props: {
                            children: $df.DirName,
                            className: 'text-nowrap clip12',
                            title: $df.DirName
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
              className: Data.Params.Dfn.Absolute != 0 ? "table table-hover" : "hidden"
            },
            _owner: null
          }],
          className: !Data.Params.Dfn.Negative ? "" : "panel panel-default"
        },
        _owner: null
      };
    }
  });

  jsdefines.define_panelif = React.createClass({
    displayName: 'define_panelif',

    mixins: [React.addons.PureRenderMixin, jsdefines.HandlerMixin],
    List: function List(data) {
      // static
      var list;
      if (data != null && data["IF"] != null && (list = data["IF"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function Reduce(data) {
      // static
      return {
        Params: data.Params,
        IF: data.IF
      };
    },
    getInitialState: function getInitialState() {
      return this.Reduce(Data); // global Data
    },
    render: function render() {
      var Data = this.state;
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
              children: {
                $$typeof: _typeofReactElement,
                type: 'a',
                key: null,
                ref: null,
                props: {
                  children: 'Interfaces',
                  href: Data.Params.Tlinks.Ifn,
                  onClick: this.handleClick
                },
                _owner: null
              },
              className: 'h4 padding-left-like-panel-heading'
            },
            _owner: null
          }, {
            $$typeof: _typeofReactElement,
            type: 'ul',
            key: null,
            ref: null,
            props: {
              children: {
                $$typeof: _typeofReactElement,
                type: 'li',
                key: null,
                ref: null,
                props: {
                  children: {
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
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'b',
                                key: null,
                                ref: null,
                                props: {
                                  children: 'Delay'
                                },
                                _owner: null
                              }, ' ', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  children: Data.Params.Ifd,
                                  className: 'badge'
                                },
                                _owner: null
                              }]
                            },
                            _owner: null
                          }, ' ', {
                            $$typeof: _typeofReactElement,
                            type: 'div',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: ['- ', Data.Params.Dlinks.Ifd.Less.Text],
                                  href: Data.Params.Dlinks.Ifd.Less.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Dlinks.Ifd.Less.ExtraClass != null ? Data.Params.Dlinks.Ifd.Less.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }, {
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: [Data.Params.Dlinks.Ifd.More.Text, ' +'],
                                  href: Data.Params.Dlinks.Ifd.More.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Dlinks.Ifd.More.ExtraClass != null ? Data.Params.Dlinks.Ifd.More.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }],
                              className: 'btn-group'
                            },
                            _owner: null
                          }]
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'li',
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
                                type: 'b',
                                key: null,
                                ref: null,
                                props: {
                                  children: 'Rows'
                                },
                                _owner: null
                              }, ' ', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  children: Data.Params.Ifn.Absolute,
                                  className: 'badge'
                                },
                                _owner: null
                              }]
                            },
                            _owner: null
                          }, ' ', {
                            $$typeof: _typeofReactElement,
                            type: 'div',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: ['- ', Data.Params.Nlinks.Ifn.Less.Text],
                                  href: Data.Params.Nlinks.Ifn.Less.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Nlinks.Ifn.Less.ExtraClass != null ? Data.Params.Nlinks.Ifn.Less.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }, {
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: [Data.Params.Nlinks.Ifn.More.Text, ' +'],
                                  href: Data.Params.Nlinks.Ifn.More.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Nlinks.Ifn.More.ExtraClass != null ? Data.Params.Nlinks.Ifn.More.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }],
                              className: 'btn-group'
                            },
                            _owner: null
                          }]
                        },
                        _owner: null
                      }],
                      className: 'list-inline'
                    },
                    _owner: null
                  },
                  className: 'list-group-item text-nowrap th'
                },
                _owner: null
              },
              className: !Data.Params.Ifn.Negative ? "hidden" : "list-group"
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
                            className: 'text-nowrap clip12',
                            title: $if.Name
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
                                    className: $if.DropsOut != null ? "" : "hidden"
                                  },
                                  _owner: null
                                }, {
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.DropsOut,
                                    className: $if.DropsOut != null ? "" : "hidden",
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
                                    className: $if.DeltaDropsOut != null ? "" : "hidden"
                                  },
                                  _owner: null
                                }, {
                                  $$typeof: _typeofReactElement,
                                  type: 'span',
                                  key: null,
                                  ref: null,
                                  props: {
                                    children: $if.DeltaDropsOut,
                                    className: $if.DeltaDropsOut != null ? "" : "hidden",
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
              className: Data.Params.Ifn.Absolute != 0 ? "table table-hover" : "hidden"
            },
            _owner: null
          }],
          className: !Data.Params.Ifn.Negative ? "" : "panel panel-default"
        },
        _owner: null
      };
    }
  });

  jsdefines.define_panelmem = React.createClass({
    displayName: 'define_panelmem',

    mixins: [React.addons.PureRenderMixin, jsdefines.HandlerMixin],
    List: function List(data) {
      // static
      var list;
      if (data != null && data["MEM"] != null && (list = data["MEM"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function Reduce(data) {
      // static
      return {
        Params: data.Params,
        MEM: data.MEM
      };
    },
    getInitialState: function getInitialState() {
      return this.Reduce(Data); // global Data
    },
    render: function render() {
      var Data = this.state;
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
              children: {
                $$typeof: _typeofReactElement,
                type: 'a',
                key: null,
                ref: null,
                props: {
                  children: 'Memory',
                  href: Data.Params.Tlinks.Memn,
                  onClick: this.handleClick
                },
                _owner: null
              },
              className: 'h4 padding-left-like-panel-heading'
            },
            _owner: null
          }, {
            $$typeof: _typeofReactElement,
            type: 'ul',
            key: null,
            ref: null,
            props: {
              children: {
                $$typeof: _typeofReactElement,
                type: 'li',
                key: null,
                ref: null,
                props: {
                  children: {
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
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'b',
                                key: null,
                                ref: null,
                                props: {
                                  children: 'Delay'
                                },
                                _owner: null
                              }, ' ', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  children: Data.Params.Memd,
                                  className: 'badge'
                                },
                                _owner: null
                              }]
                            },
                            _owner: null
                          }, ' ', {
                            $$typeof: _typeofReactElement,
                            type: 'div',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: ['- ', Data.Params.Dlinks.Memd.Less.Text],
                                  href: Data.Params.Dlinks.Memd.Less.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Dlinks.Memd.Less.ExtraClass != null ? Data.Params.Dlinks.Memd.Less.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }, {
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: [Data.Params.Dlinks.Memd.More.Text, ' +'],
                                  href: Data.Params.Dlinks.Memd.More.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Dlinks.Memd.More.ExtraClass != null ? Data.Params.Dlinks.Memd.More.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }],
                              className: 'btn-group'
                            },
                            _owner: null
                          }]
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'li',
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
                                type: 'b',
                                key: null,
                                ref: null,
                                props: {
                                  children: 'Rows'
                                },
                                _owner: null
                              }, ' ', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  children: Data.Params.Memn.Absolute,
                                  className: 'badge'
                                },
                                _owner: null
                              }]
                            },
                            _owner: null
                          }, ' ', {
                            $$typeof: _typeofReactElement,
                            type: 'div',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: ['- ', Data.Params.Nlinks.Memn.Less.Text],
                                  href: Data.Params.Nlinks.Memn.Less.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Nlinks.Memn.Less.ExtraClass != null ? Data.Params.Nlinks.Memn.Less.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }, {
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: [Data.Params.Nlinks.Memn.More.Text, ' +'],
                                  href: Data.Params.Nlinks.Memn.More.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Nlinks.Memn.More.ExtraClass != null ? Data.Params.Nlinks.Memn.More.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }],
                              className: 'btn-group'
                            },
                            _owner: null
                          }]
                        },
                        _owner: null
                      }],
                      className: 'list-inline'
                    },
                    _owner: null
                  },
                  className: 'list-group-item text-nowrap th'
                },
                _owner: null
              },
              className: !Data.Params.Memn.Negative ? "hidden" : "list-group"
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
              className: Data.Params.Memn.Absolute != 0 ? "table table-hover" : "hidden"
            },
            _owner: null
          }],
          className: !Data.Params.Memn.Negative ? "" : "panel panel-default"
        },
        _owner: null
      };
    }
  });

  jsdefines.define_panelps = React.createClass({
    displayName: 'define_panelps',

    mixins: [React.addons.PureRenderMixin, jsdefines.HandlerMixin],
    List: function List(data) {
      // static
      var list;
      if (data != null && data["PS"] != null && (list = data["PS"].List) != null) {
        return list;
      }
      return [];
    },
    Reduce: function Reduce(data) {
      // static
      return {
        Params: data.Params,
        PS: data.PS
      };
    },
    getInitialState: function getInitialState() {
      return this.Reduce(Data); // global Data
    },
    render: function render() {
      var Data = this.state;
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
              children: {
                $$typeof: _typeofReactElement,
                type: 'a',
                key: null,
                ref: null,
                props: {
                  children: 'Processes',
                  href: Data.Params.Tlinks.Psn,
                  onClick: this.handleClick
                },
                _owner: null
              },
              className: 'h4 padding-left-like-panel-heading'
            },
            _owner: null
          }, {
            $$typeof: _typeofReactElement,
            type: 'ul',
            key: null,
            ref: null,
            props: {
              children: {
                $$typeof: _typeofReactElement,
                type: 'li',
                key: null,
                ref: null,
                props: {
                  children: {
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
                          children: [{
                            $$typeof: _typeofReactElement,
                            type: 'span',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'b',
                                key: null,
                                ref: null,
                                props: {
                                  children: 'Delay'
                                },
                                _owner: null
                              }, ' ', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  children: Data.Params.Psd,
                                  className: 'badge'
                                },
                                _owner: null
                              }]
                            },
                            _owner: null
                          }, ' ', {
                            $$typeof: _typeofReactElement,
                            type: 'div',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: ['- ', Data.Params.Dlinks.Psd.Less.Text],
                                  href: Data.Params.Dlinks.Psd.Less.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Dlinks.Psd.Less.ExtraClass != null ? Data.Params.Dlinks.Psd.Less.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }, {
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: [Data.Params.Dlinks.Psd.More.Text, ' +'],
                                  href: Data.Params.Dlinks.Psd.More.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Dlinks.Psd.More.ExtraClass != null ? Data.Params.Dlinks.Psd.More.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }],
                              className: 'btn-group'
                            },
                            _owner: null
                          }]
                        },
                        _owner: null
                      }, {
                        $$typeof: _typeofReactElement,
                        type: 'li',
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
                                type: 'b',
                                key: null,
                                ref: null,
                                props: {
                                  children: 'Rows'
                                },
                                _owner: null
                              }, ' ', {
                                $$typeof: _typeofReactElement,
                                type: 'span',
                                key: null,
                                ref: null,
                                props: {
                                  children: Data.Params.Psn.Absolute,
                                  className: 'badge'
                                },
                                _owner: null
                              }]
                            },
                            _owner: null
                          }, ' ', {
                            $$typeof: _typeofReactElement,
                            type: 'div',
                            key: null,
                            ref: null,
                            props: {
                              children: [{
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: ['- ', Data.Params.Nlinks.Psn.Less.Text],
                                  href: Data.Params.Nlinks.Psn.Less.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Nlinks.Psn.Less.ExtraClass != null ? Data.Params.Nlinks.Psn.Less.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }, {
                                $$typeof: _typeofReactElement,
                                type: 'a',
                                key: null,
                                ref: null,
                                props: {
                                  children: [Data.Params.Nlinks.Psn.More.Text, ' +'],
                                  href: Data.Params.Nlinks.Psn.More.Href,
                                  className: "btn btn-default" + " " + (Data.Params.Nlinks.Psn.More.ExtraClass != null ? Data.Params.Nlinks.Psn.More.ExtraClass : ""),
                                  onClick: this.handleClick
                                },
                                _owner: null
                              }],
                              className: 'btn-group'
                            },
                            _owner: null
                          }]
                        },
                        _owner: null
                      }],
                      className: 'list-inline'
                    },
                    _owner: null
                  },
                  className: 'list-group-item text-nowrap th'
                },
                _owner: null
              },
              className: !Data.Params.Psn.Negative ? "hidden" : "list-group"
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
                                  className: Data.Params.Vlinks.Psk[1 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Psk[1 - 1].LinkHref,
                              className: Data.Params.Vlinks.Psk[1 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Psk[2 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Psk[2 - 1].LinkHref,
                              className: Data.Params.Vlinks.Psk[2 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Psk[3 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Psk[3 - 1].LinkHref,
                              className: Data.Params.Vlinks.Psk[3 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Psk[4 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Psk[4 - 1].LinkHref,
                              className: Data.Params.Vlinks.Psk[4 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Psk[5 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Psk[5 - 1].LinkHref,
                              className: Data.Params.Vlinks.Psk[5 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Psk[6 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Psk[6 - 1].LinkHref,
                              className: Data.Params.Vlinks.Psk[6 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Psk[7 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Psk[7 - 1].LinkHref,
                              className: Data.Params.Vlinks.Psk[7 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Psk[8 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Psk[8 - 1].LinkHref,
                              className: Data.Params.Vlinks.Psk[8 - 1].LinkClass,
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
                                  className: Data.Params.Vlinks.Psk[9 - 1].CaretClass
                                },
                                _owner: null
                              }],
                              href: Data.Params.Vlinks.Psk[9 - 1].LinkHref,
                              className: Data.Params.Vlinks.Psk[9 - 1].LinkClass,
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
              className: Data.Params.Psn.Absolute != 0 ? "table table-hover" : "hidden"
            },
            _owner: null
          }],
          className: !Data.Params.Psn.Negative ? "" : "panel panel-default"
        },
        _owner: null
      };
    }
  });
  return jsdefines;
});
