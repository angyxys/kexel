# Kexel Development Roadmap

Complete feature implementation roadmap with milestones, timelines, and priorities.

## 📊 Overview

| Phase | Timeline | Focus | Features |
|-------|----------|-------|----------|
| **Phase 1** | Week 1-2 | Core Foundation | Logs, Search, Temp Bans |
| **Phase 2** | Week 3-4 | Admin Features | Dashboard, Invitations, Sessions |
| **Phase 3** | Week 5-6 | Security & Integration | 2FA, Webhooks, API Keys |
| **Phase 4** | Week 7-8 | External Integration | Discord Bot, Patreon Sync |
| **Phase 5** | Ongoing | Polish | Mobile, Analytics, Tickets |

---

## 🟥 Phase 1: Core Foundation (Week 1-2)

### Milestone 1.1: Audit Logging System
**Priority:** 🔴 CRITICAL | **Estimated:** 2 days

**What:**
- New `AuditLog` model tracking all changes
- Middleware to log all mutations
- Endpoint to list/filter audit logs
- Export to CSV/JSON

**Tasks:**
- [ ] Create `AuditLog` model with fields:
  - `id`, `user_id`, `action`, `resource_type`, `resource_id`, `old_value`, `new_value`, `ip_address`, `user_agent`, `created_at`
- [ ] Create `AuditLogRepository` with filtering
- [ ] Create audit middleware that logs all POST/PUT/DELETE
- [ ] Add `GET /web/audit-logs` endpoint with filters
- [ ] Add `GET /web/audit-logs/export` for CSV export
- [ ] Update frontend with new Audit Logs page
- [ ] Add table view with filters and pagination

**Backend Files:**
```
internal/database/models/models.go (add AuditLog)
internal/repository/audit_log.go (new)
internal/service/audit.go (new)
internal/handler/audit.go (new)
internal/middleware/audit.go (new)
cmd/api/main.go (register middleware)
```

**Frontend Files:**
```
src/pages/AuditLogs.tsx (new)
src/api/audit.ts (new)
src/components/AuditTable.tsx (new)
```

**Testing:**
```bash
# Verify logs are created
curl http://localhost:8080/web/audit-logs \
  -H "Authorization: Bearer <token>"
```

---

### Milestone 1.2: Advanced Search & Filters
**Priority:** 🟠 HIGH | **Estimated:** 1.5 days

**What:**
- Full-text search on players
- Advanced filtering (role, ban status, date range)
- Saved search filters
- Export filtered results

**Tasks:**
- [ ] Add `search_index` on player tables
- [ ] Implement full-text search in repository
- [ ] Create filter API endpoints
- [ ] Add filter UI components
- [ ] Implement saved filters (store in User preferences)
- [ ] Add bulk export functionality

**Backend Files:**
```
internal/repository/player.go (add search methods)
internal/service/player.go (add search methods)
internal/handler/player.go (add search endpoints)
```

**Frontend Files:**
```
src/components/SearchBar.tsx (new)
src/components/FilterPanel.tsx (new)
src/hooks/useFilters.ts (new)
src/pages/Dashboard.tsx (update)
```

**Testing:**
```bash
curl "http://localhost:8080/web/players/search?q=test&role=mod" \
  -H "Authorization: Bearer <token>"
```

---

### Milestone 1.3: Temporary Bans
**Priority:** 🟡 MEDIUM | **Estimated:** 1.5 days

**What:**
- Ban with expiration date
- Auto-unban when expired
- Reason tracking
- Suspend/unsuspend without deletion

**Tasks:**
- [ ] Update `Player` model with:
  - `ban_reason` text field
  - `ban_expires_at` timestamp (nullable)
- [ ] Create background job to check expired bans
- [ ] Update ban/unban endpoints
- [ ] Add ban reason to UI
- [ ] Show countdown to unban
- [ ] Add scheduled job runner

**Backend Files:**
```
internal/database/models/models.go (update Player)
internal/service/player.go (add unban logic)
internal/handler/player.go (update endpoints)
internal/jobs/cleanup.go (new - scheduled tasks)
```

**Frontend Files:**
```
src/components/BanModal.tsx (new - with expiry)
src/pages/Dashboard.tsx (update - show countdown)
```

