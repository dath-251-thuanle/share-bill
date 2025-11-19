# ShareBill – Group Expense Splitting System
### DATH-251 – Team N8  
**Project 3:** A web application that helps travel/offline groups track shared expenses, automatically calculate who owes whom, and generate accurate VietQR codes for settlement.

## Project Overview
ShareBill allows a group of friends to easily record expenses during a trip, automatically computes debts/surpluses for each member, and instantly generates personalized VietQR codes so everyone can settle with one scan (including refunds for those who overpaid).

Production URL (after deployment): https://bill.thuanle.me

## Team N8 Members
| Student ID | Full Name            | Role                        
|------------|----------------------|-----------------------------|
| 2313624    | Trần Đỗ Cao Trí      | Leader                      |
| 2352708    | Đinh Cao Thiên Lộc   |                             | 
| 2311883    | Nguyễn Thị Kim Loan  |                             |
| 2052xxxx   | Lê Văn C             |                             |
| 2052xxxx   | Phạm Thị D           |                             |
| 2052xxxx   | Phạm Thị D           |                             |
| 2052xxxx   | Phạm Thị D           |                             |

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
| Backend     | Go 1.23 + Fiber + GORM + PostgreSQL + Redis                 |
| Frontend    | Next.js 14 (App Router) + TypeScript + TailwindCSS +        |
|             | TanStack Query + Zod                                        |
| Database    | PostgreSQL 16 + Redis 7                                     |
| QR Code     | vietqr.io API + go-qrcode                                   |
| Deployment  | Docker Compose + GitHub Actions (self-hosted runners)       |
| Secrets     | mozilla/sops + age encryption                               |

## Project Structure
├───.github
│   └───workflows
├───docs            # All specifications are stored in here
│   ├───M1
│   ├───M2
│   └───M3
├───logs
├───reports
├───scripts
├───src             # Backend and Frontend code are implemented
└───uploads

## Local Development
```bash
git clone https://github.com/dath-251-thuanle/share-bill.git
cd share-bill
sops -d dev.enc.json > .env          # ask leader for sops key
docker compose up --build -d

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
