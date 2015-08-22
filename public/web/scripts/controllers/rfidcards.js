'use strict';

/**
 * @ngdoc function
 * @name redqueenUiApp.controller:RfidcardsCtrl
 * @description
 * # RfidcardsCtrl
 * Controller of the redqueenUiApp
 */
angular.module('redqueenUiApp')
  .controller('RfidCardsCtrl', [ '$scope', '$location', 'RfidCard','Schedule', function ($scope, $location, RfidCardResource, ScheduleResource) {
    $scope.rfidCards = [];
    $scope.activeMenu = 'cards';

    $scope.perPage = 30;
    $scope.page = 1;
    $scope.totalItems = 0;

    ScheduleResource.all().then(function(data){ 
        $scope.has_schedules = (data.total_items > 0);
    });

    $scope.edit = function RfidCardsCtrlEdit(rfidCard) {
      $location.path('/rfidcards/' + rfidCard.id + '/edit');
    };

    var update = function() {
        RfidCardResource.all($scope.page, $scope.perPage).then(function(data) {
              $scope.rfidCards = data.items;
              $scope.totalItems = data.total_items;
        });
    };

    $scope.queryCards = function() { 
        update();
    };

    update();

  }]);
