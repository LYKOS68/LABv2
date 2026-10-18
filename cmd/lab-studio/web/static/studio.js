// LAB-STUDIO: verschiebbare Arbeitsflaeche mit Feldern (Notizen, Songs).
// Verhalten wie im Prototyp, Daten pro Konto ueber /studio/api/.
"use strict";

const MIN_SCALE = 0.25;
const MAX_SCALE = 4;
const ZOOM_STEP = 1.1;
const FIELD_SIZES = { note: [240, 160], songs: [640, 420] }; // gleich fieldSizes in studio.go

const state = {
  view: { x: 0, y: 0, scale: 1 },
  fields: [],
  songs: [],
  selected: null,
};

const viewport = document.getElementById("viewport");
const world = document.getElementById("world");
const GRID = parseFloat(getComputedStyle(document.documentElement).getPropertyValue("--grid"));

let csrfToken = document.querySelector('meta[name="csrf-token"]').content;
let renderPending = false;
let writeQueue = Promise.resolve();
let viewSaveTimer = null;

// ── Kleine Helfer ────────────────────────────────────────────────────────────────────────

function byId(id) {
  return document.getElementById(id);
}

function showStatus(text) {
  byId("status").textContent = text || "";
}

function fail(error) {
  showStatus("Fehler: " + error.message);
}

function clamp(value, low, high) {
  return Math.min(high, Math.max(low, value));
}

function make(tag, className, text) {
  const node = document.createElement(tag);
  if (className) {
    node.className = className;
  }
  if (text) {
    node.textContent = text;
  }
  return node;
}

function findField(id) {
  return state.fields.find((field) => field.id === id) || null;
}

function songTitle(id) {
  const song = state.songs.find((s) => s.id === id);
  return song ? song.title : "Song";
}

function viewportCenter() {
  return { x: viewport.clientWidth / 2, y: viewport.clientHeight / 2 };
}

// ── Server ───────────────────────────────────────────────────────────────────────────────

// Sitzung abgelaufen: zurueck zur Anmeldung.
function checkLoggedIn(response) {
  if (response.status === 401) {
    window.location.href = "/login";
    throw new Error("nicht angemeldet");
  }
}

async function load(url) {
  const response = await fetch(url);
  checkLoggedIn(response);
  if (!response.ok) {
    throw new Error(String(response.status));
  }
  return response.json();
}

async function refreshCsrfToken() {
  const data = await load("/studio/api/csrf");
  csrfToken = data.token;
}

function writeRequest(method, url, body) {
  return fetch(url, {
    method: method,
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrfToken },
    body: JSON.stringify(body),
  });
}

// Ein abgelaufenes CSRF-Token (nach 30 Minuten) wird einmal erneuert, dann wird wiederholt.
async function send(method, url, body) {
  let response = await writeRequest(method, url, body);
  if (response.status === 403) {
    await refreshCsrfToken();
    response = await writeRequest(method, url, body);
  }
  checkLoggedIn(response);
  if (!response.ok) {
    throw new Error(response.status + " " + (await response.text()));
  }
  return response.json();
}

// Schreibzugriffe laufen nacheinander, damit kein alter Stand einen neueren ueberholt.
function write(method, url, body) {
  const result = writeQueue.then(() => send(method, url, body));
  writeQueue = result.catch(() => {});
  return result;
}

function patchField(field, changes) {
  return write("PATCH", "/studio/api/fields/" + field.id, changes).catch(fail);
}

// ── Ansicht: verschieben und zoomen ──────────────────────────────────────────────────────

function setView(x, y, scale, persist) {
  state.view = { x: x, y: y, scale: clamp(scale, MIN_SCALE, MAX_SCALE) };
  scheduleRender();
  if (persist) {
    saveViewSoon();
  }
}

// Gespeichert wird erst, wenn 250 ms lang nichts mehr passiert.
function saveViewSoon() {
  clearTimeout(viewSaveTimer);
  viewSaveTimer = setTimeout(() => {
    write("PUT", "/studio/api/workspace/view", state.view).catch(fail);
  }, 250);
}

// Zoomt so, dass der Punkt unter (px, py) an seiner Stelle bleibt.
function zoomAt(px, py, factor) {
  const view = state.view;
  const newScale = clamp(view.scale * factor, MIN_SCALE, MAX_SCALE);
  if (newScale === view.scale) {
    return;
  }
  const worldX = (px - view.x) / view.scale;
  const worldY = (py - view.y) / view.scale;
  setView(px - worldX * newScale, py - worldY * newScale, newScale, true);
}

function zoomAtCenter(factor) {
  const center = viewportCenter();
  zoomAt(center.x, center.y, factor);
}

