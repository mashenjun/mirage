package consts

import "encoding/json"

var (
	OpenIDKey = "OpenIDKey"
	EmpNoKey  = "EmpNoKey"
)

type Site struct {
	SiteID   int64
	SiteName string
}

type siteSlice []Site

var Sites = siteSlice{
	{
		1,
		"食堂1",
	},
	{
		2,
		"食堂2",
	},
}

func (ss siteSlice) GetOtherSites(siteID int64) siteSlice {
	if len(ss) <= 0 {
		return nil
	}
	ret := make([]Site, 0)
	for _, v := range ss {
		if v.SiteID != siteID {
			ret = append(ret, v)
		}
	}
	return ret
}

var EmptyJson, _ = json.Marshal(struct {}{})
