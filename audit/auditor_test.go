package audit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/gojekfarm/proctor-engine/kubernetes"
	"github.com/gojekfarm/proctor-engine/storage"
	"github.com/gojekfarm/proctor-engine/utility"
)

func TestExecutionAuditor(t *testing.T) {
	mockStore := &storage.MockStore{}
	mockKubeClient := &kubernetes.MockClient{}
	testAuditor := New(mockStore, mockKubeClient)

	jobName := "any-job-name"
	executedJobName := "proctor-123"
	imageName := "any/image:name"
	jobArgs := map[string]string{"key": "value"}

	ctx := context.WithValue(context.Background(), utility.JobSubmissionStatusContextKey, utility.JobSubmissionSuccess)
	ctx = context.WithValue(ctx, utility.JobNameContextKey, jobName)
	ctx = context.WithValue(ctx, utility.JobSubmittedForExecutionContextKey, executedJobName)
	ctx = context.WithValue(ctx, utility.ImageNameContextKey, imageName)
	ctx = context.WithValue(ctx, utility.JobArgsContextKey, jobArgs)

	done := make(chan bool, 2)
	mockStore.On("JobsExecutionAuditLog", utility.JobSubmissionSuccess, jobName, executedJobName, imageName, jobArgs).Return(nil).Once()
	mockKubeClient.On("JobExecutionStatus", executedJobName).Return("SUCCEEDED", nil).Once()
	mockStore.On("UpdateJobsExecutionAuditLog", executedJobName, "SUCCEEDED").Return(nil).Run(func(args mock.Arguments) {}).Once().Run(func(args mock.Arguments) { done <- true })

	testAuditor.AuditJobsExecution(ctx)

	<-done
	mockStore.AssertExpectations(t)
	mockKubeClient.AssertExpectations(t)
}

func TestExecutionAuditorClientError(t *testing.T) {
	mockStore := &storage.MockStore{}
	testAuditor := New(mockStore, &kubernetes.MockClient{})

	ctx := context.WithValue(context.Background(), utility.JobSubmissionStatusContextKey, utility.JobSubmissionClientError)

	mockStore.On("JobsExecutionAuditLog", utility.JobSubmissionClientError, "", "", "", map[string]string{}).Return(nil).Once()

	testAuditor.AuditJobsExecution(ctx)

	mockStore.AssertExpectations(t)
}

func TestExecutionAuditorServerError(t *testing.T) {
	mockStore := &storage.MockStore{}
	testAuditor := New(mockStore, &kubernetes.MockClient{})

	ctx := context.WithValue(context.Background(), utility.JobSubmissionStatusContextKey, utility.JobSubmissionServerError)

	mockStore.On("JobsExecutionAuditLog", utility.JobSubmissionServerError, "", "", "", map[string]string{}).Return(nil).Once()

	testAuditor.AuditJobsExecution(ctx)

	mockStore.AssertExpectations(t)
}
