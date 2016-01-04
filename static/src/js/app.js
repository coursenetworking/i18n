var i18n = angular.module("i18n", ["i18n.config"]);
i18n.directive("supermode", function(){
    return {
        restrict: "AE",
        replace: true,
        transclude: true,
        template: "<div class='super-mode' ng-class=\"{'active': supermode}\" ng-transclude></div>",
        controller: function($scope){
            $scope.supermode = false;
        },
        link: function(scope, elem, attrs){
            var keyMemos = [], timer = null;
            angular.element(window).bind("keydown", function(e){
                keyMemos.push(e.keyCode);
                if(keyMemos.length == 6) {
                    scope.$apply(function() {
                        angular.equals(keyMemos, [49,50,51,49,50,51]) && (scope.supermode = !scope.supermode);
                    });
                    keyMemos = [];
                }
                if(!timer) {
                    timer = setTimeout(function(){
                        keyMemos = [];
                        timer = null;
                    }, 1.5e3);
                }
                return true;
            });
            return true;
        }
    }
});

i18n.controller("langCtrl", ["$scope", "$http", "$timeout", "API", "LANGUAGES", function($scope, $http, $timeout, API, LANGUAGES){
    $scope.langs = LANGUAGES;
    $scope.lang = "zh_CN";
    $scope.fetch = function(){
        $timeout(function(){
            $http({
                method: 'GET',
                url: API.TRANSLATETION + $scope.lang
            }).then(function(res){
                res = res.data || {};
                if(res.result) {
                    $scope.sections = res.data;
                }
            }, function(){

            });
        }, 0);
    };
    $scope.fetch();

    $scope.addSection = function() {
        $scope.sections.unshift({
            items: {},
            rename_to: 'section'+(+new Date()),
            to_lang: $scope.lang
        });
    }
}]);

i18n.directive("langselect", function(){
    return {
        restrict: "AE",
        scope: {
            lang: "=",
            langs: "@langs",
            fetch: "&changeLang"
        },
        template: '<select ng-model="lang" ng-transclude ng-change="fetch()"></select>',
        replace: true,
        transclude: true
    }
});

i18n.directive("langsectioncreate", function(){
    return {
        restrict: "AE",
        templateUrl: "tpl/lang-section-create.html",
        replace: true
    }
});

i18n.directive("langsectionedit", function(){
    return {
        restrict: "AE",
        templateUrl: "tpl/lang-section-edit.html",
        replace: true
    }
});

i18n.directive("langsectionview", function(){
    return {
        restrict: "AE",
        templateUrl: "tpl/lang-section-view.html",
        replace: true
    }
});

i18n.directive("langsectioncreate", function(){
    return {
        restrict: "AE",
        template: "<section class=\"panel panel-lang\"><div langsectioncreate></div></section>",
        replace: true
    }
});

i18n.directive("langsectionlist", function(){
    return {
        restrict: "AE",
        template: "<section class=\"panel panel-lang\" ng-repeat=\"(k, section) in sections|filter:searchSection as results\" ng-switch=\"supermode\"><div ng-switch-when=\"true\" langsectionedit></div><div ng-switch-default langsectionview></div></section>",
        replace: true
    }
});

i18n.controller("sectionCtrl", ["$scope", "$http", "API", function($scope, $http, API){
        var _cacheSection = angular.extend({}, $scope.section);
        $scope.saveSection = function() {
            //@TODO delete items.source
            var newItems = {};
            for (var s in $scope.section.items) {
                var item = $scope.section.items[s];
                if(item.is_new) {
                    newItems[item.rename_to] = item;
                    delete $scope.section.items[s];
                }
            }
            angular.extend($scope.section.items, newItems);
            $http({
                method: "POST",
                url: API.TRANSLATETION + $scope.lang + '/' + $scope.section.section,
                data: $scope.section
            }).then(function(){
                alert('SAVE SUCCESSFULLY.');
            });
            return false;
        };

        $scope.resetSection = function() {
            $scope.section = angular.extend($scope.section, _cacheSection);
            return false;
        };

        $scope.addLang = function() {
            var newSection = $scope.section.items[""];
            if(newSection) {
                $scope.section.items[newSection.rename_to] = newSection;
                delete $scope.section.items[""];
            }
            $scope.section.items[""] = {
                rename_to: "",
                translate_to: "",
                is_new: true
            }
        }

        $scope.deleteLang = function(source) {
            $scope.section.items[source] && delete $scope.section.items[source];
        }
    }]
);
