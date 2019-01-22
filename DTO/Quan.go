package DTO

import "time"

type QuanDTO struct {
	Id int
	Name string
	CategoryId int
	Gender string
	Amount int
	Price int
	Create time.Month
	Image string
}