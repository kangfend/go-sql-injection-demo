# Login API - SQL Injection Demo

⚠️ **WARNING**: Kode ini dibuat untuk tujuan edukasi tentang keamanan aplikasi. JANGAN gunakan di production!

## Tech Stack

- **Go 1.22+** - HTTP routing dengan method pattern
- **SQLite** - Database

## Endpoints

### 1. `/login` - VULNERABLE
Endpoint yang **vulnerable terhadap SQL Injection** karena:
- Menggunakan string concatenation langsung (`fmt.Sprintf`)
- Tidak menggunakan prepared statements
- Tidak melakukan sanitasi input

### 2. `/login-secure` - SECURE
Endpoint yang **aman** karena:
- Menggunakan prepared statements dengan parameterized query (`?`)
- Input otomatis di-escape oleh database driver
- Mencegah SQL injection

### 3. HTTP Method Routing (Go 1.22+)
Menggunakan fitur routing baru Go 1.22+ yang mendukung HTTP method langsung di route pattern

## Cara Menjalankan

```bash
# Install dependencies
go mod download

# Run server
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## Testing

### 1. Login Normal (Vulnerable Endpoint)
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

### 2. Login Normal (Secure Endpoint)
```bash
curl -X POST http://localhost:8080/login-secure \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

### 3. SQL Injection Attacks (Berhasil di /login, Gagal di /login-secure)

**a) Bypass authentication (OR injection):**
```bash
# Vulnerable endpoint - BERHASIL bypass
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin'\'' OR '\''1'\''='\''1", "password": "anything"}'

# Secure endpoint - GAGAL (treated as literal string)
curl -X POST http://localhost:8080/login-secure \
  -H "Content-Type: application/json" \
  -d '{"username": "admin'\'' OR '\''1'\''='\''1", "password": "anything"}'
```

**b) Comment-based injection:**
```bash
# Vulnerable - BERHASIL
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin'\''--", "password": "anything"}'

# Secure - GAGAL
curl -X POST http://localhost:8080/login-secure \
  -H "Content-Type: application/json" \
  -d '{"username": "admin'\''--", "password": "anything"}'
```

**c) UNION-based injection:**
```bash
# Vulnerable - BERHASIL inject data
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "'\'' UNION SELECT 1,'\''hacker'\'','\''hacked@evil.com'\'','\''admin'\''--", "password": "x"}'

# Secure - GAGAL
curl -X POST http://localhost:8080/login-secure \
  -H "Content-Type: application/json" \
  -d '{"username": "'\'' UNION SELECT 1,'\''hacker'\'','\''hacked@evil.com'\'','\''admin'\''--", "password": "x"}'
```

## Default Users

- Username: `admin`, Password: `admin123`, Role: `admin`
- Username: `user`, Password: `user123`, Role: `user`
- Username: `john`, Password: `john123`, Role: `user`

## Default Users

- Username: `admin`, Password: `admin123`, Role: `admin`
- Username: `user`, Password: `user123`, Role: `user`
- Username: `john`, Password: `john123`, Role: `user`

## Perbandingan Kode

### ❌ Vulnerable Version (`/login`)
```go
// BAHAYA: String concatenation langsung
query := fmt.Sprintf("SELECT id, username, email, role FROM users WHERE username='%s' AND password='%s'",
    loginReq.Username, loginReq.Password)
row := db.QueryRow(query)
```

### ✅ Secure Version (`/login-secure`)
```go
// AMAN: Prepared statement dengan parameterized query
query := "SELECT id, username, email, role FROM users WHERE username=? AND password=?"
row := db.QueryRow(query, loginReq.Username, loginReq.Password)
```

## Educational Purpose Only

Kode ini dibuat untuk:
- Belajar tentang SQL Injection
- Testing security tools
- Demonstrasi vulnerability

**JANGAN** digunakan di aplikasi production!
