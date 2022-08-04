CREATE TABLE merkle_trees (
    "root" bytea,
    "addresses" bytea[] NOT NULL,
    PRIMARY KEY ("root")
);
