window.addEventListener('load', function() {
    const navToggle = document.querySelector('.nav-toggle');
    const nav = document.querySelector('.nav');
    navToggle.addEventListener('click', () => nav.classList.toggle('active'));
    var currentUrl = window.location.pathname;
    if (currentUrl === '/') {
        home();
    };
});

function copyToClipboard(item) {
    var txt = document.getElementById(item);
    txt.select();
    txt.setSelectionRange(0, 99999);
    navigator.clipboard.writeText(txt.value);
    alert("Copied the text: " + txt.value);
};

function home() {
    var projectDrawer = document.getElementById('project-drawer');
    var summary = document.getElementById('project-drawer-summary');
    var arrows = [
        document.getElementById('project-drawer-img-up'),
        document.getElementById('project-drawer-img-down')
    ];
    projectDrawer.addEventListener('toggle', () => {
        window.scrollTo(0, document.body.scrollHeight);
        summary.innerText = "Show " + (projectDrawer.open ? "Less" : "More");
        arrows.forEach((arrow) => arrow.classList.toggle("hidden"));
    });
};
