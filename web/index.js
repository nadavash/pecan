window.onload = () => {
    const ws = new WebSocket("ws://192.168.86.203:8080/ws/controller");
    ws.onopen = () => { onWsReady(ws); };
};

function onWsReady(ws) {
    const button = document.getElementById("up-button");
    button.onmousedown = () => {
        ws.send("upKeyDown");
    };

    button.onmouseup = () => {
        ws.send("upKeyUp");
    };
}