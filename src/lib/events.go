package events

import (
	"github.com/asaskevich/EventBus"
)

const (
	COMPLETED_MAPPING_LOAD   = "mapping.load.completed"
	COMPLETED_SEARCH_LOAD    = "search.load.completed"
	COMPLETED_SEASONAL_LOAD  = "seasonal.load.completed"
	COMPLETED_ENTRY_CREATION = "entry.creation.completed"
)

var Bus = EventBus.New()
