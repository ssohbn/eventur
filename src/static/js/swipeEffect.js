var id = "Event 2";

const swipeLeft = () => {
  const card = document.getElementById(id);
  card.classList.add("swipe-left-fade");
  id = "Event 1";
};
const swipeRight = () => {
  const card = document.getElementById(id);
  card.classList.add("swipe-right-fade");
  id = "Event 1";
};

const updateBackgound = () => {
    
}
