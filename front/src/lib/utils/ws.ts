import WebSocket from "ws";

export function connectWs() {
  const ws = new WebSocket("ws://localhost:8088/ws");

  ws.on("error", console.error);

  ws.on("open", function open() {
    ws.send("something");
  });

  ws.on("message", (data: any) => {
    console.log("received: %s", data);
  });
}
