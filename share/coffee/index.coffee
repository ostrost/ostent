require.config
  shim: {bootstrap: {deps: ['jquery']}}
  baseUrl: '/js/devel'
  paths:
    domReady:  'vendor/min/requirejs-domReady/2.0.1/domReady'
    headroom:  'vendor/min/headroom/0.7.0/headroom.min'
    jquery:    'vendor/min/jquery/2.1.3/jquery-2.1.3.min'
    bootstrap: 'vendor/min/bootstrap/3.3.2/bootstrap.min'
    react:     'vendor/min/react/0.13.1/react.min'
    jscript:   'gen/jscript'

# main require
require ['jquery', 'bootstrap', 'react', 'jscript'], ($, _, React, jscript) ->
  updates = undefined # events source. set later
  neweventsource = (onmessage) ->
    conn = null
    sendSearch = (search) ->
      # conn = new EventSource('/index.sse' + search)
      # location.search = search
      console.log('SEARCH', search)
      conn.close() # if conn?
      window.setTimeout(init, 1000)
    sendClient = (client) ->
      # TODO
      return
      console.log(JSON.stringify(client), 'sendClient')
      obj = {Client: client}
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
      sendClient: sendClient
      sendSearch: sendSearch
      close: () -> conn.close()
    }
  newwebsocket = (onmessage) ->
    conn = null
    sendSearch = (search) -> sendJSON({Search: search})
    sendClient = (client) ->
      console.log(JSON.stringify(client), 'sendClient')
      return sendJSON({Client: client})
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
      sendClient: sendClient
      sendSearch: sendSearch
      close: () -> conn.close()
    }

  @IFbytesCLASS = React.createClass
    getInitialState: () -> Data.IFbytes # a global Data
    render: () ->
      Data = {IFbytes: @state}
      return jscript.ifbytes_table(Data, (jscript.ifbytes_rows(Data, $if
      ) for $if in Data?.IFbytes?.List ? []))

  @IFerrorsCLASS = React.createClass
    getInitialState: () -> Data.IFerrors # a global Data
    render: () ->
      Data = {IFerrors: @state}
      return jscript.iferrors_table(Data, (jscript.iferrors_rows(Data, $if
      ) for $if in Data?.IFerrors?.List ? []))

  @IFpacketsCLASS = React.createClass
    getInitialState: () -> Data.IFpackets # a global Data
    render: () ->
      Data = {IFpackets: @state}
      return jscript.ifpackets_table(Data, (jscript.ifpackets_rows(Data, $if
      ) for $if in Data?.IFpackets?.List ? []))

  @DFbytesCLASS = React.createClass
    getInitialState: () -> {
      Links:   Data.Links,  # a global Data
      DFbytes: Data.DFbytes # a global Data
    }
    render: () ->
      Data = @state
      return jscript.dfbytes_table(Data, (jscript.dfbytes_rows(Data, $disk
      ) for $disk in Data?.DFbytes?.List ? []))

  @DFinodesCLASS = React.createClass
    getInitialState: () -> {
      Links:    Data.Links,   # a global Data
      DFinodes: Data.DFinodes # a global Data
    }
    render: () ->
      Data = @state
      return jscript.dfinodes_table(Data, (jscript.dfinodes_rows(Data, $disk
      ) for $disk in Data?.DFinodes?.List ? []))

  @MEMtableCLASS = React.createClass
    getInitialState: () -> Data.MEM # a global Data
    render: () ->
      Data = {MEM: @state}
      return jscript.mem_table(Data, (jscript.mem_rows(Data, $mem
      ) for $mem in Data?.MEM?.List ? []))

  @CPUtableCLASS = React.createClass
    getInitialState: () -> Data.CPU # a global Data
    render: () ->
      Data = {CPU: @state}
      return jscript.cpu_table(Data, (jscript.cpu_rows(Data, $core
      ) for $core in Data?.CPU?.List ? []))

  @PStableCLASS = React.createClass
    getInitialState: () -> {
      Links:   Data.Links,  # a global Data
      PStable: Data.PStable # a global Data
    }
    render: () ->
      Data = @state
      return jscript.ps_table(Data, (jscript.ps_rows(Data, $proc
      ) for $proc in Data?.PStable?.List ? []))

  @VGtableCLASS = React.createClass
    getInitialState: () -> { # a global Data:
      VagrantMachines: Data.VagrantMachines
      VagrantError:    Data.VagrantError
      VagrantErrord:   Data.VagrantErrord
    }
    render: () ->
      Data = @state
      if Data?.VagrantErrord? and Data.VagrantErrord
        rows = [jscript.vagrant_error(Data)]
      else
        rows = (jscript.vagrant_rows(Data, $mach
        ) for $mach in Data?.VagrantMachines?.List ? [])
      return jscript.vagrant_table(Data, rows)

  @addDiv = (sel) -> sel.append('<div />').find('div').get(0)

  @HideClass = React.createClass
    statics: component: (opt) ->
      el = addDiv(opt.$button_el)
      React.render(React.createElement(HideClass, opt), el)

    reduce: (data) ->
      if data?.Client?
        value = data.Client[@props.xkey]
        return {Hide: value} if value isnt undefined
      return null
    getInitialState: () -> @reduce(Data) # a global Data
    componentDidMount: () -> @props.$button_el.click(@click)
    render: () ->
      @props.$collapse_el.collapse(if @state.Hide then 'hide' else 'show')
      buttonactive =  @state.Hide
      buttonactive = !@state.Hide if (
        @props.reverseActive? and @props.reverseActive)
      opclass = if buttonactive then 'addClass' else 'removeClass'
      @props.$button_el[opclass]('active')
      return null
    click: (e) ->
      (S = {})[@props.xkey] = !@state.Hide
      updates.sendClient(S)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined

  @ButtonClass = React.createClass
    statics: component: (opt) ->
      el = addDiv(opt.$button_el)
      React.render(React.createElement(ButtonClass, opt), el)

    reduce: (data) ->
      if data?.Client?
        S = {}
        # coffeelint: disable=max_line_length
        S.Hide = data.Client[@props.Khide] if                   data.Client[@props.Khide] isnt undefined # Khide is a required prop
        S.Able = data.Client[@props.Kable] if @props.Kable? and data.Client[@props.Kable] isnt undefined
        S.Send = data.Client[@props.Ksend] if @props.Ksend? and data.Client[@props.Ksend] isnt undefined
        S.Text = data.Client[@props.Ktext] if @props.Ktext? and data.Client[@props.Ktext] isnt undefined
        # coffeelint: enable=max_line_length
        return S
    getInitialState: () -> @reduce(Data) # a global Data
    componentDidMount: () -> @props.$button_el.click(@click)
    render: () ->
      if @props.Kable
        able = @state.Able
        able = !able if not (@props.Kable.indexOf('not') > -1) # That's a hack
        opclass = if able then 'addClass' else 'removeClass'
        @props.$button_el.prop('disabled', able)
        @props.$button_el[opclass]('disabled')
      opclass = if @state.Send then 'addClass' else 'removeClass'
      @props.$button_el[opclass]('active') if @props.Ksend?
      @props.$button_el.text(@state.Text) if @props.Ktext?
      return null
    click: (e) ->
      S = {}
      S[@props.Khide] = !@state.Hide if @state.Hide?  and @state.Hide
      S[@props.Ksend] = !@state.Send if @props.Ksend? and @state.Send?
      S[@props.Ksig]  =  @props.Vsig if @props.Ksig?
      updates.sendClient(S)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined

  @TabsClass = React.createClass
    statics: component: (opt) ->
      el = addDiv(opt.$button_el)
      React.render(React.createElement(TabsClass, opt), el)

    reduce: (data) ->
      if data?.Client?
        S = {}
        # coffeelint: disable=max_line_length
        S.Hide = data.Client[@props.Khide] if                   data.Client[@props.Khide] isnt undefined # Khide is a required prop
        S.Send = data.Client[@props.Ksend] if @props.Ksend? and data.Client[@props.Ksend] isnt undefined
        # coffeelint: enable=max_line_length
        return S
    getInitialState: () -> @reduce(Data) # a global Data
    componentDidMount: () ->
      @props.$button_el.click(@clicktab)
      @props.$hidebutton_el.click(@clickhide)
    render: () ->
      if @state.Hide
        @props.$collapse_el.collapse('hide')
        @props.$hidebutton_el.addClass('active')
        return null
      @props.$hidebutton_el.removeClass('active')
      curtabid = +@state.Send # MUST be an int
      nots = @props.$collapse_el.not('[data-tabid="'+ curtabid + '"]')
      $(el).collapse('hide') for el in nots
      $(@props.$collapse_el.not(nots)).collapse('show')
      activeClass = (el) ->
        xel = $(el)
        tabid_attr = +xel.attr('data-tabid') # an int
        opclass = if tabid_attr == curtabid then 'addClass' else 'removeClass'
        xel[opclass]('active')
        return
      activeClass(el) for el in @props.$button_el
      return null
    clicktab: (e) ->
      S = {}
      # +"STRING" to make it an int
      S[@props.Ksend] = +$( $(e.target).attr('href') ).attr('data-tabid')
      S[@props.Khide] = false if @state.Hide? and @state.Hide
      updates.sendClient(S)
      e.preventDefault()
      e.stopPropagation() # don't change checkbox/radio state
      return undefined
    clickhide: (e) ->
      (S = {})[@props.Khide] = !@state.Hide
      updates.sendClient(S)
      e.stopPropagation() # preserves checkbox/radio
      e.preventDefault()  # checked/selected state
      return undefined

  @RefreshInputClass = React.createClass
    statics: component: (opt) ->
      sel = opt.sel; delete opt.$
      opt.$input_el = sel.find('.refresh-input')
      opt.$group_el = sel.find('.refresh-group')
      el = addDiv(opt.$input_el)
      React.render(React.createElement(RefreshInputClass, opt), el)

    reduce: (data) ->
      if data?.Client? and (
        data.Client[@props.K]? or
        data.Client[@props.Kerror]?)
        S = {}
        S.Value = data.Client[@props.K]      if data.Client[@props.K]?
        S.Error = data.Client[@props.Kerror] if data.Client[@props.Kerror]?
        return S
    getInitialState: () ->
      S = @reduce(Data) # a global Data
      delete S.Value # to make input empty initially
      return S

    componentDidMount: () -> @props.$input_el.on('input', @submit)
    render: () ->
      # The initial render should not place a value.
      # The check relied on @isMounted() until it was deprecated.
      # getInitialState now deletes .Value.
      @props.$input_el.prop('value', @state.Value) if (
        @state.Value? and !@state.Error)
      opclass = if @state.Error then 'addClass' else 'removeClass'
      @props.$group_el[opclass]('has-warning')
      return null
    submit: (e) ->
      (S = {})[@props.Ksig] = $(e.target).val()
      updates.sendClient(S)
      e.preventDefault()
      e.stopPropagation() # don't change checkbox/radio state
      return undefined

  @NewTextCLASS = (reduce) -> React.createClass
    newstate: (data) ->
      v = reduce(data)
      return {Text: v} if v?
    getInitialState: () -> @newstate(Data) # a global Data
    render: () -> React.DOM.span(null, @state.Text)

  @AlertClass = React.createClass
    show: () -> return @state.Error?
    newstate: (data) ->
      error = data.Client?.DebugError # data.Error
      a = { # return
        Error: error
        ErrorText: @state?.ErrorText
        Changed: @state? and error? and error != @state.Error
      }
      a.ErrorText = a.Error if a.Changed and a.Error?
      console.log('newstate', a)
      return a
    getInitialState: () -> @newstate(Data) # a global Data
    render: () ->
      if @state.Changed
        console.log('show', @state) # , if @show() then 'show' else 'hide')
        @props.$collapse_el.collapse('show') if @show()
      return React.DOM.span(null, @state.ErrorText)

  @setState = (obj, data) ->
    if data?
      delete data[key] for key of data when !data[key]?
      return obj.setState(data)

  update = () ->
    # coffeelint: disable=max_line_length
    hideconfigmem = HideClass.component({xkey: 'HideconfigMEM', $collapse_el: $('#memconfig'), $button_el: $('header a[href="#mem"]'), reverseActive: true})
    hideconfigif  = HideClass.component({xkey: 'HideconfigIF',  $collapse_el: $('#ifconfig'),  $button_el: $('header a[href="#if"]'),  reverseActive: true})
    hideconfigcpu = HideClass.component({xkey: 'HideconfigCPU', $collapse_el: $('#cpuconfig'), $button_el: $('header a[href="#cpu"]'), reverseActive: true})
    hideconfigdf  = HideClass.component({xkey: 'HideconfigDF',  $collapse_el: $('#dfconfig'),  $button_el: $('header a[href="#df"]'),  reverseActive: true})
    hideconfigps  = HideClass.component({xkey: 'HideconfigPS',  $collapse_el: $('#psconfig'),  $button_el: $('header a[href="#ps"]'),  reverseActive: true})
    hideconfigvg  = HideClass.component({xkey: 'HideconfigVG',  $collapse_el: $('#vgconfig'),  $button_el: $('header a[href="#vg"]'),  reverseActive: true})

    hideram = HideClass.component({xkey: 'HideRAM', $collapse_el: $('#mem'), $button_el: $('#memconfig').find('.hiding')})
    hidecpu = HideClass.component({xkey: 'HideCPU', $collapse_el: $('#cpu'), $button_el: $('#cpuconfig').find('.hiding')})
    hideps  = HideClass.component({xkey: 'HidePS',  $collapse_el: $('#ps'),  $button_el: $('#psconfig') .find('.hiding')})
    hidevg  = HideClass.component({xkey: 'HideVG',  $collapse_el: $('#vg'),  $button_el: $('#vgconfig') .find('.hiding')})

    ip       = React.render(React.createElement(NewTextCLASS((data) -> data?.IP       )), $('#ip'      )   .get(0))
    hostname = React.render(React.createElement(NewTextCLASS((data) -> data?.Hostname )), $('#hostname')   .get(0))
    uptime   = React.render(React.createElement(NewTextCLASS((data) -> data?.Uptime   )), $('#uptime'  )   .get(0))
    la       = React.render(React.createElement(NewTextCLASS((data) -> data?.LA       )), $('#la'      )   .get(0))

    iftitle  = React.render(React.createElement(NewTextCLASS((data) -> data?.Client?.TabTitleIF)), $('header a[href="#if"]').get(0))
    dftitle  = React.render(React.createElement(NewTextCLASS((data) -> data?.Client?.TabTitleDF)), $('header a[href="#df"]').get(0))

    psplus   = React.render(React.createElement(NewTextCLASS((data) -> data?.Client?.PSplusText)), $('label.more[href="#psmore"]').get(0))
    psmore   = ButtonClass.component({Ksig: 'MorePsignal', Vsig: true,  Khide: 'HidePS', Kable: 'PSnotExpandable',  $button_el: $('label.more[href="#psmore"]')})
    psless   = ButtonClass.component({Ksig: 'MorePsignal', Vsig: false, Khide: 'HidePS', Kable: 'PSnotDecreasable', $button_el: $('label.less[href="#psless"]')})

    hideswap = ButtonClass.component({Khide: 'HideRAM', Ksend: 'HideSWAP', $button_el: $('label[href="#hideswap"]')})

    expandif = ButtonClass.component({Khide: 'HideIF',  Ksend: 'ExpandIF',  Ktext: 'ExpandtextIF',  Kable: 'ExpandableIF',  $button_el: $('label[href="#if"]')})
    expandcpu= ButtonClass.component({Khide: 'HideCPU', Ksend: 'ExpandCPU', Ktext: 'ExpandtextCPU', Kable: 'ExpandableCPU', $button_el: $('label[href="#cpu"]')})
    expanddf = ButtonClass.component({Khide: 'HideDF',  Ksend: 'ExpandDF',  Ktext: 'ExpandtextDF',  Kalbe: 'ExpandableDF',  $button_el: $('label[href="#df"]')})

    # NB buttons and collapses selected by class
    tabsif = TabsClass.component({Khide: 'HideIF', Ksend: 'TabIF', $collapse_el: $('.if-tab'), $button_el: $('.if-switch'), $hidebutton_el: $('#ifconfig').find('.hiding')})
    tabsdf = TabsClass.component({Khide: 'HideDF', Ksend: 'TabDF', $collapse_el: $('.df-tab'), $button_el: $('.df-switch'), $hidebutton_el: $('#dfconfig').find('.hiding')})

    refresh_mem = RefreshInputClass.component({K: 'RefreshMEM', Kerror: 'RefreshErrorMEM', Ksig: 'RefreshSignalMEM', sel: $('#memconfig')})
    refresh_if  = RefreshInputClass.component({K: 'RefreshIF',  Kerror: 'RefreshErrorIF',  Ksig: 'RefreshSignalIF',  sel: $('#ifconfig')})
    refresh_cpu = RefreshInputClass.component({K: 'RefreshCPU', Kerror: 'RefreshErrorCPU', Ksig: 'RefreshSignalCPU', sel: $('#cpuconfig')})
    refresh_df  = RefreshInputClass.component({K: 'RefreshDF',  Kerror: 'RefreshErrorDF',  Ksig: 'RefreshSignalDF',  sel: $('#dfconfig')})
    refresh_ps  = RefreshInputClass.component({K: 'RefreshPS',  Kerror: 'RefreshErrorPS',  Ksig: 'RefreshSignalPS',  sel: $('#psconfig')})
    refresh_vg  = RefreshInputClass.component({K: 'RefreshVG',  Kerror: 'RefreshErrorVG',  Ksig: 'RefreshSignalVG',  sel: $('#vgconfig')})

    memtable  = React.render(React.createElement(MEMtableCLASS),  document.getElementById('mem'       +'-'+ 'table'))
    pstable   = React.render(React.createElement(PStableCLASS),   document.getElementById('ps'        +'-'+ 'table'))
    dfbytes   = React.render(React.createElement(DFbytesCLASS),   document.getElementById('dfbytes'   +'-'+ 'table'))
    dfinodes  = React.render(React.createElement(DFinodesCLASS),  document.getElementById('dfinodes'  +'-'+ 'table'))
    cputable  = React.render(React.createElement(CPUtableCLASS),  document.getElementById('cpu'       +'-'+ 'table'))
    ifbytes   = React.render(React.createElement(IFbytesCLASS),   document.getElementById('ifbytes'   +'-'+ 'table'))
    iferrors  = React.render(React.createElement(IFerrorsCLASS),  document.getElementById('iferrors'  +'-'+ 'table'))
    ifpackets = React.render(React.createElement(IFpacketsCLASS), document.getElementById('ifpackets' +'-'+ 'table'))
    vgtable   = React.render(React.createElement(VGtableCLASS),   document.getElementById('vg'        +'-'+ 'table'))
    # coffeelint: enable=max_line_length

    # alertComp = React.render(AlertClass({
    #         $collapse_el: $('#alert-parent')
    #         }), document.getElementById('alert-message'))

    onmessage = (event) ->
      data = JSON.parse(event.data)
      return if !data?

      # alertComp.setState(alertComp.newstate(data))
      # if alertComp.show()
      #         return

      console.log('DEBUG ERROR',
        data.Client.DebugError) if data.Client?.DebugError?
      if data.Reload? and data.Reload
        window.setTimeout((() -> location.reload(true)), 5000)
        window.setTimeout(updates.close, 2000)
        console.log('in 5s: location.reload(true)')
        console.log('in 2s: updates.close()')
        return

      setState(pstable,  {PStable:  data.PStable,  Links: data.Links})
      setState(dfbytes,  {DFbytes:  data.DFbytes,  Links: data.Links})
      setState(dfinodes, {DFinodes: data.DFinodes, Links: data.Links})

      setState(hideconfigmem, hideconfigmem.reduce(data))
      setState(hideconfigif,  hideconfigif .reduce(data))
      setState(hideconfigcpu, hideconfigcpu.reduce(data))
      setState(hideconfigdf,  hideconfigdf .reduce(data))
      setState(hideconfigps,  hideconfigps .reduce(data))
      setState(hideconfigvg,  hideconfigvg .reduce(data))

      setState(hideram,       hideram      .reduce(data))
      setState(hidecpu,       hidecpu      .reduce(data))
      setState(hideps,        hideps       .reduce(data))
      setState(hidevg,        hidevg       .reduce(data))

      setState(ip,        ip      .newstate(data))
      setState(hostname,  hostname.newstate(data))
      setState(uptime,    uptime  .newstate(data))
      setState(la,        la      .newstate(data))

      setState(iftitle,   iftitle .newstate(data))
      setState(dftitle,   dftitle .newstate(data))

      setState(psplus,    psplus  .newstate(data))
      setState(psmore,    psmore  .reduce(data))
      setState(psless,    psless  .reduce(data))

      setState(hideswap,  hideswap.reduce(data))

      setState(expandif,  expandif.reduce(data))
      setState(expandcpu, expandcpu.reduce(data))
      setState(expanddf,  expanddf.reduce(data))

      setState(tabsif,    tabsif.reduce(data))
      setState(tabsdf,    tabsdf.reduce(data))

      setState(refresh_mem, refresh_mem.reduce(data))
      setState(refresh_if,  refresh_if .reduce(data))
      setState(refresh_cpu, refresh_cpu.reduce(data))
      setState(refresh_df,  refresh_df .reduce(data))
      setState(refresh_ps,  refresh_ps .reduce(data))
      setState(refresh_vg,  refresh_vg .reduce(data))

      setState(memtable,  data.MEM)
      setState(cputable,  data.CPU)
      setState(ifbytes,   data.IFbytes)
      setState(iferrors,  data.IFerrors)
      setState(ifpackets, data.IFpackets)
      setState(vgtable, {
        VagrantMachines: data.VagrantMachines,
        VagrantError:    data.VagrantError,
        VagrantErrord:   data.VagrantErrord
      })

      console.log(JSON.stringify(data.Client), 'recvClient') if data.Client?

      # update the tooltips
      $('span .tooltipable')    .popover({trigger: 'hover focus'})
      $('span .tooltipabledots').popover() # the clickable dots
      return

    updates = newwebsocket(onmessage)
  # updates = neweventsource(onmessage)
    return # end of `update'

  require ['domReady', 'jquery', 'bootstrap', 'headroom'], (domReady, $) ->
    # neither bootstrap nor headroom export anything
    domReady () ->
      (new window.Headroom(document.querySelector('nav'), {
        offset: 20 # ~padding-top of a container row
      })).init()

      $('.collapse').collapse({toggle: false}) # init collapsable objects

      $('span .tooltipable')      .popover({trigger: 'hover focus'})
      $('span .tooltipabledots')  .popover() # the clickable dots
      $('[data-toggle="popover"]').popover() # should be just #hostname
      $('#la')                    .popover({
        trigger: 'hover focus',
        placement: 'right',
        # NOT placement: 'auto right' until #la is the last element for parent
        html: true, content: () -> $('#uptime-parent').html()
      })

      $('body').on('click', (e) -> # hide the popovers on click outside
        $('span .tooltipabledots').each(() ->
          # the 'is' for buttons that trigger popups
          # the 'has' for icons within a button that triggers a popup
          $(this).popover('hide') if (!$(this).is(e.target) and
            $(this).has(e.target).length == 0 and
            $('.popover').has(e.target).length == 0)
          return)
        return)

      # referencing upper-scope `update'
      update() unless (42 for param in location.search.substr(1).split(
        '&') when (param.split('=')[0] == 'still')).length

      return # return from domReady
    return # end of sub`require'
  return # end of main `require'

# Local variables:
# coffee-tab-width: 2
# End:
