package eventpkg

import "time"

// CreateEvent creates an event and returns it
func CreateEvent() map[string]interface{} {
	currentTime := time.Now()

	event := map[string]interface{}{
		"version":     "0",
		"detail-type": "ECR Image Scan",
		"source":      "aws.ecr",
		"account":     "851725240994",
		"time":        currentTime.Format(time.RFC3339),
		"region":      "us-east-1",
		"resources":   []interface{}{},
		"detail": map[string]interface{}{
			"scan-status":     "COMPLETE",
			"repository-name": "my-app-repo",
			"image-tags":      []string{"latest"},
			"updated-date":    currentTime.Format(time.RFC3339),
		},
	}
	return event
}

