/**
 * Get and visualize the newest game state.
 */
function refresh() {
  $.ajax({
    url: 'http://localhost:8090/state',
  }).done(function(dataStr) {
    visualize(JSON.parse(dataStr));
  });
}

const minSize = 10;
const maxSize = 70;
const platforms = ['ta', 'book', 'hrs', 'hc'];

/**
 * Format a USD amount for display.
 * @param {Number} currency
 * @returns {String}
 */
function formatCurrency(currency) {
	const rounded = Math.round(currency * 100) / 100;
	return `$${rounded}`;
}

/**
 * Visualize the game state.
 * @param {Object} data Game state as received from server.
 */
function visualize(data) {
  const minMax = data.players.reduce(function(minMax, player) {
    if (minMax.min == null || player.balance < minMax.min) {
      minMax.min = player.balance;
    }
    if (minMax.max == null || player.balance > minMax.max) {
      minMax.max = player.balance;
    }
    return minMax;
  }, {'min': null, 'max': null});

  function calcSize(balance) {
    return minSize + (maxSize - minSize) * (balance - minMax.min) / (minMax.max - minMax.min);
  }

  const playersEl = $('<div>');
  data.players.forEach(function(player) {
    playersEl.append(
        $('<div class="player" style="width: {width}%"><span class="name">{name}</span><br/><span class="balance">{balance}</span></div>'
            .replace('{name}', player.name)
            .replace('{balance}', formatCurrency(player.balance))
            .replace('{width}', Math.round(calcSize(player.balance))),
        ),
    );
  });
  $('#players').html(playersEl.html());

  const offers = {
    'ta': [], 'book': [], 'hrs': [], 'hc': [],
  };
  data.offers.sort(function(o1, o2) {
    return o1.price - o2.price;
  }).forEach(function(offer) {
    offers[offer.platform].push(offer);
  });

  platforms.forEach(function(platform) {
    const platformEl = $('<ul>');
    platformEl.append('<li>TOTAL: {total}</li>'.replace('{total}', offers[platform].length));
    offers[platform].forEach(function(offer) {
      platformEl.append(
          $('<li>{price} ({player})</li>'
              .replace('{price}', formatCurrency(offer.price))
              .replace('{player}', offer.player),
          ),
      );
    });
    $('#' + platform).html(platformEl.html());
  });
}

$(document).ready(function() {
  setInterval(refresh, 1000);
});
