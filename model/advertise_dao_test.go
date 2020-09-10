package model

import (
	"encoding/json"
	"testing"
)

func TestAdvertiseDao_Find(t *testing.T) {
	adv := AdvConfig{Candidates: []AdvCandidate{
		{
			AdCode: "debug",
			Location: "http://www.baidu.com",
			Action: 2,
		},
	}}
	b, err := json.Marshal(adv)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}
