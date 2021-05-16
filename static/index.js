// INTERACTIVITY STUFF
let frontMenuEl = document.getElementById("front")
let appMenuEl = document.getElementById("app")
let settingsMenuEl = document.getElementById("settings")
let appsLinkEl = document.getElementById("apps-link")
let settingsLinkEl = document.getElementById("settings-link")
let lowerLimitEl = document.getElementById("lower-limit")
let cancelLowerLimitEl = document.getElementById("cancel-lower-limit")
let appNameEl = document.getElementById("app-name")
let limitGoalEl = document.getElementById("limit-goal")
let appLogoEl = document.getElementById("app-logo")
let appHomeEl = document.getElementById("app-home")
let settingsHomeEl = document.getElementById("settings-home")
let settingsApplyEl = document.getElementById("settings-apply-all")

let lastScreen = frontMenuEl
let currentScreen = frontMenuEl

let lastLowerLimit = false

for(let el of [...document.getElementsByClassName("input"), appNameEl])
	el.addEventListener("keydown", (e) => {
		if(e.keyCode == 13) 
			e.preventDefault()
	})

function setScreen(screen) {
	lastScreen = currentScreen
	switch (screen) {
		case "front":
			currentScreen = frontMenuEl
			break
		case "app":
			currentScreen = appMenuEl
			break
		case "settings":
			currentScreen = settingsMenuEl
			break
	}
	lastScreen.classList.remove("visible")
	lastScreen.classList.add("hiding")
	setTimeout(() => {
		currentScreen.classList.remove("hiding")
		currentScreen.classList.add("visible")
	}, 1)
}

function setLowerLimit(lowerLimit) {
	if(lowerLimit && !lastLowerLimit) {
		limitGoalEl.classList.remove("hiding")
		limitGoalEl.classList.add("visible")
		lowerLimitEl.classList.remove("visible")
		lowerLimitEl.classList.add("hiding")
	}
	else if(!lowerLimit && lastLowerLimit) {
		limitGoalEl.classList.remove("visible")
		limitGoalEl.classList.add("hiding")
		lowerLimitEl.classList.remove("hiding")
		lowerLimitEl.classList.add("visible")
	}
	lastLowerLimit = lowerLimit
}

settingsLinkEl.addEventListener("click", (e) => {
	e.preventDefault()
	setScreen("settings")
})

appHomeEl.addEventListener("click", (e) => {
	e.preventDefault()
	setScreen("front")
})

settingsHomeEl.addEventListener("click", (e) => {
	e.preventDefault()
	setScreen("front")
})

// CONFIG STUFF
let config = {}
let usedUnitEl = document.getElementById("used-unit")
let limitUnitEl = document.getElementById("limit-unit")
let cycleUnitEl = document.getElementById("cycle-unit")
let goalUnitEl = document.getElementById("goal-unit")
let goalLeftUnitEl = document.getElementById("goal-left-unit")
let cycleLengthEl = document.getElementById("cycle-length")
let measureUnitEl = document.getElementById("measure-unit")
let settingsResetAllEl = document.getElementById("settings-reset-all")

function loadConfig() {
	return new Promise((resolve, reject) => {
		fetch("./info")
			.then((body) => {
				return body.json()
			})
			.then((data) => {
				config.cycle_start = data.lastReset
				config.cycle_length = data.resetTime
				config.unit = data.unit
				switch(config.cycle_length) {
					case 2592000:
						config.cycle_length_name = "month"
						cycleUnitEl.innerHTML = " monthly "
						goalLeftUnitEl.innerHTML = " months "
						cycleLengthEl.selectedIndex = 2
						break
					case 86400:
						config.cycle_length_name = "day"
						cycleUnitEl.innerHTML = " daily "
						goalLeftUnitEl.innerHTML = " days "
						cycleLengthEl.selectedIndex = 1
						break
					default:
						config.cycle_length_name = "week"
						cycleUnitEl.innerHTML = " weekly "
						goalLeftUnitEl.innerHTML = " weeks "
						cycleLengthEl.selectedIndex = 0
						break
				}
				switch(config.unit) {
					case 1:
						config.unit_name = "second"
						usedUnitEl.innerHTML = " seconds "
						limitUnitEl.innerHTML = " second "
						goalUnitEl.innerHTML = " seconds "
						measureUnitEl.selectedIndex = 2
						break
					case 60:
						config.unit_name = "minute"
						usedUnitEl.innerHTML = " minutes "
						limitUnitEl.innerHTML = " minute "
						goalUnitEl.innerHTML = " minutes "
						measureUnitEl.selectedIndex = 1
						break
					default:
						config.unit_name = "hour"
						usedUnitEl.innerHTML = " hours "
						limitUnitEl.innerHTML = " hour "
						goalUnitEl.innerHTML = " hours "
						measureUnitEl.selectedIndex = 0
						break
				}
				resolve()
			})
			.catch((e) => {
				reject(e)
			})
	})
}

