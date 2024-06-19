# workflow
-when an image is pushed into a repository in ECR,an eventbus is created then event rule is created, it scans the image,which also adds lambda function as a target.
-when lambda function is triggered,it displays the image tag in cloudwatch logs.
-The same image tag is stored in parameter store and mongodb. 
-Everything is done through Go code,aws console is not handled manually in this application.

# setup/helpers/eventbridge

-Eventbridge contains eventbus.go,rule.go,target.go
-It creates a event bus,rule and adds a lambda function created as target,whenever an image is pushed into repository.

# setup/helpers/parameter_store

-parameter_store.go contains how to store the image tag in it,gets mongodb credentials which are stored in a secret manager (key-value pairs to be written).
with the mongo uri,image tag is stored.

# lambda
- A lambda function is created with the cli commands and it contains a zip file which includes the bootstrap file,when the zip file is uploaded in lambda function test,imagetag will be fetched into the cloudwatch logs.

# Components Used
- AWS Elastic Container Registry (ECR): For repo creation and image push.
- AWS Secrets Manager: Utilizes secrets for storing database details,image details.(Mongo DB)
- AWS Event Bridge:Event Bus,rule creation.
- AWS Lambda:Lambda Creation
- AWS System Manager:Parameter Store.
- AWS IAM:Give access to the role for secretManagerAccess,systemManagerAccess.(Add Permissions)

# commands

project dir/setup> go run main.go
project dir/lambda> ->GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go
->zip myfunction bootstrap
->touch trust-policy.json
->aws iam create-role --role-name myfunctionExecRole --assume-role-policy-document file://trust-policy.json
->aws lambda create-function --function-name myfunction --runtime provided.al2023 --handler bootstrap --architectures x86_64 --role arn:aws:iam::975050154225:role/myfunctionExecRole --zip-file fileb://myfunction.zip


# Test

-To write testcases for checking the code.
-cursor at a function>>right click>>Go:generate unit testcases for function
-write testcases for each function.

-To run: Project Directory\your folder\go test
-To save: funcname_test.go