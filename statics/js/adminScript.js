let adminScriptGetTokenInfoClick = 0;
let adminButtonClick = 0;

function getTokenInfo(name) {
    let ebg = document.getElementById("tokenEditBtnGroup")
    if (adminScriptGetTokenInfoClick === 0) {
        adminScriptGetTokenInfoClick++
        ebg.classList.remove("d-none")

    }
    let userTokensInfoDiv = document.getElementById("userTokensInfo")
    let hideDiv = document.getElementById(name);
    let userInfos = userTokensInfoDiv.querySelectorAll(".userInfo");
    for (const userInfo of userInfos) {
        userInfo.classList.remove("d-block")
        userInfo.classList.add("d-none")
    }
    hideDiv.classList.remove("d-none")
    hideDiv.classList.add("d-block")
    let updateTokenBtn = document.getElementById("updateTokenBtn");
    let removeTokenBtn = document.getElementById("removeTokenBtn");
    let editExpiredTimeBtn = document.getElementById("editExpiredTimeBtn");

    // 修改按鈕的 onclick 屬性
    updateTokenBtn.setAttribute("onclick", `adminUpdateToken('${name}')`);
    removeTokenBtn.setAttribute("onclick", `adminRemoveToken('${name}')`);
    editExpiredTimeBtn.setAttribute("onclick", `adminEditExpiredTime('${name}')`);
}

function adminUpdateToken(id) {
    // if (adminButtonClick > 0) {
    //     return
    // }
    // adminButtonClick++
    let xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/adminUpdateTokenBtn")
    xhr.onload = function () {
        const obj = JSON.parse(xhr.responseText);
        if (obj["status"] === "true") {
            document.querySelector("#" + id + " .CreateTime").text = obj["CreateTime"]
            document.querySelector("#" + id + " .UpdateTime").text = obj["UpdateTime"]
            document.querySelector("#" + id + " .ExpiredTime").text = obj["ExpiredTime"]
            document.querySelector("#" + id + " .Token").text = obj["Token"]
        } else {
            alert("更新失敗")
        }
    }
    xhr.send(id)
}

function adminRemoveToken(id) {
    // if (adminButtonClick > 0) {
    //     return
    // }
    // adminButtonClick++
}

function adminEditExpiredTime(id) {
    if (adminButtonClick > 0) {
        return
    }
    adminButtonClick++
    let t = document.getElementById("editExpiredTimeInput");
    let xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/adminEditExpiredTime")
    xhr.onload = function (ev) {
        const obj = JSON.parse(xhr.responseText);
        if (obj["status"] === "true") {
            let query = "#" + id + " .ExpiredTime"
            console.log(query)
            let a = document.querySelector(query)
            console.log(a)
            a.text = obj["msg"]
        } else {
            alert(obj["msg"])
        }
        adminButtonClick = 0
    }
    const obj = JSON.stringify({"time": t.valueAsNumber.toString(), "id": id})
    xhr.send(obj)
}
