package config

import "fmt"

type ClickHouse struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (m *ClickHouse) Dsn() string {
	return fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s?max_execution_time=60", m.Username, m.Password, m.Path, m.Port, m.Dbname)
}
