# Project: Message Streaming

## Getting Started
For running app use ./start.sh script

## api
GOlang server, that accepts message and sends to other users
- GET /messages - connection to HTTP Streaming
- POST /messages - send a message with json body { "text": "Hello" }

## bot
A NodeJS server that generates and sends messages to api server
- POST /start - start proccess of message sending
- POST /stop - stop proccess of message sending
- POST /set-ms-speed/:speed - set interval in milliseconds between each request, but mininum 1000
- GET /status - bot status

### Example of connection to streaming
```console
curl -N http://localhost:8081/messages
```
