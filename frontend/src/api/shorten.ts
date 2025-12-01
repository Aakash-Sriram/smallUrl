export async function getLongUrl(longUrl:string):Promise<{shortUrl:string}>{
    const url = import.meta.env.VITE_SMALL_URL;

    if(!url)throw new Error("URL not loaded");

    const req = await fetch(url,{
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url: longUrl })
    });

    if(!req.ok)throw new Error("Didnt get shortUrl");

    return req.json();
}