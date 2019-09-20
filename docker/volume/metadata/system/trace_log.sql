ATTACH TABLE trace_log
(
    `event_date` Date, 
    `event_time` DateTime, 
    `revision` UInt32, 
    `timer_type` Enum8('Real' = 0, 'CPU' = 1), 
    `thread_number` UInt32, 
    `query_id` String, 
    `trace` Array(UInt64)
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(event_date)
ORDER BY (event_date, event_time)
SETTINGS index_granularity = 1024
