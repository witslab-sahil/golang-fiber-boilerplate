package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/witslab-sahil/fiber-boilerplate/internal/temporal/workflows"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/logger"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/temporal"
	"go.temporal.io/sdk/client"
	"go.temporal.io/api/workflowservice/v1"
)

type WorkflowHandler struct {
	temporalClient *temporal.Client
	logger         logger.Logger
}

func NewWorkflowHandler(temporalClient *temporal.Client, logger logger.Logger) *WorkflowHandler {
	return &WorkflowHandler{
		temporalClient: temporalClient,
		logger:         logger,
	}
}

type StartWorkflowRequest struct {
	WorkflowID string                          `json:"workflow_id"`
	Input      workflows.UserOnboardingInput   `json:"input"`
}

func (h *WorkflowHandler) StartUserOnboarding(c *fiber.Ctx) error {
	if h.temporalClient == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Workflow service unavailable",
		})
	}

	var req StartWorkflowRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Start workflow
	options := client.StartWorkflowOptions{
		ID:        req.WorkflowID,
		TaskQueue: workflows.OnboardingTaskQueue,
	}

	we, err := h.temporalClient.GetClient().ExecuteWorkflow(c.Context(), options, workflows.UserOnboardingWorkflow, req.Input)
	if err != nil {
		h.logger.Error("Failed to start workflow: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start workflow",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"workflow_id": we.GetID(),
		"run_id":      we.GetRunID(),
	})
}

func (h *WorkflowHandler) GetWorkflowStatus(c *fiber.Ctx) error {
	if h.temporalClient == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Workflow service unavailable",
		})
	}

	workflowID := c.Params("id")
	if workflowID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Workflow ID is required",
		})
	}

	// Get workflow execution
	resp, err := h.temporalClient.GetClient().DescribeWorkflowExecution(c.Context(), workflowID, "")
	if err != nil {
		h.logger.Error("Failed to describe workflow: ", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Workflow not found",
		})
	}

	return c.JSON(fiber.Map{
		"workflow_id": workflowID,
		"status":      resp.WorkflowExecutionInfo.Status.String(),
		"start_time":  resp.WorkflowExecutionInfo.StartTime,
		"close_time":  resp.WorkflowExecutionInfo.CloseTime,
	})
}

func (h *WorkflowHandler) ListWorkflows(c *fiber.Ctx) error {
	if h.temporalClient == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Workflow service unavailable",
		})
	}

	// List workflow executions
	var workflows []fiber.Map
	
	request := &workflowservice.ListWorkflowExecutionsRequest{
		Namespace: h.temporalClient.GetNamespace(),
		PageSize:  10,
	}

	response, err := h.temporalClient.GetClient().ListWorkflow(c.Context(), request)
	if err != nil {
		h.logger.Error("Failed to list workflows: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list workflows",
		})
	}

	for _, execution := range response.Executions {
		workflows = append(workflows, fiber.Map{
			"workflow_id": execution.Execution.WorkflowId,
			"run_id":      execution.Execution.RunId,
			"status":      execution.Status.String(),
			"start_time":  execution.StartTime,
			"close_time":  execution.CloseTime,
		})
	}

	return c.JSON(fiber.Map{
		"workflows": workflows,
	})
}