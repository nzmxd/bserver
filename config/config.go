package config

type Server struct {
	Zap        Zap             `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis      Redis           `mapstructure:"redis" json:"redis" yaml:"redis"`
	RedisList  []Redis         `mapstructure:"redis-list" json:"redis-list" yaml:"redis-list"`
	System     System          `mapstructure:"system" json:"system" yaml:"system"`
	Mysql      Mysql           `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	ClickHouse ClickHouse      `mapstructure:"clickhouse" json:"clickhouse" yaml:"clickhouse"`
	Sqlite     Sqlite          `mapstructure:"sqlite" json:"sqlite" yaml:"sqlite"`
	DBList     []SpecializedDB `mapstructure:"db-list" json:"db-list" yaml:"db-list"`
	Local      Local           `mapstructure:"local" json:"local" yaml:"local"`
	AliyunOSS  AliyunOSS       `mapstructure:"aliyun-oss" json:"aliyun-oss" yaml:"aliyun-oss"`
	TencentCOS TencentCOS      `mapstructure:"tencent-cos" json:"tencent-cos" yaml:"tencent-cos"`
	Minio      Minio           `mapstructure:"minio" json:"minio" yaml:"minio"`
	DiskList   []DiskList      `mapstructure:"disk-list" json:"disk-list" yaml:"disk-list"`
}
