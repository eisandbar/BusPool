const {Kafka, CompressionTypes, CompressionCodecs} = require('kafkajs')
const SnappyCodec = require('kafkajs-snappy')

const express = require('express');
const app = express();
const http = require('http');
const server = http.createServer(app);
const { Server } = require("socket.io");
const io = new Server(server, {
    cors: "localhost:5173"
});

app.get('/', (req, res) => {
    res.sendFile(__dirname + '/index.html');
});

io.on('connection', (socket) => {
    console.log('a user connected');
});

server.listen(3000, () => {
    console.log('listening on *:3000');
});

io.on("connection", socket => {
    console.log("connected")
})

CompressionCodecs[CompressionTypes.Snappy] = SnappyCodec

const kafka = new Kafka({
    brokers: ["localhost:9091"]
})

const consumer = kafka.consumer({
    groupId: "Web"
  })

let buses = []

setInterval(() => {
    console.log("Updating locations", buses[0])
    io.emit("update location", buses)
}, 2000)

const main = async () => {
  await consumer.connect()

  await consumer.subscribe({
    topic: "bus-positions-elastic",
    fromBeginning: false
  })

  await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {

        let msg = message ? JSON.parse(message.value) : {}
        let Id = msg.id
        let Location = msg.location
        if (buses[Id]) {
            buses[Id].Lat = Location.lat
            buses[Id].Lon = Location.lon
        } else {
            buses[Id] = {
                Lat: Location.lat,
                Lon: Location.lon
            }
        }
        }
  })
}

main().catch(async error => {
  console.error(error)
  try {
    await consumer.disconnect()
  } catch (e) {
    console.error('Failed to gracefully disconnect consumer', e)
  }
  process.exit(1)
})
