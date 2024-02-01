let text = Deno.env.get("DATA");
let data;
data = JSON.parse(text);

try {
    const a = await Deno.resolveDns(data.host, "A");
}
catch (error){
    console.error("Could not resolve hostname")
    Deno.exit(1)
}