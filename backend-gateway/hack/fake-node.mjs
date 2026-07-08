#!/usr/bin/env node
// Backend Node FAKE para desenvolvimento local do gateway/frontend sem MPD.
// Fala o protocolo legado do server.js (parâmetros via headers, batradio-apikey).
// Uso: node backend-gateway/hack/fake-node.mjs [porta] (default 9320, apikey "dev")
import http from "node:http";

const PORT = Number(process.argv[2] ?? 9320);
const APIKEY = process.env.FAKE_NODE_APIKEY ?? "dev";

const ARTISTS = [
  ["Led Zeppelin", "rock"], ["AC/DC", "rock"], ["Cream", "rock"],
  ["Dire Straits", "rock"], ["B.B. King", "blues"], ["Stevie Ray Vaughan", "blues"],
  ["Robert Johnson", "blues"], ["Albert King", "blues"], ["ZZ Top", "rock"],
  ["Jimi Hendrix", "rock"], ["Derek and the Dominos", "rock"], ["Muddy Waters", "blues"],
];
const TITLES = [
  "Whole Lotta Love", "Back in Black", "Sunshine of Your Love", "Sultans of Swing",
  "The Thrill Is Gone", "Pride and Joy", "Cross Road Blues", "Born Under a Bad Sign",
  "Tush", "Voodoo Child (Slight Return)", "Layla", "Mannish Boy", "Black Dog",
  "Highway to Hell", "Sweet Home Chicago", "La Grange", "Texas Flood", "Kashmir",
];

// Acervo sintético grande (50k faixas) para exercitar busca/virtualização.
const files = [];
for (let i = 0; i < 50000; i++) {
  const [artist, genre] = ARTISTS[i % ARTISTS.length];
  const title = `${TITLES[i % TITLES.length]} (take ${Math.floor(i / TITLES.length) + 1})`;
  files.push({
    file: `${genre}/${artist.toLowerCase().replaceAll(/[^a-z0-9]+/g, "_")}/${String(i).padStart(5, "0")}.mp3`,
    Artist: artist,
    Title: title,
    Album: `Coletânea ${1 + (i % 40)}`,
    Genre: i % 37 === 0 ? "Vinheta" : genre === "rock" ? (i % 3 === 0 ? "Blues Rock" : "Rock") : "Blues",
    Time: String(120 + (i % 300)),
    "Last-Modified": "2026-01-01T00:00:00Z",
  });
}

// Fila inicial: 400 faixas
let queue = [];
let nextId = 1;
const resetPositions = () => queue.forEach((s, i) => { s.Pos = String(i); });
for (let i = 0; i < 400; i++) {
  queue.push({ ...files[(i * 131) % files.length], Id: String(nextId++) });
}
resetPositions();

let status = { state: "play", song: "5", songid: queue[5].Id, elapsed: "42.5", repeat: "1", random: "0", xfade: "0" };
let playlists = { "Programação Padrão": 1284, "Blues da Meia-Noite": 212, "Só Clássicos": 640, "Vinhetas MFM": 48, "Especial AC/DC": 96, "Domingo Blues": 180 };
setInterval(() => { status.elapsed = String(Number(status.elapsed) + 1); }, 1000);

const fullStatus = () => ({
  ...status,
  playlistlength: String(queue.length),
  currentSong: queue.find((s) => s.Pos === status.song),
  next: queue.find((s) => s.Pos === String(Number(status.song) + 1)),
});

http.createServer((req, res) => {
  const h = (n) => req.headers[n];
  res.setHeader("content-type", "application/json");
  if (h("batradio-apikey") !== APIKEY) {
    res.statusCode = 403;
    return res.end(JSON.stringify({ message: "Invalid API Key " }));
  }
  const route = `${req.method} ${req.url}`;
  const send = (v) => res.end(JSON.stringify(v));
  switch (route) {
    case "GET /status": return send(fullStatus());
    case "GET /list": {
      if (h("type") === "playlist")
        return send(Object.keys(playlists).map((p) => ({ playlist: p, "Last-Modified": "2026-06-20T10:00:00Z" })));
      return send(files);
    }
    case "POST /playlist": return send(queue);
    case "POST /playorpause":
      status.state = status.state === "play" ? "pause" : "play";
      return send(fullStatus());
    case "POST /play":
      status.song = h("position");
      status.state = "play";
      status.elapsed = "0";
      return send(fullStatus());
    case "POST /repeat": status.repeat = status.repeat === "1" ? "0" : "1"; return send(fullStatus());
    case "POST /shuffle": status.random = status.random === "1" ? "0" : "1"; return send(fullStatus());
    case "POST /fadein": status.xfade = status.xfade === "0" ? "1" : "0"; return send(fullStatus());
    case "POST /addtoplaylist": {
      const names = JSON.parse(h("files"));
      const pos = Number(h("position"));
      const items = names.map((n) => ({ ...files.find((f) => f.file === n) ?? { file: n, Title: n, Time: "60" }, Id: String(nextId++) }));
      queue.splice(pos, 0, ...items);
      resetPositions();
      return send(queue);
    }
    case "POST /delete": {
      const positions = JSON.parse(h("files")).map(Number).sort((a, b) => b - a);
      for (const p of positions) queue.splice(p, 1);
      resetPositions();
      return send(queue);
    }
    case "POST /move": {
      const from = Number(h("from-pos")), to = Number(h("to-pos"));
      const [item] = queue.splice(from, 1);
      queue.splice(to, 0, item);
      resetPositions();
      return send(queue);
    }
    case "POST /loadplaylist": {
      queue = [];
      for (let i = 0; i < (playlists[h("name")] ?? 100); i++)
        queue.push({ ...files[(i * 17) % files.length], Id: String(nextId++) });
      resetPositions();
      status.song = "0";
      return send(queue);
    }
    case "POST /saveplaylist": playlists[h("name")] = queue.length; return send(queue);
    case "POST /removeplaylist": delete playlists[h("name")]; return send(queue);
    default:
      res.statusCode = 404;
      return send({ message: `rota desconhecida ${route}` });
  }
}).listen(PORT, () => console.log(`fake-node ouvindo em :${PORT} (apikey: ${APIKEY}, acervo: ${files.length} faixas, fila: ${queue.length})`));
