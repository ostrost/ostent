require.config
  shim:
    bscollapse: {deps: ['jquery']}, # bootstrap:
    reactDOM:   {deps: ['react']},
  urlArgs: "bust=" + (new Date()).getTime()
  paths:
    domReady:   'vendor/requirejs-domready/2.0.1/domReady'
    jquery:     'vendor/jquery/2.1.4/jquery.min'
    bscollapse: 'vendor/bootstrap/3.3.5-collapse/bootstrap.min'
    react:      'vendor/react/0.14.0/react-with-addons.min'
    reactDOM:   'vendor/react/0.14.0/react-dom.min'
    jsdefines:  'lib/jsdefines'

# main require
require ['jquery', 'react', 'reactDOM', 'jsdefines', 'domReady', 'bscollapse'],
($, React, ReactDOM, jsdefines) ->
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

      again = (e) ->
        window.setTimeout(init, 5000) if !e.wasClean
        return

      conn.onclose   = () -> console.log('sse closed (should recover)')
      conn.onerror   = () -> console.log('sse errord (should recover)')
      conn.onmessage = onmessage
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

      again = (e) ->
        window.setTimeout(init, 5000) if !e.wasClean
        return

      conn.onclose   = again
      conn.onerror   = again
      conn.onmessage = onmessage
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
    mixins: [React.addons.PureRenderMixin, HandlerMixin]
    getInitialState: () -> @Reduce(Data) # a global Data
    Reduce: (data) -> {Params: data.Params, IF: data.IF}
    render: () ->
      Data = @state
      rows = (jsdefines.if_rows(Data, $if) for $if in Data?.IF?.List ? [])
      return (jsdefines.panelif.bind(this)(Data, rows))

  @DFClass = React.createClass
    mixins: [React.addons.PureRenderMixin, HandlerMixin]
    getInitialState: () -> @Reduce(Data) # a global Data
    Reduce: (data) -> {Params: data.Params, DF: data.DF}
    render: () ->
      Data = @state
      rows = (jsdefines.df_rows(Data, $df) for $df in Data?.DF?.List ? [])
      return (jsdefines.paneldf.bind(this)(Data, rows))

  @MEMClass = React.createClass
    mixins: [React.addons.PureRenderMixin, HandlerMixin]
    getInitialState: () -> @Reduce(Data) # a global Data
    Reduce: (data) -> {Params: data.Params, MEM: data.MEM}
    render: () ->
      Data = @state
      rows = (jsdefines.mem_rows(Data, $mem) for $mem in Data?.MEM?.List ? [])
      return (jsdefines.panelmem.bind(this)(Data, rows))

  @CPUClass = React.createClass
    mixins: [React.addons.PureRenderMixin, HandlerMixin]
    getInitialState: () -> @Reduce(Data) # a global Data
    Reduce: (data) -> {Params: data.Params, CPU: data.CPU}
    render: () ->
      Data = @state
      rows = (jsdefines.cpu_rows(Data, $cpu) for $cpu in Data?.CPU?.List ? [])
      return (jsdefines.panelcpu.bind(this)(Data, rows))

  @PSClass = React.createClass
    mixins: [React.addons.PureRenderMixin, HandlerMixin]
    getInitialState: () -> @Reduce(Data) # a global Data
    Reduce: (data) -> {Params: data.Params, PS: data.PS}
    render: () ->
      Data = @state
      rows = (jsdefines.ps_rows(Data, $ps) for $ps in Data?.PS?.List ? [])
      return (jsdefines.panelps.bind(this)(Data, rows))

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
      ReactDOM.render(React.createElement(cl), document.getElementById(id))
    HN  = render('hn', TextClass((data) -> data?.HN))
    UP  = render('up', TextClass((data) -> data?.UP))
    LA  = render('la', TextClass((data) -> data?.LA))
    MEM = render('mem', MEMClass)
    PS  = render('ps',  PSClass)
    DF  = render('df',  DFClass)
    CPU = render('cpu', CPUClass)
    IF  = render('if',  IFClass)

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

      setState(HN,  HN.Reduce(data))
      setState(UP,  UP.Reduce(data))
      setState(LA,  LA.Reduce(data))
      setState(PS,  PS.Reduce(data))
      setState(MEM, MEM.Reduce(data))
      setState(CPU, CPU.Reduce(data))
      setState(IF,  IF.Reduce(data))
      setState(DF,  DF.Reduce(data))

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
