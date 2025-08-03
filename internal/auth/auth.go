package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ValidateToken отправляет токен в Auth Service для проверки
func ValidateToken(token string) (string, error) {
	// Проверка на пустой токен
	if token == "" {
		return "", fmt.Errorf("empty token provided")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", "http://localhost:7777/api/v1/validate", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа один раз
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		var errorResult struct {
			Error string `json:"error"`
		}
		// Пытаемся декодировать ошибку
		if err := json.Unmarshal(body, &errorResult); err == nil && errorResult.Error != "" {
			return "", fmt.Errorf("authentication failed: %s", errorResult.Error)
		}
		return "", fmt.Errorf("authentication failed with status %d", resp.StatusCode)
	}

	// Для успешного ответа декодируем данные пользователя
	var result struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w, body: %s", err, string(body))
	}

	if result.UserID == "" {
		return "", fmt.Errorf("invalid token response: user_id is empty")
	}

	return result.UserID, nil
}
