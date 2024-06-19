package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
)

type eventbridgeClient interface {
	CreateEventBus(ctx context.Context, params *eventbridge.CreateEventBusInput, optFns ...func(*eventbridge.Options)) (*eventbridge.CreateEventBusOutput, error)
	PutRule(ctx context.Context, params *eventbridge.PutRuleInput, optFns ...func(*eventbridge.Options)) (*eventbridge.PutRuleOutput, error)
	PutTargets(ctx context.Context, params *eventbridge.PutTargetsInput, optFns ...func(*eventbridge.Options)) (*eventbridge.PutTargetsOutput, error)
}

func CreateOrUpdateEventBus(client eventbridgeClient) (string, error) {
	eventBusName := "eventbus"

	// Try to create the EventBus
	createInput := &eventbridge.CreateEventBusInput{
		Name: aws.String(eventBusName),
	}
	result, err := client.CreateEventBus(context.Background(), createInput)
	if err != nil {
		// Check if the error is because the EventBus already exists
		if strings.Contains(err.Error(), "ResourceAlreadyExistsException") {
			fmt.Println("Event bus already exists. Event bus Name:", eventBusName)
			return eventBusName, nil
		} else {
			fmt.Println("Error creating event bus:", err)
			return "", err
		}
	}

	eventBusName = strings.Split(aws.ToString(result.EventBusArn), "/")[1]
	fmt.Println("Event bus created successfully. Event bus Name:", eventBusName)
	return eventBusName, nil
}
