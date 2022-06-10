CREATE TABLE clients
(
    id        BIGSERIAL PRIMARY KEY,
    login     TEXT      NOT NULL UNIQUE,
    password  TEXT      NOT NULL,
    full_name TEXT      NOT NULL,
    passport  TEXT      NOT NULL,
    birthday  DATE      NOT NULL,
    status    TEXT      NOT NULL DEFAULT 'INACTIVE' CHECK (status IN ('INACTIVE', 'ACTIVE')),
    created   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- комментарий
-- ALTER TABLE clients ADD COLUMN phone TEXT NOT NULL;

CREATE TABLE cards
(
    id       BIGSERIAL PRIMARY KEY,
    number   TEXT      NOT NULL,
    balance  BIGINT    NOT NULL DEFAULT 0,
    issuer   TEXT      NOT NULL CHECK ( issuer IN ('Visa', 'MasterCard', 'MIR') ),
    holder   TEXT      NOT NULL,
    owner_id BIGINT    NOT NULL REFERENCES clients,
    status   TEXT      NOT NULL DEFAULT 'INACTIVE' CHECK (status IN ('INACTIVE', 'ACTIVE')),
    type     TEXT      NOT NULL DEFAULT 'COMMON' CHECK (type IN ('COMMON', 'VIRTUAL', 'ADDITIONAL')),
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE  TABLE transactions
(
    id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT NOT NULL references cards(id),
    amount BIGINT NOT NULL CHECK (amount > 0),
    mcc TEXT,
    status TEXT NOT NULL DEFAULT 'COMPLETED' CHECK (status IN ('COMPLETED', 'REFUNDED', 'BLOCKED')),
    reciever_id BIGINT references cards(id),
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users
(
    id BIGSERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password TEXT UNIQUE NOT NULL,
    role TEXT NOT NULL CHECK ( role IN ('CLIENT', 'ADMIN', 'VIPCLIENT', 'OPERATOR')),
    client_id BIGINT references clients(id) DEFAULT Null
);

CREATE TABLE user_to_token
(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT references users(id) NOT NULL,
    token TEXT NOT NULL UNIQUE
);

