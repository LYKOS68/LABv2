// LAB-CONNECT, Etappe 1: Oberflaechen-Prototyp nach dem Original (R:\dw\RAP-LAB-Etappe-1).
// Alles laeuft im Browser, nichts wird an den Server geschickt (hub/rules/11-lab-connect-mvp-plan.md).
"use strict";

// ── Symbole (einfache Strich-Icons) ──────────────────────────────────────────────────────

const ICONS = {
  search: '<circle cx="11" cy="11" r="7"/><path d="M20 20l-4-4"/>',
  sliders: '<path d="M4 7h10M18 7h2M4 17h4M12 17h8"/><circle cx="16" cy="7" r="2"/><circle cx="10" cy="17" r="2"/>',
  pin: '<path d="M12 21s-7-6.2-7-11a7 7 0 0 1 14 0c0 4.8-7 11-7 11z"/><circle cx="12" cy="10" r="2.5"/>',
  chevron: '<path d="M9 6l6 6-6 6"/>',
  wave: '<path d="M4 10v4M8 7v10M12 4v16M16 8v8M20 11v2"/>',
  mic: '<rect x="12" y="3" width="6" height="10" rx="3" transform="rotate(30 15 8)"/><path d="M11 13l-6 7M6 21c2 0 3-1 4-2"/>',
  disc: '<circle cx="12" cy="12" r="9"/><circle cx="12" cy="12" r="3"/><path d="M12 6a6 6 0 0 1 6 6"/>',
  pen: '<path d="M4 20l4-1 11-11-3-3L5 16l-1 4z"/><path d="M14 7l3 3"/>',
  video: '<rect x="3" y="6" width="13" height="12" rx="2"/><path d="M16 10l5-3v10l-5-3"/>',
  palette: '<path d="M12 3a9 9 0 1 0 0 18c1.5 0 2-1 2-2s-1-1.5-1-2.5 1-1.5 2-1.5h2a4 4 0 0 0 4-4c0-4.4-4-8-9-8z"/><circle cx="7.5" cy="11" r="1"/><circle cx="10" cy="7" r="1"/><circle cx="15" cy="7.5" r="1"/>',
  sound: '<path d="M3 15V9M7 18V6M11 14v-4M15 17V7M19 13v-2"/>',
  studio: '<rect x="4" y="3" width="16" height="18" rx="2"/><circle cx="12" cy="14" r="4"/><circle cx="12" cy="7" r="1"/>',
  people: '<circle cx="9" cy="8" r="3"/><circle cx="17" cy="9" r="2.5"/><path d="M3 20c.8-3.5 3.2-5.5 6-5.5s5.2 2 6 5.5M15 15c2.5 0 4.5 1.5 5.5 4.5"/>',
  heart: '<path d="M12 20s-7-4.4-9-9a4.8 4.8 0 0 1 9-3 4.8 4.8 0 0 1 9 3c-2 4.6-9 9-9 9z"/>',
  plus: '<path d="M12 5v14M5 12h14"/>',
  chat: '<path d="M4 18l1.5-4A8 8 0 1 1 9 18.6L4 20z"/>',
  user: '<circle cx="12" cy="8" r="4"/><path d="M4 21c1-4.5 4.5-7 8-7s7 2.5 8 7"/>',
  back: '<path d="M19 12H5M11 6l-6 6 6 6"/>',
  external: '<path d="M7 17L17 7M9 7h8v8"/>',
  image: '<rect x="3" y="5" width="15" height="14" rx="2"/><circle cx="8" cy="10" r="1.5"/><path d="M3 17l5-4 4 3 3-2 3 2M19 3v6M16 6h6"/>',
  close: '<path d="M6 6l12 12M18 6L6 18"/>',
};

function svg(name) {
  return '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" ' +
    'stroke-linecap="round" stroke-linejoin="round">' + ICONS[name] + "</svg>";
}

function fillIcons(root) {
  for (const node of root.querySelectorAll("[data-icon]")) {
    node.innerHTML = svg(node.dataset.icon);
  }
}

// ── Daten: Kategorien und Beispielanzeigen (wie im Original) ─────────────────────────────

