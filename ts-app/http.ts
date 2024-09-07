
import { serve } from "bun";

const htmlContent = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
  </head>
  <body>

    <img id="dynamicImage" src="" alt="Image" width="500" height="500">

    <script>
        let socket = new WebSocket("ws://localhost:3000/ws");
        console.log("Attempting Connection...");

        socket.onopen = () => {
            console.log("Successfully Connected");
            socket.send("Hi From the Client!");
        };

        socket.onmessage = (msg) => {
            console.log("Message from Server: ", msg.data);
            // imgタグのsrc属性を更新
            const img = document.getElementById("dynamicImage");
            img.src = msg.data;
        }

        socket.onclose = (event) => {
            console.log("Socket Closed Connection: ", event);
        };

        socket.onerror = (error) => {
            console.log("Socket Error: ", error);
        };

        function sendMessage() {
            const input = document.getElementById("chatInput");
            const message = input.value;
            if(message) {
                socket.send(message);
                input.value = "";
            }
        }
    </script>
  </body>
</html>
`;

// <input type="text" id="chatInput" />
// <button onclick="sendMessage()">Send</button>



serve({
  port: 3030,

  fetch(request) {

    return new Response(htmlContent, {
      headers: {
        "Content-Type": "text/html",

      },
    });

  },
});


console.log("Server is running on http://localhost:3030");
