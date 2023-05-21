package APP

import (
	"myVote/APP/model"
	"myVote/APP/router"
	"myVote/APP/tools"
)

func Start() {
	defer func() {
		model.Close()
	}()
	model.New()
	tools.NewToken("dididi")
	r := router.New()
	r.Run(":80")
}
