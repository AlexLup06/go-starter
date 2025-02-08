package middleware

import (
	"net/http"

	"alexlupatsiy.com/personal-website/backend/repository"
	"github.com/gin-gonic/gin"
)

// InjectDbHandle injects a transaction handle to the request's context.
// Every handler can then extract that transaction to access the database.
func InjectDbHandle(contextDb repository.ContextDb) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbCtx, err := contextDb.WithContext(c.Request.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		// The request now has a context which contains the transaction.
		c.Request = c.Request.WithContext(dbCtx)

		// Continue handling the request.
		c.Next()

		hasErrors := len(c.Errors.Errors()) > 0

		if hasErrors {
			if !contextDb.IsCommittedOrRolledBack(dbCtx) {
				rollBackErr := contextDb.Rollback(dbCtx)
				if rollBackErr != nil {
					c.Status(http.StatusInternalServerError)
					return
				}
			}
		} else {
			if !contextDb.IsCommittedOrRolledBack(dbCtx) {
				err = contextDb.Commit(dbCtx)
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return
				}
			}
		}
	}
}
