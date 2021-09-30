/**
 * Visualize the state we receive from the server.
 */
const visualize = (function() {
  // Keep a history from previous rounds so we can plot a graph.
  const history = {
    balance: [],
    boughtOffers: [],
  };
  let round = 0;
  return function(state) {
    $('#balance').text(state.balance);
    $('#boughtOffers').text(state.offers.length);

    history.balance.push([round, state.balance]);
    history.boughtOffers.push([round, state.offers.length]);
    $.plot('#graph', [{
      data: history.balance,
      label: 'Balance',
    }, {
      data: history.boughtOffers,
      label: 'Bought offers',
      yaxis: 2,
    }], {
      yaxes: [{}, {
        min: 0, max: 100,
        position: 'right',
      }],
    });
    round++;
  };
})();

/**
 * Calculate the offers to publish in the next round. Reimplement this with your automatic revenue manager.
 * @param {Object} state Game state as received from the server
 * @return {Object[]} List of offers to publish
 */
function makeOffers(state) {
  // make 100 random offers
  const platforms = ['ta', 'book', 'hrs', 'hc'];
  const offers = [];
  const price = parseFloat(document.getElementById('price').value);
  for (let i = 0; i < 100; i++) {
    offers.push({
      platform: platforms[Math.floor(Math.random() * platforms.length)],
      price: price + Math.random() * 20.0,
    });
  }
  return offers;
}

/**
 * Display a status message to the player.
 * @param {String} msg
 */
function status(msg) {
  $('#status').text(msg);
  console.log(msg);
}

/**
 * Join a networked game, and set up all required callbacks and
 * @param {Object} settings Name and websocket URL
 */
function joinGame(settings) {
  status('Connecting to ' + settings.url + ' ...');

  ws = new WebSocket(settings.url);
  ws.onopen = function() {
    status('Connection established! Joining as ' + settings.name + ' ...');
    ws.send(JSON.stringify(
        {
          Msg: 'join',
          Data: {
            name: settings.name,
          },
        },
    ));
  };

  ws.onmessage = function(message) {
    message = JSON.parse(message.data);
    switch (message.Msg) {
      case 'joined':
        status('Joined as player ' + message.Data.playerIndex);
        $('#setup').hide();
        $('#game').show();
        break;

      case 'state':
        console.log('Received', message);
        visualize(message.Data);

        const offers = makeOffers(message.Data);
        const input = {
          'Msg': 'input',
          'Data': {
            'offers': offers,
          },
        };
        console.log('Sending', input);
        ws.send(JSON.stringify(input));
        break;

      default:
        log('Invalid msg');
    }
  };
}

$(document).ready(function() {
  $('#join').click(function() {
    joinGame({
      name: $('#name').val(),
      url: $('#url').val(),
    });
  });
});
