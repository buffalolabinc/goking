'use strict';

/**
 * @ngdoc function
 * @name redqueenUiApp.controller:LogsCtrl
 * @description
 * # LogsCtrl
 * Controller of the redqueenUiApp
 */
angular.module('redqueenUiApp')
  .controller('LogsCtrl', [ '$scope', 'Log', function ($scope, LogResource) {
    $scope.activeMenu = 'logs';

    $scope.perPage = 30;
    $scope.page = 1;
    $scope.totalItems = 0;

    var update = function() { 
        LogResource.all($scope.page, $scope.perPage).then(function(data) {
            $scope.logs = data.items;
            $scope.totalItems = data.total_items;
        });
    };

    $scope.queryLogs = function(page) { 
        update();
    };

    update();

  }]);
