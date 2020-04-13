function onLoaded() {
    var source = new EventSource("//localhost:8081/sse");
    source.onmessage = function (event) {
        console.log("OnMessage called:");
        console.dir(event);
    };
    source.addEventListener("ping", function(event) {
        console.log("PING");
        console.dir(event);
      });
    
}