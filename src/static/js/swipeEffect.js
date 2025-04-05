let cards;
let i = 0;

const init = async () => {
  const res = await fetch("/api/events");
  cards = await res.json();
  console.log(cards);
  updateBackgound(cards[i].img);
}
init();

const updateBackgound = (img) => {
  document.body.style.backgroundImage = `url(${img})`;
};

const swipeLeft = () => {
  if (i >= cards.length)
    return;
  const card = document.getElementById(cards.length - i - 1);
  console.log(card);
  card.classList.add("swipe-left-fade");
  i += 1;
  updateBackgound(cards[i].img);
};

const swipeRight = () => {
  if (i >= cards.length)
    return;
  const card = document.getElementById(cards.length - i - 1);
  console.log(card);
  card.classList.add("swipe-right-fade");
  i += 1;
  updateBackgound(cards[i].img);
};
