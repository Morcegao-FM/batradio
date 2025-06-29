var mpd = require("mpd");
var fs = require("fs");
var logger = require("log4js").getLogger();
var express = require("express");
var timeout = require('connect-timeout');
const { exception } = require("console");
const { timingSafeEqual } = require("crypto");
var app = express();
app.use(timeout('30s'));
app.use(haltOnTimedout);

var config = JSON.parse(fs.readFileSync("config.json"));

function haltOnTimedout (req, res, next) {
  if (!req.timedout) next();  
}
// Memory database
var currentPlaylist = Array();
var currentFiles = Array();
var currentStatus;
var busy = false;
var webserverstarted = false;

logger.level = config.log.minLevel;
logger.info(`Starting application - LOGLEVEL ${config.log.minLevel} Version 1.9.2`);
logger.debug(
  `Connecting to MPD - SERVER=${config.mpd.server} PORT=${config.mpd.port}`
);
// MPD HANDLING
var client = mpd.connect({
  host: config.mpd.server,
  port: config.mpd.port,
  password: config.mpd.password,
});

var start = async function () {
  logger.trace("BEGIN main.start()");
  
  await getFiles();
  await getPlaylist();
  await getStatus();
  setInterval(Ping, 59000);
  logger.info("Inicializando o servidor WEB. PORT ", config.webserver.port);
  if(!webserverstarted)
    app.listen(config.webserver.port);

  webserverstarted = true;
  logger.trace("END main.start()");
};

function nullfunc()
{

}

var error = function (err) {
  logger.error(`MPD error - ${err}`);
  client.on('error', nullfunc);
  client.on('ready', nullfunc);
  client = mpd.connect({
    host: config.mpd.server,
    port: config.mpd.port,
    password: config.mpd.password,
  });
  client.on("error", error);
  client.on("ready", start);  
  
};

function promisedGetPlaylist() {
  logger.debug(`BEGIN promisedGetPlaylist`);
  return new Promise((resolve, reject) => {
    client.sendCommand(mpd.cmd("playlistinfo", []), (err, msg) => {
      logger.debug(`MPDCOMMAND playlistinfo sent`);
      currentPlaylist = Array();
      if (err) {
        reject(err);
        logger.error("promisedGetPlaylist error " + err);
        return;
      }
      msg = msg.split("\nfile:");
      msg.forEach((item) => {
        currentPlaylist.push(mpd.parseKeyValueMessage("file: " + item));
      });
      resolve();
      logger.debug(
        `MPDCOMMAND plalystinfo processed - ${currentPlaylist.length} SONGS`
      );
    });
  });
}


function promisedGetFiles(folder) {
  logger.debug(`BEGIN promisedGetFiles`);
  return new Promise((resolve, reject) => {
    client.sendCommand(mpd.cmd("lsinfo ", [folder]), async (err, msg) => {
      logger.debug(`MPDCOMMAND lsinfo sent`);
      if (err) {
        logger.error("promisedGetFiles error " + err);
        reject(err);
        return;
      }
      newFiles = mpd.parseListArrayMessage(msg);
      newFiles.forEach((t) => {
        t.folder = (t.directory ? t.directory : t.file);
        currentFiles.push(t);
      });
      resolve();
      logger.debug(
        `MPDCOMMAND lsinfo processed - ${currentFiles.length} files/directories and playlists found`
      );
    });
  });
}


function promisedCommand(command, params)
{
  logger.debug(`BEGIN promisedCommand command=[${command}]  params.length=[${params.length}]`);

  return new Promise((resolve, reject) => {
    client.sendCommand(mpd.cmd(command, params), (err, msg) => {
      logger.debug(`MPDCOMMAND ${command} sent`);
      if (err) {
        if(err.indexOf && err.indexOf('ServiceUnavailableError'))
        {
          console.warn("Conexão caiu, tentando reconectar");
          client = mpd.connect({
            host: config.mpd.server,
            port: config.mpd.port,
            password: config.mpd.password,
          });
        }
        logger.error(err);
        reject(err);
        logger.debug(`MPDCOMMAND ${error} sent`);
        return;
      }
      resolve(mpd.parseKeyValueMessage(msg));
      logger.debug(`MPDCOMMAND ${command} processed`);
    });
  });
}


async function getPlaylist() {
  logger.info("BEGIN getPlaylist");
  isBusy();
  await promisedGetPlaylist().then(() => {
    busy = false;
  });
  logger.info("END getPlaylist");
}

async function getFiles(currentFolder) {
  logger.info("BEGIN getFiles");
  if (!currentFolder) {
    currentFolder = "";
    currentFiles = Array();
  }
  isBusy();
  await promisedGetFiles(currentFolder).then(async () => {
    for (i = 0; i < currentFiles.length; i++) {
      if (currentFiles[i].directory) {        
        await promisedGetFiles(currentFiles[i].folder).then(() => {
          
        });
      }
    }
    busy = false;
  });
  logger.info("END getFiles");
}

async function getStatus() {
  logger.info("BEGIN getStatus");
  await promisedCommand('status',[]).then((status) => {
    busy = false;
    var currentPos = status.song;
    if(!status.xfade) 
        status.xfade = 0;
    currentStatus = status;    
    currentStatus.currentSong = currentPlaylist.find((t) => {
      return t.Pos == currentPos;
    });
    currentStatus.next = currentPlaylist.find((t) => {
      return t.Pos == status.nextsong;
    });
    logger.debug(`CURRENT SONG: ${currentStatus.currentSong.Title}`);
  });
  logger.info("END getStatus");
}

client.on("error", error);
client.on("ready", start);

function isBusy() {
  if (busy) throw "MPD is busy, try again later";
  busy = true;
}

function checkAPIKey(req, res) {
  logger.debug(
    "checkAPIKey ",
    req.header("batradio-apikey"),
    config.webserver.apiKey
  );
  if (req.header("batradio-apikey") != config.webserver.apiKey) {
    logger.warn("Invalid API KEY", req.route.path, req.header("host"));
    return false;
  }
  return true;
}



async function Ping()
{
  logger.info('Starting  ping');
  await promisedCommand('ping',[]).catch( () => { this.error( 'Fail ping '); } );
  logger.info('Ending ping')
}







app.post('/playorpause', async (req,res) =>
{
  logger.debug('/playorpause required');
  if(!checkAPIKey(req,res))
  {
    return res.status("403").send({ message: "Invalid API Key " }).end();
  }    
  res.type("application/json");
  logger.debug("/playorpause will be sent");
  await getStatus();
  var command = 'play';
  if(currentStatus.state == 'play')
    command = 'pause'  
  isBusy();
  await promisedCommand(command, []).then((data) => {busy = false;} )  
  await getStatus();
  res.status(200).send(currentStatus).end();
});


app.post('/play', async (req,res) =>
{
  logger.debug('/playorpause required');
  if(!checkAPIKey(req,res))
  {
    return res.status("403").send({ message: "Invalid API Key " }).end();
  }    
  position = req.header('position');
  res.type("application/json");
  logger.debug("/playorpause will be sent");
  isBusy();  
  await promisedCommand('play', [position]).then((data) => {busy = false;} )  
  await getStatus();
  res.status(200).send(currentStatus).end();
});


app.post('/addtoplaylist', async (req,res) =>
{
  logger.debug('/addtoplaylist required');
  if(!checkAPIKey(req,res))
  {
    return res.status("403").send({ message: "Invalid API Key " }).end();
  }    
  musics = JSON.parse(req.header("files"));
  position = req.header("position");
  musics.forEach( async  (music) => { 
      await promisedCommand('addid', [music, position]).then( async (data) => {logger.debug(`Music added ${music} at ${position} - ${ JSON.stringify(data) }`)});
  });
  await getPlaylist();
  return res.status(200).send(currentPlaylist).end();
});



app.post('/delete', async (req,res) =>
{
  logger.debug('/delete required');
  if(!checkAPIKey(req,res))
  {
    return res.status("403").send({ message: "Invalid API Key " }).end();
  }    
  musics = JSON.parse(req.header("files"));  
  musics.forEach( async  (music) => { 
      await promisedCommand('delete', [music]).then( async (data) => {logger.debug(`Music removed from playlist at ${music} - ${ JSON.stringify(data) }`)});
  });
  await getPlaylist();
  return res.status(200).send(currentPlaylist).end();
});

app.post('/move', async (req,res) =>
{
  logger.debug('/move required');
  if(!checkAPIKey(req,res))
  {
    return res.status("403").send({ message: "Invalid API Key " }).end();
  }    
  var from = JSON.parse(req.header("from-pos"));  
  var to = JSON.parse(req.header("to-pos"));  
  await promisedCommand('move', [from,to]).then( async (data) => {logger.debug(`Music moved from ${from} to ${to} - ${ JSON.stringify(data) }`)});
  await getPlaylist();
  return res.status(200).send(currentPlaylist).end();
});




app.post('/repeat', async (req,res) =>
{
  console.log(currentStatus);
  logger.debug('/repeat required');
  if(!checkAPIKey(req,res))
  {
    return res.status("403").send({ message: "Invalid API Key " }).end();
  }    
  res.type("application/json");
  logger.debug("/repeat will be sent");
  await getStatus();
  var command = '1';
  if(currentStatus.repeat == '1')
    command = '0'
  isBusy();
  await promisedCommand('repeat', [command]).then((data) => {busy = false;} )  
  await getStatus();
  res.status(200).send(currentStatus).end();
});


app.post('/shuffle', async (req,res) =>
{
  logger.debug('/shuffle required');
  if(!checkAPIKey(req,res))
  {
    return res.status("403").send({ message: "Invalid API Key " }).end();
  }    
  res.type("application/json");
  logger.debug("/shuffle will be sent");
  await getStatus();
  var command = '1';
  if(currentStatus.random == '1')
    command = '0'  
  isBusy();
  await promisedCommand('random', [command]).then((data) => {busy = false;} )  
  await getStatus();
  res.status(200).send(currentStatus).end();
});


app.post('/fadein', async (req,res) =>
{
  logger.debug('/fadein required');
  if(!checkAPIKey(req,res))
  {
    return res.status("403").send({ message: "Invalid API Key " }).end();
  }    
  res.type("application/json");
  logger.debug("/fadein will be sent");
  await getStatus();
  isBusy();  
  await promisedCommand('crossfade', [currentStatus.xfade == '0'?1:0]).then((data) => {busy = false;} )  
  await getStatus();
  res.status(200).send(currentStatus).end();
});


app.get("/status", async (req, res) => {
  logger.debug("/status required");
  if (!checkAPIKey(req, res))
    return res.status("403").send({ message: "Invalid API Key " }).end();

  await getStatus();
  logger.debug("/status will be sent");
  res.type("application/json");
  res.status(200).send(currentStatus).end();
});

app.post("/playlist", async (req, res) => {
  updatePlaylist = req.header("update") == "true";
  logger.debug("/playlist required updatePlaylist=", updatePlaylist);
  if (!checkAPIKey(req, res))
    return res.status("403").send({ message: "Invalid API Key " }).end();
  if (updatePlaylist)
    await getPlaylist().catch((err) =>
      res.status(500).send({ message: err }).end()
    );
  logger.debug("/playlist will be sent");
  res.type("application/json");
  res.status(200).send(currentPlaylist).end();
});

app.post("/files", async (req, res) => {
  updateFiles = req.header("update") == "true";
  logger.debug("/files required updateFiles=", updateFiles);
  if (!checkAPIKey(req, res))
    return res.status("403").send({ message: "Invalid API Key " }).end();
  if (updateFiles)
    await getFiles().catch((err) =>
      res.status(500).send({ message: err }).end()
    );
  logger.debug("/files will be sent");
  res.type("application/json");
  res.status(200).send(currentFiles).end();
});

app.post("/save", async(req,res) => {
    playlistName = req.header("name");
    logger.info(`Salvando playlist ${playlistName}`);
    if (!checkAPIKey(req, res))
      return res.status("403").send({ message: "Invalid API Key " }).end();
    await promisedCommand("rm", [playlistName] ).catch((err) => {});
    await promisedCommand("save", [playlistName] );
    await getFiles();
    res.status(200).send(currentStatus);
}

);


app.post("/rm", async(req,res) => {
  playlistName = req.header("name");
  logger.info(`Salvando playlist ${playlistName}`);
  if (!checkAPIKey(req, res))
    return res.status("403").send({ message: "Invalid API Key " }).end();
  await promisedCommand("rm", [playlistName] ).catch((err) => {});  
  await getFiles();
  res.status(200).send(currentStatus);
}

);

app.post("/load", async(req,res) => {
  playlistName = req.header("name");
  logger.info(`Salvando playlist ${playlistName}`);
  if (!checkAPIKey(req, res))
    return res.status("403").send({ message: "Invalid API Key " }).end();
  await promisedCommand("load", [playlistName] ).catch((err) => { res.status(500).send({msg: `Não foi possível carregar a playlist ${playlistName}`})} );    
  res.status(200).send(currentStatus);
}

);

app.get("/list", async (req, res) => {
  updateFiles = req.header("update") == "true";
  type = req.header("type");
  logger.debug("/list required list update=", updateFiles, " type=", type);
  if (!checkAPIKey(req, res))
    return res.status("403").send({ message: "Invalid API Key " }).end();
  if (updateFiles)
    await getFiles().catch((err) =>
      res.status(500).send({ message: err }).end()
    );
  if (type != "directory" && type != "file" && type != "playlist") {
    logger.error("/list type not informed type=", type);
    res
      .status(500)
      .send("Invalid type " + type)
      .end();
    return;
  }
  logger.debug("/list will be sent");
  res.type("application/json");
  if (type == "directory")
    res
      .status(200)
      .send(currentFiles.filter((t) => t.directory && !t.playlist))
      .end();
  if (type == "file")
    res
      .status(200)
      .send(currentFiles.filter((t) => t.file))
      .end();
  if (type == "playlist")
    res
      .status(200)
      .send(currentFiles.filter((t) => t.playlist && t.playlist.indexOf('.m3u') == -1 ))
      .end();
});

app.post('/loadplaylist', async (req,res) =>
{
  logger.info('/loadplaylist required');
  if(!checkAPIKey(req,res))
  {
    return res.status("403").send({ message: "Invalid API Key " }).end();
  }    
  playlistName = req.header("name");
  await promisedCommand('stop', []).then(async(data) => {logger.debug(data); logger.info(`Clear before load ${playlistName}`)} );

  await promisedCommand('clear', []).then(async(data) => {logger.debug(data); logger.info(`Clear before load ${playlistName}`)} );

  await promisedCommand('load', [playlistName]).then(async(data) => {logger.debug(data); logger.info(`Playlist loaded ${playlistName}`)} );

  await getPlaylist();
  return res.status(200).send(currentPlaylist).end();
});


app.post('/saveplaylist', async (req,res) =>
{
  logger.info('/loadplaylist required');
  if(!checkAPIKey(req,res))
  {
    return res.status("403").send({ message: "Invalid API Key " }).end();
  }    
  playlistName = req.header("name");

  await promisedCommand('rm', [playlistName]).then(async(data) => {logger.debug(data); logger.info(`Playlist loaded ${playlistName}`)} ).catch( async(e) =>{ logger.info("Playlist nova")});
  await promisedCommand('save', [playlistName]).then(async(data) => {logger.debug(data); logger.info(`Playlist loaded ${playlistName}`)} );

  await getPlaylist();
  return res.status(200).send(currentPlaylist).end();
});


app.post('/removeplaylist', async (req,res) =>
{
  logger.info('/loadplaylist required');
  if(!checkAPIKey(req,res))
  {
    return res.status("403").send({ message: "Invalid API Key " }).end();
  }    
  playlistName = req.header("name");

  await promisedCommand('rm', [playlistName]).then(async(data) => {logger.debug(data); logger.info(`Playlist loaded ${playlistName}`)} );

  await getPlaylist();
  return res.status(200).send(currentPlaylist).end();
});