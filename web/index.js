window.onload = () => {
    const ws = new WebSocket("ws://192.168.86.203:8080/ws/controller");
    ws.onopen = () => { onWsReady(ws); };
    if(/iP(hone|ad)/.test(window.navigator.userAgent)) {
        document.body.addEventListener('touchstart', function() {}, false);
    }
};

function onWsReady(ws) {
    const button = document.getElementById("up-button");
    // button.onmousedown = () => { ws.send("upKeyDown"); };
    // button.onmouseup = () => { ws.send("upKeyUp"); };
    button.addEventListener("touchstart", (e) => { e.preventDefault(); ws.send("upKeyDown"); }, false);
    button.addEventListener("touchend", (e) => { e.preventDefault(); ws.send("upKeyUp"); }, false);
}
