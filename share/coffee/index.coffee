require.config
  shim:
    bscollapse: {deps: ['jquery']}, # bootstrap:
    reactDOM:   {deps: ['react']},
  urlArgs: "bust=" + (new Date()).getTime()
  paths:
    domReady:   'vendor/requirejs-domready/2.0.1/domReady'
    jquery:     'vendor/jquery/2.1.4/jquery.min'
    bscollapse: 'vendor/bootstrap/3.3.5-collapse/bootstrap.min'
    react:      'vendor/react/0.14.2/react-with-addons.min'
    reactDOM:   'vendor/react/0.14.2/react-dom.min'
    jsdefines:  'lib/jsdefines'

# main require
require ['jquery', 'react', 'reactDOM', 'jsdefines', 'domReady', 'bscollapse'],
($, React, ReactDOM, jsdefines) ->
  # domReady, bscollapse "required" for r.js only.
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
        return null
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

  update = () ->
    render_define = (el) ->
      cl = jsdefines[$(el).attr('data-define')]
      ReactDOM.render(React.createElement(cl), el)
    els = (render_define(el) for el in $('.updates'))

    onmessage = (event) ->
      data = JSON.parse(event.data)
      return if !data?

      if data.Reload? and data.Reload
        window.setTimeout((() -> location.reload(true)), 5000)
        window.setTimeout(window.updates.close, 2000)
        console.log('in 5s: location.reload(true)')
        console.log('in 2s: window.updates.close()')
        return

      if data.Error?
        console.log 'Error', data.Error
        return

      el.NewState(data) for el in els
      return

    window.updates = newwebsocket(onmessage)
  # window.updates = neweventsource(onmessage)
    return # end of `update'

  require ['domReady', 'jquery'], (domReady, $) ->
    domReady () ->
      update() unless (42 for param in location.search.substr(1).split(
        '&') when (param.split('=')[0] == 'still')).length
      return null

# Local variables:
# coffee-tab-width: 2
# End:
