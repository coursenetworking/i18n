var i18n = angular.module("i18n", []);

i18n.controller("langCtrl", ["$scope", "$http", function($scope, $http){
    $scope.fetch = function(){
        $http({
            method: 'GET',
            url: 'http://test.com/api/translation/en'
        }).then(function(res){
            console.log(res);
        }, function(){

        });
    }
}]);

i18n.directive("langselect", function(){
    return {
        restrict: "AE",
        template: '<select name="lang" id="lang-select"><option value="cn">cn-ZH</option><option value="en">English</option></select>',
        replace: true,
        link: function(scope, elem, attrs) {
            elem.bind('change', function(){
                scope.section.to_lang = elem.value;
            });
        }
    }
});
i18n.controller("sectionCtrl", ["$scope", function($scope){
        $scope.panel = {isShow: true};
        var originSection = {
            name: "index header",
            to_lang: "cn_zh",
            langs: {
                "search": "搜索",
                "post": "微博"
            }
        }
        $scope.section = $.extend({}, originSection);
        $scope.getFormData = function(){
            console.log($scope.section);
        }
        $scope.resetFromData = function() {
            console.log(originSection);
            $scope.section = originSection;
        }
        // View controller
        $scope.panelToggle = function(){
            $scope.panel.isShow = !$scope.panel.isShow;
        }
    }]
);
