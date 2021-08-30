package area

import (
	"encoding/json"
	"sync"
)

func init() {
	once := sync.Once{}
	once.Do(func() {
		if err := json.Unmarshal(pcasRaw, &pacs); err != nil {
			panic(err)
		}
	})
}
