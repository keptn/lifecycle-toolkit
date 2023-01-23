let text = Deno.env.get("DATA");
let data;
if (text != "") {
  data = JSON.parse(text);
}

function getPrometheusURL(url: string, metrics: string): string {
  let dateTime = new Date().toISOString();
  let hasPort = url.includes(":9090");
  let hasProtocol = url.includes("http://");
  let queryURL: string = url;
  if (!hasPort) {
    queryURL = queryURL + ":9090";
  }
  if (!hasProtocol) {
    queryURL = "http://" + queryURL;
  }

  return queryURL + "/api/v1/query?query=" + metrics + "&time=" + dateTime;
}

let promtheusURL = getPrometheusURL(data.url, data.metrics);
console.log("Prometheus URL => " + promtheusURL);
let value;
try {
  let jsonResponse = await fetch(promtheusURL);
  let jsonData = await jsonResponse.json();
  value = jsonData.data.result[0].value[1];
} catch (error) {
  console.error("Could not fetch ");
  Deno.exit(1);
}

console.log("Expected Value => " + data.expected_value + ", Value => " + value);
if (value == data.expected_value) {
  console.log("Evaluation Successful");
  Deno.exit(0);
} else {
  console.log("Evaluation Failed");
  Deno.exit(1);
}
