package store

import (
	"net/http"
	"strconv"
)

type FeedPaginationQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"omitempty,oneof=asc desc"`
}

func (fpq FeedPaginationQuery) Parse(r *http.Request) (FeedPaginationQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")

	if limit != "" {
		l, err := strconv.Atoi(limit)

		if err != nil {
			return fpq, nil
		}

		fpq.Limit = l
	}

	offset := qs.Get("offset")

	if offset != "" {
		o, err := strconv.Atoi(offset)

		if err != nil {
			return fpq, nil
		}

		fpq.Offset = o
	}

	sort := qs.Get("sort")
	if sort != "" {
		fpq.Sort = sort
	}

	return fpq, nil
}
