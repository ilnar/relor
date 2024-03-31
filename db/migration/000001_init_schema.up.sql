CREATE TABLE "workflows" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "current_node" varchar NOT NULL,
  "status" varchar NOT NULL,
  "graph" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "next_action_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transitions" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "workflow_id" uuid NOT NULL REFERENCES workflows(id),
  "from_node" varchar NOT NULL,
  "to_node" varchar NOT NULL,
  "label" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "previous" uuid REFERENCES transitions(id),
  "next" uuid REFERENCES transitions(id),

  -- Ensure that a workflow can only have one transition to a given node, i.e. linked list.
  UNIQUE (workflow_id, previous),
  UNIQUE (workflow_id, "next")
);