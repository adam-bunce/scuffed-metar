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

        // include upper winds checkbox
        upperwindsCheckbox: document.getElementById('winds-checkbox'),
        upperwindsCheckboxSelected: "false",

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
    updateSubmitUrlSites(event) {
        let sites = event.target.value
        let sites_array = sites.trim().split(" ")
        let baseUrl = window.location.href.split('?')[0] // domain to update
        let queryParams = sites_array.map(code => `airport=${code.toUpperCase()}`).join('&')
        queryParams += `&winds=${Trip.$.upperwindsCheckboxSelected}`

        Trip.$.submitButton.href = `${baseUrl}?${queryParams}`
    },
    updateSubmitUrlQueryParams(include_bool) {
        let url_with_sites = Trip.$.submitButton.href.split("&")
        let base_url = url_with_sites[0]
        if (base_url.length <= 0) {
            base_url = window.location.href.split('?')[0] // domain to update
        }
        const new_url =  base_url + (include_bool === "true" ? "&winds=true" : "&winds=false")
        Trip.$.submitButton.href = new_url
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
}`

        let printStyleEl = document.createElement('style');
        printStyleEl.id = "print-formatting"
        printStyleEl.innerHTML = printStyle;
        document.head.appendChild(printStyleEl);

        window.print()

    },
    bindEvents() {
        // search
        Trip.bindInputEvent(Trip.$.inputBox, (event) => Trip.updateSubmitUrlSites(event))
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
    },
    resetPrintStyles() {
        let printStyleEl = document.getElementById("print-formatting")
        if (printStyleEl) printStyleEl.remove();
    },
    init() {
        Trip.setPrintCheckboxes()
        Trip.bindEvents()
        // Trip.getWindsCookie()

        Trip.bindClickEvent(Trip.$.upperwindsCheckbox, () => {
            if (Trip.$.upperwindsCheckboxSelected === "false") {
               Trip.$.upperwindsCheckbox.innerHTML = "x"
                Trip.$.upperwindsCheckboxSelected = "true"
            } else {
                Trip.$.upperwindsCheckbox.innerHTML = " "
                Trip.$.upperwindsCheckboxSelected = "false"
            }
            Trip.updateSubmitUrlQueryParams(Trip.$.upperwindsCheckboxSelected)
        })

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