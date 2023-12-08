const metarPrintSettings = `
.airportInfo, [id$="-section"], [id$="-divider"], #print-dialog, #header, #footer, img, .camera-container, .jump-to-top {
    display: none;
} 
.spacing-div {
    margin: 0;
}
`

const gfaPrintSettings = `
dialog, #header, #footer, img, .spacing-div {
    display: none;
} 
blockquote {
    margin: 0px;
}
.camera-container {
    display: block;
}
h3.mono {
    display: none;
}
.shadow {
    box-shadow: none;
}
`

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

const App = {
    $: {
        // time
        gmtTime: document.getElementById('gmtTime'),
        zuluTime: document.getElementById('zuluTime'),

        // dialog visibility
        printDialog: document.getElementById("print-dialog"),
        printDialogOpenTrigger: document.getElementById("print-dialog-open-trigger"),
        printDialogCloseTrigger: document.getElementById("print-dialog-close-trigger"),

        // dialog functionality
        printDialogInfoText: document.getElementById("print-dialog-info-text"),
        printDialogPrintButton: document.getElementById("print-button"),
        printDialogResetButton: document.getElementById("reset-button"),
        printDialogSelectAllButton: document.getElementById("select-all-button"),
        printOptionCheckboxes: null,
        selectedPrintItemIds: [],

        // darkmode
        isDarkMode: false,
        darkmodeButton: document.getElementById("darkmodeToggle"),

        // notam redirect button
        notamRedirectbutton: document.getElementById("notamredirect")
    },
    toZuluTimeFormat(date) {
        const day = String(date.getUTCDate()).padStart(2, '0')
        const hours = String(date.getUTCHours()).padStart(2, '0')
        const minutes = String(date.getUTCMinutes()).padStart(2, '0')

        return `${day}${hours}${minutes}Z`
    },
    updateTime() {
        // The timezone is always zero UTC offset
        const currentDate = new Date()
        App.$.gmtTime.innerText = currentDate.toISOString().slice(11, 19) + ' GMT'
        App.$.zuluTime.innerText = App.toZuluTimeFormat(currentDate)
    },

    setPrintCheckboxes() {
        App.$.printOptionCheckboxes = document.querySelectorAll("[id$='-print-checkbox']")
    },
    togglePrintOptionCheckbox(id) {
        if (App.$.selectedPrintItemIds.includes(id)) App.$.selectedPrintItemIds = App.$.selectedPrintItemIds.filter(itemId => itemId !== id )
        else App.$.selectedPrintItemIds.push(id)
        App.updateSelectedItemsUI()
    },
    updateSelectedItemsUI() {
        for (let item of App.$.printOptionCheckboxes) {
            if (App.$.selectedPrintItemIds.includes(item.id)) item.innerHTML = "X"
            else item.innerHTML = ""
        }
    },
    selectAllPrintOptionsCheckbox() {
        App.$.selectedPrintItemIds = []
        for (let item of App.$.printOptionCheckboxes) {
            App.$.selectedPrintItemIds.push(item.id)
        }
        App.updateSelectedItemsUI()
    },
    clearPrintOptionsCheckboxes() {
        App.$.selectedPrintItemIds = []
        App.updateSelectedItemsUI()
    },
    printSelectedItem() {
        // remove so i can print notams
//         if (App.$.selectedPrintItemIds.length === 0) return
        const useMetarSettings = String(App.$.selectedPrintItemIds[0]).includes("C") // all airport id's will start with a C

        let selectedElements = ""
        if (useMetarSettings) {
            selectedElements = App.$.selectedPrintItemIds.reduce((prev, curr, index) => {
                return index === 0 ? `#${curr.replace(/-print-checkbox/g, '')}-section` : `${prev}, #${curr.replace(/-print-checkbox/g, '')}-section`
            }, '')
        } else {
            selectedElements = App.$.selectedPrintItemIds.reduce((prev, curr, index) => {
                return index === 0 ? `#print-${curr.replace(/-print-checkbox/g, '')}` : `${prev}, #print-${curr.replace(/-print-checkbox/g, '')}`
            }, '')
        }

        const printStyle = `
@media print{
    :root {
        --width: 100%;
        --black: #000000;
    }
    ${useMetarSettings ? metarPrintSettings : gfaPrintSettings}

    ${selectedElements} {
         padding-top: 10px;
         display: block !important;   
         ${useMetarSettings ? 'border-bottom: 2px dotted var(--black);margin: 0 0 2em 0;' : ''}
         ${useMetarSettings ? '' : 'width: 75% !important'}
    }
}`

        let printStyleEl = document.createElement('style');
        printStyleEl.id = "print-formatting"
        printStyleEl.innerHTML = printStyle;
        document.head.appendChild(printStyleEl);

        window.print()
    },
    resetPrintStyles() {
        let printStyleEl = document.getElementById("print-formatting")
        if (printStyleEl) printStyleEl.remove();
    },

    startTimeCycle() {
        App.updateTime()
        setInterval(App.updateTime, 1000)
    },
    bindClickEvent(element, func) {
        if (element) element.addEventListener('click', func)
    },
    getCookies() {
        const cookies = document.cookie.split("=")
        const themeIndex = cookies.indexOf("theme")
        if (themeIndex === -1) {
            App.$.isDarkMode = false
        } else {
            App.$.isDarkMode = cookies[themeIndex + 1] === "true";
        }
    },
    setTheme() {
        if (App.$.isDarkMode) {
            let darkmodeEl = document.createElement('style')
            darkmodeEl.id = 'darkmode'
            darkmodeEl.innerHTML = darkmodeCSS
            document.head.appendChild(darkmodeEl)
        } else {
            const darkmodeEl= document.getElementById("darkmode")
            if (darkmodeEl) darkmodeEl.remove()
        }
        document.cookie = `theme=${App.$.isDarkMode}; path=/; max-age=2630000`;
    },
    bindEvents() {
        // dialog visibility
        App.bindClickEvent(App.$.printDialogOpenTrigger, () => App.$.printDialog.showModal())
        App.bindClickEvent(App.$.printDialogCloseTrigger, () => App.$.printDialog.close())

        // dialog print settings
        App.$.printOptionCheckboxes.forEach((el) => {App.bindClickEvent(el, () => App.togglePrintOptionCheckbox(el.id))})
        App.bindClickEvent(App.$.printDialogPrintButton, () => App.printSelectedItem())
        App.bindClickEvent(App.$.printDialogSelectAllButton, () => App.selectAllPrintOptionsCheckbox())
        App.bindClickEvent(App.$.printDialogResetButton, () => App.clearPrintOptionsCheckboxes())

        // darkmode
        App.bindClickEvent(App.$.darkmodeButton, () => {App.$.isDarkMode = !App.$.isDarkMode; App.setTheme()})

        /// notam redirect
        App.bindClickEvent(App.$.notamRedirectbutton, () => {
            console.log(App.$.selectedPrintItemIds)
                // selectedElements = App.$.selectedPrintItemIds.reduce((prev, curr, index) => {
            const queryString = App.$.selectedPrintItemIds.reduce((prev, curr) => {
                // return index === 0 ? `#print-${curr.replace(/-print-checkbox/g, '')}` : `${prev}, #print-${curr.replace(/-print-checkbox/g, '')}`
                return `${prev}&airport=${curr.replace(/-print-checkbox/g, '')}`
            }, "?")
            console.log(queryString)
            App.$.notamRedirectbutton.href = "/notam" + queryString
        })

    },
    init() {
        App.getCookies()
        App.setTheme()

        App.setPrintCheckboxes()
        App.bindEvents()
        App.startTimeCycle()

        window.onafterprint = () => App.resetPrintStyles()
        window.onload = () => {
            // can't print until page is fully loaded, after load unlock print button
            App.$.printDialogPrintButton.disabled = false
            App.$.printDialogInfoText.innerText = "Can Print."
        }
    },
}

App.init()