**Testing:**
```bash
# Ban player for 7 days
curl -X PUT http://localhost:8080/web/player/usr_123 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "is_banned": true,
    "ban_reason": "Spam",
    "ban_expires_at": "2026-09-03T12:00:00Z"
  }'
```

---

## 🟧 Phase 2: Admin Features (Week 3-4)

### Milestone 2.1: Comprehensive Dashboard with Stats
**Priority:** 🟠 HIGH | **Estimated:** 2 days

**What:**
- KPI cards (total players, VIPs, mods, owners, bans)
- Charts with activity trends
- Recent actions list
- Quick actions

**Tasks:**
- [ ] Create stats API endpoint
- [ ] Add chart library (Recharts)
- [ ] Design dashboard layout
- [ ] Implement real-time updates
- [ ] Add period selection (7d, 30d, 90d, all)
- [ ] Cache stats for performance

**Backend Files:**
```
internal/service/stats.go (new)
internal/handler/stats.go (new)
cmd/api/main.go (register stats routes)
```

**Frontend Files:**
```
src/pages/Dashboard.tsx (redesign)
src/components/StatsCard.tsx (new)
src/components/ActivityChart.tsx (new)
src/components/QuickActions.tsx (new)
src/store/statsStore.ts (new)
```

**Dependencies:**
```bash
npm install recharts date-fns
```

---

### Milestone 2.2: Invitation System
**Priority:** 🟡 MEDIUM | **Estimated:** 2 days

**What:**
- Generate invitation codes
- Set expiry and max uses
- Auto-assign role
- Track code usage
- Revoke codes

**Tasks:**
- [ ] Create `InvitationCode` model with:
  - `code`, `created_by`, `role`, `max_uses`, `uses`, `expires_at`, `created_at`
- [ ] Create invitation service and repository
- [ ] Modify register endpoint to accept code
- [ ] Create admin invitation management page
- [ ] Add code generation logic
- [ ] Track who used which code

**Backend Files:**
```
internal/database/models/models.go (add InvitationCode)
internal/repository/invitation.go (new)
internal/service/invitation.go (new)
internal/handler/invitation.go (new)
internal/service/auth.go (update register method)
```

**Frontend Files:**
```
src/pages/Invitations.tsx (new)
src/components/GenerateCodeModal.tsx (new)
src/api/invitations.ts (new)
```

---

### Milestone 2.3: Session Management
**Priority:** 🟡 MEDIUM | **Estimated:** 1.5 days

**What:**
- View active sessions
- Revoke sessions
- Device detection
- Suspicious activity alerts

**Tasks:**
- [ ] Create `UserSession` model:
  - `id`, `user_id`, `token_hash`, `device_info`, `ip_address`, `last_activity`, `created_at`, `expires_at`
- [ ] Track sessions on login
- [ ] Add session list endpoint
- [ ] Add revoke endpoint
- [ ] Show device info (browser, OS, IP)
- [ ] Alert on new device login

**Backend Files:**
```
internal/database/models/models.go (add UserSession)
internal/repository/session.go (new)
internal/service/auth.go (track sessions)
internal/handler/sessions.go (new)
```

**Frontend Files:**
```
src/pages/Sessions.tsx (new)
src/components/SessionList.tsx (new)
src/api/sessions.ts (new)
```

---

## 🟥 Phase 3: Security & Integration (Week 5-6)

### Milestone 3.1: Two-Factor Authentication (TOTP)
**Priority:** 🔴 CRITICAL | **Estimated:** 2 days

**What:**
- Google Authenticator / Authy support
- Backup codes
- Account recovery
- Enforce 2FA for admins

**Tasks:**
- [ ] Add `UserOTP` model with secret + backup codes
- [ ] Implement TOTP generation and verification
- [ ] Add 2FA setup flow in user settings
- [ ] Create backup codes generation
- [ ] Enforce 2FA for owner/mod roles
- [ ] Recovery flows

**Backend Files:**
```
internal/database/models/models.go (add UserOTP)
internal/service/otp.go (new)
internal/handler/otp.go (new)
internal/middleware/otp.go (new - verify TOTP on login)
```

**Frontend Files:**
```
src/pages/TwoFactorSetup.tsx (new)
src/components/QRCodeDisplay.tsx (new)
src/components/BackupCodes.tsx (new)
src/pages/VerifyOTP.tsx (new)
```

