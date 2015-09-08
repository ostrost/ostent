require.config
  shim: {bscollapse: {deps: ['jquery']}} #, bootstrap: {deps: ['jquery']}
  baseUrl: '/js/src'
  urlArgs: "bust=" + (new Date()).getTime()
  paths:
    domReady:  'vendor/requirejs-domready/2.0.1/domReady'
    jquery:    'vendor/jquery/2.1.4/jquery.min'
    bscollapse:'vendor/bootstrap/3.3.5-collapse/bootstrap.min'
    react:     'vendor/react/0.13.3/react.min'
    jsdefines: 'lib/jsdefines'

# main require
require ['jquery', 'react', 'jsdefines', 'domReady', 'bscollapse'],
($, React, jsdefines) ->
  # domReady, bscollapse "required" for r.js only.
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
    sendSearch = (search) ->
      console.log 'Search', search
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
        console.log('Not connected, cannot send search', search)
        return
      return conn.send(JSON.stringify({Search: search}))
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

  HandlerMixin =
    handleChange: (e) -> @handle(e, false,
     '?' + e.target.name +
     '=' + e.target.value +
     '&' + location.search.substr(1))
    handleClick: (e) ->
      href = e.target.getAttribute('href')
      href = $(e.target).parent().get(0).getAttribute('href') if !href?
      @handle(e, true, href)
    handle: (e, ps, href) ->
      history.pushState({}, '', href) if ps
      updates.sendSearch(href)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined

  @IFClass = React.createClass
    mixins: [HandlerMixin]
    getInitialState: () -> { # a global Data
      Params: Data.Params
      IF:     Data.IF
    }
    render: () ->
      Data = @state
      return jsdefines.panelif.bind(this)(Data, (jsdefines.if_rows(Data, $if
      ) for $if in Data?.IF?.List ? []))

  @DFClass = React.createClass
    mixins: [HandlerMixin]
    getInitialState: () -> { # a global Data
      Params: Data.Params
      DF:     Data.DFbytes
    }
    render: () ->
      Data = @state
      return jsdefines.paneldf.bind(this)(Data,
             (jsdefines.df_rows(Data, $df) for $df in Data?.DF?.List ? []))

  @MEMClass = React.createClass
    mixins: [HandlerMixin]
    getInitialState: () -> { # a global Data
      Params: Data.Params
      MEM:    Data.MEM
    }
    render: () ->
      Data = @state
      return jsdefines.panelmem.bind(this)(Data, (jsdefines.mem_rows(Data, $mem
      ) for $mem in Data?.MEM?.List ? []))

  @CPUClass = React.createClass
    mixins: [HandlerMixin]
    getInitialState: () -> { # a global Data
      Params: Data.Params
      CPU:    Data.CPU
    }
    render: () ->
      Data = @state
      return jsdefines.panelcpu.bind(this)(Data, (jsdefines.cpu_rows(Data, $cpu
      ) for $cpu in Data?.CPU?.List ? []))

  @PSClass = React.createClass
    mixins: [HandlerMixin]
    getInitialState: () -> { # a global Data
      Params: Data.Params
      PS:     Data.PS
    }
    render: () ->
      Data = @state
      return jsdefines.panelps.bind(this)(Data, (jsdefines.ps_rows(Data, $ps
      ) for $ps in Data?.PS?.List ? []))

  @VGClass = React.createClass
    mixins: [HandlerMixin]
    getInitialState: () -> { # a global Data:
      Params:          Data.Params
      VagrantMachines: Data.VagrantMachines
      VagrantError:    Data.VagrantError
      VagrantErrord:   Data.VagrantErrord
    }
    render: () ->
      Data = @state
      if Data?.VagrantErrord? and Data.VagrantErrord
        rows = [jsdefines.vg_error.bind(this)(Data)]
      else
        rows = (jsdefines.vg_rows.bind(this)(Data, $vgm
        ) for $vgm in Data?.VagrantMachines?.List ? [])
      return jsdefines.panelvg.bind(this)(Data, rows)

  @TextClass = (reduce) -> React.createClass
    Reduce: (data) ->
      v = reduce(data)
      return {Text: v} if v?
    getInitialState: () -> @Reduce(Data) # a global Data
    render: () -> React.DOM.span(null, @state.Text)

  @setState = (obj, data) ->
    if data?
      delete data[key] for key of data when !data[key]?
      return obj.setState(data)

  update = () ->
    render = (id, cl) ->
      React.render(React.createElement(cl), document.getElementById(id))
    IP  = render('ip', TextClass((data) -> data?.IP)) if data?.IP?
    HN  = render('hn', TextClass((data) -> data?.HN))
    UP  = render('up', TextClass((data) -> data?.UP))
    LA  = render('la', TextClass((data) -> data?.LA))
    MEM = render('mem', MEMClass)
    PS  = render('ps',  PSClass)
    DF  = render('df',  DFClass)
    CPU = render('cpu', CPUClass)
    IF  = render('if',  IFClass)
    VG  = render('vg',  VGClass)

    onmessage = (event) ->
      data = JSON.parse(event.data)
      return if !data?

      if data.Reload? and data.Reload
        window.setTimeout((() -> location.reload(true)), 5000)
        window.setTimeout(updates.close, 2000)
        console.log('in 5s: location.reload(true)')
        console.log('in 2s: updates.close()')
        return

      if data.Error?
        console.log 'Error', data.Error
        return

      setState(IP, IP.Reduce(data)) if IP?
      setState(HN, HN.Reduce(data))
      setState(UP, UP.Reduce(data))
      setState(LA, LA.Reduce(data))

      setState(PS,  {Params: data.Params, PS:  data.PS })
      setState(MEM, {Params: data.Params, MEM: data.MEM})
      setState(CPU, {Params: data.Params, CPU: data.CPU})
      setState(IF,  {Params: data.Params, IF:  data.IF })
      setState(DF,  {Params: data.Params, DF:  data.DF })
      setState(VG, {
        Params:          data.Params,
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

  require ['domReady', 'jquery'], (domReady, $) ->
    domReady () ->
      update() unless (42 for param in location.search.substr(1).split(
        '&') when (param.split('=')[0] == 'still')).length

# Local variables:
# coffee-tab-width: 2
# End:
