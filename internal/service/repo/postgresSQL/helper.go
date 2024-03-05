package postgresSQL

import (
	"fmt"
)

func getLimitAndOffset(limit, offset int) string {
	var limitQ, offsetQ string
	if limit > 0 {

		limitQ = fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset > 0 {
		offsetQ = fmt.Sprintf(" OFFSET %d", offset)
	}
	return limitQ + offsetQ
}
