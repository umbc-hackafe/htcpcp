// Coffee front end angular module
angular.module('coffee', ['ngMaterial']).controller('coffeeController', function ($scope, $http) {

  // Coffee application
  var coffee = this;

  // Array of coffee instances
  coffee.coffees = [];

  // Array of schedule instances
  coffee.schedules = [];

  // Sequential schedule instance id
  coffee.scheduleId = 0;

  // Add a new coffee instance
  coffee.addCoffee = function () {
    // Push into coffees array
    coffee.coffees.push({
      // Name of coffee
      name: '',
      // Type of coffee
      type: '',
      // Lumps of sugar
      sugar: '',
      // Globs of creamer
      creamer: ''
    });
  };

  // Save a coffee instance
  coffee.saveCoffee = function (index) {
    // TODO - Write this function
  };

  // Remove a coffee instance
  coffee.removeCoffee = function (index) {
    // Remove the coffe instance from the array
    coffee.coffees.splice(index, 1);
  };

  // Add a schedule instance
  coffee.addSchedule = function () {
    // // Post request to create new coffee object
    // $http.post('/api/create/schedule', {
    //   // Data object to send
    //   foo: 'bar'
    // }).then(function () {
    //   // TODO - Success
    //   console.log('Success');
    // }, function () {
    //   // TODO - FAILURE
    //   console.log('Failure');
    // });
    // Push into coffees array
    coffee.schedules.push({
      // Coffee instance
      coffee: '',
      // Sequential schedule instance id
      id: coffee.scheduleId,
      // Days of the week check boxes
      sunday: false,
      monday: false,
      tuesday: false,
      wednesday: false,
      thursday: false,
      friday: false,
      saturday: false
    });
    // Increase the schedule id
    coffee.scheduleId += 1;
  };

  // Save a schedule instance
  coffee.saveSchedule = function (index) {
    // TODO - Write this function
  };

  // Remove a schedule instance
  coffee.removeSchedule = function (index) {
    // Remove the schedule instance from the array
    coffee.schedules.splice(index, 1);
  };

});
