-- name: InsertMerkleTree :exec
insert into merkle_trees (root, addresses)
values ($1, $2)
on conflict (root) do nothing;

-- name: GetAddressesForMerkleTree :one
select addresses from merkle_trees where root = $1;

-- name: InsertMerkleProof :exec
insert into merkle_proofs (root, address, proof)
values ($1, $2, $3)
on conflict (root, address) do nothing;

-- name: GetMerkleProof :one
select proof from merkle_proofs where root = $1 and address = $2;
