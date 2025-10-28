-- Email leads captured from exit-intent and other forms
CREATE TABLE IF NOT EXISTS leads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    source TEXT NOT NULL, -- 'exit_intent', 'footer_form', 'pricing_page', etc.
    captured_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_email_sent_at DATETIME,
    drip_stage INTEGER DEFAULT 0, -- 0=welcome, 1=day2, 2=day4, 3=day7, 4=converted
    converted_at DATETIME,
    unsubscribed_at DATETIME,
    metadata TEXT -- JSON for additional data
);

CREATE INDEX IF NOT EXISTS idx_leads_email ON leads(email);
CREATE INDEX IF NOT EXISTS idx_leads_drip ON leads(drip_stage, captured_at);
CREATE INDEX IF NOT EXISTS idx_leads_source ON leads(source);

-- Email send log for tracking
CREATE TABLE IF NOT EXISTS email_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    lead_id INTEGER NOT NULL,
    email_type TEXT NOT NULL, -- 'welcome', 'day2', 'day7', etc.
    sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    opened_at DATETIME,
    clicked_at DATETIME,
    status TEXT DEFAULT 'sent', -- 'sent', 'opened', 'clicked', 'bounced', 'failed'
    FOREIGN KEY(lead_id) REFERENCES leads(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_email_log_lead ON email_log(lead_id);
CREATE INDEX IF NOT EXISTS idx_email_log_status ON email_log(status);

