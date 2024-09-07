import React, { useEffect, useState } from "react";

function App() {

    const [number, setNumber] = useState<number | null>(null);

    useEffect(() => {
        const socket = new WebSocket("ws://localhost:8080/ws");

        socket.onmessage = (event: MessageEvent) => {
            const receivedNumber = parseFloat(event.data);
            setNumber(receivedNumber);
        };

        socket.onclose = () => {
            console.log("WebSocket connection closed");
        };

        return () => socket.close();

    }, []);


    return (
        <div>
            <h1>Random Number: {number}</h1>
        </div>

    );
}

export default App;

