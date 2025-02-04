package main

import (
	"IM-Backend/cache"
	"IM-Backend/cache/redis"
	"IM-Backend/configs"
	"IM-Backend/controller"
	"IM-Backend/dao/pg"
	"IM-Backend/route"
	"IM-Backend/service"
	"IM-Backend/service/identity"
	"context"
	"flag"
)

var (
	flagConf string
)

func init() {
	flag.StringVar(&flagConf, "conf", "configs/conf.yaml", "config file path")
}

func main() {
	cctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nc := configs.NewNacosConfig(flagConf)
	ncClient := configs.NewNacosClient(nc)

	var ac configs.AppConf
	ac.AddNotifyer()        //添加配置通知
	ac.InitConfig(ncClient) //初始化应用配置并开启监听

	db := pg.NewDB(ac)
	client := redis.NewRedisClient(ac)

	pgtable := pg.NewPgTable()
	readRepo := pg.NewReadRepo(db, pgtable)
	writeRepo := pg.NewWriteRepo(db, pgtable)
	dbTrashCleaner := pg.NewTrashCleaner(db, pgtable)
	dbTrashFinder := pg.NewTrashFinder(db, pgtable)

	ider := redis.NewIDer(client)
	cacheReader := redis.NewReader(client)
	cacheWriter := redis.NewWriter(client)
	svcManager := cache.NewSvcManager(ac)

	var (
		pfp  = make(chan identity.PostIdentity, 30)        //待寻找的post(pending find post)
		pfc  = make(chan identity.CommentIdentity, 30)     //待寻找的comment(pending find comment)
		pdc  = make(chan identity.CommentIdentity, 30)     //待删除的comment(pending delete comment)
		pdpl = make(chan identity.PostLikeIdentity, 30)    //待删除的post like(pending delete post like)
		pdcl = make(chan identity.CommentLikeIdentity, 30) //待删除的comment like(pending delete comment like)
	)
	authSvc := service.NewAuthSvc(svcManager)
	postSvc := service.NewPostSvc(writeRepo, readRepo, cacheWriter, cacheReader, pfp, ac)
	commentSvc := service.NewCommentService(writeRepo, readRepo, cacheWriter, cacheReader, pfc, ac)
	detectSvc := service.NewDetectSvc(pfp, pfc, pdc, pdpl, pdcl, dbTrashFinder, svcManager)
	cleanSvc := service.NewCleanSvc(pdc, pdpl, pdcl, dbTrashCleaner, ac)

	postCtrl := controller.NewPostController(postSvc, ider)
	commentCtrl := controller.NewCommentController(commentSvc, ider)

	app := route.NewApp(postCtrl, commentCtrl, authSvc, detectSvc, cleanSvc)

	app.Run(cctx) //运行应用
}
