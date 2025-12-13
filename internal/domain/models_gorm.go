package domain

import "time"

// GORM модели для ORM подхода в Task1

// CustomerGorm представляет покупателя для GORM
type CustomerGorm struct {
	CustomerID int    `gorm:"primaryKey;column:customer_id"`
	Name       string `gorm:"column:name"`
	Address    string `gorm:"column:address"`
	City       string `gorm:"column:city"`
}

// TableName возвращает имя таблицы для GORM
func (CustomerGorm) TableName() string {
	return "customers"
}

// PartGorm представляет деталь для GORM
type PartGorm struct {
	PartCode  string  `gorm:"primaryKey;column:part_code"`
	PartType  string  `gorm:"column:part_type"`
	Name      string  `gorm:"column:name"`
	Unit      string  `gorm:"column:unit"`
	PlanPrice float64 `gorm:"column:plan_price"`
}

// TableName возвращает имя таблицы для GORM
func (PartGorm) TableName() string {
	return "parts"
}

// ShipmentGorm представляет отгрузку для GORM с загрузкой связей
type ShipmentGorm struct {
	WarehouseNo   int       `gorm:"primaryKey;column:warehouse_no"`
	ShipmentDocNo int       `gorm:"primaryKey;column:shipment_doc_no"`
	CustomerID    int       `gorm:"column:customer_id"`
	PartCode      string    `gorm:"column:part_code"`
	Unit          string    `gorm:"column:unit"`
	Qty           float64   `gorm:"column:qty"`
	ShipmentDate  time.Time `gorm:"column:shipment_date"`

	// Связи GORM - автоматически загружаются через Preload
	Customer CustomerGorm `gorm:"foreignKey:CustomerID;references:CustomerID"`
	Part     PartGorm     `gorm:"foreignKey:PartCode;references:PartCode"`
}

// TableName возвращает имя таблицы для GORM
func (ShipmentGorm) TableName() string {
	return "shipments"
}

