package collector

type clusterStateResponse struct {
	ClusterName           string          `json:"cluster_name"`
	ClusterUuid           string          `json:"cluster_uuid"`
	CompressedSizeInBytes int64           `json:"compressed_size_in_bytes"`
	Version               int64           `json:"version"`
	StateUuid             string          `json:"state_uuid"`
	MasterNode            string          `json:"master_node"`
	Blocks                blocks          `json:"blocks"`
	Nodes                 map[string]node `json:"nodes"`
	Metadata              metadata        `json:"metadata"`
	RoutingTable          routingTable    `json:"routing_table"`
	RoutingNodes          routingNodes    `json:"routing_nodes"`
}

//索引禁止读写信息
//集群级别禁止读写信息
type blocks struct {
	Global  map[string]blocksDetail            `json:"global"`
	Indices map[string]map[string]blocksDetail `json:"indices"` //最外层key为index
}

type blocksDetail struct {
	Description string   `json:"description"`
	Retryable   bool     `json:"retryable"`
	Levels      []string `json:"levels"`
}

type node struct {
	Name             string            `json:"name"`
	EphemeralId      string            `json:"ephemeral_id"`
	TransportAddress string            `json:"transport_address"`
	Attributes       map[string]string `json:"attributes"`
}

type metadata struct {
	ClusterUuid    string                 `json:"cluster_uuid"`
	Templates      interface{}            `json:"templates"`
	Indices        map[string]indexDetail `json:"indices"` //key为索引名称
	Repositories   interface{}            `json:"repositories"`
	IndexGraveyard interface{}            `json:"index-graveyard"`
}

//索引state
type indexDetail struct {
	State             string              `json:"state"`
	Settings          settings            `json:"settings"`
	Mappings          map[string]mappings `json:"mappings"` //key为索引名称
	Aliases           interface{}         `json:"aliases"`
	PrimaryTerms      interface{}         `json:"primary_terms"`
	InSyncAllocations map[string][]string `json:"in_sync_allocations"` //key为分片编号
}

type settings struct {
	Index indexSettings `json:"index"`
}

//分片所在节点
//索引创建时间
type indexSettings struct {
	CreationDate     string `json:"creation_date"`
	NumberOfShards   string `json:"number_of_shards"`
	NumberOfReplicas string `json:"number_of_replicas"`
	Uuid             string `json:"uuid"`
}

//是否开启动态mapping
type mappings struct {
	Dynamic    interface{}            `json:"dynamic"`
	Properties map[string]interface{} `json:"properties"`
}

type routingTable struct {
	Indices map[string]shards `json:"indices"` //key 为索引
}

type shards struct {
	Shards map[string][]shardsDetail `json:"shards"` //key为分片编号
}

type shardsDetail struct {
	State          string      `json:"state"`
	Primary        bool        `json:"primary"`
	Node           string      `json:"node"`
	RelocatingNode string      `json:"relocating_node"`
	Shard          int64       `json:"shard"`
	Index          string      `json:"index"`
	AllocationId   interface{} `json:"allocation_id"`
}

//未分配分片数目
type routingNodes struct {
	Unassigned []interface{}          `json:"unassigned"`
	Nodes      map[string]interface{} `json:"nodes"`
}
