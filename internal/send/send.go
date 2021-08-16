package send

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"uacl/model"

	"github.com/TomBowyerResearchProject/common/logger"
)

func ChatterUser(user *model.User) error {
	bodyBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", os.Getenv("CHATTER_URL"), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("CHATTER_SECRET"))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	logger.Info("Created a user... for chatter")

	return resp.Body.Close()
}
