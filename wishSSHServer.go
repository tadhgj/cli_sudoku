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
	"charm.land/log/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"

	// keyboard challenge
	gossh "golang.org/x/crypto/ssh"
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
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true // we don't care what kind of key it is
		}),

		// todo: see if I want a password; eg. sudoku
		// wish.WithPasswordAuth(func(ctx ssh.Context, password string) bool {
		// 	return password == "asd123"
		// }),

		wish.WithKeyboardInteractiveAuth(func(_ ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool {
			return true // we just want to let the user connect; we are using 'auth' simply as an identifier
		}),

		// the provided public key will be stored in a sqlite table matched to a user uuid
		// there will be a users table with the uuids corresponding to
		// - current game state
		// keys and users will be removed after 90 days of inactivity
		// what will manage *that* process? cron job that checks all rows
		// and removes ones with timestamps over limit
		// table:
		// PUB KEY; TIMESTAMP; GAMESTATE (stuffed into a string? blob?)

		// could support -t commands so you could say ssh <url> -t hard for a new hard game

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

	// deal with the ssh key here
	// only try this if ssh key is valid!
	// yourKey := string(s.PublicKey().Marshal()[:])
	// fmt.Println(yourKey)

	// don't run until we know ssh key system works
	// iden := SSHIdentity{
	// 	sshKeyMarshalled: yourKey,
	// }

	iden := SSHIdentity{
		sshKeyMarshalled: "hello",
	}

	m := model{
		term:   pty.Term,
		width:  pty.Window.Width,
		height: pty.Window.Height,
		game:   InitializeGameBasedOnIdentity(iden),
	}

	ctx := s.Context()
	go func() {
		<-ctx.Done()
		// fires when a client disconnect
		// TODO: stop their game timer; save their state to db
	}()

	// return m, []tea.ProgramOption{tea.WithAltScreen()} // moved to 'view' cmd
	return m, []tea.ProgramOption{}
}

type model struct {
	term   string
	width  int
	height int
	game   SudokuGameWrapperState
	// quitStyle lipgloss.Style
}

func (m model) Init() tea.Cmd {
	// fires when new user connects
	fmt.Println("tea.Cmd calls init")

	// check what to do
	// new user: show game selection screen
	// existing user: enter game
	// 1 sec delay
	// time.Sleep(1 * time.Second)
	// m.game.SetToFreshWrapper() // won't do anything because method doesn't have pointer of model, and can't?

	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// fmt.Println("wish Update called")
	// fmt.Println(msg)
	var retur tea.Cmd = nil
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, tea.ClearScreen
	case tea.BackgroundColorMsg:
		m.game, retur = m.game.WrapperUpdate(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
		m.game, retur = m.game.WrapperUpdate(msg)
	}
	return m, retur
}

func (m model) View() tea.View {
	return m.game.WrapperView()
}
