package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/stretchr/testify/assert"
)

type MockAWSError struct {
	CodeVal    string
	MessageVal string
	OrigErrVal error
}

func (m MockAWSError) Code() string {
	return m.CodeVal
}
//ErrCodeNoSuchBucket or ErrCodeAccessDenied.
func (m MockAWSError) Message() string {
	return m.MessageVal
}
// The error message provides additional context about the error that occurred.


func (m MockAWSError) OrigErr() error {
	return m.OrigErrVal
}
// allows you to access the original error for debugging 

func (m MockAWSError) Error() string {
	return fmt.Sprintf("%s: %s", m.CodeVal, m.MessageVal)
}
//  implements the error interface, allowing the MockAWSError struct to be treated as an error.

func TestCreateOrUpdateEventBus(t *testing.T) {
	t.Run("event bus created successfully", func(t *testing.T) {
		client := mockEventbridgeClient{}
		output, _:= CreateOrUpdateEventBus(client)
		// assert.NoError(t, err)
		assert.Equal(t, output, "myBusname")
	})
	const ErrCodeResourceAlreadyExistsException = "ResourceAlreadyExistsException"

	t.Run("event bus already exists", func(t *testing.T) {
		client := mockEventbridgeClient{
			CreateEventBusErr: &MockAWSError{
				CodeVal:    ErrCodeResourceAlreadyExistsException,
				MessageVal: "ResourceAlreadyExistsException",
			},
		}
		
		output, err := CreateOrUpdateEventBus(client)
		assert.NoError(t, err)
		assert.Equal(t, output, "eventbus8")
	})

	t.Run("error creating event bus", func(t *testing.T) {
		client := mockEventbridgeClient{
			CreateEventBusErr: fmt.Errorf("client returns errors"),
		}
		_, err := CreateOrUpdateEventBus(client)
		assert.Error(t, err)
	})
}

func TestCreateRule(t *testing.T) {
	t.Run("client returning rule name", func(t *testing.T) {
		clientA := mockEventbridgeClient{}
		outputA, _ := CreateRule(clientA, "myBusname")
		assert.Equal(t, outputA, "myRuleName")
	})

	t.Run("Error creating rule:", func(t *testing.T) {
		clientA := mockEventbridgeClient{
			PutRuleErr: fmt.Errorf("Error creating rule:"),
		}
		_, err := CreateRule(clientA, "myBusname")
		//  assert.Equal(t,output,"myBusname")
		assert.Equal(t, err.Error(), "Error creating rule:")
	})
}


func TestAddTarget(t *testing.T) {
	t.Run("Error adding target to rule:", func(t *testing.T) {
		clientB := mockEventbridgeClient{}
		err := AddTarget(clientB, "myRuleName", "myBusname")
		assert.Nil(t, err)
	})

	t.Run("Error adding target to rule:", func(t *testing.T) {
		clientB := mockEventbridgeClient{
			PutTargetsErr: fmt.Errorf("Error adding target to rule:"),
		}
		err := AddTarget(clientB, "myRuleName", "myBusname")
		//	assert.Equal(t,output,"myBusname")
		// assert.NotNil(t,err)
		assert.Equal(t, err.Error(), "Error adding target to rule:")
	})
}

type mockEventbridgeClient struct {
	CreateEventBusErr error
	PutRuleErr        error
	PutTargetsErr     error
}

func (client mockEventbridgeClient) CreateEventBus(ctx context.Context, params *eventbridge.CreateEventBusInput, optFns ...func(*eventbridge.Options)) (*eventbridge.CreateEventBusOutput, error) {
	if client.CreateEventBusErr != nil {
		return nil, client.CreateEventBusErr
	}
	return &eventbridge.CreateEventBusOutput{
		EventBusArn: aws.String("arn/myBusname"),
	}, nil
}
func (client mockEventbridgeClient) PutRule(ctx context.Context, params *eventbridge.PutRuleInput, optFns ...func(*eventbridge.Options)) (*eventbridge.PutRuleOutput, error) {
	// return nil, nil
	if client.PutRuleErr != nil {
		return nil, client.PutRuleErr
	}
	arnrule := "arn/myBusname/myRuleName"
	output := &eventbridge.PutRuleOutput{
		RuleArn: aws.String(arnrule),
	}

	return output, nil
}

func (client mockEventbridgeClient) PutTargets(ctx context.Context, params *eventbridge.PutTargetsInput, optFns ...func(*eventbridge.Options)) (*eventbridge.PutTargetsOutput, error) {
	if client.PutTargetsErr != nil {
		return nil, client.PutTargetsErr
	}

	return &eventbridge.PutTargetsOutput{}, nil
}

