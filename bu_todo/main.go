package main

import "C"

import (
    "encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"plugin"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const pluginPath = "./myplugin/myplugin.so"

// ─── スタイル定義 ──────────────────────────────────────────────
var (
	// Codex ライクな暗めのカラーパレット
	colorBg      = lipgloss.Color("#0d0d0d")
	colorSurface = lipgloss.Color("#1a1a1a")
	colorBorder  = lipgloss.Color("#2a2a2a")
	colorAccent  = lipgloss.Color("#7c6aff") // 紫アクセント
	colorGreen   = lipgloss.Color("#3ddc84")
	colorText    = lipgloss.Color("#e8e8e8")
	colorMuted   = lipgloss.Color("#666666")
	colorUser    = lipgloss.Color("#ffffff")
	colorSystem  = lipgloss.Color("#7c6aff")

	// ヘッダースタイル
	headerStyle = lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorAccent).
			Bold(true).
			Padding(0, 2)

	// モデル表示
	modelBadgeStyle = lipgloss.NewStyle().
			Foreground(colorBg).
			Background(colorAccent).
			Padding(0, 1).
			Bold(true)

	// ユーザーメッセージ
	userLabelStyle = lipgloss.NewStyle().
			Foreground(colorUser).
			Bold(true)

	userBubbleStyle = lipgloss.NewStyle().
			Foreground(colorUser).
			PaddingLeft(4)

	// アシスタントメッセージ
	assistantLabelStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true)

	assistantBubbleStyle = lipgloss.NewStyle().
				Foreground(colorText).
				PaddingLeft(4)

	// システムメッセージ
	systemStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true).
			PaddingLeft(2)

	// 入力エリア枠
	inputBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorBorder).
				Padding(0, 1)

	inputFocusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorAccent).
				Padding(0, 1)

	// ステータスバー
	statusStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			PaddingLeft(2)

	statusActiveStyle = lipgloss.NewStyle().
				Foreground(colorGreen).
				PaddingLeft(2)

	// コードブロック
	codeStyle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Background(lipgloss.Color("#111111")).
			Padding(0, 1)
)

// ─── メッセージ型 ────────────────────────────────────────────────
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

type Message struct {
	Role    Role
	Content string
	Time    time.Time
}

// ─── Bubble Tea メッセージ ────────────────────────────────────────
type responseMsg struct {
	content string
}

type errMsg struct{ err error }

// ─── モデル ──────────────────────────────────────────────────────
type model struct {
	messages  []Message
	viewport  viewport.Model
	textarea  textarea.Model
	spinner   spinner.Model
	width     int
	height    int
	loading   bool
	err       error
	modelName string
}

func initialModel() model {
	// テキストエリアの初期化
	ta := textarea.New()
	//ta.Placeholder = "メッセージを入力... (Enter で送信、Shift+Enter で改行)"
	ta.Placeholder = "Message Input... (Enter で送信、Shift+Enter で改行)"
	ta.Focus()
	ta.SetWidth(80)
	ta.SetHeight(3)
	ta.CharLimit = 4096
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetKeys("shift+enter")

	// スタイルカスタマイズ
	//ta.Styles.Base = lipgloss.NewStyle().
	//	Foreground(colorText)
	//ta.Styles.Placeholder = lipgloss.NewStyle().
	//	Foreground(colorMuted)
	//ta.Styles.CursorLine = lipgloss.NewStyle().
	//	Background(lipgloss.Color("#1f1f1f"))

	// ビューポートの初期化
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		Background(colorBg)

	// スピナーの初期化
	sp := spinner.New()
	sp.Spinner = spinner.Points
	sp.Style = lipgloss.NewStyle().Foreground(colorAccent)

	// 初期システムメッセージ
	messages := []Message{
		{
			Role:    RoleSystem,
			Content: "Chat へようこそ。何でもお気軽にどうぞ。",
			Time:    time.Now(),
		},
	}

	return model{
		messages:  messages,
		viewport:  vp,
		textarea:  ta,
		spinner:   sp,
		modelName: "model-123",
	}
}

