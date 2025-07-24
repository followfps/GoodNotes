package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ValidateToken отправляет токен в Auth Service для проверки
func ValidateToken(token string) (string, error) {
	////reqBody, _ := json.Marshal(map[string]string{"token": token})
	//resp, err := http.Get(
	//	"http://localhost:7777/api/v1/validate",
	//)
	//if err != nil {
	//	return "", err
	//}
	//defer resp.Body.Close()

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:7777/api/v1/validate", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		UserID string `json:"user_id"`
		Error  string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(result.Error)
	}
	return result.UserID, nil
}
