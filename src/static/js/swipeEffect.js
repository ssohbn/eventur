let cards;
let i = 0;

const init = async () => {
  const res = await fetch("/api/events");
  cards = await res.json();
  console.log(cards);
  updateBackgound(cards[cards.length - i - 1].Img_url);
};
init();

const updateBackgound = (img) => {
  document.body.style.backgroundImage = `url(${img})`;
  console.log("image", img);
};

const swipeLeft = () => {
  if (i >= cards.length) return;
  const card = document.getElementById(cards.length - i - 1);
  console.log(card);
  card.classList.add("swipe-left-fade");
  if (i < cards.length - 1) {
    i += 1;
  } else {
    removeControls();
    fadeInOut("Not interested!", "#f56e64");
    return;
  }
  updateBackgound(cards[cards.length - i - 1].Img_url);
  fadeInOut("Not interested!", "#f56e64");
};

const swipeRight = () => {
  if (i >= cards.length) return;
  const card = document.getElementById(cards.length - i - 1);
  console.log(card);
  card.classList.add("swipe-right-fade");
  if (i < cards.length - 1) {
    i += 1;
  } else {
    removeControls();
    fadeInOut("Sounds Fun!", "#73f564");
    return;
  }
  updateBackgound(cards[cards.length - i - 1].Img_url);
  fadeInOut("Sounds Fun!", "#73f564");
};

const fadeInOut = (message, color) => {
  const popup = document.getElementById("pop-up");
  popup.innerHTML = message;
  popup.style.backgroundColor = color;
  popup.classList.add("fade-in");
  const delay = 500; // 1/2 seconds
  setTimeout(() => {
    popup.classList.remove("fade-in");
    popup.classList.add("fade-out");
  }, delay);
  setTimeout(() => {
    popup.classList.remove("fade-out");
  }, delay * 2);
};

const removeControls = () => {
  const controls = document.getElementById("event-controls");
  controls.style.display = "none";
  updateBackgound("");
};

const showControls = () => {
  const controls = document.getElementById("event-controls");
  controls.style.display = "default";
};
