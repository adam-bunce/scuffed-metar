const metarPrintSettings = `
.airportInfo, [id$="-section"], [id$="-divider"], #print-dialog, #header, #footer, img, .camera-container, .jump-to-top {
    display: none;
} 
.spacing-div {
    margin: 0;
}
body {
    font-size: 10px;
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
body {
    font-size: 10px;
}

.notamInfo {
    line-height: 1.5em;
}
`


const App = {
    $: {
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

        // selection items (notam/winds)
        redirectbutton: document.getElementById("submission-redirect-button"),
        submissionResetButton: document.getElementById("submission-reset-button"),
        submissionSelectAllButton: document.getElementById("submission-select-all-button"),
        selectOpts: null,
        selectedSelectOptsIds: [],
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
            if (App.$.selectedPrintItemIds.includes(item.id)) item.innerHTML = "x"
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
    setSelectionItems() {
        App.$.selectOpts= document.querySelectorAll("[id$='-selection-opt']")
    },
    toggleSelectionOption(id) {
        if (App.$.selectedSelectOptsIds .includes(id)) App.$.selectedSelectOptsIds = App.$.selectedSelectOptsIds.filter(itemId => itemId !== id )
        else App.$.selectedSelectOptsIds.push(id)
        App.updateSelectedSelectionOptionsUI()
    },
    updateSelectedSelectionOptionsUI() {
        for (let item of App.$.selectOpts) {
            if (App.$.selectedSelectOptsIds.includes(item.id)) item.classList.add("selected")
            else item.classList.remove("selected")
        }
    },
    selectAllSelectionOptions() {
        App.$.selectedSelectOptsIds = []
        for (let item of App.$.selectOpts) {
            App.$.selectedSelectOptsIds.push(item.id)
        }
        App.updateSelectedSelectionOptionsUI()
    },
    clearSelectionOptions() {
        App.$.selectedSelectOptsIds = []
        App.updateSelectedSelectionOptionsUI()
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

        // selection settings
        App.$.selectOpts.forEach((el) => {App.bindClickEvent(el, () => App.toggleSelectionOption(el.id))})
        App.bindClickEvent(App.$.submissionSelectAllButton, () => App.selectAllSelectionOptions())
        App.bindClickEvent(App.$.submissionResetButton, () => App.clearSelectionOptions())
        App.bindClickEvent(App.$.redirectbutton, () => {
            const queryString = App.$.selectedSelectOptsIds.reduce((prev, curr) => {
                return `${prev}&airport=${curr.replace(/-selection-opt/g, '')}`
            }, "?")
            App.$.redirectbutton.href = App.$.redirectbutton.href + queryString
        })
    },
    init() {
        App.setPrintCheckboxes()
        App.setSelectionItems()
        App.bindEvents()

        window.onafterprint = () => App.resetPrintStyles()
        window.onload = () => {
            // can't print until page is fully loaded, after load unlock print button
            try {
                App.$.printDialogPrintButton.disabled = false
                App.$.printDialogInfoText.innerText = "Can Print."
            } catch (e) {
               // do nothing :)
            }
        }
    },
}

App.init()