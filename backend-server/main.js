var mpd = require("mpd");
var fs = require("fs");

var config = JSON.parse(fs.readFileSync("config.json"));
var currentPlaylist = Array();

console.log("Connecting");

var client = mpd.connect({
  host: config.mpd.server,
  port: config.mpd.port,
  password: config.mpd.password
});

var start = function() {
  getPlaylist();  
};

var error = function(err) {
  console.log(err);
};

function promisedGetPlaylist() {
  return new Promise((resolve, reject) => {    
    client.sendCommand(mpd.cmd("playlistinfo", []), (err, msg) => {
        console.log('playlistinfo required');
      if (err) 
        console.error("err " + err);
      msg = msg.split("\nfile:");
      msg.forEach((item) => {  currentPlaylist.push( mpd.parseKeyValueMessage('file: ' + item)); });
      resolve();
      console.log('playlistinfo finished');
    });    
  });
}

async function getPlaylist() {  
  await promisedGetPlaylist().then(() => {
    console.log(currentPlaylist[0]);
  });
  
}

client.on("error", error);
client.on("ready", start);
