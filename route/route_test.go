package route

import (
	"IM-Backend/cache"
	"IM-Backend/cache/redis"
	"IM-Backend/configs"
	"IM-Backend/controller"
	"IM-Backend/dao/pg"
	"IM-Backend/pkg"
	"IM-Backend/service"
	"IM-Backend/service/identity"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

var testR *gin.Engine

func TestMain(m *testing.M) {
	//cctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	nc := configs.NewNacosConfig("../configs/config.yaml")
	ncClient := configs.NewNacosClient(nc)

	var ac configs.AppConf

	ac.InitConfig(ncClient) //初始化读取配置

	db := pg.NewDB(ac)
	client := redis.NewRedisClient(ac)

	pgtable := pg.NewPgTable()
	readRepo := pg.NewReadRepo(db, pgtable)
	writeRepo := pg.NewWriteRepo(db, pgtable)
	//dbTrashCleaner := pg.NewTrashCleaner(db, pgtable)
	//dbTrashFinder := pg.NewTrashFinder(db, pgtable)

	ider := redis.NewIDer(client)
	cacheReader := redis.NewReader(client)
	cacheWriter := redis.NewWriter(client)
	svcManager := cache.NewSvcManager(ac)

	var (
		pfp = make(chan identity.PostIdentity, 30)    //待寻找的post(pending find post)
		pfc = make(chan identity.CommentIdentity, 30) //待寻找的comment(pending find comment)
		//pdc  = make(chan identity.CommentIdentity, 30)     //待删除的comment(pending delete comment)
		//pdpl = make(chan identity.PostLikeIdentity, 30)    //待删除的post like(pending delete post like)
		//pdcl = make(chan identity.CommentLikeIdentity, 30) //待删除的comment like(pending delete comment like)
	)

	authSvc := service.NewAuthSvc(svcManager)
	postSvc := service.NewPostSvc(writeRepo, readRepo, cacheWriter, cacheReader, pfp, ac)
	commentSvc := service.NewCommentService(writeRepo, readRepo, cacheWriter, cacheReader, pfc, ac)
	//detectSvc := service.NewDetectSvc(pfp, pfc, pdc, pdpl, pdcl, dbTrashFinder, svcManager)
	//cleanSvc := service.NewCleanSvc(pdc, pdpl, pdcl, dbTrashCleaner, ac)

	postCtrl := controller.NewPostController(postSvc, ider)
	commentCtrl := controller.NewCommentController(commentSvc, ider)
	ac.AddNotifyer(postSvc, commentSvc, svcManager) //添加配置通知
	ac.StartListen(ncClient)                        //开启监听

	testR = newRoute(postCtrl, commentCtrl, authSvc)
	code := m.Run()

	// 退出
	os.Exit(code)
}

// Test for POST for /api/v1/posts/publish
func TestPublishPost(t *testing.T) {
	w := httptest.NewRecorder()
	appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
	if err != nil {
		t.Error(err)
	}
	testurl := fmt.Sprintf("/api/v1/posts/publish?appKey=%s&svc=cc", appKey)
	payload := url.Values{}
	payload.Set("title", "test title")
	payload.Set("author", "john")
	payload.Set("content", "This is the content of the test post.")
	payload.Set("extra", "{\"image\":\"hello\"}")
	req, err := http.NewRequest("POST", testurl, strings.NewReader(payload.Encode()))
	if err != nil {
		t.Error(err)
	}
	//注意别忘了设置header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testR.ServeHTTP(w, req)
	t.Log(w.Body)
}

// Test for GET for /api/v1/posts/getinfo
func TestGetPostInfo(t *testing.T) {
	w := httptest.NewRecorder()
	appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
	if err != nil {
		t.Error(err)
	}

	var postId uint64 = 7468308864811663361

	testurl := fmt.Sprintf("/api/v1/posts/getinfo?appKey=%s&svc=cc&post_id=%d", appKey, postId)
	req, err := http.NewRequest("GET", testurl, nil)
	if err != nil {
		t.Error(err)
	}
	testR.ServeHTTP(w, req)
	t.Log(w.Body)
}

// Test for PUT for /api/v1/posts/update
func TestUpdatePost(t *testing.T) {
	t.Run("only content", func(t *testing.T) {
		w := httptest.NewRecorder()
		appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
		if err != nil {
			t.Error(err)
		}
		var postId uint64 = 7468308864811663361
		userID := "john"
		testurl := fmt.Sprintf("/api/v1/posts/update?appKey=%s&svc=cc&post_id=%d&user_id=%s", appKey, postId, userID)
		payload := url.Values{}
		payload.Set("content", "haha111")
		//payload.Set("extra", "{\"image\":\"hello\"}")
		req, err := http.NewRequest("PUT", testurl, strings.NewReader(payload.Encode()))
		if err != nil {
			t.Error(err)
		}
		//注意别忘了设置header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		testR.ServeHTTP(w, req)
		t.Log(w.Body)
	})
	t.Run("only extra", func(t *testing.T) {
		w := httptest.NewRecorder()
		appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
		if err != nil {
			t.Error(err)
		}
		var postId uint64 = 7468308864811663361
		userID := "john"
		testurl := fmt.Sprintf("/api/v1/posts/update?appKey=%s&svc=cc&post_id=%d&user_id=%s", appKey, postId, userID)
		payload := url.Values{}
		//payload.Set("content", "haha")
		payload.Set("extra", "{\"image\":\"nihao\"}")
		req, err := http.NewRequest("PUT", testurl, strings.NewReader(payload.Encode()))
		if err != nil {
			t.Error(err)
		}
		//注意别忘了设置header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		testR.ServeHTTP(w, req)
		t.Log(w.Body)
	})
	t.Run("none", func(t *testing.T) {
		w := httptest.NewRecorder()
		appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
		if err != nil {
			t.Error(err)
		}
		var postId uint64 = 7468308864811663361
		userID := "john"
		testurl := fmt.Sprintf("/api/v1/posts/update?appKey=%s&svc=cc&post_id=%d&user_id=%s", appKey, postId, userID)
		payload := url.Values{}
		//payload.Set("content", "haha")
		//payload.Set("extra", "{\"image\":\"hello\"}")
		req, err := http.NewRequest("PUT", testurl, strings.NewReader(payload.Encode()))
		if err != nil {
			t.Error(err)
		}
		//注意别忘了设置header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		testR.ServeHTTP(w, req)
		t.Log(w.Body)
	})

}
