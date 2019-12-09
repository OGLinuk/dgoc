package main

import (
	"fmt"
	"net/http"
	"strings"

	cpb "github.com/OGLinuk/dgoc/collect/proto"
	"github.com/gin-gonic/gin"
)

func CollectHandler(ctx *gin.Context) {
	s := strings.Replace(ctx.Param("seed"), "+", "/", -1)
	req := &cpb.CollectRequest{Seed: s}

	if resp, err := collectClient.Collect(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"collected": resp.Collected,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to collectClient.Collect: %s", err.Error()),
		})
	}
}
