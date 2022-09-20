let text = Deno.env.get("SECURE_DATA");
let data;
if (text != undefined) {
    data = JSON.parse(text);
}

const body = `{"text": "${data.text}"}`;

console.log(body)
let resp = await fetch("https://hooks.slack.com/services/" + data.slack_hook, {
    method: "POST",
    body,
});

console.log(resp)