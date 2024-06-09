function addHoverListener(element, enter, leave) {
    element.addEventListener('mouseenter', enter)
    element.addEventListener('mouseleave', leave)
}

function copyToClipboard(item) {
    var txt = document.getElementById(item);
    txt.select();
    txt.setSelectionRange(0, 99999);
    navigator.clipboard.writeText(txt.value);
    alert("Copied the text: " + txt.value);
};


const hamburger = {
    onClick: function() {
        const list = document.getElementsByClassName("header-right-list")[0];
        list.classList.toggle("nav-burger-on");
        document.getElementById("nav-burger").textContent =
            list.classList.contains("nav-burger-on") ? "[X]" : "[ ]";
    }
}

const imageTask = {
    catalogue: [
        //        { uri: "/static/img/stream/earth-stream.png", frames: 36 },
        { uri: "/static/img/stream/computer-stream.png", frames: 15 },
        //       { uri: "/static/img/stream/snake-stream.png", frames: 44 }
    ],
    img: "",
    block: document.getElementById("text-block"),
    index: 0,
    cursor: 0,
    getSprite: (index) => imageTask.img.substring(index * imageTask.cursor,
        (index + 1) * imageTask.cursor),
    setup: function(index) {
        const selected = imageTask.catalogue[index];
        getAsciiImage(selected.uri, {
            maxHeight: window.screen.height / 3,
            maxWidth: window.screen.width / 2,
        }).then(img => {
            imageTask.img = img;
            imageTask.cursor = img.length / selected.frames;
            setInterval(() => imageTask.block
                .textContent = imageTask
                    .getSprite(imageTask.index++ % selected.frames), 120);
        }, () => { });
    }
}

const exec = {
    text: "",
    table: {},
    hash: (text) => text.replace(/[\n\t\s]/g, '').substring(0, 15),
    mailMatch: function(key) {
        const emailMatch = exec.table[key].match(/mailto:([^?]+)/);
        const subjectMatch = exec.table[key].match(/subject=([^&]+)/);
        if (emailMatch && subjectMatch) {
            const email = emailMatch[1];
            const subject = subjectMatch[1];
            exec.table[key] = (`mail -S "${subject}" ${email}`);
        }
    },
    sublistMatch: function(sublist, dir) {
        const subAnchor = sublist.querySelector('a');
        exec.table[exec.hash(subAnchor.textContent)] =
            "cd " + dir + "/" + subAnchor.textContent.trim();
    },
    listMatch: function(nav) {
        const anchorText = nav.querySelector('a').textContent;
        exec.table[exec.hash(anchorText)] = "cd " + anchorText.trim();
        Array.from(nav.getElementsByClassName("sublist-item"))
            .forEach(sublist => exec.sublistMatch(sublist, anchorText.trim()));
    },
    processText: function() {
        const field = document.getElementById("exec-text");
        const old = field.textContent;
        const block = document.getElementsByClassName("header-left")[0]
            .querySelector('span');
        if (exec.text == old) {
            if (!block.classList.contains("blinking-cursor")) {
                block.classList.add("blinking-cursor");
            }
            return;
        }
        if (block.classList.contains("blinking-cursor")) {
            block.classList.remove("blinking-cursor");
            return;
        }
        if (exec.text.length < old.length) {
            /* if oldText is larger than new */
            field.textContent = old.substring(0, old.length - 1);
        } else if (exec.text.length >= old.length &&
            (old.length == 0 ||
                exec.text.substring(0, old.length) == old)) {
            /* if oldText is smaller or equal and share common text with new */
            field.textContent += exec.text[old.length];
        } else if (exec.text != old) {
            field.textContent = old.substring(0, old.length - 1);
        }
    },
    setup: function() {
        const setTaskText = (item) => exec.text =
            exec.table[exec.hash(item.textContent)];
        const clearTask = () => exec.text = "";

        document.querySelectorAll('a').forEach(anchor => {
            exec.table[exec.hash(anchor.textContent)] = "curl " + anchor.href;
            addHoverListener(anchor, () => setTaskText(anchor), clearTask);
        });

        const hostname = document.getElementsByClassName("header-left")[0];
        exec.table[exec.hash(hostname.querySelector('a').textContent)] = "cd ~";

        Array.from(document.getElementsByClassName("nav-item"))
            .forEach(exec.listMatch);

        Object.keys(exec.table).forEach(exec.mailMatch);
        setInterval(exec.processText, 45);
    }
}

window.addEventListener('load', exec.setup);
switch (window.location.pathname) {
    case "/":
        window.addEventListener('load', imageTask.
            setup(Math.floor(Math.random() * imageTask.catalogue.length)));
}


