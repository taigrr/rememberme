package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	answerState sessionState = iota
	questionState
)

type QNA struct {
	Question string
	Answer   string
}

var qna = []QNA{{Question: "what's 1 + 1", Answer: "2"}, {Question: "is red a warm colour?", Answer: "yes"}, {Question: "what is the best snack", Answer: "popcorn"}}
var (
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	boxstyle  = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(highlight).
			Padding(0, 1)
)

type model struct {
	state    sessionState
	viewport viewport.Model
	question string
	answer   string
}

type QNAMsg string

func (m model) getQNACmd() tea.Msg {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return QNAMsg("hello")
}

func New() model {
	m := model{state: questionState, question: "a question", answer: "an answer"}
	m.viewport = viewport.New(8, 8)
	m.viewport.SetContent(m.question)
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			// refresh new question and answer
			i := rand.Int()
			i = i % len(qna)
			m.question = qna[i].Question
			m.answer = qna[i].Answer
			m.state = questionState
			m.viewport.SetContent(m.question)
		case "enter":
			if m.state == questionState {
				m.state = answerState
				m.viewport.SetContent(m.answer)
			} else {
				m.state = questionState
				m.viewport.SetContent(m.question)
			}
		case "ctrl+c":
			fallthrough
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return boxstyle.Render(m.viewport.View())
}

func main() {
	p := tea.NewProgram(
		New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if err := p.Start(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
