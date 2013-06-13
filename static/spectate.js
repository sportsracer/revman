function refresh() {
	$.ajax({
		url: "http://localhost:8090/state"	
	}).done(function(dataStr) {
		visualize(JSON.parse(dataStr));
	});
}

const minSize = 10
, maxSize = 70;

function visualize(data) {
	
	var minMax = data.players.reduce(function(minMax, player) {
		if (minMax.min == null || player.balance < minMax.min) {
			minMax.min = player.balance;
		}
		if (minMax.max == null || player.balance > minMax.max) {
			minMax.max = player.balance;
		}
		return minMax;
	}, {"min": null, "max": null});
	
	function calcSize(balance) {
		return minSize + (maxSize - minSize) * (balance - minMax.min) / (minMax.max - minMax.min);
	}
	
	var playersEl = $("<div>");
	data.players.forEach(function(player) {
		playersEl.append(
			$('<div class="player" style="width: {width}%"><span class="name">{name}</span><br/><span class="balance">{balance}</span></div>'
				.replace("{name}", player.name)
				.replace("{balance}", player.balance)
				.replace("{width}", Math.round(calcSize(player.balance)))
			)
		);
	});
	$("#players").html(playersEl.html());
}

$(document).ready(function() {
	setInterval(refresh, 1000);
});