// ─── Init ────────────────────────────────────────────────────────
func (m model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		m.spinner.Tick,
	)
}

// ─── Update ──────────────────────────────────────────────────────
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd  tea.Cmd
		vpCmd  tea.Cmd
		spCmd  tea.Cmd
		cmds   []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			if !m.loading {
				input := strings.TrimSpace(m.textarea.Value())
				if input == "" {
					break
				}
				// ユーザーメッセージ追加
				m.messages = append(m.messages, Message{
					Role:    RoleUser,
					Content: input,
					Time:    time.Now(),
				})
				m.textarea.Reset()
				m.loading = true
				m.viewport.SetContent(m.renderMessages())
				m.viewport.GotoBottom()

				// API 呼び出し（ここでは模擬）
				cmds = append(cmds, simulateAPICall(input))
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerHeight := 3
		statusHeight := 2
		inputHeight := 6
		vpHeight := m.height - headerHeight - statusHeight - inputHeight
		if vpHeight < 1 {
			vpHeight = 1
		}
		m.viewport.Width = m.width - 2
		m.viewport.Height = vpHeight
		m.textarea.SetWidth(m.width - 6)
		m.viewport.SetContent(m.renderMessages())

	case responseMsg:
		m.loading = false
		m.messages = append(m.messages, Message{
			Role:    RoleAssistant,
			Content: msg.content,
			Time:    time.Now(),
		})
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()

	case errMsg:
		m.loading = false
		m.err = msg.err

	case spinner.TickMsg:
		m.spinner, spCmd = m.spinner.Update(msg)
		cmds = append(cmds, spCmd)
	}

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, taCmd, vpCmd)

	return m, tea.Batch(cmds...)
}

// ─── View ────────────────────────────────────────────────────────
func (m model) View() string {
	if m.width == 0 {
		return "読み込み中..."
	}

	var b strings.Builder

	// ── ヘッダー ──
	logo := headerStyle.Render("◆ App")
	badge := modelBadgeStyle.Render(m.modelName)
	spacer := strings.Repeat(" ", max(0, m.width-lipgloss.Width(logo)-lipgloss.Width(badge)-4))
	header := lipgloss.NewStyle().
		Background(colorBg).
		Width(m.width).
		Render(logo + spacer + badge)

	separator := lipgloss.NewStyle().
		Foreground(colorBorder).
		Render(strings.Repeat("─", m.width))

	b.WriteString(header + "\n")
	b.WriteString(separator + "\n")

	// ── チャットビューポート ──
	b.WriteString(m.viewport.View() + "\n")
	b.WriteString(separator + "\n")

	// ── 入力エリア ──
	var inputBox string
	if m.textarea.Focused() {
		inputBox = inputFocusedBorderStyle.
			Width(m.width - 4).
			Render(m.textarea.View())
	} else {
		inputBox = inputBorderStyle.
			Width(m.width - 4).
			Render(m.textarea.View())
	}
	b.WriteString(inputBox + "\n")

	// ── ステータスバー ──
	var status string
	if m.loading {
		status = statusActiveStyle.Render(m.spinner.View() + " 応答を生成中...")
	} else if m.err != nil {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).
			Render("✗ エラー: " + m.err.Error())
	} else {
		msgCount := fmt.Sprintf("%d メッセージ", len(m.messages))
		status = statusStyle.Render("↑↓ スクロール  Enter 送信  Esc 終了  " + msgCount)
	}
	b.WriteString(status)

	return b.String()
}

// ─── メッセージ描画 ──────────────────────────────────────────────
func (m model) renderMessages() string {
	var b strings.Builder
	for i, msg := range m.messages {
		if i > 0 {
			b.WriteString("\n")
		}
		switch msg.Role {
		case RoleSystem:
			b.WriteString(systemStyle.Render("◆ " + msg.Content))

		case RoleUser:
			ts := lipgloss.NewStyle().Foreground(colorMuted).Render(msg.Time.Format("15:04"))
			label := userLabelStyle.Render("▸ You") + "  " + ts
			b.WriteString(label + "\n")
			b.WriteString(userBubbleStyle.Render(msg.Content))

		case RoleAssistant:
			ts := lipgloss.NewStyle().Foreground(colorMuted).Render(msg.Time.Format("15:04"))
			label := assistantLabelStyle.Render("◆ App") + "  " + ts
			b.WriteString(label + "\n")
			// コードブロックの簡易レンダリング
			rendered := renderCodeBlocks(msg.Content, m.width)
			b.WriteString(assistantBubbleStyle.Render(rendered))
		}
		b.WriteString("\n")
	}
	return b.String()
}

