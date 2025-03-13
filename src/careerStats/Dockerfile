# Build stage
FROM golang:latest AS build
WORKDIR /app

# Copy source code
COPY src/careerStats/careerStats.go ./script.go

# Initialize Go module and get dependencies
RUN go mod init career_stats
RUN go get -u github.com/gocolly/colly/v2
RUN go get -u github.com/aws/aws-lambda-go/lambda
RUN go get -u github.com/aws/aws-sdk-go-v2
RUN go get -u github.com/aws/aws-sdk-go-v2/config
RUN go get -u github.com/aws/aws-sdk-go-v2/service/dynamodb
RUN go get -u github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue
RUN go mod tidy

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./script.go

# Lambda stage
FROM public.ecr.aws/lambda/go:1
WORKDIR ${LAMBDA_TASK_ROOT}

# Copy the binary from build stage
COPY --from=build /app/main .

# Set the CMD to your Lambda handler
CMD ["main"]