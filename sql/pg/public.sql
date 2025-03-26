CREATE TABLE merge_https (
                             merge_id BIGINT PRIMARY KEY, -- Auto increment field for MergeId
                             hosts VARCHAR(255) NOT NULL,
                             create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                             update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

