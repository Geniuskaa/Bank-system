INSERT INTO clients(login, password, full_name, passport, birthday, status)
VALUES ('petr',
        'password-hash',
        'Пётр Николаевич Иванов',
        '8204 95523',
        '1970.01.30',
        'ACTIVE');

INSERT INTO clients(login, password, full_name, passport, birthday, status)
VALUES ('vasya', 'password-hash', 'Василий Николаевич Иванов',
        '8205 96563', '1970.01.30', 'ACTIVE'),
       ('masha', 'password-hash', 'Мария Ивановна Петрова',
        '8205 48839', '1990.11.21', 'ACTIVE'),
       ('dasha', 'password-hash', 'Дарья Ивановна Крылова',
        '8205 94483', '1995.04.27', 'ACTIVE');

INSERT INTO cards(number, balance, issuer, holder, owner_id, status, type, created)
VALUES ('5246 2472', 12500000, 'Visa', 'Vasay Ivanov', 2, 'ACTIVE', 'COMMON',  now()),
       ('7266 2628', 195000, 'Visa', 'Masha Petrova', 3, 'ACTIVE', 'COMMON', now()),
       ('1273 9484', 700000, 'Visa', 'Dasha Krylova', 4, 'ACTIVE', 'COMMON', now());

INSERT INTO transactions(sender_id, amount, mcc, status, date)
VALUES (1, 15000, '5425','COMPLETED', '2022-06-07 22:27');