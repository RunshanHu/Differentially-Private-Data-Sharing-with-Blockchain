'use strict'
var mongoose = require('mongoose');
var Schema = mongoose.Schema;

var DataSchema = new Schema({
  name: {
    type: String,
    Required: 'Kindly enter the name of person'
  },

  salary: {
    type: Number,
    Required: 'Input salary of this person'
  }

});

module.exports = mongoose.model('Data', DataSchema);
