package collector

type clusterStatsResponse struct {
	ClusterName string                      `json:"cluster_name"`
	ClusterUuid string                      `json:"cluster_uuid"`
	Status      string                      `json:"status"`
	Timestamp   int64                       `json:"timestamp"`
	Nodes       clusterStatsNodesResponse   `json:"nodes"`
	Indices     clusterStatsIndicesResponse `json:"indices"`
}

type clusterStatsNodesResponse struct {
	Count          clusterStatsNodesCountResponse            `json:"count"`
	Versions       []string                                  `json:"versions"`
	Os             clusterStatsNodesOsResponse               `json:"os"`
	Process        clusterStatsNodesProcessResponse          `json:"process"`
	Jvm            clusterStatsNodesJvmResponse              `json:"jvm"`
	Fs             clusterStatsNodesFsResponse               `json:"fs"`
	Plugins        []clusterStatsNodesPluginResponse         `json:"plugins"`
	NetworkTypes   clusterStatsNodesNetworkTypesResponse     `json:"network_types"`
	DiscoveryTypes clusterStatsNodesDiscoveryTypesResponse   `json:"discovery_types"`
	PackagingTypes []clusterStatsNodesPackagingTypesResponse `json:"packaging_types"`
	Ingest         clusterStatsNodesIngestResponse           `json:"ingest,omitempty"`
}

type clusterStatsNodesIngestResponse struct {
	NumberOfPipelines int64                                                    `json:"number_of_pipelines,omitempty"`
	ProcessorStats    map[string]clusterStatsNodesIngestProcessorStatsResponse `json:"processor_stats,omitempty"`
}

type clusterStatsNodesIngestProcessorStatsResponse struct {
	Count        int64 `json:"count"`
	Failed       int64 `json:"failed"`
	Current      int64 `json:"current"`
	TimeInMillis int64 `json:"time_in_millis"`
}

type clusterStatsNodesPackagingTypesResponse struct {
	Flavor string `json:"flavor"`
	Type   string `json:"type"`
	Count  int64  `json:"count"`
}

type clusterStatsNodesDiscoveryTypesResponse struct {
	Zen int64 `json:"zen"`
}

type clusterStatsNodesNetworkTypesResponse struct {
	TransportTypes clusterStatsNodesNetworkTypesTypeResponse `json:"transport_types"`
	HttpTypes      clusterStatsNodesNetworkTypesTypeResponse `json:"http_types"`
}

type clusterStatsNodesNetworkTypesTypeResponse struct {
	Netty4 int64 `json:"netty4"`
}

type clusterStatsNodesPluginResponse struct {
	Name                 string   `json:"name"`
	Version              string   `json:"version"`
	ElasticsearchVersion string   `json:"elasticsearch_version"`
	JavaVersion          string   `json:"java_version"`
	Description          string   `json:"description"`
	ExtendedPlugins      []string `json:"extended_plugins"`
	Classname            string   `json:"classname"`
	HasNativeController  bool     `json:"has_native_controller"`
}

type clusterStatsNodesFsResponse struct {
	TotalInBytes     int64 `json:"total_in_bytes"`
	FreeInBytes      int64 `json:"free_in_bytes"`
	AvailableInBytes int64 `json:"available_in_bytes"`
}

type clusterStatsNodesJvmResponse struct {
	MaxUptimeInMillis int64                                 `json:"max_uptime_in_millis"`
	Versions          []clusterStatsNodesJvmVersionResponse `json:"versions"`
	Mem               clusterStatsNodesJvmMemResponse       `json:"mem"`
	Threads           int64                                 `json:"threads"`
}

type clusterStatsNodesJvmVersionResponse struct {
	Version         string `json:"version"`
	VmName          string `json:"vm_name"`
	VmVersion       string `json:"vm_version"`
	VmVendor        string `json:"vm_vendor"`
	BundledJdk      bool   `json:"bundled_jdk"`
	UsingBundledJdk bool   `json:"using_bundled_jdk"`
	Count           int64  `json:"count"`
}

type clusterStatsNodesJvmMemResponse struct {
	HeapUsedInBytes int64 `json:"heap_used_in_bytes"`
	HeapMaxInBytes  int64 `json:"heap_max_in_bytes"`
}

type clusterStatsNodesProcessResponse struct {
	Cpu                 clusterStatsNodesCpuResponse                 `json:"cpu"`
	OpenFileDescriptors clusterStatsNodesOpenFileDescriptorsResponse `json:"open_file_descriptors"`
}

type clusterStatsNodesCpuResponse struct {
	Percent float64 `json:"percent"`
}

type clusterStatsNodesOpenFileDescriptorsResponse struct {
	Min int64   `json:"min"`
	Max int64   `json:"max"`
	Avg float64 `json:"avg"`
}

type clusterStatsNodesOsResponse struct {
	AvailableProcessors int64                                    `json:"available_processors"`
	AllocatedProcessors int64                                    `json:"allocated_processors"`
	Names               []clusterStatsNodesOsNamesResponse       `json:"names"`
	PrettyNames         []clusterStatsNodesOsPrettyNamesResponse `json:"pretty_names"`
	Mem                 clusterStatsNodesOsMemResponse           `json:"mem"`
}

type clusterStatsNodesOsMemResponse struct {
	TotalInBytes int64 `json:"total_in_bytes"`
	FreeInBytes  int64 `json:"free_in_bytes"`
	UsedInBytes  int64 `json:"used_in_bytes"`
	FreePercent  int64 `json:"free_percent"`
	UsedPercent  int64 `json:"used_percent"`
}

type clusterStatsNodesOsPrettyNamesResponse struct {
	PrettyName string `json:"pretty_name"`
	Count      int64  `json:"count"`
}

