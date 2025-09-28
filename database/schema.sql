CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    "trigger" JSONB NOT NULL,
    action JSONB NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    next_run TIMESTAMPTZ

);


CREATE TABLE IF NOT EXISTS task_results (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
     run_at TIMESTAMPTZ NOT NULL,
     status_code INT NOT NULL,
     success BOOLEAN NOT NULL,
     response_headers JSONB,
     response_body JSONB,
     error_message TEXT,
     duration_ms INT NOT NULL,
     created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
