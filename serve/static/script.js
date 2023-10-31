function updateTime() {
   // The timezone is always zero UTC offset
   const currentDate = new Date()
   document.getElementById('gmtTime').innerText = currentDate.toISOString().slice(11, 19) + ' GMT'
   document.getElementById('zuluTime').innerText = toZuluTimeFormat(currentDate)
}
function toZuluTimeFormat(date) {
   const day = String(date.getUTCDate()).padStart(2, '0');
   const hours = String(date.getUTCHours()).padStart(2, '0');
   const minutes = String(date.getUTCMinutes()).padStart(2, '0');

   return `${day}${hours}${minutes}Z`;
}

updateTime();
setInterval(updateTime, 1000);