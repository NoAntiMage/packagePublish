package constant

import "errors"

const (
	LongTimeSupport = "9999-12-31"

	PlanCreated          = "PlanCreated"
	PlanPublished        = "PlanPublished"
	PlanPackageComfirmed = "PackageComfirmed"
	PlanReadyToPush      = "ReadyToPush"
	PlanPushing          = "Pushing"
	PlanPushed           = "Pushed"

	PlanReceived  = "Received"
	PlanReversion = "Reversion"
	PlanAbortion  = "Abortion"
)

var (
	ErrPackageInvalid      = errors.New("ErrPackageInvalid")
	ErrRemoteNotReady      = errors.New("ErrRemoteNotReady")
	ErrLoginTokenCheckFail = errors.New("ErrLoginTokenCheckFail")
	ErrPushFailed          = errors.New("ErrPusheFailed")
	ErrChunkUploadFail     = errors.New("ErrChunkUploadFail")
)
