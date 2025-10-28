-- Series data points
CREATE TABLE IF NOT EXISTS series_points (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    series_name TEXT NOT NULL,
    date TEXT NOT NULL,
    value REAL NOT NULL,
    meta TEXT,
    source_updated_at TEXT,
    ingested_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(series_name, date)
);

CREATE INDEX IF NOT EXISTS idx_series_date ON series_points(series_name, date);

-- Published posts tracking
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    platform TEXT NOT NULL,
    post_id TEXT,
    series_name TEXT NOT NULL,
    content TEXT NOT NULL,
    chart_path TEXT,
    published_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT DEFAULT 'draft'
);

CREATE INDEX IF NOT EXISTS idx_posts_platform ON posts(platform, published_at);
