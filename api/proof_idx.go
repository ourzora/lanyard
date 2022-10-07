package api

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

func updateIdx(ctx context.Context, db *pgxpool.Pool) {
	t, err := db.Exec(ctx, `
			insert into trees_proofs
			(select root, proofs from trees
			where root not in (select root from trees_proofs))
			`)

	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to sync proof index")
	} else if t.RowsAffected() > 0 {
		log.Ctx(ctx).Info().Int64("rows", t.RowsAffected()).Msg("synced proof index")
	}

}

func (s *Server) SyncProofIdx(ctx context.Context) {
	// For large trees, it's expensive to write the index of
	// proofs to the database, so we do it in a background task.

	go func() {
		for ; ; time.Sleep(time.Second) {
			updateIdx(ctx, s.db)
		}
	}()
}
