// LAB-CONNECT, Etappe 1: Oberflaechen-Prototyp. Alles laeuft im Browser,
// nichts wird an den Server geschickt (hub/rules/11-lab-connect-mvp-plan.md).
"use strict";

// ── Symbole (einfache Strich-Icons) ──────────────────────────────────────────────────────

const ICONS = {
  search: '<circle cx="11" cy="11" r="7"/><path d="M20 20l-4-4"/>',
  sliders: '<path d="M4 7h10M18 7h2M4 17h4M12 17h8"/><circle cx="16" cy="7" r="2"/><circle cx="10" cy="17" r="2"/>',
  pin: '<path d="M12 21s-7-6.2-7-11a7 7 0 0 1 14 0c0 4.8-7 11-7 11z"/><circle cx="12" cy="10" r="2.5"/>',
  wave: '<path d="M4 10v4M8 7v10M12 4v16M16 8v8M20 11v2"/>',
  mic: '<rect x="12" y="3" width="6" height="10" rx="3" transform="rotate(30 15 8)"/><path d="M11 13l-6 7M6 21c2 0 3-1 4-2"/>',
  disc: '<circle cx="12" cy="12" r="9"/><circle cx="12" cy="12" r="3"/><path d="M12 6a6 6 0 0 1 6 6"/>',
  pen: '<path d="M4 20l4-1 11-11-3-3L5 16l-1 4z"/><path d="M14 7l3 3"/>',
  star: '<path d="M12 3l2.7 5.6 6.1.9-4.4 4.3 1 6.1L12 17l-5.4 2.9 1-6.1-4.4-4.3 6.1-.9z"/>',
  heart: '<path d="M12 20s-7-4.4-9-9a4.8 4.8 0 0 1 9-3 4.8 4.8 0 0 1 9 3c-2 4.6-9 9-9 9z"/>',
  plus: '<path d="M12 5v14M5 12h14"/>',
  chat: '<path d="M4 18l1.5-4A8 8 0 1 1 9 18.6L4 20z"/>',
  user: '<circle cx="12" cy="8" r="4"/><path d="M4 21c1-4.5 4.5-7 8-7s7 2.5 8 7"/>',
  back: '<path d="M19 12H5M11 6l-6 6 6 6"/>',
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

// ── Daten: Kategorien und Beispielanzeigen ───────────────────────────────────────────────

const CATEGORIES = [
  { id: "alle", name: "Für dich", icon: "wave" },
  { id: "mix", name: "Mix & Master", icon: "mic" },
  { id: "beats", name: "Beats", icon: "disc" },
  { id: "songwriting", name: "Songwriting", icon: "pen" },
  { id: "features", name: "Features", icon: "star" },
];

const listings = [
  {
    kind: "biete", category: "mix", place: "Remote", price: 50, unit: "PLATZHALTER",
    title: "Mix & Master für deinen nächsten Release",
    description: "Du schickst die Spuren, ich mische und mastere deinen Song radiofertig.\nBeispieltext für den Designentwurf.",
  },
  {
    kind: "biete", category: "beats", place: "Mannheim · Remote", price: 30, unit: "PLATZHALTER",
    title: "Individuelle Beats für dein Projekt",
    description: "Boom-Bap, Trap oder Drill – auf deinen Sound zugeschnitten.\nBeispieltext für den Designentwurf.",
  },
  {
    kind: "suche", category: "songwriting", place: "Berlin", price: null, unit: "PLATZHALTER",
    title: "Suche Co-Writer für Hook und Refrain",
    description: "Strophen stehen, mir fehlt eine Hook, die hängen bleibt.\nBeispieltext für den Designentwurf.",
  },
  {
    kind: "biete", category: "features", place: "Hamburg · Remote", price: 80, unit: "PLATZHALTER",
    title: "Feature-Part (16 Bars) für deinen Track",
    description: "Ich schreibe und recorde einen Part passend zu deinem Thema.\nBeispieltext für den Designentwurf.",
  },
];

const state = {
  category: "alle",
  search: "",
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
  span.innerHTML = svg(name);
  span.dataset.icon = name;
  return span;
}

function categoryOf(listing) {
  return CATEGORIES.find((c) => c.id === listing.category) || CATEGORIES[0];
}

function priceText(listing) {
  return listing.price === null ? "Preis auf Anfrage" : listing.price + " €";
}

let toastTimer = null;

function showToast(text) {
  const toast = byId("toast");
  toast.textContent = text;
  toast.hidden = false;
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => {
    toast.hidden = true;
  }, 2600);
}

// ── Ansichten wechseln ───────────────────────────────────────────────────────────────────

function showView(name) {
  for (const view of document.querySelectorAll(".view")) {
    view.hidden = view.id !== "view-" + name;
  }
  for (const button of document.querySelectorAll("[data-view]")) {
    button.classList.toggle("active", button.dataset.view === name && !button.classList.contains("back"));
  }
  window.scrollTo(0, 0);
}

function showLater(title) {
  byId("later-title").textContent = title;
  showView("later");
  for (const button of document.querySelectorAll(".bottom [data-later]")) {
    button.classList.toggle("active", button.dataset.later === title);
  }
}

// ── Entdecken ────────────────────────────────────────────────────────────────────────────

