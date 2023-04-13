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
	PageId    int    `gorm:"primaryKey;autoIncrement:false"`
	Przypadek string `gorm:"primaryKey"`
	Liczba    string
	Text      string
}

// TableName overrides the table name used by User to `profiles`
func (Odmiana) TableName() string {
	return "odmiany"
}
