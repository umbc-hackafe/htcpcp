function postCoffee($http, coffee) {
  // Post request to create new coffee object
  $http.post('/api/update/drink', {
    // Empty name
    name: coffee.name,
    // Size
    size: coffee.size,
    // Kcup string for yes, empty string for no.
    k_cup: coffee.k_cup,
    // Empty intial sugar and creamer.
    sugar: tryParse(coffee.sugar),
    creamer: tryParse(coffee.creamer),
    // Teabag
    tea_bag: coffee.tea_bag,
    id: coffee.id
  }).then(function (response) { // Success callback
    // Set the coffee item id
    coffee.id = response.data.id;
  }, function (response) { // Failure callback
    // TODO - YOU FAILED, time to debug
    console.log('FAILURE')
    console.log(response);
  });
}

function tryParse(number) {
  if (angular.isNumber(number)) {
    return number;
  } else if (angular.isString(number)) {
    return Number.parseInt(number, 10);
  } else {
    console.log("Not a number: ", number);
    return 0;
  }
}

function convertDays(schedule) {
  var daysArr = [];
  if (schedule.sunday) daysArr.push('sunday');
  if (schedule.monday) daysArr.push('monday');
  if (schedule.tuesday) daysArr.push('tuesday');
  if (schedule.wednesday) daysArr.push('wednesday');
  if (schedule.thursday) daysArr.push('thursday');
  if (schedule.friday) daysArr.push('friday');
  if (schedule.saturday) daysArr.push('saturday');
  return daysArr;
}

function postSchedule($http, schedule) {
  /**
  * NOTE - Create Schedule
  */

  // Post request to create new coffee object
  $http.post('/api/update/schedule', {
    // Empty name
    name: schedule.name,
    // Empty days
    days: convertDays(schedule),
    // Empty time
    time: Math.floor(
      (schedule.time.getTime() - new Date(0, 1, 1, 0, 0, 0, 0).getTime()) / 1000),
    // Enabled flag
    enabled: schedule.enabled,
    // 0 for drink id
    drink: schedule.drink,
    // 0 for machine id
    machine: schedule.machine,
    id: schedule.id
  }).then(function (response) { // Success callback
    // Set the coffee item id
    schedule.id = response.data.id;
  }, function (response) { // Failure callback
    // TODO - YOU FAILED, time to debug
    console.log('FAILURE')
    console.log(response);
  });
}

// Coffee front end angular module
angular.module('coffee', ['ngMaterial']).controller('coffeeController', function ($scope, $http) {

  // Coffee application
  var coffee = this;

  // Array of coffee instances
  coffee.coffees = [];

  // Array of schedule instances
  coffee.schedules = [];

  // Add a new coffee instance
  coffee.addCoffee = function () {
    var newCoffee = {
      // Name of coffee
      name: '',
      // Lumps of sugar
      sugar: 0,
      // Globs of creamer
      creamer: 0,
      // Initial value of ID
      id: 0,
      size: 6,
      k_cup: 'yes',
      tea_bag: ''
    };

    // Push into coffees array
    coffee.coffees.push(newCoffee);

    postCoffee($http, newCoffee);
  };

  // Save a coffee instance
  coffee.saveCoffee = function (instance) {
    postCoffee($http, instance);
  };

  // Remove a coffee instance
  coffee.removeCoffee = function (index) {
    // Remove the coffe instance from the array
    coffee.coffees.splice(index, 1);
  };

  coffee.brewCoffee = function (instance) {
    $http.post('/api/brew', {
      drink: instance.id,
      machine: 1
    }).then(function (response) { // Success callback
      // Set the coffee item id
      console.log('Brew request succeeded');
    }, function (response) { // Failure callback
      // TODO - YOU FAILED, time to debug
      console.log('FAILURE')
      console.log(response);
    });
  };

  // Add a schedule instance
  coffee.addSchedule = function () {
    var newSchedule = {
      // Schedule name
      name: '',
      // Time string
      time: new Date(0, 1, 1, 0, 0, 0, 0),
      // Days of the week check boxes
      sunday: false,
      monday: false,
      tuesday: false,
      wednesday: false,
      thursday: false,
      friday: false,
      saturday: false,
      // Initial value of ID
      id: 0,

      drink: 0,
      machine: 1
    }

    // Push into coffees array
    coffee.schedules.push(newSchedule);

    postSchedule($http, newSchedule)
  };

  // Save a schedule instance
  coffee.saveSchedule = function (schedule) {
    postSchedule($http, schedule);
  };

  // Remove a schedule instance
  coffee.removeSchedule = function (index) {
    // Remove the schedule instance from the array
    coffee.schedules.splice(index, 1);
  };

});
