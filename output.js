// include /static/js/base.js
cmsjs.WithApp(function(app) {

    app.service('WidgetApiController', ['$http', function($http) {
        var apiPost = function(url, data, context, successCallback) {
            $http.post(url, data, context)
                .success(function(resp, status, headers, ctx) {
                    if (resp.Success) {
                        if (successCallback) {
                            successCallback(resp, ctx);
                        }
                    } else {
                        alert(resp.Error);
                    }
                })
                .error(function(resp) {
                    alert(resp);
                });
        };
        return {
            Post: function(endpoint, data, context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/' + endpoint), data, context, successCallback);
            },
            
            SendMessageToAllUsers: function(<[]*ast.Ident Value>, context, successCallback) {
                apiPost(ApiUrl('/SendMessageToAllUsers/' + endpoint), data, context, successCallback);
            },
            
        };
    }]);

});
