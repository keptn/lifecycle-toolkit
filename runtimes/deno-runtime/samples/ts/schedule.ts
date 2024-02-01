let text = Deno.env.get("DATA");
let data;
if (text != "") {
    data = JSON.parse(text);
}

let targetDate = new Date(data.targetDate)
let dateTime = new Date();

if(targetDate < dateTime){
    console.log("Date has passed - ok");
    Deno.exit(0);
} else {
    console.log("It's too early - failing");
    Deno.exit(1);
}

console.log(targetDate);


