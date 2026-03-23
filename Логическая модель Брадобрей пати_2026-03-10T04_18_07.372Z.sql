
CREATE TABLE IF NOT EXISTS "Пользователь" (
	"логин" ,
	"хэш пароля" ,
	"ФИО" 
);


CREATE TABLE IF NOT EXISTS "профиль (роль)" (
	"специализация" ,
	"расчетный оклад" ,
	"график " ,
	"ФИО" 
);


CREATE TABLE IF NOT EXISTS "Салон" (
	"адрес" ,
	"координаты" ,
	"режим работы" ,
	"лимит персонала" 
);


CREATE TABLE IF NOT EXISTS "Услуга" (
	"наименование" ,
	"стоимость" ,
	"длительность" 
);


CREATE TABLE IF NOT EXISTS "Бронирование" (
	"id" serial NOT NULL UNIQUE,
	"время начала" ,
	"время окончания" ,
	"статус" ,
	"итоговая стоимость" ,
	PRIMARY KEY("id")
);


CREATE TABLE IF NOT EXISTS "Материал" (
	"название" ,
	"единица измерения" 
);


CREATE TABLE IF NOT EXISTS "Складской запас" (
	"адрес салона" ,
	"количество" 
);


CREATE TABLE IF NOT EXISTS "норма расхода" (
	"услуга" ,
	"материал" ,
	"объем" 
);


CREATE TABLE IF NOT EXISTS "отзыв" (
	"текст отзыва" ,
	"оценка" 
);


CREATE TABLE IF NOT EXISTS "платеж" (
	"статус транзакции" ,
	"сумма" ,
	"время оплаты" 
);


CREATE TABLE IF NOT EXISTS "работа сотрудника" (
	"ФИО сотрудника" ,
	"Адрес салона" 
);


CREATE TABLE IF NOT EXISTS "факт бронирования" (
	"бронирование" serial NOT NULL UNIQUE,
	"услуга" ,
	PRIMARY KEY("бронирование")
);


ALTER TABLE "Материал"
ADD FOREIGN KEY("название") REFERENCES "Складской запас"("адрес салона")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "Салон"
ADD FOREIGN KEY("адрес") REFERENCES "Складской запас"("адрес салона")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "Материал"
ADD FOREIGN KEY("название") REFERENCES "норма расхода"("материал")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "Услуга"
ADD FOREIGN KEY("наименование") REFERENCES "норма расхода"("услуга")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "Пользователь"
ADD FOREIGN KEY("ФИО") REFERENCES "профиль (роль)"("ФИО")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "профиль (роль)"
ADD FOREIGN KEY("ФИО") REFERENCES "работа сотрудника"("ФИО сотрудника")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "Салон"
ADD FOREIGN KEY("адрес") REFERENCES "работа сотрудника"("Адрес салона")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "Услуга"
ADD FOREIGN KEY("наименование") REFERENCES "факт бронирования"("услуга")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "Бронирование"
ADD FOREIGN KEY("id") REFERENCES "факт бронирования"("бронирование")
ON UPDATE NO ACTION ON DELETE NO ACTION;