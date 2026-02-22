package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"golang.org/x/term"
)

// Prompt обслуживает пользовательский ввод
type Prompt struct {
	in  io.Reader
	out io.Writer
}

func NewPrompt() *Prompt {
	return &Prompt{in: os.Stdin, out: os.Stdout}
}

// SecretName ввод названия секрета
func (p *Prompt) SecretName() (string, error) {
	label := fmt.Sprintf("Enter secret name (max %d chars): ", dto.SecretNameMaxLen)
	name := p.prompt(label)
	if name == "" {
		return "", fmt.Errorf("name can not be empty")
	}
	if len([]rune(name)) > dto.SecretNameMaxLen {
		return "", fmt.Errorf("name is too long, %d max", dto.SecretNameMaxLen)
	}
	return name, nil
}

// RegisterCredentials ввод учетных данных для регистрации
func (p *Prompt) RegisterCredentials() (dto.Credentials, error) {
	cr := dto.Credentials{}
	cr.Login = p.prompt("login: ")
	if err := cr.ValidateLogin(); err != nil {
		return cr, err
	}
	cr.Password = p.promptPassword("password: ")
	if err := cr.ValidatePassword(); err != nil {
		return cr, err
	}
	return cr, nil
}

// Credentials ввод учетных данных для входа или сохранения в системе
func (p *Prompt) Credentials() (dto.Credentials, error) {
	cr := dto.Credentials{
		Login:    p.prompt("login: "),
		Password: p.promptPassword("password: "),
	}
	if cr.Login == "" || cr.Password == "" {
		return cr, fmt.Errorf("login and password can not be empty")
	}
	return cr, nil
}

// Overwrite() запрашивает у пользователя нужно ли файл переписать
func (p *Prompt) Overwrite(fileName string) bool {
	label := fmt.Sprintf("File \"%s\" allready exist, overwrite it Y/n? [Y]: ", fileName)
	answer := p.prompt(label)
	if answer == "Y" || answer == "y" || answer == "" {
		return true
	}
	return false
}

func (p *Prompt) prompt(label string) string {
	fmt.Fprint(p.out, label)
	reader := bufio.NewReader(p.in)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func (p *Prompt) promptPassword(label string) string {
	fmt.Fprint(p.out, label)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return ""
	}
	fmt.Fprintln(p.out)
	return string(bytePassword)
}
