CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price INT NOT NULL CHECK (price >= 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);

CREATE INDEX idx_user_id ON subscriptions(user_id);
CREATE INDEX idx_service_name ON subscriptions(service_name);
CREATE INDEX idx_start_date ON subscriptions(start_date);
CREATE INDEX idx_end_date ON subscriptions(end_date);