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
	{
		Name: "2022-08-05.0.rename.sql",
		SQL: `
			ALTER TABLE merkle_trees
			RENAME COLUMN addresses TO unhashed_leaves;

			ALTER TABLE merkle_trees
			ADD COLUMN ltd text[] NOT NULL;

			ALTER TABLE merkle_trees
			ADD COLUMN packed boolean NOT NULL;


			ALTER TABLE merkle_proofs
			ADD COLUMN unhashed_leaf bytea NOT NULL;

			ALTER TABLE merkle_proofs
			ALTER COLUMN address
			DROP NOT NULL;
		`,
	},
}
