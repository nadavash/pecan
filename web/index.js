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
        const connectionView = document.getElementById("connection-view");
        connectionView.style.display = "none";
        const loadingView = document.getElementById("loading-view");
        loadingView.style.display = "none";
        const controllerView = document.getElementById("controller-view");
        controllerView.style.display = "none";

        switch (activeView) {
            case "connection-view":
                connectionView.style.display = "block";
                break;
            case "loading-view":
                loadingView.style.display = "block";
                break;
            case "controller-view":
                controllerView.style.display = "block";
                break;
        }

    }
};
