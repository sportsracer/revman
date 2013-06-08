
const pollFreq = 50
, rooms = 100
, platforms = ["ta", "book"];

var ws = null
, statusP = null
, setupForm = null
, joinButton = null
, gameWindow = null
, balanceSpan = null;

function status(msg) {
	if (statusP) {
		statusP.textContent = msg;
	}
}

function log(msg) {
	console.log(msg);
	status(msg);
}

window.onload = function() {
	statusP = document.getElementById("status");
	setupForm = document.getElementById("setup");
	joinButton = document.getElementById("join");
	gameWindow = document.getElementById("game");
	balanceSpan = document.getElementById("balance");
	boughtOffersSpan = document.getElementById("boughtOffers");

	joinButton.onclick = function() {
		var url = document.getElementById("url").value
		, name = document.getElementById("name").value;

		joinGame(url, name);
	};
};

function joinGame(url, name) {
	log("Connecting to " + url + " ...");
	ws = new WebSocket(url);

	ws.onopen = function() {
		log("Connection established! Joining as " + name + " ...");
		ws.send(JSON.stringify(
				{
					Msg: "join",
					Data: {
						name: name
					}
				}
			));
	};

	ws.onmessage = function(message) {
		message = JSON.parse(message.data);
		switch (message.Msg) {

		case "joined":
			log("Joined as player " + message.Data.playerIndex);
			setupForm.style.display = "none";
			gameWindow.style.display = "block";
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

function visualize(state) {
	balanceSpan.textContent = state.balance;
	boughtOffersSpan.textContent = state.offers.length;
}

function makeOffers(state) {
	// make 100 random offers
	var offers = []
	, price = parseFloat(document.getElementById("price").value);
	for (var i = 0; i < 100; i++) {
		offers.push({
			platform: platforms[Math.floor(Math.random() * platforms.length)],
			//price: 80.0 
			price: price + Math.random() * 5.0
		});
	}
	return offers;
}
