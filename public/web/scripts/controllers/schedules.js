'use strict';

/**
 * @ngdoc function
 * @name redqueenUiApp.controller:SchedulesCtrl
 * @description
 * # SchedulesCtrl
 * Controller of the redqueenUiApp
 */
angular.module('redqueenUiApp')
  .controller('SchedulesCtrl', [ '$scope', '$location', 'Schedule', function ($scope, $location, ScheduleResource) {
    $scope.schedules = [];
    $scope.activeMenu = 'schedules';

    $scope.perPage = 30;
    $scope.page = 1;
    $scope.totalItems = 0;

    $scope.edit = function SchedulesCtrlEdit(rfidCard) {
      $location.path('/schedules/' + rfidCard.id + '/edit');
    };

    var update = function() { 
        ScheduleResource.all($scope.page, $scope.perPage).then(function(data) {
          $scope.schedules = data.items;
          $scope.totalItems = data.total_items;
        });
    };

    $scope.querySchedules = function() { 
        update();
    };

    update();

}]);
