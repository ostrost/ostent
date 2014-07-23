
@newwebsocket = (onmessage) ->
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
                hostport = window.location.hostname + (if location.port then ':' + location.port else '')
                conn = new WebSocket('ws://' + hostport + '/ws')
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
                return ifbytes_table(Data, (ifbytes_rows(Data, $if) for $if in Data?.IFbytes?.List ? []))

@IFerrorsCLASS = React.createClass
        getInitialState: () -> Data.IFerrors # a global Data
        render: () ->
                Data = {IFerrors: @state}
                return iferrors_table(Data, (iferrors_rows(Data, $if) for $if in Data?.IFerrors?.List ? []))

@IFpacketsCLASS = React.createClass
        getInitialState: () -> Data.IFpackets # a global Data
        render: () ->
                Data = {IFpackets: @state}
                return ifpackets_table(Data, (ifpackets_rows(Data, $if) for $if in Data?.IFpackets?.List ? []))

@DFbytesCLASS = React.createClass
        getInitialState: () -> {DFlinks: Data.DFlinks, DFbytes: Data.DFbytes} # a global Data
        render: () ->
                Data = @state
                return dfbytes_table(Data, (dfbytes_rows(Data, $disk) for $disk in Data?.DFbytes?.List ? []))

@DFinodesCLASS = React.createClass
        getInitialState: () -> {DFlinks: Data.DFlinks, DFinodes: Data.DFinodes} # a global Data
        render: () ->
                Data = @state
                return dfinodes_table(Data, (dfinodes_rows(Data, $disk) for $disk in Data?.DFinodes?.List ? []))

@MEMtableCLASS = React.createClass
        getInitialState: () -> Data.MEM # a global Data
        render: () ->
                Data = {MEM: @state}
                return mem_table(Data, (mem_rows(Data, $mem) for $mem in Data?.MEM?.List ? []))

@CPUtableCLASS = React.createClass
        getInitialState: () -> Data.CPU # a global Data
        render: () ->
                Data = {CPU: @state}
                return cpu_table(Data, (cpu_rows(Data, $core) for $core in Data?.CPU?.List ? []))

@PStableCLASS = React.createClass
        getInitialState: () -> {PStable: Data.PStable, PSlinks: Data.PSlinks} # a global Data
        render: () ->
                Data = @state
                return ps_table(Data, (ps_rows(Data, $proc) for $proc in Data?.PStable?.List ? []))

@VGtableCLASS = React.createClass
        getInitialState: () -> { # a global Data:
                VagrantMachines: Data.VagrantMachines
                VagrantError:    Data.VagrantError
                VagrantErrord:   Data.VagrantErrord
        }
        render: () ->
                Data = @state
                if Data?.VagrantErrord? and Data.VagrantErrord
                        rows = [vagrant_error(Data)]
                else
                        rows = (vagrant_rows(Data, $machine) for $machine in Data?.VagrantMachines?.List ? [])
                return vagrant_table(Data, rows)

@addNoscript = ($) -> $.append('<noscript />').find('noscript').get(0)

@HideClass = React.createClass
        statics: component: (opt) -> React.renderComponent(HideClass(opt), addNoscript(opt.$button_el))

        reduce: (data) ->
                if data?.Client?
                        value = data.Client[@props.key]
                        return {Hide: value} if value isnt undefined
        getInitialState: () -> @reduce(Data) # a global Data
        componentDidMount: () -> @props.$button_el.click(@click)
        render: () ->
                @props.$collapse_el.collapse(if @state.Hide then 'hide' else 'show')
                buttonactive =  @state.Hide
                buttonactive = !@state.Hide if @props.reverseActive? and @props.reverseActive
                @props.$button_el[if buttonactive then 'addClass' else 'removeClass']('active')
                return null
        click: (e) ->
                (S = {})[@props.key] = !@state.Hide
                websocket.sendClient(S)
                e.stopPropagation() # preserves checkbox/radio
                e.preventDefault()  # checked/selected state
                return undefined

