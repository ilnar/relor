CREATE TABLE "workflows" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "current_node" varchar NOT NULL,
  "status" varchar NOT NULL,
  "graph" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "next_action_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "workflow_events" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "workflow_id" uuid NOT NULL REFERENCES workflows(id),
  "from_node" varchar NOT NULL,
  "to_node" varchar NOT NULL,
  "label" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);