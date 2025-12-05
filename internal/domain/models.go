package domain

import "time"

// Part represents a part/detail in the database.
type Part struct {
	PartCode  string  `json:"part_code"`
	PartType  string  `json:"part_type"`
	Name      string  `json:"name"`
	Unit      string  `json:"unit"`
	PlanPrice float64 `json:"plan_price"`
}

// Customer represents a customer in the database.
type Customer struct {
	CustomerID int    `json:"customer_id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	City       string `json:"city"`
}

// Shipment represents a shipment record in the database.
type Shipment struct {
	WarehouseNo   int       `json:"warehouse_no"`
	ShipmentDocNo int       `json:"shipment_doc_no"`
	CustomerID    int       `json:"customer_id"`
	PartCode      string    `json:"part_code"`
	Unit          string    `json:"unit"`
	Qty           float64   `json:"qty"`
	ShipmentDate  time.Time `json:"shipment_date"`
}

// FullShipmentInfo represents the VIEW combining all three tables.
type FullShipmentInfo struct {
	WarehouseNo     int       `json:"warehouse_no"`
	ShipmentDocNo   int       `json:"shipment_doc_no"`
	ShipmentDate    time.Time `json:"shipment_date"`
	Qty             float64   `json:"qty"`
	CustomerID      int       `json:"customer_id"`
	CustomerName    string    `json:"customer_name"`
	CustomerAddress string    `json:"customer_address"`
	CustomerCity    string    `json:"customer_city"`
	PartCode        string    `json:"part_code"`
	PartName        string    `json:"part_name"`
	PartType        string    `json:"part_type"`
	Unit            string    `json:"unit"`
	PlanPrice       float64   `json:"plan_price"`
	TotalPrice      float64   `json:"total_price"`
}

// Task1Result represents the result for Task 1.
type Task1Result struct {
	WarehouseNo   int       `json:"warehouse_no"`
	PartCode      string    `json:"part_code"`
	ShipmentDate  time.Time `json:"shipment_date"`
	Qty           float64   `json:"qty"`
	CustomerName  string    `json:"customer_name"`
}

// Task2Result represents the result for Task 2 with aggregation.
type Task2Result struct {
	WarehouseNo      int     `json:"warehouse_no"`
	PartCode         string  `json:"part_code"`
	CustomerName     string  `json:"customer_name"`
	Qty              float64 `json:"qty"`
	TotalPartQty     float64 `json:"total_part_qty"`
	ShareOfTotal     float64 `json:"share_of_total"`
}

// Task3Result represents the result for Task 3.
type Task3Result struct {
	CustomerID   int    `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	CustomerCity string `json:"customer_city"`
}

// ProcedureResult represents the result of the stored procedure.
type ProcedureResult struct {
	TotalQty   float64 `json:"total_qty"`
	TotalValue float64 `json:"total_value"`
}

