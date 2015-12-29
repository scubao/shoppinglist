/**
 * Created by oliver on 28.12.15.
 */

var app = angular.module('httptest', ['ngMaterial']);

app.controller('getShoppingEntries', ['$scope', '$http', function ($scope, $http) {
    $http.get('http://127.0.0.1:8080/entries').success(function (data) {
        $scope.data = data;
        $scope.save = function(entry) {
            $scope.message = "Save " + entry.done;
        }
        $scope.checkEnter = function($event) {
            var keyCode = $event.which || $event.keyCode;
            if (keyCode == 13) {
                $scope.newproduct = $scope.hanspeter;
            }
        }
    });
}]);
