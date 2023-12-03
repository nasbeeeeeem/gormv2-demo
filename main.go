package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID             uint        `gorm:"primaryKey"`
	Name           string      `gorm:"not null"`
	Email          string      `gorm:"not null;unique"`
	GroupUsers     []GroupUser `gorm:"foreignKey:UserID"`
	PaymentsPaidBy []Payment   `gorm:"foreignKey:PaidBy"`
	PaymentsPaidTo []Payment   `gorm:"foreignKey:PaidTo"`
}

type Group struct {
	ID         uint        `gorm:"primaryKey"`
	Name       string      `gorm:"not null"`
	GroupUsers []GroupUser `gorm:"foreignKey:GroupID"`
}

type GroupUser struct {
	UserID        uint    `gorm:"unique"`
	GroupID       uint    `gorm:"unique"`
	EventsUserID  []Event `gorm:"foreignKey:CreatedBy;references:UserID"` // CreatedBy列の外部キー
	EventsGroupID []Event `gorm:"foreignKey:GroupID;references:GroupID"`  // GroupID列の外部キー
}

type Event struct {
	ID              uint   `gorm:"primaryKey"`
	Name            string `gorm:"not null"`
	CreatedBy       uint
	GroupID         uint
	PaymentsEventID []Payment `gorm:"foreignKey:EventID"`
}

type Payment struct {
	gorm.Model
	EventID uint
	PaidBy  uint
	PaidTo  uint
	Amount  uint
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
	db.AutoMigrate(&User{}, &Group{}, &GroupUser{}, &Event{}, &Payment{})
	// db.Migrator().CreateConstraint(&GroupUser{}, "User")
	// db.Migrator().CreateConstraint(&GroupUser{}, "Group")
	// db.Migrator().CreateConstraint(&Event{}, "GroupUser")

	// サンプルデータの登録
	user := User{
		Name:  "yakiu",
		Email: "yakiu@gmail.com",
	}
	db.Create(&user)

	group := Group{
		Name: "SampleGroup",
	}
	db.Create(&group)

	groupUser := GroupUser{
		UserID:  user.ID,
		GroupID: group.ID,
	}
	db.Create(&groupUser)

	event := Event{
		Name:      "SampleEvent",
		CreatedBy: groupUser.UserID,
		GroupID:   groupUser.GroupID,
	}
	db.Create(&event)

	payment := Payment{
		EventID: event.ID,
		PaidBy:  user.ID,
		PaidTo:  user.ID,
		Amount:  2000,
	}
	db.Create(&payment)

	fmt.Print("Sample data created successfully")
}
