package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/student/my-kpfu-db-app/internal/domain"
	"github.com/student/my-kpfu-db-app/internal/repository"
)

// Handler holds the repository.
type Handler struct {
	repo *repository.Repository
}

// New creates a new Handler.
func New(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// RegisterRoutes registers all routes for the application.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// Main pages
	r.GET("/", h.Home)
	r.GET("/view", h.View)
	r.GET("/dynamic", h.Dynamic)

	// Task pages
	r.GET("/task-1", h.Task1Page)
	r.GET("/task-2", h.Task2Page)
	r.GET("/task-3", h.Task3Page)

	// API endpoints for CRUD operations
	api := r.Group("/api")
	{
		// Parts
		api.POST("/parts", h.CreatePart)
		api.PUT("/parts/:code", h.UpdatePart)
		api.DELETE("/parts/:code", h.DeletePart)

		// Customers
		api.POST("/customers", h.CreateCustomer)
		api.PUT("/customers/:id", h.UpdateCustomer)
		api.DELETE("/customers/:id", h.DeleteCustomer)

		// Shipments
		api.POST("/shipments", h.CreateShipment)
		api.PUT("/shipments/:warehouse/:doc", h.UpdateShipment)
		api.DELETE("/shipments/:warehouse/:doc", h.DeleteShipment)

		// Tasks
		api.GET("/task-1/sql", h.Task1SQL)
		api.GET("/task-1/orm", h.Task1ORM)
		api.GET("/task-2", h.Task2)
		api.GET("/task-3/sql", h.Task3SQL)
		api.GET("/task-3/record", h.Task3Record)

		// Dynamic table data
		api.GET("/table/:name", h.GetTableData)

		// Procedure
		api.GET("/procedure/:customer_id", h.GetProcedureResult)
	}
}

// ============================================================================
// Main Pages
// ============================================================================

func (h *Handler) Home(c *gin.Context) {
	parts, err := h.repo.GetParts(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching parts: %v", err)
		return
	}

	customers, err := h.repo.GetCustomers(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching customers: %v", err)
		return
	}

	shipments, err := h.repo.GetShipments(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching shipments: %v", err)
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"Title":     "Главная",
		"Parts":     parts,
		"Customers": customers,
		"Shipments": shipments,
	})
}

func (h *Handler) View(c *gin.Context) {
	fullInfo, err := h.repo.GetFullShipmentInfo(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching view data: %v", err)
		return
	}

	c.HTML(http.StatusOK, "view.html", gin.H{
		"Title":    "VIEW - Полная информация об отгрузках",
		"FullInfo": fullInfo,
	})
}

func (h *Handler) Dynamic(c *gin.Context) {
	c.HTML(http.StatusOK, "dynamic.html", gin.H{
		"Title": "Динамическое отображение таблиц",
	})
}

// ============================================================================
// Task Pages
// ============================================================================

func (h *Handler) Task1Page(c *gin.Context) {
	c.HTML(http.StatusOK, "task1.html", gin.H{
		"Title": "Задача 1: Отгрузки по городу",
	})
}

func (h *Handler) Task2Page(c *gin.Context) {
	results, err := h.repo.GetTask2(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching task 2 data: %v", err)
		return
	}

	c.HTML(http.StatusOK, "task2.html", gin.H{
		"Title":   "Задача 2: Отгрузки текущего года",
		"Results": results,
	})
}

func (h *Handler) Task3Page(c *gin.Context) {
	c.HTML(http.StatusOK, "task3.html", gin.H{
		"Title": "Задача 3: Кванторный запрос",
	})
}

// ============================================================================
// CRUD API Handlers
// ============================================================================

func (h *Handler) CreatePart(c *gin.Context) {
	var part domain.Part
	if err := c.ShouldBindJSON(&part); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreatePart(c.Request.Context(), &part); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, part)
}

func (h *Handler) UpdatePart(c *gin.Context) {
	var part domain.Part
	if err := c.ShouldBindJSON(&part); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	part.PartCode = c.Param("code")

	if err := h.repo.UpdatePart(c.Request.Context(), &part); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, part)
}

func (h *Handler) DeletePart(c *gin.Context) {
	code := c.Param("code")
	if err := h.repo.DeletePart(c.Request.Context(), code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Part deleted"})
}

func (h *Handler) CreateCustomer(c *gin.Context) {
	var customer domain.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateCustomer(c.Request.Context(), &customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h *Handler) UpdateCustomer(c *gin.Context) {
	var customer domain.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	customer.CustomerID = id

	if err := h.repo.UpdateCustomer(c.Request.Context(), &customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *Handler) DeleteCustomer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.repo.DeleteCustomer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted"})
}

func (h *Handler) CreateShipment(c *gin.Context) {
	var shipment domain.Shipment
	if err := c.ShouldBindJSON(&shipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateShipment(c.Request.Context(), &shipment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shipment)
}

func (h *Handler) UpdateShipment(c *gin.Context) {
	var shipment domain.Shipment
	if err := c.ShouldBindJSON(&shipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	warehouse, _ := strconv.Atoi(c.Param("warehouse"))
	doc, _ := strconv.Atoi(c.Param("doc"))
	shipment.WarehouseNo = warehouse
	shipment.ShipmentDocNo = doc

	if err := h.repo.UpdateShipment(c.Request.Context(), &shipment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shipment)
}

func (h *Handler) DeleteShipment(c *gin.Context) {
	warehouse, _ := strconv.Atoi(c.Param("warehouse"))
	doc, _ := strconv.Atoi(c.Param("doc"))

	if err := h.repo.DeleteShipment(c.Request.Context(), warehouse, doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipment deleted"})
}

// ============================================================================
// Task API Handlers
// ============================================================================

func (h *Handler) Task1SQL(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		city = "Казань"
	}

	results, err := h.repo.GetTask1SQL(c.Request.Context(), city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) Task1ORM(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		city = "Казань"
	}

	results, err := h.repo.GetTask1ORM(c.Request.Context(), city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) Task2(c *gin.Context) {
	results, err := h.repo.GetTask2(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) Task3SQL(c *gin.Context) {
	results, err := h.repo.GetTask3SQL(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) Task3Record(c *gin.Context) {
	results, err := h.repo.GetTask3RecordBased(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) GetTableData(c *gin.Context) {
	tableName := c.Param("name")

	data, err := h.repo.GetTableData(c.Request.Context(), tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetProcedureResult(c *gin.Context) {
	customerID, err := strconv.Atoi(c.Param("customer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	result, err := h.repo.GetCustomerShipmentSummary(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
