let clicks = 0;
let testClickTime = 0;

function getNewToken() {
    if (clicks === 0) {
        clicks++;
        const request = new XMLHttpRequest();
        request.open("POST", "/api/GetNewToken");
        request.onload = function () {
            const obj = JSON.parse(request.responseText);
            // console.log(request.responseText);
            if (obj["status"] === true) {
                location.reload();
            } else {
                alert(obj["msg"]);
            }
            clicks = 0;
        }
        request.send();
    }
}

function showTestButton() {
    testClickTime++;
    if (testClickTime >= 12) {
        let testButton = document.getElementById("testButtons");
        testButton.style.display = "";
    }
}

function testWsConnect() {
    let url = document.getElementById("testURL");
    let ws = new WebSocket(url.value);
    let EnterIP = document.getElementById("testEnterIP");
    let SendMsg = document.getElementById("testSendMsg");
    let ConnectStatus = document.getElementById("testConnectStatus");
    let WsConnectClose = document.getElementById("testWsConnectClose");
    let WsSendMsg = document.getElementById("testWsSendMsg");
    let MsgInput = document.getElementById("testMsgInput");
    let ReplyMsg = document.getElementById("testReplyMsg");

    ws.onopen = function (e) {
        EnterIP.style.display = "none";
        SendMsg.style.display = "block";
        ConnectStatus.textContent = "已連線";
        ConnectStatus.className = "text-success mt-1";
        ReplyMsg.textContent="";
    }

    ws.onclose = function (e) {
        EnterIP.style.display = "block";
        SendMsg.style.display = "none";
        ConnectStatus.textContent = "尚未連線";
        ConnectStatus.className = "text-danger mt-1";
        removeAllEventListener();
    }

    ws.onmessage = function (e) {
        let message = e.data;
        ReplyMsg.insertAdjacentHTML("beforeend","<div>" + message + "</div>")
        ReplyMsg.scrollTo(0, ReplyMsg.scrollHeight);
    }

    WsConnectClose.addEventListener("click", wsClose);
    WsSendMsg.addEventListener("click", wsSendMsg);
    MsgInput.addEventListener("keypress", keypress);

    function wsSendMsg() {
        ws.send(MsgInput.value);
        MsgInput.value = "";
    }

    function wsClose(){
        ws.close();
    }

    function keypress(e) {
        if (e.key === "Enter") {
            wsSendMsg()
        }
    }

    function removeAllEventListener(){
        WsConnectClose.removeEventListener("click", wsClose);
        WsSendMsg.removeEventListener("click", wsSendMsg);
        MsgInput.removeEventListener("keypress", keypress);
    }
}