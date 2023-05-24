package main

import (
	"github.com/gin-gonic/gin"
)

const (
	msgEnter = "Entering"
	msgExit  = "Exiting"
)

func (cl Client) mlogHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	// id, err := getID(c)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// ml, err := cl.MLog.Get(c, id)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// cl.Log.Debugf("ml: %#v", ml.Messages)

	// cu, err := cl.User.Current(c)
	// if err == nil {
	// 	ml, err = cl.MLog.UpdateRead(c, ml, cu)
	// 	if err != nil {
	// 		sn.JErr(c, err)
	// 		return
	// 	}
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"messages": ml.Messages,
	// 	"unread":   0,
	// })
}

func (cl Client) mlogAddHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	// cu, err := cl.User.Current(c)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// id, err := getID(c)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// obj := struct {
	// 	Message string  `json:"message"`
	// 	Creator sn.User `json:"creator"`
	// }{}

	// err = c.ShouldBind(&obj)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// if obj.Creator.ID() != cu.ID() {
	// 	sn.JErr(c, fmt.Errorf("invalid creator: %w", sn.ErrValidation))
	// 	return
	// }

	// ml, err := cl.MLog.Get(c, id)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// m := ml.AddMessage(cu, obj.Message)
	// _, err = cl.MLog.Put(c, id, ml)
	// if err != nil {
	// 	sn.JErr(c, err)
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"message": m,
	// 	"unread":  0,
	// })
}
