# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in this project, please report it by emailing [rahmatrafiq.1999@gmail.com](mailto:rahmatrafiq.1999@gmail.com).

**Please do NOT create a public GitHub issue for security vulnerabilities.**

## Security Incident History

### 2026-02-11: Hardcoded Test Database Credentials

**Severity:** Medium (Test Environment Only)

**Issue:**
- Test database password `[REDACTED]` was accidentally committed to git history
- Found in commits:
  - `85c9ad2065c61ffbef0ad8ddae5717fb6f68b17d` (tests/setup_test_db.sh)
  - `399b4786a19fe0db301e75ed7513fa133db91761` (.env.testing)

**Impact:**
- Test database credentials exposed in public repository
- No production systems affected
- Password was only used for local development testing

**Resolution:**
- ✅ Removed `.env.testing` from repository (commit 6a007ec)
- ✅ Created `.env.testing.example` template without secrets
- ✅ Added `.env.testing` to `.gitignore`
- ✅ Created `.gitguardian.yaml` to ignore historical alerts
- ⚠️ **ACTION REQUIRED:** Rotate MySQL password (see below)

**Password Rotation Steps:**

```bash
# Option 1: Change root password
mysql -u root -p'[REDACTED]' -e "ALTER USER 'root'@'localhost' IDENTIFIED BY 'YourNewSecurePassword';"

# Option 2: Create dedicated test user (RECOMMENDED)
mysql -u root -p'[REDACTED]' -e "CREATE USER 'test_user'@'localhost' IDENTIFIED BY 'YourNewTestPassword';"
mysql -u root -p'[REDACTED]' -e "GRANT ALL PRIVILEGES ON golang_starter_kit_2025_test.* TO 'test_user'@'localhost';"
mysql -u root -p'[REDACTED]' -e "FLUSH PRIVILEGES;"
```

After rotation, update your local `.env.testing`:
```bash
MYSQL_USER=test_user
MYSQL_PASSWORD=YourNewTestPassword
```

## Security Best Practices

### Environment Variables

1. **Never commit `.env` files** - Always use `.env.example` templates
2. **Use strong, unique passwords** - Minimum 16 characters with mixed case, numbers, symbols
3. **Rotate credentials regularly** - At least every 90 days for test environments
4. **Use different credentials per environment** - Development, Testing, Staging, Production

### Secret Management

- ✅ All `.env*` files (except `.example`) are in `.gitignore`
- ✅ GitGuardian pre-commit hooks enabled (planned)
- ✅ `.gitguardian.yaml` configured to scan for secrets
- ⚠️ Consider using secret management tools:
  - [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/)
  - [HashiCorp Vault](https://www.vaultproject.io/)
  - [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault/)

### Database Security

- ✅ Use separate database users for different environments
- ✅ Grant minimum required privileges (principle of least privilege)
- ✅ Use SSL/TLS for database connections in production
- ✅ Enable MySQL audit logging for production databases

### Production Environment

**Critical Security Controls:**

1. **SKIP_AUTH Environment Variable**
   - ❌ MUST be `false` or unset in production
   - ✅ Production check enforced in `auth_middleware.go:27`
   - Fatal error if `SKIP_AUTH=true` in production

2. **JWT Secret Key**
   - ❌ DO NOT use default/example keys in production
   - ✅ Generate strong random key: `openssl rand -base64 64`
   - ✅ Store in secure secret management system

3. **Database Credentials**
   - ❌ DO NOT use root user in production
   - ✅ Create dedicated application user with limited privileges
   - ✅ Use connection pooling with reasonable limits

4. **Rate Limiting**
   - ✅ Enabled globally (1000 req/min default)
   - ✅ Auth endpoints: 5 req/15min
   - ✅ User endpoints: 10,000 req/hour

## Security Checklist for New Developers

- [ ] Copy `.env.example` to `.env` and set unique values
- [ ] Copy `.env.testing.example` to `.env.testing` for tests
- [ ] Never commit `.env` or `.env.testing` files
- [ ] Use strong, unique passwords (not example passwords)
- [ ] Install pre-commit hooks: `git config core.hooksPath .githooks`
- [ ] Review this SECURITY.md before making changes

## Dependencies Security

We use the following tools to monitor dependencies:

- **go mod tidy** - Keep dependencies clean and up-to-date
- **golangci-lint** - Static analysis including security checks
- **Dependabot** - Automated dependency updates (GitHub)

Update dependencies regularly:
```bash
go get -u ./...
go mod tidy
```

## Contact

For security concerns, contact: rahmatrafiq.1999@gmail.com
