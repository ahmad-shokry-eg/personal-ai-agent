package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"help/config"

	"github.com/Ingenimax/agent-sdk-go/pkg/agent"
	"github.com/Ingenimax/agent-sdk-go/pkg/interfaces"
	"github.com/Ingenimax/agent-sdk-go/pkg/llm/openai"
)

type Agent struct {
	agent *agent.Agent
}

func NewAgent() (*Agent, error) {
	// Silence the noisy SDK debug logs by defaulting local logger to error only
	os.Setenv("LOG_LEVEL", "error")
	
	llm := openai.NewClient(
		config.OPEN_ROUTER_KEY,
		openai.WithBaseURL("https://openrouter.ai/api/v1"),
		openai.WithModel(config.MODEL_ID),
	)

	a, err := agent.NewAgent(
		agent.WithLLM(llm),
		agent.WithSystemPrompt("You are a helpful expert developer terminal assistant."),
	)
	if err != nil {
		return nil, err
	}

	return &Agent{agent: a}, nil
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
