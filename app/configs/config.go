package configs

type PostgresConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	DataBaseName string `yaml:"dataBaseName"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
}

type NatsConfig struct {
	clusterID string `yaml:"clusterID"`
	clientID  string `yaml:"clientID"`
	natsURL   string `yaml:"natsURL"`
}
