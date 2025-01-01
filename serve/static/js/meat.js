
const Meat = {
    $: {
        title: document.getElementsByTagName("title").item(0),
        titleHeader: document.getElementById("top"),
    },
    updateText() {
        if (window.location.origin.toLowerCase().includes("meat")) {
            const pageTitle = this.$.title.innerText.split("|")
            this.$.title.innerText = pageTitle[0] + " | " + "Scuffed MEAT"
            this.$.titleHeader.innerText = "Scuffed MEAT"
        }
    },
    init() {
        this.updateText()
    }
}

Meat.init()