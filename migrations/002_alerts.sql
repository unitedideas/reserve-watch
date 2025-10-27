-- Alerts table
CREATE TABLE IF NOT EXISTS alerts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_email TEXT NOT NULL,
    name TEXT NOT NULL,
    series_id TEXT NOT NULL,
    condition TEXT NOT NULL, -- 'above' or 'below'
    threshold REAL NOT NULL,
    webhook_url TEXT,
    is_active INTEGER DEFAULT 1,
    last_triggered_at TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    UNIQUE(user_email, series_id, condition, threshold)
);

CREATE INDEX IF NOT EXISTS idx_alerts_active ON alerts(is_active);
CREATE INDEX IF NOT EXISTS idx_alerts_user ON alerts(user_email);

-- Alert history table
CREATE TABLE IF NOT EXISTS alert_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alert_id INTEGER NOT NULL,
    series_id TEXT NOT NULL,
    value REAL NOT NULL,
    threshold REAL NOT NULL,
    triggered_at TEXT DEFAULT (datetime('now')),
    webhook_status TEXT, -- 'success', 'failed', 'skipped'
    FOREIGN KEY(alert_id) REFERENCES alerts(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_alert_history_alert ON alert_history(alert_id);
CREATE INDEX IF NOT EXISTS idx_alert_history_time ON alert_history(triggered_at);

