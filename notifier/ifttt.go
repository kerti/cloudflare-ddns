package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/spf13/viper"
)

// IFTTT is the notifier service for IfThisThenThat
type IFTTT struct {
	V1 string `json:"value1"`
	V2 string `json:"value2"`
	V3 string `json:"value3"`
}

// Notify invokes the notifier
func (i *IFTTT) Notify() error {
	isActive := viper.GetBool("notifier.ifttt.webhook.active")
	if !isActive {
		logger.Info("[IFTTT] Notification not activated, skipping...")
		return nil
	}

	makerKey := viper.GetString("notifier.ifttt.webhook.makerKey")
	eventName := viper.GetString("notifier.ifttt.webhook.eventName")

	url := fmt.Sprintf("https://maker.ifttt.com/trigger/%s/with/key/%s", eventName, makerKey)

	body, err := json.Marshal(i)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	buffer := bytes.NewBuffer([]byte(body))
	client := &http.Client{}
	resp, err := client.Post(url, "application/json", buffer)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		logger.Debug("[IFTTT] Notification sent...")
		return nil
	}

	return fmt.Errorf("ifttt notifier API responded with HTTP/%v", resp.StatusCode)
}
