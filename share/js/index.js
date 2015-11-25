var $         = require('jquery'),
    React     = require('react'),
    ReactDOM  = require('react-dom'),
    jsdefines = require('./jsdefines.js');

function neweventsource(onmessage) {
  var conn = null;
  var init;
  var sendSearch = function(search) {
    // conn = new EventSource('/index.sse' + search)
    // location.search = search
    console.log('SEARCH', search);
    if (true) { // conn !== null
      conn.close();
    }
    return window.setTimeout(init, 1000);
  };
  init = function() {
    conn = new EventSource('/index.sse' + location.search);
    conn.onopen = function() {
      $(window).bind('popstate', (function() { sendSearch(location.search); }));
    };
    var again = function(e) {
      if (!e.wasClean) {
        window.setTimeout(init, 5000);
      }
    };
    conn.onclose = function() { console.log('sse closed (should recover)'); }; // again;
    conn.onerror = function() { console.log('sse errord (should recover)'); }; // again;
    conn.onmessage = onmessage;
  };
  init();
  return {
    sendSearch: sendSearch,
    close: function() { return conn.close(); }
  };
};

function newwebsocket(onmessage) {
  var conn = null;
  var init;
  var sendSearch = function(search) {
    console.log('Search', search);
    // 0 conn.CONNECTING
    // 1 conn.OPEN
    // 2 conn.CLOSING
    // 3 conn.CLOSED
    if (conn == null ||
        conn.readyState === conn.CLOSING ||
        conn.readyState === conn.CLOSED) {
      init();
    }
    if (conn == null ||
        conn.readyState !== conn.OPEN) {
      console.log('Not connected, cannot send search', search);
      return;
    }
    conn.send(JSON.stringify({Search: search}));
  };
  init = function() {
    var hostport = window.location.hostname + (location.port ? ':' + location.port : '');
    conn = new WebSocket('ws://' + hostport + '/index.ws');
    conn.onopen = function() {
      sendSearch(location.search);
      $(window).bind('popstate', (function() { sendSearch(location.search); }));
    };
    var again = function(e) {
      if (!e.wasClean) {
        window.setTimeout(init, 5000);
      }
    };
    conn.onclose = again;
    conn.onerror = again;
    conn.onmessage = onmessage;
  };
  init();
  return {
    sendSearch: sendSearch,
    close: function() { return conn.close(); }
  };
};

function render_define(el) {
  var cl = jsdefines[$(el).attr('data-define')];
  return ReactDOM.render(React.createElement(cl), el);
}

function main() {
  for (var i = 0, loc = location.search.substr(1).split('&'); i < loc.length; i++) {
    if (loc[i].split('=')[0] === 'still') {
      return;
    }
  }

  var els = [];
  for (var i = 0, sel = $('.updates'); i < sel.length; i++) {
    els.push(render_define(sel[i]));
  }

  var onmessage = function(event) {
    var data = JSON.parse(event.data);
    if (data == null) {
      return;
    }
    if ((data.Reload != null) && data.Reload) {
      window.setTimeout((function() {
        location.reload(true);
      }), 5000);
      window.setTimeout(window.updates.close, 2000);
      console.log('in 5s: location.reload(true)');
      console.log('in 2s: window.updates.close()');
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

  window.updates = newwebsocket(onmessage); // neweventsource(onmessage);
}

main();

// Local variables:
// js-indent-level: 2
// js2-basic-offset: 2
// End:
