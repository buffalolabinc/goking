'use strict';

/**
 * @ngdoc service
 * @name redqueenUiApp.Schedule
 * @description
 * # Schedule
 * Service in the redqueenUiApp.
 */
angular.module('redqueenUiApp')
  .service('Schedule', [ '$q', '$timeout', '$http', 'underscore', 'moment', function($q, $timeout, $http, _, moment) {

    function Schedule(data) {
      angular.extend(this, data);

      this.$isNew = (typeof(this.id) === 'undefined' || !this.id);
    }

    Schedule.all = function ScheduleResourceAll(page, per_page) {
      var deferred = $q.defer();

      $http({
          url: '/api/schedules',
          method: 'GET',
          params: {
              'page': page,
              'per_page': per_page
          }
      }).then(function(data) {
          var schedules = data.data;
        deferred.resolve(schedules);
      }, function() {
        deferred.reject();
      });

      return deferred.promise;
    };

    Schedule.find = function ScheduleResourceFind(id) {
      var deferred = $q.defer();

      $http.get('/api/schedules/' + id).then(function(data) {
        var schedule = new Schedule(data.data);

	schedule.start_time = moment(schedule.start_time).toDate();
	schedule.end_time = moment(schedule.end_time).toDate();

        deferred.resolve(schedule);
      }, function() {
        deferred.reject();
      });

      return deferred.promise;
    };

    Schedule.prototype.$save = function ScheduleSave() {
      var deferred = $q.defer();
      var self = this;
      var url = null;
      var method = null;

      var data = {
        name: self.name,
        mon: self.mon === true,
        tue: self.tue === true,
        wed: self.wed === true,
        thu: self.thu === true,
        fri: self.fri === true,
        sat: self.sat === true,
        sun: self.sun === true,
        start_time: self.start_time,
        end_time: self.end_time
      };

      if (self.$isNew) {
        url = '/api/schedules';
        method = 'POST';
      } else {
        url = '/api/schedules/' + self.id;
        method = 'PUT';
      }

      $http({
        url: url,
        method:  method,
        data: data
      }).then(function(data) {
        var schedule = new Schedule(data.data);

        deferred.resolve(schedule);
      }, function() {
        deferred.reject();
      });

      return deferred.promise;
    };

    return Schedule;
  }]);
