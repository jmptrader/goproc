console.log('Starting longrunning process...');

var i = 0;

setInterval(function() {
	if(i === 5) {
		throw new Error("test")
		// process.exit(1);
	}
	console.log(i);
	i++;
}, 1000);