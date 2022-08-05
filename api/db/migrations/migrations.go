package migrations

import "github.com/contextart/al/migrate"

var Migrations = []migrate.Migration{
	{
		Name: "2022-08-04.0.init.sql",
		SQL: `
		CREATE TABLE merkle_trees (
			root bytea,
			addresses bytea[] NOT NULL,
			PRIMARY KEY (root)
		);
		CREATE TABLE merkle_proofs (
			root bytea NOT NULL,
			address bytea NOT NULL,
			proof bytea[] NOT NULL
		);
		CREATE UNIQUE INDEX  ON merkle_proofs(root, address);
	`,
	},
}
