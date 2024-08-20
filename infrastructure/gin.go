package infrastructure

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func startHTTPServer(config Config, r *gin.Engine) (err error) {
	log.Printf("http server started: http://localhost:%s/webhook/event\n\n", config.Env.HttpPort)
	err = r.Run(fmt.Sprintf(":%s", config.Env.HttpPort))
	if err != nil {
		return fmt.Errorf("failed to start http server: %v", err)
	}
	return nil
}

func StartServer(config Config, r *gin.Engine) (err error) {
	//if config.UseHttps {
	//	err = startHTTPSServer(config, r)
	//} else {
	//	err = startHTTPServer(config, r)
	//}
	err = startHTTPServer(config, r)
	return
}
