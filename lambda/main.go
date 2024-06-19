package main
 
import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/exec"
    "runtime"
    "time"
 
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/lambda"
    "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)
 
func main() {
    // Run the build script to prepare the deployment package
    if err := runBuildScript(); err != nil {
        log.Fatalf("failed to execute build script: %v", err)
    }
 
    // Load the AWS configuration
    cfg, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
        log.Fatalf("unable to load SDK config: %v", err)
    }
 
    // Create a Lambda client
    svc := lambda.NewFromConfig(cfg)
 
    // Define the Lambda function name
    functionName := "MyGoLambdaFunction"
 
    // Check if the Lambda function already exists
    if _, err := getFunction(svc, functionName); err != nil {
        // If the function does not exist, create it
        fmt.Println("Creating new Lambda function")
        if createErr := createFunction(svc, functionName); createErr != nil {
            log.Fatalf("failed to create function: %v", createErr)
        }
    } else {
        // If the function exists, update it
        fmt.Println("Updating existing Lambda function")
        if updateErr := updateFunction(svc, functionName); updateErr != nil {
            log.Fatalf("failed to update function code: %v", updateErr)
        }
    }
 
    // Invoke the Lambda function
    if err := invokeFunction(svc, functionName); err != nil {
        log.Fatalf("failed to invoke function: %v", err)
    }
}
 
func runBuildScript() error {
    // Determine the appropriate command to run the build script based on the OS
    var cmd *exec.Cmd
    if runtime.GOOS == "windows" {
        // Use Git Bash for running the build script on Windows
        gitBashPath := `C:\Program Files\Git\bin\bash.exe`
        cmd = exec.Command(gitBashPath, "-c", "bash build.sh")
    } else {
        cmd = exec.Command("/bin/bash", "build.sh")
    }
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
 
func getFunction(svc *lambda.Client, functionName string) (*lambda.GetFunctionOutput, error) {
    // Get the Lambda function details
    return svc.GetFunction(context.Background(), &lambda.GetFunctionInput{
        FunctionName: aws.String(functionName),
    })
}
 
func createFunction(svc *lambda.Client, functionName string) error {
    // Read the function code from the zip file
    zipFile, err := os.ReadFile("lambda_function/lambdaFunction.zip")
    if err != nil {
        return fmt.Errorf("failed to read file lambda_function/lambdaFunction.zip: %v", err)
    }
 
    // Create a new Lambda function
    input := &lambda.CreateFunctionInput{
        Code: &types.FunctionCode{
            ZipFile: zipFile,
        },
        FunctionName: aws.String(functionName),
        Handler:      aws.String("bootstrap"), // Specifies the handler for the function
        Role:         aws.String("arn:aws:iam::975050154225:role/service-role/pujitha-role-5crlj5xc"),
        Runtime:      types.RuntimeProvidedal2,
    }
 
    _, err = svc.CreateFunction(context.Background(), input)
    return err
}
 
func updateFunction(svc *lambda.Client, functionName string) error {
    // Read the function code from the zip file
    zipFile, err := os.ReadFile("lambda_function/lambdaFunction.zip")
    if err != nil {
        return fmt.Errorf("failed to read file lambda_function/lambdaFunction.zip: %v", err)
    }
 
    // Update the existing Lambda function code
    input := &lambda.UpdateFunctionCodeInput{
        FunctionName: aws.String(functionName),
        ZipFile:      zipFile,
    }
 
    result, err := svc.UpdateFunctionCode(context.Background(), input)
    if err != nil {
        return err
    }
 
    // Print the updated function ARN
    fmt.Printf("Function updated: %s\n", *result.FunctionArn)
    return nil
}
 
func invokeFunction(svc *lambda.Client, functionName string) error {
    // Read the event JSON from the file
    eventFile, err := os.ReadFile("./event.json")
    if err != nil {
        return fmt.Errorf("failed to read event.json file: %v", err)
    }
 
    // Unmarshal the event JSON
    var event map[string]interface{}
    err = json.Unmarshal(eventFile, &event)
    if err != nil {
        return fmt.Errorf("failed to unmarshal event JSON: %v", err)
    }
 
    // Modify event to add the current date
    event["detail"].(map[string]interface{})["date"] = time.Now().Format("2006-01-02")
 
    // Marshal the event JSON into a payload
    payload, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %v", err)
    }
 
    // Invoke the Lambda function
    invokeInput := &lambda.InvokeInput{
        FunctionName: aws.String(functionName),
        Payload:      payload,
    }
 
    invokeResult, err := svc.Invoke(context.Background(), invokeInput)
    if err != nil {
        return err
    }
 
    // Print the invoke result
    fmt.Printf("Lambda function invoked successfully: %s\n", string(invokeResult.Payload))
    return nil
}