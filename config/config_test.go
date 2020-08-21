package config

import (
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/qingyuanz/intelliecard/consts"
)

func TestInitOption(t *testing.T) {

}

var bs = `
---
server:
  addr: 0.0.0.0:9999
  request_timeout: 3 # second
  response_timeout: 5 # second, for pprof change it to 120

mysql:
  driver: mysql
  dsn: root:root@tcp(localhost:3306)/intelliecard?charset=utf8&parseTime=True&loc=Local
  max_open_conns: 20
  max_idle_conns: 10
  max_life_conns: 60
  trace_include_not_found: false

sql_server:
  driver: mysql
  dsn: root:root@tcp(localhost:3306)/intelliecard?charset=utf8&parseTime=True&loc=Local
  max_open_conns: 20
  max_idle_conns: 10
  max_life_conns: 60
  trace_include_not_found: false

deadline:
  breakfast: "07:30:00"
  lunch: "11:30:00"
  dinner: "17:30:00"
  supper: "19:30:00"

recent_days: 7

sites:
  - site_id: 1
    site_name: "食堂1"
  - site_id: 2
    site_name: "食堂2"
  - site_id: 3
    site_name: "食堂3"

log:
  level: "DEBUG"
`

func TestAppOptionsStruct(t *testing.T) {
	b := []byte(bs)
	option := Options
	if err := yaml.Unmarshal(b, &option); err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", option)
	FillSites(option)
	t.Logf("%+v", consts.Sites)
}
