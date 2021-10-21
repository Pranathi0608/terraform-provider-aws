// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package cloudwatchevents

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// ListTags lists cloudwatchevents service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(conn *cloudwatchevents.CloudWatchEvents, identifier string) (tftags.KeyValueTags, error) {
	input := &cloudwatchevents.ListTagsForResourceInput{
		ResourceARN: aws.String(identifier),
	}

	output, err := conn.ListTagsForResource(input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.Tags), nil
}

// []*SERVICE.Tag handling

// Tags returns cloudwatchevents service tags.
func Tags(tags tftags.KeyValueTags) []*cloudwatchevents.Tag {
	result := make([]*cloudwatchevents.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &cloudwatchevents.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from cloudwatchevents service tags.
func KeyValueTags(tags []*cloudwatchevents.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}

// UpdateTags updates cloudwatchevents service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn *cloudwatchevents.CloudWatchEvents, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &cloudwatchevents.UntagResourceInput{
			ResourceARN: aws.String(identifier),
			TagKeys:     aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.UntagResource(input)

		if err != nil {
			return fmt.Errorf("error untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &cloudwatchevents.TagResourceInput{
			ResourceARN: aws.String(identifier),
			Tags:        Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.TagResource(input)

		if err != nil {
			return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}