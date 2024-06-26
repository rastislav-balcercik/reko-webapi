/*
 * Reconvalescence Support API
 *
 * Reconvalescence support
 *
 * API version: 1.0.0
 * Contact: rasto.balcercik@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package reconvalescence

import (
    "github.com/gin-gonic/gin"
)

func AddRoutes(engine *gin.Engine) {
  group := engine.Group("/api")
  
  {
    api := newReconvalescenceTicketListAPI()
    api.addRoutes(group)
  }
  
}