package entity

import "time"

type Employee struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Department  string `json:"department"`
}

type Customer struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}

type Product struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Unit  string `json:"unit"`
}

type BillDetails struct {
	Id           int    `json:"id"`
	BillId       string `json:"billId"`
	Product      int    `json:"product"`
	ProductPrice int    `json:"productPrice"`
	Qty          int    `json:"qty"`
}

type Bill struct {
	Id         int       `json:"id"`
	BillId     string    `json:"billData"`
	EntryDate  time.Time `json:"entryDate"`
	FinishDate time.Time `json:"finishDate"`
	Employee   int       `json:"employee"`
	Customer   int       `json:"customer"`
	TotalBill  int       `json:"totalBill"`
}
