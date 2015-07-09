
// include /static/js/base.js
cmsjs.WithApp(function(app) {
    
    app.service('WidgetApi', ['$http', function($http) {
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
        var _p = function(queryParam) {
            if (!queryParam) {
                return "";
            } else {
                return queryParam;
            }
        };
        return {
            
            SendMessageToAllUsers: function(content, context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/SendMessageToAllUsers'+'?'+_p(content)), null, context, successCallback);
            },
            
            SendMessage: function(toUserID, content, relatedModel, relatedModelID, context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/SendMessage'+'?'+_p(toUserID)+'&'+_p(content)+'&'+_p(relatedModel)+'&'+_p(relatedModelID)), null, context, successCallback);
            },
            
            GetUserMessages: function(context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/GetUserMessages'+''), null, context, successCallback);
            },
            
            GetUserUnreadMessages: function(fromUserID, context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/GetUserUnreadMessages'+'?'+_p(fromUserID)), null, context, successCallback);
            },
            
            GetUserMessageConversations: function(context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/GetUserMessageConversations'+''), null, context, successCallback);
            },
            
            MarkReadForMessages: function(fromUserID, context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/MarkReadForMessages'+'?'+_p(fromUserID)), null, context, successCallback);
            },
            
            Post: function(endpoint, data, context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/' + endpoint), data, context, successCallback);
            }
        };
    }]);
});
