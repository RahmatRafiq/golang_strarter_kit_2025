package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"
	"unicode"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var MakeFactoryCommand = &cli.Command{
	Name:  "make:factory",
	Usage: "Generate a new model factory",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Factory name (e.g., User, Product)",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		name := c.String("name")
		if name == "" {
			return fmt.Errorf("factory name is required. Example: make:factory --name=User")
		}

		// Generate factory file
		fileName := fmt.Sprintf("%s_factory.go", strings.ToLower(name))
		dir := path.Join("app", "database", "factories")

		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create factories directory: %w", err)
		}

		filePath := path.Join(dir, fileName)

		// Check if file already exists
		if _, err := os.Stat(filePath); err == nil {
			return fmt.Errorf("factory already exists: %s", filePath)
		}

		// Capitalize first letter
		runes := []rune(name)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		structName := string(runes)
		content := fmt.Sprintf(`package factories

import (
	"fmt"
	"time"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// %[1]sFactory creates %[1]s instances for testing and seeding
type %[1]sFactory struct {
	*BaseFactory
}

// New%[1]sFactory creates a new %[1]sFactory
func New%[1]sFactory(db *gorm.DB) *%[1]sFactory {
	return &%[1]sFactory{
		BaseFactory: NewBaseFactory(db),
	}
}

// Make creates a %[1]s instance without saving to database
func (f *%[1]sFactory) Make(overrides ...map[string]interface{}) *models.%[1]s {
	reference, err := helpers.GenerateReference("%[2]s")
	if err != nil {
		// Factory fallback: use timestamp-based reference if UUID fails
		reference = fmt.Sprintf("%[2]s-FALLBACK-%%d", time.Now().UnixNano())
		log.Warn().Err(err).Msg("Failed to generate reference in factory, using fallback")
	}

	item := &models.%[1]s{
		Reference: reference,
		// TODO: Add your model fields here
		// Name:      RandomName(),
		// Email:     RandomEmail(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Apply overrides if provided
	if len(overrides) > 0 {
		override := overrides[0]
		// TODO: Apply overrides for your fields
		// if val, ok := override["name"]; ok {
		//     item.Name = val.(string)
		// }
	}

	return item
}

// MakeMany creates multiple %[1]s instances without saving
func (f *%[1]sFactory) MakeMany(count int, overrides ...map[string]interface{}) []*models.%[1]s {
	items := make([]*models.%[1]s, count)
	for i := 0; i < count; i++ {
		items[i] = f.Make(overrides...)
	}
	return items
}

// Create creates and saves a %[1]s to the database
func (f *%[1]sFactory) Create(overrides ...map[string]interface{}) (*models.%[1]s, error) {
	item := f.Make(overrides...)
	if err := f.DB.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// CreateMany creates and saves multiple %[1]s to the database
func (f *%[1]sFactory) CreateMany(count int, overrides ...map[string]interface{}) ([]*models.%[1]s, error) {
	items := f.MakeMany(count, overrides...)
	if err := f.DB.Create(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CreateInBatches creates multiple %[1]s in batches for better performance
func (f *%[1]sFactory) CreateInBatches(count, batchSize int, overrides ...map[string]interface{}) ([]*models.%[1]s, error) {
	items := f.MakeMany(count, overrides...)
	if err := f.DB.CreateInBatches(&items, batchSize).Error; err != nil {
		return nil, err
	}
	return items, nil
}
`, structName, strings.ToUpper(name[:3]))

		if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
			return fmt.Errorf("failed to create factory file: %w", err)
		}

		log.Info().
			Str("file", filePath).
			Msg("Factory created successfully")
		log.Info().
			Msg("Next steps: 1) Update the Make() method with your model fields")
		log.Info().
			Msg("Next steps: 2) Add override handling for your fields")
		log.Info().
			Str("usage", "factory := factories.New"+structName+"Factory(db)").
			Msg("Next steps: 3) Use in seeders")

		return nil
	},
}
