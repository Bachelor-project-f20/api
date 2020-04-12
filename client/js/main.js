function onLoaded() {
    var source = new EventSource("//localhost:8081/api/getUser");
    source.onmessage = function (event) {
        console.log("OnMessage called:");
        console.dir(event);
    };
    
}