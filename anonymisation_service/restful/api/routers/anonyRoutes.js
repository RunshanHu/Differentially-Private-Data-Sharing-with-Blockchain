'use strict'
module.exports = function(app) {
  var data_controller = require('../controllers/anonyController.js');

  app.route('/dataset/all')
    .get(data_controller.list_all_records)
   
  app.route('/dataset/create')
    .post(data_controller.create_a_record);

  app.route('/dataset/sum')
    .post(data_controller.get_noise_sum)
    .get(data_controller.getsum);

  app.route('/dataset/avg')
    .post(data_controller.get_noise_avg)
    .get(data_controller.getavg);

  app.route('/dataset/max')
    .post(data_controller.get_noise_max)
    .get(data_controller.getmax)
  
  app.route('/dataset/min')
    .post(data_controller.get_noise_min)
    .get(data_controller.getmin)

  app.route('/dataset/:Id')
    .get(data_controller.read_a_record)
    .put(data_controller.update_a_record)
    .delete(data_controller.delete_a_record);

};
