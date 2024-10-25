package events

import (
	"fmt"
)

func Listen() {
	Bus.Subscribe(COMPLETED_MAPPING_LOAD, func() {
		fmt.Println("Entry creation completed!")
	})

	Bus.Subscribe(COMPLETED_ENTRY_CREATION, func() {
		fmt.Println("Entry creation completed!")
	})

	Bus.Subscribe(COMPLETED_SEARCH_LOAD, func() {
		fmt.Println("Search load completed!")
	})

	Bus.Subscribe(COMPLETED_SEASONAL_LOAD, func() {
		fmt.Println("Seasonal load completed!")
	})
}
