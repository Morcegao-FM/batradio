var mpd = require('mpd')
var fs = require('fs')

var config =JSON.parse(  fs.readFileSync('config.json') );

console.log( config.mpd.server );
