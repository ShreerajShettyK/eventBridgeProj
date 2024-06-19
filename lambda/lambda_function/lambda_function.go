package main
 
import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "time"
 
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
    "github.com/aws/aws-sdk-go-v2/service/ssm"
    "github.com/aws/aws-sdk-go-v2/service/ssm/types"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "errors"
)
 
// LambdaHandler defines the interface for the Lambda function handler.
type LambdaHandler interface {
    HandleRequest(ctx context.Context, event events.CloudWatchEvent) error
}
 
// SSMClient defines the interface for the AWS Systems Manager client.
type SSMClient interface {
    PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
}
 
// SecretsManagerClient defines the interface for the AWS Secrets Manager client.
type SecretsManagerClient interface {
    GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}
 
// MongoCollection defines the interface for the MongoDB collection.
type MongoCollection interface {
    InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}
 
// LambdaFunction implements the LambdaHandler interface.
type LambdaFunction struct {
    SSMClient            SSMClient
    SecretsManagerClient SecretsManagerClient
    MongoCollection      MongoCollection
}
 
// ParameterStoreData defines the structure for the data stored in Parameter Store.
type ParameterStoreData struct {
    ImageTag       string `json:"imageTag"`
    RepositoryName string `json:"repositoryName"`
    Date           string `json:"date"`
}
 
// NewLambdaFunction creates a new instance of LambdaFunction with AWS SDK clients initialized.
func NewLambdaFunction() *LambdaFunction {
    cfg, _ := config.LoadDefaultConfig(context.Background())
    return &LambdaFunction{
        SSMClient:            ssm.NewFromConfig(cfg),
        SecretsManagerClient: secretsmanager.NewFromConfig(cfg),
    }
}
 
// HandleRequest function to process the CloudWatch Event.
func (lf *LambdaFunction) HandleRequest(ctx context.Context, event events.CloudWatchEvent) error {
    // Log the incoming event for debugging purposes.
    log.Println("Received event:", event)
 
    // Unmarshal the event detail to extract the ECR event information.
    var ecrEvent events.ECRScanEventDetailType
    if err := json.Unmarshal(event.Detail, &ecrEvent); err != nil {
        return fmt.Errorf("error unmarshalling event detail: %v", err)
    }
 
    // Check if image tags are present in the event detail.
    if len(ecrEvent.ImageTags) == 0 {
        return fmt.Errorf("no image tags found in event detail")
    }
 
    // Generate a random 6-character string.
    randomString := generateRandomString(6)
 
    // Extract the latest image tag from the list.
    imageTag := fmt.Sprintf("latest,%s", randomString)
    repositoryName := ecrEvent.RepositoryName
    date := time.Now().Format("2006-01-02")
    log.Printf("Received image tag: %s from repository: %s", imageTag, repositoryName)
 
    // Prepare data for Parameter Store.
    psData := ParameterStoreData{
        ImageTag:       imageTag,
        RepositoryName: repositoryName,
        Date:           date,
    }
 
    psDataJSON, err := json.Marshal(psData)
    if err != nil {
        log.Printf("Error marshalling parameter store data: %v", err)
        return err
    }
 
    // Store data in Parameter Store.
    err = lf.storeInParameterStore(ctx, repositoryName, string(psDataJSON))
    if err != nil {
        return err
    }
 
    log.Println("Successfully updated the parameter store.")
 
    // Retrieve MongoDB connection string from Secrets Manager.
    secretValue, err := lf.SecretsManagerClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
        SecretId: aws.String("myApp/mongo-db-credentials"),
    })
    if err != nil {
        log.Printf("Error retrieving secret: %v", err)
        return err
    }
 
    var secretsMap map[string]string
    if err := json.Unmarshal([]byte(*secretValue.SecretString), &secretsMap); err != nil {
        log.Printf("Error unmarshalling secret: %v", err)
        return err
    }
 
    // Store image tag in MongoDB.
    err = lf.storeInMongoDB(ctx, imageTag, repositoryName, date)
    if err != nil {
        return err
    }
 
    fmt.Println("Image tag stored in MongoDB successfully!")
    return nil
}
 
// storeInParameterStore stores the given data in AWS Parameter Store.
func (lf *LambdaFunction) storeInParameterStore(ctx context.Context, name, value string) error {
    _, err := lf.SSMClient.PutParameter(ctx, &ssm.PutParameterInput{
        Name:      aws.String(name),
        Value:     aws.String(value),
        Type:      types.ParameterTypeString,
        Overwrite: aws.Bool(true),
    })
 
    if err != nil {
        log.Printf("Error storing parameter in SSM: %v", err)
        return err
    }
 
    return nil
}
 
// storeInMongoDB stores the image tag information in MongoDB.
func (lf *LambdaFunction) storeInMongoDB(ctx context.Context, imageTag, repositoryName, date string) error {
    if lf.MongoCollection == nil {
        return errors.New("MongoCollection not set")
    }
 
    _, err := lf.MongoCollection.InsertOne(ctx, map[string]interface{}{
        "imageTag":       imageTag,
        "repositoryName": repositoryName,
        "date":           date,
    })
 
    if err != nil {
        log.Printf("Error inserting image tag in MongoDB: %v", err)
        return err
    }
 
    return nil
}
 
// generateRandomString generates a random string of specified length.
func generateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}
 
func main() {
    // Start the Lambda function.
    lambdaFunction := NewLambdaFunction()
    lambda.Start(lambdaFunction.HandleRequest)
}