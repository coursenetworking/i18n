var config = {
    API: {
        TRANSLATETION: "http://localhost:8080/translation/"
    }
}

var mConfig = angular.module("i18n.config", []);
angular.forEach(config, function(key, value) {
    mConfig.constant(value, key);
});
