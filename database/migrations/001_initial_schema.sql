-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    external_id text UNIQUE NOT NULL, -- User ID from external authentication system (e.g., Keycloak sub)
    email text UNIQUE NOT NULL,
    name text NOT NULL,
    role text NOT NULL, -- User role (employee, approver, admin) - as per R-8
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE limit_requests (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL REFERENCES users (id), -- Request initiator (R-1)
    amount decimal NOT NULL, -- Request amount (R-1)
    currency text NOT NULL, -- Request currency (R-1)
    justification text NOT NULL, -- Justification (R-1)
    desired_date date NOT NULL, -- Desired start date (R-1)
    status text NOT NULL, -- Request status (e.g., PENDING_TEAM_LEAD, PENDING_RISK_OFFICER, PENDING_CFO, APPROVED, REJECTED, COMPLETED) (R-7)
    current_approver_id uuid REFERENCES users (id), -- Current approver (R-2, R-4)
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE approval_steps (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    request_id uuid NOT NULL REFERENCES limit_requests (id), -- Link to the request (R-5)
    approver_id uuid NOT NULL REFERENCES users (id), -- User who performed the action (R-5)
    action text NOT NULL, -- Action (e.g., APPROVED, REJECTED) (R-5)
    comment text, -- Comment (R-5)
    assigned_approver_role text NOT NULL, -- Role assigned for this approval step (e.g., TEAM_LEAD, RISK_OFFICER, CFO)
    "timestamp" timestamp NOT NULL DEFAULT now() -- Action date and time (R-5)
);

CREATE TABLE audit_log (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    request_id uuid REFERENCES limit_requests (id), -- Link to the request (R-10)
    action text NOT NULL, -- Action performed (R-10)
    actor_id uuid NOT NULL REFERENCES users (id), -- User who performed the action (R-10)
    "timestamp" timestamp NOT NULL DEFAULT now(), -- Action time (R-10)
    payload_json jsonb -- Action details (R-10)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS audit_log;

DROP TABLE IF EXISTS approval_steps;

DROP TABLE IF EXISTS limit_requests;

DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "pgcrypto";