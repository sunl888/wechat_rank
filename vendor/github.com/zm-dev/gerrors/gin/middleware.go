package gin

import (
	"github.com/zm-dev/gerrors"
	"github.com/gin-gonic/gin"
	"encoding/json"
)

func NewHandleErrorMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execute all the handlers

		// at this point, all the handlers finished. Let's read the errors!
		// in this example we only will use the **last error typed as public**
		// but you could iterate over all them since c.Errors is a slice!
		errorToPrint := c.Errors.Last()
		if errorToPrint != nil {
			ge := &gerrors.GlobalError{}
			if json.Unmarshal([]byte(errorToPrint.Err.Error()), ge) == nil {
				if ge.ServiceName == "" {
					ge.ServiceName = serviceName
				}
				c.JSON(ge.StatusCode, gin.H{
					"code":    ge.Code,
					"message": ge.Message,
				})
			} else {
				c.JSON(500, gin.H{
					"message": errorToPrint.Error(),
				})
			}
		}
	}
}
