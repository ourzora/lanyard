package migrate

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
)

// Migration describes a PostgreSQL database migration.
// At least one of SQL and Hash must be set.
// If both are set, Hash must be the SHA-256 of SQL.
type Migration struct {
	Name        string
	SQL         string // unavailable for applied migrations
	Hash        string // internally computed if unset
	RequiredSHA string
	OutsideTx   bool
}

func (m *Migration) hash() string {
	if m.SQL == "" {
		return m.Hash
	}
	h := sha256.Sum256([]byte(withoutWhiteSpace(m.SQL)))
	return hex.EncodeToString(h[:])
}

// withoutWhiteSpace removes all whitespace, interior and exterior so that
// migration specifications can be tweaked for readability without changing the
// semantics.
func withoutWhiteSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func (m Migration) String() string {
	return fmt.Sprintf("%s - %s", m.Name, m.hash()[:5])
}

// Error records an error
// and the migration that caused it.
// Index is the index of Mig in
// the given list of migrations,
// or the index where it would have
// been if it's not in the list.
type Error struct {
	Mig   Migration
	Index int
	Err   error
}

func (e *Error) Error() string {
	return e.Mig.String() + " at " + strconv.Itoa(e.Index) + ": " + e.Err.Error()
}

// Run runs all unapplied migrations in m.
func Run(ctx context.Context, db beginner, m []Migration) error {
	for {
		// keep going until there are no more to run (or an error)
		ran, err := run1(ctx, db, m)
		if !ran {
			return err
		}
	}
}

type beginner interface {
	execer
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}
type execer interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// run1 runs a single unapplied migration.
// It returns whether it ran successfully
// along with any error.
//
// Note that a return value of false, nil
// means that there were no unapplied migrations to run.
func run1(ctx context.Context, db beginner, migrations []Migration) (ran bool, err error) {
	// Begin a SQL transaction for all changes to the migrations
	// table. We use it to acquire an exclusive lock on the
	// table to ensure we're the only process migrating this
	// database.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false, xerrors.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Acquire a tx-level lock. 4 is an arbitrary bigint that should
	// uniquely identify this lock.
	_, err = tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(4)`)
	if err != nil {
		return false, xerrors.Errorf("advisory lock: %w", err)
	}

	// Create the migrations table if not yet created.
	const q = `
		CREATE SEQUENCE IF NOT EXISTS migration_seq
			START WITH 1
			INCREMENT BY 1
			NO MINVALUE
			NO MAXVALUE
			CACHE 1;
		CREATE TABLE IF NOT EXISTS migrations (
			filename text NOT NULL,
			hash text NOT NULL,
			applied_at timestamp with time zone DEFAULT now() NOT NULL,
			index int DEFAULT nextval('migration_seq') NOT NULL,
			PRIMARY KEY(filename)
		);
	`
	_, err = tx.ExecContext(ctx, q)
	if err != nil {
		return false, xerrors.Errorf("creating migration table: %w", err)
	}

	// Query migrations for the set of unapplied transactions.
	unapplied, err := FilterApplied(tx, migrations)
	if err != nil {
		return false, err
	}

	if len(unapplied) == 0 {
		return false, nil // all up to date!
	}

	m := unapplied[0]

	// Some migrations contain SQL that PostgreSQL does not support
	// within a transaction, so we flag those and run them outside
	// of the migration transaction.
	//
	// This means it is possible that the migration can be applied
	// successfully, but then we fail to insert a row in the
	// migrations table. We'll need to keep an eye on any error
	// messages from migratedb, to avoid attempting to apply
	// the same migration again.
	if m.OutsideTx {
		_, err = db.ExecContext(ctx, m.SQL)
	} else {
		_, err = tx.ExecContext(ctx, m.SQL)
	}
	if err != nil {
		return false, xerrors.Errorf("migration %s: %w", m.Name, err)
	}

	err = insertAppliedMigration(ctx, tx, m)
	if err != nil {
		return false, err
	}
	err = tx.Commit()
	if err == nil {
		log.Ctx(ctx).Info().Str("migration", m.Name).Msg("success")
	} else {
		log.Ctx(ctx).Error().Err(err).Str("migration", m.Name).Msg("failed")
	}
	return err == nil, err
}

// GetApplied returns the list of currently-applied migrations.
func GetApplied(ctx context.Context, db execer) ([]Migration, error) {
	const q1 = `
		SELECT count(*) FROM pg_tables
		WHERE schemaname=current_schema() AND tablename='migrations'
	`
	var n int
	err := db.QueryRowContext(ctx, q1).Scan(&n)
	if err != nil {
		return nil, xerrors.Errorf("checking for migrations table: %w", err)
	}
	if n == 0 {
		return nil, nil
	}

	const q2 = `
		SELECT filename, hash
		FROM migrations
		ORDER BY index
	`
	var a []Migration
	rows, err := db.QueryContext(ctx, q2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := Migration{}
		err := rows.Scan(&m.Name, &m.Hash)
		if err != nil {
			return nil, err
		}
		a = append(a, m)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return a, nil
}

// FilterApplied returns the slice of ms containing
// all migrations in ms that haven't yet been applied.
func FilterApplied(db execer, ms []Migration) ([]Migration, error) {
	applied, err := GetApplied(context.TODO(), db)
	if err != nil {
		return nil, err
	}
	for i, app := range applied {
		if i >= len(ms) {
			return nil, &Error{Mig: app, Index: i, Err: errors.New("applied but not requested")}
		}
		m := ms[i]
		if app.hash() != m.hash() {
			return nil, &Error{Mig: app, Index: i, Err: errors.New("hash mismatch")}
		}
	}
	return ms[len(applied):], nil
}

func insertAppliedMigration(ctx context.Context, db *sql.Tx, m Migration) error {
	const q = `
		INSERT INTO migrations (filename, hash, applied_at)
		VALUES($1, $2, NOW())
	`
	_, err := db.ExecContext(ctx, q, m.Name, m.hash())
	if err != nil {
		return xerrors.Errorf("recording applied migration: %w", err)
	}
	return nil
}

var validNameRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\.\d\.[a-z][a-z0-9_-]+\.sql$`)

// Validity returns an error if the provided list of migrations
// isn't valid. The required properties are:
//   names are well formed
//   names are in order
//   keys (YYYY-MM-DD.N) are not duplicated
func Validity(migrations []Migration) error {
	list := make([]Migration, len(migrations))
	copy(list, migrations)

	for i, m := range list {
		if !validNameRegex.MatchString(m.Name) {
			return fmt.Errorf("bad name: %s", m.Name)
		}
		if i > 0 && list[i-1].Name >= m.Name {
			return errors.New("out of order: " + m.Name)
		}
	}

	// Fail if we have more than one of any index
	// on the same day. YYYY-MM-DD.N
	a := make([]string, len(list))
	for i, m := range list {
		a[i] = m.Name[:12]
		if i > 0 && a[i-1] == a[i] {
			return fmt.Errorf("duplicate indexes %s %s", list[i-1].Name, m.Name)
		}
	}
	return nil
}
