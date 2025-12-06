/*
Задача: 15:Покупатели <–>> 14:Учет отгрузки <<–> 02:Детали

Таблицы:
1. Покупатели (customers)
2. Учет отгрузки готовой продукции (shipments)
3. Справочник деталей (parts)
*/

BEGIN;

-- Удаление существующих объектов
DROP TABLE IF EXISTS shipments_audit CASCADE;
DROP TABLE IF EXISTS shipments CASCADE;
DROP TABLE IF EXISTS customers CASCADE;
DROP TABLE IF EXISTS parts CASCADE;

DROP FUNCTION IF EXISTS fn_cascade_delete_shipments() CASCADE;
DROP FUNCTION IF EXISTS fn_log_shipment_insert() CASCADE;
DROP PROCEDURE IF EXISTS p_customer_shipment_summary(INT, OUT DECIMAL(10,2), OUT DECIMAL(10,2));
DROP FUNCTION IF EXISTS fn_customer_count_by_city(TEXT);
DROP FUNCTION IF EXISTS fn_shipments_in_range(DATE, DATE);
DROP VIEW IF EXISTS v_full_shipment_info CASCADE;

-- ============================================================================
-- ТАБЛИЦЫ
-- ============================================================================

-- Справочник деталей (Файл02)
CREATE TABLE parts (
    part_code            TEXT PRIMARY KEY,
    part_type            TEXT NOT NULL CHECK (part_type IN ('покупная', 'собственного производства')),
    name                 TEXT NOT NULL,
    unit                 TEXT NOT NULL CHECK (unit IN ('шт','кг','м','компл')),
    plan_price           DECIMAL(10,2) NOT NULL CHECK (plan_price >= 0),
    CONSTRAINT chk_part_code_not_empty CHECK (LENGTH(part_code) > 0)
);

-- Покупатели (Файл15)
CREATE TABLE customers (
    customer_id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name                 TEXT NOT NULL,
    address              TEXT NOT NULL DEFAULT 'Не указан',
    city                 TEXT NOT NULL
);

-- Учет отгрузки готовой продукции (Файл14)
CREATE TABLE shipments (
    warehouse_no         INT NOT NULL CHECK (warehouse_no > 0),
    shipment_doc_no      INT NOT NULL CHECK (shipment_doc_no > 0),
    customer_id          INT NOT NULL,
    part_code            TEXT NOT NULL,
    unit                 TEXT NOT NULL CHECK (unit IN ('шт','кг','м','компл')),
    qty                  DECIMAL(10,2) NOT NULL CHECK (qty > 0),
    shipment_date        DATE NOT NULL DEFAULT CURRENT_DATE,
    PRIMARY KEY (warehouse_no, shipment_doc_no),
    
    -- Связь с customers БЕЗ системного каскада (триггер будет создан)
    CONSTRAINT fk_shipment_customer FOREIGN KEY (customer_id) 
        REFERENCES customers(customer_id),
    
    -- Связь с parts С системным каскадом
    CONSTRAINT fk_shipment_part FOREIGN KEY (part_code) 
        REFERENCES parts(part_code)
        ON DELETE CASCADE ON UPDATE CASCADE
);

