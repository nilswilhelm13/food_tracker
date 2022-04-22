package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"strings"
)

var Client *Handler

type Handler struct {
	elasticClient elasticsearch.Client
	r             map[string]interface{}
}

func NewElasticHandler(elasticClient elasticsearch.Client) *Handler {
	return &Handler{elasticClient: elasticClient}
}

func (e Handler) Info() {
	log.SetFlags(0)

	e.r = make(map[string]interface{})

	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//

	// 1. Get cluster info
	//
	res, err := e.elasticClient.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&e.r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	fmt.Println("Client: %s", elasticsearch.Version)
	fmt.Println("Server: %s", e.r["version"].(map[string]interface{})["number"])
	fmt.Println(strings.Repeat("~", 37))
}
func (e Handler) IndexDocument(document interface{}, index string) {
	// TODO duplicates
	docJSON, _ := json.Marshal(document)
	// Set up the request object.
	req := esapi.IndexRequest{
		Index: index,
		//DocumentID: strconv.Itoa(1),
		Body:    strings.NewReader(string(docJSON)),
		Refresh: "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), &e.elasticClient)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), 1)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
}

func (e Handler) Search(query map[string]interface{}) (map[string]interface{}, error) {
	// 3. Search for the indexed documents
	//
	// Build the request body.
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := e.elasticClient.Search(
		e.elasticClient.Search.WithContext(context.Background()),
		e.elasticClient.Search.WithIndex("foodlist"),
		e.elasticClient.Search.WithBody(&buf),
		e.elasticClient.Search.WithTrackTotalHits(true),
		e.elasticClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&e.r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(e.r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(e.r["took"].(float64)),
	)
	return e.r, nil

}
