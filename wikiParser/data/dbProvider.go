package data

import (
	"fmt"
	"log"
	"regexp"
	"wiki-parser/utilities"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbProvider struct {
	logger     *log.Logger
	connection *gorm.DB
	batchSize  int
}

func NewDbProvider(l *log.Logger) (*DbProvider, error) {
	connection, err := getDbConnection("root", "Test1234!")

	if err != nil {
		l.Printf("New. Failed to get DB connection: %v\n", err)
		return nil, err
	}

	return &DbProvider{l, connection, 10}, nil
}

func (d *DbProvider) Close() error {
	sqlDB, err := d.connection.DB()
	if err != nil {
		d.logger.Printf("Closing. Failed to get DB connection: %v\n", err)
	}
	err = sqlDB.Close()
	if err != nil {
		d.logger.Printf("Closing. Failed to close DB connection: %v\n", err)
		panic(err)
	}
	return nil
}

func getDbConnection(name string, pasword string) (*gorm.DB, error) {
	dsn := "root:Test1234!@tcp(127.0.0.1:3306)/wikidb?charset=utf8mb4&parseTime=True&loc=Local"
	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDB, err := connection.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connection successful!")

	return connection, nil
}

func (d *DbProvider) ParsePages() error {
	for {
		batch, err := utilities.StopWatch(d.logger, d.getNextPageIDsBatch, "")
		if err != nil {
			d.logger.Printf("ParsePages. Failed to read the batch: %v\n", err)
		}

		if len(batch.([]Page)) != 0 {
			break
		}
	}
	return nil
}

func (d *DbProvider) getNextPageIDsBatch() (interface{}, error) {
	var pages []Page
	err := d.connection.Where("parsed = 0").Limit(d.batchSize).Find(&pages).Error
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		d.logger.Println(page.PageId, " - ", page.Parsed)
	}
	return pages, nil
}

func (d *DbProvider) writeData(tag *Text) {
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
		d.connection.Create(&odmiana)
		i--
	}
}
