package logout

import (
	"github.com/gin-gonic/gin"
)

func (s *Service) Logout(ctx *gin.Context) error {
	ctx.SetCookie(
		"token", 
		"",      
		-1,      
		"/",    
		"",      
		false,   
		true,  
	)

	return nil
}
