'use strict';

/**
 * @ngdoc service
 * @name redqueenUiApp.moment
 * @description
 * # moment
 * Service in the redqueenUiApp.
 */
angular.module('redqueenUiApp')
  .service('moment', function moment() {
    return window.moment;
  });