// ─── コードブロック簡易レンダリング ─────────────────────────────
func renderCodeBlocks(content string, width int) string {
	lines := strings.Split(content, "\n")
	var result []string
	inCode := false
	var codeLines []string
	lang := ""

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if !inCode {
				inCode = true
				lang = strings.TrimPrefix(line, "```")
				codeLines = nil
			} else {
				// コードブロック終了
				inCode = false
				header := ""
				if lang != "" {
					header = lipgloss.NewStyle().
						Foreground(colorMuted).
						Render("  " + lang + "\n")
				}
				codeContent := codeStyle.
					Width(width-8).
					Render(strings.Join(codeLines, "\n"))
				result = append(result, header+codeContent)
			}
		} else if inCode {
			codeLines = append(codeLines, "  "+line)
		} else {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

type TodoItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
// ─── API モック（実際には Anthropic API を呼ぶ） ───────────────────
func simulateAPICall(input string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(800 * time.Millisecond) // API レイテンシを模擬

		// 入力に応じたデモレスポンス
		var response string
		inputLower := strings.ToLower(input)
        p, err := plugin.Open(pluginPath)
        if err != nil {
                fmt.Fprintf(os.Stderr, "❌ プラグインのロードに失敗: %v\n", err)
                os.Exit(1)
        }						
		switch {
        case strings.HasPrefix(inputLower, "add"):
            var target_buff = strings.ReplaceAll(inputLower, "add ", "")
            todoAdd, err := p.Lookup("TodoAdd")
            if err != nil {
                fmt.Fprintf(os.Stderr, "❌ TodoAdd のルックアップに失敗: %v\n", err)
                os.Exit(1)
            }
            todoAddString, ok := todoAdd.(func(string) string)
            if !ok {
                fmt.Fprintf(os.Stderr, "❌ TodoAdd の型が不正\n")
                os.Exit(1)
            }
            result := todoAddString(target_buff)                      
            response = result

        case strings.Contains(inputLower, "list"):
            todoList, err := p.Lookup("TodoList")
            if err != nil {
                fmt.Fprintf(os.Stderr, "❌ TodoList のルックアップに失敗: %v\n", err)
                os.Exit(1)
            }
            todoListString, ok := todoList.(func() string)
            if !ok {
                fmt.Fprintf(os.Stderr, "❌ TodoList の型が不正\n")
                os.Exit(1)
            }
            result := todoListString()
            var items []TodoItem
            err = json.Unmarshal([]byte(result), &items)
            if err != nil {
                fmt.Println("JSONパースエラー:", err)
                os.Exit(1)
            }
            fmt.Println("パース後のデータ:")
            var out_str = ""
            for _, item := range items {
                tmp := fmt.Sprintf("id= %d , %s \n", item.ID, item.Title)
                out_str += tmp
            }
			response = out_str

		case strings.Contains(inputLower, "hello") || strings.Contains(inputLower, "こんにちは"):
			response = "こんにちは！Chat App です。コードの作成、デバッグ、質問など何でもどうぞ。"

		default:
			response = fmt.Sprintf(`「%s」についてですね。

もう少し具体的に教えていただけると、より的確なサポートができます。例えば：

- 実装したい機能の詳細
- 使用している言語やフレームワーク
- エラーメッセージがあれば内容

何でもお聞きください！`, input)
		}

		return responseMsg{content: response}
	}
}

// ─── ユーティリティ ──────────────────────────────────────────────
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ─── エントリポイント ────────────────────────────────────────────
func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),       // フルスクリーンモード
		tea.WithMouseCellMotion(), // マウスサポート
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}
