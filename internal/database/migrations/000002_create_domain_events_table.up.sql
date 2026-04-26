CREATE TABLE IF NOT EXISTS domain_events (
    id           UUID PRIMARY KEY,
    event_type   VARCHAR(255) NOT NULL,
    aggregate_id VARCHAR(255) NOT NULL,
    payload      JSONB NOT NULL,
    occurred_at  TIMESTAMPTZ NOT NULL,
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON domain_events (published_at) WHERE published_at IS NULL;
