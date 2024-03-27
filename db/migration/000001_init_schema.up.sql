CREATE TABLE "workflows" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "current_node" varchar NOT NULL,
  "status" varchar NOT NULL,
  "graph" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "next_action_at" timestamptz NOT NULL DEFAULT (now())
);