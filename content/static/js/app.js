window.addEventListener('load', function() {
    var currentUrl = window.location.pathname;

    if (currentUrl === '/') {
        home();
    }
});

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

function home() {
    var projectDrawer = document.getElementById('project-drawer');
    var summary = document.getElementById('project-drawer-summary');
    var img = document.getElementById('project-drawer-img');

    projectDrawer.addEventListener('toggle',
        () => {
            window.scrollTo(0, document.body.scrollHeight);
            summary.innerText = "Show " + (projectDrawer.open ? "Less" : "More");
            img.classList.toggle("project-drawer-animation")
        });
}


