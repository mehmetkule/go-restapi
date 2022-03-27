CREATE TABLE document(
    id uuid DEFAULT uuid_generate_v4 (),
    parent_id character varying(100) NOT NULL,
    name character varying(150),
    data BYTEA NOT NULL,
    created date,
    CONSTRAINT document_pkey PRIMARY KEY (id)
) WITH(OIDS = FALSE);

CREATE TABLE users(
    id uuid DEFAULT uuid_generate_v4 (),
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    is_2fa bool,
    token text,
    created date,
    CONSTRAINT users_pkey PRIMARY KEY (id)
) WITH (OIDS = FALSE);