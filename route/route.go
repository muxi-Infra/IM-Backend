package route

import "github.com/gin-gonic/gin"

const port = ":8181"

type App struct {
	r *gin.Engine
}

func NewApp() *App {
	r := gin.Default()
	return &App{r: r}
}

func (a *App) Run() {
	if err := a.r.Run(port); err != nil {
		panic(err)
	}
}
