# Stripe Integration Setup Guide

## Overview
The Reserve Watch dashboard now includes Stripe checkout integration for Pro and Team subscriptions. This guide explains how to:
1. Set up Stripe products and pricing
2. Configure environment variables
3. Test the checkout flow
4. Handle webhooks (future)

## 1. Create Stripe Account
1. Go to https://stripe.com
2. Sign up or log in
3. Complete business verification if needed

## 2. Create Products & Prices

### Pro Plan ($9/month)
1. Go to **Products** → **Add Product**
2. Set:
   - Name: `Reserve Watch Pro`
   - Description: `Real-time data, live indices, custom alerts, CSV exports`
   - Pricing model: `Recurring`
   - Price: `$9.00 USD`
   - Billing period: `Monthly`
3. Click **Save product**
4. Copy the **Price ID** (starts with `price_...`)
5. Update `internal/web/pricing.go` line 231:
   ```go
   <button onclick="checkout('price_XXXXXXXXXXXXX', 'Pro')" class="cta-button" id="pro-btn">Start Pro Plan - $9/mo</button>
   ```

### Team Plan ($39/user/month)
1. Go to **Products** → **Add Product**
2. Set:
   - Name: `Reserve Watch Team`
   - Description: `Team dashboard, SSO, custom reports, priority support`
   - Pricing model: `Recurring`
   - Price: `$39.00 USD`
   - Billing period: `Monthly`
3. Click **Save product**
4. Copy the **Price ID** (starts with `price_...`)
5. Update `internal/web/pricing.go` line 249:
   ```go
   <button onclick="checkout('price_XXXXXXXXXXXXX', 'Team')" class="cta-button" id="team-btn">Contact Sales for Team</button>
   ```

## 3. Get API Keys
1. Go to **Developers** → **API keys**
2. Copy **Secret key** (starts with `sk_test_...` for test mode or `sk_live_...` for production)
3. Copy **Publishable key** (starts with `pk_test_...` or `pk_live_...`)

## 4. Configure Railway Environment Variables
You've already added these to Railway:
- `STRIPE_SECRET_KEY=sk_test_...` (or `sk_live_...`)
- `STRIPE_PUBLISHABLE_KEY=pk_test_...` (or `pk_live_...`)

## 5. Test the Checkout Flow

### Test Mode (Recommended First)
1. Use test API keys (`sk_test_...`, `pk_test_...`)
2. Update price IDs in `pricing.go` with your test product price IDs
3. Commit and push:
   ```bash
   git add internal/web/pricing.go
   git commit -m "Update Stripe price IDs for Pro and Team plans"
   git push origin main
   ```
4. Visit https://web-production-4c1d00.up.railway.app/pricing
5. Click **Start Pro Plan - $9/mo**
6. Use Stripe test card:
   - Card number: `4242 4242 4242 4242`
   - Expiry: Any future date (e.g., `12/34`)
   - CVC: Any 3 digits (e.g., `123`)
   - ZIP: Any 5 digits (e.g., `12345`)
7. Complete checkout
8. You should be redirected to `/success` with a success message

### Production Mode
1. Switch to live API keys (`sk_live_...`, `pk_live_...`) in Railway
2. Create live products in Stripe (follow same steps as test, but in live mode)
3. Update price IDs in `pricing.go` with live product price IDs
4. Test with a real card (charge will be real!)

## 6. Verify Checkout in Stripe Dashboard
1. Go to **Payments** → **All payments**
2. You should see the test payment
3. Go to **Customers** → **All customers**
4. You should see the test customer

## 7. Success Page
After successful checkout, users are redirected to `/success?session_id=...` which shows:
- ✅ Success message
- List of Pro features
- Session ID for reference
- Link back to dashboard

## 8. Error Handling
The checkout handles:
- Failed API calls (shows alert to user)
- Network errors (button re-enables)
- Missing Stripe keys (checkout will fail gracefully)

## 9. Next Steps (Future Features)
- [ ] Stripe webhooks to handle subscription events (created, canceled, failed payment)
- [ ] User database to track subscriptions
- [ ] Login/authentication system
- [ ] API key management for Pro users
- [ ] Usage tracking and limits

## 10. Testing Checklist
- [ ] Pro plan price ID updated in `pricing.go`
- [ ] Team plan price ID updated in `pricing.go`
- [ ] Stripe secret key added to Railway
- [ ] Code deployed to Railway
- [ ] Test checkout with `4242 4242 4242 4242`
- [ ] Redirected to success page
- [ ] Payment appears in Stripe dashboard

## Support
- Stripe docs: https://stripe.com/docs/checkout/quickstart
- Stripe test cards: https://stripe.com/docs/testing#cards
- Railway support: https://railway.app/help

