## docker install
curl -fsSL https://get.docker.com/ | sh

## aws cli v2 install
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && unzip awscliv2.zip && sudo ./aws/install
aws configure

## amazon-ecr-credential-helper install
sudo apt-get install amazon-ecr-credential-helper

add to ~/.docker/config.json
```
{
    "credHelpers": {
        "418272765255.dkr.ecr.us-east-1.amazonaws.com": "ecr-login"
    }
}
```

## login n push container image to ecr
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 418272765255.dkr.ecr.us-east-1.amazonaws.com
docker tag go-gin-lambda:latest 418272765255.dkr.ecr.us-east-1.amazonaws.com/go-gin-lambda:latest
docker push 418272765255.dkr.ecr.us-east-1.amazonaws.com/go-gin-lambda:latest

## create role n policy
aws iam create-role \
  --role-name lambda-execution-role \
  --assume-role-policy-document file://trust-policy.json
aws iam attach-role-policy \
  --role-name lambda-execution-role \
  --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

## create lambda from go gin image
aws lambda create-function \
  --function-name go-gin-lambda \
  --package-type Image \
  --code ImageUri=418272765255.dkr.ecr.us-east-1.amazonaws.com/go-gin-lambda:latest \
  --role arn:aws:iam::418272765255:role/lambda-execution-role

## create API gateway
aws apigatewayv2 create-api \
  --name "GoGinLambdaAPI" \
  --protocol-type "HTTP" \
  --target arn:aws:lambda:us-east-1:418272765255:function:go-gin-lambda

## bind API gateway to lambda
aws apigatewayv2 create-stage \
  --api-id 18mklx3y47 \
  --stage-name prod \
  --auto-deploy
aws lambda add-permission \
  --function-name go-gin-lambda \
  --statement-id apigateway-test \
  --action lambda:InvokeFunction \
  --principal apigateway.amazonaws.com \
  --source-arn "arn:aws:execute-api:us-east-1:418272765255:18mklx3y47/*/*/*"

## test
aws lambda invoke \
  --function-name go-gin-lambda \
  --payload '{"httpMethod":"GET","path":"/hello"}' \
  --cli-binary-format raw-in-base64-out \
  response.json

cat response.json

{"statusCode":200,"headers":null,"multiValueHeaders":{"Content-Type":["application/json; charset=utf-8"]},"body":"{\"message\":\"Hello from AWS Lambda with Go Gin!\"}"}