type clusterStatsNodesOsNamesResponse struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type clusterStatsNodesCountResponse struct {
	Total               int64 `json:"total"`
	CoordinatingOnly    int64 `json:"coordinating_only"`
	Data                int64 `json:"data"`
	Ingest              int64 `json:"ingest"`
	Master              int64 `json:"master"`
	RemoteClusterClient int64 `json:"remote_cluster_client"`
}

type clusterStatsIndicesResponse struct {
	Count      int64                                 `json:"count"`
	Shards     clusterStatsIndicesShardsResponse     `json:"shards"`
	Docs       clusterStatsIndicesDocsResponse       `json:"docs"`
	Store      clusterStatsIndicesStoreResponse      `json:"store"`
	FieldData  clusterStatsIndicesFieldDataResponse  `json:"fielddata"`
	QueryCache clusterStatsIndicesQueryCacheResponse `json:"query_cache"`
	Completion clusterStatsIndicesCompletionResponse `json:"completion"`
	Segments   clusterStatsIndicesSegmentsResponse   `json:"segments"`
	Mappings   clusterStatsIndicesMappingsResponse   `json:"mappings"`
	Analysis   clusterStatsIndicesAnalysisResponse   `json:"analysis"`
}

type clusterStatsIndicesAnalysisResponse struct {
	CharFilterTypes    []clusterStatsIndicesFieldTypesResponse `json:"char_filter_types"`
	TokenizerTypes     []clusterStatsIndicesFieldTypesResponse `json:"tokenizer_types"`
	FilterTypes        []clusterStatsIndicesFieldTypesResponse `json:"filter_types"`
	AnalyzerTypes      []clusterStatsIndicesFieldTypesResponse `json:"analyzer_types"`
	BuiltInCharFilters []clusterStatsIndicesFieldTypesResponse `json:"built_in_char_filters"`
	BuiltInTokenizers  []clusterStatsIndicesFieldTypesResponse `json:"built_in_tokenizers"`
	BuiltInFilters     []clusterStatsIndicesFieldTypesResponse `json:"built_in_filters"`
	BuiltInAnalyzers   []clusterStatsIndicesFieldTypesResponse `json:"built_in_analyzers"`
}

type clusterStatsIndicesMappingsResponse struct {
	FieldTypes []clusterStatsIndicesFieldTypesResponse `json:"field_types"`
}

type clusterStatsIndicesFieldTypesResponse struct {
	Name       string `json:"name"`
	Count      int64  `json:"count"`
	IndexCount int64  `json:"index_count"`
}

type clusterStatsIndicesSegmentsResponse struct {
	Count                     int64                                                   `json:"count"`
	MemorySizeInBytes         int64                                                   `json:"memory_size_in_bytes"`
	TermsMemoryInBytes        int64                                                   `json:"terms_memory_in_bytes"`
	StoredFieldsMemoryInBytes int64                                                   `json:"stored_fields_memory_in_bytes"`
	TermVectorsMemoryInBytes  int64                                                   `json:"term_vectors_memory_in_bytes"`
	NormsMemoryInBytes        int64                                                   `json:"norms_memory_in_bytes"`
	PointsMemoryInBytes       int64                                                   `json:"points_memory_in_bytes"`
	DocValuesMemoryInBytes    int64                                                   `json:"doc_values_memory_in_bytes"`
	IndexWriterMemoryInBytes  int64                                                   `json:"index_writer_memory_in_bytes"`
	VersionMapMemoryInBytes   int64                                                   `json:"version_map_memory_in_bytes"`
	FixedBitSetMemoryInBytes  int64                                                   `json:"fixed_bit_set_memory_in_bytes"`
	MaxUnsafeAutoIdTimestamp  int64                                                   `json:"max_unsafe_auto_id_timestamp"`
	FileSizes                 map[string]clusterStatsIndicesSegmentsFileSizesResponse `json:"file_sizes,omitempty"`
}

type clusterStatsIndicesSegmentsFileSizesResponse struct {
	SizeInBytes int64 `json:"size_in_bytes,omitempty"`
	Description int64 `json:"description,omitempty"`
}

type clusterStatsIndicesCompletionResponse struct {
	SizeInBytes int64 `json:"size_in_bytes"`
}

type clusterStatsIndicesQueryCacheResponse struct {
	MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
	TotalCount        int64 `json:"total_count"`
	HitCount          int64 `json:"hit_count"`
	MissCount         int64 `json:"miss_count"`
	CacheSize         int64 `json:"cache_size"`
	CacheCount        int64 `json:"cache_count"`
	Evictions         int64 `json:"evictions"`
}

type clusterStatsIndicesFieldDataResponse struct {
	MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
	Evictions         int64 `json:"evictions"`
}

type clusterStatsIndicesStoreResponse struct {
	SizeInBytes     int64 `json:"size_in_bytes"`
	ReservedInBytes int64 `json:"reserved_in_bytes"`
}

type clusterStatsIndicesDocsResponse struct {
	Count   int64 `json:"count"`
	Deleted int64 `json:"deleted"`
}

type clusterStatsIndicesShardsResponse struct {
	Total       int64                                  `json:"total"`
	Primaries   int64                                  `json:"primaries"`
	Replication float64                                `json:"replication"`
	Index       clusterStatsIndicesShardsIndexResponse `json:"index"`
}

type clusterStatsIndicesShardsIndexResponse struct {
	Shards      clusterStatsIndicesShardsIndexShardResponse `json:"shards"`
	Primaries   clusterStatsIndicesShardsIndexShardResponse `json:"primaries"`
	Replication clusterStatsIndicesShardsIndexShardResponse `json:"replication"`
}

type clusterStatsIndicesShardsIndexShardResponse struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
	Avg float64 `json:"avg"`
}
