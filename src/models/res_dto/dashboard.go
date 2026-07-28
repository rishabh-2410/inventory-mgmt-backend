package models

type DashboardData struct {
	CurrentInventoryValue float64
	LowStockProducts      int
	OutOfStockProducts    int
}

type InventoryStats struct {
	UnitsSold     float64
	UnitsReceived float64
	UnitsDamaged  float64
	UnitsReturned float64
}
