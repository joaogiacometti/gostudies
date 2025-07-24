package reviews

type RequestReviewFlashcard struct {
	Remembered *bool `json:"remembered" validate:"required"`
}
