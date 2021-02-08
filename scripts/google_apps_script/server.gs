//For Slack
function doPost(e) {
    console.log("doPost");
    var token = "your_slack_token";

    if (token != e.parameter.token) {
        return;
    }
    var text = e.parameter.text.replace(/<@[a-zA-Z0-9].*?>/, '').slice(0, 50);
    return request(text);
}

//For Test
function doGet(e) {
    var params = JSON.stringify(e);
    //return HtmlService.createHtmlOutput(params);
    return ContentService.createTextOutput(JSON.stringify(params))
        .setMimeType(ContentService.MimeType.JSON);
}

function request(text) {
    console.log("request");
    var googleHomeURL = 'https://xxxxx.ngrok.io/google-home-notifier';
    var urlFetchOption = {
        'method' : 'post',
        'contentType' : 'application/x-www-form-urlencoded',
        'payload' : { 'text' : text}
    };

    var response = UrlFetchApp.fetch(googleHomeURL, urlFetchOption);
    return response;
}
