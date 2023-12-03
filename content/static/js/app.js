function setNavbar() {
    document.addEventListener('DOMContentLoaded', () => {
        const navToggle = document.querySelector('.nav-toggle');
        const nav = document.querySelector('.nav');
        navToggle.addEventListener('click',
            () => nav.classList.toggle('active'));
    });
}
setNavbar();

function copyToClipboard(item) {
    var txt = document.getElementById(item);
    txt.select();
    txt.setSelectionRange(0, 99999);
    navigator.clipboard.writeText(txt.value);
    alert("Copied the text: " + txt.value);
}
