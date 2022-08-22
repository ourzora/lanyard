package migrations

import "github.com/contextwtf/lanyard/migrate"

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
			CREATE UNIQUE INDEX ON merkle_proofs(root, address);
		`,
	},
	{
		Name: "2022-08-05.0.rename.sql",
		SQL: `
			ALTER TABLE merkle_trees
			RENAME COLUMN addresses TO unhashed_leaves;

			ALTER TABLE merkle_trees
			ADD COLUMN ltd text[];

			ALTER TABLE merkle_trees
			ADD COLUMN packed boolean;

			ALTER TABLE merkle_proofs
			ADD COLUMN unhashed_leaf bytea NOT NULL;

			ALTER TABLE merkle_proofs
			ALTER COLUMN address
			DROP NOT NULL;

			DROP INDEX merkle_proofs_root_address_idx;
			CREATE UNIQUE INDEX on merkle_proofs(root, unhashed_leaf);
		`,
	},
	{
		Name: "2022-08-17.0.proofs.sql",
		SQL: `
			ALTER TABLE merkle_trees
			ADD COLUMN proofs jsonb;
		`,
	},
	{
		Name: "2022-08-18.0.drop-proofs.sql",
		SQL: `
			DROP TABLE merkle_proofs;
		`,
	},
	{
		Name: "2022-08-18.1.rename-trees.sql",
		SQL: `
			ALTER TABLE merkle_trees RENAME TO trees;
		`,
	},
	{
		Name: "2022-08-22.0.add-proof-idx.sql",
		SQL: `
		    CREATE INDEX on trees USING gin(proofs jsonb_path_ops);
		`,
	},
	{
		Name: "2022-08-22.0.add-inserted-at.sql",
		SQL: `
		ALTER TABLE trees 
		ADD COLUMN "inserted_at" timestamptz NOT NULL DEFAULT now();
		`,
	},
}
