package stack

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	dynamodb "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	con "github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoCdkStackProps struct {
	cdk.StackProps
}

func NewGoCdkStack(scope con.Construct, id string, props *GoCdkStackProps) cdk.Stack {
	var sprops cdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	
	stack := cdk.NewStack(scope, &id, &sprops)

	table := dynamodb.NewTable(stack, jsii.String("testUserDB"), &dynamodb.TableProps{
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("username"),
			Type: dynamodb.AttributeType_STRING,	
		},
		TableName: jsii.String("user"),
	})

	lambdaFunc := lambda.NewFunction(stack, jsii.String("testFunction"), &lambda.FunctionProps{
		Runtime: lambda.Runtime_PROVIDED_AL2023(),
		Code: lambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"),nil),
		Handler: jsii.String("main"),
	})
	
	table.GrantReadWriteData(lambdaFunc)
	
	api := apigateway.NewRestApi(stack, jsii.String("testAPI"), &apigateway.RestApiProps{
		DefaultCorsPreflightOptions: &apigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type","Authorization"),
			AllowMethods: jsii.Strings("OPTIONS","GET","POST","PUT","PATCH","DELETE"),
			AllowOrigins: jsii.Strings("*"),
		},
		DeployOptions: &apigateway.StageOptions{
			LoggingLevel: apigateway.MethodLoggingLevel_INFO,
		},
	})

	integration := apigateway.NewLambdaIntegration(lambdaFunc,nil)

	registerResource := api.Root().AddResource(jsii.String("register"),nil)
	registerResource.AddMethod(jsii.String("POST"),integration,nil)

	loginResource := api.Root().AddResource(jsii.String("login"),nil)
	loginResource.AddMethod(jsii.String("POST"),integration,nil)

	protectedResource := api.Root().AddResource(jsii.String("protected"),nil)
	protectedResource.AddMethod(jsii.String("GET"),integration,nil)
	
	return stack
}