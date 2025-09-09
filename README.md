# Auth Microservice Example

โครงสร้างระบบ microservices สำหรับ authentication และ authorization ด้วย Go

## โครงสร้างโปรเจกต์
```
auth-microservice/
├── api-gateway/
# Auth Microservice Example

ตัวอย่างระบบ microservices ด้วย Go — โฟกัสเรื่อง Authentication (JWT) และ API Gateway

## โครงสร้างโปรเจกต์ (ย่อ)

auth-microservice/
- api-gateway/         # edge service (ตรวจ JWT แล้ว forward)
- auth-service/        # register / login / validate (เชื่อม Postgres)
- user-service/        # user profile (ตัวอย่าง, in-memory)
- todo-service/        # todo API (ตัวอย่าง, in-memory)
- docker-compose.yml
- postgres/init.sql

## Tech stack
- Go + Echo framework
- JWT via github.com/golang-jwt/jwt/v5
- PostgreSQL (auth-service)

## Ports (default in this repo)
- API Gateway: 8080
- Auth Service: 8081
- Todo Service: 8083
- User Service: 8084

## Prerequisites
- Go 1.18+
- (สำหรับ DB) Docker & docker-compose หรือ Postgres พร้อมข้อมูลใน `.env`

## Environment (example)
วางไฟล์ `auth-microservice/.env` หรือใช้ env vars ใน shell

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=Microservice

JWT_SECRET=changeme
```

> หมายเหตุ: อย่าเก็บค่า secret/รหัสผ่านจริงใน repo — `.gitignore` ของโฟลเดอร์ `auth-microservice` มีรายการสำหรับ `.env` และ `*.exe` แล้ว

## Run (docker)
ในโฟลเดอร์ `auth-microservice`:

```powershell
docker-compose up --build -d
```

## Run (local, no Docker)
เปิด terminal แยกสำหรับแต่ละ service แล้วรัน:

```powershell
# ตั้ง JWT secret ใน session
$env:JWT_SECRET = "secret123"

# auth-service (ต้องมี DB or .env present)
cd auth-microservice/auth-service
$env:JWT_SECRET='secret123'; go run .

# api-gateway
cd ../api-gateway
$env:JWT_SECRET='secret123'; go run .

# todo-service
cd ../todo-service
go run .

# user-service
cd ../user-service
go run .
```

## Quick API examples

- Register (auth-service)

```bash
curl -X POST http://localhost:8081/register -H 'Content-Type: application/json' -d '{"email":"alice@example.com","password":"pass"}'
```

- Login → get access token

```bash
curl -X POST http://localhost:8081/login -H 'Content-Type: application/json' -d '{"email":"alice@example.com","password":"pass"}'
# response: { "access_token": "...", "refresh_token": "..." }
```

- Call todo-service directly (gateway forwards header `X-User-Email`)

```bash
curl -H 'X-User-Email: alice@example.com' http://localhost:8083/todos
```

- Call via gateway (gateway validates JWT then forwards)

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/todo/todos
```

> หมายเหตุ: ปัจจุบัน gateway forward path `/todo/*` ตรงๆ — ถ้าเรียกผ่าน gateway แล้วได้ 404 คุณอาจต้องเรียกตรงไปยัง `http://localhost:8083/todos` หรือแก้ gateway ให้ strip prefix `/todo` ก่อน forward

## Database
- ไฟล์ `postgres/init.sql` มีตัวอย่าง schema สำหรับ `users` table — ใช้เมื่อรัน Postgres ด้วย docker-compose

## Security / Notes
- อย่า commit `.env` หรือไฟล์ไบนารีไปใน repo
- สำหรับ production แนะนำให้ใช้ secret manager และ migration tool สำหรับ DB

## Contributing
- Fork → create feature branch → PR

License: MIT
