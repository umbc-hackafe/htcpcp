angular.module('coffee', []).controller('coffeeController', function () {

  var coffee = this;

  coffee.coffees = [
    {
      name: 'Blueberry Blast',
      type: 'Blueberry Grind',
      sugar: '2',
      creamer: '3'
    },
  ];

  coffee.schedules = [
    {name: 'Morning Coffee'},
  ];

  // Add a new coffe instance
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
  coffee.saveCoffee = function (item) {
    // TODO - Write this function
  };

  // Remove a coffee instance
  coffee.removeCoffee = function (item) {

  };

  coffee.addSchedule = function () {
    coffee.schedules.push({name:'Name of Schedule'});
  };

});
