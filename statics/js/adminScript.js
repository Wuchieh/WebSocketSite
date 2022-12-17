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
    if (adminButtonClick > 0) {
        return
    }
    adminButtonClick++
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
    if (adminButtonClick > 0) {
        return
    }
    adminButtonClick++
    let xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/adminRemoveTokenBtn")
    xhr.onload = function () {
        const obj = JSON.parse(xhr.responseText);
        if (obj["status"] === "true") {
            let a = document.getElementById("tokenList-" + id)
            a.remove()

            let b = document.getElementById(id)
            b.remove()

            adminScriptGetTokenInfoClick = 0

            let ebg = document.getElementById("tokenEditBtnGroup")
            ebg.classList.remove("d-block")
            ebg.classList.add("d-none")
        } else {
            alert(obj["msg"])
        }
    }
    xhr.send(id)
}

function adminEditExpiredTime(id) {
    if (adminButtonClick > 0) {
        return
    }
    adminButtonClick++
    let t = document.getElementById("editExpiredTimeInput");
    let xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/adminEditExpiredTime")
    xhr.onload = function () {
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

function setSetting() {
    // const setting = [
    //     "ServerIP",
    //     "Port",
    //     "ScheduleTime",
    //     "ExpiredTime",
    //     "Mode",
    //     "AdminPWD",
    // ]
    // let settingInput = [];
    // for (const string of setting) {
    //     let a = document.getElementById(string)
    //     settingInput.push(a)
    // }
    // const obj = JSON;
    // for (let i = 0; i < setting.length; i++) {
    //     obj[setting[i]] = settingInput[i].value
    // }
    function showToast(title, content, color) {
        const toastDiv = document.getElementById('ToastDiv')
        toastDiv.querySelector("#toast-header").classList.remove("bg-primary", "bg-danger")

        toastDiv.querySelector("#toast-header").classList.add(color)
        toastDiv.querySelector("div > strong").textContent = title
        toastDiv.querySelector(".toast-body").textContent = content

        const toast = new bootstrap.Toast(toastDiv)
        toast.show()
    }

    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/adminSetSetting")
    xhr.onload = function () {
        const obj = JSON.parse(xhr.responseText)
        if (obj["status"] === "true") {
            showToast("Success", obj["msg"], "bg-primary")
        } else {
            showToast("Fail", obj["msg"], "bg-danger")
        }
    }

    xhr.send(JSON.stringify(
        {
            "ServerIP": document.getElementById("ServerIP").value,
            "Port": document.getElementById("Port").value,
            "ScheduleTime": Number(document.getElementById("ScheduleTime").value),
            "ExpiredTime": Number(document.getElementById("ExpiredTime").value),
            "Mode": Number(document.getElementById("Mode").value),
            "AdminPWD": document.getElementById("AdminPWD").value,
        }
    ))
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function reboot() {
    if (adminButtonClick > 0) {
        return
    }
    adminButtonClick++
    const xhr = new XMLHttpRequest()
    let modalBody = document.getElementById("modal-body")
    let progress = document.getElementById("progress")
    modalBody.classList.add("d-none")
    progress.classList.remove("d-none")


    xhr.open("POST", "/api/adminReboot")
    xhr.send()
    for (let i = 0; i < 100; i += 10) {
        if (progress.querySelector("div").style.width === "100%") {
            break
        }
        progress.querySelector("div").style.width = i + "%"
        getServerStatus().then(r => {
        })
        await sleep(1000)
    }
}

async function getServerStatus() {
    const xhr = new XMLHttpRequest()
    xhr.open("GET", location.href)
    xhr.onload = async function () {
        document.getElementById("progress").querySelector("div").style.width = 100 + "%"
        await sleep(1000)
        location.replace(location.href)
    }
    xhr.send()
}