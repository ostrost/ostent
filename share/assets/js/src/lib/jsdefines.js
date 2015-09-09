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
          { className: "text-right bg-usepct", "data-usepct": $mem.UsePct
          },
          $mem.UsePct,
          "%"
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
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
        React.createElement(
          "div",
          { className: "h4 padding-left-like-panel-heading"
          },
          React.createElement(
            "a",
            { href: Data.Params.Tlinks.Memn, onClick: this.handleClick
            },
            "Memory"
          )
        ),
        React.createElement(
          "ul",
          { className: !Data.Params.Memn.Negative ? "collapse-hidden" : "list-group"
          },
          React.createElement(
            "li",
            { className: "list-group-item text-nowrap th"
            },
            React.createElement(
              "ul",
              { className: "list-inline"
              },
              React.createElement(
                "li",
                null,
                React.createElement(
                  "span",
                  null,
                  React.createElement(
                    "b",
                    null,
                    "Delay"
                  ),
                  " ",
                  React.createElement(
                    "span",
                    { className: "badge"
                    },
                    Data.Params.Memd
                  )
                ),
                " ",
                React.createElement(
                  "div",
                  { className: "btn-group"
                  },
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Dlinks.Memd.Less.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Memd.Less.ExtraClass != null ? Data.Params.Dlinks.Memd.Less.ExtraClass : "")
                    },
                    "- ",
                    Data.Params.Dlinks.Memd.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Dlinks.Memd.More.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Memd.More.ExtraClass != null ? Data.Params.Dlinks.Memd.More.ExtraClass : "")
                    },
                    Data.Params.Dlinks.Memd.More.Text,
                    " +"
                  )
                )
              ),
              React.createElement(
                "li",
                null,
                React.createElement(
                  "span",
                  null,
                  React.createElement(
                    "b",
                    null,
                    "Rows"
                  ),
                  " ",
                  React.createElement(
                    "span",
                    { className: "badge"
                    },
                    Data.Params.Memn.Absolute
                  )
                ),
                " ",
                React.createElement(
                  "div",
                  { className: "btn-group"
                  },
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Nlinks.Memn.Less.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Memn.Less.ExtraClass != null ? Data.Params.Nlinks.Memn.Less.ExtraClass : "")
                    },
                    "- ",
                    Data.Params.Nlinks.Memn.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Nlinks.Memn.More.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Memn.More.ExtraClass != null ? Data.Params.Nlinks.Memn.More.ExtraClass : "")
                    },
                    Data.Params.Nlinks.Memn.More.Text,
                    " +"
                  )
                )
              )
            )
          )
        ),
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
                "Use%"
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

    if_rows: function if_rows(Data, $if) {
      return React.createElement(
        "tr",
        { key: "if-rowby-name-" + $if.Name
        },
        React.createElement(
          "td",
          { className: "text-nowrap clip12", title: $if.Name
          },
          $if.Name
        ),
        React.createElement(
          "td",
          { className: "text-right text-nowrap"
          },
          React.createElement(
            "span",
            { className: "text-graylighter", title: "Total BYTES modulo 4G"
            },
            $if.BytesIn,
            "/",
            $if.BytesOut
          ),
          " ",
          $if.DeltaBitsIn,
          "/",
          $if.DeltaBitsOut
        ),
        React.createElement(
          "td",
          { className: "text-right text-nowrap"
          },
          React.createElement(
            "span",
            { className: "text-graylighter", title: "Total packets modulo 4G"
            },
            $if.PacketsIn,
            "/",
            $if.PacketsOut
          ),
          " ",
          $if.DeltaPacketsIn,
          "/",
          $if.DeltaPacketsOut
        ),
        React.createElement(
          "td",
          { className: "text-right text-nowrap"
          },
          React.createElement(
            "span",
            { className: "text-graylighter", title: "Total errors modulo 4G"
            },
            $if.ErrorsIn,
            "/",
            $if.ErrorsOut
          ),
          " ",
          $if.DeltaErrorsIn,
          "/",
          $if.DeltaErrorsOut
        )
      );
    },
    panelif: function panelif(Data, rows) {
      return React.createElement(
        "div",
        { className: !Data.Params.Ifn.Negative ? "" : "panel panel-default"
        },
        React.createElement(
          "div",
          { className: "h4 padding-left-like-panel-heading"
          },
          React.createElement(
            "a",
            { href: Data.Params.Tlinks.Ifn, onClick: this.handleClick
            },
            "Interfaces"
          )
        ),
        React.createElement(
          "ul",
          { className: !Data.Params.Ifn.Negative ? "collapse-hidden" : "list-group"
          },
          React.createElement(
            "li",
            { className: "list-group-item text-nowrap th"
            },
            React.createElement(
              "ul",
              { className: "list-inline"
              },
              React.createElement(
                "li",
                null,
                React.createElement(
                  "span",
                  null,
                  React.createElement(
                    "b",
                    null,
                    "Delay"
                  ),
                  " ",
                  React.createElement(
                    "span",
                    { className: "badge"
                    },
                    Data.Params.Ifd
                  )
                ),
                " ",
                React.createElement(
                  "div",
                  { className: "btn-group"
                  },
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Dlinks.Ifd.Less.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Ifd.Less.ExtraClass != null ? Data.Params.Dlinks.Ifd.Less.ExtraClass : "")
                    },
                    "- ",
                    Data.Params.Dlinks.Ifd.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Dlinks.Ifd.More.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Ifd.More.ExtraClass != null ? Data.Params.Dlinks.Ifd.More.ExtraClass : "")
                    },
                    Data.Params.Dlinks.Ifd.More.Text,
                    " +"
                  )
                )
              ),
              React.createElement(
                "li",
                null,
                React.createElement(
                  "span",
                  null,
                  React.createElement(
                    "b",
                    null,
                    "Rows"
                  ),
                  " ",
                  React.createElement(
                    "span",
                    { className: "badge"
                    },
                    Data.Params.Ifn.Absolute
                  )
                ),
                " ",
                React.createElement(
                  "div",
                  { className: "btn-group"
                  },
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Nlinks.Ifn.Less.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Ifn.Less.ExtraClass != null ? Data.Params.Nlinks.Ifn.Less.ExtraClass : "")
                    },
                    "- ",
                    Data.Params.Nlinks.Ifn.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Nlinks.Ifn.More.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Ifn.More.ExtraClass != null ? Data.Params.Nlinks.Ifn.More.ExtraClass : "")
                    },
                    Data.Params.Nlinks.Ifn.More.Text,
                    " +"
                  )
                )
              )
            )
          )
        ),
        React.createElement(
          "table",
          { className: Data.Params.Ifn.Absolute != 0 ? "table table-hover" : "collapse-hidden"
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
                { className: "text-right col-md-3", title: "Bits per second"
                },
                "Bits In/Out"
              ),
              React.createElement(
                "th",
                { className: "text-right col-md-3", title: "Packets per second"
                },
                "Packets In/Out"
              ),
              React.createElement(
                "th",
                { className: "text-right col-md-3", title: "Errors per second"
                },
                "Errors In/Out"
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

    cpu_rows: function cpu_rows(Data, $cpu) {
      return React.createElement(
        "tr",
        { key: "cpu-rowby-N-" + $cpu.N
        },
        React.createElement(
          "td",
          { className: "text-right text-nowrap"
          },
          $cpu.N
        ),
        React.createElement(
          "td",
          { className: "text-right bg-usepct", "data-usepct": $cpu.UserPct
          },
          $cpu.UserPct,
          "%"
        ),
        React.createElement(
          "td",
          { className: "text-right bg-usepct", "data-usepct": $cpu.SysPct
          },
          $cpu.SysPct,
          "%"
        ),
        React.createElement(
          "td",
          { className: "text-right bg-usepct", "data-usepct": $cpu.WaitPct
          },
          $cpu.WaitPct,
          "%"
        ),
        React.createElement(
          "td",
          { className: "text-right bg-usepct-inverse", "data-usepct": $cpu.IdlePct
          },
          $cpu.IdlePct,
          "%"
        )
      );
    },
    panelcpu: function panelcpu(Data, rows) {
      return React.createElement(
        "div",
        { className: !Data.Params.CPUn.Negative ? "" : "panel panel-default"
        },
        React.createElement(
          "div",
          { className: "h4 padding-left-like-panel-heading"
          },
          React.createElement(
            "a",
            { href: Data.Params.Tlinks.CPUn, onClick: this.handleClick
            },
            "CPU"
          )
        ),
        React.createElement(
          "ul",
          { className: !Data.Params.CPUn.Negative ? "collapse-hidden" : "list-group"
          },
          React.createElement(
            "li",
            { className: "list-group-item text-nowrap th"
            },
            React.createElement(
              "ul",
              { className: "list-inline"
              },
              React.createElement(
                "li",
                null,
                React.createElement(
                  "span",
                  null,
                  React.createElement(
                    "b",
                    null,
                    "Delay"
                  ),
                  " ",
                  React.createElement(
                    "span",
                    { className: "badge"
                    },
                    Data.Params.CPUd
                  )
                ),
                " ",
                React.createElement(
                  "div",
                  { className: "btn-group"
                  },
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Dlinks.CPUd.Less.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.CPUd.Less.ExtraClass != null ? Data.Params.Dlinks.CPUd.Less.ExtraClass : "")
                    },
                    "- ",
                    Data.Params.Dlinks.CPUd.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Dlinks.CPUd.More.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.CPUd.More.ExtraClass != null ? Data.Params.Dlinks.CPUd.More.ExtraClass : "")
                    },
                    Data.Params.Dlinks.CPUd.More.Text,
                    " +"
                  )
                )
              ),
              React.createElement(
                "li",
                null,
                React.createElement(
                  "span",
                  null,
                  React.createElement(
                    "b",
                    null,
                    "Rows"
                  ),
                  " ",
                  React.createElement(
                    "span",
                    { className: "badge"
                    },
                    Data.Params.CPUn.Absolute
                  )
                ),
                " ",
                React.createElement(
                  "div",
                  { className: "btn-group"
                  },
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Nlinks.CPUn.Less.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.CPUn.Less.ExtraClass != null ? Data.Params.Nlinks.CPUn.Less.ExtraClass : "")
                    },
                    "- ",
                    Data.Params.Nlinks.CPUn.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Nlinks.CPUn.More.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.CPUn.More.ExtraClass != null ? Data.Params.Nlinks.CPUn.More.ExtraClass : "")
                    },
                    Data.Params.Nlinks.CPUn.More.Text,
                    " +"
                  )
                )
              )
            )
          )
        ),
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
                { className: "text-right"
                },
                "User"
              ),
              React.createElement(
                "th",
                { className: "text-right"
                },
                "Sys"
              ),
              React.createElement(
                "th",
                { className: "text-right"
                },
                "Wait"
              ),
              React.createElement(
                "th",
                { className: "text-right"
                },
                "Idle"
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

    df_rows: function df_rows(Data, $df) {
      return React.createElement(
        "tr",
        { key: "df-rowby-dirname-" + $df.DirName
        },
        "  ",
        React.createElement(
          "td",
          { className: "text-nowrap clip12", title: $df.DevName
          },
          $df.DevName
        ),
        "  ",
        React.createElement(
          "td",
          { className: "text-nowrap clip12", title: $df.DirName
          },
          $df.DirName
        ),
        React.createElement(
          "td",
          { className: "text-right text-nowrap"
          },
          React.createElement(
            "span",
            { className: "text-graylighter", title: "Inodes free"
            },
            $df.Ifree
          ),
          " ",
          $df.Avail
        ),
        React.createElement(
          "td",
          { className: "text-right text-nowrap bg-usepct", "data-usepct": $df.UsePct
          },
          React.createElement(
            "span",
            { className: "text-graylighter", title: "Inodes use%"
            },
            $df.IusePct,
            "%"
          ),
          " ",
          $df.UsePct,
          "%"
        ),
        React.createElement(
          "td",
          { className: "text-right text-nowrap"
          },
          React.createElement(
            "span",
            { className: "text-graylighter", title: "Inodes used"
            },
            $df.Iused
          ),
          " ",
          $df.Used
        ),
        React.createElement(
          "td",
          { className: "text-right text-nowrap"
          },
          React.createElement(
            "span",
            { className: "text-graylighter", title: "Inodes total"
            },
            $df.Inodes
          ),
          " ",
          $df.Total
        )
      );
    },
    paneldf: function paneldf(Data, rows) {
      return React.createElement(
        "div",
        { className: !Data.Params.Dfn.Negative ? "" : "panel panel-default"
        },
        React.createElement(
          "div",
          { className: "h4 padding-left-like-panel-heading"
          },
          React.createElement(
            "a",
            { href: Data.Params.Tlinks.Dfn, onClick: this.handleClick
            },
            "Disk usage"
          )
        ),
        React.createElement(
          "ul",
          { className: !Data.Params.Dfn.Negative ? "collapse-hidden" : "list-group"
          },
          React.createElement(
            "li",
            { className: "list-group-item text-nowrap th"
            },
            React.createElement(
              "ul",
              { className: "list-inline"
              },
              React.createElement(
                "li",
                null,
                React.createElement(
                  "span",
                  null,
                  React.createElement(
                    "b",
                    null,
                    "Delay"
                  ),
                  " ",
                  React.createElement(
                    "span",
                    { className: "badge"
                    },
                    Data.Params.Dfd
                  )
                ),
                " ",
                React.createElement(
                  "div",
                  { className: "btn-group"
                  },
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Dlinks.Dfd.Less.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Dfd.Less.ExtraClass != null ? Data.Params.Dlinks.Dfd.Less.ExtraClass : "")
                    },
                    "- ",
                    Data.Params.Dlinks.Dfd.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Dlinks.Dfd.More.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Dfd.More.ExtraClass != null ? Data.Params.Dlinks.Dfd.More.ExtraClass : "")
                    },
                    Data.Params.Dlinks.Dfd.More.Text,
                    " +"
                  )
                )
              ),
              React.createElement(
                "li",
                null,
                React.createElement(
                  "span",
                  null,
                  React.createElement(
                    "b",
                    null,
                    "Rows"
                  ),
                  " ",
                  React.createElement(
                    "span",
                    { className: "badge"
                    },
                    Data.Params.Dfn.Absolute
                  )
                ),
                " ",
                React.createElement(
                  "div",
                  { className: "btn-group"
                  },
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Nlinks.Dfn.Less.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Dfn.Less.ExtraClass != null ? Data.Params.Nlinks.Dfn.Less.ExtraClass : "")
                    },
                    "- ",
                    Data.Params.Nlinks.Dfn.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Nlinks.Dfn.More.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Dfn.More.ExtraClass != null ? Data.Params.Nlinks.Dfn.More.ExtraClass : "")
                    },
                    Data.Params.Nlinks.Dfn.More.Text,
                    " +"
                  )
                )
              )
            )
          )
        ),
        React.createElement(
          "table",
          { className: Data.Params.Dfn.Absolute != 0 ? "table table-hover" : "collapse-hidden"
          },
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              { className: "text-nowrap"
              },
              React.createElement(
                "th",
                { className: "header "
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Dfk[1 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[1 - 1].LinkClass
                  },
                  "Device",
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
                  "Mounted",
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
                  "Avail",
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
                  "Use%",
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
                  "Used",
                  React.createElement("span", { className: Data.Params.Vlinks.Dfk[5 - 1].CaretClass
                  })
                )
              ),
              React.createElement(
                "th",
                { className: "header text-right"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Dfk[6 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[6 - 1].LinkClass
                  },
                  "Total",
                  React.createElement("span", { className: Data.Params.Vlinks.Dfk[6 - 1].CaretClass
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

    ps_rows: function ps_rows(Data, $ps) {
      return React.createElement(
        "tr",
        { key: "ps-rowby-pid-" + $ps.PID
        },
        React.createElement(
          "td",
          { className: "text-right"
          },
          " ",
          $ps.PID
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          " ",
          $ps.UID
        ),
        React.createElement(
          "td",
          null,
          $ps.User
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          " ",
          $ps.Priority
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          " ",
          $ps.Nice
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          " ",
          $ps.Size
        ),
        React.createElement(
          "td",
          { className: "text-right"
          },
          " ",
          $ps.Resident
        ),
        React.createElement(
          "td",
          { className: "text-center"
          },
          $ps.Time
        ),
        React.createElement(
          "td",
          null,
          $ps.Name
        )
      );
    },
    panelps: function panelps(Data, rows) {
      return React.createElement(
        "div",
        { className: !Data.Params.Psn.Negative ? "" : "panel panel-default"
        },
        React.createElement(
          "div",
          { className: "h4 padding-left-like-panel-heading"
          },
          React.createElement(
            "a",
            { href: Data.Params.Tlinks.Psn, onClick: this.handleClick
            },
            "Processes"
          )
        ),
        React.createElement(
          "ul",
          { className: !Data.Params.Psn.Negative ? "collapse-hidden" : "list-group"
          },
          React.createElement(
            "li",
            { className: "list-group-item text-nowrap th"
            },
            React.createElement(
              "ul",
              { className: "list-inline"
              },
              React.createElement(
                "li",
                null,
                React.createElement(
                  "span",
                  null,
                  React.createElement(
                    "b",
                    null,
                    "Delay"
                  ),
                  " ",
                  React.createElement(
                    "span",
                    { className: "badge"
                    },
                    Data.Params.Psd
                  )
                ),
                " ",
                React.createElement(
                  "div",
                  { className: "btn-group"
                  },
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Dlinks.Psd.Less.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Psd.Less.ExtraClass != null ? Data.Params.Dlinks.Psd.Less.ExtraClass : "")
                    },
                    "- ",
                    Data.Params.Dlinks.Psd.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Dlinks.Psd.More.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Psd.More.ExtraClass != null ? Data.Params.Dlinks.Psd.More.ExtraClass : "")
                    },
                    Data.Params.Dlinks.Psd.More.Text,
                    " +"
                  )
                )
              ),
              React.createElement(
                "li",
                null,
                React.createElement(
                  "span",
                  null,
                  React.createElement(
                    "b",
                    null,
                    "Rows"
                  ),
                  " ",
                  React.createElement(
                    "span",
                    { className: "badge"
                    },
                    Data.Params.Psn.Absolute
                  )
                ),
                " ",
                React.createElement(
                  "div",
                  { className: "btn-group"
                  },
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Nlinks.Psn.Less.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Psn.Less.ExtraClass != null ? Data.Params.Nlinks.Psn.Less.ExtraClass : "")
                    },
                    "- ",
                    Data.Params.Nlinks.Psn.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Nlinks.Psn.More.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Psn.More.ExtraClass != null ? Data.Params.Nlinks.Psn.More.ExtraClass : "")
                    },
                    Data.Params.Nlinks.Psn.More.Text,
                    " +"
                  )
                )
              )
            )
          )
        ),
        React.createElement(
          "table",
          { className: Data.Params.Psn.Absolute != 0 ? "table table-hover" : "collapse-hidden"
          },
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              { className: "text-nowrap"
              },
              React.createElement(
                "th",
                { className: "header text-right"
                },
                React.createElement(
                  "a",
                  { href: Data.Params.Vlinks.Psk[1 - 1].LinkHref, className: Data.Params.Vlinks.Psk[1 - 1].LinkClass
                  },
                  "PID",
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
                  "UID",
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
                  "USER",
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
                  "PR",
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
                  "NI",
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
                  "VIRT",
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
                  "RES",
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
                  "TIME",
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
                  "COMMAND",
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

    vg_rows: function vg_rows(Data, $vgm) {
      return React.createElement(
        "tr",
        { key: "vagrant-rowby-uuid-" + $vgm.UUID
        },
        React.createElement(
          "td",
          null,
          $vgm.UUID
        ),
        React.createElement(
          "td",
          null,
          $vgm.Name
        ),
        React.createElement(
          "td",
          null,
          $vgm.Provider
        ),
        React.createElement(
          "td",
          null,
          $vgm.State
        ),
        React.createElement(
          "td",
          null,
          $vgm.Vagrantfile_path
        )
      );
    },
    vg_error: function vg_error(Data) {
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
        React.createElement(
          "div",
          { className: "h4 padding-left-like-panel-heading"
          },
          React.createElement(
            "a",
            { href: Data.Params.Tlinks.Vgn, onClick: this.handleClick
            },
            "Vagrant"
          )
        ),
        React.createElement(
          "ul",
          { className: !Data.Params.Vgn.Negative ? "collapse-hidden" : "list-group"
          },
          React.createElement(
            "li",
            { className: "list-group-item text-nowrap th"
            },
            React.createElement(
              "ul",
              { className: "list-inline"
              },
              React.createElement(
                "li",
                null,
                React.createElement(
                  "span",
                  null,
                  React.createElement(
                    "b",
                    null,
                    "Delay"
                  ),
                  " ",
                  React.createElement(
                    "span",
                    { className: "badge"
                    },
                    Data.Params.Vgd
                  )
                ),
                " ",
                React.createElement(
                  "div",
                  { className: "btn-group"
                  },
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Dlinks.Vgd.Less.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Vgd.Less.ExtraClass != null ? Data.Params.Dlinks.Vgd.Less.ExtraClass : "")
                    },
                    "- ",
                    Data.Params.Dlinks.Vgd.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Dlinks.Vgd.More.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Vgd.More.ExtraClass != null ? Data.Params.Dlinks.Vgd.More.ExtraClass : "")
                    },
                    Data.Params.Dlinks.Vgd.More.Text,
                    " +"
                  )
                )
              ),
              React.createElement(
                "li",
                null,
                React.createElement(
                  "span",
                  null,
                  React.createElement(
                    "b",
                    null,
                    "Rows"
                  ),
                  " ",
                  React.createElement(
                    "span",
                    { className: "badge"
                    },
                    Data.Params.Vgn.Absolute
                  )
                ),
                " ",
                React.createElement(
                  "div",
                  { className: "btn-group"
                  },
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Nlinks.Vgn.Less.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Vgn.Less.ExtraClass != null ? Data.Params.Nlinks.Vgn.Less.ExtraClass : "")
                    },
                    "- ",
                    Data.Params.Nlinks.Vgn.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { onClick: this.handleClick, href: Data.Params.Nlinks.Vgn.More.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Vgn.More.ExtraClass != null ? Data.Params.Nlinks.Vgn.More.ExtraClass : "")
                    },
                    Data.Params.Nlinks.Vgn.More.Text,
                    " +"
                  )
                )
              )
            )
          )
        ),
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
