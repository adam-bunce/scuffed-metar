const defaultPrintSettings= `
/* obvi stuff*/

#metars .spacing-div:not(:last-of-type) {
    border-bottom: 2px dotted var(--black);
    margin: 0 0 2em 0;
    display: block;
}

dialog, #header, #footer, .spacing-div, img, [id$="-gfa-header"] {
    display: none;
} 
.camera-container {display:block}

body {
    font-size: 10px;
}

.notamInfo {
    line-height: 1.5em;
}
`

const Trip = {
    $: {
        // search
        inputBox: document.getElementById('site-input'),
        submitButton: document.getElementById('submission-redirect-button'),

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

        notamPrintCheckbox: null,
    },
    bindInputEvent(element, func) {
        if (element) element.addEventListener('input', func)
    },
    bindEnterEvent(element, func) {
        if (element) element.addEventListener('keydown', func)
    },
    bindClickEvent(element, func) {
        if (element) element.addEventListener('click', func)
    },
    updateSubmitUrl(event) {
        let sites = event.target.value
        let sites_array = sites.trim().split(" ")
        let baseUrl = window.location.href.split('?')[0] // domain to update
        let queryParams = sites_array.map(code => `airport=${code.toUpperCase()}`).join('&')

        Trip.$.submitButton.href = `${baseUrl}?${queryParams}`
    },
    enterSubmit(event) {
        if (event.key === 'Enter') {
            Trip.$.submitButton.click()
        }
    },
    setPrintCheckboxes() {
        Trip.$.printOptionCheckboxes = document.querySelectorAll("[id$='-print-checkbox']")
        Trip.$.notamPrintCheckbox = document.getElementById("notam-print-toggle")
    },
    togglePrintOptionCheckbox(id) {
        if (Trip.$.selectedPrintItemIds.includes(id)) Trip.$.selectedPrintItemIds = Trip.$.selectedPrintItemIds.filter(itemId => itemId !== id)
        else Trip.$.selectedPrintItemIds.push(id)
        Trip.updateSelectedItemsUI()
    },
    toggleNotamCheckbox() {
        if (Trip.$.notamPrintCheckbox.innerHTML === "x") {
            Trip.$.notamPrintCheckbox.innerHTML = ""
        } else Trip.$.notamPrintCheckbox.innerHTML = "x"
    },
    updateSelectedItemsUI() {
        for (let item of Trip.$.printOptionCheckboxes) {
            if (Trip.$.selectedPrintItemIds.includes(item.id)) item.innerHTML = "x"
            else item.innerHTML = ""
        }
    },
    selectAllPrintOptionsCheckbox() {
        Trip.$.selectedPrintItemIds = []
        for (let item of Trip.$.printOptionCheckboxes) {
            Trip.$.selectedPrintItemIds.push(item.id)
        }
        Trip.updateSelectedItemsUI()
    },
    clearPrintOptionsCheckboxes() {
        Trip.$.selectedPrintItemIds = []
        Trip.updateSelectedItemsUI()
    },
    printSelectedItem() {
        let selectedElements = ""
        selectedElements = Trip.$.selectedPrintItemIds.reduce((prev, curr, index) => {
            return index === 0 ? `#print-${curr.replace(/-print-checkbox/g, '')}` : `${prev}, #print-${curr.replace(/-print-checkbox/g, '')}`
        }, '')

        const printStyle = `
@media print{
    :root {
        --width: 100%;
        --black: #000000;
    }
    
    ${defaultPrintSettings}

    ${selectedElements} {
         padding-top: 10px;
         display: block !important;   
         width: 75% !important
    }
    
    .requested-at {
        display: block;
    }
    
    ${Trip.$.notamPrintCheckbox.innerHTML !== 'x' ? "#notam-section{display: none}"  : ""}
}`

        let printStyleEl = document.createElement('style');
        printStyleEl.id = "print-formatting"
        printStyleEl.innerHTML = printStyle;
        document.head.appendChild(printStyleEl);

        window.print()

    },
    bindEvents() {
        // search
        Trip.bindInputEvent(Trip.$.inputBox, (event) => Trip.updateSubmitUrl(event))
        Trip.bindEnterEvent(Trip.$.inputBox, (event) => Trip.enterSubmit(event))

        // dialog visibility
        Trip.bindClickEvent(Trip.$.printDialogOpenTrigger, () => Trip.$.printDialog.showModal())
        Trip.bindClickEvent(Trip.$.printDialogCloseTrigger, () => Trip.$.printDialog.close())

        // dialog print settings
        Trip.$.printOptionCheckboxes.forEach((el) => {
            Trip.bindClickEvent(el, () => Trip.togglePrintOptionCheckbox(el.id))
        })
        Trip.bindClickEvent(Trip.$.printDialogPrintButton, () => Trip.printSelectedItem())
        Trip.bindClickEvent(Trip.$.printDialogSelectAllButton, () => Trip.selectAllPrintOptionsCheckbox())
        Trip.bindClickEvent(Trip.$.printDialogResetButton, () => Trip.clearPrintOptionsCheckboxes())

        Trip.bindClickEvent(Trip.$.notamPrintCheckbox, () => Trip.toggleNotamCheckbox())
    },
    resetPrintStyles() {
        let printStyleEl = document.getElementById("print-formatting")
        if (printStyleEl) printStyleEl.remove();
    },
    init() {
        Trip.setPrintCheckboxes()
        Trip.bindEvents()

        window.onafterprint = () => Trip.resetPrintStyles()
        window.onload = () => {
            // can't print until page is fully loaded, after load unlock print button
            try {
                Trip.$.printDialogPrintButton.disabled = false
                Trip.$.printDialogInfoText.innerText = "Can Print."
            } catch (e) {
                // do nothing :)
            }
        }

    },
}

Trip.init()