function centerOn(field) {
  const center = viewportCenter();
  const scale = state.view.scale;
  const x = center.x - (field.x + field.w / 2) * scale;
  const y = center.y - (field.y + field.h / 2) * scale;
  setView(x, y, scale, true);
}

function select(id) {
  if (state.selected === id) {
    return;
  }
  state.selected = id;
  scheduleRender();
}

// ── Linke Leiste ─────────────────────────────────────────────────────────────────────────

async function addField(type) {
  const center = viewportCenter();
  const view = state.view;
  const [w, h] = FIELD_SIZES[type];
  const field = {
    type: type,
    x: Math.round((center.x - view.x) / view.scale - w / 2),
    y: Math.round((center.y - view.y) / view.scale - h / 2),
  };
  try {
    const created = await write("POST", "/studio/api/fields", field);
    state.fields.push(created);
    select(created.id);
    scheduleRender();
  } catch (error) {
    fail(error);
  }
}

// SONGS: vorhandenes SONGS-Fenster in die Mitte holen, sonst eins anlegen.
function showSongsField() {
  const field = state.fields.find((f) => f.type === "songs");
  if (!field) {
    addField("songs");
    return;
  }
  select(field.id);
  centerOn(field);
}

const actions = {
  "songs": showSongsField,
  "add-field": () => addField("note"),
  "zoom-in": () => zoomAtCenter(ZOOM_STEP),
  "zoom-out": () => zoomAtCenter(1 / ZOOM_STEP),
  "reset-view": () => setView(0, 0, 1, true),
};

byId("nav").addEventListener("click", (event) => {
  const button = event.target.closest("[data-action]");
  if (button && actions[button.dataset.action]) {
    actions[button.dataset.action]();
  }
});

// ── Felder zeichnen ──────────────────────────────────────────────────────────────────────

function scheduleRender() {
  if (renderPending) {
    return;
  }
  renderPending = true;
  requestAnimationFrame(render);
}

function render() {
  renderPending = false;
  applyView();
  const shown = new Set();
  for (const field of state.fields) {
    renderField(field);
    shown.add(String(field.id));
  }
  for (const node of [...world.children]) {
    if (!shown.has(node.dataset.id)) {
      node.remove();
    }
  }
}

function applyView() {
  const view = state.view;
  const grid = GRID * view.scale;
  world.style.transform = `translate(${view.x}px, ${view.y}px) scale(${view.scale})`;
  viewport.style.backgroundPosition = `${view.x}px ${view.y}px`;
  viewport.style.backgroundSize = `${grid}px ${grid}px`;
  byId("zoomlabel").textContent = Math.round(view.scale * 100) + " %";
}

function renderField(field) {
  let node = world.querySelector(`[data-id="${field.id}"]`);
  if (!node) {
    node = createFieldNode(field);
    world.appendChild(node);
  }
  // Inhalt nur neu bauen, wenn sich die Art geaendert hat (z. B. Liste -> Editor).
  const mode = field.type + ":" + field.song_id;
  if (node.dataset.mode !== mode) {
    node.dataset.mode = mode;
    fillField(node, field);
  }
  node.style.left = field.x + "px";
  node.style.top = field.y + "px";
  node.style.width = field.w + "px";
  node.style.height = field.h + "px";
  node.setAttribute("aria-selected", state.selected === field.id ? "true" : "false");
}

function createFieldNode(field) {
  const node = make("div", "field");
  const head = make("div", "head");
  head.appendChild(make("span", "title"));
  node.append(head, make("div", "body"));
  node.dataset.id = String(field.id);
  return node;
}

// Leert den Inhalt eines Felds, setzt den Titel und gibt den leeren Inhaltsbereich zurueck.
function resetField(node, title) {
  node.querySelector(".title").textContent = title;
  const body = node.querySelector(".body");
  body.replaceChildren();
  return body;
}

function fillField(node, field) {
  if (field.type === "note") {
    buildNote(node, field);
  } else if (field.song_id) {
    buildEditor(node, field);
  } else {
    buildSongList(node, field);
  }
}

// ── Notiz: speichert 0,5 s nach dem letzten Tastendruck ──────────────────────────────────

function buildNote(node, field) {
  const body = resetField(node, "NOTIZ");
  const area = make("textarea", "note");
  area.spellcheck = false;
  area.value = field.text || "";
  body.appendChild(area);

  let timer = null;
  const saveNote = () => {
    clearTimeout(timer);
    if ((field.text || "") === area.value) {
      return;
    }
    field.text = area.value;
    patchField(field, { text: area.value });
  };
  area.addEventListener("input", () => {
    clearTimeout(timer);
    timer = setTimeout(saveNote, 500);
  });
  area.addEventListener("change", saveNote);
}

