CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT ('0001-01-01 00:00:00+00'),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "clients" ADD user_id bigint NOT NULL;
ALTER TABLE "clients" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "accounts" ADD CONSTRAINT "client_currency_key" UNIQUE (client_id, currency);