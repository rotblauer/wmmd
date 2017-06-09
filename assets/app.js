var ws;
var toggleWikiStatus;
var localStorageWikiKey = "wub_wiki_setting";

function getWikiStatus() {
	var r = localStorage.getItem(localStorageWikiKey);
	console.log("wiki status stored", r);
	if (r === null || r === "" || typeof(r) === "undefined") {
		r = "false";
		return setWikiStatus(r);
	}
	return r === "false" ? false : true;
}
function setWikiStatus(bool) {
	localStorage.setItem(localStorageWikiKey, bool);
	return bool;
}
var load = function () {
	console.log("hello");

	var body = document.getElementById("content-body");
	var headerHolder = document.getElementById("content-header-holder");
	var header = document.getElementById("content-header");
	var sidebar = document.getElementById("content-sidebar");
	var footer = document.getElementById("content-footer");
	var lastEdited = document.getElementById("last-edited");
	var hudFilename = document.getElementById("hud-filename");
	var hudWikiStatus = document.getElementById("wiki-status");

	emojify.setConfig({img_dir : 'assets/emoji'});

	var wikiStatus = getWikiStatus();
	setDisplayFromWikiStatus(wikiStatus);
	var lastWikiStatusBody = "";

	ws = new WebSocket("ws://" + window.location.host + "/x/0");

	toggleWikiStatus = function () {
		wikiStatus = setWikiStatus(!wikiStatus);
		console.log("wiki", wikiStatus);
		setDisplayFromWikiStatus(wikiStatus);
	}

	function setDisplayFromWikiStatus(wikiStatus) {
		if (wikiStatus) {
			headerHolder.style.display = "block";
			sidebar.style.display = "block";
			footer.style.display = "block";
			if (!body.classList.contains("four-fifths")) {
				body.classList.add("four-fifths");
			}
			body.innerHTML = lastWikiStatusBody;
			emojify.run(body);
		} else {
			if (body.classList.contains("four-fifths")) {
				body.classList.remove("four-fifths");
			}
			headerHolder.style.display = "none";
			sidebar.style.display = "none";
			footer.style.display = "none";
		}
		
		hudWikiStatus.innerHTML = wikiStatus;
	}
	

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
			emojify.run(footer);
			break;
		case "_Sidebar.md":
			sidebar.innerHTML = parsed.body;
			emojify.run(sidebar);
			break;
		default:
			body.innerHTML = parsed.body;
			header.innerHTML = parsed.title;
			hudFilename.innerHTML = parsed.title;
			emojify.run(body);
			emojify.run(header);
			if (lastWikiStatusBody == "") { 
				lastWikiStatusBody = parsed.body;
			}
			if (wikiStatus) {
				lastWikiStatusBody = parsed.body;
			}
		}	
		if (!wikiStatus) {
			body.innerHTML = parsed.body;
			header.innerHTML = parsed.title;
			hudFilename.innerHTML = parsed.title;
			emojify.run(body);
			emojify.run(header);
		}

		var d = new Date();
		var n = d.toTimeString();
		lastEdited.innerHTML = "you last updated this at " + n ;
		
	}
}
