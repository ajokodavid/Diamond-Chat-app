var alerter = document.getElementById("alerter");

class MySocket{
    constructor() {
        this.mysocket = null;
        this.msgcontainer = document.getElementById("msgcontainer");
        this.msginput = document.getElementById("input1");
    }

    showMessage(text, myself) {
        var div = document.createElement("div")
        div.innerHTML = text;
        var cself = (myself)? "self" : "";
        div.className = "msg " + cself;
        div.classList.add("d-flex", "justify-content-center")
        this.msgcontainer.appendChild(div)
    }

    send() {
        var txt = this.msginput.value;
        this.showMessage("<b>Me:</b> " + ` <b> ${txt} </b>`,true);
        this.mysocket.send(txt);
        this.msginput.value = "";
        let msgArr = []
        msgArr.unshift(txt.split(" "))
        localStorage.setItem("Clients", `${JSON.stringify(msgArr)}`)
        console.log(JSON.stringify(msgArr));
    }

    click() {
        this.send();
    }

    connectSocket() {
        console.log("David's Socket");
        var socket = new WebSocket("ws://:PORT/echo");
        this.mysocket = socket;

        socket.onmessage = (e) => {
            this.showMessage(e.data, false)
        }

        socket.onopen = () => {
            console.log("Socket Opened");
            alerter.innerHTML = "Socket Opened"
        }

        socket.onclose = () => {
            console.log("Socket Closed");
            alerter.innerHTML = "Socket Closed";
        }
    }
}

var mysocket = new MySocket()
mysocket.connectSocket();