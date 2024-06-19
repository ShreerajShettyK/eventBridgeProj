package helpers
 
import (
    "context"
    "fmt"
    "strings"
 
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/eventbridge"
)
 
func CreateRule(client eventbridgeClient, eventBusName string) (string, error) {
    eventPattern := `{
        "source": ["aws.ecr"],
        "detail-type": ["ECR Image Scan"]
    }`
 
    input := &eventbridge.PutRuleInput{
        Name:         aws.String("Rule-ECRPushEvent"),
        EventPattern: aws.String(eventPattern),
        EventBusName: aws.String(eventBusName),
    }
 
    result, err := client.PutRule(context.Background(), input)
    if err != nil {
        fmt.Println("Error creating rule:", err)
        return "", err
    }
    ruleName := strings.Split(aws.ToString(result.RuleArn), "/")[2]
 
    fmt.Println("Rule created successfully. Rule ARN:", *result.RuleArn)
    return ruleName, nil
}


