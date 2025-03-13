package dynamo

type Config struct {
	TableName string `json:"table_name"`
	Region    string `json:"region"`
	Endpoint  string `json:"endpoint"`
}
