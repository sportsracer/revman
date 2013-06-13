
const rooms = 100
, platforms = ["ta", "book"];

visualize = (function() {
	var history = {
		balance: [],
		boughtOffers: []
	}
	, round = 0;
	return function (state) {
		$("#balance").text(state.balance);
		$("#boughtOffers").text(state.offers.length);
		
		history.balance.push([round, state.balance]);
		history.boughtOffers.push([round, state.offers.length]);
		$.plot("#graph", [{
			data: history.balance,
			label: "Balance"
		}, {
			data: history.boughtOffers,
			label: "Bought offers",
			yaxis: 2
		}], {
			yaxes: [ {}, {
				min: 0, max: 100,
				position: "right"
			}]
		});
		round++;
	}
})();

function makeOffers(state) {
	// make 100 random offers
	var offers = []
	, price = parseFloat(document.getElementById("price").value);
	for (var i = 0; i < 100; i++) {
		offers.push({
			platform: platforms[Math.floor(Math.random() * platforms.length)],
			price: price + Math.random() * 5.0
		});
	}
	return offers;
}

function status(msg) {
	$("#status").text(msg);
	console.log(msg);
}

function joinGame(settings) {

	status("Connecting to " + settings.url + " ...");
	
	ws = new WebSocket(settings.url);
	ws.onopen = function() {
		status("Connection established! Joining as " + settings.name + " ...");
		ws.send(JSON.stringify(
				{
					Msg: "join",
					Data: {
						name: settings.name
					}
				}
			));
	};
	
	ws.onmessage = function(message) {
		message = JSON.parse(message.data);
		switch (message.Msg) {
		
		case "joined":
			status("Joined as player " + message.Data.playerIndex);
			$("#setup").hide();
			$("#game").show();
			break;
			
		case "state":
			console.log("Received", message);
			visualize(message.Data);
			
			var offers = makeOffers(message.Data)
			, input = {
				"Msg": "input",
				"Data": {
					"offers": offers
				}
			};
			console.log("Sending", input);
			ws.send(JSON.stringify(input));
			break;
			
		default:
			log("Invalid msg");
		}
	};

}

$(document).ready(function() {
	var statusEl = $("#status");
	
	function getSettings() {
		return {
			name: $("#name").val(),
			url: $("#url").val()
		}
	}
	
	$("#join").click(function() {
		joinGame(getSettings());
	});
	
});
