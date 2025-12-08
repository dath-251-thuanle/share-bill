-- 1. users
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(50) NOT NULL,
    email           VARCHAR(100) UNIQUE NOT NULL,
    password_hash   TEXT NOT NULL,
    bank_name       VARCHAR(50),
    account_number  VARCHAR(30),
    account_name    VARCHAR(100),              -- thêm để hiển thị QR đẹp hơn
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. events
CREATE TABLE events (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    owner_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_settled  BOOLEAN DEFAULT FALSE,      -- đánh dấu đã thanh toán xong chưa
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. participants (bảng trung gian user ↔ event)
CREATE TABLE participants (
    event_id    UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (event_id, user_id)
);

-- 4. expenses
CREATE TABLE expenses (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    description VARCHAR(255) NOT NULL,
    total       NUMERIC(18,2) NOT NULL CHECK (total >= 0),
    expense_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. expense_payers (nhiều người cùng trả 1 expense)
CREATE TABLE expense_payers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expense_id      UUID NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
    participant_id  UUID NOT NULL, -- sẽ dùng composite (event_id, user_id) ở tầng service
    user_id         UUID NOT NULL REFERENCES users(id),
    amount          NUMERIC(18,2) NOT NULL CHECK (amount > 0),
    UNIQUE(expense_id, user_id)
);

-- 6. expense_beneficiaries (tỉ lệ thụ hưởng)
CREATE TABLE expense_beneficiaries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expense_id      UUID NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
    participant_id  UUID NOT NULL, -- composite (event_id, user_id)
    user_id         UUID NOT NULL REFERENCES users(id),
    ratio           NUMERIC(12,4) NOT NULL CHECK (ratio > 0), -- ví dụ 1.5, 1.5, 3.0
    UNIQUE(expense_id, user_id)
);