package customerrors

import "errors"

var ErrCantParseDataIssue = errors.New("Can't parser or use data")
var ErrMetricNotExistIssue = errors.New("Metric with this name or type is does't exist")
