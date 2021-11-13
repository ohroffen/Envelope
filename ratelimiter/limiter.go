package ratelimiter

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

var initialCapacity = int64(5000)
var createRate = int64(5000)

var limiter = ratelimit.NewBucketWithQuantum(
	time.Second, initialCapacity, createRate)

func TokenRateLimiter() gin.HandlerFunc {
	// fmt.Println("token create rate:", limiter.Rate())
	// fmt.Println("available token :", limiter.Available())
	return func(c *gin.Context) {
		if limiter.TakeAvailable(1) < 1 {
			log.Printf("available token :%d", limiter.Available())
			c.AbortWithStatusJSON(http.StatusTooManyRequests, "Too Many Request")
		}
		// if limiter.TakeAvailable(1) == 0 {
		// 	log.Printf("available token :%d", limiter.Available())
		// 	context.AbortWithStatusJSON(http.StatusTooManyRequests, "Too Many Request")
		// } else {
		// 	context.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Available()))
		// 	context.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Capacity()))
		// 	context.Next()
		// }
	}
}
