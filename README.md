# go-gin-lambda

This project demonstrates how to deploy a Go Gin application on AWS Lambda using ECR for the container image and integrate it with AWS API Gateway. The application provides two endpoints: a `GET /hello` and a `POST /echo`.

## Installation

### Docker Installation
```bash
curl -fsSL https://get.docker.com/ | sh
```

### AWS CLI v2 Installation
```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && unzip awscliv2.zip && sudo ./aws/install
aws configure
```

### Amazon ECR Credential Helper Installation
```bash
sudo apt-get install amazon-ecr-credential-helper
```

Add the following to your `~/.docker/config.json`:
```json
{
    "credHelpers": {
        "<your-account-id>.dkr.ecr.<region>.amazonaws.com": "ecr-login"
    }
}
```

## Build and Push Container Image to ECR

### Authenticate Docker to ECR
```bash
aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <your-account-id>.dkr.ecr.<region>.amazonaws.com
```

### Tag and Push the Image
```bash
docker tag go-gin-lambda:latest <your-account-id>.dkr.ecr.<region>.amazonaws.com/go-gin-lambda:latest
docker push <your-account-id>.dkr.ecr.<region>.amazonaws.com/go-gin-lambda:latest
```

## IAM Role Setup

### Create Execution Role for Lambda
```bash
aws iam create-role \
  --role-name lambda-execution-role \
  --assume-role-policy-document file://trust-policy.json
```

### Attach Basic Execution Policy
```bash
aws iam attach-role-policy \
  --role-name lambda-execution-role \
  --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

## Lambda Function Setup

### Create Lambda Function from ECR Image
```bash
aws lambda create-function \
  --function-name go-gin-lambda \
  --package-type Image \
  --code ImageUri=<your-account-id>.dkr.ecr.<region>.amazonaws.com/go-gin-lambda:latest \
  --role arn:aws:iam::<your-account-id>:role/lambda-execution-role
```

## API Gateway Setup

### Create REST API
```bash
aws apigateway create-rest-api \
  --name "GoGinLambdaAPI" \
  --description "API Gateway for GoGinLambda" \
  --region <region>
```

### Configure `/hello` Endpoint
#### Create Resource
```bash
aws apigateway create-resource \
  --rest-api-id <api-id> \
  --parent-id $(aws apigateway get-resources --rest-api-id <api-id> --query "items[?path=='/'].id" --output text) \
  --path-part hello
```

#### Add GET Method
```bash
aws apigateway put-method \
  --rest-api-id <api-id> \
  --resource-id <hello-resource-id> \
  --http-method GET \
  --authorization-type NONE
```

#### Integrate with Lambda
```bash
aws apigateway put-integration \
  --rest-api-id <api-id> \
  --resource-id <hello-resource-id> \
  --http-method GET \
  --type AWS_PROXY \
  --integration-http-method POST \
  --uri arn:aws:apigateway:<region>:lambda:path/2015-03-31/functions/arn:aws:lambda:<region>:<your-account-id>:function:go-gin-lambda/invocations
```

### Add Permission for `/hello`
```bash
aws lambda add-permission \
  --function-name go-gin-lambda \
  --statement-id rest-apigateway-permission \
  --action lambda:InvokeFunction \
  --principal apigateway.amazonaws.com \
  --source-arn "arn:aws:execute-api:<region>:<your-account-id>:<api-id>/*/GET/hello"
```

### Configure `/echo` Endpoint
#### Create Resource
```bash
aws apigateway create-resource \
  --rest-api-id <api-id> \
  --parent-id $(aws apigateway get-resources --rest-api-id <api-id> --query "items[?path=='/'].id" --output text) \
  --path-part echo
```

#### Add POST Method
```bash
aws apigateway put-method \
  --rest-api-id <api-id> \
  --resource-id <echo-resource-id> \
  --http-method POST \
  --authorization-type NONE
```

#### Integrate with Lambda
```bash
aws apigateway put-integration \
  --rest-api-id <api-id> \
  --resource-id <echo-resource-id> \
  --http-method POST \
  --type AWS_PROXY \
  --integration-http-method POST \
  --uri arn:aws:apigateway:<region>:lambda:path/2015-03-31/functions/arn:aws:lambda:<region>:<your-account-id>:function:go-gin-lambda/invocations
```

### Add Permission for `/echo`
```bash
aws lambda add-permission \
  --function-name go-gin-lambda \
  --statement-id rest-apigateway-permission-post-echo \
  --action lambda:InvokeFunction \
  --principal apigateway.amazonaws.com \
  --source-arn "arn:aws:execute-api:<region>:<your-account-id>:<api-id>/*/POST/echo"
```

### Deploy API Gateway
```bash
aws apigateway create-deployment \
  --rest-api-id <api-id> \
  --stage-name prod
```

## Testing Endpoints

### Test `/hello`
```bash
curl -X GET https://<api-id>.execute-api.<region>.amazonaws.com/prod/hello
# Expected Output: {"message":"Hello from AWS Lambda with Go Gin!"}
```

### Test `/echo`
```bash
curl -X POST https://<api-id>.execute-api.<region>.amazonaws.com/prod/echo \
  -H "Content-Type: application/json" \
  -d '{"key1": "value1", "key2": "value2"}'
# Expected Output: {"received":{"key1":"value1","key2":"value2"}}
```
