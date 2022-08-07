CREATE TABLE IF NOT EXISTS events (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title text,
    start_date timestamp with time zone,
    end_date timestamp with time zone,
    description text,
    user_id uuid,
    notify_at timestamp with time zone
);
