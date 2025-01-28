const defaultPrintSettings= `
/* obvi stuff*/

#metars .spacing-div:not(:last-of-type) {
    border-bottom: 2px dotted var(--black);
    display: block;
    padding: 0;
}
.source-link, dialog, #header, #footer, .spacing-div, img, [id$="-gfa-header"] {
    display: none;
} 
.camera-container {display:block}
.pt {padding-top: 0px;}
body { font-size: 10px; }
.notamInfo { line-height: 1.5em;}
.divider, .spacing-div {margin: 0 0 .5em 0;}
`

const Trip = {
    $: {
        // search
        inputBox: document.getElementById('site-input'),
        submitButton: document.getElementById('submission-redirect-button'),
        inputBoxContent: "",


        // include upper winds checkbox
        upperwindsCheckbox: document.getElementById('winds-checkbox'),
        upperwindsCheckboxSelected: "false",

        // include mets checkbox
        metsCheckbox: document.getElementById('mets-checkbox'),
        metsCheckboxSelected: "false",

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
        windPrintCheckbox: null,
        metsPrintCheckbox: null,
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
    buildSubmitUpdateUrl() {
        const sites =  Trip.$.inputBoxContent.trim().split(" ")
        const includeWinds =  Trip.$.upperwindsCheckboxSelected === "true"
        const includeMets =  Trip.$.metsCheckboxSelected === "true"

        let base_url = window.location.href.split('?')[0]

        Trip.$.submitButton.href = base_url
            + "?"
            + sites.reduce((prev, curr, idx) => idx === 0 ? `airport=${curr.toUpperCase()}` : prev + `&airport=${curr.toUpperCase()}`, "")
            + `&winds=${includeWinds}`
            + `&mets=${includeMets}`
    },
    enterSubmit(event) {
        if (event.key === 'Enter') {
            Trip.$.submitButton.click()
        }
    },
    setPrintCheckboxes() {
        Trip.$.printOptionCheckboxes = document.querySelectorAll("[id$='-print-checkbox']")
        Trip.$.notamPrintCheckbox = document.getElementById("notam-print-toggle")
        Trip.$.windPrintCheckbox = document.getElementById("winds-print-toggle")
        Trip.$.metsPrintCheckbox= document.getElementById("mets-print-toggle")
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
    toggleWindCheckbox() {
        if (Trip.$.windPrintCheckbox.innerHTML === "x") {
            Trip.$.windPrintCheckbox.innerHTML = ""
        } else Trip.$.windPrintCheckbox.innerHTML = "x"
    },
    toggleMetsCheckbox() {
        if (Trip.$.metsPrintCheckbox.innerHTML === "x") {
            Trip.$.metsPrintCheckbox.innerHTML = ""
        } else Trip.$.metsPrintCheckbox.innerHTML = "x"
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
    ${Trip.$.windPrintCheckbox.innerHTML !== 'x' ? "#wind-section{display: none}"  : ""}
    ${Trip.$.metsPrintCheckbox.innerHTML !== 'x' ? "#mets-section{display: none}"  : ""}
}`

        let printStyleEl = document.createElement('style');
        printStyleEl.id = "print-formatting"
        printStyleEl.innerHTML = printStyle;
        document.head.appendChild(printStyleEl);

        window.print()

    },
    bindEvents() {
        Trip.bindInputEvent(Trip.$.inputBox, () => Trip.$.inputBoxContent = event.target.value)
        Trip.bindInputEvent(Trip.$.inputBox, () => Trip.buildSubmitUpdateUrl())
        Trip.bindClickEvent(Trip.$.upperwindsCheckbox, () => Trip.buildSubmitUpdateUrl())
        Trip.bindClickEvent(Trip.$.metsCheckbox, () => Trip.buildSubmitUpdateUrl())

        // search
        // Trip.bindInputEvent(Trip.$.inputBox, (event) => Trip.updateSubmitUrlSites(event))
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
        Trip.bindClickEvent(Trip.$.windPrintCheckbox, () => Trip.toggleWindCheckbox())
        Trip.bindClickEvent(Trip.$.metsPrintCheckbox, () => Trip.toggleMetsCheckbox())
    },
    resetPrintStyles() {
        let printStyleEl = document.getElementById("print-formatting")
        if (printStyleEl) printStyleEl.remove();
    },
    init() {
        Trip.setPrintCheckboxes()
        // Trip.getWindsCookie()

        Trip.bindClickEvent(Trip.$.upperwindsCheckbox, () => {
            if (Trip.$.upperwindsCheckboxSelected === "false") {
               Trip.$.upperwindsCheckbox.innerHTML = "x"
                Trip.$.upperwindsCheckboxSelected = "true"
            } else {
                Trip.$.upperwindsCheckbox.innerHTML = " "
                Trip.$.upperwindsCheckboxSelected = "false"
            }
        })

        Trip.bindClickEvent(Trip.$.metsCheckbox, () => {
            if (Trip.$.metsCheckboxSelected=== "false") {
                Trip.$.metsCheckbox.innerHTML = "x"
                Trip.$.metsCheckboxSelected= "true"
            } else {
                Trip.$.metsCheckbox.innerHTML = " "
                Trip.$.metsCheckboxSelected = "false"
            }
        })

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