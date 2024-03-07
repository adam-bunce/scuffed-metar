const darkmodeCSS = `
:root {
    --black: #D8D8D8;
    --grey: #9F9F9F;
    --otherGrey: #D8D8D8;
    --greyog: #808080;
    --white: #151515;
    --blue: #091F58;
    --red: #DD4E3E;
}

.red {
    background-color: var(--red);
    color: var(--black);
}

.green {
    background-color: var(--blue);
    color: var(--black);
}

.btn {
    border: 1px solid var(--greyog) !important;
    color: var(--otherGrey);
}
button {
    background-color: var(--white);
    color: var(--black);
}

.shadow {
    box-shadow: 2px 2px var(--greyog);
}

.current-page {
    background-color: var(--blue);
    color: var(--otherGrey);
}

.version {
    border: 1px solid var(--greyog);
    background-color: var(--blue);
    color: var(--otherGrey);
}

dialog {
    border: 1px solid var(--greyog);
    background-color: var(--white);
    color: var(--black);
}

.jump-to-top {
    border: 1px solid var(--greyog);
    color: var(--black)
}
`



const DarkMode = {
    $: {
        isDarkMode: false,
        darkmodeButton: document.getElementById("darkmode-toggle"),
    },
    bindClickEvent(element, func) {
        if (element) element.addEventListener('click', func)
    },
    getThemeCookie() {
        const cookies = document.cookie.split("=")
        const themeIndex = cookies.indexOf("theme")
        if (themeIndex === -1) {
            DarkMode.$.isDarkMode = false
        } else {
            DarkMode.$.isDarkMode = cookies[themeIndex + 1] === "true";
        }
    },
    setTheme(){
        if (DarkMode.$.isDarkMode) {
            let darkmodeEl = document.createElement('style')
            darkmodeEl.id = 'darkmode'
            darkmodeEl.innerHTML = darkmodeCSS
            document.head.appendChild(darkmodeEl)
        } else {
            const darkmodeEl= document.getElementById("darkmode")
            if (darkmodeEl) darkmodeEl.remove()
        }
        document.cookie = `theme=${DarkMode.$.isDarkMode}; path=/; max-age=2630000`;
    },
    init() {
        DarkMode.getThemeCookie()
        DarkMode.setTheme()

        DarkMode.bindClickEvent(DarkMode.$.darkmodeButton, () => {
            DarkMode.$.isDarkMode = !DarkMode.$.isDarkMode
            DarkMode.setTheme()
        })
    }
}

DarkMode.init()