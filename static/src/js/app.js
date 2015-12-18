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

i18n.controller("langCtrl", ["$scope", "$http", "$timeout", "API", function($scope, $http, $timeout, API){
    $scope.langs = {
        "af_ZA": "Afrikaans",
        "ar_AR": "العربية",
        "az_AZ": "Azərbaycan dili",
        "be_BY": "Беларуская",
        "bg_BG": "Български",
        "bn_IN": "বাংলা",
        "br_FR": "Brezhoneg",
        "bs_BA": "Bosanski",
        "ca_ES": "Català",
        "cb_IQ": "کوردیی ناوەندی",
        "cs_CZ": "Čeština",
        "cx_PH": "Bisaya",
        "cy_GB": "Cymraeg",
        "da_DK": "Dansk",
        "de_DE": "Deutsch",
        "el_GR": "Ελληνικά",
        "en_GB": "English (UK)",
        "en_PI": "English (Pirate)",
        "en_UD": "English (Upside Down)",
        // "en_US": "English (US)",
        "eo_EO": "Esperanto",
        "es_CO": "Español (Colombia)",
        "es_ES": "Español (España)",
        "es_LA": "Español",
        "et_EE": "Eesti",
        "eu_ES": "Euskara",
        "fa_IR": "فارسی",
        "fb_LT": "Leet Speak",
        "fi_FI": "Suomi",
        "fo_FO": "Føroyskt",
        "fr_CA": "Français (Canada)",
        "fr_FR": "Français (France)",
        "fy_NL": "Frysk",
        "ga_IE": "Gaeilge",
        "gl_ES": "Galego",
        "gn_PY": "Guarani",
        "gu_IN": "ગુજરાતી",
        "he_IL": "עברית",
        "hi_IN": "हिन्दी",
        "hr_HR": "Hrvatski",
        "hu_HU": "Magyar",
        "hy_AM": "Հայերեն",
        "id_ID": "Bahasa Indonesia",
        "is_IS": "Íslenska",
        "it_IT": "Italiano",
        "ja_JP": "日本語",
        "ja_KS": "日本語(関西)",
        "jv_ID": "Basa Jawa",
        "ka_GE": "ქართული",
        "kk_KZ": "Қазақша",
        "km_KH": "ភាសាខ្មែរ",
        "kn_IN": "ಕನ್ನಡ",
        "ko_KR": "한국어",
        "ku_TR": "Kurdî (Kurmancî)",
        "la_VA": "lingua latina",
        "lt_LT": "Lietuvių",
        "lv_LV": "Latviešu",
        "mk_MK": "Македонски",
        "ml_IN": "മലയാളം",
        "mn_MN": "Монгол",
        "mr_IN": "मराठी",
        "ms_MY": "Bahasa Melayu",
        "my_MM": "မြန်မာဘာသာ",
        "nb_NO": "Norsk (bokmål)",
        "ne_NP": "नेपाली",
        "nl_BE": "Nederlands (België)",
        "nl_NL": "Nederlands",
        "nn_NO": "Norsk (nynorsk)",
        "or_IN": "ଓଡ଼ିଆ",
        "pa_IN": "ਪੰਜਾਬੀ",
        "pl_PL": "Polski",
        "ps_AF": "پښتو",
        "pt_BR": "Português (Brasil)",
        "pt_PT": "Português (Portugal)",
        "ro_RO": "Română",
        "ru_RU": "Русский",
        "rw_RW": "Ikinyarwanda",
        "si_LK": "සිංහල",
        "sk_SK": "Slovenčina",
        "sl_SI": "Slovenščina",
        "sq_AL": "Shqip",
        "sr_RS": "Српски",
        "sv_SE": "Svenska",
        "sw_KE": "Kiswahili",
        "ta_IN": "தமிழ்",
        "te_IN": "తెలుగు",
        "tg_TJ": "Тоҷикӣ",
        "th_TH": "ภาษาไทย",
        "tl_PH": "Filipino",
        "tr_TR": "Türkçe",
        "uk_UA": "Українська",
        "ur_PK": "اردو",
        "uz_UZ": "O'zbek",
        "vi_VN": "Tiếng Việt",
        "zh_CN": "中文(简体)",
        "zh_HK": "中文(香港)",
        "zh_TW": "中文(台灣)"
    };
    $scope.lang = "zh_CN";
    console.log(API);
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

i18n.directive("langsectionform", function(){
    return {
        restrict: "AE",
        templateUrl: "tpl/lang-section-form.html",
        replace: true
    }
});

i18n.directive("langsectioncreate", function(){
    return {
        restrict: "AE",
        template: "<section class=\"panel panel-lang super-mode-show\"><div langsectionform></div></section>",
        replace: true
    }
});

i18n.directive("langsectionlist", function(){
    return {
        restrict: "AE",
        template: "<section class=\"panel panel-lang\" ng-repeat=\"(k, section) in sections|filter:searchSection as results\"><div langsectionform></div></section>",
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
