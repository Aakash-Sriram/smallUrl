import { useState } from "react";
import { getLongUrl } from "../api/shorten" 
function Form() {
    const [url,getUrl] = useState("")
    const [shortUrl, setShortUrl] = useState("");
    const handler = async(e:React.FormEvent) =>{
        e.preventDefault();
        const resp = await getLongUrl(url);
        if(!resp)throw new Error("Did'nt get shortUrl response");
        console.log(""+resp.shortUrl)
        setShortUrl(resp.shortUrl)
    }
    return (
        <>
            <form onSubmit={handler}>
                <label htmlFor="longUrl">
                    Enter URL
                </label>
                <input id="longUrl" type="text" required value={url}onChange={(e)=>{getUrl(e.target.value)}}></input>
                <button type="submit">get shortUrl</button>
            </form>
            {shortUrl && (
                <div>
                    <a href="{shortUrl}">{shortUrl}</a>
                </div>
                )

            }
        </> 
    )
}

export default Form
