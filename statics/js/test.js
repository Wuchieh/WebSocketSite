function getContent(path, callback) {
    const xhr = new XMLHttpRequest();
    const content = document.getElementById("content");
    xhr.open("POST", "/api/getContent");
    xhr.onload = function () {
        content.innerHTML = xhr.responseText;
        if (typeof callback === 'function') {
            callback();
        }
    }
    xhr.onerror = function () {
        content.innerHTML = "<h1>加載錯誤</h1>";
    }
    xhr.send(path);
}

getContent('your-html-file.html', function () {
    eval('(function () { alert("Hello World") })()');
});