@ButtonClass = React.createClass
        statics: component: (opt) -> React.renderComponent(ButtonClass(opt), addNoscript(opt.$button_el))

        reduce: (data) ->
                if data?.Client?
                        S = {}
                        S.Hide = data.Client[@props.Khide] if                   data.Client[@props.Khide] isnt undefined # Khide is a required prop
                        S.Able = data.Client[@props.Kable] if @props.Kable? and data.Client[@props.Kable] isnt undefined
                        S.Send = data.Client[@props.Ksend] if @props.Ksend? and data.Client[@props.Ksend] isnt undefined
                        S.Text = data.Client[@props.Ktext] if @props.Ktext? and data.Client[@props.Ktext] isnt undefined
                        return S
        getInitialState: () -> @reduce(Data) # a global Data
        componentDidMount: () -> @props.$button_el.click(@click)
        render: () ->
                if @props.Kable
                        able = @state.Able
                        able = !able if not (@props.Kable.indexOf('not') > -1) # That's a hack
                        @props.$button_el.prop('disabled', able)
                        @props.$button_el[if able then 'addClass' else 'removeClass']('disabled')
                @props.$button_el[if @state.Send then 'addClass' else 'removeClass']('active') if @props.Ksend?
                @props.$button_el.text(@state.Text) if @props.Ktext?
                return null
        click: (e) ->
                S = {}
                S[@props.Khide] = !@state.Hide if @state.Hide?  and @state.Hide # if the panel was hidden
                S[@props.Ksend] = !@state.Send if @props.Ksend? and @state.Send? # Q is a @state.Send? check excessive?
                S[@props.Ksig]  =  @props.Vsig if @props.Ksig?
                websocket.sendClient(S)
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

