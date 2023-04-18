package data

type Tabler interface {
	TableName() string
}

type Text struct {
	OldId   uint `gorm:"primaryKey"`
	OldText string
}

// TableName overrides the table name used by User to `profiles`
func (Text) TableName() string {
	return "text"
}

type Odmiana struct {
	Guid      string `gorm:"primaryKey"`
	PageId    int
	Przypadek string
	Liczba    string
	Text      string
}

// TableName overrides the table name used by User to `profiles`
func (Odmiana) TableName() string {
	return "odmiany"
}

type ParseTypes uint8

const (
	NotParsed uint8 = 0
	Parsed    uint8 = 1
	Error     uint8 = 2
)

type Page struct {
	PageId int `gorm:"primaryKey;autoIncrement:false"`
	Parsed ParseTypes
}

// TableName overrides the table name used by User to `profiles`
func (Page) TableName() string {
	return "page"
}
