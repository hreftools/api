package emails

import (
	"bytes"
	_ "embed"
	templateHtml "html/template"
	templateTxt "text/template"
)

//go:embed templates/auth-signup.html
var AuthSignupTemplateHtml string

//go:embed templates/auth-signup.txt
var AuthSignupTemplateTxt string

//go:embed templates/auth-resend-verification.html
var AuthResendVerificationTemplateHtml string

//go:embed templates/auth-resend-verification.txt
var AuthResendVerificationTemplateTxt string

type AuthSignupParams struct {
	Username string
	Email    string
	Token    string
}

type AuthResendVerificationParams struct {
	Token string
}

func RenderTemplateHtml(template string, data any) (string, error) {
	tmpl, err := templateHtml.New("email").Parse(template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func RenderTemplateTxt(template string, data any) (string, error) {
	tmpl, err := templateTxt.New("email").Parse(template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
