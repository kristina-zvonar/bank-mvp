CREATE TABLE "countries" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "clients" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "active" boolean NOT NULL DEFAULT true,
  "country_id" bigint NOT NULL
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "balance" decimal NOT NULL,
  "currency" varchar NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "locked" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "client_id" bigint NOT NULL
);

CREATE TABLE "services" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "type" int NOT NULL
);

CREATE TABLE "cards" (
  "id" bigserial PRIMARY KEY,
  "number" varchar NOT NULL,
  "valid_through" date NOT NULL,
  "cvc" char(3) NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "account_id" bigint NOT NULL
);

CREATE TABLE "banks" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "country_id" bigint NOT NULL
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "amount" decimal NOT NULL,
  "source_account_id" bigint,
  "dest_account_id" bigint,
  "ext_source_account_id" varchar,
  "ext_dest_account_id" varchar,
  "category" int NOT NULL,
  "service_id" bigint
);

CREATE UNIQUE INDEX ON "countries" ("name");

CREATE UNIQUE INDEX ON "banks" ("name");

ALTER TABLE "clients" ADD FOREIGN KEY ("country_id") REFERENCES "countries" ("id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("client_id") REFERENCES "clients" ("id");

ALTER TABLE "cards" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "banks" ADD FOREIGN KEY ("country_id") REFERENCES "countries" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("source_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("dest_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id");
