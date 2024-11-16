ALTER TABLE "users"
    ADD COLUMN "is_email_verified" bool NOT NULL DEFAULT false;

CREATE TABLE "verify_emails"
(
    "id"          bigserial PRIMARY KEY,
    "username"    text        NOT NULL,
    "email"       text        NOT NULL,
    "secret_code" text        NOT NULL,
    "is_used"     bool        NOT NULL DEFAULT false,
    "expires_at"  timestamptz NOT NULL DEFAULT (now() + interval '15 minutes'),
    "created_at"  timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "verify_emails"
    ADD FOREIGN KEY ("username") REFERENCES "users" ("username");