angular.module('coffee', []).controller('coffeeController', function () {

  var coffee = this;

  coffee.coffees = [];

  coffee.schedules = [];

  // Add a new coffee instance
  coffee.addCoffee = function () {
    // Push into coffees array
    coffee.coffees.push({
      name: '',
      type: '',
      sugar: '',
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

  coffee.addSchedule = function () {
    // Push into coffees array
    coffee.schedules.push({
      coffee: '',
      sunday: false,
      monday: false,
      tuesday: false,
      wednesday: false,
      thursday: false,
      friday: false,
      saturday: false
    });
  };

  coffee.saveSchedule = function () {
    // TODO - Write this function
  };

  coffee.removeSchedule = function () {

  };

});
