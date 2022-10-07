package api

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

func (s *Server) SyncProofIdx(ctx context.Context) {
	// It's expensive to write the index of proofs to the database, so we do it in
	// a background task.

	var lock sync.Mutex

	go func() {
		for ; ; time.Sleep(time.Second) {
			unlocked := lock.TryLock()
			if !unlocked {
				continue
			}
			t, err := s.db.Exec(ctx, `
			insert into trees_proofs
			(select root, proofs from trees
			where root not in (select root from trees_proofs))
			`)
			lock.Unlock()

			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("failed to sync proof index")
			} else if t.RowsAffected() > 0 {
				log.Ctx(ctx).Info().Int64("rows", t.RowsAffected()).Msg("synced proof index")
			}
		}
	}()
}
