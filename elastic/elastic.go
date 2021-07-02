package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	es "github.com/olivere/elastic/v7"
	"reflect"
	"time"
)

func onError(err error) {
	fmt.Println(err)
}

type ElasticClient interface {
	ElasticBasicActions
	ElasticBulkActions
}

type ElasticBasicActions interface {
	Search(ctx context.Context, indexName string, option SearchOption) (result []byte, err error)
	Store(ctx context.Context, name string, doc interface{}, template *DynamicTemplate) (res *es.IndexResponse, err error)
	Remove(ctx context.Context, indexName string, id string) (res *es.DeleteResponse, err error)
	RemoveIndex(ctx context.Context, indexName ...string) (res *es.IndicesDeleteResponse, err error)
	Ping(ctx context.Context, nodeURL string) (*es.PingResult, int, error)
}

type ElasticBulkActions interface {
	AddBulkProcessor(bulkProcessor BulkProcessor) (err error)
	BulkStore(ctx context.Context, indexName string, processorName string, docs []interface{}, template *DynamicTemplate) (err error)
}

type BulkProcessor struct {
	Name          string
	Workers       int
	BulkActions   int
	BulkSize      int
	FlushInterval time.Duration
}

type Client struct {
	esclient *es.Client
	Config   ClientConfig
}

type ClientConfig struct {
	BulkProcessors map[string]*es.BulkProcessor
}

type SearchOption struct {
	Query es.Query
	Sort  map[string]bool
	From  int
	Size  int
}

func NewClient(urls ...string) (c ElasticClient, err error) {
	esclient, err := es.NewClient(
		es.SetURL(urls...),
		es.SetSniff(false),
		es.SetHealthcheck(false),
		es.SetRetrier(NewElasticRetrier(3*time.Second, onError)),
	)
	if err != nil {
		return
	}

	client := &Client{
		esclient: esclient,
	}

	client.defaultBulkProcessor()
	c = client

	return
}

func (c *Client) Search(ctx context.Context, indexName string, option SearchOption) (result []byte, err error) {
	if option.Size == 0 {
		option.Size = 10 // default size to 10
	}

	search := c.esclient.Search().
		Index(indexName).
		Query(option.Query).
		From(option.From).
		Size(option.Size)

	for k, v := range option.Sort {
		search = search.Sort(k, v)
	}

	searchResult, err := search.Do(ctx)
	if err != nil {
		return
	}

	var source map[string]interface{}
	if searchResult.Hits.TotalHits.Value > 0 {
		for _, hit := range searchResult.Hits.Hits {
			err = json.Unmarshal(hit.Source, &source)
			if err != nil {
				return
			}
		}
	}

	result, err = json.Marshal(source)
	if err != nil {
		return
	}

	return
}

func (c *Client) createMappings(ctx context.Context, indexName string, template *DynamicTemplate) (err error) {
	exist, err := c.esclient.IndexExists(indexName).Do(ctx)
	if err != nil {
		return
	}

	if !exist {
		// create mapping
		_, err = c.esclient.CreateIndex(indexName).
			BodyJson(template).
			Index(indexName).
			Do(ctx)
		if err != nil {
			return
		}
	}

	return

}

// Index index(upsert) document to elastic
func (c *Client) Store(ctx context.Context, name string, doc interface{}, template *DynamicTemplate) (res *es.IndexResponse, err error) {
	err = c.createMappings(ctx, name, template)
	if err != nil {
		return
	}

	id := getDocumentID(doc)

	res, err = c.esclient.Index().
		Index(name).
		Id(id).
		BodyJson(doc).
		Do(ctx)

	return
}

// Delete delete document to elastic
func (c *Client) Remove(ctx context.Context, indexName string, id string) (res *es.DeleteResponse, err error) {
	res, err = c.esclient.Delete().
		Index(indexName).
		Id(id).
		Do(ctx)
	return
}

// Delete delete document to elastic
func (c *Client) RemoveIndex(ctx context.Context, indexName ...string) (res *es.IndicesDeleteResponse, err error) {
	res, err = c.esclient.DeleteIndex(indexName...).
		Do(ctx)
	return
}

func findFieldID(fieldName string) bool {
	return fieldName == "ID" || fieldName == "id" || fieldName == "Id"
}

func getDocumentID(doc interface{}) (id string) {
	val := reflect.ValueOf(doc)
	switch val.Kind() {
	case reflect.Struct:
		fieldByname := val.FieldByNameFunc(findFieldID)

		if fieldByname.IsValid() {
			id = fmt.Sprintf("%s", fieldByname)
		}

		return
	case reflect.Map:
		value := val.MapIndex(reflect.ValueOf("id"))
		if value.IsValid() {
			id = fmt.Sprintf("%s", value)
		}
	}

	return
}

// BulkStore bulk index(upsert) document to elastic
func (c *Client) BulkStore(ctx context.Context, indexName string, processorName string, docs []interface{}, template *DynamicTemplate) (err error) {
	processor, err := c.GetBulkProcessor(processorName)
	if err != nil {
		return
	}

	err = c.createMappings(ctx, indexName, template)
	if err != nil {
		return
	}

	for _, doc := range docs {
		id := getDocumentID(doc)
		bulkUpdateReq := es.NewBulkIndexRequest().Type("_doc").Index(indexName).Id(id).Doc(doc)
		processor.Add(bulkUpdateReq)
	}

	return
}

func (c *Client) newBulkProcessor(processor BulkProcessor) (bulkProcessor *es.BulkProcessor, err error) {
	bulkProcessor, err = c.esclient.BulkProcessor().
		Name(processor.Name).
		Workers(processor.Workers).
		BulkActions(processor.BulkActions).     // commit if # requests reach certain of number
		BulkSize(processor.BulkSize).           // commit when document size reach certain size
		FlushInterval(processor.FlushInterval). // commit every interval of time
		Do(context.Background())
	return
}

func (c *Client) defaultBulkProcessor() {
	bulkProcessor := BulkProcessor{
		Name:          "default",
		Workers:       10,
		BulkActions:   1000,            // flush when reach 1000 requests
		BulkSize:      2 << 20,         // flush when reach 2 MB
		FlushInterval: 1 * time.Second, // flush every 1 seconds
	}

	if c.Config.BulkProcessors == nil {
		c.Config.BulkProcessors = make(map[string]*es.BulkProcessor)
	}

	defaultBulkProcessor, _ := c.newBulkProcessor(bulkProcessor)
	c.Config.BulkProcessors[bulkProcessor.Name] = defaultBulkProcessor
}

func (c *Client) AddBulkProcessor(bulkProcessor BulkProcessor) (err error) {
	processor, err := c.newBulkProcessor(bulkProcessor)
	if err != nil {
		return err
	}

	c.Config.BulkProcessors[bulkProcessor.Name] = processor
	return
}

func (c *Client) GetBulkProcessor(name string) (processor *es.BulkProcessor, err error) {
	if name == "" {
		name = "default"
	}

	processor, ok := c.Config.BulkProcessors[name]
	if !ok {
		err = fmt.Errorf("Bulk processor with name %s not found", name)
		return
	}

	return
}

func (c *Client) Ping(ctx context.Context, nodeURL string) (*es.PingResult, int, error) {
	return c.esclient.Ping(nodeURL).Do(ctx)
}
