package models

type DashboardData struct {
	CurrentInventory   float64 `json:"current_inventory"`
	LowStockProducts   int     `json:"low_stock_products"`
	OutOfStockProducts int     `json:"out_of_stock_products"`
}

type InventoryStats struct {
	UnitsSold     float64 `json:"units_sold"`
	UnitsReceived float64 `json:"units_received"`
	UnitsDamaged  float64 `json:"units_damaged"`
	UnitsReturned float64 `json:"units_returned"`
}
