package main
 
import (
    "context"
    "fmt"
    "log"
    "time"
 
    "github.com/aws/aws-sdk-go-v2/config"
    // "github.com/aws/aws-sdk-go-v2/service/ecr"
    "github.com/aws/aws-sdk-go-v2/service/eventbridge"
 
    helpers "setup/helpers/eventbridge"
)
 
func main() {
    fmt.Println("AWS EventBridge Setup")
 
    // Load the AWS SDK configuration
    cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }
 
    // Create an EventBridge client
    eventBridgeClient := eventbridge.NewFromConfig(cfg)
 
    // Setup EventBridge
    eventBusName, err := helpers.CreateOrUpdateEventBus(eventBridgeClient)
    if err != nil {
        log.Fatalf("Failed to create event bus: %v", err)
    }
 
    ruleName, err := helpers.CreateRule(eventBridgeClient, eventBusName)
    if err != nil {
        log.Fatalf("Failed to create rule: %v", err)
    }
 
    log.Println(ruleName)
    time.Sleep(5 * time.Second)
 
    err = helpers.AddTarget(eventBridgeClient, ruleName, eventBusName)
    if err != nil {
        log.Fatalf("Failed to add target: %v", err)
    }
 
    fmt.Println("EventBridge setup completed successfully.")
 
}