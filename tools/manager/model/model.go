package model

import (
	"DarkFlameMaster/seat"
	"DarkFlameMaster/serverinfo"
	"DarkFlameMaster/tools/manager/config"
	"DarkFlameMaster/web"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

type errMsg error

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}

type Model struct {
	// for list
	list   list.Model
	choice string
	// for text input
	textInput textinput.Model
	title     string
	// normal
	step     Step
	success  bool
	err      error
	quitting bool
	about    bool
}

type Step string

var (
	StepChooseProgram      Step = "choose_program"
	StepDumpData           Step = "dump_data"
	StepDeleteSeat         Step = "delete_seat"
	StepDeleteSeatBySeat   Step = "delete_seat_by_seat"
	StepDeleteSeatByUserID Step = "delete_seat_by_user_id"
)

func NewModel() *Model {
	items := []list.Item{
		item("删除座位"),
		item("导出数据"),
		item("关于"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "爆裂吧，现实！粉碎吧，精神！消失吧！这个世界！"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &Model{
		list:     l,
		choice:   "",
		title:    "",
		step:     StepChooseProgram,
		success:  false,
		err:      nil,
		quitting: false,
		about:    false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) initTextInput(placeholder string) (textinput.Model, tea.Msg) {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.CursorStart()
	ti.Cursor.Blink = true

	return ti, cursor.Blink()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case StepChooseProgram:
		return m.updateChooseProgram(msg)
	case StepDeleteSeat:
		return m.updateDeleteSeat(msg)
	case StepDeleteSeatBySeat:
		return m.updateDeleteSeatBySeat(msg)
	case StepDeleteSeatByUserID:
		return m.updateDeleteSeatByUserID(msg)
	}
	return m, tea.Quit
}

func (m Model) View() string {
	switch m.step {
	case StepChooseProgram:
		return m.viewChooseProgram()
	case StepDumpData:
		return m.viewDumpData()
	case StepDeleteSeat:
		return m.viewDeleteSeat()
	case StepDeleteSeatBySeat:
		return m.viewDeleteSeatBySeat()
	case StepDeleteSeatByUserID:
		return m.viewDeleteSeatByUserID()
	}
	if m.quitting {
		return quitTextStyle.Render("See you later!\n\n")
	}
	return "\n" + m.list.View()
}

func (m Model) updateChooseProgram(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			switch m.choice {
			case "删除座位":
				itemSize := len(m.list.Items())
				// Remove all item
				for i := 0; i < itemSize; i++ {
					m.list.RemoveItem(0)
				}
				// Add item
				m.list.SetItems([]list.Item{
					item("根据座位号删除"),
					item("根据用户ID删除"),
				})
				m.step = StepDeleteSeat
			case "导出数据":
				m.step = StepDumpData
				return m.updateDumpData(msg)
			case "关于":
				m.quitting = true
				m.about = true
				return m, tea.Quit
			default:
				panic("invalid choice")
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) updateDumpData(_ tea.Msg) (tea.Model, tea.Cmd) {
	cfg := config.GetGlobalConfig()
	resp, err := http.Get(cfg.Server + "/dump_tickets")
	if err != nil {
		m.err = errMsg(fmt.Errorf("导出数据失败: %s", err.Error()))
		return m, tea.Quit
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		m.err = errMsg(fmt.Errorf("导出数据失败: %s", err.Error()))
		return m, tea.Quit
	}
	// write ticket information (csv data) to file.
	fileName := fmt.Sprintf("dump_tickets_%s.csv",
		time.Now().Format("20060102150405"))
	// create dump file and write data
	err = os.WriteFile(fileName, body, 0644)
	if err != nil {
		m.err = errMsg(fmt.Errorf("导出数据失败: %s", err.Error()))
		return m, tea.Quit
	}
	return m, tea.Quit
}

func (m Model) updateDeleteSeat(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			switch m.choice {
			case "根据座位号删除":
				m.step = StepDeleteSeatBySeat
				model, msg := m.initTextInput("input: 行号:列号（英文冒号）")
				m.title = "根据座位号删除座位"
				m.textInput = model
				m.textInput.Reset()
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			case "根据用户ID删除":
				m.step = StepDeleteSeatByUserID
				model, msg := m.initTextInput("input: 用户ID")
				m.title = "根据用户ID删除座位"
				m.textInput = model
				m.textInput.Reset()
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			default:
				panic("invalid choice")
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) updateDeleteSeatBySeat(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cfg := config.GetGlobalConfig()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			val := m.textInput.Value()
			seatStr := strings.Split(val, ":")

			handleErr := func(err error) (tea.Model, tea.Cmd) {
				m.textInput.Reset()
				m.err = errMsg(err)
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
			if len(seatStr) != 2 {
				return handleErr(fmt.Errorf("输入格式错误"))
			}
			row, err := strconv.Atoi(seatStr[0])
			if err != nil {
				return handleErr(fmt.Errorf("输入内容有误：%s", err))
			}
			col, err := strconv.Atoi(seatStr[1])
			if err != nil {
				return handleErr(fmt.Errorf("输入内容有误：%s", err))
			}
			seats := make([]*seat.Seat, 0, 1)
			seats = append(seats, &seat.Seat{
				Row:    uint(row),
				Column: uint(col),
			})
			req := web.DeleteTicketsReq{
				Mode:  "seat",
				Seats: seats,
			}
			data, err := json.Marshal(req)
			if err != nil {
				return handleErr(err)
			}
			resp, err := http.Post(cfg.Server+"/delete_tickets",
				"application/json", strings.NewReader(string(data)))
			if err != nil {
				m.err = errMsg(err)
				return m, tea.Quit
			}
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				m.err = errMsg(err)
				return m, tea.Quit
			}
			if string(body) == "ok" {
				m.success = true
				return m, tea.Quit
			} else {
				m.err = errMsg(fmt.Errorf("删除失败，%s", string(body)))
				return m, tea.Quit
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) updateDeleteSeatByUserID(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cfg := config.GetGlobalConfig()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			userId := m.textInput.Value()

			handleErr := func(err error) (tea.Model, tea.Cmd) {
				m.textInput.Reset()
				m.err = errMsg(err)
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}

			req := web.DeleteTicketsReq{
				Mode:  "user",
				Users: []string{userId},
			}
			data, err := json.Marshal(req)
			if err != nil {
				return handleErr(err)
			}
			resp, err := http.Post(cfg.Server+"/delete_tickets",
				"application/json", strings.NewReader(string(data)))
			if err != nil {
				m.err = errMsg(err)
				return m, tea.Quit
			}
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				m.err = errMsg(err)
				return m, tea.Quit
			}
			if string(body) == "ok" {
				m.success = true
				return m, tea.Quit
			} else {
				m.err = errMsg(fmt.Errorf("删除失败，%s", string(body)))
				return m, tea.Quit
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) viewChooseProgram() string {
	if m.quitting {
		if m.about {
			aboutStr := "MIT License\n\nCopyright (c) 2023 poppinzhang\n" +
				serverinfo.Get().String() + "\n\n"
			return quitTextStyle.Render(aboutStr)
		}
		return quitTextStyle.Render("See you later!\n\n")
	}
	return "\n" + m.list.View()
}

func (m Model) viewDumpData() string {
	if m.err != nil {
		return quitTextStyle.Render(fmt.Sprintf("%s\n\n", m.err))
	} else {
		return quitTextStyle.Render(fmt.Sprintf("导出成功！\n\n"))
	}
}

func (m Model) viewDeleteSeat() string {
	if m.quitting {
		return quitTextStyle.Render("See you later!\n\n")
	}
	return "\n" + m.list.View()
}

func (m Model) viewDeleteSeatBySeat() string {
	if m.quitting {
		return quitTextStyle.Render("See you later!\n\n")
	}
	if m.err != nil {
		str := fmt.Sprintf("Error: %v\n\n%s\n\n%s",
			m.err,
			m.textInput.View(),
			"(esc to quit)",
		) + "\n"
		m.err = nil
		return str
	}
	if m.success {
		return quitTextStyle.Render("删除成功！", "\n\n")
	}
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.title,
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func (m Model) viewDeleteSeatByUserID() string {
	return m.viewDeleteSeatBySeat()
}
