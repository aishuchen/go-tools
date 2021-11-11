package s3

import (
	"errors"
)

var (
	bucketIsNotString     = errors.New("bucket is not type string")
	endpointIsNotString   = errors.New("endpoint_url is not type string")
	usePathStyleIsNotBool = errors.New("use_path_style is not type bool")
)
