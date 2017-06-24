'use strict'

var mongoose = require('mongoose'),
  Data = mongoose.model('Data');

const UTILITY_BOUND = 1500;

function Lap(mu, b, min, max) {
  var x = Math.random() * (max - min) + min;
  var f = 0.5 / b * Math.exp( -Math.abs(x-mu) / b);
  return f;
}

function generateNoise(budget, global_sensitivity, utility_bound) {
  var min = -10;
  var max = 10;
  var noise = Lap(0, budget*global_sensitivity, min, max);
  while (noise > utility_bound) {
      noise = Lap(0, budget*global_sensitivity, min, max);
  }
  console.log("---> generated noise: ", noise );
  return noise;
}

function calculate_sum(records) {
  var sum = 0.0;
  for (var index = 0, length = records.length; index < length; ++index) {
    sum = sum +  records[index].salary; 
  }
  return sum;
}

function calculate_avg(records) {
  return calculate_sum(records) / records.length;
}

function calculate_max(records) {
    var max = -10000.0; 
    for (var index = 0, length = records.length; index < length; ++index) {
      if( records[index].salary > max)
        max  = records[index].salary;
    }
    return max;
}

function calculate_min(records) {
    var min = 1000000.0; 
    for (var index = 0, length = records.length; index < length; ++index) {
      if( records[index].salary < min)
        min  = records[index].salary;
    }
    return min;
}

exports.getsum = function(req, res) {
  console.log('---> get real sum ');
  Data.find({}, function(err, records){
    var sum = calculate_sum(records); 
    res.json({"Result": sum});
  });
}

exports.get_noise_sum = function(req, res) {
  console.log('---> get noise sum')
  Data.find({}, function(err, records){
    var global_sensitivity = calculate_max(records);
    var sum = calculate_sum(records) + generateNoise(req.body.budget, global_sensitivity, UTILITY_BOUND); 
    res.json({"Result": sum});
  });
  
}

exports.getavg = function(req, res) {
  console.log('---> get real avg ');
  Data.find({}, function(err, records){
    var avg = calculate_avg(records); 
    res.json({"Result": avg});
  });
}

exports.get_noise_avg = function(req, res) {
  console.log('---> get noise avg ')
  Data.find({}, function(err, records){
    var global_sensitivity = calculate_max(records) / (records.length - 1);
    var avg = calculate_avg(records) + generateNoise(req.body.budget, global_sensitivity, UTILITY_BOUND);
    res.json({"Result": avg});
  });
}


exports.getmax = function(req, res) {
  console.log('--> get real max ');
  Data.find({}, function(err, records){
    var max = calculate_max(records);
    res.json({"Result": max});
  });
}
exports.get_noise_max = function(req, res) {
  console.log('---> get noise max ');
  Data.find({}, function(err, records){
    var global_sensitivity = calculate_max(records);
    var max = calculate_max(records) + generateNoise(req.body.budget, global_sensitivity, UTILITY_BOUND);
    res.json({"Result": max});
  });
}

exports.getmin = function(req, res) {
  console.log('---> get real min ');
  Data.find({}, function(err, records){
    var min = calculate_min(records);
    res.json({"Result": min});
  });
}

exports.get_noise_min = function(req, res) {
  console.log('---> get noise min ')
  Data.find({}, function(err, records){
    var global_sensitivity = calculate_min(records);
    var min = calculate_min(records) + generateNoise(req.body.budget, global_sensitivity, UTILITY_BOUND);
    res.json({"Result": min});
  });
}

exports.list_all_records = function(req, res) {
  console.log('-->list all records');
  Data.find({}, function(err, records){
    if(err)
      res.send(err);
    res.json(records);
  });
};

exports.create_a_record = function(req, res) {
  console.log('-->create a record');
  console.log('---->request body: ' + req.body);
  var new_record = new Data(req.body);
  new_record.save(function(err, record) {
    if (err)
      res.send(err);
    res.json(record);
  });
};

exports.read_a_record = function(req, res) {
  console.log('-->read a record');
  Data.findById(req.params.Id, function(err, record){
    if(err)
      res.send(err);
    res.json(record);
  });
};

exports.update_a_record = function(req, res) {
  console.log('-->update a record');
  Data.findOneAndUpdate(req.params.Id, req.body, {new: true}, function(err, record){
    if (err) 
      res.send(err);
    res.json(record);
  });
};

exports.delete_a_record = function(req, res) {
  console.log('-->delete a record');
  Data.remove({
    _id: req.params.Id
  }, function(err, record) {
    if(err) 
      res.send(err);
    res.json({ message: 'Record successfully deleted' });
  });
};
