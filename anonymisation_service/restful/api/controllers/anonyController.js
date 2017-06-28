'use strict'

var mongoose = require('mongoose'),
  Data = mongoose.model('Data');

const UTILITY_BOUND = 1500;
const SMALL_BUDGET = 0.05;
const SENSI = 1;

function Lap(mu, b) {
  var x = Math.random() - 0.5;
  var f;
  if (x > 0) {
      f = mu - b * Math.log(1 - 2*x);
  } else {
      f = mu + b * Math.log(1 + 2*x);
  }
  return f;
}

function generateNoise(budget, global_sensitivity, utility_bound) {
  var noise = Lap(0, global_sensitivity / budget);
  while (Math.abs(noise) > utility_bound) {
      noise = Lap(0, global_sensitivity / budget);
  }
  console.log("---> generating noise, budget: ", budget, " global_sensitivity: ", global_sensitivity);
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

exports.getavg = function(req, res) {
  console.log('---> get real avg ');
  Data.find({}, function(err, records){
    var avg = calculate_avg(records); 
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

exports.getmin = function(req, res) {
  console.log('---> get real min ');
  Data.find({}, function(err, records){
    var min = calculate_min(records);
    res.json({"Result": min});
  });
}

exports.get_noise_sum = function(req, res) {
  console.log('---> get noise sum')
  Data.find({}, function(err, records){
    var global_sensitivity, budget;
    console.log("---> received request body: ", req.body);
    if(req.body.flag == 1) {
        global_sensitivity = SENSI;
        budget = SMALL_BUDGET;
    } else {
        global_sensitivity = calculate_max(records);
        budget = req.body.budget;
    }
    var sum = calculate_sum(records) + generateNoise(budget, global_sensitivity, UTILITY_BOUND); 
    res.json({"Result": sum});
  });
}

exports.get_noise_avg = function(req, res) {
  console.log('---> get noise avg ')
  Data.find({}, function(err, records){
    var global_sensitivity, budget;
    console.log("---> received request body: ", req.body);
    if(req.body.flag == 1) {
        global_sensitivity = SENSI;
        budget = SMALL_BUDGET;
    } else {
        global_sensitivity = calculate_max(records) / records.length;
        budget = req.body.budget;
    }
    var avg = calculate_avg(records) + generateNoise(budget, global_sensitivity, UTILITY_BOUND);
    res.json({"Result": avg});
  });
}

exports.get_noise_max = function(req, res) {
  console.log('---> get noise max ');
  Data.find({}, function(err, records){
    var global_sensitivity, budget;
    console.log("---> received request body: ", req.body);
    if(req.body.flag == 1) {
        global_sensitivity = SENSI;
        budget = SMALL_BUDGET;
    } else {
        global_sensitivity = calculate_max(records);
        budget = req.body.budget;
    }
    var max = calculate_max(records) + generateNoise(budget, global_sensitivity, UTILITY_BOUND);
    res.json({"Result": max});
  });
}

exports.get_noise_min = function(req, res) {
  console.log('---> get noise min ')
  Data.find({}, function(err, records){
    var global_sensitivity, budget;
    console.log("---> received request body: ", req.body);
    if(req.body.flag == 1) {
        global_sensitivity = SENSI;
        budget = SMALL_BUDGET;
    } else {
        global_sensitivity = calculate_min(records);
        budget = req.body.budget;
    }
    var min = calculate_min(records) + generateNoise(budget, global_sensitivity, UTILITY_BOUND);
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
