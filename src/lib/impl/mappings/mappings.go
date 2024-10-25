package mappings

import (
	database_fetch "example/user/hello/src/database/impl/fetch"
	events "example/user/hello/src/lib"
	"example/user/hello/src/types"
	"log"
)

func LoadMappings(data struct {
	ID      string
	Type    types.Type
	Formats []types.Format
}) ([]types.Anime, []types.Manga, error) {
	existing, err := database_fetch.Get(data.ID, data.Type)
	if err != nil {
		log.Println("Failed to fetch existing data:", err)
		return nil, nil, err
	}

	if existing != nil {
		events.Bus.Publish(events.COMPLETED_MAPPING_LOAD)
		return nil, nil, nil
	}

	log.Println("No existing data found, fetching mappings")

	return nil, nil, nil
}
