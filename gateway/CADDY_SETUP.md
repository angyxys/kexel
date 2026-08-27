# Caddy & HTTPS Setup Guide

This guide explains how to configure Caddy with automatic HTTPS using Cloudflare DNS validation.

## Overview

Caddy is now used as the reverse proxy and gateway for Kexel. It automatically:
- Issues and renews SSL/TLS certificates
- Validates domain ownership using Cloudflare DNS API
- Redirects HTTP to HTTPS
- Proxies requests to backend and frontend services

## Prerequisites

1. A domain name (e.g., `api.example.com`)
2. The domain must be managed by Cloudflare
3. A Cloudflare API token with DNS edit permissions

## Getting Cloudflare API Token

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com/profile/api-tokens)
2. Click **Create Token**
3. Use the **Edit zone DNS** template (recommended)
4. Configure permissions:
   - **Zone → DNS → Edit** ✓
   - Zone Resources: Select your specific domain
5. Copy the generated token

## Configuration

### 1. Create `.env` file

Copy `.env.example` and configure it:

```env
# ... existing database configuration ...

# Your domain
DOMAIN=api.example.com

# Cloudflare API Token (from above)
CLOUDFLARE_API_TOKEN=your_token_here

# Email for certificate notifications
ACME_EMAIL=your-email@example.com
```

### 2. Point Domain DNS

In Cloudflare:
1. Add an A record pointing to your server's IP:
   - Name: `api` (or your subdomain)
   - Type: `A`
   - Content: Your server's public IP
   - Proxy status: DNS only (gray cloud)

### 3. Start Services

```bash
docker-compose up -d
```

Caddy will:
- Automatically request an SSL certificate from Let's Encrypt
- Use Cloudflare DNS API to validate domain ownership
- Store certificates in `caddy_data` volume
- Proxy requests to backend/frontend

## Verification

### Check Caddy Logs
```bash
docker-compose logs -f gateway
```

Look for messages like:
```
{"level":"info","msg":"autosaved config"}
{"level":"info","msg":"tls.dns.cloudflare","msg":"cleaning up"}
```

### Test HTTPS
```bash
curl https://api.example.com/api/vrc/list/vip
```

Should return JSON with VIP players.

## Troubleshooting

### "DNS Validation Failed"
- Verify Cloudflare API token is correct
- Check domain is managed by Cloudflare
- Ensure API token has DNS edit permissions
- Wait a few minutes for DNS propagation

### "Certificate Renewal Failed"
- Check Caddy logs: `docker-compose logs gateway`
- Verify API token is still valid
- Try restarting: `docker-compose restart gateway`

### Port Already in Use
- Ensure ports 80 and 443 are available
- Check: `netstat -tlnp | grep -E ':80|:443'`
- Firewall rules must allow these ports

## File Structure

```
gateway/
├── Caddyfile          # Caddy configuration
├── nginx.conf         # (Legacy - can be removed)
└── CADDY_SETUP.md     # This file
```

## Security Notes

1. **Keep API Token Secret** - Never commit to Git
2. **Rotate Regularly** - Regenerate API token periodically
3. **Limit Permissions** - Only grant DNS edit rights, not account admin
4. **HTTPS Only** - Caddy enforces HTTPS redirects automatically
5. **Rate Limiting** - Caddy includes built-in rate limiting for ACME

## Advanced Configuration

### Custom Headers
Edit `gateway/Caddyfile` to add custom response headers:

```
{$DOMAIN} {
  header Strict-Transport-Security "max-age=31536000; includeSubDomains"
  header X-Content-Type-Options "nosniff"
  # ... rest of configuration
}
```

### Environment Substitution
Variables in Caddyfile are replaced by values from `.env`:
- `{$DOMAIN}` → Value of `DOMAIN` env var
- `{$CLOUDFLARE_API_TOKEN}` → Value of `CLOUDFLARE_API_TOKEN`
- `{$ACME_EMAIL}` → Value of `ACME_EMAIL`

### Disable HTTPS (Development)
In development, you can use HTTP instead:
```bash
# Edit Caddyfile
http://{$DOMAIN} {
  # ... rest of configuration
}
```

Then restart: `docker-compose restart gateway`

## References

- [Caddy Documentation](https://caddyserver.com/docs/)
- [Caddy DNS Module](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuers/acme/challenges/dns)
- [Cloudflare API Tokens](https://developers.cloudflare.com/api/tokens/create/)
- [Let's Encrypt](https://letsencrypt.org/)

## Support

If you encounter issues:
1. Check Caddy logs: `docker-compose logs -f gateway`
2. Verify all environment variables are set
3. Test DNS resolution: `nslookup api.example.com`
4. Verify firewall allows ports 80/443
