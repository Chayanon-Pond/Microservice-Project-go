# Auth Microservice Example

โครงสร้างระบบ microservices สำหรับ authentication และ authorization ด้วย Go

## โครงสร้างโปรเจกต์
```
auth-microservice/
├── api-gateway/
│   └── main.go
├── auth-service/
│   ├── handler.go
│   ├── service.go
│   ├── model.go
│   └── main.go
├── user-service/
│   ├── handler.go
│   ├── repository.go
│   ├── model.go
│   └── main.go
├── todo-service/
│   ├── handler.go
│   └── main.go
├── docker-compose.yml
└── postgres/
    └── init.sql
```

## Features
- Register, Login, JWT, Protected Route, Logout, Role/Permission
- ใช้ Echo framework
- ใช้ PostgreSQL เป็นฐานข้อมูล
- รันทุก service ด้วย Docker Compose

## วิธีเริ่มต้น
1. ติดตั้ง Go และ Docker
2. สั่ง `docker-compose up --build`
3. พัฒนาแต่ละ service ตามไฟล์ที่สร้างไว้

## License
MIT
