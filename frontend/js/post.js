
let menu = document.querySelector('#menu-icon');
let navbar = document.querySelector('.content .right_side');
let searchbar = document.querySelector('.search_bar')

menu.onclick = () => {
	menu.classList.toggle('menu');
	navbar.classList.toggle('open');
    searchbar.classList.toggle('open')
}