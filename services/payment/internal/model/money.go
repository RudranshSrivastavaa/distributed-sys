package model

import "errors"



type Money struct {
    Amount   float64  `gorm:"column:amount;not null"`
    Currency string `gorm:"column:currency;size:3;not null"`
}

func NewMoney(amount float64,currency string) (Money, error) {

	if amount <= 0 {
		return Money{}, errors.New("amount must be greater than zero")
	}

	if currency == "" {
		return Money{}, errors.New("currency is required")
	}

	return Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}