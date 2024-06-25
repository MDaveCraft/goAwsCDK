package config

import "github.com/aws/aws-cdk-go/awscdk/v2"

func Env() *awscdk.Environment {

	// default environment 
	return nil 
	
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
