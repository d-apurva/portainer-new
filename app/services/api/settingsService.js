angular.module('click2cloud.services')
.factory('SettingsService', ['$q', 'Settings', function SettingsServiceFactory($q, Settings) {
  'use strict';
  var service = {};

  service.settings = function() {
    var deferred = $q.defer();

    Settings.get().$promise
    .then(function success(data) {
      var settings = new SettingsViewModel(data);
      deferred.resolve(settings);
    })
    .catch(function error(err) {
      deferred.reject({ msg: 'Unable to retrieve application settings', err: err });
    });

    return deferred.promise;
  };

  service.update = function(settings) {
    return Settings.update({}, settings).$promise;
  };

  return service;
}]);
