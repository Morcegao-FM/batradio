const fs = require('fs');
const path = require('path');
const request = require('supertest');
const { expect } = require('chai');
const EventEmitter = require('events');
const Module = require('module');

const backendDir = path.join(__dirname, '..');

// sample configuration used by the server during tests
const config = {
  mpd: { server: 'localhost', port: 6600, password: 'empty' },
  webserver: { port: 3005, apiKey: 'testkey' },
  log: { minLevel: 'info' }
};

let app;

// stub implementation for the mpd library
const stubMPD = {
  connect: () => {
    const client = new EventEmitter();
    client.sendCommand = (cmdObj, cb) => {
      const cmd = cmdObj.command.trim();
      if (cmd === 'status') {
        const msg = 'state: play\nsong: 1\nnextsong: 2\nxfade: 0\n';
        return process.nextTick(() => cb(null, msg));
      }
      if (cmd === 'playlistinfo') {
        const msg =
          'file: track1.mp3\nTitle: Song1\nPos: 1\nfile: track2.mp3\nTitle: Song2\nPos: 2\n';
        return process.nextTick(() => cb(null, msg));
      }
      if (cmd.startsWith('lsinfo')) {
        return process.nextTick(() => cb(null, ''));
      }
      process.nextTick(() => cb(null, ''));
    };
    process.nextTick(() => client.emit('ready'));
    return client;
  },
  cmd: (command, params) => ({ command, params }),
  parseKeyValueMessage: msg => {
    const obj = {};
    msg
      .split('\n')
      .filter(Boolean)
      .forEach(line => {
        const idx = line.indexOf(':');
        const key = line.slice(0, idx).trim();
        const val = line.slice(idx + 1).trim();
        obj[key] = val;
      });
    return obj;
  },
  parseListArrayMessage: msg => {
    const items = [];
    let current = {};
    msg
      .split('\n')
      .filter(Boolean)
      .forEach(line => {
        const idx = line.indexOf(':');
        const key = line.slice(0, idx).trim();
        const val = line.slice(idx + 1).trim();
        if (key === 'file' && Object.keys(current).length) {
          items.push(current);
          current = {};
        }
        current[key] = val;
      });
    if (Object.keys(current).length) items.push(current);
    return items;
  }
};

// patched express that exposes the created app and avoids listening on a port
const originalExpress = require('express');
function patchedExpress() {
  const appInstance = originalExpress();
  appInstance.listen = () => appInstance; // prevent real listen
  app = appInstance;
  return appInstance;
}
Object.assign(patchedExpress, originalExpress);

// override require for mpd and express before loading the server
const originalRequire = Module.prototype.require;
Module.prototype.require = function (modulePath) {
  if (modulePath === 'mpd') return stubMPD;
  if (modulePath === 'express') return patchedExpress;
  return originalRequire.call(this, modulePath);
};

// ensure timers created by the server do not keep the event loop alive
const originalSetInterval = global.setInterval;
global.setInterval = (...args) => {
  const timer = originalSetInterval(...args);
  timer.unref();
  return timer;
};

before(function (done) {
  fs.writeFileSync(path.join(backendDir, 'config.json'), JSON.stringify(config));
  require('../server');
  // wait a moment for the mocked MPD client to emit 'ready'
  setTimeout(done, 50);
});

after(function () {
  fs.unlinkSync(path.join(backendDir, 'config.json'));
  // restore patched functions
  Module.prototype.require = originalRequire;
  global.setInterval = originalSetInterval;
});

describe('GET /status', function () {
  it('should return 200 and status fields', function (done) {
    request(app)
      .get('/status')
      .set('batradio-apikey', config.webserver.apiKey)
      .expect('Content-Type', /json/)
      .expect(200)
      .end((err, res) => {
        if (err) return done(err);
        expect(res.body).to.have.property('state');
        expect(res.body).to.have.property('song');
        expect(res.body).to.have.property('nextsong');
        done();
      });
  });
});