loadConfig()

function setConfig(config) {
	for(let [name, value] of Object.entries(config)) {
		switch(name) {
			case "cycle_length_name":
				switch(value) {
					case "month":
						fetch("./command?cyclelength+2592000")
						break
					case "day":
						fetch("./command?cyclelength+86400")
						break
					default:
						fetch("./command?cyclelength+604800")
						break
				}
				break
			case "unit_name":
				switch(value) {
					case "second":
						fetch("./command?unit+1")
						break
					case "minute":
						fetch("./command?unit+60")
						break
					default:
						fetch("./command?unit+3600")
						break
				}
			default:
				console.log("config not found")
		}
	}
	loadConfig()
}

cycleLengthEl.addEventListener("change", (e) => {
	switch(cycleLengthEl.selectedIndex) {
		case 0:
			setConfig({cycle_length_name: "week"})
			break
		case 1:
			setConfig({cycle_length_name: "day"})
			break
		case 2:
			setConfig({cycle_length_name: "month"})
			break
	}
})

measureUnitEl.addEventListener("change", (e) => {
	switch(measureUnitEl.selectedIndex) {
		case 0:
			setConfig({unit_name: "hour"})
			break
		case 1:
			setConfig({unit_name: "minute"})
			break
		case 2:
			setConfig({unit_name: "second"})
			break
	}
})

settingsResetAllEl.addEventListener("click", (e)=>{
	e.preventDefault()
	if(!confirm("Are you sure you want to reset all settings?"))
		return
	for(let el of document.getElementsByTagName("select")) {
		el.selectedIndex = 0
		el.dispatchEvent(new Event("change"))
	}
})

// APP MENU STUFF
let appIndex = 0
let applications = []

let usedTimeEl = document.getElementById("used-time")
let limitTimeEl = document.getElementById("limit-time")
let goalTimeEl = document.getElementById("goal-time")
let goalLeftTimeEl = document.getElementById("goal-left-time")
let processNameEl = document.getElementById("process-name")
let appNextEl = document.getElementById("app-next")
let appBackEl = document.getElementById("app-back")
let appDeleteEl = document.getElementById("app-delete")

function refresh() {
	return new Promise((resolve, reject) => {
		fetch("./command?info")
			.then((body) => {
				return body.json()
			})
			.then(async (data) => {
				applications = data.programs
				if(applications.length == 0) {
					await fetch("./command?add")
					await refresh()
				}
				resolve()
			})
			.catch(() => {
				reject()
			})
	})
}

async function loadApp(index) {
	await refresh()
	if(index < 0)
		index = applications.length
	appIndex = index

	if(index == 0)
		appBackEl.innerHTML = "add"
	else 
		appBackEl.innerHTML = "back"

	if(index >= applications.length - 1)
		appNextEl.innerHTML = "add"
	else
		appNextEl.innerHTML = "next"

	if(index >= applications.length) {
		appNameEl.innerHTML = ""
		usedTimeEl.innerHTML = ""
		limitTimeEl.innerHTML = ""
		goalTimeEl.innerHTML = ""
		goalLeftTimeEl.innerHTML = ""
		processNameEl.innerHTML = ""
		appLogoEl.setAttribute("src", "./image.svg")
	}
	else {
		appNameEl.innerHTML = applications[index].friendly_name
		usedTimeEl.innerHTML =
			Math.floor(applications[index].used / config.unit)
		limitTimeEl.innerHTML = 
			Math.floor(applications[index].limit / config.unit)
		goalTimeEl.innerHTML = 
			Math.floor(applications[index].goal / config.unit)
		goalLeftTimeEl.innerHTML = applications[index].goal_cycles
		processNameEl.innerHTML = applications[index].name
		appLogoEl.setAttribute("src", applications[index].image_url)
	}
	if(applications[index].goal_cycles == 0
	|| applications[index].goal == applications[index].limit)
		setLowerLimit(false)
	else
		setLowerLimit(true)
	setScreen("app")
}

