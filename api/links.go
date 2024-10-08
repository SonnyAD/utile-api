package api

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"utile.space/api/utils"
)

// @Summary		Get Recommended Links Page
// @Description	Returns a page of recommended links by SonnyAD
// @Tags			links
// @Produce		json,xml,application/yaml,plain
// @Param			start	query		string	false	"Start cursor for pagination"
// @Param			search	query		string	false	"Search filter"
// @Success		200		{object}	LinksPage
// @Router			/links [get]
func GetLinksPage(w http.ResponseWriter, r *http.Request) {
	databaseID := os.Getenv("NOTION_DATABASE_ID")
	notionAPISecret := os.Getenv("NOTION_SECRET")

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	start := r.URL.Query().Get("start")
	searchQuery := r.URL.Query().Get("search")
	var filter string

	// TODO: This looks really bad, I should likely make some struct for that
	if searchQuery != "" {
		if start != "" {
			filter = "{\"filter\":{\"and\":[{\"property\":\"On Telegram?\",\"checkbox\":{\"equals\":true}},{\"or\":[{\"property\":\"URL\",\"text\":{\"contains\":\"" + searchQuery + "\"}},{\"property\":\"URL\",\"text\":{\"contains\":\"" + searchQuery + "\"}}]}]},\"start_cursor\":\"" + start + "\"}"
		} else {
			filter = "{\"filter\":{\"and\":[{\"property\":\"On Telegram?\",\"checkbox\":{\"equals\":true}},{\"or\":[{\"property\":\"Description\",\"text\":{\"contains\":\"" + searchQuery + "\"}},{\"property\":\"URL\",\"text\":{\"contains\":\"" + searchQuery + "\"}}]}]}}"
		}
	} else {
		if start != "" {
			filter = "{\"filter\":{\"property\":\"On Telegram?\",\"checkbox\":{\"equals\":true}},\"start_cursor\":\"" + start + "\"}"
		} else {
			filter = "{\"filter\":{\"property\":\"On Telegram?\",\"checkbox\":{\"equals\":true}}}"
		}
	}

	request, err := http.NewRequestWithContext(r.Context(), "POST", "https://api.notion.com/v1/databases/"+databaseID+"/query", strings.NewReader(filter))

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Authorization", "Bearer "+notionAPISecret)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Notion-Version", "2021-08-16")

	resp, err := client.Do(request)

	if err != nil {
		log.Print(err)
		return
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	var res NotionResponse

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Print(err)
		return
	}

	var reply LinksPage

	reply.NextPage = res.NextCursor

	links := make([]Link, len(res.Results))

	for i, v := range res.Results {
		links[i].URL = v.Properties.URL.URL

		if len(v.Properties.Description.Title) > 0 {
			links[i].Description = v.Properties.Description.Title[0].PlainText
		}

		tags := make([]Tag, len(v.Properties.Tags.MultiSelect))

		for j, tag := range v.Properties.Tags.MultiSelect {
			tags[j].Name = tag.Name
			tags[j].Color = tag.Color
		}
		links[i].Tags = tags
	}

	reply.Links = links

	utils.Output(w, r.Header["Accept"], reply, "Please use JSON")
}

type LinksPage struct {
	XMLName  xml.Name `json:"-" xml:"page" yaml:"-"`
	Links    []Link   `json:"links" xml:"link" yaml:"links"`
	NextPage string   `json:"next" xml:"next" yaml:"next"`
}

type Link struct {
	URL         string `json:"url" xml:"url" yaml:"url"`
	Description string `json:"description" xml:"description" yaml:"description"`
	Tags        []Tag  `json:"tags" xml:"tag" yaml:"tags"`
}

type Tag struct {
	Name  string `json:"name" xml:"name" yaml:"name"`
	Color string `json:"color" xml:"color" yaml:"color"`
}

type NotionResponse struct {
	Object  string `json:"object"`
	Results []struct {
		Object         string      `json:"object"`
		ID             string      `json:"id"`
		CreatedTime    time.Time   `json:"created_time"`
		LastEditedTime time.Time   `json:"last_edited_time"`
		Cover          interface{} `json:"cover"`
		Icon           interface{} `json:"icon"`
		Parent         struct {
			Type       string `json:"type"`
			DatabaseID string `json:"database_id"`
		} `json:"parent"`
		Archived   bool `json:"archived"`
		Properties struct {
			Tags struct {
				ID          string      `json:"id"`
				Type        string      `json:"type"`
				MultiSelect []NotionTag `json:"multi_select"`
			} `json:"Tags"`
			URL struct {
				ID   string `json:"id"`
				Type string `json:"type"`
				URL  string `json:"url"`
			} `json:"URL"`
			Description struct {
				ID    string `json:"id"`
				Type  string `json:"type"`
				Title []struct {
					Type string `json:"type"`
					Text struct {
						Content string      `json:"content"`
						Link    interface{} `json:"link"`
					} `json:"text"`
					Annotations struct {
						Bold          bool   `json:"bold"`
						Italic        bool   `json:"italic"`
						Strikethrough bool   `json:"strikethrough"`
						Underline     bool   `json:"underline"`
						Code          bool   `json:"code"`
						Color         string `json:"color"`
					} `json:"annotations"`
					PlainText string      `json:"plain_text"`
					Href      interface{} `json:"href"`
				} `json:"title"`
			} `json:"Description"`
		} `json:"properties"`
		URL string `json:"url"`
	} `json:"results"`
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

type NotionTag struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
