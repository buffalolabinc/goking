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
      this.$isNew = (typeof(this.Id) === 'undefined' || !this.Id);

    }

    Schedule.all = function ScheduleResourceAll() {
      var deferred = $q.defer();

      $http.get('/api/schedules').then(function(data) {
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

        var sTime = moment(schedule.StartTime);
        var eTime = moment(schedule.EndTime);

        schedule.StartTime = sTime.format('HH:mm:ss');
        schedule.EndTime = eTime.format('HH:mm:ss');

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

      var fixTime = function(time) {
        var t = time.length < 8 ? time + ':00' : time;

        // make up some date because we dont care anyway
        return moment('2015-05-12 '+t).toISOString();
      };

      var data = {
        Name: self.Name,
        Mon: self.Mon === true,
        Tue: self.Tue === true,
        Wed: self.Wed === true,
        Thu: self.Thu === true,
        Fri: self.Fri === true,
        Sat: self.Sat === true,
        Sun: self.Sun === true,
        StartTime: fixTime(self.StartTime),
        EndTime: fixTime(self.EndTime)
      };

      if (self.$isNew) {
        url = '/api/schedules';
        method = 'POST';
      } else {
        url = '/api/schedules/' + self.Id;
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
