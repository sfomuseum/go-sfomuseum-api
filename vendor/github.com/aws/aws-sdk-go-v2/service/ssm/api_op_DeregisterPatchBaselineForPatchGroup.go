// Code generated by smithy-go-codegen DO NOT EDIT.

package ssm

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Removes a patch group from a patch baseline.
func (c *Client) DeregisterPatchBaselineForPatchGroup(ctx context.Context, params *DeregisterPatchBaselineForPatchGroupInput, optFns ...func(*Options)) (*DeregisterPatchBaselineForPatchGroupOutput, error) {
	if params == nil {
		params = &DeregisterPatchBaselineForPatchGroupInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DeregisterPatchBaselineForPatchGroup", params, optFns, c.addOperationDeregisterPatchBaselineForPatchGroupMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DeregisterPatchBaselineForPatchGroupOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DeregisterPatchBaselineForPatchGroupInput struct {

	// The ID of the patch baseline to deregister the patch group from.
	//
	// This member is required.
	BaselineId *string

	// The name of the patch group that should be deregistered from the patch baseline.
	//
	// This member is required.
	PatchGroup *string

	noSmithyDocumentSerde
}

type DeregisterPatchBaselineForPatchGroupOutput struct {

	// The ID of the patch baseline the patch group was deregistered from.
	BaselineId *string

	// The name of the patch group deregistered from the patch baseline.
	PatchGroup *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDeregisterPatchBaselineForPatchGroupMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpDeregisterPatchBaselineForPatchGroup{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpDeregisterPatchBaselineForPatchGroup{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addOpDeregisterPatchBaselineForPatchGroupValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDeregisterPatchBaselineForPatchGroup(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opDeregisterPatchBaselineForPatchGroup(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ssm",
		OperationName: "DeregisterPatchBaselineForPatchGroup",
	}
}
