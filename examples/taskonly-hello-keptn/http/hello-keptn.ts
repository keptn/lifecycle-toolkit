let text = Deno.env.get("DATA");
let data;
let name;
data = JSON.parse(text);

name = data.name
console.log("Hello, " + name );