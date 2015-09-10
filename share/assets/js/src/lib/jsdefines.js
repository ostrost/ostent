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
                    { href: Data.Params.Dlinks.Memd.Less.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Memd.Less.ExtraClass != null ? Data.Params.Dlinks.Memd.Less.ExtraClass : ""), onClick: this.handleClick
                    },
                    "- ",
                    Data.Params.Dlinks.Memd.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { href: Data.Params.Dlinks.Memd.More.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Memd.More.ExtraClass != null ? Data.Params.Dlinks.Memd.More.ExtraClass : ""), onClick: this.handleClick
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
                    { href: Data.Params.Nlinks.Memn.Less.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Memn.Less.ExtraClass != null ? Data.Params.Nlinks.Memn.Less.ExtraClass : ""), onClick: this.handleClick
                    },
                    "- ",
                    Data.Params.Nlinks.Memn.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { href: Data.Params.Nlinks.Memn.More.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Memn.More.ExtraClass != null ? Data.Params.Nlinks.Memn.More.ExtraClass : ""), onClick: this.handleClick
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
            { className: "mutext", title: "Inodes free"
            },
            $df.Ifree
          ),
          " ",
          $df.Avail
        ),
        React.createElement(
          "td",
          { className: "text-right bg-usepct text-nowrap", "data-usepct": $df.UsePct
          },
          React.createElement(
            "span",
            { className: "mutext", title: "Inodes use%"
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
            { className: "mutext", title: "Inodes used"
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
            { className: "mutext", title: "Inodes total"
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
                    { href: Data.Params.Dlinks.Dfd.Less.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Dfd.Less.ExtraClass != null ? Data.Params.Dlinks.Dfd.Less.ExtraClass : ""), onClick: this.handleClick
                    },
                    "- ",
                    Data.Params.Dlinks.Dfd.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { href: Data.Params.Dlinks.Dfd.More.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Dfd.More.ExtraClass != null ? Data.Params.Dlinks.Dfd.More.ExtraClass : ""), onClick: this.handleClick
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
                    { href: Data.Params.Nlinks.Dfn.Less.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Dfn.Less.ExtraClass != null ? Data.Params.Nlinks.Dfn.Less.ExtraClass : ""), onClick: this.handleClick
                    },
                    "- ",
                    Data.Params.Nlinks.Dfn.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { href: Data.Params.Nlinks.Dfn.More.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Dfn.More.ExtraClass != null ? Data.Params.Nlinks.Dfn.More.ExtraClass : ""), onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Dfk[1 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[1 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Dfk[2 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[2 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Dfk[3 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[3 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Dfk[4 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[4 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Dfk[5 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[5 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Dfk[6 - 1].LinkHref, className: Data.Params.Vlinks.Dfk[6 - 1].LinkClass, onClick: this.handleClick
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
          { className: "text-right bg-usepct",
            "data-usepct": $cpu.UserPct
          },
          $cpu.UserPct,
          "%"
        ),
        React.createElement(
          "td",
          { className: "text-right bg-usepct",
            "data-usepct": $cpu.SysPct
          },
          $cpu.SysPct,
          "%"
        ),
        React.createElement(
          "td",
          { className: "text-right bg-usepct",
            "data-usepct": $cpu.WaitPct
          },
          $cpu.WaitPct,
          "%"
        ),
        React.createElement(
          "td",
          { className: "text-right bg-usepct-inverse",
            "data-usepct": $cpu.IdlePct
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
                    { href: Data.Params.Dlinks.CPUd.Less.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.CPUd.Less.ExtraClass != null ? Data.Params.Dlinks.CPUd.Less.ExtraClass : ""), onClick: this.handleClick
                    },
                    "- ",
                    Data.Params.Dlinks.CPUd.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { href: Data.Params.Dlinks.CPUd.More.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.CPUd.More.ExtraClass != null ? Data.Params.Dlinks.CPUd.More.ExtraClass : ""), onClick: this.handleClick
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
                    { href: Data.Params.Nlinks.CPUn.Less.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.CPUn.Less.ExtraClass != null ? Data.Params.Nlinks.CPUn.Less.ExtraClass : ""), onClick: this.handleClick
                    },
                    "- ",
                    Data.Params.Nlinks.CPUn.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { href: Data.Params.Nlinks.CPUn.More.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.CPUn.More.ExtraClass != null ? Data.Params.Nlinks.CPUn.More.ExtraClass : ""), onClick: this.handleClick
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
          { className: "text-right"
          },
          $if.IP
        ),
        React.createElement(
          "td",
          { className: "text-right text-nowrap"
          },
          React.createElement(
            "span",
            { className: "mutext", title: "Total BYTES modulo 4G"
            },
            $if.BytesIn,
            "/",
            $if.BytesOut
          ),
          " ",
          React.createElement(
            "span",
            { title: "BITS per second"
            },
            $if.DeltaBitsIn,
            "/",
            $if.DeltaBitsOut
          )
        ),
        React.createElement(
          "td",
          { className: "text-right text-nowrap"
          },
          React.createElement(
            "span",
            { className: "mutext", title: "Total packets modulo 4G"
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
            { className: "mutext", title: "Total drops/errors modulo 4G"
            },
            $if.DropsIn,
            "/",
            $if.DropsOut,
            "/",
            $if.ErrorsIn,
            "/",
            $if.ErrorsOut
          ),
          " ",
          React.createElement(
            "span",
            { className: $if.DeltaDropsIn == "0" ? "mutext" : ""
            },
            $if.DeltaDropsIn
          ),
          React.createElement(
            "span",
            { className: $if.DeltaDropsOut == "0" || $if.DeltaDropsOut == "0" ? "mutext" : ""
            },
            "/"
          ),
          React.createElement(
            "span",
            { className: $if.DeltaDropsOut == "0" ? "mutext" : ""
            },
            $if.DeltaDropsOut
          ),
          React.createElement(
            "span",
            { className: $if.DeltaErrorsIn == "0" || $if.DeltaErrorsIn == "0" ? "mutext" : ""
            },
            "/"
          ),
          React.createElement(
            "span",
            { className: $if.DeltaErrorsIn == "0" ? "mutext" : ""
            },
            $if.DeltaErrorsIn
          ),
          React.createElement(
            "span",
            { className: $if.DeltaErrorsOut == "0" || $if.DeltaErrorsOut == "0" ? "mutext" : ""
            },
            "/"
          ),
          React.createElement(
            "span",
            { className: $if.DeltaErrorsOut == "0" ? "mutext" : ""
            },
            $if.DeltaErrorsOut
          )
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
                    { href: Data.Params.Dlinks.Ifd.Less.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Ifd.Less.ExtraClass != null ? Data.Params.Dlinks.Ifd.Less.ExtraClass : ""), onClick: this.handleClick
                    },
                    "- ",
                    Data.Params.Dlinks.Ifd.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { href: Data.Params.Dlinks.Ifd.More.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Ifd.More.ExtraClass != null ? Data.Params.Dlinks.Ifd.More.ExtraClass : ""), onClick: this.handleClick
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
                    { href: Data.Params.Nlinks.Ifn.Less.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Ifn.Less.ExtraClass != null ? Data.Params.Nlinks.Ifn.Less.ExtraClass : ""), onClick: this.handleClick
                    },
                    "- ",
                    Data.Params.Nlinks.Ifn.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { href: Data.Params.Nlinks.Ifn.More.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Ifn.More.ExtraClass != null ? Data.Params.Nlinks.Ifn.More.ExtraClass : ""), onClick: this.handleClick
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
                { className: "text-right"
                },
                "IP"
              ),
              React.createElement(
                "th",
                { className: "text-right text-nowrap col-md-3", title: "Bits In/Out per second"
                },
                "IO ",
                React.createElement(
                  "i",
                  null,
                  "b"
                ),
                "ps"
              ),
              React.createElement(
                "th",
                { className: "text-right text-nowrap col-md-3", title: "Packets In/Out per second"
                },
                "Packets IO/s"
              ),
              React.createElement(
                "th",
                { className: "text-right text-nowrap col-md-3", title: "Drops/Errors In/Out per second"
                },
                "Loss IO/s"
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
                    { href: Data.Params.Dlinks.Psd.Less.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Psd.Less.ExtraClass != null ? Data.Params.Dlinks.Psd.Less.ExtraClass : ""), onClick: this.handleClick
                    },
                    "- ",
                    Data.Params.Dlinks.Psd.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { href: Data.Params.Dlinks.Psd.More.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Psd.More.ExtraClass != null ? Data.Params.Dlinks.Psd.More.ExtraClass : ""), onClick: this.handleClick
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
                    { href: Data.Params.Nlinks.Psn.Less.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Psn.Less.ExtraClass != null ? Data.Params.Nlinks.Psn.Less.ExtraClass : ""), onClick: this.handleClick
                    },
                    "- ",
                    Data.Params.Nlinks.Psn.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { href: Data.Params.Nlinks.Psn.More.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Psn.More.ExtraClass != null ? Data.Params.Nlinks.Psn.More.ExtraClass : ""), onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Psk[1 - 1].LinkHref, className: Data.Params.Vlinks.Psk[1 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Psk[2 - 1].LinkHref, className: Data.Params.Vlinks.Psk[2 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Psk[3 - 1].LinkHref, className: Data.Params.Vlinks.Psk[3 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Psk[4 - 1].LinkHref, className: Data.Params.Vlinks.Psk[4 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Psk[5 - 1].LinkHref, className: Data.Params.Vlinks.Psk[5 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Psk[6 - 1].LinkHref, className: Data.Params.Vlinks.Psk[6 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Psk[7 - 1].LinkHref, className: Data.Params.Vlinks.Psk[7 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Psk[8 - 1].LinkHref, className: Data.Params.Vlinks.Psk[8 - 1].LinkClass, onClick: this.handleClick
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
                  { href: Data.Params.Vlinks.Psk[9 - 1].LinkHref, className: Data.Params.Vlinks.Psk[9 - 1].LinkClass, onClick: this.handleClick
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
                    { href: Data.Params.Dlinks.Vgd.Less.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Vgd.Less.ExtraClass != null ? Data.Params.Dlinks.Vgd.Less.ExtraClass : ""), onClick: this.handleClick
                    },
                    "- ",
                    Data.Params.Dlinks.Vgd.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { href: Data.Params.Dlinks.Vgd.More.Href, className: "btn btn-default" + " " + (Data.Params.Dlinks.Vgd.More.ExtraClass != null ? Data.Params.Dlinks.Vgd.More.ExtraClass : ""), onClick: this.handleClick
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
                    { href: Data.Params.Nlinks.Vgn.Less.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Vgn.Less.ExtraClass != null ? Data.Params.Nlinks.Vgn.Less.ExtraClass : ""), onClick: this.handleClick
                    },
                    "- ",
                    Data.Params.Nlinks.Vgn.Less.Text
                  ),
                  React.createElement(
                    "a",
                    { href: Data.Params.Nlinks.Vgn.More.Href, className: "btn btn-default" + " " + (Data.Params.Nlinks.Vgn.More.ExtraClass != null ? Data.Params.Nlinks.Vgn.More.ExtraClass : ""), onClick: this.handleClick
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
