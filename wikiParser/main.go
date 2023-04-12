package main

import (
	"fmt"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Text struct {
	OldId   uint `gorm:"primaryKey"`
	OldText string
}
type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by User to `profiles`
func (Text) TableName() string {
	return "text"
}

type Odmiana struct {
	PageId    int    `gorm:"primaryKey;autoIncrement:false"`
	Przypadek string `gorm:"primaryKey"`
	Liczba    string
	Text      string
}

// TableName overrides the table name used by User to `profiles`
func (Odmiana) TableName() string {
	return "odmiany"
}

func main() {
	dsn := "root:Test1234!@tcp(127.0.0.1:3306)/wikidb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}
	var tag Text
	db.First(&tag, 445)
	writeData(&tag, db)
}

func writeData(tag *Text, db *gorm.DB) {
	r := regexp.MustCompile(`\|(\S+)\s((?:\S)+)\s=\s([^|}\\]+)`)
	matches := r.FindAllStringSubmatch(tag.OldText, -1) // matches is [][]string
	i := 14
	for _, match := range matches {
		if i == 0 {
			break
		}
		fmt.Printf(
			"%s, %s, %s\n", match[1], match[2], match[3])

		odmiana := Odmiana{
			PageId:    int(tag.OldId),
			Przypadek: match[1],
			Liczba:    match[2],
			Text:      match[3],
		}
		db.Create(&odmiana)
		i--
	}
}
