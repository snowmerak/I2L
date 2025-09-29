package i2l

import (
	"context"
	"fmt"
	"i2l/models"
	"os"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/ollama"
)

type I2L struct {
	g                   *genkit.Genkit
	generativeModel     ai.Model
	classificationModel ai.Model
	embeddingModel      ai.Embedder
}

func DefaultOllamaRAG(ctx context.Context) (*I2L, error) {
	o := &ollama.Ollama{
		ServerAddress: "http://localhost:11434",
		Timeout:       300,
	}

	g := genkit.Init(ctx,
		genkit.WithPlugins(o, &googlegenai.GoogleAI{
			APIKey: os.Getenv("GEMINI_API_KEY"),
		}),
	)

	gos, err := models.OllamaGptOss20b(g)
	if err != nil {
		return nil, fmt.Errorf("failed to define ollama gpt-oss:20b model: %w", err)
	}

	eg, err := models.OllamaEmbeddingGemma(g, o.ServerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to define ollama gemma embedding model: %w", err)
	}

	return &I2L{
		g:                   g,
		generativeModel:     gos,
		classificationModel: gos,
		embeddingModel:      eg,
	}, nil
}

func DefaultGoogleAIRAG(ctx context.Context) (*I2L, error) {
	o := &ollama.Ollama{
		ServerAddress: "http://localhost:11434",
		Timeout:       300,
	}

	g := genkit.Init(ctx,
		genkit.WithPlugins(o, &googlegenai.GoogleAI{
			APIKey: os.Getenv("GEMINI_API_KEY"),
		}),
	)

	gm, err := models.GoogleAI(g, "gemini-flash-latest")
	if err != nil {
		return nil, fmt.Errorf("failed to define googleai gemma-3-12b-it model: %w", err)
	}

	eg, err := models.OllamaEmbeddingGemma(g, o.ServerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to define googleai gemma embedding model: %w", err)
	}

	return &I2L{
		g:                   g,
		generativeModel:     gm,
		classificationModel: gm,
		embeddingModel:      eg,
	}, nil
}

func GraphRAGWithGenkit(g *genkit.Genkit) *I2L {
	return &I2L{g: g}
}
