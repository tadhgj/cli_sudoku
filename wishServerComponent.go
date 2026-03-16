package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"

	// "strings"

	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

func main() {
	// Create a new SSH server.
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		// needs keys
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		// todo: leaves terminal in un-ideal state; fix
		// wish.WithIdleTimeout(1*time.Minute),
		// todo: probably don't want a banner if we enter a full screen thing
		// wish.WithBannerHandler(func(ctx ssh.Context) string {
		// 	return fmt.Sprintf(banner, ctx.User())
		// }),
		// todo: see if I want a password; eg. sudoku
		// wish.WithPasswordAuth(func(ctx ssh.Context, password string) bool {
		// 	return password == "asd123"
		// }),

		// actually, no, what I want is public key auth.
		// your public key will be stored in a sqlite table matched to a user uuid
		// there will be a users table with the uuids corresponding to
		// - current game state
		// keys and users will be removed after 90 days of inactivity
		// what will manage *that* process? cron job that checks all rows
		// and removes ones with timestamps over limit
		// table:
		// PUB KEY; TIMESTAMP; GAMESTATE (stuffed into a string? blob?)

		// could support -t commands so you could say ssh <url> -t hard for a new hard game

		// how to switch between screens?

		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
			// consider elapsed.Middleware()? why? it was included on the 'banner'
			// example but I'm not sure why I need it.
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {

	fmt.Println("A new client is connecting, initalizing instance...")

	pty, _, _ := s.Pty()

	m := model{
		term:       pty.Term,
		width:      pty.Window.Width,
		height:     pty.Window.Height,
		style:      lipgloss.NewStyle(),
		bg:         "light", // default
		boardState: GenerateSudokuBoardState(""),
		quitStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}

	// return m, []tea.ProgramOption{tea.WithAltScreen()} // moved to 'view' cmd
	return m, []tea.ProgramOption{}
}

type model struct {
	term       string
	width      int
	height     int
	style      lipgloss.Style
	bg         string
	boardState SudokuBoardInteractionState
	quitStyle  lipgloss.Style
}

func (m model) Init() tea.Cmd {
	// fires when new user connects
	fmt.Println("tea.Cmd calls init")

	// check what to do
	// new user: show game selection screen
	// existing user: enter game

	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		if msg.IsDark() {
			m.bg = "dark"
		} else {
			m.bg = "light"
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		// todo make nice keymap
		case "up":
			m.boardState.MoveCursorUp()
		case "down":
			m.boardState.MoveCursorDown()
		case "left":
			m.boardState.MoveCursorLeft()
		case "right":
			m.boardState.MoveCursorRight()
		// todo: make this neater
		case "1":
			m.boardState.SetNumberAtCursor(1)
		case "2":
			m.boardState.SetNumberAtCursor(2)
		case "3":
			m.boardState.SetNumberAtCursor(3)
		case "4":
			m.boardState.SetNumberAtCursor(4)
		case "5":
			m.boardState.SetNumberAtCursor(5)
		case "6":
			m.boardState.SetNumberAtCursor(6)
		case "7":
			m.boardState.SetNumberAtCursor(7)
		case "8":
			m.boardState.SetNumberAtCursor(8)
		case "9":
			m.boardState.SetNumberAtCursor(9)

		case "delete", "backspace":
			m.boardState.SetNumberAtCursor(0)
		}

	}
	return m, nil
}

func (m model) View() tea.View {
	// todo: create sudoku board game and render

	outStr := "Tadhg-doku\n\n"

	// print random number
	// randInt := rand.Int()
	// randIntStr := strconv.Itoa(randInt)

	// outStr += randIntStr + "\n"

	// print board
	outStr += RenderSudokuBoardState(m.boardState, m.style, m.bg)

	outStr += "\n"

	outStr += m.quitStyle.Render("Press q to quit")
	v := tea.NewView(outStr)
	v.AltScreen = true
	return v
}
