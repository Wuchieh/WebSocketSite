let clicks = 0;

function getNewToken() {
    if (clicks === 0) {
        const request = new XMLHttpRequest();
        request.open("POST", "/api/GetNewToken");
        request.onload = function () {
            const obj = JSON.parse(request.responseText);
            // console.log(request.responseText);
            if (obj["status"] === true) {
                location.reload();
            }else{
                alert(obj["msg"]);
            }
            clicks = 0;
        }
        request.send();
    }
}