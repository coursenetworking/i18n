var i18n = angular.module("i18n", ["i18n.config"]);
i18n.directive("supermode", function() {
    return {
        restrict: "AE",
        replace: true,
        transclude: true,
        template: "<div class='super-mode' ng-class=\"{'active': supermode}\" ng-transclude></div>",
        controller: function($scope) {
            $scope.supermode = false;
        },
        link: function(scope, elem, attrs) {
            var keyMemos = [],
                timer = null;
            angular.element(window).bind("keydown", function(e) {
                keyMemos.push(e.keyCode);
                if (keyMemos.length == 6) {
                    scope.$apply(function() {
                        angular.equals(keyMemos, [49, 50, 51, 49, 50, 51]) && (scope.supermode = !scope.supermode);
                    });
                    keyMemos = [];
                }
                if (!timer) {
                    timer = setTimeout(function() {
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

i18n.controller("langCtrl", ["$scope", "$http", "$timeout", "API", "LANGUAGES", 'GOOGLE_MODE', function($scope, $http, $timeout, API, LANGUAGES, GOOGLE_MODE) {

    function setCookie(key, value, expiredays) {
        var exdate = new Date();
        exdate.setDate(exdate.getDate() + expiredays);
        document.cookie = key + "=" + escape(value) +
            ((expiredays === null) ? "" : ";expires=" + exdate.toGMTString());
    }

    function getCookie(key) {
        if (document.cookie.length > 0) {
            var start = document.cookie.indexOf(key + "=");
            if (start != -1) {
                start = start + key.length + 1;
                var end = document.cookie.indexOf(";", start);
                if (end == -1) end = document.cookie.length;
                return unescape(document.cookie.substring(start, end));
            }
        }
        return false;
    }

    $scope.langs = LANGUAGES;
    $scope.isIndex = GOOGLE_MODE;
    $scope.lang = getCookie('i18n.lang') || "zh_CN";

    $scope.fetch = function() {
        $timeout(function() {
            setCookie('i18n.lang', $scope.lang, 30);
            $http({
                method: 'GET',
                url: API.TRANSLATETION + $scope.lang
            }).then(function(res) {
                res = res.data || {};
                if (res.result) {
                    $scope.sections = res.data;
                }
            }, function() {

            });
        }, 0);
    };
    $scope.fetch();

    $scope.addSection = function() {
        $scope.sections.unshift({
            items: {},
            to_lang: $scope.lang,
            is_new: true
        });
    }

    $scope.$watch('searchTerm', function(newValue, oldValue) {
        if (oldValue == void(0) && newValue) {
            $scope.isIndex = false;
        }
    });

    $scope.search = function(seactionName) {
        $scope.searchTerm = seactionName;
    }

    $scope.sectionFilter = function(section) {
        if (section.is_new || $scope.searchTerm == "") {
            return true;
        }

        var reg = new RegExp($scope.searchTerm);
        if (section && reg.test(section.section)) {
            return true;
        }
        return false;
    }
}]);

i18n.directive("langselect", function() {
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

i18n.directive("langsectioncreate", function() {
    return {
        restrict: "AE",
        templateUrl: "tpl/lang-section-create.html",
        replace: true
    }
});

i18n.directive("langsectionedit", function() {
    return {
        restrict: "AE",
        templateUrl: "tpl/lang-section-edit.html",
        replace: true
    }
});

i18n.directive("langsectionview", function() {
    return {
        restrict: "AE",
        templateUrl: "tpl/lang-section-view.html",
        replace: true
    }
});

i18n.directive("langsectionform", function() {
    return {
        restrict: "AE",
        template: "<div ng-switch=\"section.is_new\"><div ng-switch-when=\"true\" langsectioncreate></div><div ng-switch-default langsectionedit></div></div>",
        replace: true
    }
});

i18n.directive("langsectionlist", function() {
    return {
        restrict: "AE",
        template: "<section class=\"panel panel-lang\" ng-repeat=\"(k, section) in sections|filter:sectionFilter\" ng-switch=\"supermode\"><div ng-switch-when=\"true\" langsectionform></div><div ng-switch-default langsectionview></div></section>",
        replace: true
    }
});

i18n.controller("sectionCtrl", ["$scope", "$http", "API", function($scope, $http, API) {
    $scope.section.rename_to = $scope.section.section;
    for (var source in $scope.section.items) {
        $scope.section.items[source].rename_to = source;
    }

    var _cacheSection = JSON.stringify($scope.section);
    $scope.saveSection = function() {
        //@TODO delete items.source
        var newItems = {};
        for (var s in $scope.section.items) {
            var item = $scope.section.items[s];
            if (item.is_new) {
                newItems[item.rename_to] = item;
                delete $scope.section.items[s];
            }
        }
        angular.extend($scope.section.items, newItems);
        $http({
            method: "POST",
            url: API.TRANSLATETION + $scope.lang + '/' + $scope.section.section + '?is_admin=1',
            data: $scope.section
        }).then(function(res) {
            if(!res.data.result) {
                window.alert(res.data.err);
            }
            $scope.section.is_new = false;
        }, function (res) {
            window.alert(res.status + ' : ' + res.statusText);
        });
        return false;
    };

    $scope.resetSection = function() {
        $scope.section = angular.extend($scope.section, JSON.parse(_cacheSection));
        return false;
    };

    $scope.addLang = function() {
        var newSection = $scope.section.items[""];
        if (newSection) {
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
}]);
