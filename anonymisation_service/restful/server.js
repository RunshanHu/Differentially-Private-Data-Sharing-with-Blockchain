var express = require('express'),
    app = express(),
    port = process.env.PORT || 3000,
    mongoose = require('mongoose'),
    Task = require('./api/models/anonyModel.js'),
    bodyParser = require('body-parser');

mongoose.Promise = global.Promise;
mongoose.connect('mongodb://localhost/SalaryDb');

app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());

var routes = require('./api/routers/anonyRoutes.js');
routes(app);

app.use(function(req, res){
  res.status(404).send({url: req.originalUrl + ' not found'});
});

app.listen(port);
console.log('anonymisation service RESTful API server start on: ' + port);
