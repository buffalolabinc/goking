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

    ScheduleResource.all().then(function(data){ 
        $scope.has_schedules = (data.length > 0);

    });

    RfidCardResource.all().then(function(data) {
      $scope.rfidCards = data;
    });
    
    $scope.edit = function RfidCardsCtrlEdit(rfidCard) {
      $location.path('/rfidcards/' + rfidCard.Id + '/edit');
    };

  }]);
