package main_test

import (
    "a21hc3NpZ25tZW50/service"
    "bytes"
    "io/ioutil"
    "net/http"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("FileService", func() {
    var fileService *service.FileService

    BeforeEach(func() {
        fileService = &service.FileService{}
    })

    Describe("ProcessFile", func() {
        It("should return the correct result for valid CSV data", func() {
            fileContent := "header1,header2\nvalue1,value2\nvalue3,value4"
            expected := map[string][]string{
                "header1": {"value1", "value3"},
                "header2": {"value2", "value4"},
            }
            result, err := fileService.ProcessFile(fileContent)
            Expect(err).ToNot(HaveOccurred())
            Expect(result).To(Equal(expected))
        })

        It("should return an error for empty CSV data", func() {
            fileContent := ""
            result, err := fileService.ProcessFile(fileContent)
            Expect(err).To(HaveOccurred())
            Expect(result).To(BeNil())
        })

        It("should return an error for invalid CSV data", func() {
            fileContent := "header1,header2\nvalue1,value2\nvalue3"
            result, err := fileService.ProcessFile(fileContent)
            Expect(err).To(HaveOccurred())
            Expect(result).To(BeNil())
        })

        It("should handle CSV data with extra spaces", func() {
            fileContent := "header1 , header2 \n value1 , value2 \n value3 , value4 "
            expected := map[string][]string{
                "header1": {"value1", "value3"},
                "header2": {"value2", "value4"},
            }
            result, err := fileService.ProcessFile(fileContent)
            Expect(err).ToNot(HaveOccurred())
            Expect(result).To(Equal(expected))
        })
    })
})

type MockClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
    return m.DoFunc(req)
}

var _ = Describe("AIService", func() {
    var (
        mockClient         *MockClient
        aiService          *service.AIService
        translationService *service.TranslationService
    )

    BeforeEach(func() {
        mockClient = &MockClient{}
        aiService = &service.AIService{Client: mockClient}
        translationService = service.NewTranslationService("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")
    })

    Describe("AnalyzeData", func() {
        It("should return the correct result for a valid response", func() {
            mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
                return &http.Response{
                    StatusCode: http.StatusOK,
                    Body:       ioutil.NopCloser(bytes.NewBufferString(`{"cells":["result"]}`)),
                }, nil
            }
            table := map[string][]string{
                "header1": {"value1", "value2"},
            }
            query := "query"
            token := "token"
            result, err := aiService.AnalyzeData(table, query, token, translationService)
            Expect(err).ToNot(HaveOccurred())
            Expect(result).To(Equal("hasil"))
        })

        It("should return an error for an empty table", func() {
            table := map[string][]string{}
            query := "query"
            token := "token"
            result, err := aiService.AnalyzeData(table, query, token, translationService)
            Expect(err).To(HaveOccurred())
            Expect(result).To(BeEmpty())
        })

        It("should return an error for an error response", func() {
            mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
                return &http.Response{
                    StatusCode: http.StatusInternalServerError,
                    Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal error"}`)),
                }, nil
            }
            table := map[string][]string{
                "header1": {"value1", "value2"},
            }
            query := "query"
            token := "token"
            result, err := aiService.AnalyzeData(table, query, token, translationService)
            Expect(err).To(HaveOccurred())
            Expect(result).To(BeEmpty())
        })
    })

    Describe("ChatWithAI", func() { 
        It("should return the correct response for a valid request (array response)", func() { 
            mockClient.DoFunc = func(req *http.Request) (*http.Response, error) { 
                return &http.Response{ 
                    StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBufferString(`{"choices":[{"message":{"content":"response"}}]}`)), 
                    }, nil } 
            context := "context" 
            query := "query" 
            token := "token" 
            result, err := aiService.ChatWithAI(context, query, token, translationService) 
            Expect(err).ToNot(HaveOccurred()) 
            Expect(result.GeneratedText).To(Equal("respon")) 
        }) 
        It("should return the correct response for a valid request (object response)", func() { 
            mockClient.DoFunc = func(req *http.Request) (*http.Response, error) { 
                return &http.Response{ 
                    StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBufferString(`{"message":{"content":"response"}}`)), 
                    }, nil } 
            context := "context" 
            query := "query" 
            token := "token" 
            result, err := aiService.ChatWithAI(context, query, token, translationService) 
            Expect(err).ToNot(HaveOccurred()) 
            Expect(result.GeneratedText).To(Equal("respon")) 
        }) 
        It("should return an error for an error response", func() { 
            mockClient.DoFunc = func(req *http.Request) (*http.Response, error) { 
                return &http.Response{ 
                    StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal error"}`)), 
                    }, nil } 
            context := "context" 
            query := "query" 
            token := "token" 
            result, err := aiService.ChatWithAI(context, query, token, translationService) 
            Expect(err).To(HaveOccurred()) 
            Expect(result.GeneratedText).To(BeEmpty()) 
        })
    })
})
