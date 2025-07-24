package flashcards

type RequestCreateFlashcard struct {
	Question string `json:"question" validate:"required,min=1,max=225"`
	Answer   string `json:"answer" validate:"required,min=1,max=225"`
}
