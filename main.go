package main

import (
	"anify/eltik/go/src/database"
	events "anify/eltik/go/src/lib"
	proxies "anify/eltik/go/src/lib/impl/proxies"
	"anify/eltik/go/src/lib/impl/request"
	"anify/eltik/go/src/mappings/impl/base"
	"anify/eltik/go/src/types"
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

	res, err := base.NewMangaDexBaseProvider().SearchAdvanced("", types.TypeManga, []types.Format{types.FormatManga}, 0, 25, []string{"Aliens"}, []string{"Harem"}, types.SeasonUnknown, 0, []string{"Isekai"}, []string{"Harem"})
	if err != nil {
		panic(err)
	}

	for _, r := range res {
		// print r.title.english
		println(*r.Title.English)
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
