const express = require('express');
const axios = require('axios');
const { faker } = require('@faker-js/faker');

const app = express();
const PORT = process.env.APP_PORT || 3000;
const REQUEST_URL = process.env.REQUEST_URL || '';

let interval = null;
let speedInMs = 20000

async function sendRequest() {
    const randomSentence = faker.lorem.sentence();
    try {
        await axios.post(REQUEST_URL, {
            text: randomSentence
        })
    } catch (e) {
        console.error('Cannot be sended', e)
    }
}

app.get('/status', async (req, res) => {
    res.json({ isWorking: !!interval })
});

app.post('/start', async (req, res) => {
  interval = setInterval(sendRequest, speedInMs);
  res.json({ message: 'started' })
});

app.post('/stop', async (req, res) => {
  clearInterval(interval);
  interval = null;

  res.json({ message: 'stopped' })
});

app.post('/set-ms-speed/:speed', async (req, res) => {
    const { speed } = req.params;
    const nSpeed = parseInt(speed);
    if (Number.isNaN(nSpeed) || nSpeed < 1000) {
        res.status(404).json({ error: '1000ms is minimum' });
        return;
    }

    speedInMs = nSpeed;
    if (!!interval) {
        clearInterval(interval);
        interval = setInterval(sendRequest, speedInMs);
    }
   

    res.json({ message: `Message is sended every ${speedInMs}ms` });
});

app.listen(PORT, () => {
  console.log(`Express bot is running on port ${PORT}`);
});
