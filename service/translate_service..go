package service

import (
    "context"
    "errors"
    "log"
    "strings"

    hf "github.com/hupe1980/go-huggingface"
)

type TranslationService struct {
    Client *hf.InferenceClient
}

func NewTranslationService(apiKey string) *TranslationService {
    return &TranslationService{
        Client: hf.NewInferenceClient(apiKey),
    }
}

func (s *TranslationService) Translate(text, sourceLang, targetLang string) (string, error) {
    // Handling empty input
    if text == "" {
        return "", errors.New("input text is empty")
    }

    // Split text into smaller chunks if it's too long
    chunks := splitTextIntoChunks(text, 500) // Adjust chunk size as needed
    var translatedChunks []string

    for _, chunk := range chunks {
        res, err := s.Client.Translation(context.Background(), &hf.TranslationRequest{
            Inputs: []string{chunk},
            Model:  "Helsinki-NLP/opus-mt-" + sourceLang + "-" + targetLang,
        })
        if err != nil {
            log.Printf("Translation error for chunk: %v\n", err)
            return "", err
        }

        if len(res) == 0 {
            log.Println("No translation result for chunk")
            return "", errors.New("no translation result for chunk")
        }

        translatedChunks = append(translatedChunks, res[0].TranslationText)
    }

    return strings.Join(translatedChunks, " "), nil
}

func splitTextIntoChunks(text string, chunkSize int) []string {
    var chunks []string
    runes := []rune(text)
    for i := 0; i < len(runes); i += chunkSize {
        end := i + chunkSize
        if end > len(runes) {
            end = len(runes)
        }
        chunks = append(chunks, string(runes[i:end]))
    }
    return chunks
}
