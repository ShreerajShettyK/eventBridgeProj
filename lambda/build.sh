# #!/bin/bash

# # Build the Go Lambda function
# go build -o lambda_function/lambda_function lambda_function/lambda_function.go

# # Create the bootstrap file
# cat << 'EOF' > bootstrap
# #!/bin/sh

# chmod +x /var/task/lambda_function

# exec /var/task/lambda_function

# zip -r lambda_function_payload.zip lambda_function bootstrap

# # Remove the executable
# rm lambda_function/lambda_function

cd lambda_function
GOOS=linux GOARCH=amd64 go build -o bootstrap .
zip lambdaFunction.zip bootstrap

