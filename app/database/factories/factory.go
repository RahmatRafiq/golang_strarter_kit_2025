package factories

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/gorm"
)

// Factory is the base factory interface
type Factory interface {
	Create(overrides ...map[string]interface{}) error
	CreateMany(count int, overrides ...map[string]interface{}) error
	Make(overrides ...map[string]interface{}) interface{}
	MakeMany(count int, overrides ...map[string]interface{}) []interface{}
}

// BaseFactory provides common factory functionality
type BaseFactory struct {
	DB *gorm.DB
}

// NewBaseFactory creates a new base factory
func NewBaseFactory(db *gorm.DB) *BaseFactory {
	return &BaseFactory{DB: db}
}

// InitFaker initializes the faker with a seed for reproducible data
func InitFaker(seed int64) {
	if seed > 0 {
		gofakeit.Seed(seed)
	} else {
		gofakeit.Seed(0) // Random seed
	}
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int) int {
	return gofakeit.Number(min, max)
}

// RandomString generates a random string of specified length
func RandomString(length int) string {
	return gofakeit.LetterN(uint(length))
}

// RandomEmail generates a random email address
func RandomEmail() string {
	return gofakeit.Email()
}

// RandomUsername generates a random username
func RandomUsername() string {
	return gofakeit.Username()
}

// RandomName generates a random full name
func RandomName() string {
	return gofakeit.Name()
}

// RandomPhone generates a random phone number
func RandomPhone() string {
	return gofakeit.Phone()
}

// RandomAddress generates a random address
func RandomAddress() string {
	return gofakeit.Address().Address
}

// RandomURL generates a random URL
func RandomURL() string {
	return gofakeit.URL()
}

// RandomParagraph generates a random paragraph
func RandomParagraph() string {
	return gofakeit.Paragraph(3, 5, 10, " ")
}

// RandomSentence generates a random sentence
func RandomSentence() string {
	return gofakeit.Sentence(10)
}

// RandomDate generates a random date
func RandomDate() string {
	return gofakeit.Date().Format("2006-01-02")
}

// RandomDateTime generates a random datetime
func RandomDateTime() string {
	return gofakeit.Date().Format("2006-01-02 15:04:05")
}

// RandomBool generates a random boolean
func RandomBool() bool {
	return gofakeit.Bool()
}

// RandomPrice generates a random price between min and max
func RandomPrice(min, max float64) float64 {
	return gofakeit.Float64Range(min, max)
}

// RandomFromSlice picks a random element from a slice
func RandomFromSlice(slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	return slice[gofakeit.Number(0, len(slice)-1)]
}

// Sequence generates a sequence number (useful for unique fields)
var sequenceCounter = make(map[string]int)

func Sequence(key string) int {
	sequenceCounter[key]++
	return sequenceCounter[key]
}

// ResetSequence resets a sequence counter
func ResetSequence(key string) {
	sequenceCounter[key] = 0
}

// ResetAllSequences resets all sequence counters
func ResetAllSequences() {
	sequenceCounter = make(map[string]int)
}

// Times repeats a function n times and returns the results
func Times(n int, fn func(i int) interface{}) []interface{} {
	results := make([]interface{}, n)
	for i := 0; i < n; i++ {
		results[i] = fn(i)
	}
	return results
}

// BatchCreate creates multiple records in batches
func BatchCreate(db *gorm.DB, records interface{}, batchSize int) error {
	return db.CreateInBatches(records, batchSize).Error
}

// WithContext wraps a factory function with context for relationships
type FactoryContext struct {
	Data map[string]interface{}
}

// NewFactoryContext creates a new factory context
func NewFactoryContext() *FactoryContext {
	return &FactoryContext{
		Data: make(map[string]interface{}),
	}
}

// Set sets a value in the context
func (fc *FactoryContext) Set(key string, value interface{}) {
	fc.Data[key] = value
}

// Get gets a value from the context
func (fc *FactoryContext) Get(key string) (interface{}, bool) {
	val, ok := fc.Data[key]
	return val, ok
}

// MustGet gets a value from the context or panics
func (fc *FactoryContext) MustGet(key string) interface{} {
	val, ok := fc.Data[key]
	if !ok {
		panic(fmt.Sprintf("key '%s' not found in factory context", key))
	}
	return val
}
