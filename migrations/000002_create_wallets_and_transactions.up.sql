-- Создаем таблицу кошельков
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Создаем таблицу истории операций
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    sender_wallet_id INT,
    receiver_wallet_id INT,
    amount NUMERIC(15, 2) NOT NULL,
    type VARCHAR(20) NOT NULL, -- 'deposit', 'withdraw', 'transfer'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_wallet_id) REFERENCES wallets(id) ON DELETE SET NULL,
    FOREIGN KEY (receiver_wallet_id) REFERENCES wallets(id) ON DELETE SET NULL
);
