export async function getLongUrl(longUrl:string):Promise<{shortUrl:string}>{

    const req = await fetch("/api/shorten",{
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url: longUrl })
    });

    if(!req.ok)throw new Error("Didnt get shortUrl");

    return req.json();
}