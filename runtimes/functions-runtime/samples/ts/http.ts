let text = Deno.env.get("DATA");
let data;
data = JSON.parse(text);

try {
    let resp = await fetch(data.url);
}
catch (error){
    console.error("Could not fetch url");
    Deno.exit(1);
}