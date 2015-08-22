'use strict';

/**
 * @ngdoc service
 * @name redqueenUiApp.RfidCard
 * @description
 * # RfidCard
 * Service in the redqueenUiApp.
 */
angular.module('redqueenUiApp')
  .service('RfidCard', [ '$q', '$timeout', '$http', 'underscore', function($q, $timeout, $http, _) {

    function RfidCard(data) {
      angular.extend(this, data);

      this.$isNew = (typeof(this.id) === 'undefined' || !this.id);
    }

    RfidCard.all = function RfidCardResourceAll(page, per_page) {
      var deferred = $q.defer();

      $http({
          url: '/api/cards',
          method: 'GET',
          params: { 
              'page' : page,
              'per_page': per_page
          }
      }).then(function(data) {
        deferred.resolve(data.data);
      }, function() {
        deferred.reject();
      });

      return deferred.promise;
    };

    RfidCard.find = function RfidCardResourceFind(id) {
      var deferred = $q.defer();

      $http.get('/api/cards/' + id).then(function(data) {
        var rfidCard = new RfidCard(data.data);

        deferred.resolve(rfidCard);
      }, function() {
        deferred.reject();
      });

      return deferred.promise;
    };

    RfidCard.prototype.$save = function RfidCardSave() {
      var deferred = $q.defer();
      var self = this;
      var url = null;
      var method = null;

      var data = {
        name: self.name,
        is_active: self.is_active,
        schedules: _.map(self.schedules, function(s) {
          return { 'id': s.id };
        })
      };

      if (self.$isNew) {
        url = '/api/cards';
        method = 'POST';

        data.pin = self.pin;
        data.code = self.code;
      } else {
        url = '/api/cards/' + self.id;
        method = 'PUT';

        if (self.pin) {
          data.pin = self.pin;
        }
      }

      $http({
        url: url,
        method:  method,
        data: data
      }).then(function(data) {
        var rfidCard = new RfidCard(data.data);

        deferred.resolve(rfidCard);
      }, function() {
        deferred.reject();
      });

      return deferred.promise;
    };

    //RfidCard.prototype.$remove = function RfidCardRemove() {
    //  var deferred = $q.defer();
    //  var self = this;
    //
    //  if (!self.$isNew) {
    //    $http.delete('/api/cards/' + self.id).then(function() {
    //      deferred.resolve();
    //    }, function() {
    //      deferred.reject();
    //    });
    //  }
    //
    //  return deferred.promise;
    //};

    return RfidCard;
  }]);