// ── SONGS-Fenster: neuer Song + Kacheln ──────────────────────────────────────────────────

function buildSongList(node, field) {
  const body = resetField(node, "SONGS");
  body.append(buildNewSongForm(field), buildTiles(field));
}

function buildNewSongForm(field) {
  const form = make("form", "new-song");
  const input = make("input");
  input.placeholder = "Neuer Song";
  input.maxLength = 100;
  const button = make("button", "", "Anlegen");
  button.type = "submit";
  form.append(input, button);
  form.addEventListener("submit", (event) => {
    event.preventDefault();
    createSong(input.value, field);
  });
  return form;
}

function buildTiles(field) {
  const tiles = make("div", "tiles");
  if (state.songs.length === 0) {
    tiles.appendChild(make("p", "empty", "Noch keine Songs. Leg oben deinen ersten an."));
    return tiles;
  }
  for (const song of state.songs) {
    const tile = make("button", "tile");
    tile.append(make("div", "thumb", "Bild"), make("h2", "", song.title));
    tile.addEventListener("click", () => openSong(field, song.id));
    tiles.appendChild(tile);
  }
  return tiles;
}

async function createSong(title, field) {
  const trimmed = title.trim();
  if (!trimmed) {
    return;
  }
  try {
    const song = await write("POST", "/studio/api/songs", { title: trimmed });
    state.songs.push(song);
    state.songs.sort((a, b) => a.title.localeCompare(b.title, "de"));
    openSong(field, song.id);
  } catch (error) {
    fail(error);
  }
}

function openSong(field, songId) {
  field.song_id = songId;
  patchField(field, { song_id: songId });
  scheduleRender();
}

// ── Song-Editor: Speichern per Knopf oder Strg+S ─────────────────────────────────────────

function buildEditor(node, field) {
  const body = resetField(node, songTitle(field.song_id));
  node.dataset.unsaved = "false";

  const back = make("button", "", "Zurück");
  const save = make("button", "", "Speichern");
  const info = make("span", "estatus", "lädt …");
  const bar = make("div", "bar");
  bar.append(back, save, info);

  const editor = make("textarea", "lyrics");
  editor.spellcheck = false;
  editor.disabled = true;

  const wrap = make("div", "editor");
  wrap.append(bar, editor);
  body.appendChild(wrap);

  loadSongText(field.song_id, editor, info);
  editor.addEventListener("input", () => {
    node.dataset.unsaved = "true";
    info.textContent = "ungespeichert";
  });
  save.addEventListener("click", () => saveSong(node, field));
  back.addEventListener("click", () => closeSong(node, field));
}

async function loadSongText(songId, editor, info) {
  try {
    const song = await load("/studio/api/songs/" + songId);
    editor.value = song.text;
    editor.disabled = false;
    info.textContent = "";
  } catch (error) {
    info.textContent = "Fehler beim Laden: " + error.message;
  }
}

async function saveSong(node, field) {
  const editor = node.querySelector("textarea.lyrics");
  const info = node.querySelector(".estatus");
  if (!editor || editor.disabled) {
    return;
  }
  try {
    await write("PUT", "/studio/api/songs/" + field.song_id, { text: editor.value });
    node.dataset.unsaved = "false";
    info.textContent = "Gespeichert";
  } catch (error) {
    info.textContent = "Fehler " + error.message;
  }
}

function closeSong(node, field) {
  const unsaved = node.dataset.unsaved === "true";
  if (unsaved && !confirm("Ungespeicherte Änderungen verwerfen?")) {
    return;
  }
  field.song_id = 0;
  patchField(field, { song_id: 0 });
  scheduleRender();
}

document.addEventListener("keydown", (event) => {
  const isSave = (event.ctrlKey || event.metaKey) && event.key.toLowerCase() === "s";
  if (!isSave) {
    return;
  }
  event.preventDefault();
  const node = focusedFieldNode();
  if (!node) {
    return;
  }
  const field = findField(Number(node.dataset.id));
  if (field && field.type === "songs" && field.song_id) {
    saveSong(node, field);
  }
});

// Feld, in dessen Editor gerade getippt wird, sonst das ausgewaehlte Feld.
function focusedFieldNode() {
  const editor = document.activeElement && document.activeElement.closest(".editor");
  if (editor) {
    return editor.closest(".field");
  }
  if (state.selected === null) {
    return null;
  }
  return world.querySelector(`[data-id="${state.selected}"]`);
}

