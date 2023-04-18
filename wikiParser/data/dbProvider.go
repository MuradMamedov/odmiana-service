package data

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"wiki-parser/utilities"

	"github.com/google/uuid"

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
		// get next batch
		batch, err := utilities.StopWatch(d.logger, d.getNextPageIDsBatch, "")
		if err != nil {
			d.logger.Printf("ParsePages. Failed to read the batch: %v\n", err)
			break
		}
		ids := mapIDsToInts(batch)

		// read  text table
		texts, err := utilities.StopWatchParametrized(d.logger, d.getTexts, ids, "getTexts")
		if err != nil {
			d.logger.Printf("ParsePages. Failed to read the texts: %v\n", err)
			break
		}

		// write batch
		_, err = utilities.StopWatchParametrized(d.logger, d.writeTexts, texts, "writeTexts")
		if err != nil {
			d.logger.Printf("ParsePages. Failed to write the texts: %v\n", err)
			break
		} else {
			// mark pages as parsed
			_, err := utilities.StopWatchParametrized(d.logger, d.markPagesAsParsed, ids, "markPagesAsParsed")
			if err != nil {
				d.logger.Printf("ParsePages. Failed to mark pages as parsed: %v\n", err)
				break
			}
		}

		if len(batch) == 0 {
			break
		}
	}
	return nil
}

func (d *DbProvider) getNextPageIDsBatch() ([]Page, error) {
	var pages []Page
	err := d.connection.Where("parsed = 0").Limit(d.batchSize).Find(&pages).Error
	if err != nil {
		return nil, err
	}

	return pages, nil
}

func (d *DbProvider) getTexts(ids []int) ([]Text, error) {
	var texts []Text
	err := d.connection.Where("old_id in (?)", ids).Find(&texts).Error
	if err != nil {
		return nil, err
	}

	return texts, nil
}

func (d *DbProvider) markPagesAsParsed(ids []int) ([]int, error) {
	err := d.connection.Where("page_id in (?)", ids).Updates(Page{Parsed: ParseTypes(Parsed)}).Error
	if err != nil {
		return ids, err
	}

	return ids, nil
}

func mapIDsToInts(pages []Page) []int {
	ids := make([]int, len(pages))
	for i, obj := range pages {
		ids[i] = int(obj.PageId)
	}
	return ids
}

func (d *DbProvider) writeTexts(texts []Text) ([]Odmiana, error) {
	r := regexp.MustCompile(`\|(\S+)\s((?:\S)+)\s=\s([^|}\\]+)`)
	var results []Odmiana
	for _, text := range texts {
		matches := r.FindAllStringSubmatch(text.OldText, -1) // matches is [][]string
		i := 14
		if len(matches) < 7 {
			continue
		}

		for _, match := range matches {
			if i == 0 {
				continue
			}

			min := int(math.Min(float64(len(match[3])), 100))
			odmiana := Odmiana{
				Guid:      uuid.New().String(),
				PageId:    int(text.OldId),
				Przypadek: match[1],
				Liczba:    match[2],
				Text:      match[3][:min],
			}
			results = append(results, odmiana)
			i--
		}
	}

	err := d.connection.CreateInBatches(&results, len(results)).Error
	if err != nil {
		return results, err
	}
	return results, nil
}