// chip: als runder Knopf in Entdecken sichtbar; alle stehen im Inserieren-Formular.
const CATEGORIES = [
  { id: "alle", name: "Für dich", icon: "wave", chip: true },
  { id: "mix", name: "Mix & Master", icon: "mic", chip: true },
  { id: "beats", name: "Beats", icon: "disc", chip: true },
  { id: "songwriting", name: "Songwriting", icon: "pen", chip: true },
  { id: "video", name: "Video", icon: "video", chip: true },
  { id: "artwork", name: "Artwork", icon: "palette", chip: true },
  { id: "presets", name: "Presets & Sounds", icon: "sound", chip: false },
  { id: "recording", name: "Recording & Studio", icon: "studio", chip: false },
  { id: "features", name: "Features & Zusammenarbeit", icon: "people", chip: false },
];

const PRICE_TYPES = {
  fixed: { label: "Festpreis", withAmount: true },
  from: { label: "Ab-Preis", withAmount: true },
  negotiable: { label: "Verhandelbar", withAmount: true },
  request: { label: "Auf Anfrage", withAmount: false },
  collab: { label: "Zusammenarbeit", withAmount: false },
};

const PLACEHOLDER_TEXT = "Hier steht später die Beschreibung der inserierenden Person: " +
  "Leistungsumfang, Ablauf und Voraussetzungen. Dieser Text ist ausdrücklich ein Layout-Platzhalter, kein echtes Angebot.";

let nextId = 1;

function listing(fields) {
  return { id: nextId++, kind: "biete", own: false, description: PLACEHOLDER_TEXT, unit: "", ...fields };
}

const listings = [
  listing({ title: "Mix & Master für deinen nächsten Release", category: "mix", priceType: "fixed", amount: "50", unit: "pro Song", place: "Remote", remote: true }),
  listing({ title: "Individuelle Beats für deinen Sound", category: "beats", priceType: "request", place: "Mannheim · Remote", remote: true }),
  listing({ title: "Hooks und Parts gemeinsam schreiben", category: "songwriting", priceType: "collab", place: "Remote", remote: true }),
  listing({ title: "Cover-Art für deine nächste Single", category: "artwork", priceType: "request", place: "Remote", remote: true }),
];

const state = {
  view: "entdecken",
  category: "alle",
  search: "",
  remoteOnly: false,
  favorites: new Set(),
  current: listings[0],
};

// ── Kleine Helfer ────────────────────────────────────────────────────────────────────────

function byId(id) {
  return document.getElementById(id);
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

function iconSpan(name) {
  const span = make("span");
  span.dataset.icon = name;
  span.innerHTML = svg(name);
  return span;
}

function categoryOf(item) {
  return CATEGORIES.find((c) => c.id === item.category) || CATEGORIES[0];
}

function priceText(item) {
  const amount = item.amount + " €";
  switch (item.priceType) {
    case "fixed":
      return amount;
    case "from":
      return "ab " + amount;
    case "negotiable":
      return amount + " VB";
    default:
      return PRICE_TYPES[item.priceType].label;
  }
}

let toastTimer = null;

function showToast(text) {
  const toast = byId("toast");
  toast.textContent = text;
  toast.hidden = false;
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => {
    toast.hidden = true;
  }, 3200);
}

// ── Ansichten ────────────────────────────────────────────────────────────────────────────

// Welche <section> zu welcher Ansicht gehoert.
const SECTION_OF_VIEW = {
  entdecken: "view-entdecken",
  favoriten: "view-entdecken",
  anzeige: "view-anzeige",
  nachrichten: "view-info",
  profil: "view-info",
};

function showView(view) {
  state.view = view;
  document.querySelector(".app").dataset.view = view;
  for (const section of document.querySelectorAll(".view")) {
    section.hidden = section.id !== SECTION_OF_VIEW[view];
  }
  for (const button of document.querySelectorAll(".tabs [data-view], .bottom [data-view]")) {
    button.classList.toggle("active", button.dataset.view === view);
  }
  if (view === "entdecken" || view === "favoriten") {
    renderExplore();
  } else if (view === "anzeige") {
    renderDetail();
  } else {
    byId("info-title").textContent = view === "nachrichten" ? "Deine Gespräche" : "Dein Künstlerprofil";
  }
  window.scrollTo(0, 0);
}

// ── Entdecken / Favoriten ────────────────────────────────────────────────────────────────

function renderExplore() {
  const favoritesOnly = state.view === "favoriten";
  byId("explore-title").textContent = favoritesOnly ? "Deine Merkliste" : "Was fehlt deinem Projekt?";
  byId("results-title").textContent = favoritesOnly ? "Vorgemerkte Beispiele" : "Entdecke Möglichkeiten";
  renderCategories();
  renderCards();
}

