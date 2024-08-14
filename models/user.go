package models

type User struct {
	ID          uint   `gorm:"primaryKey;autoIncrement;index"`
	FirstName   string `gorm:"column:first_name;size:255;not null"`
	LastName    string `gorm:"column:last_name;size:255;not null"`
	CompanyName string `gorm:"column:company_name;size:255"`
	Address     string `gorm:"column:address;size:255"`
	City        string `gorm:"column:city;size:255"`
	County      string `gorm:"column:county;size:255"`
	Postal      string `gorm:"column:postal;size:20"`
	Phone       string `gorm:"column:phone;size:20"`
	Email       string `gorm:"column:email;size:255;not null"`
	Web         string `gorm:"column:web;size:255"`
}
