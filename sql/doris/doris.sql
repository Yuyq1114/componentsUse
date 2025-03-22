CREATE TABLE https_logs (
                            log_id BIGINT(20) NOT NULL COMMENT "日志ID",
                            type VARCHAR(50) NOT NULL COMMENT "请求/响应类型",
                            url VARCHAR(2048) NULL COMMENT "URL地址",
                            method VARCHAR(10) NULL COMMENT "HTTP方法",
                            status_code INT NULL COMMENT "HTTP状态码",
                            headers_json JSON NULL COMMENT "HTTP头信息",
                            content TEXT NULL COMMENT "HTTP内容",
                            create_time DATETIME DEFAULT CURRENT_TIMESTAMP
)
    ENGINE=OLAP
    DUPLICATE KEY(log_id, type)
DISTRIBUTED BY HASH(log_id) BUCKETS 10
PROPERTIES (
    "replication_num" = "1"
);
