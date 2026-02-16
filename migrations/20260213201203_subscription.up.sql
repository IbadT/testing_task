CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    service_name VARCHAR(100) NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    --     user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE CHECK (end_date >= start_date),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);