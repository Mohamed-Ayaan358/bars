package main

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// https://www.youtube.com/watch?v=wCBwcrWqIDs

// Styling
var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

type status int

// Indices to determine the column
const (
	todo status = iota
	inProgress
	done
)

const divisor = 4

type Task struct {
	status      status
	title       string
	description string
	// loaded used so that the list content would be loaded before runtime error: index out of range [0] with length 0

}

func (t *Task) NextTask() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

func (t *Task) PrevTask() {
	if t.status == todo {
		t.status = done
	} else {
		t.status--
	}
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

type Model struct {
	focused status
	loaded  bool
	lists   []list.Model
	err     error
	qutting bool
}

func New() *Model {
	return &Model{}
}

func (m *Model) MoveNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	selectedTask := selectedItem.(Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	// this shows that the status is followed for the tasks as well, it signifies the list the task is in
	selectedTask.NextTask()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	return nil
}
func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}
func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m *Model) initLists(width, height int) {
	// That /divisor is to establish division of widths for each defaultList
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height-divisor)

	m.lists = []list.Model{defaultList, defaultList, defaultList}
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "buy strawbery", description: "strawberies 1kg"},
		Task{status: todo, title: "get chocolate", description: "hersheies"},
	})

	m.lists[inProgress].Title = "In progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: todo, title: "get this done", description: "charm"},
	})

	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: todo, title: "Get started", description: "Have to try and get done with bars"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			// This reinitialises the whole thing, so its good for like updates
			m.initLists(msg.Width, msg.Height)
			m.loaded = true

		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.qutting = true
			return m, tea.Quit
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		case "enter":

			return m, m.MoveNext
		}

	}
	var cmd tea.Cmd
	// Kind of important because of deciding between lists
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	// The below code clear terminal when killed
	if m.qutting {
		return ""
	}
	// .View() gives the string representation of this
	if m.loaded {
		todoView := m.lists[todo].View()
		inProgView := m.lists[inProgress].View()
		doneView := m.lists[done].View()
		switch m.focused {
		case inProgress:
			return lipgloss.JoinHorizontal(lipgloss.Left,
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgView),
				columnStyle.Render(doneView),
			)

		case done:
			return lipgloss.JoinHorizontal(lipgloss.Left,
				columnStyle.Render(todoView),
				columnStyle.Render(inProgView),
				focusedStyle.Render(doneView),
			)
		default:
			return lipgloss.JoinHorizontal(lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgView),
				columnStyle.Render(doneView),
			)
		}

	} else {
		return "loading....."
	}

}

func main() {
	m := New()
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		os.Exit(1)
	}
}
