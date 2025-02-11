package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
)

// CreateJob handles the creation of a new job
func CreateJobH(ctx *gin.Context) {
    var job models.Job
    if err := ctx.ShouldBindJSON(&job); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := services.CreateJob(ctx, &job); err != nil {

        if err == services.ErrJobExists {
            ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Bad request", "error": err.Error()})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create job", "error": err.Error()})
        return

    }

    ctx.JSON(http.StatusCreated, job)
}


// UpdateJob updates a specific job
func UpdateJobH(ctx *gin.Context) {
    
    var updateJob models.Job

    // Extract :jobId from URL parameters and inject it to body, since in models.Job job_id is required
    jobID := ctx.Param("jobId")
    updateJob.JobID = jobID

    if err := ctx.ShouldBindJSON(&updateJob); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := services.UpdateJob(ctx, &updateJob); err != nil {

        if err == services.ErrJobDoesNotExist {
            ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Bad request", "error": err.Error()})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update job", "error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, updateJob)
}

// GetJob retrieves a specific job by ID
func GetJobByIdH(ctx *gin.Context) {
    jobId := ctx.Param("jobId")
    
    job, err := services.GetJobById(ctx, jobId)
    if err != nil {

        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, gin.H{"message": "Job not found"})
            return
        }

        ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to retrieve job", "error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, job)
}


// GetJobsByTitleH retrieves jobs by job title
func GetJobsByTitleH(ctx *gin.Context) {
    jobtitle := ctx.Param("jobtitle")
    
    jobs, err := services.GetJobsByTitle(ctx, jobtitle)
    if err != nil {

        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve jobs"})
        return
    }

    ctx.JSON(http.StatusOK, jobs)
}

// GetJobsByStatus retrieves jobs by status
func GetJobsByStatusH(ctx *gin.Context) {
    status := ctx.Param("status")
    
    jobs, err := services.GetJobsByStatus(ctx, status)
    if err != nil {

        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve jobs"})
        return
    }

    ctx.JSON(http.StatusOK, jobs)
}

// ListUserJobs retrieves all jobs for a user
func ListUserJobsH(ctx *gin.Context) {
    
    jobs, err := services.GetJobsByUserId(ctx)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve jobs"})
        return
    }

    ctx.JSON(http.StatusOK, jobs)
}

// DeleteJob deletes a specific job
func DeleteJobH(ctx *gin.Context) {
    jobId := ctx.Param("jobId")
    
    if err := services.DeleteJob(ctx, jobId); err != nil {

        if err == services.ErrJobDoesNotExist {
            ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Bad request", "error": err.Error()})
            return
        }

        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
}

