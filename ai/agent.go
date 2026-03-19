package ai

import (
	"context"
	"fmt"
	"strings"

	"help/config"

	"github.com/Ingenimax/agent-sdk-go/pkg/agent"
	"github.com/Ingenimax/agent-sdk-go/pkg/interfaces"
	"github.com/Ingenimax/agent-sdk-go/pkg/llm/openai"
	"github.com/Ingenimax/agent-sdk-go/pkg/logging"
)

type Agent struct {
	agent *agent.Agent
}

func NewAgent() (*Agent, error) {
	// Silence the noisy SDK debug logs by setting the logger explicitly to error level
	logger := logging.New()
	logging.WithLevel("error")(logger)

	llm := openai.NewClient(
		config.OPEN_ROUTER_KEY,
		openai.WithBaseURL("https://openrouter.ai/api/v1"),
		openai.WithModel(config.MODEL_ID),
		openai.WithLogger(logger),
	)

	a, err := agent.NewAgent(
		agent.WithLLM(llm),
		agent.WithSystemPrompt("You are a helpful expert developer terminal assistant."),
		agent.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	return &Agent{agent: a}, nil
}

func (a *Agent) RewriteCommand(ctx context.Context, failedCmd string) (string, error) {
	prompt := fmt.Sprintf("The user typed a terminal command that was not found: `%s`. Please deduce what they meant and provide the correct terminal command. ONLY output the corrected command, without any quotes, markdown formatting, or explanation.", failedCmd)
	resp, err := a.agent.Run(ctx, prompt)
	if err != nil {
		return "", err
	}
	// Clean up any stray markdown ticks just in case the AI ignored instructions
	resp = strings.TrimPrefix(strings.TrimSpace(resp), "`")
	resp = strings.TrimSuffix(resp, "`")
	return strings.TrimSpace(resp), nil
}

func (a *Agent) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	prompt := fmt.Sprintf("Generate a concise, professional Git commit message for the following diff. Only output the commit message, no markdown formatting, markdown quotes, or extra text.\n\nDiff:\n%s", diff)
	resp, err := a.agent.Run(ctx, prompt)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

func (a *Agent) StreamDebug(ctx context.Context, failedCmd, errorOutput string) (<-chan interfaces.AgentStreamEvent, error) {
	prompt := fmt.Sprintf("The following command failed:\n`%s`\n\nError output:\n```\n%s\n```\n\nPlease efficiently explain why it failed and suggest exactly how to fix it.", failedCmd, errorOutput)
	return a.agent.RunStream(ctx, prompt)
}
