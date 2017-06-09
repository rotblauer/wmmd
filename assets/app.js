var load = function () {
	console.log("hello");

	var body = document.getElementById("content-body");
	var header = document.getElementById("content-header");
	var sidebar = document.getElementById("content-sidebar");
	var footer = document.getElementById("content-footer");
	var ws = new WebSocket("ws://" + window.location.host + "/0");

	ws.onconnect = function (msg) {
		console.log("connected");
	}
	ws.ondisconnect = function (msg) {
		console.log("disconnected");
	}
	ws.onmessage = function (msg) {
		// {title: "", body: ""}
		var parsed = JSON.parse(msg.data);
		console.log("got message", msg);
		switch (parsed.title) {
			case "_Footer.md":
				footer.innerHTML = parsed.body;	
				break;
			case "_Sidebar.md":
				sidebar.innerHTML = parsed.body;
				break;
			default:
				body.innerHTML = parsed.body;
				header.innerHTML = parsed.title;
		}
	}	
}
