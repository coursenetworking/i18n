var i18n = angular.module("i18n", []);
i18n.controller("sectionCtrl", ["$scope",
    function($scope){
        $scope.data = {
            message: "Hello"
        };
    }]
);
