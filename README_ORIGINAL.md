# Automated Content Engine (Go) — README

Build a **mostly automated Go-based content engine** for two products:

1) **De-Dollarization Dashboard** — aggregates official signals (foreign-exchange reserves, Renminbi payment share, Cross-Border Interbank Payment System activity, central-bank gold demand), exposes charts, indices, and alert webhooks.

2) **Crash-Drill Autopilot** — watches a few stress gauges and produces **approved, auditable checklists** (for example, build a short-term United States Treasury bill ladder; switch settlement rails; generate a gold “proof-of-holdings” pack).

---

## Non-Goals (to avoid conflict with 1Password)
- No password or secrets management.
- No access or device-gating.
- Read-only ingestion of identity or device sources can come later; current focus is **macro signals → content/alerts** and **signals → playbooks**.

---

## Stack (Go first)
- Language: **Go** (v1.22+)
- Scheduler: `robfig/cron` (or GitHub Actions / Cloudflare Workers Cron)
- HTTP/HTML: `net/http`, `goquery`
- PDF parsing: `ledongthuc/pdf`
- Images (charts / Open Graph): `fogleman/gg`
- OAuth 2.0: `golang.org/x/oauth2`
- Data store: SQLite (can swap to PostgreSQL)
- Optional site: static MDX or lightweight Next.js for posts (not required for v1)

---

## Repo Layout
```
/cmd/runner                 # main entrypoint (cron jobs)
/internal/config            # env + secrets
/internal/ingest            # cofer.go, fred.go, cips.go, rmbtracker.go, wgc.go
/internal/compose           # templates, chart & OG generators
/internal/publish           # linkedin.go, mailchimp.go, youtube.go
/internal/store             # sqlite helpers + migrations
/internal/util              # logging, retry, rate-limit
/templates                  # MDX/markdown with disclosure slot
/.github/workflows          # optional scheduled runs
```

---

## Environment Variables (example)
```
APP_ENV=dev
DB_DSN=file:content.db?_fk=1

# APIs
FRED_API_KEY=...
LINKEDIN_ACCESS_TOKEN=...      # Org or user
LINKEDIN_ORG_URN=urn:li:organization:XXXX
MAILCHIMP_API_KEY=...          # or Substack token
MAILCHIMP_SERVER_PREFIX=usXX
YOUTUBE_CLIENT_ID=...
YOUTUBE_CLIENT_SECRET=...
YOUTUBE_REFRESH_TOKEN=...

# Feature flags
PUBLISH_LINKEDIN=true
PUBLISH_MAILCHIMP=true
PUBLISH_YOUTUBE=true
AUTOPUBLISH=true               # set false for manual review
```

---

## Data Sources (officials; expand once, then acronym)
- **International Monetary Fund – COFER** (Currency Composition of Official Foreign Exchange Reserves): reserve-currency shares (quarterly).
- **Society for Worldwide Interbank Financial Telecommunication – Renminbi Tracker**: Renminbi global payment share (monthly PDF).
- **Cross-Border Interbank Payment System – CIPS**: participants and annual/daily throughput (site stats).
- **World Gold Council – WGC**: central-bank gold purchases and demand (monthly/quarterly posts).
- **Federal Reserve Economic Data – FRED**:
  - `VIXCLS` (Cboe Volatility Index close; market “nervousness” gauge).
  - `BAMLC0A4CBBB` (BBB corporate bond option-adjusted spread; credit-stress gauge).

---

## Core Cron Jobs (suggested)
```
# COFER check (weekly; actual updates quarterly)
0 9 * * 1

# SWIFT RMB Tracker check (monthly)
0 10 * * 3

# CIPS stats scrape (daily)
15 9 * * *

# Weekly Trigger Watch (VIX + BBB)
0 8 * * 5
```

---

## Output Artifacts (auto-generated per run)
- **Blog Note** (120–180 words).
- **LinkedIn caption** (1–2 short paragraphs).
- **Newsletter intro** (90–120 words).
- **20-second video script** (for a YouTube Short).
- **1 PNG chart** and **1 Open Graph image** (brand frame).
- **Disclosure footer** (advertising/endorsement + “not investment advice”).

---

## Product Logic

### A) De-Dollarization Dashboard
**Pipelines:** Normalize COFER, RMB Tracker, CIPS, WGC. Track `source_updated_at` and `ingested_at`.

**Indices:**
- **RMB Penetration Score** = function(payments share × reserve share × CIPS reach).
- **Reserve Diversification Pressure** = function(gold share trend + net central-bank buying).

**Alerts:**
- Threshold rules (for example, “RMB > 3.5%”, “CIPS daily average > X”, “Gold purchases > Y tonnes”).
- Webhooks and email delivery.

### B) Crash-Drill Autopilot
**Signals:** `VIXCLS` threshold; `BAMLC0A4CBBB` weekly change threshold.

