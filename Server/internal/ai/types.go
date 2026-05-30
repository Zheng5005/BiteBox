package ai

// AiIngredient represents a single recipe ingredient with quantity.
type AiIngredient struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

// AiStep represents an ordered cooking instruction.
type AiStep struct {
	Order       int    `json:"order"`
	Instruction string `json:"instruction"`
}

// AiRecipe is the strictly-typed output from AI generation.
type AiRecipe struct {
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	Ingredients      []AiIngredient `json:"ingredients"`
	Steps            []AiStep       `json:"steps"`
	Tags             []string       `json:"tags"`
	EstimatedMinutes int            `json:"estimated_minutes"`
	Difficulty       string         `json:"difficulty"`
	Cuisine          string         `json:"cuisine"`
}

// Validate checks that all required fields are populated.
func (r *AiRecipe) Validate() error {
	if r.Title == "" {
		return &ValidationError{Field: "title", Message: "title is required"}
	}
	if r.Description == "" {
		return &ValidationError{Field: "description", Message: "description is required"}
	}
	if len(r.Ingredients) == 0 {
		return &ValidationError{Field: "ingredients", Message: "at least one ingredient required"}
	}
	if len(r.Steps) == 0 {
		return &ValidationError{Field: "steps", Message: "at least one step required"}
	}
	if r.EstimatedMinutes <= 0 {
		return &ValidationError{Field: "estimated_minutes", Message: "must be positive"}
	}
	switch r.Difficulty {
	case "easy", "medium", "hard":
		// valid
	default:
		return &ValidationError{Field: "difficulty", Message: "must be easy, medium, or hard"}
	}
	return nil
}

// ValidationError is returned when AI output fails schema validation.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return "validation error on " + e.Field + ": " + e.Message
}