function renderCategories() {
  const row = byId("categories");
  row.replaceChildren();
  for (const category of CATEGORIES.filter((c) => c.chip)) {
    const button = make("button", "category");
    button.classList.toggle("active", category.id === state.category);
    const bubble = make("span", "bubble");
    bubble.appendChild(iconSpan(category.icon));
    button.append(bubble, make("span", "", category.name));
    button.addEventListener("click", () => {
      state.category = category.id;
      renderExplore();
    });
    row.appendChild(button);
  }
}

function matchesFilters(item) {
  const search = state.search.trim().toLocaleLowerCase("de");
  if (state.category !== "alle" && item.category !== state.category) {
    return false;
  }
  if (search && !item.title.toLocaleLowerCase("de").includes(search)) {
    return false;
  }
  if (state.remoteOnly && !item.remote) {
    return false;
  }
  if (state.view === "favoriten" && !state.favorites.has(item.id)) {
    return false;
  }
  return true;
}

function renderCards() {
  const shown = listings.filter(matchesFilters);
  const cards = byId("cards");
  cards.replaceChildren();
  for (const item of shown) {
    cards.appendChild(buildCard(item));
  }
  byId("count").textContent = shown.length + (shown.length === 1 ? " Beispiel" : " Beispiele");
  byId("empty").hidden = shown.length > 0;
}

function buildPicture(item, caption, label) {
  const picture = make("div", "picture");
  const badge = make("span", item.own ? "badge own" : "badge", item.own ? "DEIN ENTWURF" : label);
  picture.append(badge, iconSpan(categoryOf(item).icon), make("span", "caption", caption));
  return picture;
}

function buildFavoriteButton(item, className) {
  const button = make("button", className);
  button.type = "button";
  const saved = state.favorites.has(item.id);
  button.classList.toggle("saved", saved);
  button.setAttribute("aria-pressed", String(saved));
  button.setAttribute("aria-label", saved ? "Beispiel entfernen" : "Beispiel merken");
  button.appendChild(iconSpan("heart"));
  button.addEventListener("click", (event) => {
    event.stopPropagation();
    toggleFavorite(item);
  });
  return button;
}

function toggleFavorite(item) {
  if (state.favorites.has(item.id)) {
    state.favorites.delete(item.id);
  } else {
    state.favorites.add(item.id);
  }
  showView(state.view);
}

function buildCard(item) {
  const card = make("article", "card");
  const picture = buildPicture(item, "Vorschaubild folgt", "BEISPIEL");
  picture.appendChild(buildFavoriteButton(item, "heart"));

  const info = make("div", "card-info");
  info.append(
    make("p", "card-place", item.place),
    make("p", "card-title", item.title),
    buildPrice(item, "card-price"),
  );
  card.append(picture, info);
  card.addEventListener("click", () => openListing(item));
  return card;
}

function buildPrice(item, className) {
  const price = make("p", className, priceText(item));
  if (item.unit && PRICE_TYPES[item.priceType].withAmount) {
    price.appendChild(make("small", "", " " + item.unit));
  }
  return price;
}

// ── Anzeige ──────────────────────────────────────────────────────────────────────────────

function openListing(item) {
  state.current = item;
  showView("anzeige");
}

function renderDetail() {
  const item = state.current;
  const detail = byId("detail");
  detail.replaceChildren();

  const gallery = buildPicture(item, "Galerie · Bildmaterial folgt", "BEISPIELANZEIGE");
  gallery.classList.add("gallery");

  const meta = make("div", "detail-row");
  const place = make("span", "place");
  place.append(iconSpan("pin"), document.createTextNode(item.place));
  meta.append(place, buildFavoriteButton(item, "icon-button"));

  const info = make("div", "detail-info");
  info.append(
    meta,
    make("h1", "", item.title),
    buildPrice(item, "price"),
    buildFacts(item),
    buildSection("Über das Angebot", make("p", "description", item.description)),
    buildSection("Arbeitsproben", buildSampleLink()),
    buildProviderCard(),
    make("p", "fine-print", "RAP-LAB vermittelt Kontakte. Keine Zahlung oder Buchung über die Plattform."),
    buildMessageAction(),
  );
  detail.append(gallery, info);
}

function buildFacts(item) {
  const facts = make("dl", "facts");
  const rows = [
    ["Kategorie", categoryOf(item).name],
    ["Zusammenarbeit", item.remote ? "Remote möglich" : "Vor Ort"],
    ["Anzeigentyp", item.kind === "biete" ? "Ich biete" : "Ich suche"],
  ];
  for (const [label, value] of rows) {
    const row = make("div");
    row.append(make("dt", "", label), make("dd", "", value));
    facts.appendChild(row);
  }
  return facts;
}

