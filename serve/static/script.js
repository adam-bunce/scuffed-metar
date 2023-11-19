const App = {
   $: {
      gmtTime: document.getElementById('gmtTime'),
      zuluTime: document.getElementById('zuluTime'),

      printButton: document.getElementById("print-button"),
      resetButton: document.getElementById("reset-button"),
      selectAllButton: document.getElementById("select-all-button"),
      gafCheckboxes: document.getElementsByClassName("gaf-checkbox"),
       selectedGafItems: [],

      printDialog: document.getElementById("print-dialog"),
      printDialogOpenTrigger: document.getElementById("print-dialog-open-trigger"),
      printDialogCloseTrigger: document.getElementById("print-dialog-close-trigger"),
      metarSelectAllButton: document.getElementById("metar-select-all-button"),
      metarResetAllButton: document.getElementById("metar-reset-button"),
      metarPrintButton: document.getElementById("metar-print-button"),
      metarCheckboxes: document.getElementsByClassName("metar-checkbox"),
      selectedMetars: [],
      setPrintStyles() {
         let targetIds = App.$.selectedGafItems.reduce((prev, curr, index) =>
         { return index === 0 ? `#img-${curr}` : `${prev}, #img-${curr}`}, '')

         let printStyles = `
  @media print {
     #header, #footer, img {
         display: none;
      } 
      .camera-container {
        display: block;
      }
      
      .shadow {
        box-shadow: none;
      }
      
      /* h3 {display: none;} */
    
    ${targetIds} {      
      padding-top: 10px;
      display: block !important;   
      width: 100vw !important;
      /* height: 90vh;  */
      /* object-fit: contain; */
    }
  }
`;


         // Insert the print styles into the document
         let printStyleEl = document.createElement('style');
         printStyleEl.id = "printStyle"
         printStyleEl.innerHTML = printStyles;
         document.head.appendChild(printStyleEl);
      },
      setMetarPrintStyles() {
           let selectedAirports = App.$.selectedMetars.reduce((prev, curr, index) =>
           { return index === 0 ? `#${curr}-section` : `${prev}, #${curr}-section`}, '')

           let printStyles = `
  @media print {
     #print-dialog, #header, #footer, img, .camera-container, .shadow, .jump-to-top {
         display: none;
      } 
      
      .airportInfo {
        display: none;
      }
      
    [id$="-section"] {
        display: none;  
    }
     
    [id$="-divider"] {
        display: none;
    }
    
    .spacing-div {
        margin: 0;
    }
      
    ${selectedAirports} {      
      display: block;
      padding-top: 10px;
      display: block !important;   
      width: 100vw !important;
      /* height: 90vh;  */
      /* object-fit: contain; */
    }
  }
`;

           // Insert the print styles into the document
           let printStyleEl = document.createElement('style');
           printStyleEl.id = "printStyle"
           printStyleEl.innerHTML = printStyles;
           document.head.appendChild(printStyleEl);
       },
      resetPrintStyles() {
          let printStyleEl = document.getElementById("printStyle")
          if (printStyleEl) printStyleEl.remove();
      },
      toggleSelectedItem(id) {
         if (App.$.selectedGafItems.includes(id)) App.$.selectedGafItems = App.$.selectedGafItems.filter(itemid => itemid !== id )
         else App.$.selectedGafItems.push(id)

         App.$.updateSelectedItemsUI()
      },
      toggleSelectedMetarItem(id) {
           if (App.$.selectedMetars.includes(id)) App.$.selectedMetars= App.$.selectedMetars.filter(itemid => itemid !== id )
           else App.$.selectedMetars.push(id)

           App.$.updateSelectedMetarItemsUI()
       },
      updateSelectedItemsUI() {
         // update ui
         for (let item of App.$.gafCheckboxes) {
            if (App.$.selectedGafItems.includes(item.id)) item.innerHTML = "X"
            else item.innerHTML = ""
         }
      },
      selectAllItems() {
         // avoid duplicates
         App.$.resetSelectedGafItems()
         for (let item of App.$.gafCheckboxes) {
           App.$.selectedGafItems.push(item.id)
         }
         App.$.updateSelectedItemsUI()
      },
       updateSelectedMetarItemsUI() {
           // update ui
           for (let item of App.$.metarCheckboxes) {
               if (App.$.selectedMetars.includes(item.id)) item.innerHTML = "X"
               else item.innerHTML = ""
           }
       },
       selectAllMetarItems() {
           // avoid duplicates
           App.$.resetSelectedMetarItems()
           for (let item of App.$.metarCheckboxes) {
               App.$.selectedMetars.push(item.id)
           }
           App.$.updateSelectedMetarItemsUI()
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
         App.$.zuluTime.innerText = App.$.toZuluTimeFormat(currentDate)
      },
      resetSelectedGafItems() {
         App.$.selectedGafItems = []
         App.$.updateSelectedItemsUI()
      },
      resetSelectedMetarItems() {
          App.$.selectedMetars = []
          App.$.updateSelectedMetarItemsUI()
      },
   },
   startTimeCycle() {
      App.$.updateTime()
      setInterval(App.$.updateTime, 1000)
   },
   bindEvents() {
      if (App.$.printDialogOpenTrigger) {
          App.$.printDialogOpenTrigger.addEventListener('click', function () {
             if (App.$.printDialog) App.$.printDialog.show()
          })
      }
       if (App.$.printDialogCloseTrigger) {
           App.$.printDialogCloseTrigger.addEventListener('click', function () {
               if (App.$.printDialog) App.$.printDialog.close()
           })
       }
      if (App.$.printButton) {
         App.$.printButton.addEventListener('click', function () {
            // NOTE might need to do something about print order due to the order they get added in
            // actually maybe not beacuse the css would jsut handle toggles not moving things
            App.$.setPrintStyles()
            window.print()
         })
      }
       if (App.$.metarPrintButton) {
           App.$.metarPrintButton.addEventListener('click', function () {
               App.$.setMetarPrintStyles()
               window.print()
           })
       }
      if (App.$.resetButton) {
         App.$.resetButton.addEventListener('click', function () {
            App.$.resetSelectedGafItems()
         })
      }
      if (App.$.metarResetAllButton) {
          App.$.metarResetAllButton.addEventListener('click', function () {
              App.$.resetSelectedMetarItems()
          })
      }
      if (App.$.selectAllButton) {
         App.$.selectAllButton.addEventListener('click', () => App.$.selectAllItems())
      }
      if (App.$.metarSelectAllButton) {
          App.$.metarSelectAllButton.addEventListener('click', () => App.$.selectAllMetarItems())
      }
      if (App.$.gafCheckboxes) {
         for (let item of App.$.gafCheckboxes) {
            item.addEventListener('click', () => App.$.toggleSelectedItem(item.id))
         }
      }
     if (App.$.metarCheckboxes)  {
         for (let item of App.$.metarCheckboxes) {
             item.addEventListener('click', () => App.$.toggleSelectedMetarItem(item.id))
         }
     }

   },
   init() {
      App.bindEvents()
      App.startTimeCycle()
      window.onafterprint = () => App.$.resetPrintStyles()
       window.addEventListener('load', function () {
          if (App.$.printButton)  {
             App.$.printDialog.innerText = "Print"
          }

           if (App.$.metarPrintButton)  {
               App.$.metarPrintButton.innerText = "Print"
           }
       })
   },
}

App.init()