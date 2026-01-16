# ShareBill – Group Expense Splitting System
### DATH-251 – Team N8  
**Project 3:** A web application that helps travel/offline groups track shared expenses, automatically calculate who owes whom, and generate accurate VietQR codes for settlement.

## Project Overview
ShareBill allows a group of friends to easily record expenses during a trip, automatically computes debts/surpluses for each member, and instantly generates personalized VietQR codes so everyone can settle with one scan (including refunds for those who overpaid).

## Team N8 Members
| Student ID | Full Name            | Role                        
|------------|----------------------|-----------------------------|
| 2313624    | Trần Đỗ Cao Trí      | Backend                     |
| 2352708    | Đinh Cao Thiên Lộc   | Backend                     | 
| 2311883    | Nguyễn Thị Kim Loan  | Backend                     |
| 2352918    | Nguyễn Lê Đức Phú    | Frontend                    |
| 2352770    | Trần Hà My           | Backend                     |
| 2153485    | Nguyễn Quang Khởi    | Frontend                    |
| 2311987    | Nguyễn Song Minh Luân| Frontend                    |

## Core Features
1. Create Event → auto-generated URL `/event/:eventId`
2. Manage participants (name + optional bank info)
3. Add expenses/transactions:
   - Description, amount
   - Payer(s)
   - Beneficiaries + custom split ratios (e.g., 1, 1, 1.5, 3…)
4. Real-time debt calculation (who owes / who is owed)
5. Assign “Collection Leader” → auto-generate QR for each member:
   - Owes money → QR to pay the exact amount to leader
   - Overpaid → QR for leader to refund the exact amount
6. Complete Docker + CI/CD pipeline with self-hosted runners

## Tech Stack
| Layer       | Technology                                                  |
|-------------|-------------------------------------------------------------|
| Backend     | Go 1.23 + Fiber + PostgreSQL + Redis                        |
| Frontend    | Next.js 14 (App Router) + TypeScript + TailwindCSS          |
| Database    | PostgreSQL 16 + Redis 7                                     |
| QR Code     | vietqr.io API + go-qrcode                                   |
| Deployment  | Docker Compose + GitHub Actions (self-hosted runners)       |

## Project Structure
├───share-bill-backend-v2 # Backend code is implemented
├───sharever-frontend     # Frontend code is implemented            
├───DATH__Final.pdf       # Report of the project
└───README.md

## Local Development
```bash
git clone https://github.com/dath-251-thuanle/share-bill.git
cd share-bill
sops -d dev.enc.json > .env          # ask leader for sops key
docker compose up --build -d

Backend (API, PostgreSQL, Redis, Nginx):

cd sharever-backend
docker compose up --build -d

Frontend (React + Vite)

cd sharever-frontend
npm install
npm run dev

Access:

Frontend : http://localhost
Backend API : http://localhost:8080
PostgreSQL : localhost:5432

Production Deployment

Server + 2 self-hosted runners already provided
Merge to main → GitHub Actions auto build & deploy
Command on prod: docker compose -f docker-compose.yml up -d --build

API Documentation

OpenAPI spec: /docs/api-spec.yaml
Swagger UI: https://bill.thuanle.me/api/docs

## Commit convention
feat:  fix:  ui:  test:  docs:  refactor:  chore:
