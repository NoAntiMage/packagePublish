package variable

import "PackageServer/dto"

var (
	FilePushJobChan = make(chan dto.PackagePushJob, 3)
)
