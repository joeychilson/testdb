-- migrate:up
CREATE TABLE people (
    "id" SERIAL PRIMARY KEY,
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "full_name" VARCHAR(255) NOT NULL,
    "age" INTEGER NOT NULL,
    "salary" DECIMAL(10,2),
    "start_date" DATE NOT NULL,
    "phone" JSON NOT NULL,
    "languages" TEXT[] NOT NULL
);

-- migrate:down
DROP TABLE people;
