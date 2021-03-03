package pg

import (
	"context"
	"fmt"

	"github.com/joerx/lolcatz-backend/db"
)

func loadSchema(db db.DB) error {

	ctx := context.Background()

	// Needs postgres >= 9.5!
	sql := []string{
		`CREATE SEQUENCE IF NOT EXISTS uploads_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1`,
		`CREATE TABLE IF NOT EXISTS "uploads" (
			"id" integer DEFAULT nextval('uploads_id_seq') NOT NULL,
			"username" character varying(100) NOT NULL,
			"filename" character varying(100) NOT NULL,
			"s3key" character varying(100) NOT NULL,
			"timestamp" time without time zone NOT NULL,
			CONSTRAINT "uploads_id" PRIMARY KEY ("id")
		) WITH (oids = false);`,
		`CREATE INDEX IF NOT EXISTS "uploads_username" ON "public"."uploads" USING btree ("username")`,
	}

	for _, stmt := range sql {
		if _, err := db.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("error creating db schema: %v", err)
		}
	}

	return nil
}
