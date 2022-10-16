const WebSocket = require("ws");
const ws = new WebSocket("ws://localhost:3000/ws/" + process.argv[2]);

ws.onmessage = (event) => {
  try {
    const data = JSON.parse(event.data);
    console.log(data.from, ":", data.data);
  } catch (e) {
    console.log(event.data);
  }
};

const readline = require("readline").createInterface({
  input: process.stdin,
  output: process.stdout,
});

const recursiveReadline = () => {
  readline.question("Message: ", (message) => {
    ws.send(
      JSON.stringify({
        from: process.argv[2],
        data: message,
        to: process.argv[3],
      })
    );
    recursiveReadline();
  });
};

recursiveReadline();
