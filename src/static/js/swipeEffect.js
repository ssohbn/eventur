let cards;
let i = 0;

const init = async () => {
  const res = await fetch("/api/events");
  cards = await res.json();
  updateBackgound(cards[cards.length - i - 1].Img_url);
  revealCard();
};
init();

const updateBackgound = (img) => {
  document.body.style.backgroundImage = `url(${img})`;
};

const swipeLeft = () => {
  if (i >= cards.length) return;
  const card = document.getElementById(cards.length - i - 1);
  card.classList.add("swipe-left-fade");
  if (i < cards.length - 1) {
    i += 1;
  } else {
    removeControls();
    fadeInOut("Not interested!", "#f56e64");
    return;
  }
  revealCard();
  updateBackgound(cards[cards.length - i - 1].Img_url);
  fadeInOut("Not interested!", "#f56e64");
};

const swipeRight = () => {
  if (i >= cards.length) return;
  const card = document.getElementById(cards.length - i - 1);
  card.classList.add("swipe-right-fade");
  showInterest(cards[cards.length - i - 1].Title);
  if (i < cards.length - 1) {
    i += 1;
  } else {
    removeControls();
    fadeInOut("Sounds Fun!", "#73f564");
    return;
  }
  revealCard();
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

const showInterest = async (event) => {
  const config = {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      eventName: event,
    }),
  };
  const res = await fetch("/api/interested", config);
  console.log(event, res);
};

const revealCard = () => {
  id = cards.length - i - 1;
  const card = document.getElementById(id);
  card.style.opacity = 1;
};

const removeControls = () => {
  const controls = document.getElementById("event-controls");
  controls.style.display = "none";
  updateBackgound("");
  //show message
  const message = document.getElementById("no-events");
  message.style.opacity = 1;
};

const revealControls = () => {
  const controls = document.getElementById("event-controls");
  controls.style.display = "default";
};
