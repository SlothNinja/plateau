package plateau

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (cl *Client) homeHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.User.Current(c)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	cl.Log.Debugf("cu: %#v", cu)

	c.JSON(http.StatusOK, gin.H{"cu": cu})
}
