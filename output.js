
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
        return {
            
            SendMessageToAllUsers: function(content, context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/SendMessageToAllUsers'+'?'+content), null, context, successCallback);
            },
            
            SendMessage: function(toUserID, content, relatedModel, relatedModelID, context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/SendMessage'+'?'+toUserID+'&'+content+'&'+relatedModel+'&'+relatedModelID), null, context, successCallback);
            },
            
            GetUserMessages: function(context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/GetUserMessages'+''), null, context, successCallback);
            },
            
            GetUserUnreadMessages: function(fromUserID, context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/GetUserUnreadMessages'+'?'+fromUserID), null, context, successCallback);
            },
            
            GetUserMessageConversations: function(context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/GetUserMessageConversations'+''), null, context, successCallback);
            },
            
            MarkReadForMessages: function(fromUserID, context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/MarkReadForMessages'+'?'+fromUserID), null, context, successCallback);
            },
            
            Post: function(endpoint, data, context, successCallback) {
                apiPost(ApiUrl('/WidgetApi/' + endpoint), data, context, successCallback);
            }
        };
    }]);
});
