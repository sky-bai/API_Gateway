package util

import "math"

type PageList struct {
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	Count     int         `json:"count"`
	TotalPage int         `json:"total_page"`
	Data      interface{} `json:"data"`
}

func CutPage(total int, page int, limit int, data interface{}) (pager PageList) {
	pager.Count = total
	pager.TotalPage = int(math.Ceil(float64(total) / float64(limit)))
	pager.Page = page
	pager.Limit = limit
	pager.Data = data
	return
}
