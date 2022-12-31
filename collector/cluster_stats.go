package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	defaultClusterStatsLabels = []string{"es_cluster"}
)

type clusterStatsMetric struct {
	Type  prometheus.ValueType
	Desc  *prometheus.Desc
	Value func(clusterHealth clusterStatsResponse) float64
}

// ClusterStats type defines the collector struct
type ClusterStats struct {
	logger                          log.Logger
	client                          *http.Client
	url                             *url.URL
	up                              prometheus.Gauge
	totalScrapes, jsonParseFailures prometheus.Counter
	metrics                         []*clusterStatsMetric
}

// NewClusterStats returns a new Collector exposing NewClusterStats stats.
func NewClusterStats(logger log.Logger, client *http.Client, url *url.URL) *ClusterStats {
	subsystem := "cluster_stats"

	return &ClusterStats{
		logger: logger,
		client: client,
		url:    url,

		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, "up"),
			Help: "Was the last scrape of the ElasticSearch cluster stats endpoint successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, "total_scrapes"),
			Help: "Current total ElasticSearch cluster stats scrapes.",
		}),
		jsonParseFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, "json_parse_failures"),
			Help: "Number of errors while parsing JSON.",
		}),

		metrics: []*clusterStatsMetric{
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_count"),
					"indices count.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Count)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_docs_count"),
					"indices docs count.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Docs.Count)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_docs_deleted"),
					"indices docs deleted.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Docs.Deleted)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_store_size_in_bytes"),
					"indices store size_in_bytes.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Store.SizeInBytes)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_fielddata_memory"),
					"indices fielddata memory.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.FieldData.MemorySizeInBytes)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_fielddata_evictions"),
					"indices fielddata evictions.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.FieldData.Evictions)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_query_cache_memory"),
					"indices query_cache memory.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.QueryCache.MemorySizeInBytes)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_query_cache_cache_count"),
					"indices query_cache_cache count.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.QueryCache.CacheCount)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_query_cache_cache_size"),
					"indices query_cache cache_size.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.QueryCache.CacheSize)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_query_cache_evictions"),
					"indices query_cache evictions.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.QueryCache.Evictions)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_query_cache_total_count"),
					"indices query_cache total_count.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.QueryCache.TotalCount)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_query_cache_miss_count"),
					"indices query_cache miss_count.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.QueryCache.MissCount)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_query_cache_hit_count"),
					"indices query_cache hit_count.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.QueryCache.HitCount)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_segment_memory"),
					"indices segment memory.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Segments.MemorySizeInBytes)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_segment_count"),
					"indices segment count.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Segments.Count)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_segment_doc_value_memory"),
					"indices segment doc_value memory.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Segments.DocValuesMemoryInBytes)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_segment_fixed_bitset_memory"),
					"indices segment fixedBitSet memory.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Segments.FixedBitSetMemoryInBytes)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_segment_index_writer_memory"),
					"indices segment indexWriter memory.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Segments.IndexWriterMemoryInBytes)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_segment_norms_memory"),
					"indices segment norms memory.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Segments.NormsMemoryInBytes)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_segment_points_memory"),
					"indices segment points memory.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Segments.PointsMemoryInBytes)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_segment_term_vectors_memory"),
					"indices segment term_vectors memory.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Segments.TermVectorsMemoryInBytes)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "indices_segment_term_memory"),
					"indices segment term memory.",
					defaultClusterStatsLabels, nil,
				),
				Value: func(clusterStats clusterStatsResponse) float64 {
					return float64(clusterStats.Indices.Segments.TermsMemoryInBytes)
				},
			},
		},
	}
}

// Describe set Prometheus metrics descriptions.
func (c *ClusterStats) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.Desc
	}

	ch <- c.up.Desc()
	ch <- c.totalScrapes.Desc()
	ch <- c.jsonParseFailures.Desc()
}

func (c *ClusterStats) fetchAndDecodeClusterHealth() (clusterStatsResponse, error) {
	var chr clusterStatsResponse

	u := *c.url
	u.Path = path.Join(u.Path, "/_cluster/stats")
	res, err := c.client.Get(u.String())
	if err != nil {
		return chr, fmt.Errorf("failed to get cluster stats from %s://%s:%s%s: %s",
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
		return chr, fmt.Errorf("HTTP Request failed with code %d", res.StatusCode)
	}

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.jsonParseFailures.Inc()
		return chr, err
	}

	if err := json.Unmarshal(bts, &chr); err != nil {
		c.jsonParseFailures.Inc()
		return chr, err
	}

	return chr, nil
}

// Collect collects ClusterStats metrics.
func (c *ClusterStats) Collect(ch chan<- prometheus.Metric) {
	var err error
	c.totalScrapes.Inc()
	defer func() {
		ch <- c.up
		ch <- c.totalScrapes
		ch <- c.jsonParseFailures
	}()

	clusterStatsResp, err := c.fetchAndDecodeClusterHealth()
	if err != nil {
		c.up.Set(0)
		_ = level.Warn(c.logger).Log(
			"msg", "failed to fetch and decode cluster stats",
			"err", err,
		)
		return
	}
	c.up.Set(1)

	for _, metric := range c.metrics {
		ch <- prometheus.MustNewConstMetric(
			metric.Desc,
			metric.Type,
			metric.Value(clusterStatsResp),
			clusterStatsResp.ClusterName,
		)
	}
}
