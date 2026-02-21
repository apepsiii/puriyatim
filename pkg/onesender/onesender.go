package onesender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	APIKey  string
	GroupID string
	Client  *http.Client
}

type MessageRequest struct {
	RecipientType string      `json:"recipient_type"`
	RecipientID   string      `json:"recipient_id"`
	Message       Message     `json:"message"`
	ScheduleTime  *time.Time  `json:"schedule_time,omitempty"`
}

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		MessageID string `json:"message_id"`
		Status    string `json:"status"`
	} `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
}

func NewClient(baseURL, apiKey, groupID string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		GroupID: groupID,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) SendTextMessage(recipientID, text string) (*MessageResponse, error) {
	req := &MessageRequest{
		RecipientType: "group",
		RecipientID:   recipientID,
		Message: Message{
			Type: "text",
			Text: text,
		},
	}

	return c.sendMessage(req)
}

func (c *Client) SendGroupMessage(text string) (*MessageResponse, error) {
	return c.SendTextMessage(c.GroupID, text)
}

func (c *Client) sendMessage(req *MessageRequest) (*MessageResponse, error) {
	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Send request
	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("request failed: %s", errorResp.Message)
	}

	// Parse success response
	var successResp MessageResponse
	if err := json.Unmarshal(body, &successResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &successResp, nil
}

// Helper function to format Jumat Berkah announcement message
func FormatJumatBerkahMessage(penerima []string, tanggal string) string {
	message := fmt.Sprintf("📢 *PENGUMUMAN PROGRAM JUMAT BERKAH*\n\n")
	message += fmt.Sprintf("Assalamu'alaikum Warahmatullahi Wabarakatuh\n\n")
	message += fmt.Sprintf("Bismillahirrahmannirrahim\n\n")
	message += fmt.Sprintf("Alhamdulillah, berikut adalah daftar penerima bantuan Program Jumat Berkah untuk hari %s:\n\n", tanggal)
	
	for i, nama := range penerima {
		message += fmt.Sprintf("%d. %s\n", i+1, nama)
	}
	
	message += "\nMohon kepada para wali/warga yang namanya tercantum untuk hadir tepat waktu dalam pengambilan bantuan.\n\n"
	message += "Terima kasih atas perhatian dan dukungannya.\n\n"
	message += "Wassalamu'alaikum Warahmatullahi Wabarakatuh\n\n"
	message += "🏠 *Panti Asuhan*"
	
	return message
}