CREATE TABLE "public"."merkle_proofs" (
    "root" bytea NOT NULL,
    "address" bytea NOT NULL,
    "proof" bytea[] NOT NULL
);

CREATE UNIQUE INDEX merkle_proofs_unique_proof_for_root_and_address
  ON public.merkle_proofs
  USING btree (root, address);
