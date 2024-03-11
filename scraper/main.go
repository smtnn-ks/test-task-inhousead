package scraper

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/smtnn-ks/test-task-inhousead/db"
)

type content_t struct {
	Item_title      string `json:"name"`
	Category1_title string `json:"category"`
	Category2_title string `json:"group"`
}

func Init() {
	ticker := time.NewTicker(time.Hour)
	url := "https://emojihub.yurace.pro/api/all"

	for {
		log.Println("[SCRAPER INFO]  Starts scraping...")
		res, err := http.Get(url)
		if err != nil {
			log.Println("[SCRAPER ERROR] ", err)
		}

		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println("[SCRAPER ERROR] ", err)
		}

		var content []content_t
		err = json.Unmarshal(resBody, &content)
		if err != nil {
			log.Println("[SCRAPER ERROR] ", err)
		}

		for _, thing := range content {
			if _, err := db.Client.Exec(
				"SELECT update_content($1, $2, $3)",
				thing.Item_title,
				thing.Category1_title,
				thing.Category2_title,
			); err != nil {
				log.Println("[SCRAPER ERROR] ", err)
			}
		}

		log.Println("[SCRAPER INFO]  Scraping cycle done")
		<-ticker.C
	}
}