**Playbooks (v1):**
1) Build or resize short-term United States Treasury bill ladder (steps, roles, forms).
2) Toggle settlement to Renminbi / CIPS for China-linked payables (bank cut-offs, counterparty checks).
3) Gold exposure update with **proof-of-holdings** PDF (vault statements + tokenized-gold reconciliation, if any).

**Runner:** Step list, role assignments, evidence uploads, immutable timestamps, Slack or email export.

---

## Minimal Interfaces (signatures)
```go
// Ingest
type SeriesPoint struct{ Date string; Value float64; Meta map[string]string }
type FetchResult struct{ Name string; Points []SeriesPoint; Changed bool; Err error }

func FetchCOFER() FetchResult
func FetchRMBTrackerPDF() FetchResult
func FetchCIPS() FetchResult
func FetchWGC() FetchResult
func FetchFRED(seriesID string) FetchResult // call for VIXCLS, BAMLC0A4CBBB

// Compose
type ComposeInput struct{
  Topic string            // "cofer"|"rmb"|"cips"|"gold"|"triggers"
  Data  map[string]any    // stats for templates
}
func EmitNote(ci ComposeInput) (paths struct{
  Blog, LinkedIn, Newsletter, Script, ChartPNG, OGPNG string
}, err error)

// Publish
func PostLinkedInUGC(text, link string) error
func SendMailchimp(html string) error
func UploadYouTubeShort(mp4Path, title, desc string) error
```

---

## Compliance (templates + linter)
Every template must render a **Disclosure** block:

> “This content uses official sources (International Monetary Fund COFER, SWIFT Renminbi Tracker, Cross-Border Interbank Payment System, World Gold Council, Federal Reserve Economic Data). Promotions are labeled ads; we may have commercial relationships. This is not investment advice.”

**Build step fails** if the disclosure placeholder is missing (`templates/lint.go`).

---

## 14-Day Build Plan (acceptance tests inline)

**Day 1–2 — Repo + cron skeleton; SQLite; env wiring**
- Create `/cmd/runner` with `robfig/cron`.
- Add feature flags and kill-switch (`AUTOPUBLISH=false`).

**Day 3 — FRED ingest + Trigger Watch compose**
- Fetch `VIXCLS` and `BAMLC0A4CBBB`; compute triggers and chart/OG images.
- **Test:** JSON fetched; PNG chart saved; markdown note rendered.

**Day 4 — CIPS scraper**
- Scrape participants and annual volume; persist snapshot diff.
- **Test:** fields populated; compose a note + image.

**Day 5 — COFER ingest (SDMX JSON)**
- Parse latest reserve-share series; compute deltas.
- **Test:** DB rows written; note + chart emitted.

**Day 6 — SWIFT RMB Tracker watcher**
- Detect new monthly PDF; extract “global payments share” and ranking.
- **Test:** new PDF → parsed values → draft note.

**Day 7 — World Gold Council ingest**
- Pull latest central-bank purchases; store and compose.
- **Test:** totals parsed; note + chart emitted.

**Day 8 — Indices calculator**
- Compute **RMB Penetration** and **Reserve Diversification Pressure**; expose `/indices`.

**Day 9 — Publishers (LinkedIn + Mailchimp)**
- Implement LinkedIn UGC post and Mailchimp draft send.
- **Test:** private audience post; test-list email.

**Day 10 — YouTube Short generator**
- Convert 20-sec script → TTS audio → stitch with chart → upload unlisted Short.
- **Test:** unlisted video visible.

**Day 11 — Disclosure linter + feature flags**
- Enforce disclosure block; add `#manual-review` label to pause autoposts.

**Day 12 — Alert webhooks**
- Simple HTTP POST with payload when thresholds trip.

**Day 13 — Analytics + rollback**
- Log opens/clicks and post IDs; add global pause flag.

**Day 14 — Full rehearsal**
- Simulate COFER + RMB Tracker + weekly Trigger Watch.
- Expect: site note + LinkedIn + newsletter + Short within ~10 minutes of detection.

---

## Acceptance Tests (v1 summary)
- **Pipelines:** each source writes points and deltas; charts render; `Changed==true` only on new data.
- **Composer:** given `ComposeInput`, writes 4 text assets + 2 images in ≤10s.
- **Publishers:** LinkedIn (201/202), Mailchimp draft send, YouTube unlisted upload.
- **Compliance:** build fails without disclosure block.

---

## KPIs (week 3 targets)
- End-to-end publish within **10 minutes** of detected release.
- ≥ **35%** newsletter open rate; ≥ **4%** LinkedIn click-through.
- **≥ 3** design-partner signups for Dashboard alerts or Crash-Drill checklists.

---

## Nice-to-Haves (after v1)
- Treasury system “tile” (Kyriba/Coupa) via webhooks for distribution.
- Admin UI for thresholds, approvals, and content review.
- Playbook packs per bank/broker (licensed content with automatic updates).

---
