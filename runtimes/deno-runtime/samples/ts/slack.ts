let text = Deno.env.get("SECURE_DATA");
let context = Deno.env.get("CONTEXT");
let data;
let body;
let contextdata;

if (text != undefined) {
    data = JSON.parse(text);
}

if (context != undefined) {
    contextdata = JSON.parse(context);
}

if (contextdata.objectType == "Application") {
    body = `{
            "text": "Application ${contextdata.appName}, Version ${contextdata.appVersion} has been deployed"
          }`
}

if (contextdata.objectType == "Workload") {
    body = `{
            "text": "Workload ${contextdata.workloadName}, Version ${contextdata.workloadVersion} in App ${contextdata.appName} has been deployed"
          }`
}

console.log(body)
let resp = await fetch("https://hooks.slack.com/services/" + data.slack_hook, {
    method: "POST",
    body,
});

console.log(resp)