"use strict";

define(function (require) {
  var React = require("react");
  return {
    mem_rows: function mem_rows(Data, $mem) {
      return React.createElement(
        "tr",
        { key: "mem-rowby-kind-" + $mem.Kind
        },
        React.createElement(
          "td",
          null,
          $mem.Kind
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $mem.Free
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          React.createElement(
            "span",
            { className: "label", "data-usepercent": $mem.UsePercent
            },
            $mem.UsePercent,
            "%"
          ),
          " ",
          $mem.Used
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $mem.Total
        )
      );
    },
    panelmem: function panelmem(Data, rows) {
      return React.createElement(
        "div",
        { className: !Data.Params.Memn.Negative ? "" : "panel panel-default"
        },
        "  ",
        React.createElement(
          "div",
          { className: !Data.Params.Memn.Negative ? "" : "panel-heading"
          },
          "    ",
          React.createElement(
            "a",
            { href: Data.Params.Tlinks.Memn, onClick: this.handleClick, className: "panel-title btn-block"
            },
            "      ",
            React.createElement(
              "b",
              { className: !Data.Params.Memn.Negative ? "h4" : "h4 bg-info"
              },
              "Memory"
            ),
            "    "
          ),
          "  "
        ),
        "  ",
        React.createElement(
          "table",
          { className: !Data.Params.Memn.Negative ? "table collapse-hidden" : "table"
          },
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "nowrap text-right"
                },
                "Delay ",
                React.createElement(
                  "span",
                  { className: "badge"
                  },
                  Data.Params.Memd
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Dlinks.Memd.Less.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Dlinks.Memd.Less.ExtraClass != null ? Data.Params.Dlinks.Memd.Less.ExtraClass : "")

                  },
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "-"
                  ),
                  " ",
                  Data.Params.Dlinks.Memd.Less.Text
                ),
                React.createElement(
                  "a",
                  { href: Data.Params.Dlinks.Memd.More.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Dlinks.Memd.More.ExtraClass != null ? Data.Params.Dlinks.Memd.More.ExtraClass : "")

                  },
                  Data.Params.Dlinks.Memd.More.Text,
                  " ",
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "+"
                  )
                )
              )
            ),
            React.createElement("td", { className: "col-md-10", colSpan: "2"
            })
          ),
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "nowrap text-right"
                },
                "Rows ",
                React.createElement(
                  "span",
                  { className: "badge"
                  },
                  Data.Params.Memn.Absolute
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Memn.Less.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Memn.Less.ExtraClass != null ? Data.Params.Nlinks.Memn.Less.ExtraClass : "")

                  },
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "-"
                  ),
                  " ",
                  Data.Params.Nlinks.Memn.Less.Text
                ),
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Memn.More.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Memn.More.ExtraClass != null ? Data.Params.Nlinks.Memn.More.ExtraClass : "")

                  },
                  Data.Params.Nlinks.Memn.More.Text,
                  " ",
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "+"
                  )
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Memn.Zero.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Memn.Zero.ExtraClass != null ? Data.Params.Nlinks.Memn.Zero.ExtraClass : "")

                  },
                  Data.Params.Nlinks.Memn.Zero.Text,
                  " ",
                  React.createElement("span", { className: "xlabel xlabel-default"
                  })
                )
              )
            ),
            React.createElement("td", { className: "col-md-10"
            })
          ),
          "  "
        ),
        "  ",
        React.createElement(
          "table",
          { className: Data.Params.Memn.Absolute != 0 ? "table table-hover" : "collapse-hidden"
          },
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              null,
              React.createElement("th", null),
              React.createElement(
                "th",
                { className: "text-right"
                },
                "Free"
              ),
              React.createElement(
                "th",
                { className: "text-right"
                },
                "Used"
              ),
              React.createElement(
                "th",
                { className: "text-right"
                },
                "Total"
              )
            )
          ),
          React.createElement(
            "tbody",
            null,
            rows
          )
        )
      );
    },

    ifbytes_rows: function ifbytes_rows(Data, $if) {
      return React.createElement(
        "tr",
        { key: "ifbytes-rowby-name-" + $if.Name
        },
        React.createElement(
          "td",
          null,
          $if.Name
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $if.DeltaIn
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $if.DeltaOut
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $if.In
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $if.Out
        )
      );
    },
    iferrors_rows: function iferrors_rows(Data, $if) {
      return React.createElement(
        "tr",
        { key: "iferrors-rowby-name-" + $if.Name
        },
        React.createElement(
          "td",
          null,
          $if.Name
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $if.DeltaIn
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $if.DeltaOut
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $if.In
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $if.Out
        )
      );
    },
    ifpackets_rows: function ifpackets_rows(Data, $if) {
      return React.createElement(
        "tr",
        { key: "ifpackets-rowby-name-" + $if.Name
        },
        React.createElement(
          "td",
          null,
          $if.Name
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $if.DeltaIn
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $if.DeltaOut
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $if.In
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $if.Out
        )
      );
    },
    panelif: function panelif(Data, r1, r2, r3) {
      return React.createElement(
        "div",
        { className: !Data.Params.Ifn.Negative ? "" : "panel panel-default"
        },
        "  ",
        React.createElement(
          "div",
          { className: !Data.Params.Ifn.Negative ? "" : "panel-heading"
          },
          "    ",
          React.createElement(
            "a",
            { href: Data.Params.Tlinks.Ifn, onClick: this.handleClick, className: "panel-title btn-block"
            },
            "      ",
            React.createElement(
              "b",
              { className: !Data.Params.Ifn.Negative ? "h4" : "h4 bg-info"
              },
              "Interfaces"
            ),
            "    "
          ),
          "  "
        ),
        "  ",
        React.createElement(
          "table",
          { className: !Data.Params.Ifn.Negative ? "table collapse-hidden" : "table"
          },
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "nowrap text-right"
                },
                "Delay ",
                React.createElement(
                  "span",
                  { className: "badge"
                  },
                  Data.Params.Ifd
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Dlinks.Ifd.Less.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Dlinks.Ifd.Less.ExtraClass != null ? Data.Params.Dlinks.Ifd.Less.ExtraClass : "")

                  },
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "-"
                  ),
                  " ",
                  Data.Params.Dlinks.Ifd.Less.Text
                ),
                React.createElement(
                  "a",
                  { href: Data.Params.Dlinks.Ifd.More.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Dlinks.Ifd.More.ExtraClass != null ? Data.Params.Dlinks.Ifd.More.ExtraClass : "")

                  },
                  Data.Params.Dlinks.Ifd.More.Text,
                  " ",
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "+"
                  )
                )
              )
            ),
            React.createElement("td", { className: "col-md-10", colSpan: "2"
            })
          ),
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "nowrap text-right"
                },
                "Rows ",
                React.createElement(
                  "span",
                  { className: "badge"
                  },
                  Data.Params.Ifn.Absolute
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Ifn.Less.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Ifn.Less.ExtraClass != null ? Data.Params.Nlinks.Ifn.Less.ExtraClass : "")

                  },
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "-"
                  ),
                  " ",
                  Data.Params.Nlinks.Ifn.Less.Text
                ),
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Ifn.More.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Ifn.More.ExtraClass != null ? Data.Params.Nlinks.Ifn.More.ExtraClass : "")

                  },
                  Data.Params.Nlinks.Ifn.More.Text,
                  " ",
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "+"
                  )
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Ifn.Zero.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Ifn.Zero.ExtraClass != null ? Data.Params.Nlinks.Ifn.Zero.ExtraClass : "")

                  },
                  Data.Params.Nlinks.Ifn.Zero.Text,
                  " ",
                  React.createElement("span", { className: "xlabel xlabel-default"
                  })
                )
              )
            ),
            React.createElement("td", { className: "col-md-10"
            })
          ),
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "text-right"
                },
                "Select"
              )
            ),
            React.createElement(
              "td",
              { colSpan: "3"
              },
              React.createElement(
                "ul",
                { className: "nav nav-pills"
                },
                React.createElement(
                  "li",
                  { className: Data.Params.Ifn.Absolute != 0 && Data.Params.Ift.Absolute == 1 ? "active" : ""
                  },
                  React.createElement(
                    "a",
                    { href: Data.Params.Vlinks.Ift[1 - 1].LinkHref, onClick: this.handleClick
                    },
                    "Packets"
                  )
                ),
                React.createElement(
                  "li",
                  { className: Data.Params.Ifn.Absolute != 0 && Data.Params.Ift.Absolute == 2 ? "active" : ""
                  },
                  React.createElement(
                    "a",
                    { href: Data.Params.Vlinks.Ift[2 - 1].LinkHref, onClick: this.handleClick
                    },
                    "Errors"
                  )
                ),
                React.createElement(
                  "li",
                  { className: Data.Params.Ifn.Absolute != 0 && Data.Params.Ift.Absolute == 3 ? "active" : ""
                  },
                  React.createElement(
                    "a",
                    { href: Data.Params.Vlinks.Ift[3 - 1].LinkHref, onClick: this.handleClick
                    },
                    "Bytes"
                  )
                )
              )
            )
          ),
          "  "
        ),
        "  ",
        React.createElement(
          "table",
          { className: Data.Params.Ifn.Absolute != 0 && Data.Params.Ift.Absolute == 1 ? "table table-hover" : "collapse-hidden"
          },
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              null,
              React.createElement(
                "th",
                null,
                "Interface"
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap", title: "per second"
                },
                "In ",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  "ps"
                )
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap", title: "per second"
                },
                "Out ",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  "ps"
                )
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap", title: "total modulo 4G"
                },
                "In ",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  "%4G"
                )
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap", title: "total modulo 4G"
                },
                "Out ",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  "%4G"
                )
              )
            )
          ),
          React.createElement(
            "tbody",
            null,
            r1
          )
        ),
        React.createElement(
          "table",
          { className: Data.Params.Ifn.Absolute != 0 && Data.Params.Ift.Absolute == 2 ? "table table-hover" : "collapse-hidden"
          },
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              null,
              React.createElement(
                "th",
                null,
                "Interface"
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap", title: "per second"
                },
                "In ",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  "ps"
                )
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap", title: "per second"
                },
                "Out ",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  "ps"
                )
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap", title: "modulo 4G"
                },
                "In ",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  "%4G"
                )
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap", title: "modulo 4G"
                },
                "Out ",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  "%4G"
                )
              )
            )
          ),
          React.createElement(
            "tbody",
            null,
            r2
          )
        ),
        React.createElement(
          "table",
          { className: Data.Params.Ifn.Absolute != 0 && Data.Params.Ift.Absolute == 3 ? "table table-hover" : "collapse-hidden"
          },
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              null,
              React.createElement(
                "th",
                null,
                "Interface"
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap", title: "BITS per second"
                },
                "In",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  React.createElement(
                    "i",
                    null,
                    "b"
                  ),
                  "ps"
                )
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap", title: "BITS per second"
                },
                "Out",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  React.createElement(
                    "i",
                    null,
                    "b"
                  ),
                  "ps"
                )
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap", title: "total BYTES modulo 4G"
                },
                "In",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  React.createElement(
                    "i",
                    null,
                    "B"
                  ),
                  "%4G"
                )
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap", title: "total BYTES modulo 4G"
                },
                "Out",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  React.createElement(
                    "i",
                    null,
                    "B"
                  ),
                  "%4G"
                )
              )
            )
          ),
          React.createElement(
            "tbody",
            null,
            r3
          )
        )
      );
    },

    cpu_rows: function cpu_rows(Data, $core) {
      return React.createElement(
        "tr",
        { key: "cpu-rowby-N-" + $core.N
        },
        React.createElement(
          "td",
          { className: "text-right nowrap"
          },
          $core.N
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          React.createElement(
            "span",
            { className: "usepercent-text", "data-usepercent": $core.User
            },
            $core.User
          )
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          React.createElement(
            "span",
            { className: "usepercent-text", "data-usepercent": $core.Sys
            },
            $core.Sys
          )
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          React.createElement(
            "span",
            { className: "usepercent-text", "data-usepercent": $core.Wait
            },
            $core.Wait
          )
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          React.createElement(
            "span",
            { className: "usepercent-text-inverse", "data-usepercent": $core.Idle
            },
            $core.Idle
          )
        )
      );
    },
    panelcpu: function panelcpu(Data, rows) {
      return React.createElement(
        "div",
        { className: !Data.Params.CPUn.Negative ? "" : "panel panel-default"
        },
        "  ",
        React.createElement(
          "div",
          { className: !Data.Params.CPUn.Negative ? "" : "panel-heading"
          },
          "    ",
          React.createElement(
            "a",
            { href: Data.Params.Tlinks.CPUn, onClick: this.handleClick, className: "panel-title btn-block"
            },
            "      ",
            React.createElement(
              "b",
              { className: !Data.Params.CPUn.Negative ? "h4" : "h4 bg-info"
              },
              "CPU"
            ),
            "    "
          ),
          "  "
        ),
        "  ",
        React.createElement(
          "table",
          { className: !Data.Params.CPUn.Negative ? "table collapse-hidden" : "table"
          },
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "nowrap text-right"
                },
                "Delay ",
                React.createElement(
                  "span",
                  { className: "badge"
                  },
                  Data.Params.CPUd
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Dlinks.CPUd.Less.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Dlinks.CPUd.Less.ExtraClass != null ? Data.Params.Dlinks.CPUd.Less.ExtraClass : "")

                  },
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "-"
                  ),
                  " ",
                  Data.Params.Dlinks.CPUd.Less.Text
                ),
                React.createElement(
                  "a",
                  { href: Data.Params.Dlinks.CPUd.More.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Dlinks.CPUd.More.ExtraClass != null ? Data.Params.Dlinks.CPUd.More.ExtraClass : "")

                  },
                  Data.Params.Dlinks.CPUd.More.Text,
                  " ",
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "+"
                  )
                )
              )
            ),
            React.createElement("td", { className: "col-md-10", colSpan: "2"
            })
          ),
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "nowrap text-right"
                },
                "Rows ",
                React.createElement(
                  "span",
                  { className: "badge"
                  },
                  Data.Params.CPUn.Absolute
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.CPUn.Less.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.CPUn.Less.ExtraClass != null ? Data.Params.Nlinks.CPUn.Less.ExtraClass : "")

                  },
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "-"
                  ),
                  " ",
                  Data.Params.Nlinks.CPUn.Less.Text
                ),
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.CPUn.More.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.CPUn.More.ExtraClass != null ? Data.Params.Nlinks.CPUn.More.ExtraClass : "")

                  },
                  Data.Params.Nlinks.CPUn.More.Text,
                  " ",
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "+"
                  )
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.CPUn.Zero.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.CPUn.Zero.ExtraClass != null ? Data.Params.Nlinks.CPUn.Zero.ExtraClass : "")

                  },
                  Data.Params.Nlinks.CPUn.Zero.Text,
                  " ",
                  React.createElement("span", { className: "xlabel xlabel-default"
                  })
                )
              )
            ),
            React.createElement("td", { className: "col-md-10"
            })
          ),
          "  "
        ),
        "  ",
        React.createElement(
          "table",
          { className: Data.Params.CPUn.Absolute != 0 ? "table table-hover" : "collapse-hidden"
          },
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              null,
              React.createElement("th", null),
              React.createElement(
                "th",
                { className: "text-right nowrap"
                },
                "User",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  "%"
                )
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap"
                },
                "Sys",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  "%"
                )
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap"
                },
                "Wait",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  "%"
                )
              ),
              React.createElement(
                "th",
                { className: "text-right nowrap"
                },
                "Idle",
                React.createElement(
                  "span",
                  { className: "unit"
                  },
                  "%"
                )
              )
            )
          ),
          React.createElement(
            "tbody",
            null,
            rows
          )
        )
      );
    },

    dfbytes_rows: function dfbytes_rows(Data, $disk) {
      return React.createElement(
        "tr",
        { key: "dfbytes-rowby-dirname-" + $disk.DirName
        },
        React.createElement(
          "td",
          null,
          $disk.DevName
        ),
        React.createElement(
          "td",
          null,
          $disk.DirName
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $disk.Avail
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          React.createElement(
            "span",
            { className: "label", "data-usepercent": $disk.UsePercent
            },
            $disk.UsePercent,
            "%"
          ),
          " ",
          $disk.Used
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $disk.Total
        )
      );
    },
    dfinodes_rows: function dfinodes_rows(Data, $disk) {
      return React.createElement(
        "tr",
        { key: "dfinodes-rowby-dirname-" + $disk.DirName
        },
        React.createElement(
          "td",
          null,
          $disk.DevName
        ),
        React.createElement(
          "td",
          null,
          $disk.DirName
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $disk.Ifree
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          React.createElement(
            "span",
            { className: "label", "data-usepercent": $disk.IusePercent
            },
            $disk.IusePercent,
            "%"
          ),
          " ",
          $disk.Iused
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          $disk.Inodes
        )
      );
    },
    paneldf: function paneldf(Data, r1, r2) {
      return React.createElement(
        "div",
        { className: !Data.Params.Dfn.Negative ? "" : "panel panel-default"
        },
        "  ",
        React.createElement(
          "div",
          { className: !Data.Params.Dfn.Negative ? "" : "panel-heading"
          },
          "    ",
          React.createElement(
            "a",
            { href: Data.Params.Tlinks.Dfn, onClick: this.handleClick, className: "panel-title btn-block"
            },
            "      ",
            React.createElement(
              "b",
              { className: !Data.Params.Dfn.Negative ? "h4" : "h4 bg-info"
              },
              "Disk usage"
            ),
            "    "
          ),
          "  "
        ),
        "  ",
        React.createElement(
          "table",
          { className: !Data.Params.Dfn.Negative ? "table collapse-hidden" : "table"
          },
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "nowrap text-right"
                },
                "Delay ",
                React.createElement(
                  "span",
                  { className: "badge"
                  },
                  Data.Params.Dfd
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Dlinks.Dfd.Less.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Dlinks.Dfd.Less.ExtraClass != null ? Data.Params.Dlinks.Dfd.Less.ExtraClass : "")

                  },
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "-"
                  ),
                  " ",
                  Data.Params.Dlinks.Dfd.Less.Text
                ),
                React.createElement(
                  "a",
                  { href: Data.Params.Dlinks.Dfd.More.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Dlinks.Dfd.More.ExtraClass != null ? Data.Params.Dlinks.Dfd.More.ExtraClass : "")

                  },
                  Data.Params.Dlinks.Dfd.More.Text,
                  " ",
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "+"
                  )
                )
              )
            ),
            React.createElement("td", { className: "col-md-10", colSpan: "2"
            })
          ),
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "nowrap text-right"
                },
                "Rows ",
                React.createElement(
                  "span",
                  { className: "badge"
                  },
                  Data.Params.Dfn.Absolute
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Dfn.Less.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Dfn.Less.ExtraClass != null ? Data.Params.Nlinks.Dfn.Less.ExtraClass : "")

                  },
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "-"
                  ),
                  " ",
                  Data.Params.Nlinks.Dfn.Less.Text
                ),
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Dfn.More.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Dfn.More.ExtraClass != null ? Data.Params.Nlinks.Dfn.More.ExtraClass : "")

                  },
                  Data.Params.Nlinks.Dfn.More.Text,
                  " ",
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "+"
                  )
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Dfn.Zero.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Dfn.Zero.ExtraClass != null ? Data.Params.Nlinks.Dfn.Zero.ExtraClass : "")

                  },
                  Data.Params.Nlinks.Dfn.Zero.Text,
                  " ",
                  React.createElement("span", { className: "xlabel xlabel-default"
                  })
                )
              )
            ),
            React.createElement("td", { className: "col-md-10"
            })
          ),
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "text-right"
                },
                "Select"
              )
            ),
            React.createElement(
              "td",
              { colSpan: "3"
              },
              React.createElement(
                "ul",
                { className: "nav nav-pills"
                },
                React.createElement(
                  "li",
                  { className: Data.Params.Dfn.Absolute != 0 && Data.Params.Dft.Absolute == 1 ? "active" : ""
                  },
                  React.createElement(
                    "a",
                    { href: Data.Params.Vlinks.Dft[1 - 1].LinkHref, onClick: this.handleClick
                    },
                    "Inodes"
                  )
                ),
                React.createElement(
                  "li",
                  { className: Data.Params.Dfn.Absolute != 0 && Data.Params.Dft.Absolute == 2 ? "active" : ""
                  },
                  React.createElement(
                    "a",
                    { href: Data.Params.Vlinks.Dft[2 - 1].LinkHref, onClick: this.handleClick
                    },
                    "Bytes"
                  )
                )
              )
            )
          ),
          "  "
        ),
        "  ",
        React.createElement(
          "table",
          { className: Data.Params.Dfn.Absolute != 0 && Data.Params.Dft.Absolute == 1 ? "table table-hover" : "collapse-hidden"
          },
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              null,
              React.createElement(
                "th",
                { className: "header"
                },
                "Device"
              ),
              React.createElement(
                "th",
                { className: "header"
                },
                "Mounted"
              ),
              React.createElement(
                "th",
                { className: "header text-right"
                },
                "Avail"
              ),
              React.createElement(
                "th",
                { className: "header text-right"
                },
                "Used"
              ),
              React.createElement(
                "th",
                { className: "header text-right"
                },
                "Total"
              )
            )
          ),
          React.createElement(
            "tbody",
            null,
            r1
          )
        ),
        React.createElement(
          "table",
          { className: Data.Params.Dfn.Absolute != 0 && Data.Params.Dft.Absolute == 2 ? "table table-hover" : "collapse-hidden"
          },
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              null,
              React.createElement(
                "th",
                { className: "header "
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Dfk[1 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[1 - 1].LinkClass
                  },
                  "  Device",
                  React.createElement("span", { className: Data.Params.Vlinks.Dfk[1 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header "
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Dfk[2 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[2 - 1].LinkClass
                  },
                  "  Mounted",
                  React.createElement("span", { className: Data.Params.Vlinks.Dfk[2 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header text-right"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Dfk[3 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[3 - 1].LinkClass
                  },
                  "  Avail",
                  React.createElement("span", { className: Data.Params.Vlinks.Dfk[3 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header text-right"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Dfk[4 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[4 - 1].LinkClass
                  },
                  "  Used",
                  React.createElement("span", { className: Data.Params.Vlinks.Dfk[4 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header text-right"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Dfk[5 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[5 - 1].LinkClass
                  },
                  "  Total",
                  React.createElement("span", { className: Data.Params.Vlinks.Dfk[5 - 1].CaretClass
                  })
                )
              )
            )
          ),
          React.createElement(
            "tbody",
            null,
            r2
          )
        )
      );
    },

    ps_rows: function ps_rows(Data, $proc) {
      return React.createElement(
        "tr",
        { key: "ps-rowby-pid-" + $proc.PID
        },
        React.createElement(
          "td",
          { className: "text-right"
          },
          " ",
          $proc.PID
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          " ",
          $proc.UID
        ),
        React.createElement(
          "td",
          null,
          $proc.User
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          " ",
          $proc.Priority
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          " ",
          $proc.Nice
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          " ",
          $proc.Size
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          " ",
          $proc.Resident
        ),
        React.createElement(
          "td",
          { className: "text-center"
          },
          $proc.Time
        ),
        React.createElement(
          "td",
          null,
          $proc.Name
        )
      );
    },
    panelps: function panelps(Data, rows) {
      return React.createElement(
        "div",
        { className: !Data.Params.Psn.Negative ? "" : "panel panel-default"
        },
        "  ",
        React.createElement(
          "div",
          { className: !Data.Params.Psn.Negative ? "" : "panel-heading"
          },
          "    ",
          React.createElement(
            "a",
            { href: Data.Params.Tlinks.Psn, onClick: this.handleClick, className: "panel-title btn-block"
            },
            "      ",
            React.createElement(
              "b",
              { className: !Data.Params.Psn.Negative ? "h4" : "h4 bg-info"
              },
              "Processes"
            ),
            "    "
          ),
          "  "
        ),
        "  ",
        React.createElement(
          "table",
          { className: !Data.Params.Psn.Negative ? "table collapse-hidden" : "table"
          },
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "nowrap text-right"
                },
                "Delay ",
                React.createElement(
                  "span",
                  { className: "badge"
                  },
                  Data.Params.Psd
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Dlinks.Psd.Less.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Dlinks.Psd.Less.ExtraClass != null ? Data.Params.Dlinks.Psd.Less.ExtraClass : "")

                  },
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "-"
                  ),
                  " ",
                  Data.Params.Dlinks.Psd.Less.Text
                ),
                React.createElement(
                  "a",
                  { href: Data.Params.Dlinks.Psd.More.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Dlinks.Psd.More.ExtraClass != null ? Data.Params.Dlinks.Psd.More.ExtraClass : "")

                  },
                  Data.Params.Dlinks.Psd.More.Text,
                  " ",
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "+"
                  )
                )
              )
            ),
            React.createElement("td", { className: "col-md-10", colSpan: "2"
            })
          ),
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "nowrap text-right"
                },
                "Rows ",
                React.createElement(
                  "span",
                  { className: "badge"
                  },
                  Data.Params.Psn.Absolute
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Psn.Less.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Psn.Less.ExtraClass != null ? Data.Params.Nlinks.Psn.Less.ExtraClass : "")

                  },
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "-"
                  ),
                  " ",
                  Data.Params.Nlinks.Psn.Less.Text
                ),
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Psn.More.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Psn.More.ExtraClass != null ? Data.Params.Nlinks.Psn.More.ExtraClass : "")

                  },
                  Data.Params.Nlinks.Psn.More.Text,
                  " ",
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "+"
                  )
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Psn.Zero.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Psn.Zero.ExtraClass != null ? Data.Params.Nlinks.Psn.Zero.ExtraClass : "")

                  },
                  Data.Params.Nlinks.Psn.Zero.Text,
                  " ",
                  React.createElement("span", { className: "xlabel xlabel-default"
                  })
                )
              )
            ),
            React.createElement("td", { className: "col-md-10"
            })
          ),
          "  "
        ),
        "  ",
        React.createElement(
          "table",
          { className: Data.Params.Psn.Absolute != 0 ? "table table-hover" : "collapse-hidden"
          },
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              null,
              React.createElement(
                "th",
                { className: "header text-right"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Psk[1 - 1].LinkHref, className: Data.Params.Vlinks.Psk[1 - 1].LinkClass
                  },
                  "  PID",
                  React.createElement("span", { className: Data.Params.Vlinks.Psk[1 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header text-right"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Psk[2 - 1].LinkHref, className: Data.Params.Vlinks.Psk[2 - 1].LinkClass
                  },
                  "  UID",
                  React.createElement("span", { className: Data.Params.Vlinks.Psk[2 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header "
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Psk[3 - 1].LinkHref, className: Data.Params.Vlinks.Psk[3 - 1].LinkClass
                  },
                  "  USER",
                  React.createElement("span", { className: Data.Params.Vlinks.Psk[3 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header text-right"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Psk[4 - 1].LinkHref, className: Data.Params.Vlinks.Psk[4 - 1].LinkClass
                  },
                  "  PR",
                  React.createElement("span", { className: Data.Params.Vlinks.Psk[4 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header text-right"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Psk[5 - 1].LinkHref, className: Data.Params.Vlinks.Psk[5 - 1].LinkClass
                  },
                  "  NI",
                  React.createElement("span", { className: Data.Params.Vlinks.Psk[5 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header text-right"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Psk[6 - 1].LinkHref, className: Data.Params.Vlinks.Psk[6 - 1].LinkClass
                  },
                  "  VIRT",
                  React.createElement("span", { className: Data.Params.Vlinks.Psk[6 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header text-right"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Psk[7 - 1].LinkHref, className: Data.Params.Vlinks.Psk[7 - 1].LinkClass
                  },
                  "  RES",
                  React.createElement("span", { className: Data.Params.Vlinks.Psk[7 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header text-center"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Psk[8 - 1].LinkHref, className: Data.Params.Vlinks.Psk[8 - 1].LinkClass
                  },
                  "  TIME",
                  React.createElement("span", { className: Data.Params.Vlinks.Psk[8 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header "
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Psk[9 - 1].LinkHref, className: Data.Params.Vlinks.Psk[9 - 1].LinkClass
                  },
                  "  COMMAND",
                  React.createElement("span", { className: Data.Params.Vlinks.Psk[9 - 1].CaretClass
                  })
                )
              )
            )
          ),
          React.createElement(
            "tbody",
            null,
            rows
          )
        )
      );
    },

    vagrant_rows: function vagrant_rows(Data, $mach) {
      return React.createElement(
        "tr",
        { key: "vagrant-rowby-uuid-" + $mach.UUID
        },
        React.createElement(
          "td",
          null,
          $mach.UUID
        ),
        React.createElement(
          "td",
          null,
          $mach.Name
        ),
        React.createElement(
          "td",
          null,
          $mach.Provider
        ),
        React.createElement(
          "td",
          null,
          $mach.State
        ),
        React.createElement(
          "td",
          null,
          $mach.Vagrantfile_path
        )
      );
    },
    vagrant_error: function vagrant_error(Data) {
      return React.createElement(
        "tr",
        { key: "vgerror"
        },
        React.createElement(
          "td",
          { colSpan: "5"
          },
          Data.VagrantError
        )
      );
    },
    panelvg: function panelvg(Data, rows) {
      return React.createElement(
        "div",
        { className: !Data.Params.Vgn.Negative ? "" : "panel panel-default"
        },
        "  ",
        React.createElement(
          "div",
          { className: !Data.Params.Vgn.Negative ? "" : "panel-heading"
          },
          "    ",
          React.createElement(
            "a",
            { href: Data.Params.Tlinks.Vgn, onClick: this.handleClick, className: "panel-title btn-block"
            },
            "      ",
            React.createElement(
              "b",
              { className: !Data.Params.Vgn.Negative ? "h4" : "h4 bg-info"
              },
              "Vagrant"
            ),
            "    "
          ),
          "  "
        ),
        "  ",
        React.createElement(
          "table",
          { className: !Data.Params.Vgn.Negative ? "table collapse-hidden" : "table"
          },
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "nowrap text-right"
                },
                "Delay ",
                React.createElement(
                  "span",
                  { className: "badge"
                  },
                  Data.Params.Vgd
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Dlinks.Vgd.Less.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Dlinks.Vgd.Less.ExtraClass != null ? Data.Params.Dlinks.Vgd.Less.ExtraClass : "")

                  },
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "-"
                  ),
                  " ",
                  Data.Params.Dlinks.Vgd.Less.Text
                ),
                React.createElement(
                  "a",
                  { href: Data.Params.Dlinks.Vgd.More.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Dlinks.Vgd.More.ExtraClass != null ? Data.Params.Dlinks.Vgd.More.ExtraClass : "")

                  },
                  Data.Params.Dlinks.Vgd.More.Text,
                  " ",
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "+"
                  )
                )
              )
            ),
            React.createElement("td", { className: "col-md-10", colSpan: "2"
            })
          ),
          React.createElement(
            "tr",
            { className: "panel-config"
            },
            React.createElement(
              "td",
              { className: "col-md-2"
              },
              React.createElement(
                "div",
                { className: "nowrap text-right"
                },
                "Rows ",
                React.createElement(
                  "span",
                  { className: "badge"
                  },
                  Data.Params.Vgn.Absolute
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Vgn.Less.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Vgn.Less.ExtraClass != null ? Data.Params.Nlinks.Vgn.Less.ExtraClass : "")

                  },
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "-"
                  ),
                  " ",
                  Data.Params.Nlinks.Vgn.Less.Text
                ),
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Vgn.More.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Vgn.More.ExtraClass != null ? Data.Params.Nlinks.Vgn.More.ExtraClass : "")

                  },
                  Data.Params.Nlinks.Vgn.More.Text,
                  " ",
                  React.createElement(
                    "span",
                    { className: "xlabel xlabel-default"
                    },
                    "+"
                  )
                )
              )
            ),
            React.createElement(
              "td",
              null,
              React.createElement(
                "div",
                { className: "btn-group nowrap-group", role: "group"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Nlinks.Vgn.Zero.Href, onClick: this.handleClick, className: "btn btn-default" + " " + (Data.Params.Nlinks.Vgn.Zero.ExtraClass != null ? Data.Params.Nlinks.Vgn.Zero.ExtraClass : "")

                  },
                  Data.Params.Nlinks.Vgn.Zero.Text,
                  " ",
                  React.createElement("span", { className: "xlabel xlabel-default"
                  })
                )
              )
            ),
            React.createElement("td", { className: "col-md-10"
            })
          ),
          "  "
        ),
        "  ",
        React.createElement(
          "table",
          { className: Data.Params.Vgn.Absolute != 0 ? "table table-hover" : "collapse-hidden"
          },
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              null,
              React.createElement(
                "th",
                null,
                "ID"
              ),
              React.createElement(
                "th",
                null,
                "Name"
              ),
              React.createElement(
                "th",
                null,
                "Provider"
              ),
              React.createElement(
                "th",
                null,
                "State"
              ),
              React.createElement(
                "th",
                null,
                "Directory"
              )
            )
          ),
          React.createElement(
            "tbody",
            null,
            rows
          )
        )
      );
    }
  };
});