// ── Maus und Touch: Flaeche schieben, Zwei-Finger-Zoom, Feld an Kopfzeile ziehen ─────────

const pointers = {};
let pan = null;
let pinch = null;
let drag = null;

function pointerPosition(event) {
  const rect = viewport.getBoundingClientRect();
  return { x: event.clientX - rect.left, y: event.clientY - rect.top };
}

function startPan(point) {
  pan = { startX: point.x, startY: point.y, viewX: state.view.x, viewY: state.view.y };
}

function startPinch() {
  const [a, b] = Object.values(pointers);
  const view = state.view;
  pinch = {
    distance: Math.hypot(a.x - b.x, a.y - b.y) || 1,
    scale: view.scale,
    worldX: ((a.x + b.x) / 2 - view.x) / view.scale,
    worldY: ((a.y + b.y) / 2 - view.y) / view.scale,
  };
}

function startDrag(event, node) {
  const field = findField(Number(node.dataset.id));
  const point = pointerPosition(event);
  drag = { field: field, pointer: event.pointerId, startX: point.x, startY: point.y, fieldX: field.x, fieldY: field.y };
  viewport.setPointerCapture(event.pointerId);
  event.preventDefault();
}

viewport.addEventListener("pointerdown", (event) => {
  const node = event.target.closest(".field");
  const head = event.target.closest(".head");
  select(node ? Number(node.dataset.id) : null);
  if (head) {
    startDrag(event, node);
    return;
  }
  if (node) {
    return; // Klick in den Feldinhalt: Flaeche nicht verschieben
  }
  pointers[event.pointerId] = pointerPosition(event);
  viewport.setPointerCapture(event.pointerId);
  const count = Object.keys(pointers).length;
  if (count === 1) {
    startPan(pointers[event.pointerId]);
    viewport.style.cursor = "grabbing";
  } else if (count === 2) {
    pan = null;
    startPinch();
  }
});

viewport.addEventListener("pointermove", (event) => {
  const point = pointerPosition(event);
  if (drag && drag.pointer === event.pointerId) {
    moveDraggedField(point);
    return;
  }
  if (!pointers[event.pointerId]) {
    return;
  }
  pointers[event.pointerId] = point;
  if (pinch && Object.keys(pointers).length === 2) {
    applyPinch();
  } else if (pan) {
    setView(pan.viewX + point.x - pan.startX, pan.viewY + point.y - pan.startY, state.view.scale, false);
  }
});

function moveDraggedField(point) {
  const scale = state.view.scale;
  drag.field.x = Math.round(drag.fieldX + (point.x - drag.startX) / scale);
  drag.field.y = Math.round(drag.fieldY + (point.y - drag.startY) / scale);
  scheduleRender();
}

function applyPinch() {
  const [a, b] = Object.values(pointers);
  const distance = Math.hypot(a.x - b.x, a.y - b.y);
  const scale = clamp(pinch.scale * distance / pinch.distance, MIN_SCALE, MAX_SCALE);
  const x = (a.x + b.x) / 2 - pinch.worldX * scale;
  const y = (a.y + b.y) / 2 - pinch.worldY * scale;
  setView(x, y, scale, false);
}

function endPointer(event) {
  if (drag && drag.pointer === event.pointerId) {
    patchField(drag.field, { x: drag.field.x, y: drag.field.y });
    drag = null;
    return;
  }
  if (!pointers[event.pointerId]) {
    return;
  }
  delete pointers[event.pointerId];
  pinch = null;
  pan = null;
  const remaining = Object.values(pointers);
  if (remaining.length === 1) {
    startPan(remaining[0]);
  }
  if (remaining.length === 0) {
    viewport.style.cursor = "";
  }
  setView(state.view.x, state.view.y, state.view.scale, true);
}

viewport.addEventListener("pointerup", endPointer);
viewport.addEventListener("pointercancel", endPointer);

viewport.addEventListener("wheel", (event) => {
  event.preventDefault();
  const point = pointerPosition(event);
  const factor = event.deltaY < 0 ? ZOOM_STEP : 1 / ZOOM_STEP;
  zoomAt(point.x, point.y, factor);
}, { passive: false });

// ── Start ────────────────────────────────────────────────────────────────────────────────

async function start() {
  try {
    const [workspace, songs] = await Promise.all([
      load("/studio/api/workspace"),
      load("/studio/api/songs"),
    ]);
    state.view = workspace.view;
    state.fields = workspace.fields;
    state.songs = songs;
    scheduleRender();
  } catch (error) {
    showStatus("Fehler: Arbeitsfläche nicht ladbar – " + error.message);
  }
}

start();
