package stack

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apiGateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	dynamodb "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoCdkStackProps struct {
	cdk.StackProps
}

func NewGoCdkStack(scope constructs.Construct, id string, props *GoCdkStackProps) cdk.Stack {
	var sprops cdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	// create DB table here
	table := dynamodb.NewTable(stack, jsii.String("myUserTable"), &dynamodb.TableProps{
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("username"),
			Type: dynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("userTable"),
	})

	myFunction := lambda.NewFunction(stack, jsii.String("myLambdaFunction"), &lambda.FunctionProps{
		Runtime: lambda.Runtime_PROVIDED_AL2023(),
		Code:    lambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil),
		Handler: jsii.String("main"),
	})

	table.GrantReadWriteData(myFunction)

	api := apiGateway.NewRestApi(stack, jsii.String("myAPIGateway"), &apiGateway.RestApiProps{
		DefaultCorsPreflightOptions: &apiGateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
			AllowMethods: jsii.Strings("GET", "POST", "DELETE", "PUT", "OPTIONS"),
			AllowOrigins: jsii.Strings("*"),
		},
		DeployOptions: &apiGateway.StageOptions{
			LoggingLevel: apiGateway.MethodLoggingLevel_INFO,
		},
	})

	integration := apiGateway.NewLambdaIntegration(myFunction, nil)

	// Define the routes
	registerResource := api.Root().AddResource(jsii.String("register"), nil)
	registerResource.AddMethod(jsii.String("POST"), integration, nil)

	// Define the routes
	loginResource := api.Root().AddResource(jsii.String("login"), nil)
	loginResource.AddMethod(jsii.String("POST"), integration, nil)

	protectedResource := api.Root().AddResource(jsii.String("protected"), nil)
	protectedResource.AddMethod(jsii.String("GET"), integration, nil)

	return stack
}
