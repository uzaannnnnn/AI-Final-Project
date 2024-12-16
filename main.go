package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"a21hc3NpZ25tZW50/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Initialize the services
var fileService = &service.FileService{
    Repo: &repository.FileRepository{},
}
var aiService *service.AIService
var store = sessions.NewCookieStore([]byte("my-key"))

func getSession(r *http.Request) *sessions.Session {
    session, _ := store.Get(r, "chat-session")
    return session
}

func main() {
    // Load the .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Retrieve the Hugging Face token from the environment variables
    token := os.Getenv("HUGGINGFACE_TOKEN")
    if token == "" {
        log.Fatal("HUGGINGFACE_TOKEN is not set in the .env file")
    }

    // Initialize AIService
    aiService = &service.AIService{
        Client: &http.Client{},
    }

    // Set up the router
    router := mux.NewRouter()
    translationService := service.NewTranslationService(token)

    // File upload endpoint
    router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
        // Ensure the content type is multipart/form-data
        if r.Header.Get("Content-Type") == "" || !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
            http.Error(w, "Content-Type header is not multipart/form-data", http.StatusBadRequest)
            log.Println("Content-Type header is not multipart/form-data")
            return
        }

        file, header, err := r.FormFile("file")
        if err != nil {
            http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
            log.Println("Failed to read file:", err)
            return
        }
        defer file.Close()

        fmt.Println("File name:", header.Filename)

        content, err := ioutil.ReadAll(file)
        if err != nil {
            http.Error(w, "Failed to read file content: "+err.Error(), http.StatusInternalServerError)
            log.Println("Failed to read file content:", err)
            return
        }

        log.Println("File content:", string(content))

        table, err := fileService.ProcessFile(string(content))
        if err != nil {
            http.Error(w, "Failed to process file: "+err.Error(), http.StatusInternalServerError)
            log.Println("Failed to process file:", err)
            return
        }

        query := r.FormValue("query")
        if query == "" {
            http.Error(w, "Query is empty", http.StatusBadRequest)
            log.Println("Query is empty")
            return
        }

        session := getSession(r)
        session.Values["query"] = query
        session.Save(r, w)

        response, err := aiService.AnalyzeData(table, query, token, translationService)
        if err != nil {
            http.Error(w, "Failed to analyze data: "+err.Error(), http.StatusInternalServerError)
            log.Println("Failed to analyze data:", err)
            return
        }

        jsonResponse(w, map[string]string{"status": "success", "answer": response})
    }).Methods("POST")

    // Chat endpoint
    router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
        var input struct {
            Query string `json:"query"`
        }
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
            log.Println("Invalid request:", err)
            return
        }

        log.Println("Chat query:", input.Query)

        session := getSession(r)
        context := session.Values["context"]

        if context == nil {
            context = ""
        }

        response, err := aiService.ChatWithAI(context.(string), input.Query, token, translationService)
        if err != nil {
            http.Error(w, "Failed to get chat response: "+err.Error(), http.StatusInternalServerError)
            log.Println("Failed to get chat response:", err)
            return
        }

        log.Println("Chat response:", response.GeneratedText)

        session.Values["context"] = context.(string) + "\n" 
        if err := session.Save(r, w); err != nil {
            http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
            log.Println("Failed to save session:", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(map[string]string{"status": "success", "answer": response.GeneratedText}); err != nil {
            http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
            log.Println("Failed to encode response:", err)
            return
        }
    }).Methods("POST")

    // Enable CORS
    corsHandler := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
    }).Handler(router)

    // Start the server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server running on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
