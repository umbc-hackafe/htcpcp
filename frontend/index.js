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

    var index = coffee.schedules.length - 1;

    // Post request to create new coffee object
    $http.post('/api/create/schedule', {
      name: '',
      days: [],
      time: '',
      drink: 0,
      machine: 0,
    }).then(function () {
      // TODO - Success
      console.log('Success');
    }, function () {
      // TODO - FAILURE
      console.log('Failure');
    });

  };

  // Save a schedule instance
  coffee.saveSchedule = function (index) {

    var daysArr = [];
    if (coffee.schedules[index].sunday) daysArr.push('sunday');
    if (coffee.schedules[index].monday) daysArr.push('monday');
    if (coffee.schedules[index].tuesday) daysArr.push('tuesday');
    if (coffee.schedules[index].wednesday) daysArr.push('wednesday');
    if (coffee.schedules[index].thursday) daysArr.push('thursday');
    if (coffee.schedules[index].friday) daysArr.push('friday');
    if (coffee.schedules[index].saturday) daysArr.push('saturday');

    function formatAMPM (date) {
      var hours = date.getHours();
      var minutes = date.getMinutes();
      var ampm = hours >= 12 ? 'pm' : 'am';
      hours = hours % 12;
      hours = hours ? hours : 12; // the hour '0' should be '12'
      minutes = minutes < 10 ? '0'+minutes : minutes;
      var strTime = hours + ':' + minutes + ' ' + ampm;
      return strTime;
    }

    console.log({
      name: coffee.schedules[index].name,
      days: daysArr,
      time: Math.floor(coffee.schedules[index].time.getTime() / 1000)
    });
  };

  // Remove a schedule instance
  coffee.removeSchedule = function (index) {
    // Remove the schedule instance from the array
    coffee.schedules.splice(index, 1);
  };

});
