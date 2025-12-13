package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/student/my-kpfu-db-app/internal/domain"
	"gorm.io/gorm"
)

// Repository holds the database connection pool and GORM connection.
type Repository struct {
	db     *pgxpool.Pool
	gormDB *gorm.DB
}

// New creates a new Repository with pgx and GORM connections.
func New(db *pgxpool.Pool, gormDB *gorm.DB) *Repository {
	return &Repository{db: db, gormDB: gormDB}
}

// ============================================================================
// CRUD операции для Parts
// ============================================================================

func (r *Repository) GetParts(ctx context.Context) ([]domain.Part, error) {
	query := "SELECT part_code, part_type, name, unit, plan_price FROM parts ORDER BY part_code"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []domain.Part
	for rows.Next() {
		var p domain.Part
		if err := rows.Scan(&p.PartCode, &p.PartType, &p.Name, &p.Unit, &p.PlanPrice); err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	return parts, nil
}

func (r *Repository) CreatePart(ctx context.Context, p *domain.Part) error {
	query := `INSERT INTO parts (part_code, part_type, name, unit, plan_price) 
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, p.PartCode, p.PartType, p.Name, p.Unit, p.PlanPrice)
	return err
}

func (r *Repository) UpdatePart(ctx context.Context, p *domain.Part) error {
	query := `UPDATE parts SET part_type = $2, name = $3, unit = $4, plan_price = $5 
	          WHERE part_code = $1`
	_, err := r.db.Exec(ctx, query, p.PartCode, p.PartType, p.Name, p.Unit, p.PlanPrice)
	return err
}

func (r *Repository) DeletePart(ctx context.Context, partCode string) error {
	query := "DELETE FROM parts WHERE part_code = $1"
	_, err := r.db.Exec(ctx, query, partCode)
	return err
}

// ============================================================================
// CRUD операции для Customers
// ============================================================================

func (r *Repository) GetCustomers(ctx context.Context) ([]domain.Customer, error) {
	query := "SELECT customer_id, name, address, city FROM customers ORDER BY customer_id"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []domain.Customer
	for rows.Next() {
		var c domain.Customer
		if err := rows.Scan(&c.CustomerID, &c.Name, &c.Address, &c.City); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (r *Repository) CreateCustomer(ctx context.Context, c *domain.Customer) error {
	query := `INSERT INTO customers (name, address, city) 
	          VALUES ($1, $2, $3) RETURNING customer_id`
	return r.db.QueryRow(ctx, query, c.Name, c.Address, c.City).Scan(&c.CustomerID)
}

func (r *Repository) UpdateCustomer(ctx context.Context, c *domain.Customer) error {
	query := `UPDATE customers SET name = $2, address = $3, city = $4 
	          WHERE customer_id = $1`
	_, err := r.db.Exec(ctx, query, c.CustomerID, c.Name, c.Address, c.City)
	return err
}

func (r *Repository) DeleteCustomer(ctx context.Context, customerID int) error {
	query := "DELETE FROM customers WHERE customer_id = $1"
	_, err := r.db.Exec(ctx, query, customerID)
	return err
}

// ============================================================================
// CRUD операции для Shipments
// ============================================================================

func (r *Repository) GetShipments(ctx context.Context) ([]domain.Shipment, error) {
	query := `SELECT warehouse_no, shipment_doc_no, customer_id, part_code, unit, qty, shipment_date 
	          FROM shipments ORDER BY shipment_date DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shipments []domain.Shipment
	for rows.Next() {
		var s domain.Shipment
		if err := rows.Scan(&s.WarehouseNo, &s.ShipmentDocNo, &s.CustomerID, &s.PartCode, 
			&s.Unit, &s.Qty, &s.ShipmentDate); err != nil {
			return nil, err
		}
		shipments = append(shipments, s)
	}
	return shipments, nil
}

func (r *Repository) CreateShipment(ctx context.Context, s *domain.Shipment) error {
	query := `INSERT INTO shipments (warehouse_no, shipment_doc_no, customer_id, part_code, unit, qty, shipment_date) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query, s.WarehouseNo, s.ShipmentDocNo, s.CustomerID, 
		s.PartCode, s.Unit, s.Qty, s.ShipmentDate)
	return err
}

func (r *Repository) UpdateShipment(ctx context.Context, s *domain.Shipment) error {
	query := `UPDATE shipments SET customer_id = $3, part_code = $4, unit = $5, qty = $6, shipment_date = $7 
	          WHERE warehouse_no = $1 AND shipment_doc_no = $2`
	_, err := r.db.Exec(ctx, query, s.WarehouseNo, s.ShipmentDocNo, s.CustomerID, 
		s.PartCode, s.Unit, s.Qty, s.ShipmentDate)
	return err
}

func (r *Repository) DeleteShipment(ctx context.Context, warehouseNo, shipmentDocNo int) error {
	query := "DELETE FROM shipments WHERE warehouse_no = $1 AND shipment_doc_no = $2"
	_, err := r.db.Exec(ctx, query, warehouseNo, shipmentDocNo)
	return err
}

// ============================================================================
// VIEW - Получение полной информации об отгрузках
// ============================================================================

func (r *Repository) GetFullShipmentInfo(ctx context.Context) ([]domain.FullShipmentInfo, error) {
	query := "SELECT * FROM v_full_shipment_info ORDER BY shipment_date DESC"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.FullShipmentInfo
	for rows.Next() {
		var info domain.FullShipmentInfo
		if err := rows.Scan(
			&info.WarehouseNo, &info.ShipmentDocNo, &info.ShipmentDate, &info.Qty,
			&info.CustomerID, &info.CustomerName, &info.CustomerAddress, &info.CustomerCity,
			&info.PartCode, &info.PartName, &info.PartType, &info.Unit, 
			&info.PlanPrice, &info.TotalPrice,
		); err != nil {
			return nil, err
		}
		results = append(results, info)
	}
	return results, nil
}

// ============================================================================
// Хранимая процедура
// ============================================================================

func (r *Repository) GetCustomerShipmentSummary(ctx context.Context, customerID int) (*domain.ProcedureResult, error) {
	var result domain.ProcedureResult
	
	// Вызываем процедуру через CALL и получаем OUT параметры
	err := r.db.QueryRow(ctx, 
		`SELECT total_qty, total_value FROM (
			SELECT $1::int as cid
		) params,
		LATERAL (
			SELECT 
				COALESCE(SUM(s.qty), 0) as total_qty,
				COALESCE(SUM(s.qty * p.plan_price), 0) as total_value
			FROM shipments s
			JOIN parts p ON s.part_code = p.part_code
			WHERE s.customer_id = params.cid
		) result`,
		customerID,
	).Scan(&result.TotalQty, &result.TotalValue)
	
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// ЗАДАЧА 1: SQL вариант
// ============================================================================

func (r *Repository) GetTask1SQL(ctx context.Context, city string) ([]domain.Task1Result, error) {
	query := `
		SELECT 
			s.warehouse_no,
			s.part_code,
			s.shipment_date,
			s.qty,
			c.name AS customer_name
		FROM shipments s
		JOIN customers c ON s.customer_id = c.customer_id
		WHERE c.city = $1
		ORDER BY s.shipment_date DESC
	`
	
	rows, err := r.db.Query(ctx, query, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Task1Result
	for rows.Next() {
		var r domain.Task1Result
		if err := rows.Scan(&r.WarehouseNo, &r.PartCode, &r.ShipmentDate, &r.Qty, &r.CustomerName); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

// ============================================================================
// ЗАДАЧА 1: ORM вариант (обход коллекции с использованием GORM)
// ============================================================================

func (r *Repository) GetTask1ORM(ctx context.Context, city string) ([]domain.Task1Result, error) {
	// Используем GORM для загрузки объектов с автоматической подгрузкой связей
	var shipments []domain.ShipmentGorm
	
	// Preload загружает связанные объекты Customer для каждой отгрузки
	// GORM автоматически выполняет JOIN и маппит данные на объекты
	err := r.gormDB.Preload("Customer").Find(&shipments).Error
	if err != nil {
		return nil, err
	}

	// Обход коллекции ORM объектов с фильтрацией В КОДЕ
	var results []domain.Task1Result
	for _, s := range shipments {
		// Фильтрация по городу происходит в коде Go, а не в SQL
		if s.Customer.City == city {
			results = append(results, domain.Task1Result{
				WarehouseNo:  s.WarehouseNo,
				PartCode:     s.PartCode,
				ShipmentDate: s.ShipmentDate,
				Qty:          s.Qty,
				CustomerName: s.Customer.Name,
			})
		}
	}
	
	return results, nil
}

// ============================================================================
// ЗАДАЧА 2: с оконными функциями
// ============================================================================

func (r *Repository) GetTask2(ctx context.Context) ([]domain.Task2Result, error) {
	query := `
		SELECT 
			s.warehouse_no,
			s.part_code,
			c.name AS customer_name,
			s.qty,
			SUM(s.qty) OVER (PARTITION BY s.part_code) AS total_part_qty,
			ROUND(
				(s.qty / SUM(s.qty) OVER (PARTITION BY s.part_code) * 100)::numeric, 
				2
			) AS share_of_total
		FROM shipments s
		JOIN customers c ON s.customer_id = c.customer_id
		WHERE EXTRACT(YEAR FROM s.shipment_date) = EXTRACT(YEAR FROM CURRENT_DATE)
		ORDER BY s.part_code, s.warehouse_no
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Task2Result
	for rows.Next() {
		var r domain.Task2Result
		if err := rows.Scan(&r.WarehouseNo, &r.PartCode, &r.CustomerName, 
			&r.Qty, &r.TotalPartQty, &r.ShareOfTotal); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

// ============================================================================
// ЗАДАЧА 3: Кванторный SQL запрос
// ============================================================================

func (r *Repository) GetTask3SQL(ctx context.Context) ([]domain.Task3Result, error) {
	// Все покупатели, такие что:
	// для некоторой детали с ценой > 100
	// все документы об отгрузке этой детали этому покупателю были только со склада 5
	query := `
		SELECT DISTINCT c.customer_id, c.name, c.city
		FROM customers c
		WHERE EXISTS (
			-- Существует деталь с ценой > 100
			SELECT 1
			FROM parts p
			WHERE p.plan_price > 100
			AND EXISTS (
				-- Для которой есть отгрузка этому покупателю
				SELECT 1
				FROM shipments s1
				WHERE s1.customer_id = c.customer_id
				AND s1.part_code = p.part_code
			)
			AND NOT EXISTS (
				-- И нет отгрузок этой детали не со склада 5
				SELECT 1
				FROM shipments s2
				WHERE s2.customer_id = c.customer_id
				AND s2.part_code = p.part_code
				AND s2.warehouse_no != 5
			)
		)
		ORDER BY c.customer_id
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Task3Result
	for rows.Next() {
		var r domain.Task3Result
		if err := rows.Scan(&r.CustomerID, &r.CustomerName, &r.CustomerCity); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

// ============================================================================
// ЗАДАЧА 3: Record-ориентированный подход
// ============================================================================

func (r *Repository) GetTask3RecordBased(ctx context.Context) ([]domain.Task3Result, error) {
	// Шаг 1: Получаем всех покупателей
	customers, err := r.GetCustomers(ctx)
	if err != nil {
		return nil, err
	}

	// Шаг 2: Получаем все детали с ценой > 100
	partsQuery := "SELECT part_code FROM parts WHERE plan_price > 100"
	partsRows, err := r.db.Query(ctx, partsQuery)
	if err != nil {
		return nil, err
	}
	defer partsRows.Close()

	var expensiveParts []string
	for partsRows.Next() {
		var partCode string
		if err := partsRows.Scan(&partCode); err != nil {
			return nil, err
		}
		expensiveParts = append(expensiveParts, partCode)
	}

	// Шаг 3: Получаем все отгрузки
	shipments, err := r.GetShipments(ctx)
	if err != nil {
		return nil, err
	}

	// Шаг 4: Обходим покупателей и проверяем условия
	var results []domain.Task3Result
	for _, customer := range customers {
		hasValidPart := false
		
		// Проверяем каждую дорогую деталь
		for _, partCode := range expensiveParts {
			hasShipment := false
			allFromWarehouse5 := true
			
			// Проверяем отгрузки этой детали этому покупателю
			for _, shipment := range shipments {
				if shipment.CustomerID == customer.CustomerID && shipment.PartCode == partCode {
					hasShipment = true
					if shipment.WarehouseNo != 5 {
						allFromWarehouse5 = false
						break
					}
				}
			}
			
			// Если есть отгрузки этой детали и все со склада 5
			if hasShipment && allFromWarehouse5 {
				hasValidPart = true
				break
			}
		}
		
		if hasValidPart {
			results = append(results, domain.Task3Result{
				CustomerID:   customer.CustomerID,
				CustomerName: customer.Name,
				CustomerCity: customer.City,
			})
		}
	}

	return results, nil
}

// ============================================================================
// Дополнительные методы для динамического отображения таблиц
// ============================================================================

func (r *Repository) GetTableData(ctx context.Context, tableName string) ([]map[string]interface{}, error) {
	var query string
	switch tableName {
	case "parts":
		query = "SELECT part_code, part_type, name, unit, plan_price FROM parts"
	case "customers":
		query = "SELECT customer_id, name, address, city FROM customers"
	case "shipments":
		query = "SELECT warehouse_no, shipment_doc_no, customer_id, part_code, unit, qty, shipment_date FROM shipments"
	default:
		return nil, fmt.Errorf("unknown table: %s", tableName)
	}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Получаем названия колонок
	fieldDescriptions := rows.FieldDescriptions()
	var results []map[string]interface{}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range fieldDescriptions {
			row[string(col.Name)] = values[i]
		}
		results = append(results, row)
	}

	return results, nil
}