function renderCategories() {
  const row = byId("categories");
  row.replaceChildren();
  for (const category of CATEGORIES) {
    const button = make("button", "category");
    button.classList.toggle("active", category.id === state.category);
    const bubble = make("span", "bubble");
    bubble.appendChild(iconSpan(category.icon));
    button.append(bubble, make("span", "", category.name));
    button.addEventListener("click", () => selectCategory(category.id));
    row.appendChild(button);
  }
}

function selectCategory(id) {
  state.category = id;
  renderCategories();
  renderCards();
}

function visibleListings() {
  const search = state.search.trim().toLowerCase();
  return listings.filter((listing) => {
    const categoryMatches = state.category === "alle" || listing.category === state.category;
    const searchMatches = search === "" || listing.title.toLowerCase().includes(search);
    return categoryMatches && searchMatches;
  });
}

function renderCards() {
  const cards = byId("cards");
  const shown = visibleListings();
  cards.replaceChildren();
  byId("count").textContent = shown.length + (shown.length === 1 ? " Anzeige" : " Anzeigen");
  if (shown.length === 0) {
    cards.appendChild(make("p", "empty", "Keine Anzeigen in dieser Auswahl."));
    return;
  }
  for (const listing of shown) {
    cards.appendChild(buildCard(listing));
  }
}

function buildPicture(listing, caption) {
  const picture = make("div", "picture");
  const badge = make("span", listing.own ? "badge own" : "badge", listing.own ? "DEIN ENTWURF" : "BEISPIEL");
  picture.append(badge, iconSpan(categoryOf(listing).icon), make("span", "caption", caption));
  return picture;
}

function buildCard(listing) {
  const card = make("button", "card");
  const picture = buildPicture(listing, "Vorschaubild folgt");
  const heart = iconSpan("heart");
  heart.className = "heart";
  picture.appendChild(heart);
  card.append(
    picture,
    make("p", "card-place", listing.place),
    make("p", "card-title", listing.title),
    make("p", "card-price", priceText(listing)),
  );
  card.addEventListener("click", () => openListing(listing));
  return card;
}

// ── Anzeige ──────────────────────────────────────────────────────────────────────────────

function openListing(listing) {
  state.current = listing;
  renderDetail();
  showView("anzeige");
}

function renderDetail() {
  const listing = state.current;
  const detail = byId("detail");
  detail.replaceChildren();

  const gallery = buildPicture(listing, "Galerie · Bildmaterial folgt");
  gallery.classList.add("gallery");

  const row = make("div", "detail-row");
  const place = make("span", "place");
  place.append(iconSpan("pin"), document.createTextNode(listing.place));
  const heart = make("button", "icon-button");
  heart.dataset.icon = "heart";
  heart.innerHTML = svg("heart");
  heart.setAttribute("aria-label", "Merken");
  heart.addEventListener("click", () => showToast("Favoriten folgen in Etappe 2."));
  row.append(place, heart);

  const price = make("p", "price", priceText(listing));
  if (listing.price !== null) {
    price.appendChild(make("small", "", " " + listing.unit));
  }

  const message = make("button", "primary");
  message.append(iconSpan("chat"), document.createTextNode("Nachricht schreiben"), make("small", "", "Entwurf"));
  message.addEventListener("click", () => showToast("Nachrichten folgen in Etappe 2."));

  detail.append(
    gallery,
    row,
    make("h1", "", listing.title),
    price,
    buildFacts(listing),
    make("p", "description", listing.description),
    message,
  );
}

function buildFacts(listing) {
  const facts = make("dl", "facts");
  const rows = [
    ["Kategorie", categoryOf(listing).name],
    ["Art", listing.kind === "biete" ? "Ich biete" : "Ich suche"],
  ];
  for (const [label, value] of rows) {
    const row = make("div");
    row.append(make("dt", "", label), make("dd", "", value));
    facts.appendChild(row);
  }
  return facts;
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

// Die neue Anzeige landet nur in dieser Ansicht; Speichern kommt mit Etappe 2.
function submitListing(event) {
  event.preventDefault();
  const form = event.target;
  const data = new FormData(form);
  const place = data.get("place").trim() || "Ohne Ort";
  const listing = {
    own: true,
    kind: data.get("kind"),
    category: data.get("category"),
    title: data.get("title").trim(),
    description: data.get("description").trim(),
    price: data.get("price") === "" ? null : Number(data.get("price")),
    unit: data.get("unit"),
    place: data.get("remote") ? place + " · Remote" : place,
  };
  listings.unshift(listing);
  form.reset();
  closeSheet();
  renderCards();
  openListing(listing);
  showToast("Nur in dieser Ansicht – Speichern folgt in Etappe 2.");
}

// ── Start ────────────────────────────────────────────────────────────────────────────────

function wireButtons() {
  for (const button of document.querySelectorAll("[data-view]")) {
    button.addEventListener("click", () => {
      if (button.dataset.view === "anzeige") {
        renderDetail();
      }
      showView(button.dataset.view);
    });
  }
  for (const button of document.querySelectorAll("[data-later]")) {
    button.addEventListener("click", () => showLater(button.dataset.later));
  }
  for (const button of document.querySelectorAll("[data-open-sheet]")) {
    button.addEventListener("click", openSheet);
  }
  for (const button of document.querySelectorAll("[data-close-sheet]")) {
    button.addEventListener("click", closeSheet);
  }
}

function start() {
  fillIcons(document);
  fillCategorySelect();
  wireButtons();
  byId("search").addEventListener("input", (event) => {
    state.search = event.target.value;
    renderCards();
  });
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
  renderCategories();
  renderCards();
}

start();
