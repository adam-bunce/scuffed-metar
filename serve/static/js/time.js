const Time = {
    $: {
        gmtTime: document.getElementById('gmtTime'),
        zuluTime: document.getElementById('zuluTime'),
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
        Time.$.gmtTime.innerText = currentDate.toISOString().slice(11, 19) + ' GMT'
        Time.$.zuluTime.innerText = Time.toZuluTimeFormat(currentDate)
    },
    startTimeCycle() {
        Time.updateTime()
        setInterval(Time.updateTime, 1000)
    },
    init() {
        Time.startTimeCycle()
    }
}

Time.init()