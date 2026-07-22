package models

type UserRole string

const (
	RoleOwner    UserRole = "owner"
	RoleEmployee UserRole = "employee"
)
