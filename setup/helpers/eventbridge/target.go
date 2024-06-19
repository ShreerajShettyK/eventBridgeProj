package helpers

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func AddTarget(client eventbridgeClient, ruleName string, eventBusName string) error {
	target := types.Target{
		Arn: aws.String("arn:aws:lambda:us-east-1:975050154225:function:MyGoLambdaFunction"),
		Id:  aws.String("Lambda"),
	}

	log.Println(ruleName)
	// Add the target to the existing rule
	_, err := client.PutTargets(context.Background(), &eventbridge.PutTargetsInput{
		Rule:         aws.String(ruleName),
		Targets:      []types.Target{target},
		EventBusName: aws.String(eventBusName),
	})

	if err != nil {
		fmt.Println("Error adding target to rule:", err.Error())
		return err
	}

	fmt.Println("Target added to rule successfully")
	return nil
}
