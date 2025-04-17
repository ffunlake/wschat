package services

import (
    "bytes"
    "encoding/json"
    "net/http"
    "time"
    "github.com/spf13/viper"
    "log"
    "fmt"
)

// API docs：https://platform.deepseek.com/api-docs
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type RequestBody struct {
    Model        string    `json:"model"`
    Messages     []Message `json:"messages"`
    Temperature  float32   `json:"temperature,omitempty"`
    MaxTokens    int       `json:"max_tokens,omitempty"`
}

type ResponseBody struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
    Error struct {
        Message string `json:"message"`
    } `json:"error"`
}
type AiService struct {
   apiKey string
   apiURL string
   model string
   temperature float32
   maxTokens int
   timeout int
   transport *http.Transport
}
var llmModel = "ai.deepseek"
func GetAiService() *AiService {
    // Get API configuration
    return &AiService{
        apiKey: viper.GetString(llmModel + ".api_key"),
        apiURL: viper.GetString(llmModel + ".api_url"),
        model: viper.GetString(llmModel + ".model"),
        temperature: float32(viper.GetFloat64(llmModel + ".settings.temperature")),
        maxTokens: viper.GetInt(llmModel + ".settings.max_tokens"),
        timeout: viper.GetInt(llmModel + ".timeout"),
        transport: &http.Transport{
            MaxIdleConns:        viper.GetInt(llmModel + ".http_client.max_idle_conns"),              // Maximum number of idle connections
            MaxIdleConnsPerHost: viper.GetInt(llmModel + ".http_client.max_idle_conns_per_host"),              // Maximum idle connections per host
            IdleConnTimeout:     time.Duration(viper.GetInt(llmModel + ".http_client.idle_conn_timeout")) * time.Second, // How long to keep idle connections alive
            TLSHandshakeTimeout: time.Duration(viper.GetInt(llmModel + ".http_client.tls_handshake_timeout")) * time.Second, // TLS handshake timeout
            DisableCompression:  viper.GetBool(llmModel + ".http_client.disable_compression"),             // Disable compression for better performance
        },
    }
}
func (ais *AiService) SentimentAnalysis(content string) string {
    // 确保有日志输出
    log.Println("Starting sentiment analysis...")
    if content == "" {
        log.Println("No content provided for sentiment analysis")
        return ""
    }
    
    // 记录非敏感配置信息
    log.Printf("Using API URL: %s", ais.apiURL)
    log.Printf("Using model: %s", ais.model)
    
    // set http client
    client := &http.Client{
        Transport: ais.transport,
        Timeout: time.Duration(ais.timeout) * time.Second,
    }

    // set request body
    requestBody := RequestBody{
        Model: ais.model,
        Messages: []Message{
            {
                Role:    "system",
                Content: "You are an AI that strictly classifies sentiment as Positive, Neutral, or Negative. Respond ONLY with the classification word.",
            },
            {
                Role:    "user",
                Content: content,
            },
        },
        Temperature: ais.temperature,
        MaxTokens:   ais.maxTokens,
    }

    // 序列化请求主体
    bodyBytes, err := json.Marshal(requestBody)
    if err != nil {
        log.Printf("Error marshaling request: %v", err)
        return ""
    }
    
    // 创建HTTP请求
    req, err := http.NewRequest("POST", ais.apiURL, bytes.NewBuffer(bodyBytes))
    if err != nil {
        log.Printf("Error creating request: %v", err)
        return ""
    }
    
    // 设置请求头
    contentType := viper.GetString(llmModel + ".headers.content_type")
    accept := viper.GetString(llmModel + ".headers.accept")
    req.Header.Set("Content-Type", contentType)
    req.Header.Set("Authorization", "Bearer "+ais.apiKey)
    req.Header.Set("Accept", accept)

    resp, err := client.Do(req)
    if err != nil {
        log.Printf("API request failed: %v", err)
        return ""
    }
    defer resp.Body.Close()


    if resp.StatusCode != http.StatusOK {
        var errorResp ResponseBody
        _ = json.NewDecoder(resp.Body).Decode(&errorResp)
        log.Printf("API Error %d: %s", resp.StatusCode, errorResp.Error.Message)
        return ""
    }


    var response ResponseBody
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        log.Printf("Failed to parse response: %v", err)
        return ""
    }

    if len(response.Choices) > 0 {
        result := response.Choices[0].Message.Content
        log.Printf("Sentiment analysis result: %s", result)
        return result
    }
    
    log.Println("No result returned from API")
    return ""
}
// GenerateAIResponse generates a response using DeepSeek AI
func (ai *AiService) GenerateAIResponse(input string) string {
    // Get API configuration


    // Create HTTP client with timeout
    client := &http.Client{
        Transport: ai.transport,
        Timeout: time.Duration(ai.timeout) * time.Second,
    }

    // Prepare request body
    requestBody := RequestBody{
        Model: ai.model,
        Messages: []Message{
            {
                Role:    "user",
                Content: input,
            },
        },
        Temperature: ai.temperature,
        MaxTokens:   ai.maxTokens,
    }

    // Convert request body to JSON
    bodyBytes, err := json.Marshal(requestBody)
    if err != nil {
        log.Printf("Error marshaling request body: %v", err)
        return ""
    }

    // Create HTTP request
    req, err := http.NewRequest("POST", ai.apiURL, bytes.NewBuffer(bodyBytes))
    if err != nil {
        log.Printf("Error creating request: %v", err)
        return ""
    }

    // Set headers
    contentType := viper.GetString(llmModel + ".headers.content_type")
    accept := viper.GetString(llmModel + ".headers.accept")
    req.Header.Set("Content-Type", contentType)
    req.Header.Set("Authorization", "Bearer "+ai.apiKey)
    req.Header.Set("Accept", accept)

    // Send request
    log.Println("Sending request to DeepSeek API for response generation...")
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("API request failed: %v", err)
        return ""
    }
    defer resp.Body.Close()

    // Handle response
    if resp.StatusCode != http.StatusOK {
        var errorResp ResponseBody
        _ = json.NewDecoder(resp.Body).Decode(&errorResp)
        log.Printf("API Error %d: %s", resp.StatusCode, errorResp.Error.Message)
        return ""
    }

    // Parse response
    var response ResponseBody
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        log.Printf("Failed to parse response: %v", err)
        return ""
    }
    fmt.Println(response.Choices)
    // Return the generated response                    
    return response.Choices[0].Message.Content      
}   
