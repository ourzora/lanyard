-- name: InsertTree :exec
insert into merkle_trees (root, unhashed_leaves, ltd, packed)
values ($1, $2, $3, $4)
on conflict (root) do nothing;

-- name: SelectTree :one
select unhashed_leaves, ltd, packed
from merkle_trees
where root = $1;

-- name: InsertProof :exec
insert into merkle_proofs (root, unhashed_leaf, address, proof)
values ($1, $2, $3, $4)
on conflict (root, unhashed_leaf) do nothing;

-- name: SelectProofByUnhashedLeaf :one
select proof
from merkle_proofs
where root = $1
and unhashed_leaf = $2;

-- name: SelectProofByAddress :many
select proof, unhashed_leaf
from merkle_proofs
where root = $1
and address = $2;

-- name: SelectTreeExists :one
select exists(select 1 from merkle_trees where root = $1);
