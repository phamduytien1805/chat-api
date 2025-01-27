CREATE TABLE IF NOT EXISTS messages (
    id BIGINT,                           -- Unique ID for each message
    channel_id UUID,                     -- Channel ID
    content TEXT,                        -- Message content
    sender_id UUID,                      -- User ID of the sender
    created_at TIMESTAMP,                -- Channel creation timestamp
    updated_at TIMESTAMP,                -- Message last update timestamp
    PRIMARY KEY ((channel_id, bucket), message_id)
) WITH CLUSTERING ORDER BY (message_id DESC);