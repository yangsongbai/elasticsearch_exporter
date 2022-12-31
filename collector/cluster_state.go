package collector

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

//    ,
//    ,
//    METADATA_READ,
//    ;

const (
	READ           = "read"
	WRITE          = "write"
	METADATA_READ  = "metadata_read"
	METADATA_WRITE = "metadata_write"
	TRUE           = "true"
	FALSE          = "false"
	RED            = "red"
	YELLOW         = "yellow"
	GREEN          = "green"
)

type IndexStateStatus struct {
	EsCluster        string `json:"es_cluster"`
	Index            string `json:"index"`
	Color            string `json:"color"`
	State            string `json:"state"`
	NumberOfShards   string `json:"number_of_shards"`
	NumberOfReplicas string `json:"number_of_replicas"`
	Uuid             string `json:"uuid"`
	Dynamic          string `json:"dynamic"`
	CreationDate     string `json:"creation_date"`
	Read             string `json:"read"`
	Write            string `json:"write"`
	MetadataRead     string `json:"metadata_read"`
	MetadataWrite    string `json:"metadata_write"`
	value            float64
}

type Labels struct {
	keys   func(...string) []string
	values func(...string) []string
}

type indexStateMetric struct {
	Type   prometheus.ValueType
	Desc   *prometheus.Desc
	Value  func(result float64) float64
	Labels Labels
}

type ClusterState struct {
	logger                          log.Logger
	client                          *http.Client
	url                             *url.URL
	up                              prometheus.Gauge
	totalScrapes, jsonParseFailures prometheus.Counter

	metrics []*indexStateMetric
}

func NewClusterState(logger log.Logger, client *http.Client, url *url.URL) *ClusterState {
	indexStateLabels := Labels{
		keys: func(...string) []string {
			//status: green yellow red
			//state:open stated
			//dynamic：是否开启动态mapping
			//creation_date: 索引创建时间
			return []string{"index", "color", "state", "dynamic", "creation_date", "number_of_shards", "number_of_replicas", READ, WRITE, METADATA_READ, METADATA_WRITE, "es_cluster"}
		},
		values: func(s ...string) []string {
			return s
		},
	}
	subsystem := "index_health"
	return &ClusterState{
		logger: logger,
		client: client,
		url:    url,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, "up"),
			Help: "Was the last scrape of the ElasticSearch index health endpoint successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, "total_scrapes"),
			Help: "Current total ElasticSearch index health scrapes.",
		}),
		jsonParseFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, "json_parse_failures"),
			Help: "Number of errors while parsing JSON.",
		}),

		metrics: []*indexStateMetric{
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "status"),
					"The index status and state.",
					indexStateLabels.keys(), nil,
				),
				Value: func(result float64) float64 {
					return result
				},
				Labels: indexStateLabels,
			},
		},
	}
}

func (c *ClusterState) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.Desc
	}
	ch <- c.up.Desc()
	ch <- c.totalScrapes.Desc()
	ch <- c.jsonParseFailures.Desc()
}

