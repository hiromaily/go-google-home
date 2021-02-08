// For Slack
// https://api.slack.com/events/url_verification
// https://developers.google.com/apps-script/guides/web
//
// [POST]
// "body": {
// 	 "type": "url_verification",
// 	 "token": "xxxxxxxxxxxxxx",
// 	 "challenge": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
// }
function doPost(e) {
    console.log("doPost");
    const token = "your_slack_token";

    // validate token
    if (token != e.parameter.token) {
        return;
    }
    // request
    const text = e.parameter.text.replace(/<@[a-zA-Z0-9].*?>/, '').slice(0, 50);
    request(text);

    // response
    return HtmlService.createHtmlOutput(e.parameter.challenge);
}

// For Test
function doGet(e) {
    var params = JSON.stringify(e);
    return ContentService.createTextOutput(JSON.stringify(params))
        .setMimeType(ContentService.MimeType.JSON);
}

function request(text) {
    console.log("request");
    const googleHomeURL = 'https://xxxxx.ngrok.io/google-home-notifier';
    const urlFetchOption = {
        'method' : 'post',
        'contentType' : 'application/x-www-form-urlencoded',
        'payload' : { 'text' : text}
    };

    var response = UrlFetchApp.fetch(googleHomeURL, urlFetchOption);
    return response;
}
