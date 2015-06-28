require.config
  shim: {bscollapse: {deps: ['jquery']}} #, bootstrap: {deps: ['jquery']}
  baseUrl: '/js/src'
  urlArgs: "bust=" + (new Date()).getTime()
  paths:
    domReady:  'vendor/requirejs-domready/2.0.1/domReady'
    headroom:  'vendor/headroom/0.7.0/headroom.min'
    jquery:    'vendor/jquery/2.1.4/jquery-2.1.4.min'
    bscollapse:'vendor/bootstrap/3.3.5-collapse/bootstrap.min'
    react:     'vendor/react/0.13.3/react.min'
    jsdefines:   'lib/jsdefines'

# main require
require ['jquery', 'react', 'jsdefines', 'domReady', 'headroom', 'bscollapse'], ($, React, jsdefines) ->
  # domReady, headroom, bscollapse "required" for r.js only.
  updates = undefined # events source. set later
  neweventsource = (onmessage) ->
    conn = null
    sendSearch = (search) ->
      # conn = new EventSource('/index.sse' + search)
      # location.search = search
      console.log('SEARCH', search)
      conn.close() # if conn?
      window.setTimeout(init, 1000)
    init = () ->
      conn = new EventSource('/index.sse' + location.search)
      conn.onopen = () ->
        $(window).bind('popstate', (() ->
          sendSearch(location.search)
          return))
        return

      statesel = 'table thead tr .header a.state'
      again = (e) ->
        $(statesel).unbind('click')
        window.setTimeout(init, 5000) if !e.wasClean
        return

      conn.onclose   = () -> console.log('sse closed (should recover)')
      conn.onerror   = () -> console.log('sse errord (should recover)')
      conn.onmessage = onmessage

      $(statesel).click(() ->
        history.pushState({path: @path}, '', @href)
        sendSearch(@search)
        return false)
      return

    init()
    return {
      sendSearch: sendSearch
      close: () -> conn.close()
    }
  newwebsocket = (onmessage) ->
    conn = null
    sendSearch = (search) -> sendJSON({Search: search})
    sendJSON = (obj) ->
      console.log(JSON.stringify(obj), 'sendJSON')
      # 0 conn.CONNECTING
      # 1 conn.OPEN
      # 2 conn.CLOSING
      # 3 conn.CLOSED
      if !conn? ||
         conn.readyState == conn.CLOSING ||
         conn.readyState == conn.CLOSED
        init()
      if !conn? ||
         conn.readyState != conn.OPEN
        console.log('Not connected, cannot send', obj)
        return
      return conn.send(JSON.stringify(obj))
    init = () ->
      hostport = window.location.hostname +
        (if location.port then ':' + location.port else '')
      conn = new WebSocket('ws://' + hostport + '/index.ws')
      conn.onopen = () ->
        sendSearch(location.search)
        $(window).bind('popstate', (() ->
          sendSearch(location.search)
          return))
        return

      statesel = 'table thead tr .header a.state'
      again = (e) ->
        $(statesel).unbind('click')
        window.setTimeout(init, 5000) if !e.wasClean
        return

      conn.onclose   = again
      conn.onerror   = again
      conn.onmessage = onmessage

      $(statesel).click(() ->
        history.pushState({path: @path}, '', @href)
        sendSearch(@search)
        return false)
      return

    init()
    return {
      sendSearch: sendSearch
      close: () -> conn.close()
    }

  @IFCLASS = React.createClass
    getInitialState: () -> { # a global Data
      Links:        Data.Links
      IFbytes:      Data.IFbytes
      IFerrors:     Data.IFerrors
      IFpackets:    Data.IFpackets
      ExpandableIF: Data.ExpandableIF
      ExpandtextIF: Data.ExpandtextIF
    }
    render: () ->
      Data = @state
      return jsdefines.panelif.bind(this)(Data,
            (jsdefines.ifpackets_rows(Data, $if) for $if in Data?.IFpackets?.List ? []),
            (jsdefines.iferrors_rows(Data, $if) for $if in Data?.IFerrors?.List ? []),
            (jsdefines.ifbytes_rows(Data, $if) for $if in Data?.IFbytes?.List ? []))
    handleChange: (e) ->
      href = '?' + e.target.name + '=' + e.target.value + '&' + location.search.substr(1)
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined
    handleClick: (e) ->
      href = e.target.getAttribute('href')
      history.pushState({}, '', href)
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined

  @DFCLASS = React.createClass
    getInitialState: () -> { # a global Data
      Links:        Data.Links
      DFbytes:      Data.DFbytes
      DFinodes:     Data.DFinodes
      ExpandableDF: Data.ExpandableDF
      ExpandtextDF: Data.ExpandtextDF
    }
    render: () ->
      Data = @state
      return jsdefines.paneldf.bind(this)(Data,
             (jsdefines.dfinodes_rows(Data, $disk) for $disk in Data?.DFinodes?.List ? []),
             (jsdefines.dfbytes_rows(Data, $disk) for $disk in Data?.DFbytes?.List ? []))
    handleChange: (e) ->
      href = '?' + e.target.name + '=' + e.target.value + '&' + location.search.substr(1)
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined
    handleClick: (e) ->
      href = e.target.getAttribute('href')
      history.pushState({}, '', href)
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined

  @LabelClassColorPercent = (p) ->
    return "label label-danger"  if p.length > 2
    return "label label-danger"  if p.length > 1 && p[0] == '9'
    return "label label-warning" if p.length > 1 && p[0] == '8'
    return "label label-success" if p.length > 1 && p[0] == '1'
    return "label label-info"    if p.length > 1
    return "label label-success"

  @MEMtableCLASS = React.createClass
    getInitialState: () -> {
      Links:  Data.Links,  # a global Data
      MEM:    Data.MEM,    # a global Data
    }
    render: () ->
      Data = @state
      return jsdefines.panelmem.bind(this)(Data, (jsdefines.mem_rows(Data, $mem
      ) for $mem in Data?.MEM?.List ? []))
    handleChange: (e) ->
      href = '?' + e.target.name + '=' + e.target.value + '&' + location.search.substr(1)
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined
    handleClick: (e) ->
      href = e.target.getAttribute('href')
      history.pushState({}, '', href)
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined

  @CPUtableCLASS = React.createClass
    getInitialState: () -> {
      Links: Data.Links, # a global Data
      CPU:   Data.CPU,   # a global Data
    }
    render: () ->
      Data = @state
      return jsdefines.panelcpu.bind(this)(Data, (jsdefines.cpu_rows(Data, $core
      ) for $core in Data?.CPU?.List ? []))
    handleChange: (e) ->
      href = '?' + e.target.name + '=' + e.target.value + '&' + location.search.substr(1)
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined
    handleClick: (e) ->
      href = e.target.getAttribute('href')
      history.pushState({}, '', href)
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined

  @PStableCLASS = React.createClass
    getInitialState: () -> {
      Links:   Data.Links,  # a global Data
      PStable: Data.PStable # a global Data
    }
    render: () ->
      Data = @state
      return jsdefines.panelps.bind(this)(Data, (jsdefines.ps_rows(Data, $proc
      ) for $proc in Data?.PStable?.List ? []))
    handleChange: (e) ->
      href = '?' + e.target.name + '=' + e.target.value + '&' + location.search.substr(1)
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined
    handleClick: (e) ->
      href = e.target.getAttribute('href')
      history.pushState({}, '', href)
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined

  @VGtableCLASS = React.createClass
    getInitialState: () -> { # a global Data:
      Links:           Data.Links
      VagrantMachines: Data.VagrantMachines
      VagrantError:    Data.VagrantError
      VagrantErrord:   Data.VagrantErrord
    }
    render: () ->
      Data = @state
      if Data?.VagrantErrord? and Data.VagrantErrord
        rows = [jsdefines.vagrant_error.bind(this)(Data)]
      else
        rows = (jsdefines.vagrant_rows.bind(this)(Data, $mach
        ) for $mach in Data?.VagrantMachines?.List ? [])
      return jsdefines.panelvg.bind(this)(Data, rows)
    handleChange: (e) ->
      href = '?' + e.target.name + '=' + e.target.value + '&' + location.search.substr(1)
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined
    handleClick: (e) ->
      href = e.target.getAttribute('href')
      history.pushState({}, '', href)
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined

  @NewTextCLASS = (reduce) -> React.createClass
    newstate: (data) ->
      v = reduce(data)
      return {Text: v} if v?
    getInitialState: () -> @newstate(Data) # a global Data
    render: () -> React.DOM.span(null, @state.Text)

  @setState = (obj, data) ->
    if data?
      delete data[key] for key of data when !data[key]?
      return obj.setState(data)

  update = () ->
    # coffeelint: disable=max_line_length
    ip       = React.render(React.createElement(NewTextCLASS((data) -> data?.IP       )), $('#ip'      )   .get(0)) if data?.IP?
    hostname = React.render(React.createElement(NewTextCLASS((data) -> data?.Hostname )), $('#hostname')   .get(0))
    uptime   = React.render(React.createElement(NewTextCLASS((data) -> data?.Uptime   )), $('#uptime'  )   .get(0))
    la       = React.render(React.createElement(NewTextCLASS((data) -> data?.LA       )), $('#la'      )   .get(0))

    memtable  = React.render(React.createElement(MEMtableCLASS), document.getElementById('mem' +'-'+ 'table'))
    pstable   = React.render(React.createElement(PStableCLASS),  document.getElementById('ps'  +'-'+ 'table'))
    dftable   = React.render(React.createElement(DFCLASS),       document.getElementById('df'  +'-'+ 'table'))
    cputable  = React.render(React.createElement(CPUtableCLASS), document.getElementById('cpu' +'-'+ 'table'))
    iftable   = React.render(React.createElement(IFCLASS),       document.getElementById('if'  +'-'+ 'table'))
    vgtable   = React.render(React.createElement(VGtableCLASS),  document.getElementById('vg'  +'-'+ 'table'))
    # coffeelint: enable=max_line_length

    onmessage = (event) ->
      data = JSON.parse(event.data)
      return if !data?

      if data.Reload? and data.Reload
        window.setTimeout((() -> location.reload(true)), 5000)
        window.setTimeout(updates.close, 2000)
        console.log('in 5s: location.reload(true)')
        console.log('in 2s: updates.close()')
        return

      setState(ip,        ip      .newstate(data)) if ip?
      setState(hostname,  hostname.newstate(data))
      setState(uptime,    uptime  .newstate(data))
      setState(la,        la      .newstate(data))

      setState(pstable,  {Links: data.Links, PStable:  data.PStable})
      setState(memtable, {Links: data.Links, MEM: data.MEM})
      setState(cputable, {Links: data.Links, CPU: data.CPU})
      setState(iftable, {
        Links:        data.Links
        IFbytes:      data.IFbytes
        IFerrors:     data.IFerrors
        IFpackets:    data.IFpackets
        ExpandableIF: data.ExpandableIF
        ExpandtextIF: data.ExpandtextIF
      })
      setState(dftable, {
        Links:        data.Links
        DFbytes:      data.DFbytes
        DFinodes:     data.DFinodes
        ExpandableDF: data.ExpandableDF
        ExpandtextDF: data.ExpandtextDF
      })
      setState(vgtable, {
        Links:           data.Links,
        VagrantMachines: data.VagrantMachines,
        VagrantError:    data.VagrantError,
        VagrantErrord:   data.VagrantErrord
      })

      if data.Location?
        history.pushState({}, '', data.Location)

      return

    updates = newwebsocket(onmessage)
  # updates = neweventsource(onmessage)
    return # end of `update'

  require ['domReady', 'jquery', 'headroom'], (domReady, $) ->
    # headroom does not export anything
    domReady () ->
      (new window.Headroom(document.querySelector('nav'), {
        offset: 20 # ~padding-top of a container row
      })).init()

      # referencing upper-scope `update'
      update() unless (42 for param in location.search.substr(1).split(
        '&') when (param.split('=')[0] == 'still')).length

      return # return from domReady
    return # end of sub`require'
  return # end of main `require'

# Local variables:
# coffee-tab-width: 2
# End:
