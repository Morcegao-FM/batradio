var mpd = require("mpd");
var fs = require("fs");
var logger = require('log4js').getLogger();
var webServer = require('express');

var config = JSON.parse(fs.readFileSync("config.json"));


// Memory database
var currentPlaylist = Array();
var currentFiles = Array();
var currentStatus;

logger.level = config.log.minLevel;
logger.info(`Starting application - LOGLEVEL ${config.log.minLevel}`);


logger.debug(`Connecting to MPD - SERVER=${config.mpd.server} PORT=${config.mpd.port}`);
// MPD HANDLING
var client = mpd.connect({
    host: config.mpd.server,
    port: config.mpd.port,
    password: config.mpd.password
});

var start = async function () {
    logger.trace('BEGIN main.start()');
    await getFiles();
    await getPlaylist();
    await getStatus();
    logger.trace('END main.start()')
    
};

var error = function (err) {
    logger.error(`MPD error - ${err}`);    
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
            msg.forEach(item => {
                currentPlaylist.push(mpd.parseKeyValueMessage("file: " + item));
            });
            resolve();
            logger.debug(`MPDCOMMAND plalystinfo processed - ${currentPlaylist.length} SONGS`);
        });
    });
}



function promisedGetFiles(folder) {    
    logger.debug(`BEGIN promisedGetFiles`);
    return new Promise((resolve, reject) => {
        client.sendCommand(mpd.cmd("lsinfo ", [folder]),async (err, msg) => {
            logger.debug(`MPDCOMMAND lsinfo sent`);
            if (err) {                
                logger.error("promisedGetFiles error " + err);
                reject(err);
                return;
            }            
            newFiles = mpd.parseArrayMessage(msg);
            newFiles.forEach((t)=> { t.folder = folder + (t.directory?t.directory:t.file); currentFiles.push(t);})
            resolve();
            logger.debug(`MPDCOMMAND lsinfo processed - ${currentFiles.length} files/directories and playlists found`);
        });
    });
}

function promisedStatus() {
    logger.debug('BEGIN promisedStatus');
    return new Promise((resolve, reject) => {
        client.sendCommand(mpd.cmd("status", []), (err, msg) => {
            logger.debug('MPDCOMMAND status sent');
            if (err) {
                reject(err);
                logger.error(`MPDCOMMAND error ${err}`);
                return;
            }
            resolve(mpd.parseKeyValueMessage(msg));
            logger.debug(`MPDCOMMAND status processed ${currentStatus}`);
        });
    });
}

async function getPlaylist() {
    logger.info('BEGIN getPlaylist');
    await promisedGetPlaylist().then(() => {
        
    });
    logger.info('END getPlaylist');
}

async function getFiles(currentFolder) {    
    logger.info('BEGIN getFiles');
    if(!currentFolder)
    { 
        currentFolder = '';
        currentFiles = Array();
    }
    
     await promisedGetFiles(currentFolder).then(async () => {         
         for(i = 0; i < currentFiles.length;i++)
         {            
             if(currentFiles[i].directory)
             {
                await promisedGetFiles(currentFiles[i].folder).then( () => {}); 
             }
         }
     });
     logger.info('END getFiles');
}

async function getStatus()
{
    logger.info('BEGIN getStatus');
    await promisedStatus().then( (status) => {
        var currentPos = status.song;
        currentStatus = status;
        currentStatus.currentSong = currentPlaylist.find(t=> { return t.Pos == currentPos; } );
        currentStatus.next = currentPlaylist.find(t=> { return t.Pos == status.nextsong; } );
        logger.debug(`CURRENT SONG: ${currentStatus.currentSong.Title}`);

    });
    logger.info('END getStatus');
}


client.on("error", error);
client.on("ready", start);
