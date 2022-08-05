package migrate

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"blake.io/pqx/pqxtest"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

func TestMain(m *testing.M) {
	pqxtest.TestMain(m)
}

func TestMigrateConcurrent(t *testing.T) {
	db := pqxtest.CreateDB(t, ``)

	migrations := []Migration{
		{
			Name: "2000-01-01.0.migrate.select1.sql",
			SQL:  "CREATE TABLE foo (id text NOT NULL PRIMARY KEY);",
		},
	}

	var group errgroup.Group
	for i := 0; i < 5; i++ {
		group.Go(func() error {
			return Run(context.Background(), db, migrations)
		})
	}
	err := group.Wait()
	if err != nil {
		t.Error(err)
	}

	// Check that there are no remaining, unapplied migrations.
	unapplied, err := FilterApplied(db, migrations)
	if err != nil {
		t.Fatal(err)
	}
	if len(unapplied) > 0 {
		t.Errorf("len(FilterApplied(migrations)) = %d, want 0", len(unapplied))
	}
}

func TestFilterApplied(t *testing.T) {
	const migrationTable = `
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
	const oneMigration = `
		INSERT INTO migrations (filename, hash, applied_at)
		VALUES ('x', 'ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb', '2016-02-09T23:21:55 US/Pacific');
	`

	cases := []struct {
		initSQL string
		migs    []Migration
		want    []Migration
	}{
		{
			initSQL: `SELECT 1;`,
			migs:    nil,
			want:    nil,
		},
		{
			initSQL: migrationTable,
			migs:    nil,
			want:    nil,
		},
		{
			initSQL: `SELECT 1;`,
			migs:    []Migration{{Name: "x", SQL: "a"}},
			want:    []Migration{{Name: "x", SQL: "a"}},
		},
		{
			initSQL: migrationTable + oneMigration,
			migs:    []Migration{{Name: "x", SQL: "a"}},
			want:    []Migration{},
		},
		{
			initSQL: `SELECT 1;`,
			migs: []Migration{
				{Name: "x", SQL: "a"},
				{Name: "y", SQL: "b"},
			},
			want: []Migration{
				{Name: "x", SQL: "a"},
				{Name: "y", SQL: "b"},
			},
		},
		{
			initSQL: migrationTable + oneMigration,
			migs: []Migration{
				{Name: "x", SQL: "a"},
				{Name: "y", SQL: "b"},
			},
			want: []Migration{{Name: "y", SQL: "b"}},
		},
		{
			initSQL: migrationTable + oneMigration,
			migs:    []Migration{{Name: "x1", SQL: "a"}},
			want:    []Migration{},
		},
	}

	for _, test := range cases {
		db := pqxtest.CreateDB(t, ``)
		_, err := db.ExecContext(context.TODO(), test.initSQL)
		if err != nil {
			t.Error(err)
			continue
		}
		got, err := FilterApplied(db, test.migs)
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("FilterApplied() = %#v, want %#v", got, test.want)
		}
	}
}

func TestFilterAppliedError(t *testing.T) {
	const migrationTable = `
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
	const oneMigration = `
		INSERT INTO migrations (filename, hash, applied_at)
		VALUES ('x', 'ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb', '2016-02-09T23:21:55 US/Pacific');
	`

	cases := []struct {
		initSQL string
		migs    []Migration
		want    *Error
	}{
		{
			initSQL: migrationTable + oneMigration,
			migs:    nil,
			want: &Error{
				Mig: Migration{
					Name: "x",
					Hash: "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
				},
				Index: 0,
				Err:   errors.New("applied but not requested"),
			},
		},
		{
			initSQL: migrationTable + oneMigration,
			migs:    []Migration{{Name: "x", SQL: "a1"}},
			want: &Error{
				Mig: Migration{
					Name: "x",
					Hash: "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
				},
				Index: 0,
				Err:   errors.New("hash mismatch"),
			},
		},
	}

	for _, test := range cases {
		db := pqxtest.CreateDB(t, ``)
		_, err := db.ExecContext(context.TODO(), test.initSQL)
		if err != nil {
			t.Error(err)
			continue
		}
		_, got := FilterApplied(db, test.migs)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("FilterApplied(db, %v) err = %v, want %v", test.migs, got, test.want)
			t.Logf("got:  %#v", got)
			t.Logf("want: %#v", test.want)
		}
	}
}