@update = (currentClient, model) ->
        return if (42 for param in location.search.substr(1).split('&') when param.split('=')[0] == 'still').length

        hideconfigmem = HideClass.component({key: 'HideconfigMEM', $collapse_el: $('#memconfig'), $button_el: $('header a[href="#mem"]'), reverseActive: true})
        hideconfigif  = HideClass.component({key: 'HideconfigIF',  $collapse_el: $('#ifconfig'),  $button_el: $('header a[href="#if"]'),  reverseActive: true})
        hideconfigcpu = HideClass.component({key: 'HideconfigCPU', $collapse_el: $('#cpuconfig'), $button_el: $('header a[href="#cpu"]'), reverseActive: true})
        hideconfigdf  = HideClass.component({key: 'HideconfigDF',  $collapse_el: $('#dfconfig'),  $button_el: $('header a[href="#df"]'),  reverseActive: true})
        hideconfigps  = HideClass.component({key: 'HideconfigPS',  $collapse_el: $('#psconfig'),  $button_el: $('header a[href="#ps"]'),  reverseActive: true})
        hideconfigvg  = HideClass.component({key: 'HideconfigVG',  $collapse_el: $('#vgconfig'),  $button_el: $('header a[href="#vg"]'),  reverseActive: true})

        hidemem = HideClass.component({key: 'HideMEM', $collapse_el: $('#mem'), $button_el: $('#memconfig').find('.hiding')})
        hidecpu = HideClass.component({key: 'HideCPU', $collapse_el: $('#cpu'), $button_el: $('#cpuconfig').find('.hiding')})
        hideps  = HideClass.component({key: 'HidePS',  $collapse_el: $('#ps'),  $button_el: $('#psconfig') .find('.hiding')})
        hidevg  = HideClass.component({key: 'HideVG',  $collapse_el: $('#vg'),  $button_el: $('#vgconfig') .find('.hiding')})

        ip       = React.renderComponent(NewTextCLASS((data) -> data?.Generic?.IP       )(), $('#generic-ip'      )   .get(0))
        hostname = React.renderComponent(NewTextCLASS((data) -> data?.Generic?.Hostname )(), $('#generic-hostname')   .get(0))
        uptime   = React.renderComponent(NewTextCLASS((data) -> data?.Generic?.Uptime   )(), $('#generic-uptime'  )   .get(0))
        la       = React.renderComponent(NewTextCLASS((data) -> data?.Generic?.LA       )(), $('#generic-la'      )   .get(0))

        iftitle  = React.renderComponent(NewTextCLASS((data) -> data?.Client?.TabTitleIF)(), $('header a[href="#if"]').get(0))
        dftitle  = React.renderComponent(NewTextCLASS((data) -> data?.Client?.TabTitleDF)(), $('header a[href="#df"]').get(0))

        psplus   = React.renderComponent(NewTextCLASS((data) -> data?.Client?.PSplusText)(), $('label.more[href="#psmore"]').get(0))
        psmore   = ButtonClass.component({Ksig: 'MorePsignal', Vsig: true,  Khide: 'HidePS', Kable: 'PSnotExpandable',  $button_el: $('label.more[href="#psmore"]')})
        psless   = ButtonClass.component({Ksig: 'MorePsignal', Vsig: false, Khide: 'HidePS', Kable: 'PSnotDecreasable', $button_el: $('label.less[href="#psless"]')})

        hideswap = ButtonClass.component({Khide: 'HideMEM', Ksend: 'HideSWAP', $button_el: $('label[href="#hideswap"]')})

        expandif = ButtonClass.component({Khide: 'HideIF',  Ksend: 'ExpandIF',  Ktext: 'ExpandtextIF',  Kable: 'ExpandableIF',  $button_el: $('label[href="#if"]')})
        expandcpu= ButtonClass.component({Khide: 'HideCPU', Ksend: 'ExpandCPU', Ktext: 'ExpandtextCPU', Kable: 'ExpandableCPU', $button_el: $('label[href="#cpu"]')})
        expanddf = ButtonClass.component({Khide: 'HideDF',  Ksend: 'ExpandDF',  Ktext: 'ExpandtextDF',  Kalbe: 'ExpandableDF',  $button_el: $('label[href="#df"]')})

        memtable  = React.renderComponent(MEMtableCLASS(),  document.getElementById('mem'       +'-'+ 'table'))
        pstable   = React.renderComponent(PStableCLASS(),   document.getElementById('ps'        +'-'+ 'table'))
        dfbytes   = React.renderComponent(DFbytesCLASS(),   document.getElementById('dfbytes'   +'-'+ 'table'))
        dfinodes  = React.renderComponent(DFinodesCLASS(),  document.getElementById('dfinodes'  +'-'+ 'table'))
        cputable  = React.renderComponent(CPUtableCLASS(),  document.getElementById('cpu'       +'-'+ 'table'))
        ifbytes   = React.renderComponent(IFbytesCLASS(),   document.getElementById('ifbytes'   +'-'+ 'table'))
        iferrors  = React.renderComponent(IFerrorsCLASS(),  document.getElementById('iferrors'  +'-'+ 'table'))
        ifpackets = React.renderComponent(IFpacketsCLASS(), document.getElementById('ifpackets' +'-'+ 'table'))
        vgtable   = React.renderComponent(VGtableCLASS(),   document.getElementById('vg'        +'-'+ 'table'))

        onmessage = (event) ->
                data = JSON.parse(event.data)
                return if !data?

                console.log('DEBUG ERROR', data.Client.DebugError) if data.Client?.DebugError?
                if data.Reload? and data.Reload
                        window.setTimeout((() -> location.reload(true)), 5000)
                        window.setTimeout(websocket.close, 2000)
                        console.log('in 5s: location.reload(true)')
                        console.log('in 2s: websocket.close()')
                        return

                setState(pstable,  {PStable:  data.PStable,  PSlinks: data.PSlinks})
                setState(dfbytes,  {DFbytes:  data.DFbytes,  DFlinks: data.DFlinks})
                setState(dfinodes, {DFinodes: data.DFinodes, DFlinks: data.DFlinks})

                setState(hideconfigmem, hideconfigmem.reduce(data))
                setState(hideconfigif,  hideconfigif .reduce(data))
                setState(hideconfigcpu, hideconfigcpu.reduce(data))
                setState(hideconfigdf,  hideconfigdf .reduce(data))
                setState(hideconfigps,  hideconfigps .reduce(data))
                setState(hideconfigvg,  hideconfigvg .reduce(data))

                setState(hidemem,       hidemem      .reduce(data))
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

                currentClient = React.addons.update(currentClient, {$merge: data.Client}) if data.Client?
                data.Client = currentClient
                model.set(Model.attributes(data))

                # update the tooltips
                $('span .tooltipable')    .popover({trigger: 'hover focus'})
                $('span .tooltipabledots').popover() # the clickable dots
                return

        @websocket = newwebsocket(onmessage)
        return

@Model = Backbone.Model.extend({})
@Model.attributes = (data) ->
        return data.Client if !data.Generic?
        return data.Generic if !data.Client?
        return React.addons.update(data.Generic, {$merge: data.Client})

