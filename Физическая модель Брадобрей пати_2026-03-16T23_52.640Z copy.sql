-- Очистка
DROP TABLE IF EXISTS "booking_service", "works", "payment", "feedback", "service_consumption", 
                     "supplies", "booking", "service", "salon", "workers", "users", 
                     "profiles", "booking_statuses", "evaluations", "transaction_statuses", 
                     "materials", "roles" CASCADE;

-- Справочники
CREATE TABLE "roles" (
    "id" serial PRIMARY KEY,
    "name" varchar(50) UNIQUE NOT NULL
);

CREATE TABLE "profiles" ( -- Специализации для мастеров
    "id" serial PRIMARY KEY,
    "specialization" varchar(50)
);

CREATE TABLE "booking_statuses" ("id" serial PRIMARY KEY, "status_name" varchar(20));
CREATE TABLE "transaction_statuses" ("id" serial PRIMARY KEY, "name" varchar(20));
CREATE TABLE "evaluations" ("id" serial PRIMARY KEY, "name" varchar(20));

-- Сущности
CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "login" varchar(50) UNIQUE NOT NULL,
    "pass_hash" varchar(255) NOT NULL,
    "first_name" varchar(100),
    "last_name" varchar(100),
    "email" varchar(100),
    "phone_number" varchar(20),
    "role_id" int REFERENCES "roles"("id"),
    "loyalty_status" varchar(20) DEFAULT 'Base'
);

CREATE TABLE "workers" (
    "id" int PRIMARY KEY REFERENCES "users"("id"),
    "specialization_id" int REFERENCES "profiles"("id"),
    "base_salary" decimal(12,2),
    "schedule" jsonb
);

CREATE TABLE "salon" (
    "id" serial PRIMARY KEY,
    "address" text NOT NULL,
    "location" text, -- Для координат
    "work_hours" varchar(100),
    "staff_limit" int
);

CREATE TABLE "service" (
    "id" serial PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "cost" decimal(10,2) NOT NULL,
    "duration_min" int NOT NULL -- Длительность в минутах
);

CREATE TABLE "booking" (
    "id" serial PRIMARY KEY,
    "client_id" int REFERENCES "users"("id"),
    "master_id" int REFERENCES "workers"("id"),
    "salon_id" int REFERENCES "salon"("id"),
    "time_start" TIMESTAMPTZ NOT NULL,
    "time_end" TIMESTAMPTZ,
    "status_id" int REFERENCES "booking_statuses"("id"),
    "total_cost" decimal(10,2)
);

CREATE TABLE "materials" (
    "id" serial PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "measure" varchar(20)
);

-- Таблицы связей
CREATE TABLE "supplies" (
    "salon_id" int REFERENCES "salon"("id"),
    "material_id" int REFERENCES "materials"("id"),
    "quantity" decimal(12,3),
    PRIMARY KEY ("salon_id", "material_id")
);

CREATE TABLE "booking_service" (
    "booking_id" int REFERENCES "booking"("id"),
    "service_id" int REFERENCES "service"("id"),
    PRIMARY KEY ("booking_id", "service_id")
);

-- ЗАПОЛНЕНИЕ (DML)
-- Роли по ТЗ
INSERT INTO roles (name) VALUES 
('Administrator'), ('Client'), ('Basic Master'), ('Advanced Master'), 
('HR Specialist'), ('Accountant'), ('Network Manager');

INSERT INTO profiles (specialization) VALUES ('Барбер'), ('Топ-барбер'), ('Стажер');
INSERT INTO booking_statuses (status_name) VALUES ('Ожидание'), ('Завершено'), ('Отменено');
INSERT INTO transaction_statuses (name) VALUES ('Оплачено'), ('Ошибка');

-- 5 Салонов
INSERT INTO salon (address, work_hours, staff_limit) VALUES 
('ул. Тверская, 1', '09:00-21:00', 5), ('ул. Арбат, 10', '10:00-22:00', 3),
('пр-т Мира, 45', '09:00-21:00', 8), ('ул. Ленина, 5', '10:00-20:00', 4),
('ул. Пушкина, 12', '09:00-22:00', 6);

-- 5 Услуг
INSERT INTO service (name, cost, duration_min) VALUES 
('Мужская стрижка', 1500.00, 60), ('Стрижка бороды', 800.00, 30),
('Бритье опасной бритвой', 2000.00, 45), ('Детская стрижка', 1200.00, 60),
('Комплекс ухода', 3500.00, 90);

-- 5 Материалов
INSERT INTO materials (name, measure) VALUES 
('Шампунь', 'мл'), ('Масло для бороды', 'мл'), ('Лезвия', 'шт'), ('Гель', 'мл'), ('Тальк', 'гр');

-- 5 Пользователей (разные роли)
INSERT INTO users (login, pass_hash, first_name, role_id, loyalty_status) VALUES 
('admin_main', 'hash', 'Иван', 1, 'Base'),
('master_1', 'hash', 'Петр', 3, 'Base'),
('master_2', 'hash', 'Олег', 4, 'Base'),
('client_1', 'hash', 'Алексей', 2, 'Gold'),
('client_2', 'hash', 'Мария', 2, 'Silver');

-- Делаем двоих работниками
INSERT INTO workers (id, specialization_id, base_salary) VALUES (2, 1, 50000), (3, 2, 80000);

-- 5 Бронирований
INSERT INTO booking (client_id, master_id, salon_id, time_start, status_id, total_cost) VALUES 
(4, 2, 1, '2026-03-10 10:00', 2, 1500),
(4, 2, 1, '2026-03-15 12:00', 2, 1500),
(5, 3, 2, '2026-03-12 14:00', 2, 3500),
(4, 3, 1, '2026-03-16 16:00', 1, 800),
(5, 2, 3, '2026-03-17 11:00', 1, 2000);