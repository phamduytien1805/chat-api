CREATE TABLE IF NOT EXISTS channels (
    id BIGINT PRIMARY KEY,                 -- Unique ID for each channel
    name TEXT,                           -- Channel name
    created_by UUID,                     -- User ID of the channel creator
    created_at TIMESTAMP,                -- Channel creation timestamp
    updated_at TIMESTAMP,                -- Channel last update timestamp
);