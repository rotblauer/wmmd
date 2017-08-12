var ws;
var toggleWikiStatus;
var toggleTypewriterStatus;
var localStorageWikiKey = "wmd_wiki_setting";
var localStorageTypewriterKey = "wmd_typewriter_setting";

emojify.setConfig({img_dir : 'assets/emoji'});

function getWikiStatus() {
	var r = "false";
	var gr = localStorage.getItem(localStorageWikiKey);
	console.log("wiki status stored", gr);
	if (!(gr === null || gr === "" || typeof(gr) === "undefined")) {
		r = gr;

	}
	r = r === "true" ? true : false;
	return setWikiStatus(r);
}
function setWikiStatus(bool) {
	localStorage.setItem(localStorageWikiKey, bool);
	return bool;
}
function getTypewriterStatus() {
	var r = "false";
	var gr = localStorage.getItem(localStorageTypewriterKey);
	console.log("typewriter status stored", gr);
	if (!(gr === null || gr === "" || typeof(gr) === "undefined")) {
		r = gr;

	}
	r = r === "true" ? true : false;
	return setTypewriterStatus(r);
}
function setTypewriterStatus(bool) {
	localStorage.setItem(localStorageTypewriterKey, bool);
	return bool;
}
// https://stackoverflow.com/questions/8917921/cross-browser-javascript-not-jquery-scroll-to-top-animation
function scrollToo(el, duration) {
  // var elOffset = $(el).offset().top;
  // console.log("elOffset", elOffset);
  // var elHeight = $(el).height();
  // console.log("elHeight", elHeight);
  // var windowHeight = window.innerHeight;
  // console.log("windowHeight", windowHeight);
  // var offset;
  // console.log("offset", offset);

  // if (elHeight < windowHeight) {
  //   offset = elOffset - ((windowHeight / 2) - (elHeight / 2));
  // } else {
  //   offset = elOffset;
  // }
  // var offset =  +
  var speed = 700;
  $('html, body').animate({scrollTop:$(el).offset().top - (window.innerHeight/2)}, speed);
  el.style.display = 'none';
  // $(el).hide();
}
// https://stackoverflow.com/questions/8024102/javascript-compare-strings-and-get-end-difference
function getDiff(string_a, string_b) {
	var first_occurance = string_b.indexOf(string_a);

  	if (!(first_occurance == -1)) {
    	// alert('Search string Not found');
  	} else {
    	string_a_length = string_a.length;
    	if (first_occurance == 0) {
      		new_string = string_b.substring(string_a_length);
    	} else {
	      	new_string = string_b.substring(0, first_occurance);
	      	new_string += string_b.substring(first_occurance + string_a_length);
    	}

    	var diffFirstIndex = string_b.indexOf(new_string);
    	var firstInsertable = string_b.substring(diffFirstIndex, string_b.length);
    	firstInsertable = firstInsertable.indexOf(">")+1;
    	diffFirstIndex = firstInsertable + diffFirstIndex;
    	string_b = string_b.substring(0, diffFirstIndex) + "<span class='change'>CHANGE</span>" + string_b.substring(diffFirstIndex, string_b.length+1)
 	}
 	return string_b
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
	var hudTypewriterStatus = document.getElementById("typewriter-status");
	var hudUpdating = document.getElementById("hud-center");

	var parsed = "",
	 diff = "",
	 lastBody = "",
	 lastWikiStatusBody = "";

	ws = new WebSocket("ws://" + window.location.host + "/x/0");

	toggleWikiStatus = function () {
		wikiStatus = setWikiStatus(!wikiStatus);
		console.log("wiki", wikiStatus);
		setDisplayFromWikiStatus(wikiStatus);
	};

	toggleTypewriterStatus = function () {
		typewriterStatus = setTypewriterStatus(!typewriterStatus);
		console.log("type", typewriterStatus);
		setDisplayFromTypewriterStatus(typewriterStatus);
	}

	function stripSuffix(title) {
		return title.replace(/\.[^/.]+$/, "")
	}

	function setDisplayFromWikiStatus(wikiStatus) {
		if (wikiStatus) {
			headerHolder.style.display = "block";
			sidebar.style.display = "block";
			footer.style.display = "block";
			if (!body.classList.contains("four-fifths")) {
				// body.classList.add("four-fifths");
			}
            if (body.classList.contains("centered")) {
                body.classList.remove("centered");
            }
			body.innerHTML = lastWikiStatusBody;
			emojify.run(body);
		} else {
			if (body.classList.contains("four-fifths")) {
				// body.classList.remove("four-fifths");
			}
            if (!body.classList.contains("centered")) {
                // body.classList.remove("four-fifths");
				body.classList.add("centered");
            }

			headerHolder.style.display = "none";
			sidebar.style.display = "none";
			footer.style.display = "none";
		}
		hudWikiStatus.innerHTML = wikiStatus;
	}
	function setDisplayFromTypewriterStatus(typewriterStatus) {
		if (typewriterStatus) {
			$("#content-body").addClass("typewriter");
			$("#content-header").addClass("typewriter");
		} else {
			$("#content-body").removeClass("typewriter");
			$("#content-header").removeClass("typewriter");
		}
		hudTypewriterStatus.innerHTML = typewriterStatus;
	}
	var wikiStatus = getWikiStatus();
	setDisplayFromWikiStatus(wikiStatus);
	var typewriterStatus = getTypewriterStatus();
	setDisplayFromTypewriterStatus(typewriterStatus);

	ws.onopen = function (msg) {
		console.log("connected");
	}
	ws.onclose = function (msg) {
		console.log("disconnected");
		hudUpdating.innerHTML = "<span>PLEASE RELOAD</span>";
		hudUpdating.classList.add("disconnected");

	}
	function showOrHideFandS() {
		var haveFooter = footer.innerHTML === "" ? false : true;
		var haveSidebar = sidebar.innerHTML === "" ? false : true;

		if (!haveFooter) {
			footer.style.display = "none";
		}
		if (!haveSidebar) {
			sidebar.style.display = "none";
		}
	}
	ws.onmessage = function (msg) {

		// hudUpdating.classList.remove("disconnected");
        $(hudUpdating).fadeIn(200);
        setTimeout(function () {
        	$(hudUpdating).fadeOut(400);
        }, 200);


		// {title: "", body: ""}
		parsed = JSON.parse(msg.data);
		console.log("got message", msg);

		if (parsed.body !== "" && parsed.title !== "") {
			$('#no-info').hide();
			$("#box").show();
		}

		// diff = getDiff(lastBody, parsed.body);
		// console.log("diff", diff);
		// lastBody = parsed.body;
		// parsed.body = diff;


		switch (stripSuffix(parsed.title)) {
		case "_Footer":
			footer.innerHTML = parsed.body;
			emojify.run(footer);
			break;
		case "_Sidebar":
			sidebar.innerHTML = parsed.body;
			emojify.run(sidebar);
			break;
		default:


			body.innerHTML = parsed.body;
			header.innerHTML = stripSuffix(parsed.title);
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
			header.innerHTML = stripSuffix(parsed.title);
			hudFilename.innerHTML = parsed.title;
			emojify.run(body);
			emojify.run(header);
		}
		// if (stripSuffix(parsed.title) != "_Footer" && stripSuffix(parsed.title) != "_Sidebar") {
		// 	debounce(scrollIt, 300);
		// }

		scrollIt();

		var d = new Date();
		var n = d.toTimeString();
		lastEdited.innerHTML = "you last updated this at " + n;
		debounce(showOrHideFandS,200);

        // $('div.highlight').each(function(i, block) {
        //     hljs.highlightBlock(block);
        //   });



	}

	function scrollIt() {
		console.log("scorllit called!");
		var changes = document.getElementsByClassName("suffix-change");
		if (changes.length > 0) {
			console.log("scrolling");
			scrollToo(changes[0], 600);
		} else {
			console.log("found no change marker");
		}
	}

	document.onkeypress = function (e) {
		console.log("keypress", e);
		e = e || window.event;
		// 119 == 'w'
		if (e.keyCode === 119) {
			toggleWikiStatus();
		}
		if (e.keyCode === 116) {
			toggleTypewriterStatus();
		}
		showOrHideFandS();
	}


	$(hudUpdating).on("click", function (e) {
		window.location.reload(true);
	})
}

// https://gist.github.com/sagiavinash/5c9084b79f68553c4b7d
// Returns a function, that, as long as it continues to be invoked, will not
// be triggered. The function will be called after it stops being called for
// N milliseconds. If `immediate` is passed, trigger the function on the
// leading edge, instead of the trailing.
function debounce(func, wait, immediate) {
	var timeout;
	return function() {
		var context = this, args = arguments;
		var later = function() {
			timeout = null;
			if (!immediate) func.apply(context, args);
		};
		var callNow = immediate && !timeout;
		clearTimeout(timeout);
		timeout = setTimeout(later, wait);
		if (callNow) func.apply(context, args);
	};
};
