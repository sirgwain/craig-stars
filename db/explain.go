package db

import (
	"context"
	"fmt"
)

type QueryPlan struct {
	ID      int    `json:"id,omitempty"`
	Parent  int    `json:"parent,omitempty"`
	Notused int    `json:"notused,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func (c *client) explainQuery(ctx context.Context) ([]QueryPlan, error) {

	rows, err := c.readConn.QueryContext(ctx, `
	EXPLAIN QUERY PLAN
	SELECT
		p.id,
		d.id
	FROM players p
	LEFT JOIN shipDesigns d ON p.gameId = d.gameId AND p.num = d.playerNum
	WHERE p.gameId = ?;
`, 1) // or any gameId
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []QueryPlan

	for rows.Next() {
		var item QueryPlan

		if err := rows.Scan(&item.ID, &item.Parent, &item.Notused, &item.Detail); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

	}

	return items, nil

}
