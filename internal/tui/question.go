package tui

type QuestionType string

const (
	ShortQuestion       QuestionType = "short"
	LongQuestion        QuestionType = "long"
	ShortSecretQuestion QuestionType = "shortsecreat"
)

type Question struct {
	Question string
	Answer   string
	Input    Input
}

func NewQuestion(question string, placeholder string, questionType QuestionType) Question {
	return Question{Question: question, Input: newInputField(questionType, placeholder)}
}

func newInputField(questionType QuestionType, placeholder string) Input {
	switch questionType {
	case ShortQuestion:
		return NewShortAnswerField(placeholder)
	case LongQuestion:
		return NewLongAnswerField()
	case ShortSecretQuestion:
		return NewShortAnswerFieldForSecrets(placeholder)
	default:
		return nil
	}
}
