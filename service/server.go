package service

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfenv"
	cftools "github.com/cloudnativego/cf-tools"
	"github.com/cloudnativego/cfmgo"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
)

func NewServer(appEnv *cfenv.App) *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	//配置 formatter 依赖注入即可 格式化 json xml yaml 格式化容器
	n := negroni.Classic() //NginxController
	mx := mux.NewRouter()  //路由容器 http.NewMux ->gin.handler echo.handl;er iris.Handler

	repo := initRepository(appEnv) //对应的初始化对应的仓库 返回对应的DbRepose inmeoriReopose fackrRopose

	initRoutes(mx, formatter, repo) //初始化路由 注入 mx formmater repo

	n.UseHandler(mx) //使用nGINX ->REQUEST->usehndr使用即可

	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, repo matchRepository) {
	mx.HandleFunc("/test", testHandler(formatter)).Methods("GET")
	mx.HandleFunc("/matches", createMatchHandler(formatter, repo)).Methods("POST")
	mx.HandleFunc("/matches", getMatchListHandler(formatter, repo)).Methods("GET")
	mx.HandleFunc("/matches/{id}", getMatchDetailsHandler(formatter, repo)).Methods("GET")
	mx.HandleFunc("/matches/{id}/moves", addMoveHandler(formatter, repo)).Methods("POST")
}
func testHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct {
			Test string
		}{"This is a test"})
	}
}
func initRepository(appEnv *cfenv.App) (repo matchRepository) {
	dbServiceURI, err := cftools.GetVCAPServiceProperty(dbServiceName, "url", appEnv)
	//如果有错误返回 或者返回对应的 dbServiceURI 为 ""
	// 则对应的打印错误信息

	//这个函数可以 都是可以进行返回对应的接口可以返回 fakeMatcherroqie 也可以返回 mfakr
	//我们只需要注入对应的repo 是接口的实例接口 返回对应的数据存储的截止即可
	if err != nil || dbServiceURI == "" {
		if err != nil {
			//这个是错误信息
			fmt.Printf("\nError retriveing databse configuration: %v\n", err)
		}
		fmt.Println("MongoDB was not detected; configurating inMeoryRepository")
		repo = NewInMemoryRepository()
		return
	}
	matchCollection := cfmgo.Connect(cfmgo.NewCollectionDialer, dbServiceURI, MatchesCollectionName)
	fmt.Printf("Connecting to MongoDB service: %s...\n", dbServiceName)
	repo = newMongoMatchRepository(matchCollection)
	return

}
