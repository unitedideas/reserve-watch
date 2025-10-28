-- Referral program tracking
CREATE TABLE IF NOT EXISTS referrals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    referrer_email TEXT NOT NULL, -- User who refers
    referred_email TEXT NOT NULL, -- User who was referred
    referral_code TEXT NOT NULL UNIQUE,
    status TEXT DEFAULT 'pending', -- 'pending', 'converted', 'credited'
    referred_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    converted_at DATETIME,
    credited_at DATETIME,
    credit_amount_cents INTEGER DEFAULT 1000, -- $10 in cents
    FOREIGN KEY(referrer_email) REFERENCES leads(email)
);

CREATE INDEX IF NOT EXISTS idx_referrals_referrer ON referrals(referrer_email);
CREATE INDEX IF NOT EXISTS idx_referrals_code ON referrals(referral_code);
CREATE INDEX IF NOT EXISTS idx_referrals_status ON referrals(status);

-- Referral credits ledger
CREATE TABLE IF NOT EXISTS referral_credits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_email TEXT NOT NULL,
    amount_cents INTEGER NOT NULL, -- Can be positive (earned) or negative (used)
    reason TEXT NOT NULL, -- 'referral_conversion', 'applied_to_subscription', etc.
    referral_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    balance_after_cents INTEGER NOT NULL,
    FOREIGN KEY(referral_id) REFERENCES referrals(id)
);

CREATE INDEX IF NOT EXISTS idx_credits_user ON referral_credits(user_email);
CREATE INDEX IF NOT EXISTS idx_credits_created ON referral_credits(created_at);

-- Social posts log (for auto-tweeting agent)
CREATE TABLE IF NOT EXISTS social_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    platform TEXT NOT NULL, -- 'twitter', 'linkedin'
    signal_key TEXT NOT NULL, -- 'dtwexbgs', 'vix', etc.
    signal_status TEXT NOT NULL, -- 'watch', 'crisis'
    content TEXT NOT NULL,
    posted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    post_id TEXT, -- External platform post ID
    engagement_count INTEGER DEFAULT 0,
    UNIQUE(signal_key, signal_status, posted_at) -- Prevent duplicate posts
);

CREATE INDEX IF NOT EXISTS idx_social_signal ON social_posts(signal_key, signal_status);
CREATE INDEX IF NOT EXISTS idx_social_posted ON social_posts(posted_at);

