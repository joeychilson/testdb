-- migrate:up
CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "first_name" varchar(255),
    "last_name" varchar(255),
    "email" varchar(255)
);

-- migrate:down
DROP TABLE "users";