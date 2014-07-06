
@newwebsocket = (onmessage) ->
        conn = null
        sendSearch = (search) -> sendJSON({Search: search})
        sendClient = (client) ->
                console.log(JSON.stringify(client), 'sendClient')
                sendJSON({Client: client})
        sendJSON = (obj) ->
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
                conn.send(JSON.stringify(obj))
        init = () ->
                hostport = window.location.hostname + (
                        if location.port
                                ':' + location.port
                        else '')
                conn = new WebSocket('ws://' + hostport + '/ws')
                conn.onopen = () ->
                        sendSearch(location.Search)
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
                        false)
                return

        init()
        return {
                sendClient: sendClient
                sendSearch: sendSearch
                close: () -> conn.close()
                getconn: () -> conn
        }

@HideClass = React.createClass
        reduce: (data) ->
                if data? and data.Client?
                        value = data.Client[@props.key]
                        return {Hide: value} if value isnt undefined
        getInitialState:   () -> @reduce(Data) # a global Data
        componentDidMount: () -> @props.$click_el.click(@click)
        render: () ->
                @props.$collapse_el.collapse(if @state.Hide then 'hide' else 'show')
                return React.DOM.span(null, null)
        click: (e) ->
                (S = {})[@props.key] = !@state.Hide
                websocket.sendClient(S)
                e.stopPropagation() # preserves checkbox/radio
                e.preventDefault()  # checked/selected state
                undefined

@ShowSwapClass = React.createClass
        getInitialState: () -> ShowSwapClass.reduce(Data) # a global Data
        statics:
                reduce: (data) ->
                        if data? and data.Client?
                                S = {}
                                S.HideSWAP = data.Client.HideSWAP if data.Client.HideSWAP isnt undefined
                                S.HideMEM  = data.Client.HideMEM  if data.Client.HideMEM  isnt undefined
                                S

        componentDidMount: () -> @props.$el.click(@click)
        render: () ->
                @props.$el[if !@state.HideSWAP then 'addClass' else 'removeClass']('active')
                return React.DOM.span(null, @props.$el.text())
        click: (e) ->
                S = {HideSWAP: !@state.HideSWAP}
                S.HideMEM = false if @state.HideMEM
                websocket.sendClient(S)
                e.stopPropagation() # preserves checkbox/radio
                e.preventDefault()  # checked/selected state
                undefined