appsLinkEl.addEventListener("click", (e) => {
	e.preventDefault()
	loadApp(0);
})

appNextEl.addEventListener("click", async (e) => {
	e.preventDefault()
	if(appNextEl.innerHTML == "add") {
		await fetch("./command?add+process.exe+36000+36000")
	}
	loadApp(appIndex + 1);
})

appBackEl.addEventListener("click", (e) => {
	e.preventDefault()
	loadApp(appIndex - 1);
})

limitTimeEl.addEventListener("focusout", async () => {
	let processName = applications[appIndex].name
	let newLimit = Number(limitTimeEl.innerHTML)
	if(Number.isNaN(newLimit))
		return
	newLimit *= config.unit
	await fetch(`./command?change+${processName}+secs_limit+${newLimit}`)
	refresh()
})

goalTimeEl.addEventListener("focusout", async () => {
	let processName = applications[appIndex].name
	let newGoal = Number(goalTimeEl.innerHTML)
	if(Number.isNaN(newGoal))
		return
	newGoal *= config.unit
	await fetch(`./command?change+${processName}+limit_goal+${newGoal}`)
	refresh()
})

goalLeftTimeEl.addEventListener("focusout", async () => {
	let processName = applications[appIndex].name
	let newGoalLeft = Number(goalLeftTimeEl.innerHTML)
	if(Number.isNaN(newGoalLeft))
		return
	await fetch(`./command?change+${processName}+goal_cycles+${newGoalLeft}`)
	refresh()
})

appNameEl.addEventListener("focusout", async () => {
	let processName = applications[appIndex].name
	let newFriendlyName = appNameEl.innerHTML.replace(/ /g, "_")
	await fetch(
		`./command?change+${processName}+friendly_name+${newFriendlyName}`)
	refresh()
})

processNameEl.addEventListener("focusout", async () => {
	let processName = applications[appIndex].name
	let newProcessName = processNameEl.innerHTML.replace(/ /g, "_")
	if(newProcessName == "")
		return
	await fetch(
		`./command?change+${processName}+name+${newProcessName}`)
	refresh()
})

appLogoEl.addEventListener("click", async (e) => {
	e.preventDefault()
	let processName = applications[appIndex].name
	let newUrl = prompt("Please enter an image url", "")
	try {
		new URL(newUrl)
		appLogoEl.setAttribute("src", newUrl)
	}
	catch {
		newUrl = "./image.svg"
		appLogoEl.setAttribute("src", newUrl)
	}
	await fetch(`./command?change+${processName}+image_url+${newUrl}`)
})

lowerLimitEl.addEventListener("click", async (e) => {
	e.preventDefault()
	let processName = applications[appIndex].name
	let newGoal = Number(limitTimeEl.innerHTML)
	if(Number.isNaN(newGoal))
		return
	newGoal = Math.floor(newGoal * config.unit / 2)
	fetch(`./command?change+${processName}+limit_goal+${newGoal}`)
	fetch(`./command?change+${processName}+goal_cycles+8`)
	goalTimeEl.innerHTML = Math.floor(newGoal / config.unit)
	goalLeftTimeEl.innerHTML = 8
	refresh()
	setLowerLimit(true)
})

cancelLowerLimitEl.addEventListener("click", async (e) => {
	e.preventDefault()
	let processName = applications[appIndex].name
	let newGoal = Number(limitTimeEl.innerHTML)
	if(Number.isNaN(newGoal))
		return
	newGoal *= config.unit
	fetch(`./command?change+${processName}+limit_goal+${newGoal}`)
	fetch(`./command?change+${processName}+goal_cycles+0`)
	setLowerLimit(false)
})

appDeleteEl.addEventListener("click", async (e) => {
	e.preventDefault()
	let processName = applications[appIndex].name
	await fetch(`./command?remove+${processName}`)
	await refresh()
	if(appIndex == applications.length)
		return loadApp(appIndex - 1)
	loadApp(appIndex)
})
