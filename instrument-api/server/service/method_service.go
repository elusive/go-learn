//
//  2016 October 17
//  John Gilliland [john.gilliland@rndgroup.com]
//

package service


import (
    "errors"
    "fmt"

    "github.com/gin-gonic/gin"
    "v.io/x/lib/vlog"
)

// LoadMethod function handles loading an med or its contents using the com executor.
func LoadMethod(c *gin.Context) (interface{}, *ApiError) {
    // get the instructions/path to method
    pathToMethod := c.Params.ByName("instructions")
    if len(pathToMethod) == 0 {
        return false, badRequest(errors.New("Missing instructions for method loading task."))
    }

    fmt.Print("Path to method", pathToMethod)
    vlog.Info("Path to method: ", pathToMethod)

    return true, nil
}