-- Таблица аудита для логирования операций
CREATE TABLE shipments_audit (
    audit_id             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    warehouse_no         INT,
    shipment_doc_no      INT,
    customer_id          INT,
    part_code            TEXT,
    qty                  DECIMAL(10,2),
    shipment_date        DATE,
    action               TEXT NOT NULL,
    action_time          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- ТРИГГЕРЫ
-- ============================================================================

-- Триггер 1: Каскадное удаление отгрузок при удалении покупателя
-- BEFORE DELETE - сначала удаляем дочерние записи, потом родительскую
CREATE OR REPLACE FUNCTION fn_cascade_delete_shipments() 
RETURNS TRIGGER 
LANGUAGE plpgsql 
AS $$
BEGIN
    DELETE FROM shipments 
    WHERE customer_id = OLD.customer_id;
    RETURN OLD;
END;
$$;

CREATE TRIGGER trg_customers_before_delete
BEFORE DELETE ON customers
FOR EACH ROW
EXECUTE FUNCTION fn_cascade_delete_shipments();

-- Триггер 2: Логирование вставок в shipments
CREATE OR REPLACE FUNCTION fn_log_shipment_insert() 
RETURNS TRIGGER 
LANGUAGE plpgsql 
AS $$
BEGIN
    INSERT INTO shipments_audit(
        warehouse_no, shipment_doc_no, customer_id, part_code, 
        qty, shipment_date, action
    )
    VALUES (
        NEW.warehouse_no, NEW.shipment_doc_no, NEW.customer_id, NEW.part_code,
        NEW.qty, NEW.shipment_date, 'INSERT'
    );
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_shipments_after_insert
AFTER INSERT ON shipments
FOR EACH ROW
EXECUTE FUNCTION fn_log_shipment_insert();

-- ============================================================================
-- ХРАНИМАЯ ПРОЦЕДУРА
-- ============================================================================

-- Процедура с выходными параметрами: суммарная информация по отгрузкам покупателя
CREATE OR REPLACE PROCEDURE p_customer_shipment_summary(
    IN p_customer_id INT,
    OUT total_qty DECIMAL(10,2),
    OUT total_value DECIMAL(10,2)
)
LANGUAGE plpgsql
AS $$
BEGIN
    -- Суммарное количество и стоимость отгрузок
    SELECT 
        COALESCE(SUM(s.qty), 0),
        COALESCE(SUM(s.qty * p.plan_price), 0)
    INTO total_qty, total_value
    FROM shipments s
    JOIN parts p ON s.part_code = p.part_code
    WHERE s.customer_id = p_customer_id;
END;
$$;

-- ============================================================================
-- ФУНКЦИИ
-- ============================================================================

-- Скалярная функция: количество покупателей в городе
CREATE OR REPLACE FUNCTION fn_customer_count_by_city(p_city TEXT) 
RETURNS INT 
LANGUAGE sql 
AS $$
    SELECT COUNT(*)::INT
    FROM customers
    WHERE city = p_city;
$$;

-- Табличная функция: список отгрузок в интервале дат
CREATE OR REPLACE FUNCTION fn_shipments_in_range(p_start DATE, p_end DATE)
RETURNS TABLE(
    warehouse_no INT,
    shipment_doc_no INT,
    customer_id INT,
    customer_name TEXT,
    part_code TEXT,
    part_name TEXT,
    qty DECIMAL(10,2),
    shipment_date DATE
)
LANGUAGE sql
AS $$
    SELECT 
        s.warehouse_no, 
        s.shipment_doc_no, 
        s.customer_id,
        c.name,
        s.part_code,
        p.name,
        s.qty, 
        s.shipment_date
    FROM shipments s
    JOIN customers c ON s.customer_id = c.customer_id
    JOIN parts p ON s.part_code = p.part_code
    WHERE s.shipment_date BETWEEN p_start AND p_end
    ORDER BY s.shipment_date;
$$;

-- ============================================================================
-- VIEW: Объединение трех таблиц
-- ============================================================================

CREATE VIEW v_full_shipment_info AS
SELECT 
    s.warehouse_no,
    s.shipment_doc_no,
    s.shipment_date,
    s.qty,
    c.customer_id,
    c.name AS customer_name,
    c.address AS customer_address,
    c.city AS customer_city,
    p.part_code,
    p.name AS part_name,
    p.part_type,
    p.unit,
    p.plan_price,
    (s.qty * p.plan_price) AS total_price
FROM shipments s
JOIN customers c ON s.customer_id = c.customer_id
JOIN parts p ON s.part_code = p.part_code;

-- ============================================================================
-- ТЕСТОВЫЕ ДАННЫЕ
-- ============================================================================

-- Детали (15 штук, разные типы и цены)
INSERT INTO parts (part_code, part_type, name, unit, plan_price) VALUES
('D001', 'покупная', 'Болт М10', 'шт', 5.50),
('D002', 'покупная', 'Гайка М10', 'шт', 3.20),
('D003', 'собственного производства', 'Корпус редуктора', 'шт', 150.00),
('D004', 'собственного производства', 'Вал приводной', 'шт', 280.00),
('D005', 'покупная', 'Подшипник 205', 'шт', 85.00),
('D006', 'собственного производства', 'Шестерня Z=20', 'шт', 120.00),
('D007', 'покупная', 'Прокладка резиновая', 'шт', 12.00),
('D008', 'собственного производства', 'Фланец соединительный', 'шт', 95.00),
('D009', 'покупная', 'Болт М16', 'шт', 8.50),
('D010', 'покупная', 'Шайба плоская', 'шт', 1.20),
('D011', 'собственного производства', 'Муфта соединительная', 'шт', 185.00),
('D012', 'покупная', 'Сальник 40x60', 'шт', 45.00),
('D013', 'собственного производства', 'Крышка подшипника', 'шт', 75.00),
('D014', 'покупная', 'Смазка литиевая', 'кг', 320.00),
('D015', 'собственного производства', 'Вал-шестерня', 'шт', 450.00);

-- Покупатели (12 штук, разные города)
INSERT INTO customers (name, address, city) VALUES
('ООО "Техноком"', 'ул. Баумана, 15', 'Казань'),
('ЗАО "Механика"', 'пр. Ленина, 42', 'Москва'),
('ИП Иванов С.П.', 'ул. Пушкина, 7', 'Казань'),
('ООО "СтройМаш"', 'ул. Советская, 123', 'Самара'),
('АО "ПромТех"', 'бул. Победы, 88', 'Казань'),
('ООО "МеталлПром"', 'пер. Заводской, 5', 'Нижний Новгород'),
('ООО "Автодеталь"', 'ул. Гагарина, 33', 'Москва'),
('ИП Петров А.В.', 'ул. Кремлевская, 18', 'Казань'),
('ООО "ТехСервис"', 'пр. Мира, 56', 'Екатеринбург'),
('АО "Завод Точмаш"', 'ул. Индустриальная, 10', 'Челябинск'),
('ООО "РемМаш"', 'ул. Московская, 78', 'Казань'),
('ИП Сидоров К.Л.', 'пер. Солнечный, 3', 'Самара');

-- Отгрузки (35+ записей, разные склады и годы)
INSERT INTO shipments (warehouse_no, shipment_doc_no, customer_id, part_code, unit, qty, shipment_date) VALUES
-- 2023 год
(1, 1000, 1, 'D001', 'шт', 50, '2023-06-10'),
(2, 2000, 2, 'D002', 'шт', 100, '2023-07-15'),
(3, 3000, 3, 'D005', 'шт', 20, '2023-09-20'),

-- 2024 год
(1, 1001, 1, 'D001', 'шт', 100, '2024-01-15'),
(1, 1002, 2, 'D003', 'шт', 5, '2024-02-20'),
(2, 2001, 1, 'D004', 'шт', 3, '2024-03-10'),
(3, 3001, 3, 'D005', 'шт', 10, '2024-04-05'),
(5, 5001, 5, 'D003', 'шт', 2, '2024-05-12'),
(5, 5002, 5, 'D004', 'шт', 4, '2024-06-18'),
(4, 4000, 4, 'D009', 'шт', 150, '2024-07-22'),
(2, 2004, 6, 'D010', 'шт', 500, '2024-08-14'),
(3, 3004, 7, 'D011', 'шт', 8, '2024-09-28'),
(1, 1006, 8, 'D001', 'шт', 200, '2024-10-05'),
(5, 5006, 5, 'D006', 'шт', 15, '2024-11-12'),
(4, 4003, 9, 'D012', 'шт', 30, '2024-12-01'),

-- 2025 год (текущий год)
(1, 1003, 1, 'D006', 'шт', 8, '2025-01-10'),
(2, 2002, 2, 'D007', 'шт', 50, '2025-02-14'),
(3, 3002, 3, 'D003', 'шт', 6, '2025-03-20'),
(4, 4001, 4, 'D004', 'шт', 2, '2025-04-25'),
(5, 5003, 5, 'D006', 'шт', 10, '2025-05-30'),
(1, 1004, 1, 'D008', 'шт', 7, '2025-06-15'),
(2, 2003, 6, 'D005', 'шт', 15, '2025-07-08'),
(3, 3003, 7, 'D003', 'шт', 4, '2025-08-12'),
(5, 5004, 5, 'D004', 'шт', 3, '2025-09-18'),
(1, 1005, 1, 'D001', 'шт', 200, '2025-10-22'),
(4, 4002, 4, 'D006', 'шт', 5, '2025-11-05'),
(5, 5005, 5, 'D003', 'шт', 1, '2025-12-01'),
(1, 1007, 8, 'D013', 'шт', 12, '2025-01-25'),
(2, 2005, 9, 'D014', 'кг', 5, '2025-02-18'),
(3, 3005, 10, 'D015', 'шт', 3, '2025-03-05'),
(4, 4004, 11, 'D011', 'шт', 6, '2025-04-12'),
(5, 5007, 5, 'D015', 'шт', 2, '2025-05-08'),
(1, 1008, 12, 'D009', 'шт', 80, '2025-06-20'),
(2, 2006, 1, 'D010', 'шт', 300, '2025-07-15'),
(3, 3006, 2, 'D012', 'шт', 25, '2025-08-28'),
(4, 4005, 3, 'D013', 'шт', 10, '2025-09-10'),
(5, 5008, 8, 'D003', 'шт', 4, '2025-10-05'),
(1, 1009, 11, 'D004', 'шт', 2, '2025-11-18'),
(2, 2007, 10, 'D006', 'шт', 7, '2025-12-02');

COMMIT;