function buildSection(title, content) {
  const section = make("section", "detail-section");
  section.append(make("h2", "", title), content);
  return section;
}

function buildSampleLink() {
  const link = make("button", "link-row");
  link.append(iconSpan("wave"), make("span", "", "Externer Link folgt"), iconSpan("external"));
  link.addEventListener("click", () => showToast("Arbeitsproben-Links folgen in Etappe 2."));
  return link;
}

function buildProviderCard() {
  const card = make("button", "provider");
  const avatar = make("span", "avatar");
  avatar.appendChild(iconSpan("user"));
  const text = make("span", "provider-text");
  text.append(make("strong", "", "Künstlerprofil"), make("span", "", "Profilangaben folgen"));
  card.append(avatar, text, iconSpan("chevron"));
  card.addEventListener("click", () => showToast("Künstlerprofile folgen in Etappe 2."));
  return card;
}

function buildMessageAction() {
  const bar = make("div", "detail-action");
  const button = make("button", "primary");
  button.append(iconSpan("chat"), document.createTextNode("Nachricht schreiben"), make("small", "", "Entwurf"));
  button.addEventListener("click", () => showToast("Nachrichten folgen in Etappe 2."));
  bar.appendChild(button);
  return bar;
}

// ── Inserieren ───────────────────────────────────────────────────────────────────────────

function openSheet() {
  byId("sheet").hidden = false;
  document.body.style.overflow = "hidden";
}

function closeSheet() {
  byId("sheet").hidden = true;
  document.body.style.overflow = "";
}

function fillCategorySelect() {
  const select = byId("category-select");
  for (const category of CATEGORIES.slice(1)) {
    const option = make("option", "", category.name);
    option.value = category.id;
    select.appendChild(option);
  }
}

// Betragsfeld nur bei Festpreis, Ab-Preis und Verhandelbar.
function updateAmountField() {
  const type = byId("price-type").value;
  byId("amount-field").hidden = !PRICE_TYPES[type].withAmount;
}

// Die neue Anzeige landet nur in dieser Ansicht; Speichern kommt mit Etappe 2.
function submitListing(event) {
  event.preventDefault();
  const form = event.target;
  const data = new FormData(form);
  let priceType = data.get("priceType");
  const amount = data.get("amount").trim();
  if (PRICE_TYPES[priceType].withAmount && amount === "") {
    priceType = "request";
  }
  const remote = data.get("remote") !== null;
  const place = data.get("place").trim() || "Ohne Ort";
  const item = listing({
    own: true,
    kind: data.get("kind"),
    category: data.get("category"),
    title: data.get("title").trim(),
    description: data.get("description").trim(),
    priceType: priceType,
    amount: amount,
    remote: remote,
    place: remote ? place + " · Remote" : place,
  });
  listings.unshift(item);
  form.reset();
  updateAmountField();
  closeSheet();
  openListing(item);
  showToast("Nur in dieser Ansicht – Speichern folgt in Etappe 2.");
}

// ── Start ────────────────────────────────────────────────────────────────────────────────

function wireButtons() {
  for (const button of document.querySelectorAll("[data-view]")) {
    button.addEventListener("click", () => showView(button.dataset.view));
  }
  for (const button of document.querySelectorAll("[data-open-sheet]")) {
    button.addEventListener("click", openSheet);
  }
  for (const button of document.querySelectorAll("[data-close-sheet]")) {
    button.addEventListener("click", closeSheet);
  }
  for (const button of document.querySelectorAll("[data-toggle-filters]")) {
    button.addEventListener("click", () => {
      byId("filters").hidden = !byId("filters").hidden;
    });
  }
  for (const button of document.querySelectorAll("[data-notice]")) {
    button.addEventListener("click", () => showToast(button.dataset.notice));
  }
}

function wireInputs() {
  byId("search").addEventListener("input", (event) => {
    state.search = event.target.value;
    renderCards();
  });
  byId("remote-only").addEventListener("change", (event) => {
    state.remoteOnly = event.target.checked;
    renderCards();
  });
  byId("price-type").addEventListener("change", updateAmountField);
  byId("listing-form").addEventListener("submit", submitListing);
  byId("sheet").addEventListener("click", (event) => {
    if (event.target.id === "sheet") {
      closeSheet();
    }
  });
  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      closeSheet();
    }
  });
}

function start() {
  fillIcons(document);
  fillCategorySelect();
  wireButtons();
  wireInputs();
  updateAmountField();
  showView("entdecken");
}

start();
