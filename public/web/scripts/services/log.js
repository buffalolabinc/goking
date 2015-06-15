'use strict';

/**
 * @ngdoc service
 * @name redqueenUiApp.log
 * @description
 * # log
 * Service in the redqueenUiApp.
 */
angular.module('redqueenUiApp')
  .service('Log', [ '$q', '$timeout', '$http', 'underscore', function($q, $timeout, $http, _) {

    function Log(data) {
      angular.extend(this, data);
    }

    Log.all = function LogResourceAll(page, per_page) {
      var deferred = $q.defer();

      $http({ 
          url: '/api/logs', 
          method: 'GET',
          params: { 
              'page': page,
              'per_page' : per_page
          }
      }).then(function(data) {
        deferred.resolve(data.data);
      }, function() {
        deferred.reject();
      });

      return deferred.promise;
    };

    return Log;
  }]);
