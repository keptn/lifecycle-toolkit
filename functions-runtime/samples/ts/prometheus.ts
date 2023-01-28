const text = Deno.env.get("DATA")!;
let data;
if (text && text != "") {
  try {
    data = JSON.parse(text);
    if (!data.metrics || !data.url || !data.expected_value) {
      console.log("Missing mandatory arguments.");
      printUsage();
    }
  } catch {
    console.error("Error Parsing Json => ", text);
    Deno.exit(1);
  }
} else {
  console.log("Missing mandatory enviornment variable DATA.");
  printUsage();
}

function printUsage() {
  console.log(
    'Expecting environment variable DATA in this format => DATA=\'{ "url":"<PROMETHEUS_URL>", "metrics": "<PROMETHEUS_QUERY>", "expected_value": "<EXPECTED_VALUE>" }\' ',
  );
  console.log(
    'Example: export DATA=\'{ "url":"http://localhost:9090", "metrics": "up{service=\\"kubernetes\\"}", "expected_value": "1" }\' ',
  );
  Deno.exit(1);
}

function getPrometheusURL(url: string, metrics: string): string {
  const dateTime = new Date().toISOString();
  const hasPort = url.includes(":9090");
  const hasProtocol = url.includes("http://");
  let queryURL: string = url;
  if (!hasPort) {
    queryURL = queryURL + ":9090";
  }
  if (!hasProtocol) {
    queryURL = "http://" + queryURL;
  }

  return queryURL + "/api/v1/query?query=" + metrics + "&time=" + dateTime;
}

const promtheusURL = getPrometheusURL(data.url, data.metrics);
console.log("Prometheus URL => " + promtheusURL);
let value;
try {
  const jsonResponse = await fetch(promtheusURL);
  const jsonData = await jsonResponse.json();
  value = jsonData.data.result[0].value[1];
} catch (error) {
  console.error("Could not fetch data from Prometheus.\n", error);
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
