package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type Proxy struct {
	ProviderID string `json:"providerId"`
	IP         string `json:"ip"`
}

var (
	BaseProxies  []Proxy
	AnimeProxies []Proxy
	MangaProxies []Proxy
	MetaProxies  []Proxy
)

func FetchCorsProxies() (map[string][]Proxy, error) {
	var wg sync.WaitGroup
	wg.Add(4)

	// Load proxies concurrently
	go func() {
		defer wg.Done()
		base, err := LoadProxies("./baseProxies.json")
		if err == nil {
			BaseProxies = append(BaseProxies, base...)
		}
	}()

	go func() {
		defer wg.Done()
		anime, err := LoadProxies("./animeProxies.json")
		if err == nil {
			AnimeProxies = append(AnimeProxies, anime...)
		}
	}()

	go func() {
		defer wg.Done()
		manga, err := LoadProxies("./mangaProxies.json")
		if err == nil {
			MangaProxies = append(MangaProxies, manga...)
		}
	}()

	go func() {
		defer wg.Done()
		meta, err := LoadProxies("./metaProxies.json")
		if err == nil {
			MetaProxies = append(MetaProxies, meta...)
		}
	}()

	wg.Wait()

	return map[string][]Proxy{
		"base":  BaseProxies,
		"anime": AnimeProxies,
		"manga": MangaProxies,
		"meta":  MetaProxies,
	}, nil
}

func LoadProxies(fileName string) ([]Proxy, error) {
	var proxies []Proxy

	// Check if the file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return proxies, nil
	}

	// Read the file content
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// Parse the JSON content
	var proxyData []Proxy
	if err := json.Unmarshal(data, &proxyData); err != nil {
		return nil, err
	}

	const batchSize = 100
	totalProxies := len(proxyData)
	currentIndex := 0

	// Process proxies in batches
	for currentIndex < totalProxies {
		batch := make([]Proxy, 0, batchSize)

		for i := 0; i < batchSize && currentIndex < totalProxies; i, currentIndex = i+1, currentIndex+1 {
			proxy := proxyData[currentIndex]

			if !strings.HasPrefix(proxy.IP, "http") {
				proxy.IP = "http://" + proxy.IP
			}

			batch = append(batch, proxy)
		}

		proxies = append(proxies, batch...)
	}

	fmt.Printf("Finished importing %d proxies from %s.\n", totalProxies, fileName)
	return proxies, nil
}
