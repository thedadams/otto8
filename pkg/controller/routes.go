package controller

import (
	"github.com/acorn-io/baaah/pkg/conditions"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/gptscript-ai/otto/pkg/controller/handlers/agents"
	"github.com/gptscript-ai/otto/pkg/controller/handlers/runs"
	"github.com/gptscript-ai/otto/pkg/controller/handlers/threads"
	"github.com/gptscript-ai/otto/pkg/controller/handlers/workflow"
	"github.com/gptscript-ai/otto/pkg/controller/handlers/workflowexecution"
	"github.com/gptscript-ai/otto/pkg/controller/handlers/workflowstep"
	"github.com/gptscript-ai/otto/pkg/services"
	v1 "github.com/gptscript-ai/otto/pkg/storage/apis/otto.gptscript.ai/v1"
)

func routes(router *router.Router, services *services.Services) error {
	workflows := workflow.Handler{
		WorkspaceClient:   services.WorkspaceClient,
		WorkspaceProvider: "directory",
	}
	workflowExecution := workflowexecution.Handler{
		WorkspaceClient: services.WorkspaceClient,
	}
	workflowStep := workflowstep.Handler{
		Invoker: services.Invoker,
	}

	root := router.Middleware(conditions.ErrorMiddleware())

	// Runs
	root.Type(&v1.Run{}).FinalizeFunc(v1.RunFinalizer, runs.DeleteRunState)
	root.Type(&v1.Run{}).HandlerFunc(runs.Cleanup)

	// Threads
	root.Type(&v1.Thread{}).FinalizeFunc(v1.ThreadFinalizer, threads.RemoveWorkspace(services.WorkspaceClient, services.KnowledgeBin))
	root.Type(&v1.Thread{}).HandlerFunc(threads.Cleanup)
	root.Type(&v1.Thread{}).HandlerFunc(threads.IngestKnowledge(services.KnowledgeBin))

	// Workflows
	root.Type(&v1.Workflow{}).FinalizeFunc(v1.WorkflowFinalizer, workflows.Finalize)
	root.Type(&v1.Workflow{}).HandlerFunc(workflows.CreateWorkspace)

	// WorkflowExecutions
	root.Type(&v1.WorkflowExecution{}).FinalizeFunc(v1.WorkflowExecutionFinalizer, workflowExecution.Finalize)
	root.Type(&v1.WorkflowExecution{}).HandlerFunc(workflowExecution.Cleanup)
	root.Type(&v1.WorkflowExecution{}).HandlerFunc(workflowExecution.Run)

	// WorkflowSteps
	steps := root.Type(&v1.WorkflowStep{})
	steps.HandlerFunc(workflowStep.SetRunning)
	steps.HandlerFunc(workflowStep.Cleanup)

	running := steps.Middleware(workflowstep.Running)
	running.HandlerFunc(workflowStep.RunInvoke)
	running.HandlerFunc(workflowStep.RunIf)
	running.HandlerFunc(workflowStep.RunForEach)
	running.HandlerFunc(workflowStep.RunWhile)

	// Agents
	root.Type(&v1.Agent{}).FinalizeFunc(v1.AgentFinalizer, agents.RemoveWorkspaces(services.WorkspaceClient, services.KnowledgeBin))
	root.Type(&v1.Agent{}).HandlerFunc(agents.IngestKnowledge(services.KnowledgeBin))

	return nil
}