func (c *ClusterState) fetchAndDecodeClusterState() (clusterStateResponse, error) {
	var csr clusterStateResponse
	u := *c.url
	//_cluster/state/routing_table,metadata
	u.Path = path.Join(u.Path, "_cluster/state")
	res, err := c.client.Get(u.String())
	if err != nil {
		return csr, fmt.Errorf("failed to get cluster state from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			_ = level.Warn(c.logger).Log(
				"msg", "failed to close http.Client",
				"err", err,
			)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return csr, fmt.Errorf("HTTP Request failed with code %d", res.StatusCode)
	}

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.jsonParseFailures.Inc()
		return csr, err
	}

	if err := json.Unmarshal(bts, &csr); err != nil {
		c.jsonParseFailures.Inc()
		return csr, err
	}
	return csr, nil
}
func (c *ClusterState) Collect(ch chan<- prometheus.Metric) {
	c.totalScrapes.Inc()
	defer func() {
		ch <- c.up
		ch <- c.totalScrapes
		ch <- c.jsonParseFailures
	}()

	clusterStateResp, err := c.fetchAndDecodeClusterState()
	if err != nil {
		c.up.Set(0)
		_ = level.Warn(c.logger).Log(
			"msg", "failed to fetch and decode cluster state",
			"err", err,
		)
		return
	}
	c.up.Set(1)
	cluster := clusterStateResp.ClusterName
	indexStateStatusMap := make(map[string]*IndexStateStatus)
	for index, indexDetail := range clusterStateResp.Metadata.Indices {
		dynamic := "true"
		for _, mappings := range indexDetail.Mappings {
			if mappings.Dynamic != nil {
				mapping := mappings.Dynamic.(string)
				if mapping != "" {
					dynamic = mapping
					break
				}
			}
		}
		value := 0.0
		indexStateStatusMap[index] = &IndexStateStatus{
			EsCluster:        cluster,
			Index:            index,
			State:            indexDetail.State,
			CreationDate:     indexDetail.Settings.Index.CreationDate,
			NumberOfReplicas: indexDetail.Settings.Index.NumberOfReplicas,
			NumberOfShards:   indexDetail.Settings.Index.NumberOfShards,
			Uuid:             indexDetail.Settings.Index.Uuid,
			Dynamic:          dynamic,
			Read:             FALSE,
			Write:            FALSE,
			MetadataRead:     FALSE,
			MetadataWrite:    FALSE,
			value:            value,
		}
	}

	//标记索引是否禁止读写
	if len(clusterStateResp.Blocks.Indices) > 0 {
		for index, blocksMap := range clusterStateResp.Blocks.Indices {
			indexStateStatus, ok := indexStateStatusMap[index]
			if !ok {
				continue
			}
			if len(blocksMap) <= 0 {
				continue
			}
			blocks := make([]string, 0)
			for _, blocksDetail := range blocksMap {
				blocks = append(blocks, blocksDetail.Levels...)
			}
			block := strings.Join(blocks, ",")
			if strings.Contains(block, READ) {
				indexStateStatus.Read = TRUE
			}
			if strings.Contains(block, WRITE) {
				indexStateStatus.Write = TRUE
			}
			if strings.Contains(block, METADATA_WRITE) {
				indexStateStatus.MetadataWrite = TRUE
			}
			if strings.Contains(block, METADATA_READ) {
				indexStateStatus.MetadataRead = TRUE
			}
			indexStateStatusMap[index] = indexStateStatus
		}
	}

	//标记索引的status
	if len(clusterStateResp.RoutingTable.Indices) > 0 {
		for index, shards := range clusterStateResp.RoutingTable.Indices {
			color := GREEN
			indexStateStatus, ok := indexStateStatusMap[index]
			if !ok {
				continue
			}
   	        stopFound:
			for _, shardsDetails := range shards.Shards {
				for _, shardsDetail := range shardsDetails {
					//如果索引已经被标识为red,则直接退出循坏
					if color == RED {
						break stopFound
					}
					if shardsDetail.State == "UNASSIGNED" {
						color = YELLOW
						if shardsDetail.Primary {
							color = RED
						}
					}
				}
			}
			indexStateStatus.Color = color
			if color == YELLOW {
				indexStateStatus.value = 1
			}
			if color == RED {
				indexStateStatus.value = 2
			}
			indexStateStatusMap[index] = indexStateStatus
		}
	}

	for _, indexStateStatus := range indexStateStatusMap {
		//"index", "color", "state", "dynamic", "creation_date", READ, WRITE, METADATA_READ, METADATA_WRITE,"es_cluster"
		for _, metric := range c.metrics {
			ch <- prometheus.MustNewConstMetric(
				metric.Desc,
				metric.Type,
				metric.Value(indexStateStatus.value),
				metric.Labels.values(indexStateStatus.Index, indexStateStatus.Color, indexStateStatus.State,
					indexStateStatus.Dynamic, indexStateStatus.CreationDate, indexStateStatus.NumberOfShards,
					indexStateStatus.NumberOfReplicas, indexStateStatus.Read, indexStateStatus.Write,
					indexStateStatus.MetadataRead, indexStateStatus.MetadataWrite, indexStateStatus.EsCluster)...,
			)
		}
	}

}
