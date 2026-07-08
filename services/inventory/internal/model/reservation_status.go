package model

type ReservationStatus string

const (
	StatusReserved  ReservationStatus = "RESERVED"
	StatusConfirmed ReservationStatus = "CONFIRMED"
	StatusReleased  ReservationStatus = "RELEASED"
	StatusExpired   ReservationStatus = "EXPIRED"
)