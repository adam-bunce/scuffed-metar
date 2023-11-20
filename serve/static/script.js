const metarPrintSettings = `
.airportInfo, [id$="-section"], [id$="-divider"], #print-dialog, #header, #footer, img, .camera-container, .jump-to-top {
    display: none;
} 

.spacing-div {
    margin: 0;
}
`

const gfaPrintSettings = `
dialog, #header, #footer, img {
    display: none;
} 
.camera-container {
    display: block;
}
.shadow {
    box-shadow: none;
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
        if (App.$.selectedPrintItemIds.length === 0) return
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
    ${useMetarSettings ? metarPrintSettings : gfaPrintSettings}

    ${selectedElements} {
         padding-top: 10px;
         display: block !important;   
         width: 100vw !important;
         ${useMetarSettings ? 'border-bottom: 2px dotted var(--black);margin: 0 0 2em 0;' : ''}
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
    bindEvents() {
        // dialog visibility
        App.bindClickEvent(App.$.printDialogOpenTrigger, () => App.$.printDialog.showModal())
        App.bindClickEvent(App.$.printDialogCloseTrigger, () => App.$.printDialog.close())

        // dialog print settings
        App.$.printOptionCheckboxes.forEach((el) => {App.bindClickEvent(el, () => App.togglePrintOptionCheckbox(el.id))})
        App.bindClickEvent(App.$.printDialogPrintButton, () => App.printSelectedItem())
        App.bindClickEvent(App.$.printDialogSelectAllButton, () => App.selectAllPrintOptionsCheckbox())
        App.bindClickEvent(App.$.printDialogResetButton, () => App.clearPrintOptionsCheckboxes())
    },
    init() {
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