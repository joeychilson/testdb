-- migrate:up
CREATE TABLE "test" (
    "id" bigserial PRIMARY KEY,
    "uuid" uuid,
    "varchar" varchar,
    "text" text,
    "int" int,
    "bigint" bigint,
    "float" float,
    "double" double precision,
    "decimal" decimal,
    "boolean" boolean,
    "inet" inet,
    "macaddr" macaddr,
    "json" json,
    "jsonb" jsonb,
    "xml" xml,
    "date" date,
    "time" time,
    "timez" time with time zone,
    "timestamp" timestamp,
    "timestampz" timestamp with time zone

);

-- migrate:down
DROP TABLE "test";