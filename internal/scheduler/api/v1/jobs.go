package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mochivi/ctljob/internal/core"
	"github.com/mochivi/ctljob/internal/scheduler"
	"github.com/mochivi/ctljob/internal/store"
)

type JobHandler struct {
	scheduler *scheduler.Scheduler
}

func NewJobHandler(scheduler *scheduler.Scheduler) *JobHandler {
	return &JobHandler{
		scheduler: scheduler,
	}
}

func (h *JobHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var submission scheduler.JobSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.scheduler.Submit(submission); err != nil {
		writeInternalError(w)
	}

	writeJSON(w, http.StatusAccepted, "Job queued")
}

func (h *JobHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Query param 'id' is required")
		return
	}

	job, err := h.scheduler.Get(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, fmt.Sprintf("job with id %s not found", id))
			return
		}
		writeInternalError(w)
		return
	}

	writeJSON(w, http.StatusOK, job)
}

func (h *JobHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Query param 'id' is required")
		return
	}

	var update scheduler.StatusUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	changed, err := h.scheduler.UpdateStatus(id, update)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, fmt.Sprintf("job with id %s not found", id))
			return
		}
		if errors.Is(err, core.ErrInvalidTransition) {
			writeError(w, http.StatusConflict, fmt.Sprintf("%s", err))
			return
		}
		if errors.Is(err, core.ErrInvalidStatus) {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("%s", err))
			return
		}
		writeInternalError(w)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":  update.Status,
		"changed": changed,
	})
}
