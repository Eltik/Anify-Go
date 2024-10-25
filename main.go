package main

import (
	"example/user/hello/src/database"
	events "example/user/hello/src/lib"
	proxies "example/user/hello/src/lib/impl/proxies"
	"example/user/hello/src/lib/impl/request"
	"example/user/hello/src/mappings/impl/manga"
	"example/user/hello/src/types"
)

func main() {
	events.Listen()
	database.Connect()
	database.CreateTables()

	proxies.FetchCorsProxies()

	proxy := request.GetRandomUnbannedProxy("mangadex")
	if proxy != nil {
		println(*proxy)
	}

	res, err := manga.NewMangaDexProvider().Search("Mushoku Tensei", types.FormatManga, 0)
	if err != nil {
		panic(err)
	}

	for _, r := range res {
		println(r.Title)
	}

	/*
		resp, err := request.Request("novelupdates", false, http.Request{
			URL: &url.URL{
				Scheme: "https",
				Host:   "www.novelupdates.com",
			},
			Method: "GET",
		}, true)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
		log.Printf("Response Body: %s\n", string(body))
	*/
}
