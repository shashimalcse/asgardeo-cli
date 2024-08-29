package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Blink() tea.Msg
	Blur() tea.Msg
	Focus() tea.Cmd
	SetValue(string)
	Value() string
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
}

type ShortAnswerField struct {
	textinput textinput.Model
}

func NewShortAnswerField(placeholder string) *ShortAnswerField {
	a := ShortAnswerField{}

	model := textinput.New()
	model.Placeholder = placeholder
	model.Focus()

	a.textinput = model
	return &a
}

func NewShortAnswerFieldForSecrets(placeholder string) *ShortAnswerField {
	a := ShortAnswerField{}

	model := textinput.New()
	model.Placeholder = placeholder
	model.EchoMode = textinput.EchoPassword
	model.EchoCharacter = '•'
	model.Focus()

	a.textinput = model
	return &a
}

func (a *ShortAnswerField) Blink() tea.Msg {
	return textinput.Blink()
}

func (a *ShortAnswerField) Init() tea.Cmd {
	return nil
}

func (a *ShortAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	a.textinput, cmd = a.textinput.Update(msg)
	return a, cmd
}

func (a *ShortAnswerField) View() string {
	return a.textinput.View()
}

func (a *ShortAnswerField) Focus() tea.Cmd {
	return a.textinput.Focus()
}

func (a *ShortAnswerField) SetValue(s string) {
	a.textinput.SetValue(s)
}

func (a *ShortAnswerField) Blur() tea.Msg {
	return a.textinput.Blur
}

func (a *ShortAnswerField) Value() string {
	return a.textinput.Value()
}

type LongAnswerField struct {
	textarea textarea.Model
}

func NewLongAnswerField() *LongAnswerField {
	a := LongAnswerField{}

	model := textarea.New()
	model.Placeholder = "Your answer here"
	model.Focus()

	a.textarea = model
	return &a
}

func (a *LongAnswerField) Blink() tea.Msg {
	return textarea.Blink()
}

func (a *LongAnswerField) Init() tea.Cmd {
	return nil
}

func (a *LongAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	a.textarea, cmd = a.textarea.Update(msg)
	return a, cmd
}

func (a *LongAnswerField) View() string {
	return a.textarea.View()
}

func (a *LongAnswerField) Focus() tea.Cmd {
	return a.textarea.Focus()
}

func (a *LongAnswerField) SetValue(s string) {
	a.textarea.SetValue(s)
}

func (a *LongAnswerField) Blur() tea.Msg {
	return a.textarea.Blur
}

func (a *LongAnswerField) Value() string {
	return a.textarea.Value()
}
