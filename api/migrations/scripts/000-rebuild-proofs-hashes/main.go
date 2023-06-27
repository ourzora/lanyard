package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"sync"

	"github.com/contextwtf/lanyard/merkle"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/sync/errgroup"
)

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "processor error: %s", err)
		debug.PrintStack()
		os.Exit(1)
	}
}

func hashProof(p [][]byte) []byte {
	return crypto.Keccak256(p...)
}

func migrateTree(
	ctx context.Context,
	tx pgx.Tx,
	leaves [][]byte,
) error {
	tree := merkle.New(leaves)

	var (
		proofHashes = [][]any{}
		eg          errgroup.Group
		pm          sync.Mutex
	)
	eg.SetLimit(runtime.NumCPU())

	for _, l := range leaves {
		l := l //avoid capture
		eg.Go(func() error {
			pf := tree.Proof(l)
			if !merkle.Valid(tree.Root(), pf, l) {
				return errors.New("invalid proof for tree")
			}
			proofHash := hashProof(pf)
			pm.Lock()
			proofHashes = append(proofHashes, []any{tree.Root(), proofHash})
			pm.Unlock()
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return err
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"proofs_hashes"},
		[]string{"root", "hash"},
		pgx.CopyFromRows(proofHashes),
	)

	return err
}

func main() {
	ctx := context.Background()
	const defaultPGURL = "postgres:///al"
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		dburl = defaultPGURL
	}
	dbc, err := pgxpool.ParseConfig(dburl)
	check(err)

	db, err := pgxpool.ConnectConfig(ctx, dbc)
	check(err)

	log.Println("fetching roots from db")
	const q = `
		SELECT unhashed_leaves
		FROM trees
		WHERE root not in (select root from proofs_hashes group by 1)
	`
	rows, err := db.Query(ctx, q)
	check(err)
	defer rows.Close()

	trees := [][][]byte{}

	for rows.Next() {
		var t [][]byte
		err := rows.Scan(&t)
		trees = append(trees, t)
		check(err)
	}

	log.Printf("migrating %d trees", len(trees))

	tx, err := db.Begin(ctx)
	check(err)
	defer tx.Rollback(ctx)

	var count int

	for _, tree := range trees {
		err = migrateTree(ctx, tx, tree)
		check(err)
		count++
		if count%1000 == 0 {
			log.Printf("migrated %d/%d trees", count, len(trees))
		}
	}

	log.Printf("committing %d trees", len(trees))
	err = tx.Commit(ctx)
	check(err)
	log.Printf("done")
}
