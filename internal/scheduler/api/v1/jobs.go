package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	submission := scheduler.JobSubmission{}
	json.NewDecoder(r.Body).Decode(&submission)

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
		if err == store.ErrNotFound {
			writeError(w, http.StatusNotFound, fmt.Sprintf("job with id %s not found", id))
			return
		}
		writeInternalError(w)
		return
	}

	writeJSON(w, http.StatusOK, job)
}
