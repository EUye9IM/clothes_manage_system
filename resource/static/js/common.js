xmlhttp = new XMLHttpRequest();
function toURL(data) {
	var tmparr = [];
	for (var i in data) {
		var key = encodeURIComponent(i);
		var value = encodeURIComponent(data[i]);
		tmparr.push(key + "=" + value)
	}
	return tmparr.join('&')
}