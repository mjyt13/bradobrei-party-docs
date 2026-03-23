-- Сначала удалим таблицы, если они были, чтобы создать "чистую" базу
DROP TABLE IF EXISTS "booking_service", "works", "payment", "feedback", "service_consumption", 
                     "supplies", "booking", "service", "salon", "workers", "users", 
                     "profiles", "booking_statuses", "evaluations", "transaction_statuses", "materials" CASCADE;

-- 1. СПРАВОЧНИКИ (Родительские таблицы)
CREATE TABLE "profiles" (
    "id" serial PRIMARY KEY,
    "specialization" varchar(50)
);

CREATE TABLE "booking_statuses" (
    "id" serial PRIMARY KEY,
    "status_name" varchar(20)
);

CREATE TABLE "evaluations" (
    "id" serial PRIMARY KEY,
    "name" varchar(20)
);

CREATE TABLE "transaction_statuses" (
    "id" serial PRIMARY KEY,
    "name" varchar(20)
);

CREATE TABLE "materials" (
    "id" serial PRIMARY KEY,
    "name" varchar(255),
    "measure" varchar(255)
);

CREATE TABLE "service" (
    "id" serial PRIMARY KEY,
    "name" varchar(255),
    "cost" decimal(10,2),
    "duration" int
);

-- 2. ОСНОВНЫЕ СУЩНОСТИ
CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "login" varchar(20) UNIQUE,
    "pass_hash" varchar(255),
    "first_name" varchar(255),
    "last_name" varchar(255),
    "phone_number" varchar(255),
	"loyalty_status" varchar(50) DEFAULT 'Base'
);

-- Наследование пользователя (Мастер)
CREATE TABLE "workers" (
    "id" int PRIMARY KEY REFERENCES "users"("id"),
    "specialization_id" int REFERENCES "profiles"("id"),
    "base_salary" decimal(10,2) NOT NULL,
    "schedule" json
);

CREATE TABLE "salon" (
    "id" serial PRIMARY KEY,
    "address" text,
    "location" text,
    "work_hours" varchar(255),
    "staff_limit" int
);

-- 3. ОПЕРАЦИОННЫЕ ТАБЛИЦЫ
CREATE TABLE "booking" (
    "id" serial PRIMARY KEY,
    "booking_user_id" int REFERENCES "users"("id"),
    "time_start" TIMESTAMPTZ,
    "time_end" TIMESTAMPTZ,
    "status_id" int REFERENCES "booking_statuses"("id"),
    "total_cost" decimal(10,2)
);

CREATE TABLE "feedback" (
    "id" serial PRIMARY KEY,
    "feedback_text" text,
    "evaluation_id" int REFERENCES "evaluations"("id"),
    "booking_id" int REFERENCES "booking"("id") -- добавил связь с бронью
);

CREATE TABLE "payment" (
    "id" serial PRIMARY KEY,
    "status_id" int REFERENCES "transaction_statuses"("id"),
    "total_cost" decimal(10,2),
    "payment_time" TIMESTAMPTZ,
    "booking_id" int REFERENCES "booking"("id") -- связь с бронью
);

-- 4. ТАБЛИЦЫ СВЯЗЕЙ (Многие-ко-многим)
CREATE TABLE "supplies" (
    "salon_id" int REFERENCES "salon"("id"),
    "material_id" int REFERENCES "materials"("id"),
    "quantity" int,
    PRIMARY KEY ("salon_id", "material_id")
);

CREATE TABLE "service_consumption" (
    "service_id" int REFERENCES "service"("id"),
    "material_id" int REFERENCES "materials"("id"),
    "amount_per_service" decimal(12,3),
    PRIMARY KEY ("service_id", "material_id")
);

CREATE TABLE "works" (
    "worker_id" int REFERENCES "workers"("id"),
    "salon_id" int REFERENCES "salon"("id"),
    PRIMARY KEY ("worker_id", "salon_id")
);

CREATE TABLE "booking_service" (
    "booking_id" int REFERENCES "booking"("id"),
    "service_id" int REFERENCES "service"("id"),
    PRIMARY KEY ("booking_id", "service_id")
);