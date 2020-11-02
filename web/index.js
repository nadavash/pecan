window.onload = () => {
    const connectButton = document.getElementById("connect");
    connectButton.onclick = () => {
        const ipAddress = document.getElementById("ip-address");
        const wsAddress = `ws://${ipAddress.value}:8080/ws/controller`
        const ws = new WebSocket(wsAddress);
        ws.onopen = onWsReady;
        ws.onclose = () => setViewState("connection-view");
        setViewState("loading-view");
    };

    if(/iP(hone|ad)/.test(window.navigator.userAgent)) {
        document.body.addEventListener('touchstart', function() {}, false);
    }

    function onWsReady() {
        setViewState("controller-view");
        const button = document.getElementById("up-button");
        // button.onmousedown = () => { ws.send("upKeyDown"); };
        // button.onmouseup = () => { ws.send("upKeyUp"); };
        button.addEventListener(
            "touchstart", 
            (e) => { e.preventDefault(); this.send("upKeyDown"); }, 
            false);
        button.addEventListener(
            "touchend", 
            (e) => { e.preventDefault(); this.send("upKeyUp"); }, 
            false);
    }

    function setViewState(activeView) {
        const views = ["connection-view", "loading-view", "controller-view"];
        views.forEach((viewId) => {
            const view = document.getElementById(viewId);
            view.style.display = viewId === activeView ? "block" : "none";
        });
    }
};
