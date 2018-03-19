package audit

import (
	"context"

	"github.com/gojekfarm/proctor-engine/kubernetes"
	"github.com/gojekfarm/proctor-engine/logger"
	"github.com/gojekfarm/proctor-engine/storage"
	"github.com/gojekfarm/proctor-engine/utility"
)

type Auditor interface {
	AuditJobsExecution(context.Context)
}

type auditor struct {
	store      storage.Store
	kubeClient kubernetes.Client
}

func New(store storage.Store, kubeClient kubernetes.Client) Auditor {
	return &auditor{
		store:      store,
		kubeClient: kubeClient,
	}
}

func (auditor *auditor) AuditJobsExecution(ctx context.Context) {
	jobSubmissionStatus := ctx.Value(utility.JobSubmissionStatusContextKey).(string)

	if jobSubmissionStatus != utility.JobSubmissionSuccess {
		err := auditor.store.JobsExecutionAuditLog(jobSubmissionStatus, "", "", "", map[string]string{})
		if err != nil {
			logger.Error("Error auditing jobs execution", err)
		}
		return
	}
	jobName := ctx.Value(utility.JobNameContextKey).(string)
	jobSubmittedForExecution := ctx.Value(utility.JobSubmittedForExecutionContextKey).(string)
	imageName := ctx.Value(utility.ImageNameContextKey).(string)
	jobArgs := ctx.Value(utility.JobArgsContextKey).(map[string]string)

	err := auditor.store.JobsExecutionAuditLog(jobSubmissionStatus, jobName, jobSubmittedForExecution, imageName, jobArgs)
	if err != nil {
		logger.Error("Error auditing jobs execution", err)
	}

	go auditor.auditJobExecutionStatus(jobSubmittedForExecution)
}

func (auditor *auditor) auditJobExecutionStatus(jobSubmittedForExecution string) {
	status, err := auditor.kubeClient.JobExecutionStatus(jobSubmittedForExecution)
	if err != nil {
		logger.Error("Error getting job execution status", err)
	}

	err = auditor.store.UpdateJobsExecutionAuditLog(jobSubmittedForExecution, status)
	if err != nil {
		logger.Error("Error auditing jobs execution", err)
	}
}
