-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- Insert test users for development and testing
-- These users represent different roles in the approval workflow:
-- 1. EMPLOYEE - creates limit requests
-- 2. TEAM_LEAD - first level approval
-- 3. RISK_OFFICER - second level approval for high amounts
-- 4. CFO - final approval for very high amounts

INSERT INTO
    users (
        id,
        external_id,
        email,
        name,
        role
    )
VALUES
    -- Primary test user for creating requests
    (
        'a005d32d-6190-477c-b23e-38c44eaaaae0',
        'test-user-1',
        'test@example.com',
        'Test User',
        'EMPLOYEE'
    ),

-- Approval workflow users
(
    'b123e567-e89b-12d3-a456-426614174000',
    'test-team-lead',
    'teamlead@example.com',
    'Team Lead User',
    'TEAM_LEAD'
),
(
    'c789f012-e89b-12d3-a456-426614174001',
    'test-risk-officer',
    'risk@example.com',
    'Risk Officer User',
    'RISK_OFFICER'
),
(
    'd456c789-e89b-12d3-a456-426614174002',
    'test-cfo',
    'cfo@example.com',
    'CFO User',
    'CFO'
) ON CONFLICT (id) DO NOTHING;
-- Avoid inserting duplicates if migration runs multiple times

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

-- Remove test users (in reverse order to avoid foreign key issues)
DELETE FROM users
WHERE
    external_id IN (
        'test-user-1',
        'test-team-lead',
        'test-risk-officer',
        'test-cfo'
    );