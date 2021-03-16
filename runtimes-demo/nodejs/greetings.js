const express = require('express');
const http = require('http');
const { collectDefaultMetrics, register } = require('prom-client');

collectDefaultMetrics();
const app = express();

app.use('/greetings', (request, response) => {
  const from = request.query.from;
  if (!from) {
    response.send("Hello, I'm NodeJS!");
    ['vertx', 'quarkus', 'springboot'].forEach(peer => {
      console.log('client call to ' + peer);
      http.get(`http://${peer}:8080/greetings?from=nodejs`, {}, res => {
        let str = ''
        res.on('data', chunk => str += chunk);
        res.on('end', () => {
          console.log(str);
        });
      }).on('error', e => {
        console.log("Got error: " + e.message);
      });
    });
  } else {
    const resp = (`Hello ${from}, I'm NodeJS!`);
    console.log(resp);
    response.send(resp);
  }
});

// expose our metrics at the default URL for Prometheus
app.get('/metrics', (_, response) => register.metrics().then(metrics => {
  response.set('Content-Type', register.contentType);
  response.send(metrics);
}));

app.listen(8080, () => console.log(`Greetings app listening on port 8080!`));
