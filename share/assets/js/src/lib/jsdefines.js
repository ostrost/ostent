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
          { className: !Data.Params.Memn.Negative ? "hidden" : "list-group"
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
          { className: Data.Params.Memn.Absolute != 0 ? "table table-hover" : "hidden"
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
          { className: !Data.Params.Dfn.Negative ? "hidden" : "list-group"
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
          { className: Data.Params.Dfn.Absolute != 0 ? "table table-hover" : "hidden"
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
          { className: !Data.Params.CPUn.Negative ? "hidden" : "list-group"
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
          { className: Data.Params.CPUn.Absolute != 0 ? "table table-hover" : "hidden"
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
            { className: "mutext"
            },
            React.createElement(
              "span",
              { title: "Total BYTES In modulo 4G"
              },
              $if.BytesIn
            ),
            "/",
            React.createElement(
              "span",
              { title: "Total BYTES Out modulo 4G"
              },
              $if.BytesOut
            )
          ),
          " ",
          React.createElement(
            "span",
            { title: "BITS In per second"
            },
            $if.DeltaBitsIn
          ),
          "/",
          React.createElement(
            "span",
            { title: "BITS Out per second"
            },
            $if.DeltaBitsOut
          )
        ),
        React.createElement(
          "td",
          { className: "text-right text-nowrap"
          },
          React.createElement(
            "span",
            { className: "mutext"
            },
            React.createElement(
              "span",
              { title: "Total packets In modulo 4G"
              },
              $if.PacketsIn
            ),
            "/",
            React.createElement(
              "span",
              { title: "Total packets Out modulo 4G"
              },
              $if.PacketsOut
            )
          ),
          " ",
          React.createElement(
            "span",
            { title: "Packets In per second"
            },
            $if.DeltaPacketsIn
          ),
          "/",
          React.createElement(
            "span",
            { title: "Packets Out per second"
            },
            $if.DeltaPacketsOut
          )
        ),
        React.createElement(
          "td",
          { className: "text-right text-nowrap"
          },
          React.createElement(
            "span",
            { className: "mutext", title: "Total drops,errors modulo 4G"
            },
            React.createElement(
              "span",
              { title: "Total drops In modulo 4G"
              },
              $if.DropsIn
            ),
            React.createElement(
              "span",
              { className: $if.DropsOut != null ? "" : "hidden"
              },
              "/"
            ),
            React.createElement(
              "span",
              { className: $if.DropsOut != null ? "" : "hidden", title: "Total drops Out modulo 4G"
              },
              $if.DropsOut
            ),
            ",",
            React.createElement(
              "span",
              { title: "Total errors In modulo 4G"
              },
              $if.ErrorsIn
            ),
            "/",
            React.createElement(
              "span",
              { title: "Total errors Out modulo 4G"
              },
              $if.ErrorsOut
            )
          ),
          " ",
          React.createElement(
            "span",
            { className: ($if.DeltaDropsIn == null || $if.DeltaDropsIn == "0") && ($if.DeltaDropsOut == null || $if.DeltaDropsOut == "0") && ($if.DeltaErrorsIn == null || $if.DeltaErrorsIn == "0") && ($if.DeltaErrorsOut == null || $if.DeltaErrorsOut == "0") ? "mutext" : ""
            },
            React.createElement(
              "span",
              { title: "Drops In per second"
              },
              $if.DeltaDropsIn
            ),
            React.createElement(
              "span",
              { className: $if.DeltaDropsOut != null ? "" : "hidden"
              },
              "/"
            ),
            React.createElement(
              "span",
              { className: $if.DeltaDropsOut != null ? "" : "hidden", title: "Drops Out per second"
              },
              $if.DeltaDropsOut
            ),
            ",",
            React.createElement(
              "span",
              { title: "Errors In per second"
              },
              $if.DeltaErrorsIn
            ),
            "/",
            React.createElement(
              "span",
              { title: "Errors Out per second"
              },
              $if.DeltaErrorsOut
            )
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
          { className: !Data.Params.Ifn.Negative ? "hidden" : "list-group"
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
          { className: Data.Params.Ifn.Absolute != 0 ? "table table-hover" : "hidden"
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
                "Packets IO ps"
              ),
              React.createElement(
                "th",
                { className: "text-right text-nowrap col-md-3", title: "Drops,Errors In/Out per second"
                },
                "Loss IO ps"
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
          { className: !Data.Params.Psn.Negative ? "hidden" : "list-group"
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
          { className: Data.Params.Psn.Absolute != 0 ? "table table-hover" : "hidden"
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
    }
  };
});