**Dependencies:**
```bash
go get github.com/pquerna/otp
npm install qrcode.react
```

---

### Milestone 3.2: API Key Management
**Priority:** 🟠 HIGH | **Estimated:** 1.5 days

**What:**
- Generate API keys for external access
- Scopes/permissions per key
- Rate limiting by key
- Key rotation
- IP whitelist

**Tasks:**
- [ ] Create `APIKey` model with scopes
- [ ] Implement key validation middleware
- [ ] Create API key management endpoints
- [ ] Add IP whitelist support
- [ ] Implement rate limiter per key
- [ ] Add audit logging for API calls

**Backend Files:**
```
internal/database/models/models.go (add APIKey)
internal/repository/api_key.go (new)
internal/service/api_key.go (new)
internal/handler/api_key.go (new)
internal/middleware/api_key.go (new)
```

**Frontend Files:**
```
src/pages/APIKeys.tsx (new)
src/components/KeyGenerator.tsx (new)
src/api/api_keys.ts (new)
```

---

### Milestone 3.3: Webhooks System
**Priority:** 🟠 HIGH | **Estimated:** 2 days

**What:**
- Custom webhooks for events
- Event types (ban, unban, user.created, etc.)
- Retry logic
- Webhook history/logs
- Signature verification

**Tasks:**
- [ ] Create `Webhook` and `WebhookEvent` models
- [ ] Implement event publishing system
- [ ] Add webhook dispatcher
- [ ] Implement retry logic (exponential backoff)
- [ ] Add HMAC signature verification
- [ ] Create webhook management endpoints
- [ ] Add webhook testing tool

**Backend Files:**
```
internal/database/models/models.go (add Webhook, WebhookEvent)
internal/repository/webhook.go (new)
internal/service/webhook.go (new)
internal/handler/webhook.go (new)
internal/events/dispatcher.go (new)
internal/jobs/webhook_dispatcher.go (new)
```

**Frontend Files:**
```
src/pages/Webhooks.tsx (new)
src/components/WebhookForm.tsx (new)
src/components/WebhookHistory.tsx (new)
src/api/webhooks.ts (new)
```

---

## 🟠 Phase 4: External Integration (Week 7-8)

### Milestone 4.1: Discord Bot
**Priority:** 🟡 MEDIUM | **Estimated:** 2.5 days

**What:**
- Slash commands for player management
- Automatic logs in Discord channel
- Role synchronization
- Member notifications

**Tasks:**
- [ ] Create Discord bot application
- [ ] Implement command handlers
- [ ] Add commands: `/player ban`, `/player list`, `/stats`
- [ ] Sync Discord roles to Kexel
- [ ] Send logs to designated channel
- [ ] Handle DM notifications
- [ ] Environment variables for bot config

**Files:**
```
cmd/discord-bot/main.go (new)
internal/discord/bot.go (new)
internal/discord/commands.go (new)
internal/discord/sync.go (new)
```

**Dependencies:**
```bash
go get github.com/bwmarrin/discordgo
```

---

### Milestone 4.2: Patreon Integration
**Priority:** 🟡 MEDIUM | **Estimated:** 1.5 days

**What:**
- Sync Patreon supporters
- Auto-assign VIP role
- Tier-based role assignment
- Monthly sync jobs

**Tasks:**
- [ ] Patreon API integration
- [ ] Create sync service
- [ ] Map Patreon tiers to Kexel roles
- [ ] Scheduled sync job (daily/weekly)
- [ ] Handle unpatrons
- [ ] Add Patreon user ID field

**Backend Files:**
```
internal/integrations/patreon/client.go (new)
internal/integrations/patreon/sync.go (new)
internal/jobs/patreon_sync.go (new)
```

---

### Milestone 4.3: Rate Limiting & DoS Protection
**Priority:** 🔴 CRITICAL | **Estimated:** 1.5 days

**What:**
- Global rate limiting
- Per-user rate limiting
- IP-based limiting
- Bypass for admins
- DDoS protection

**Tasks:**
- [ ] Implement rate limiter middleware
- [ ] Add Redis for distributed rate limiting
- [ ] Create configurable limits per endpoint
- [ ] Add admin bypass
- [ ] Implement exponential backoff
- [ ] Add metrics/monitoring

**Backend Files:**
```
internal/middleware/rate_limit.go (new)
internal/ratelimit/limiter.go (new)
```

**Dependencies:**
```bash
go get github.com/redis/go-redis/v9
go get github.com/JGLTechnologies/gin-rate-limit
```

---

## 🟡 Phase 5: Polish & Optimization (Ongoing)

### Milestone 5.1: Mobile Responsive Design
**Priority:** 🟡 MEDIUM | **Estimated:** 1.5 days

**What:**
- Mobile-optimized layouts
- Touch-friendly controls
- Responsive tables
- Mobile navigation

**Tasks:**
- [ ] Test on mobile devices
- [ ] Optimize table views for mobile
- [ ] Add hamburger menu
- [ ] Touch-friendly buttons
- [ ] Swipe gestures for actions

**Frontend Files:**
```
src/components/MobileNav.tsx (new)
src/hooks/useIsMobile.ts (new)
Update all pages for responsiveness
```

---

### Milestone 5.2: Analytics Dashboard
**Priority:** 🟡 MEDIUM | **Estimated:** 2 days

**What:**
- Advanced charts and graphs
- Trend analysis
- Player behavior analytics
- Predictive models
- PDF exports

**Tasks:**
- [ ] Create analytics service
- [ ] Implement data aggregation
- [ ] Add chart variety (line, bar, pie, heatmap)
- [ ] Time-series analysis
- [ ] PDF export functionality

**Frontend Files:**
```
src/pages/Analytics.tsx (new)
src/components/TrendChart.tsx (new)
src/components/HeatmapChart.tsx (new)
```

**Dependencies:**
```bash
npm install pdfkit chart.js
```

---

### Milestone 5.3: Support Tickets System
**Priority:** 🟡 MEDIUM | **Estimated:** 2 days

**What:**
- Players can report issues
- Ticket tracking
- Assignment to mods
- Categories and priority
- SLA tracking

**Tasks:**
- [ ] Create `Ticket` model
- [ ] Ticket CRUD endpoints
- [ ] Add ticket list UI
- [ ] Implement assignment logic
- [ ] Add comments to tickets
- [ ] Auto-resolve old tickets

**Backend Files:**
```
internal/database/models/models.go (add Ticket, TicketComment)
internal/repository/ticket.go (new)
internal/service/ticket.go (new)
internal/handler/ticket.go (new)
```

**Frontend Files:**
```
src/pages/Tickets.tsx (new)
src/pages/TicketDetail.tsx (new)
src/components/TicketForm.tsx (new)
src/api/tickets.ts (new)
```

---

## 📋 Implementation Checklist

### Pre-Implementation
- [ ] Review requirements
- [ ] Set up development environment (`docker-compose.dev.yml`)
- [ ] Create feature branch
- [ ] Write tests first (TDD)

### Per Milestone
- [ ] Implement backend changes
- [ ] Write unit tests (backend)
- [ ] Create/update API documentation
- [ ] Implement frontend changes
- [ ] Write integration tests
- [ ] Test locally
- [ ] Code review
- [ ] Update changelog
- [ ] Deploy to staging
- [ ] User acceptance testing
- [ ] Merge to main

### Quality Standards
- ✅ 80%+ test coverage
- ✅ No TypeScript errors
- ✅ No Go lint warnings
- ✅ Security scanning passed
- ✅ Performance benchmarks met
- ✅ Documentation updated

---

## 🚀 Getting Started

### Start Phase 1 Now:

```bash
# 1. Set up development environment
docker-compose -f docker-compose.dev.yml up

# 2. Create feature branch for Milestone 1.1
git checkout -b feat/audit-logging

# 3. Start implementing audit log model
# See Milestone 1.1 tasks above

# 4. Test locally
curl http://localhost:8080/web/audit-logs \
  -H "Authorization: Bearer <your-token>"

# 5. Commit and push
git add .
git commit -m "feat: implement audit logging system"
git push origin feat/audit-logging
```

---

## 📞 Questions?

- Check DEVELOPMENT.md for setup help
- Review existing code patterns
- Ask in issues or discussions
- Update this roadmap as priorities change

**Last Updated:** 2026-08-27  
**Current Phase:** 1 (Core Foundation)  
**Next Milestone:** 1.1 (Audit Logging)
