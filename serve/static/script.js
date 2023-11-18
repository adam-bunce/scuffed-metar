const App = {
   $: {
      gmtTime: document.getElementById('gmtTime'),
      zuluTime: document.getElementById('zuluTime'),
      printButton: document.getElementById("print-button"),
      resetButton: document.getElementById("reset-button"),
      selectAllButton: document.getElementById("select-all-button"),
      gafCheckboxes: document.getElementsByClassName("checkbox"),
      selectedGafItems: [],
      setPrintStyles() {
         let targetIds = App.$.selectedGafItems.reduce((prev, curr, index) =>
         { return index === 0 ? `#img-${curr}` : `${prev}, #img-${curr}`}, '')
         console.log(targetIds)

         let printStyles = `
  @media print {
     #header, #footer {
     display: none;
      } 
      img {
          display: none;
      }
      .camera-container {
        display: block;
      }
      .main {
        width: 80%;
      }
    
    ${targetIds} {
      padding-top: 10px;
      display: block !important;   
      width: 100% !important;
    }
  }
`;

         console.log(printStyles)

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
      updateSelectedGafItems(itemId) {

      },
   },
   startTimeCycle() {
      App.$.updateTime()
      setInterval(App.$.updateTime, 1000)
   },
   bindEvents() {
      if (App.$.printButton) {
         App.$.printButton.addEventListener('click', function () {
            // NOTE might need to do something about print order due to the order they get added in
            // actually maybe not beacuse the css would jsut handle toggles not moving things
            console.log('printing images with ids: ', App.$.selectedGafItems)
            App.$.setPrintStyles()
            window.print()
         })
      }
      if (App.$.resetButton) {
         App.$.resetButton.addEventListener('click', function () {
            App.$.resetSelectedGafItems()
         })
      }
      if (App.$.selectAllButton) {
         App.$.selectAllButton.addEventListener('click', () => App.$.selectAllItems())
      }
      if (App.$.gafCheckboxes) {
         console.log(App.$.gafCheckboxes)
         for (let item of App.$.gafCheckboxes) {
            console.log(item.id);
            item.addEventListener('click', () => App.$.toggleSelectedItem(item.id))
         }
      }
   },
   init() {
      App.bindEvents()
      App.startTimeCycle()
      window.onafterprint = () => App.$.resetPrintStyles()
   },
}

App.init()