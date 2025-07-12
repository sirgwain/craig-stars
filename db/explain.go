package db

import "context"

type QueryPlan struct {
	ID      int    `json:"id,omitempty"`
	Parent  int    `json:"parent,omitempty"`
	Notused int    `json:"notused,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func (c *client) explainQuery(ctx context.Context, gameID int64) ([]QueryPlan, error) {

	query := `
EXPLAIN QUERY PLAN
SELECT
 p.id,
 d.id
FROM
    players p
    LEFT JOIN shipDesigns d ON p.gameId = d.gameId
    AND p.num = d.playerNum
WHERE
    p.gameId = ?;	
`
	results := []QueryPlan{}
	if err := c.readConn.Select(&results, query, gameID); err != nil {
		return nil, err
	}

	return results, nil
}