@View = Backbone.View.extend({
        initialize: () ->
              # @listentext('IP',       $('#generic-ip'))
              # @listentext('Hostname', $('#generic-hostname'))
              # @listentext('Uptime',   $('#uptime #generic-uptime'))
              # @listentext('LA',       $('#generic-la'))

              # $hswapb = $('label[href="#showswap"]')
              # @listenactivate('HideSWAP', $hswapb, true)

              # $section_* were here

              # the $config_{mem,cpu,ps,vg} were here
                $config_if  = $('#ifconfig')
                $config_df  = $('#dfconfig')

              # the $hidden_{mem,cpu,ps,vg} were here
                $hidden_if  = $config_if .find('.hiding')
                $hidden_df  = $config_df .find('.hiding')

              # the 4th argument to @listenhide used to be optional and `false' by default
              # @listenhide('HideMEM', $section_mem, $hidden_mem, false)
              # @listenhide('HideCPU', $section_cpu, $hidden_cpu, false)
              # @listenhide('HidePS',  $section_ps,  $hidden_ps,  false) # $section_ps used to be $('#ps')
              # @listenhide('HideVG',  $section_vg,  $hidden_vg,  false) # $section_vg used to be $('#vg')

              # $header_mem = $('header a[href="'+ $section_mem.selector + '"]')
              # $header_if  = $('header a[href="'+ $section_if .selector + '"]')
              # $header_cpu = $('header a[href="'+ $section_cpu.selector + '"]')
              # $header_df  = $('header a[href="'+ $section_df .selector + '"]')

              # $header_ps  = $('header a[href="#ps"]') # remember $section_ps
              # $header_vg  = $('header a[href="#vg"]') # remember $section_vg

              # @listentext('TabTitleIF', $header_if)
              # @listentext('TabTitleDF', $header_df)

              # @listenhide('HideconfigMEM', $config_mem, $header_mem) #, true)
              # @listenhide('HideconfigIF',  $config_if,  $header_if)  #, true)
              # @listenhide('HideconfigCPU', $config_cpu, $header_cpu) #, true)
              # @listenhide('HideconfigDF',  $config_df,  $header_df)  #, true)
              # @listenhide('HideconfigPS',  $config_ps,  $header_ps)  #, true)
              # @listenhide('HideconfigVG',  $config_vg,  $header_vg)  #, true)

                # NB by class
                $tab_if    = $('.if-switch')
                $tab_df    = $('.df-switch')
                $panels_if = $('.if-tab')
                $panels_df = $('.df-tab')

                @listenTo(@model, 'change:HideIF', @change_collapsetabfunc('HideIF', 'TabIF', $panels_if, $tab_if, $hidden_if))
                @listenTo(@model, 'change:HideDF', @change_collapsetabfunc('HideDF', 'TabDF', $panels_df, $tab_df, $hidden_df))
                @listenTo(@model, 'change:TabIF',  @change_collapsetabfunc('HideIF', 'TabIF', $panels_if, $tab_if, $hidden_if))
                @listenTo(@model, 'change:TabDF',  @change_collapsetabfunc('HideDF', 'TabDF', $panels_df, $tab_df, $hidden_df))

              # $psmore = $('label.more[href="#psmore"]')
              # $psless = $('label.less[href="#psless"]')
              # @listentext('PSplusText',         $psmore)
              # @listenenable('PSnotExpandable',  $psmore)
              # @listenenable('PSnotDecreasable', $psless)

              # $config_{if,df} defined previously
                $config_mem = $('#memconfig')
                $config_cpu = $('#cpuconfig')
                $config_ps  = $('#psconfig')
                $config_vg  = $('#vgconfig')

                @listenrefresherror('RefreshErrorMEM', $config_mem.find('.refresh-group'))
                @listenrefresherror('RefreshErrorIF',  $config_if .find('.refresh-group'))
                @listenrefresherror('RefreshErrorCPU', $config_cpu.find('.refresh-group'))
                @listenrefresherror('RefreshErrorDF',  $config_df .find('.refresh-group'))
                @listenrefresherror('RefreshErrorPS',  $config_ps .find('.refresh-group'))
                @listenrefresherror('RefreshErrorVG',  $config_vg .find('.refresh-group'))

                @listenrefreshvalue('RefreshMEM',      $config_mem.find('.refresh-input'))
                @listenrefreshvalue('RefreshIF',       $config_if .find('.refresh-input'))
                @listenrefreshvalue('RefreshCPU',      $config_cpu.find('.refresh-input'))
                @listenrefreshvalue('RefreshDF',       $config_df .find('.refresh-input'))
                @listenrefreshvalue('RefreshPS',       $config_ps .find('.refresh-input'))
                @listenrefreshvalue('RefreshVG',       $config_vg .find('.refresh-input'))

                # var B = _.bind(function(c) { return _.bind(c, this); }, this);
                # B = _.bind(((c) -> _.bind(c, @)), @)
                B = (c) -> c # _.bind(((c) -> _.bind(c, @)), @)

                # B((e) -> return e)
                # click_expandfunc: (H, H2) -> (e) ->
                #     $b.click( B(@click_expandfunc(E, H)) )


              # $section_mem = $('#mem') # other $section_* were here
              # $section_if  = $('#if')
              # $section_cpu = $('#cpu')
              # $section_df  = $('#df')

                expandable_sections = [
                  # [$section_if,  'ExpandIF',  'HideIF',  'ExpandableIF',  'ExpandtextIF'],
                  # [$section_cpu, 'ExpandCPU', 'HideCPU', 'ExpandableCPU', 'ExpandtextCPU'],
                  # [$section_df,  'ExpandDF',  'HideDF',  'ExpandableDF',  'ExpandtextDF']
                ]
                doexpandable = (sections) =>
                    S = sections[0]
                    E = sections[1]
                    H = sections[2]
                    L = sections[3]
                    T = sections[4]
                    $b = $('label[href="'+ S.selector + '"]')

                    @listentext(T, $b) # Expandtext*, the text of label[href="#{if,cpu,df}"]
                    @listenenable(L, $b, true)
                    @listenactivate(E, $b)
                    $b.click( B(@click_expandfunc(E, H)) )
                    return

              # doexpandable(sections) for sections in expandable_sections

              # $hswapb    .click( B(@click_expandfunc('HideSWAP', 'HideMEM')) )
                $tab_if    .click( B(@click_tabfunc('TabIF', 'HideIF')) )
                $tab_df    .click( B(@click_tabfunc('TabDF', 'HideDF')) )

              # $header_mem.click( B(@click_expandfunc('HideconfigMEM')) )
              # $header_if .click( B(@click_expandfunc('HideconfigIF' )) )
              # $header_cpu.click( B(@click_expandfunc('HideconfigCPU')) )
              # $header_df .click( B(@click_expandfunc('HideconfigDF' )) )
              # $header_ps .click( B(@click_expandfunc('HideconfigPS' )) )
              # $header_vg .click( B(@click_expandfunc('HideconfigVG' )) )

              # $hidden_mem.click( B(@click_expandfunc('HideMEM')) )
                $hidden_if .click( B(@click_expandfunc('HideIF' )) )
              # $hidden_cpu.click( B(@click_expandfunc('HideCPU')) )
                $hidden_df .click( B(@click_expandfunc('HideDF' )) )
              # $hidden_ps .click( B(@click_expandfunc('HidePS' )) )
              # $hidden_vg .click( B(@click_expandfunc('HideVG' )) )

              # $psmore    .click( B(@click_psignalfunc('HidePS', true )) )
              # $psless    .click( B(@click_psignalfunc('HidePS', false)) )

                $config_mem.find('.refresh-input').on('input', B(@submit_rsignalfunc('RefreshSignalMEM')) )
                $config_if .find('.refresh-input').on('input', B(@submit_rsignalfunc('RefreshSignalIF' )) )
                $config_cpu.find('.refresh-input').on('input', B(@submit_rsignalfunc('RefreshSignalCPU')) )
                $config_df .find('.refresh-input').on('input', B(@submit_rsignalfunc('RefreshSignalDF' )) )
                $config_ps .find('.refresh-input').on('input', B(@submit_rsignalfunc('RefreshSignalPS' )) )
                $config_vg .find('.refresh-input').on('input', B(@submit_rsignalfunc('RefreshSignalVG' )) )
                return

        submit_rsignalfunc: (R) ->
                return (e) ->
                        (S = {})[R] = $(e.target).val()
                        websocket.sendClient(S)
                        return

        listentext: (K, $el) -> @listenTo(@model, 'change:'+ K, @_text(K, $el))
#       listenHTML: (K, $el) -> @listenTo(@model, 'change:'+ K, @_HTML(K, $el))
        _text:      (K, $el) -> () -> $el.text(@model.attributes[K])
#       _HTML:      (K, $el) -> () -> $el.html(@model.attributes[K])

        listenrefresherror: (E, $el) ->
                @listenTo(@model, 'change:'+ E, () ->
                        $el[if @model.attributes[E] then 'addClass' else 'removeClass']('has-warning'))

        listenrefreshvalue: (E, $el) ->
                @listenTo(@model, 'change:'+ E, () ->
                        $el.prop('value', @model.attributes[E]))

        listenenable: (K, $el, reverse) ->
                @listenTo(@model, 'change:'+ K, () ->
                    V = @model.attributes[K]? and @model.attributes[K]
                    V = !V if reverse? && reverse
                    $el.prop('disabled', V)
                    $el[if V then 'addClass' else 'removeClass']('disabled')
                    return)

        listenactivate: (K, $el, reverse) ->
                @listenTo(@model, 'change:'+ K, () ->
                    V = @model.attributes[K]? and @model.attributes[K]
                    V = !V if !reverse? && reverse
                    $el[if V then 'addClass' else 'removeClass']('active')
                    return)

        listenhide: (H, $el, $button_el) ->
                # the 4th argument used to be `reverse'
                @listenTo(@model, 'change:'+ H, () ->
                        V = @model.attributes[H]? and @model.attributes[H]
                        $el.collapse(if V then 'hide' else 'show') # do what change_collapsefunc does

                        V = !V # if !reverse? && reverse
                        $button_el[if V then 'addClass' else 'removeClass']('active') # do what listenactivate does
                        return)

        change_collapsefunc: (H, $el) -> () -> $el.collapse(if @model.attributes[H] then 'hide' else 'show')

        change_collapsetabfunc: (H, T, $el, $tabel, $buttonel) -> () ->
                A = @model.attributes
                if A[H] # hiding all
                        $el.collapse('hide') # do what change_collapsefunc does
                        $buttonel.addClass('active')
                        return
                $buttonel.removeClass('active')

                # $el is $('.if-tab')
                # $tabel is $('.if-switch')

                curtabid = +A[T] # MUST be an int
                nots = $el.not('[data-tabid="'+ curtabid + '"]')

                $(el).collapse('hide') for el in nots
                $($el.not(nots)).collapse('show')

                activeClass = (el) ->
                        xel = $(el)
                        tabid_attr = +xel.attr('data-tabid') # an int
                        xel[if tabid_attr == curtabid then 'addClass' else 'removeClass']('active')
                        return
                activeClass(el) for el in $tabel
                return

        click_expandfunc: (H, H2) -> (e) =>
                A = @model.attributes
                (S = {})[H] = !A[H]
                S[H2] = !A[H2] if H2? and A[H2] # if was hidden
                websocket.sendClient(S)
                e.preventDefault()
                e.stopPropagation() # don't change checkbox/radio state
                return undefined

        click_tabfunc: (T, H) -> (e) =>
                newtabid = +$( $(e.target).attr('href') ).attr('data-tabid') # THIS. +string makes an int
                (S = {})[T] = newtabid
                V = @model.attributes[H]
                S[H] = !V if V # if was hidden
                websocket.sendClient(S)
                e.preventDefault()
                e.stopPropagation() # don't change checkbox/radio state
                return undefined

        click_psignalfunc: (H, v) -> (e) =>
                S = {MorePsignal: v}
                V = @model.attributes[H]
                S[H] = !V if V # if was hidden
                websocket.sendClient(S)
                e.preventDefault()
                e.stopPropagation() # don't change checkbox/radio state
                return undefined
})

@ready = () ->
        (new Headroom(document.querySelector('nav'), {
                offset: 71 - 51
                # "relative" padding-top of the toprow
                # 71 is the absolute padding-top of the toprow
                # 51 is the height of the nav (50 +1px bottom border)
        })).init()

        $('.collapse').collapse({toggle: false}) # init collapsable objects

        $('span .tooltipable')      .popover({trigger: 'hover focus'})
        $('span .tooltipabledots')  .popover() # the clickable dots
        $('[data-toggle="popover"]').popover() # should be just #generic-hostname
        $('#generic-la')            .popover({
                trigger: 'hover focus',
                placement: 'right', # not 'auto right' until #generic-la is the last element for it's parent
                html: true, content: () -> $('#uptime').html()
        })

        $('body').on('click', (e) -> # hide the popovers on click outside
                $('span .tooltipabledots').each(() ->
                        # the 'is' for buttons that trigger popups
                        # the 'has' for icons within a button that triggers a popup
                        $(this).popover('hide') if !$(this).is(e.target) and $(this).has(e.target).length == 0 && $('.popover').has(e.target).length == 0
                        return)
                return)

        model = new Model(Model.attributes(Data))
        new View({model: model})

        update(Data.Client, model)
        return

