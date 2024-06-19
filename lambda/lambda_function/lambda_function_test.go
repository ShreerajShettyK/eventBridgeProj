package main
 
import (
    "context"
    "errors"
    "testing"
 
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
    "github.com/aws/aws-sdk-go-v2/service/ssm"
    "github.com/stretchr/testify/assert"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)
 
// MockSSMClient is a mock implementation of the SSMClient interface.
type MockSSMClient struct{}
 
// PutParameter is the mock implementation of the PutParameter method.
func (m *MockSSMClient) PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error) {
    // Simulate successful parameter storage
    return &ssm.PutParameterOutput{}, nil
}
 
// MockSecretsManagerClient is a mock implementation of the SecretsManagerClient interface.
type MockSecretsManagerClient struct{}
 
// GetSecretValue is the mock implementation of the GetSecretValue method.
func (m *MockSecretsManagerClient) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
    // Simulate successful secret retrieval
    secretString := `{"connectionString": "mongodb+srv://task3-shreeraj:YIXZaFDnEmHXC3PS@cluster0.0elhpdy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"}`
    return &secretsmanager.GetSecretValueOutput{
        SecretString: &secretString,
    }, nil
}
 
type MockMongoCollection struct{}
 
// InsertOne is the mock implementation of the InsertOne method.
func (m *MockMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
    // Simulate successful document insertion
    return &mongo.InsertOneResult{}, nil
}
 
// MockMongoCollectionError is a mock implementation of the MongoDB collection interface that returns an error.
type MockMongoCollectionError struct{}
 
// InsertOne is the mock implementation of the InsertOne method that returns an error.
func (m *MockMongoCollectionError) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
    // Simulate an error during document insertion
    return nil, errors.New("MongoDB error")
}
 
func TestHandleRequest(t *testing.T) {
    // Create a new instance of the LambdaFunction with mock clients.
    lf := &LambdaFunction{
        SSMClient:            &MockSSMClient{},
        SecretsManagerClient: &MockSecretsManagerClient{},
        MongoCollection:      &MockMongoCollection{},
    }
 
    // Create a sample CloudWatchEvent to use in the test.
    event := events.CloudWatchEvent{
        Detail: []byte(`{"image-tags":["tag1", "tag2"],"repository-name":"test-repo"}`),
    }
 
    // Call the HandleRequest method and assert the result.
    err := lf.HandleRequest(context.Background(), event)
    assert.NoError(t, err)
}
 
// Additional mock implementation to test error scenarios
 
// MockSSMClientError is a mock implementation of the SSMClient interface that returns an error.
type MockSSMClientError struct{}
 
// PutParameter is the mock implementation of the PutParameter method.
func (m *MockSSMClientError) PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error) {
    // Simulate an error during parameter storage
    return nil, errors.New("SSM error")
}
 
func TestHandleRequestMongoError(t *testing.T) {
    // Create a new instance of the LambdaFunction with a mock MongoDB collection that returns an error.
    lf := &LambdaFunction{
        SSMClient:            &MockSSMClient{},
        SecretsManagerClient: &MockSecretsManagerClient{},
        MongoCollection:      &MockMongoCollectionError{},
    }
 
    // Create a sample CloudWatchEvent to use in the test.
    event := events.CloudWatchEvent{
        Detail: []byte(`{"image-tags":["tag1", "tag2"],"repository-name":"test-repo"}`),
    }
 
    // Call the HandleRequest method and assert the error.
    errExpected := errors.New("MongoDB error")
    err := lf.HandleRequest(context.Background(), event)
    assert.EqualError(t, err, errExpected.Error())
}
 
func TestHandleRequestSSMError(t *testing.T) {
    // Create a new instance of the LambdaFunction with a mock SSM client that returns an error.
    lf := &LambdaFunction{
        SSMClient:            &MockSSMClientError{},
        SecretsManagerClient: &MockSecretsManagerClient{},
    }
 
    // Create a sample CloudWatchEvent to use in the test.
    event := events.CloudWatchEvent{
        Detail: []byte(`{"image-tags":["tag1", "tag2"],"repository-name":"test-repo"}`),
    }
 
    // Call the HandleRequest method and assert the error.
    errExpected := errors.New("SSM error")
    err := lf.HandleRequest(context.Background(), event)
    assert.EqualError(t, err, errExpected.Error())
}
 
// MockSecretsManagerClientError is a mock implementation of the SecretsManagerClient interface that returns an error.
type MockSecretsManagerClientError struct{}
 
// GetSecretValue is the mock implementation of the GetSecretValue method.
func (m *MockSecretsManagerClientError) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
    // Simulate an error during secret retrieval
    return nil, errors.New("Secrets Manager error")
}
 
func TestHandleRequestSecretsManagerError(t *testing.T) {
    // Create a new instance of the LambdaFunction with a mock Secrets Manager client that returns an error.
    lf := &LambdaFunction{
        SSMClient:            &MockSSMClient{},
        SecretsManagerClient: &MockSecretsManagerClientError{},
    }
 
    // Create a sample CloudWatchEvent to use in the test.
    event := events.CloudWatchEvent{
        Detail: []byte(`{"image-tags":["tag1", "tag2"],"repository-name":"test-repo"}`),
    }
 
    // Call the HandleRequest method and assert the error.
    errExpected := errors.New("Secrets Manager error")
    err := lf.HandleRequest(context.Background(), event)
    assert.EqualError(t, err, errExpected.Error())
}
 
func TestHandleRequestError(t *testing.T) {
    // Create a new instance of the LambdaFunction with mock clients.
    lf := &LambdaFunction{
        SSMClient:            &MockSSMClient{},
        SecretsManagerClient: &MockSecretsManagerClient{},
    }
 
    // Create a sample CloudWatchEvent to use in the test.
    event := events.CloudWatchEvent{
        Detail: []byte(`invalid-json`),
    }
 
    // Call the HandleRequest method and assert the error.
    err := lf.HandleRequest(context.Background(), event)
    assert.Error(t, err)
}
 
func TestHandleRequestNoImageTags(t *testing.T) {
    // Create a new instance of the LambdaFunction with mock clients.
    lf := &LambdaFunction{
        SSMClient:            &MockSSMClient{},
        SecretsManagerClient: &MockSecretsManagerClient{},
    }
 
    // Create a sample CloudWatchEvent with no image tags.
    event := events.CloudWatchEvent{
        Detail: []byte(`{"image-tags":[],"repository-name":"test-repo"}`),
    }
 
    // Call the HandleRequest method and assert the error.
    err := lf.HandleRequest(context.Background(), event)
    assert.Error(t, err)
}