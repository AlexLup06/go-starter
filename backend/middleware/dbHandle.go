package middleware

import (
	"fmt"
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
			fmt.Printf("here with rollback\n")
			if !contextDb.IsCommittedOrRolledBack(dbCtx) {
				fmt.Printf("rolling back\n")
				rollBackErr := contextDb.Rollback(dbCtx)
				if rollBackErr != nil {
					c.Status(http.StatusInternalServerError)
					return
				}
			}
		} else {
			fmt.Printf("here with commit\n")
			if !contextDb.IsCommittedOrRolledBack(dbCtx) {

				fmt.Printf("committing\n")

				err = contextDb.Commit(dbCtx)
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return
				}
			}
		}
	}
}
