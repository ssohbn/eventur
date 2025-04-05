var cards = [
  {
    Title: "Event 2",
    Blurb: "This is the second event.",
    img: "https://images.pexels.com/photos/1105666/pexels-photo-1105666.jpeg",
  },
  {
    Title: "Event 1",
    Blurb: "This is the first event.",
    img: "https://as2.ftcdn.net/v2/jpg/04/96/15/83/1000_F_496158338_SgDd7OQQC2QVfN7U5Qijl2muktM0LjjG.jpg",
  },
];
var i = 0;

const updateBackgound = (img) => {
  document.body.style.backgroundImage = `url(${img})`;
};
updateBackgound(cards[i].img);

const swipeLeft = () => {
  const card = document.getElementById(cards[i].Title);
  card.classList.add("swipe-left-fade");
  i += 1;
  updateBackgound(cards[i].img);
};

const swipeRight = () => {
  const card = document.getElementById(cards[i].Title);
  card.classList.add("swipe-right-fade");
  i += 1;
  updateBackgound(cards[i].img);
};
