# 🚀 Quick Start - CI/CD Setup

## ⚡ Setup dalam 5 Menit

### 1. Docker Hub Setup (2 menit)

```bash
# Login ke Docker Hub
docker login

# Atau buat access token:
# 1. Buka https://hub.docker.com/settings/security
# 2. Klik "New Access Token"
# 3. Copy token untuk digunakan di secrets
```

### 2. GitHub Actions Setup (2 menit)

```bash
# 1. Buka repository di GitHub
# 2. Settings > Secrets and variables > Actions > New repository secret
# 3. Tambahkan secrets:

DOCKER_USERNAME: your_dockerhub_username
DOCKER_PASSWORD: your_dockerhub_token_or_password
```

### 3. GitLab CI Setup (1 menit)

```bash
# 1. Buka project di GitLab
# 2. Settings > CI/CD > Variables > Add variable
# 3. Tambahkan variables yang sama seperti GitHub
```

### 4. Test Locally

```bash
# Install tools
make install-tools

# Run tests
make test

# Run linter
make lint

# Test Docker build
make docker-build
```

## 📋 Quick Commands

### Development
```bash
make help              # Lihat semua commands
make dev-setup         # Setup development environment
make run               # Run application
make test              # Run tests
make lint              # Run linter
```

### CI/CD Testing
```bash
make ci-test           # Run full CI tests locally
make docker-build      # Build Docker image
make docker-run        # Run dengan Docker Compose
```

### Git Workflow
```bash
# Sebelum commit
make fmt               # Format code
make lint              # Check linting
make test              # Run tests

# Atau jalankan semuanya
make ci-test
```

## 🎯 What Happens After Push?

### Push ke `main` branch:
```
1. ✅ Run tests (MySQL + PostgreSQL)
2. ✅ Run linter
3. ✅ Build Docker image
4. ✅ Push to Docker Hub (tagged: latest, main-<sha>)
5. ✅ Security scan
6. ✅ Notify success
```

### Push ke `develop` branch:
```
1. ✅ Run tests
2. ✅ Run linter
3. ✅ Build Docker image
4. ✅ Push to Docker Hub (tagged: develop-<sha>)
```

### Create Pull Request:
```
1. ✅ Quick validation (formatting, go.mod)
2. ✅ Unit tests
3. ✅ Build check
4. ⏸️  No deployment
```

## 🔍 Monitoring

### View Pipeline Status

**GitHub Actions:**
```
Repository > Actions tab
```

**GitLab CI:**
```
Project > CI/CD > Pipelines
```

### View Docker Images

```bash
# List images
docker images | grep golang-starter-kit-2025

# Pull latest image
docker pull your_username/golang-starter-kit-2025:latest

# Run pulled image
docker run -p 8080:8080 your_username/golang-starter-kit-2025:latest
```

## 🐛 Common Issues

### Issue 1: Tests failing with database error
```bash
# Solution: Check database services are running
docker ps | grep mysql
docker ps | grep postgres

# Or run locally with Docker Compose
make docker-run
```

### Issue 2: Linter errors
```bash
# Fix formatting
make fmt

# Check what's wrong
make lint

# Fix common issues
go fmt ./...
go mod tidy
```

### Issue 3: Docker build fails
```bash
# Clean and rebuild
make clean
make docker-build

# Check Docker daemon
docker info
```

### Issue 4: Permission denied on secrets
```bash
# Verify secrets are set correctly:
# GitHub: Settings > Secrets > Check names match exactly
# GitLab: Settings > CI/CD > Variables > Check masked/protected
```

## 📊 Reading Pipeline Results

### GitHub Actions Output
```
✅ Green checkmark = Success
❌ Red X = Failed
🟡 Yellow dot = Running
⚪ Gray circle = Pending
```

### GitLab CI Output
```
✓ passed = Success
✗ failed = Failed
⟳ running = Running
○ pending = Pending
```

## 🎓 Next Steps

1. **Add badges to README**: See `CI_CD_SETUP.md`
2. **Setup Codecov**: For coverage reporting
3. **Configure VPS deployment**: Uncomment VPS job in GitLab CI
4. **Add Slack notifications**: For deployment alerts
5. **Setup staging environment**: Add staging branch workflow

## 📚 More Info

- Full documentation: `CI_CD_SETUP.md`
- Makefile commands: `make help`
- GitHub Actions: `.github/workflows/`
- GitLab CI: `.gitlab-ci.yml`

---

**Ready to go!** 🎉 Push your code and watch the magic happen!
