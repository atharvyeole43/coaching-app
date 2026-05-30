CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE coach_status AS ENUM ('active', 'inactive');
CREATE TYPE transaction_status AS ENUM ('pending', 'completed', 'failed', 'refunded');

CREATE TABLE coaches (
    id          SERIAL PRIMARY KEY,
    uuid        UUID NOT NULL DEFAULT gen_random_uuid(),
    full_name   VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    status      coach_status NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL DEFAULT NULL,
    CONSTRAINT uq_coaches_email UNIQUE (email)
);

CREATE TABLE income_transactions (
    id                SERIAL PRIMARY KEY,
    uuid              UUID NOT NULL DEFAULT gen_random_uuid(),
    coach_id          INT NOT NULL,
    amount            NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    status            transaction_status NOT NULL DEFAULT 'completed',
    transaction_date  DATE NOT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at        TIMESTAMP NULL DEFAULT NULL,
    CONSTRAINT fk_income_coach FOREIGN KEY (coach_id) REFERENCES coaches(id)
);

CREATE INDEX idx_income_coach_date ON income_transactions (coach_id, transaction_date);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_coaches_updated_at
    BEFORE UPDATE ON coaches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_income_updated_at
    BEFORE UPDATE ON income_transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

INSERT INTO coaches (uuid, full_name, email, status) VALUES
(gen_random_uuid(), 'Rahul Sharma', 'rahul.sharma@gmail.com', 'active'),
(gen_random_uuid(), 'Priya Mehta',  'priya.mehta@gmail.com',  'active'),
(gen_random_uuid(), 'Amit Verma',   'amit.verma@gmail.com',   'active'),
(gen_random_uuid(), 'Sneha Patil',  'sneha.patil@gmail.com',  'inactive'),
(gen_random_uuid(), 'Vikas Joshi',  'vikas.joshi@gmail.com',  'active');

INSERT INTO income_transactions (uuid, coach_id, amount, status, transaction_date) VALUES
(gen_random_uuid(), 1, 80000,  'completed', CURRENT_DATE - INTERVAL '6 days'),
(gen_random_uuid(), 1, 120000, 'completed', CURRENT_DATE - INTERVAL '5 days'),
(gen_random_uuid(), 1, 95000,  'completed', CURRENT_DATE - INTERVAL '4 days'),
(gen_random_uuid(), 1, 150000, 'completed', CURRENT_DATE - INTERVAL '2 days'),
(gen_random_uuid(), 1, 200000, 'completed', CURRENT_DATE - INTERVAL '1 day'),
(gen_random_uuid(), 1, 175000, 'completed', CURRENT_DATE),
(gen_random_uuid(), 1, 60000,  'completed', CURRENT_DATE - INTERVAL '13 days'),
(gen_random_uuid(), 1, 90000,  'completed', CURRENT_DATE - INTERVAL '12 days'),
(gen_random_uuid(), 1, 75000,  'completed', CURRENT_DATE - INTERVAL '11 days'),
(gen_random_uuid(), 1, 110000, 'completed', CURRENT_DATE - INTERVAL '10 days'),
(gen_random_uuid(), 1, 85000,  'completed', CURRENT_DATE - INTERVAL '9 days'),
(gen_random_uuid(), 1, 130000, 'completed', CURRENT_DATE - INTERVAL '8 days'),
(gen_random_uuid(), 1, 100000, 'completed', CURRENT_DATE - INTERVAL '7 days');
