var React       = require('react'),
    ReactDOM    = require('react-dom'),
    ReconnectWS = require('reconnectingWebsocket'),
    jsdefines   = require('./jsdefines.jsx');

/*
function neweventsource(onmessage) {
  var init, conn = null;
  var sendSearch = function(search) {
    location.search = search;
    console.log('new search', search);
    conn.close();
    window.setTimeout(init, 100);
  };
  var sendLocationSearch = function() { sendSearch(location.search); };
  init = function() {
    conn = new EventSource('/index.sse' + location.search);
    conn.onopen = function() {
      console.log('sse opened');
      window.addEventListener('popstate', sendLocationSearch);
    };
    conn.onclose = function() {
      console.log('sse closed');
      window.removeEventListener('popstate', sendLocationSearch);
    };
    conn.onerror = function() {
      console.log('sse errord');
      window.removeEventListener('popstate', sendLocationSearch);
    };
    conn.onmessage = onmessage;
  };
  init();
  return {
    sendSearch: sendSearch,
    close: function() { return conn.close(); }
  };
}; // */

function main(data) {
  if (data.params.Still.Absolute != 0) {
    return;
  }

  var els = [];
  for (var i = 0, sel = document.querySelectorAll('.updates'); i < sel.length; i++) {
    var cl = jsdefines[sel[i].getAttribute('data-define')];
    els.push(ReactDOM.render(React.createElement(cl), sel[i]));
  }

  var hostport = location.hostname + (location.port ? ':' + location.port : '');
  var ws = new ReconnectWS('ws://' + hostport + '/index.ws');

  // sendSearch is also referenced as window.updates.sendSearch
  ws.sendSearch = function(search) {
    console.log('ws send', search);
    ws.send(JSON.stringify({Search: search}));
  };
  ws.sendLocationSearch = function() { ws.sendSearch(location.search); };

  ws.onclose = function() {
    console.log('ws closed');
    window.removeEventListener('popstate', ws.sendLocationSearch);
  };
  ws.onopen = function() {
    console.log('ws opened');
    ws.sendLocationSearch();
    window.addEventListener('popstate', ws.sendLocationSearch);
  };
  ws.onmessage = function(event) { // the onmessage
    var data = JSON.parse(event.data);
    if (data == null) {
      return;
    }
    if ((data.Reload != null) && data.Reload) {
      window.setTimeout((function() { location.reload(true); }), 5000);
      window.setTimeout(ws.close, 2000);
      console.log('in 5s: location.reload(true)');
      console.log('in 2s: ws.close()');
      return;
    }
    if (data.Error != null) {
      console.log('Error', data.Error);
      return;
    }
    for (var i = 0; i < els.length; i++) {
      els[i].NewState(data);
    }
  };

  window.updates = ws; // neweventsource(onmessage);
}

main(Data); // global Data

// Local variables:
// js-indent-level: 2
// js2-basic-offset: 2
// End:
