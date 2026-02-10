# CI/CD Setup Documentation

## 📋 Overview

Proyek ini menggunakan CI/CD pipeline untuk otomatisasi testing, building, dan deployment dengan support untuk **GitHub Actions** dan **GitLab CI**.

## 🚀 Features

### GitHub Actions (`.github/workflows/ci-cd.yml`)
- ✅ Automated testing dengan MySQL & PostgreSQL
- ✅ Code linting dengan golangci-lint
- ✅ Docker image build dan push ke Docker Hub
- ✅ Security scanning dengan Trivy
- ✅ Code coverage reporting dengan Codecov
- ✅ Deployment notifications

### GitLab CI (`.gitlab-ci.yml`)
- ✅ Multi-stage pipeline (test, build, deploy, security)
- ✅ Parallel testing dengan database services
- ✅ Docker image build dan push
- ✅ Security scanning (Trivy & govulncheck)
- ✅ Code formatting checks
- ✅ Dependency vulnerability scanning

## 🔧 Setup Instructions

### 1. GitHub Actions Setup

#### Required Secrets
Tambahkan secrets berikut di GitHub repository settings (`Settings` > `Secrets and variables` > `Actions`):

```
DOCKER_USERNAME     # Docker Hub username
DOCKER_PASSWORD     # Docker Hub password or access token
```

#### Optional Secrets (untuk Codecov)
```
CODECOV_TOKEN       # Token untuk upload coverage (optional)
```

#### Cara menambahkan secrets:
1. Buka repository di GitHub
2. Klik `Settings` > `Secrets and variables` > `Actions`
3. Klik `New repository secret`
4. Tambahkan name dan value untuk setiap secret

### 2. GitLab CI Setup

#### Required Variables
Tambahkan variables berikut di GitLab project settings (`Settings` > `CI/CD` > `Variables`):

```
DOCKER_USERNAME     # Docker Hub username
DOCKER_PASSWORD     # Docker Hub password or access token
```

#### Optional Variables (untuk VPS deployment)
```
SSH_PRIVATE_KEY     # SSH private key untuk akses VPS
VPS_HOST           # IP address atau hostname VPS
VPS_USER           # SSH username untuk VPS
```

#### Cara menambahkan variables:
1. Buka project di GitLab
2. Klik `Settings` > `CI/CD` > `Variables`
3. Klik `Add variable`
4. Tambahkan key dan value untuk setiap variable
5. Pastikan untuk protect dan mask sensitive variables

### 3. Docker Hub Setup

#### Membuat Docker Hub Access Token
1. Login ke [Docker Hub](https://hub.docker.com)
2. Klik profile picture > `Account Settings`
3. Klik `Security` tab
4. Klik `New Access Token`
5. Beri nama token (contoh: "github-actions-ci")
6. Pilih permissions: `Read, Write, Delete`
7. Copy token dan simpan sebagai `DOCKER_PASSWORD` secret

#### Membuat Repository
1. Login ke Docker Hub
2. Klik `Create Repository`
3. Nama repository: `golang-starter-kit-2025`
4. Pilih visibility: Public atau Private
5. Klik `Create`

## 📊 Pipeline Workflow

### GitHub Actions Pipeline

```
┌─────────────┐
│   Push/PR   │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────┐
│  Test & Lint Job            │
│  - Setup MySQL & PostgreSQL │
│  - Run unit tests           │
│  - Run golangci-lint        │
│  - Upload coverage          │
└──────────┬──────────────────┘
           │
           ▼
    ┌──────────────┐
    │ Tests Pass?  │
    └──┬───────┬───┘
       │ No    │ Yes
       │       │
       ▼       ▼
    [FAIL]  ┌──────────────────┐
            │ Build Docker Job │
            │ - Build image    │
            │ - Push to Hub    │
            └────────┬─────────┘
                     │
                     ▼
            ┌─────────────────┐
            │ Security Scan   │
            │ - Trivy scan    │
            └────────┬────────┘
                     │
                     ▼
            ┌─────────────────┐
            │ Notify Success  │
            └─────────────────┘
```

### GitLab CI Pipeline

```
Stage 1: TEST
├── test:unit          # Unit tests dengan coverage
├── lint               # Code linting
└── format:check       # Format checking

Stage 2: BUILD
└── build:docker       # Build dan push Docker image

Stage 3: DEPLOY
└── deploy:dockerhub   # Deployment notification
    └── deploy:vps     # (Manual) Deploy ke VPS

Stage 4: SECURITY
├── security:trivy     # Container security scan
└── security:deps      # Dependency vulnerability check
```

## 🧪 Testing Locally

### Test dengan Docker
```bash
# Build Docker image
docker build -t golang-starter-kit-2025:test .

# Run container
docker run -p 8080:8080 golang-starter-kit-2025:test
```

### Test dengan Docker Compose
```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

### Run linter locally
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run --timeout=5m
```

### Run tests locally
```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

## 🎯 Triggered Events

### GitHub Actions
- **Push** ke branch `main` atau `develop`: Run full pipeline + deploy
- **Pull Request** ke `main` atau `develop`: Run test + build only
- **Manual**: Bisa trigger manual via GitHub Actions UI

### GitLab CI
- **Push** ke branch `main`: Run full pipeline + deploy
- **Push** ke branch `develop`: Run test + build only
- **Merge Request**: Run test + build only
- **Manual Jobs**: Deploy ke VPS (manual trigger)

## 📈 Monitoring & Badges

### GitHub Actions Badge
Tambahkan badge di README.md:

```markdown
![CI/CD](https://github.com/USERNAME/REPO/workflows/CI/CD%20Pipeline/badge.svg)
```

### GitLab CI Badge
Tambahkan badge di README.md:

```markdown
[![pipeline status](https://gitlab.com/USERNAME/PROJECT/badges/main/pipeline.svg)](https://gitlab.com/USERNAME/PROJECT/-/commits/main)
[![coverage report](https://gitlab.com/USERNAME/PROJECT/badges/main/coverage.svg)](https://gitlab.com/USERNAME/PROJECT/-/commits/main)
```

### Codecov Badge (GitHub)
```markdown
[![codecov](https://codecov.io/gh/USERNAME/REPO/branch/main/graph/badge.svg)](https://codecov.io/gh/USERNAME/REPO)
```

## 🔒 Security Best Practices

1. **Never commit secrets**: Gunakan environment variables dan secrets
2. **Use access tokens**: Jangan gunakan password langsung
3. **Limit token scope**: Berikan minimal permissions yang dibutuhkan
4. **Rotate tokens**: Update tokens secara berkala
5. **Enable branch protection**: Require CI pass sebelum merge
6. **Review security scan results**: Check Trivy dan govulncheck reports

## 🐛 Troubleshooting

### GitHub Actions

**Problem**: Tests failing dengan database connection error
```
Solution:
- Pastikan health checks di service configuration berjalan
- Tunggu database ready sebelum run tests
- Check database credentials di test .env
```

**Problem**: Docker push failed
```
Solution:
- Verify DOCKER_USERNAME dan DOCKER_PASSWORD secrets
- Pastikan Docker Hub repository sudah dibuat
- Check Docker Hub access token masih valid
```

### GitLab CI

**Problem**: Pipeline stuck di "waiting for runner"
```
Solution:
- Pastikan GitLab Runner aktif
- Check runner tags match dengan job tags
- Verify runner has docker executor
```

**Problem**: Cache not working
```
Solution:
- Clear cache: CI/CD > Pipelines > Clear runner caches
- Verify cache key configuration
- Check runner has write permission untuk cache
```

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitLab CI/CD Documentation](https://docs.gitlab.com/ee/ci/)
- [Docker Hub Documentation](https://docs.docker.com/docker-hub/)
- [golangci-lint Configuration](https://golangci-lint.run/usage/configuration/)
- [Trivy Documentation](https://aquasecurity.github.io/trivy/)

## 🤝 Contributing

Untuk menambahkan atau memodifikasi CI/CD pipeline:

1. Test changes locally terlebih dahulu
2. Update documentation jika ada perubahan significant
3. Test di feature branch sebelum merge ke main
4. Monitor first run setelah merge

## 📝 Notes

- Pipeline akan otomatis run untuk setiap push dan PR
- Docker images di-tag dengan commit SHA dan branch name
- Security scans hanya run di branch `main`
- Manual deployment ke VPS tersedia (commented by default)
- Test coverage dilaporkan ke Codecov (GitHub) atau GitLab Coverage

---

**Last Updated**: 2026-02-10
**Maintained by**: DevOps Team
