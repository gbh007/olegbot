package tgusecases

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"slices"
	"text/template"
	"time"

	"app/internal/domain"

	"github.com/samber/lo"
	"github.com/samber/lo/mutable"
)

const defaultPromptTemplate = `Тебя зовут именами
{{- range $index, $element := $.names -}}
	{{- if ne $index 0 -}}{{ "," }}{{ end -}}
	{{- " " -}}
	{{- $element -}}
{{- end }}.
Ты пишешь смешные фразы, на основании следующих фраз придумай новый смешной ответ для пользователя из одной фразы:
{{ range choiceString $.quotes 50 -}}
- "{{ . }}"
{{ end }}`

func (u *UseCases) randomQuote(ctx context.Context, botInfo *domain.Bot, msgHistory []string, useLLM bool) (string, error) {
	quotes, err := u.repo.Quotes(ctx, u.botID)
	if err != nil {
		return "", fmt.Errorf("use case: random quote: %w", err)
	}

	if len(quotes) == 0 {
		return "", fmt.Errorf("use case: random quote: no quotes")
	}

	if u.llm != nil && botInfo != nil && useLLM {
		tStart := time.Now()
		q, err := u.llmQuote(ctx, botInfo, quotes, msgHistory)
		if err != nil {
			u.logger.ErrorContext(ctx, "get llm quote", "err", err.Error())
		} else {
			u.logger.DebugContext(ctx, "use llm quote", "quote", q, "duration", time.Since(tStart).String())
			return q, nil
		}
	}

	return quotes[rand.Intn(len(quotes))].Text, nil
}

func (u UseCases) llmQuote(ctx context.Context, botInfo *domain.Bot, quotes []domain.Quote, msgHistory []string) (string, error) {
	timeout := time.Minute

	if u.llmTimeout > 0 {
		timeout = u.llmTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	qts := lo.Map(quotes, func(q domain.Quote, _ int) string { return q.Text })

	prompt, err := makePrompt(botInfo.LLMPrompt, botInfo.Tags, qts)
	if err != nil {
		return "", fmt.Errorf("make prompt: %w", err)
	}

	q, err := u.llm.GetQuote(ctx, prompt, msgHistory)
	if err != nil {
		return "", fmt.Errorf("get quote: %w", err)
	}

	return q, nil
}

func makePrompt(promptTemplate string, names, quotes []string) (string, error) {
	t := template.New("")

	t.Funcs(template.FuncMap{
		"choiceString": func(values []string, limit int) []string {
			values = slices.Clone(values)
			mutable.Shuffle(values)
			if len(values) > limit {
				values = values[:limit]
			}

			return values
		},
	})

	if promptTemplate == "" {
		promptTemplate = defaultPromptTemplate
	}

	t, err := t.Parse(promptTemplate)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	buff := &bytes.Buffer{}

	err = t.Execute(buff, map[string]any{
		"names":  names,
		"quotes": quotes,
	})
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buff.String(), nil
}
