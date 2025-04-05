var cards = [
  {
    Title: "Midnight Concert",
    Blurb: "This is the first event.",
    img: "https://images.pexels.com/photos/1105666/pexels-photo-1105666.jpeg",
  },
  {
    Title: "Picnic at the park",
    Blurb: "This is the second event.",
    img: "https://as2.ftcdn.net/v2/jpg/04/96/15/83/1000_F_496158338_SgDd7OQQC2QVfN7U5Qijl2muktM0LjjG.jpg",
  },
  {
    Title: "Midnight Concert",
    Blurb: "This is the first event.",
    img: "",
  },
];
var i = 0;

const updateBackgound = (img) => {
  document.body.style.backgroundImage = `url(${img})`;
};
updateBackgound(cards[i].img);

const swipeLeft = () => {
  const card = document.getElementById(cards[i].Title);
  console.log(card);
  card.classList.add("swipe-left-fade");
  i += 1;
  updateBackgound(cards[i].img);
  fadeInOut("Not interested!", "#f56e64");
};

const swipeRight = () => {
  const card = document.getElementById(cards[i].Title);
  console.log(card);
  card.classList.add("swipe-right-fade");
  i += 1;
  updateBackgound(cards[i].img);
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
