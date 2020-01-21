function SocketGO(remoteAddr, opts) {
  if (opts === undefined) {
    var opts = {}
  }

  this.remoteAddr = remoteAddr;
  this.handlers = {};
  this.debug = opts.debug || false
  this.sock = null;
}

SocketGO.prototype.connect = function() {
  var sock = new WebSocket(this.remoteAddr);
  var that = this;

  sock.addEventListener('open', function(event) {
    console.info("Websocket connection opened");
  });

  sock.addEventListener('message', function(event) {
    that.onMessage(event);
  });

  sock.addEventListener('error', function(event) {
    setTimeout(function() {
      console.info("Error in websocket connection, reconnecting");
      that.connect();
    }, 1000);
  });
  sock.addEventListener('close', function(event) {
    setTimeout(function() {
      console.info("Websocket connection closed, trying to reconnect");
      that.connect();
    }, 1000);
  });

  this.sock = sock;
};

SocketGO.prototype.onMessage = function(event) {
  var data = JSON.parse(event.data);
  var event = data.Event;
  var payload = data.Payload;

  if (!this.handlers.hasOwnProperty(event)) {
    return;
  }

  this.handlers[event](this, payload);
}

SocketGO.prototype.send = function(name, payload) {
  this.sock.send(JSON.stringify({
    Event: name,
    Payload: payload
  }));
}

SocketGO.prototype.handle = function(name, callback) {
  this.handlers[name] = callback;
}
