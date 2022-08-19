

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;


CREATE SEQUENCE public.migration_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


SET default_tablespace = '';

SET default_table_access_method = heap;


CREATE TABLE public.migrations (
    filename text NOT NULL,
    hash text NOT NULL,
    applied_at timestamp with time zone DEFAULT now() NOT NULL,
    index integer DEFAULT nextval('public.migration_seq'::regclass) NOT NULL
);



CREATE TABLE public.trees (
    root bytea NOT NULL,
    unhashed_leaves bytea[] NOT NULL,
    ltd text[],
    packed boolean,
    proofs jsonb
);



ALTER TABLE ONLY public.trees
    ADD CONSTRAINT merkle_trees_pkey PRIMARY KEY (root);



ALTER TABLE ONLY public.migrations
    ADD CONSTRAINT migrations_pkey PRIMARY KEY (filename);




