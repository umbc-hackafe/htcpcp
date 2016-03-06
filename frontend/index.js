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

    /**
    * NOTE - Create Schedule
    */

    // New index for coffee instance is last index
    var index = coffee.coffees.length - 1;

    // Post request to create new coffee object
    $http.post('/api/create/coffee', {
      // Empty name
      name: '',
      // Empty days
      type: '',
      // Size
      size: '',
      // Kcup
      k_cup: '',
      // Empty time
      sugar: '',
      // 0 for drink id
      creamer: '',
    }).then(function (response) { // Success callback
      // Set the coffee item id
      coffee.schedules[index].id = response.data.id;
    }, function (response) { // Failure callback
      // TODO - YOU FAILED, time to debug
      console.log('FAILURE')
      console.log(response);
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

    // Push into coffees array
    coffee.schedules.push({
      // Schedule name
      name: '',
      // Coffee instance
      coffee: '',
      // Time string
      time: '',
      // Days of the week check boxes
      sunday: false,
      monday: false,
      tuesday: false,
      wednesday: false,
      thursday: false,
      friday: false,
      saturday: false
    });

    /**
    * NOTE - Create Schedule
    */

    // New index for schedule instance is last index
    var index = coffee.schedules.length - 1;

    // Post request to create new coffee object
    $http.post('/api/create/schedule', {
      // Empty name
      name: '',
      // Empty days
      days: [],
      // Empty time
      time: '',
      // 0 for drink id
      drink: 0,
      // 0 for machine id
      machine: 0,
    }).then(function (response) { // Success callback
      // Set the coffee item id
      coffee.schedules[index].id = response.data.id;
    }, function (response) { // Failure callback
      // TODO - YOU FAILED, time to debug
      console.log('FAILURE')
      console.log(response);
    });

  };

  // Save a schedule instance
  coffee.saveSchedule = function (index) {

    /**
    * NOTE - Update Schedule
    */

    // Initialize days array
    var daysArr = [];
    if (coffee.schedules[index].sunday) daysArr.push('sunday');
    if (coffee.schedules[index].monday) daysArr.push('monday');
    if (coffee.schedules[index].tuesday) daysArr.push('tuesday');
    if (coffee.schedules[index].wednesday) daysArr.push('wednesday');
    if (coffee.schedules[index].thursday) daysArr.push('thursday');
    if (coffee.schedules[index].friday) daysArr.push('friday');
    if (coffee.schedules[index].saturday) daysArr.push('saturday');

    // Post request to create new coffee object
    $http.post('/api/create/schedule', {
      // ID given to specify update and not create
      id: coffee.schedules[index].id,
      // Empty name
      name: coffee.schedules[index].name,
      // Empty days
      days: daysArr,
      // Empty time
      time: Math.floor(coffee.schedules[index].time.getTime() / 1000),
      // 0 for drink id
      drink: 0,
      // 0 for machine id
      machine: 0,
    }).then(function (response) { // Success callback
      // Set the coffee item id
      coffee.schedules[index].id = response.data.id;
    }, function (response) { // Failure callback
      // TODO - YOU FAILED, time to debug
      console.log('FAILURE')
      console.log(response);
    });
  };

  // Remove a schedule instance
  coffee.removeSchedule = function (index) {
    // Remove the schedule instance from the array
    coffee.schedules.splice(index, 1);
  };

});
