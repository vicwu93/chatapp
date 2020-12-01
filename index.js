// console.log("reached js file")
const wsLOCALHOST = `ws://localhost:8888/ws?user=`
let myApp = {}
myApp.websocket
myApp.container

myApp.append = (msg) => {
    let p = document.createElement('p')
    p.innerHTML = msg 
    myApp.container.append(p)
    var elem = document.getElementById("container")
    elem.scrollTop = elem.scrollHeight
}

/**
 * 1. Sets the user window message as "me: msg"
 * 2. Sends the msg to server
 */
myApp.sendmsg = () => {
    let msg = document.getElementById('user-input').value 
    console.log(msg)
    // maybe better to append after ws send
    myApp.append(`<b>me</b>: ${msg}`)
    console.log("after send msg")
    myApp.websocket.send(JSON.stringify({
        Msg: msg
    }))
    document.getElementById('user-input').value = ''
}

myApp.init = () => {
    let user = prompt(`What's your name?`)
    myApp.websocket = new WebSocket(`${wsLOCALHOST}${user}`)
    console.log(myApp.websocket)
    myApp.container = document.getElementById('container')

    myApp.websocket.onmessage = (event) => {
        let msg = ''

        console.log("received onmsg")
        let res = JSON.parse(event.data)
        console.log("reached here...")

        switch(res.Option) {
            case 'connect':
                msg = `<b>${res.User}</b> has joined the room.`
                break
            case 'disconnect':
                msg = `<b>${res.User}</b> has left the room.`
                break
            default:
                console.log(res.Option);
                msg = `<b>${res.User}</b>: ${res.Message}`
                break
        }
        myApp.append(msg)
    }
}

const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');

function switchTheme(e) {
    if (e.target.checked) {
        document.documentElement.setAttribute('data-theme', 'dark');
        localStorage.setItem('theme', 'dark');
    } else {
        document.documentElement.setAttribute('data-theme', 'light');
        localStorage.setItem('theme', 'light');
    }
}

toggleSwitch.addEventListener('change', switchTheme, false);

const currentTheme = localStorage.getItem('theme') ? localStorage.getItem('theme') : null;

if (currentTheme) {
    document.documentElement.setAttribute('data-theme', currentTheme);

    if (currentTheme === 'dark') {
        toggleSwitch.checked = true;
    }
}

window.onload = myApp.init