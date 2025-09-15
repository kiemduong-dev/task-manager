# Task Manager API

A RESTful API built with **Go**, **Gin**, **GORM**, and **PostgreSQL**.  
Features user authentication with JWT, role-based access control, and full CRUD operations for tasks.

---

## 🚀 Features
- User registration & login with **JWT token** authentication.
- Create, read, update, delete (CRUD) tasks.
- Assign due dates and categories to tasks.
- Role-based access control (admin/user).
- Logging and basic error handling.

---

## 🛠 Tech Stack
- **Go** (Golang)  
- **Gin** (HTTP web framework)  
- **GORM** (ORM for Go)  
- **PostgreSQL** (Database)

---

## ⚙️ Setup Instructions

### 1️⃣ Prerequisites
- Go 1.21+
- PostgreSQL 16+

### 2️⃣ Database
Create a database and user in PostgreSQL:
```sql
CREATE DATABASE taskdb;
CREATE USER taskuser WITH PASSWORD 'taskpass';
GRANT ALL PRIVILEGES ON DATABASE taskdb TO taskuser;
