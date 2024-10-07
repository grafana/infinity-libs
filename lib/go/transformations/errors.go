package transformations

import "errors"

var (
	ErrSummarizeByFieldNotFound = errors.New("summarize by field not found. Not applying summarize")
	ErrNotUniqueFieldNames = errors.New("field names are not unique. Not applying filter")
	ErrEvaluatingFilterExpression = errors.New("error evaluating filter expression")
	ErrInvalidToken = errors.New("invalid token")
	
	ErrMergeTransformationNoFrameSupplied = errors.New("no frames supplied for merge frame transformation")
	ErrMergeTransformationDifferentFields = errors.New("unable to merge fields due to different fields")
	ErrMergeTransformationDifferentFieldNames = errors.New("unable to merge field due to different field names")
	ErrMergeTransformationDifferentFieldTypes = errors.New("unable to merge fields due to different field types")
)