const swipeLeft = (id) => {
  const card = document.getElementById(id);
  card.classList.add("swipe-left-fade");
};
const swipeRight = (id) => {
  const card = document.getElementById(id);
  card.classList.add("swipe-right-fade");
};
