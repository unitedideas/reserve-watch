# Create GitHub Personal Access Token

## Quick Steps:

1. **Go to**: https://github.com/settings/tokens/new

2. **Note**: `reserve-watch deployment`

3. **Expiration**: Choose "7 days" (or shorter for security)

4. **Select scopes** (check these boxes):
   - ✅ `repo` (Full control of private repositories)
     - This includes: repo:status, repo_deployment, public_repo, repo:invite, security_events

5. **Click**: "Generate token" (green button at bottom)

6. **Copy the token** (it looks like: `ghp_xxxxxxxxxxxxxxxxxxxx`)

7. **Give me the token** - I'll use it once to push your code, then you can delete it

## Security Notes:
- The token acts like a password
- You can delete it immediately after I push
- It only has access to repositories (not account settings)
- Set short expiration for safety

That's it! Once you give me the token, I can push the code immediately.

