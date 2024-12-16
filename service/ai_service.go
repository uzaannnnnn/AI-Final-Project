package service

import (
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
) 

type HTTPClient interface { 
        Do(req *http.Request) (*http.Response, error) 
} 
type AIService struct { 
    Client HTTPClient 
} 
func (s *AIService) ChatWithAI(context, query, token string, translationService *TranslationService) (model.ChatResponse, error) {
    translated, err := translationService.Translate(query, "id", "en")
    if err != nil {
        return model.ChatResponse{}, err
    }

    input := map[string]interface{}{
        "model": "microsoft/Phi-3.5-mini-instruct",
        "messages": []map[string]string{
            {
                "role":    "user",
                "content": translated,
            },
        },
        "max_tokens": 600,
        "stream":     false,
    }
    body, err := json.Marshal(input)
    if err != nil {
        return model.ChatResponse{}, err
    }

    fmt.Println("ChatWithAI request body:", string(body))

    req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct/v1/chat/completions", bytes.NewBuffer(body))
    if err != nil {
        return model.ChatResponse{}, err
    }
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := s.Client.Do(req)
    if err != nil {
        return model.ChatResponse{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        respBody, _ := ioutil.ReadAll(resp.Body)
        fmt.Println("Error response body:", string(respBody))
        return model.ChatResponse{}, errors.New("failed to get chat response")
    }

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return model.ChatResponse{}, err
    }

    fmt.Println("ChatWithAI response body:", string(respBody))

    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        return model.ChatResponse{}, err
    }

    fmt.Printf("Parsed JSON result: %+v\n", result)

    // Coba ekstraksi data berdasarkan struktur respons JSON aktual
    var generatedText string

    // Cek apakah ada 'choices' dalam respons
    if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
        if choice, ok := choices[0].(map[string]interface{}); ok {
            if message, ok := choice["message"].(map[string]interface{}); ok {
                if text, ok := message["content"].(string); ok {
                    generatedText = text
                }
            }
        }
    } else if message, ok := result["message"].(map[string]interface{}); ok {
        if text, ok := message["content"].(string); ok {
            generatedText = text
        }
    }

    if generatedText == "" {
        return model.ChatResponse{}, errors.New("failed to extract generated text from response")
    }

    translatedAnswer, err := translationService.Translate(generatedText, "en", "id")
    if err != nil {
        return model.ChatResponse{}, err
    }

    var translatedResponse model.ChatResponse
    translatedResponse.GeneratedText = translatedAnswer

    return translatedResponse, nil
}

 
func (s *AIService) AnalyzeData(table map[string][]string, query, token string, translationService *TranslationService) (string, error) { 
    if len(table) == 0 { 
        return "", errors.New("table is empty") 
    } 
    translated, err := translationService.Translate(query, "id", "en") 
    if err != nil { 
        return "", err 
    } 
    input := model.AIRequest{ 
        Inputs: model.Inputs{ 
            Table: table, Query: translated, 
        }, 
    }
    body, err := json.Marshal(input) 
    if err != nil { 
        return "", err 
    } 
    fmt.Println("AnalyzeData request body:", string(body)) 
    req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq", bytes.NewBuffer(body)) 
    if err != nil { 
        return "", err 
    } 
    req.Header.Set("Authorization", "Bearer "+token) 
    req.Header.Set("Content-Type", "application/json") 
    resp, err := s.Client.Do(req) 
    if err != nil { 
        return "", err 
    } 
    defer resp.Body.Close() 
    if resp.StatusCode != http.StatusOK { 
        respBody, _ := ioutil.ReadAll(resp.Body) 
        fmt.Println("Error response body:", string(respBody)) 
        return "", errors.New("failed to analyze data") 
    } 
    var result model.TapasResponse 
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil { 
        return "", err 
    } 
    fmt.Println("AnalyzeData response:", result) // Extract the answer from the cells field 
    answer := "" 
    if len(result.Cells) > 0 { 
        answer = result.Cells[0] 
    }
    translatedAnswer, err := translationService.Translate(answer, "en", "id") 
    if err != nil { 
        return "", err 
    } 
    return translatedAnswer, nil 
}