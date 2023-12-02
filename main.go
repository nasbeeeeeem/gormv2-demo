package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID         uint        `gorm:"primaryKey"`
	Name       string      `gorm:"not null"`
	Email      string      `gorm:"not null;unique"`
	GroupUsers []GroupUser `gorm:"foreignKey:UserID;references:ID"`
	Events     []Event     `gorm:"foreignKey:CreatedBy;references:ID"`
}

type Group struct {
	ID         uint        `gorm:"primaryKey"`
	Name       string      `gorm:"not null"`
	GroupUsers []GroupUser `gorm:"foreignKey:GroupID;references:ID"`
	Events     []Event     `gorm:"foreignKey:GroupID;references:ID"`
}

type GroupUser struct {
	UserID  uint `gorm:"primaryKey"`
	GroupID uint `gorm:"primaryKey"`
	// EventsUserID  []Event `gorm:"foreignKey:CreatedBy;references:ID"`
	// EventsGroupID []Event `gorm:"foreignKey:GroupID;references:ID"`
}

type Event struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	CreatedBy uint
	GroupID   uint
}

func main() {
	// dbのコネクション
	dsn := "host=localhost user=postgres password=root dbname=gorm_db port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get db instance")
	}

	defer func() {
		closeErr := sqlDB.Close()
		if closeErr != nil {
			panic("failed to close connection")
		}
	}()

	// dbのマイグレーション
	db.AutoMigrate(&User{}, &Group{}, &GroupUser{}, &Event{})
	// db.Migrator().CreateConstraint(&GroupUser{}, "User")
	// db.Migrator().CreateConstraint(&GroupUser{}, "Group")
	// db.Migrator().CreateConstraint(&Event{}, "GroupUser")
	fmt.Print("migrate successfully")
}
