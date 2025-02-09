package route

import (
	"IM-Backend/cache"
	"IM-Backend/cache/redis"
	"IM-Backend/configs"
	"IM-Backend/controller"
	"IM-Backend/dao/pg"
	"IM-Backend/middleware"
	"IM-Backend/pkg"
	"IM-Backend/service"
	"IM-Backend/service/identity"
	"bytes"
	"encoding/json"
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

var testR = gin.Default()

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

	//加载中间件
	loadMiddleware(testR, middleware.LockMiddleware(),
		middleware.AuthMiddleware(authSvc))
	loadRoute(testR, postCtrl, commentCtrl)
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

	var postId uint64 = 7469293073747410955

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

// Test for DELETE for /api/v1/posts/delete
func TestDeletePost(t *testing.T) {
	w := httptest.NewRecorder()
	appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
	if err != nil {
		t.Error(err)
	}
	var postId uint64 = 7468309405977542658
	userID := "john"
	testurl := fmt.Sprintf("/api/v1/posts/delete?appKey=%s&svc=cc&post_id=%d&user_id=%s", appKey, postId, userID)
	req, err := http.NewRequest("DELETE", testurl, nil)
	if err != nil {
		t.Error(err)
	}
	testR.ServeHTTP(w, req)
	t.Log(w.Body)
}

// Test for PUT for /api/v1/posts/like
func TestLikePost(t *testing.T) {

	t.Run("like", func(t *testing.T) {
		w := httptest.NewRecorder()

		appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
		if err != nil {
			t.Error(err)
		}

		var postId uint64 = 7468308864811663361
		userID := "john"

		js := map[string]interface{}{
			"like": true,
		}

		jsdata, err := json.Marshal(&js)
		if err != nil {
			t.Error(err)
		}

		testurl := fmt.Sprintf("/api/v1/posts/like?appKey=%s&svc=cc&post_id=%d&user_id=%s", appKey, postId, userID)
		req, err := http.NewRequest("PUT", testurl, bytes.NewReader(jsdata))
		if err != nil {
			t.Error(err)
		}

		// 设置 Content-Type 为 application/json
		req.Header.Set("Content-Type", "application/json")

		testR.ServeHTTP(w, req)

		t.Log(w.Body)
	})
	t.Run("cancel like", func(t *testing.T) {
		w := httptest.NewRecorder()

		appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
		if err != nil {
			t.Error(err)
		}

		var postId uint64 = 7468308864811663361
		userID := "john"

		js := map[string]interface{}{
			"like": false,
		}

		jsdata, err := json.Marshal(&js)
		if err != nil {
			t.Error(err)
		}

		testurl := fmt.Sprintf("/api/v1/posts/like?appKey=%s&svc=cc&post_id=%d&user_id=%s", appKey, postId, userID)
		req, err := http.NewRequest("PUT", testurl, bytes.NewReader(jsdata))
		if err != nil {
			t.Error(err)
		}

		// 设置 Content-Type 为 application/json
		req.Header.Set("Content-Type", "application/json")

		testR.ServeHTTP(w, req)

		t.Log(w.Body)
	})
}

// Test for GET for /api/v1/posts/getlike
func TestGetPostLike(t *testing.T) {
	w := httptest.NewRecorder()
	appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
	if err != nil {
		t.Error(err)
	}
	var postId uint64 = 7468308864811663361
	testurl := fmt.Sprintf("/api/v1/posts/getlike?appKey=%s&svc=cc&post_id=%d", appKey, postId)
	req, err := http.NewRequest("GET", testurl, nil)
	if err != nil {
		t.Error(err)
	}
	testR.ServeHTTP(w, req)
	t.Log(w.Body)
}

// Test for GET for /api/v1/posts/getlist
func TestGetPostList(t *testing.T) {
	w := httptest.NewRecorder()
	appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
	if err != nil {
		t.Error(err)
	}
	cursor := time.Now().Format("2006-01-02T15:04:05")
	var limit uint = 5
	testurl := fmt.Sprintf("/api/v1/posts/getlist?appKey=%s&svc=cc&cursor=%s&limit=%d", appKey, cursor, limit)
	req, err := http.NewRequest("GET", testurl, nil)
	if err != nil {
		t.Error(err)
	}
	testR.ServeHTTP(w, req)
	t.Log(w.Body)
}

// Test for POST for /api/v1/posts/comments/publish
func TestPublishComment(t *testing.T) {
	w := httptest.NewRecorder()
	appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
	if err != nil {
		t.Error(err)
	}
	var postId uint64 = 7468308864811663361
	userID := "john"
	testurl := fmt.Sprintf("/api/v1/posts/comments/publish?appKey=%s&svc=cc&post_id=%d&user_id=%s", appKey, postId, userID)
	payload := url.Values{}
	payload.Set("content", "this is my first comment")
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

// Test for POST for /api/v1/posts/comments/reply
func TestReplyComment(t *testing.T) {

	t.Run("comment_id is right", func(t *testing.T) {
		w := httptest.NewRecorder()
		appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
		if err != nil {
			t.Error(err)
		}
		var postId uint64 = 7468308864811663361
		var commentID uint64 = 7468556624295100421
		userID := "john"
		testurl := fmt.Sprintf("/api/v1/posts/comments/reply?appKey=%s&svc=cc&post_id=%d&user_id=%s&root_comment_id=%d&father_comment_id=%d", appKey, postId, userID, commentID, commentID)
		payload := url.Values{}
		payload.Set("content", "this is my first comment")
		payload.Set("extra", "{\"image\":\"hello\"}")
		req, err := http.NewRequest("POST", testurl, strings.NewReader(payload.Encode()))
		if err != nil {
			t.Error(err)
		}
		//注意别忘了设置header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		testR.ServeHTTP(w, req)

		t.Log(w.Body)
	})
	t.Run("comment_id is wrong", func(t *testing.T) {
		w := httptest.NewRecorder()
		appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
		if err != nil {
			t.Error(err)
		}
		var postId uint64 = 7468308864811663361
		var commentID uint64 = 0
		userID := "john"
		testurl := fmt.Sprintf("/api/v1/posts/comments/reply?appKey=%s&svc=cc&post_id=%d&user_id=%s&comment_id=%d", appKey, postId, userID, commentID)
		payload := url.Values{}
		payload.Set("content", "this is my first comment")
		payload.Set("extra", "{\"image\":\"hello\"}")
		req, err := http.NewRequest("POST", testurl, strings.NewReader(payload.Encode()))
		if err != nil {
			t.Error(err)
		}
		//注意别忘了设置header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		testR.ServeHTTP(w, req)

		t.Log(w.Body)
	})
}

// Test for PUT for /api/v1/posts/comments/update
func TestUpdateComment(t *testing.T) {
	t.Run("only content", func(t *testing.T) {
		w := httptest.NewRecorder()
		appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
		if err != nil {
			t.Error(err)
		}
		var commentID uint64 = 7468556624295100421
		userID := "john"
		testurl := fmt.Sprintf("/api/v1/posts/comments/update?appKey=%s&svc=cc&user_id=%s&comment_id=%d", appKey, userID, commentID)
		payload := url.Values{}
		payload.Set("content", "update new comment")
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
		var commentID uint64 = 7468556624295100421
		userID := "john"
		testurl := fmt.Sprintf("/api/v1/posts/comments/update?appKey=%s&svc=cc&user_id=%s&comment_id=%d", appKey, userID, commentID)
		payload := url.Values{}
		//payload.Set("content", "update new comment")
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
		var commentID uint64 = 7468556624295100421
		userID := "john"
		testurl := fmt.Sprintf("/api/v1/posts/comments/update?appKey=%s&svc=cc&user_id=%s&comment_id=%d", appKey, userID, commentID)
		payload := url.Values{}
		//payload.Set("content", "update new comment")
		//payload.Set("extra", "{\"image\":\"nihao\"}")
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

// Test for DELETE for /api/v1/posts/comments/delete
func TestDeleteComment(t *testing.T) {
	w := httptest.NewRecorder()
	appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
	if err != nil {
		t.Error(err)
	}
	var commentID uint64 = 7468559583527567369
	userID := "john"
	testurl := fmt.Sprintf("/api/v1/posts/comments/delete?appKey=%s&svc=cc&user_id=%s&comment_id=%d", appKey, userID, commentID)
	req, err := http.NewRequest("DELETE", testurl, nil)
	if err != nil {
		t.Error(err)
	}
	testR.ServeHTTP(w, req)

	t.Log(w.Body)
}

// Test for GET for /api/v1/posts/comments/getinfo
func TestGetCommentInfo(t *testing.T) {
	t.Run("get root comments", func(t *testing.T) {
		w := httptest.NewRecorder()
		appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
		if err != nil {
			t.Error(err)
		}
		var postID uint64 = 7468308864811663361
		var rootCommentID uint64 = 0
		cursor := "2000-09-21T00:00:00"
		limit := 2
		testurl := fmt.Sprintf("/api/v1/posts/comments/getinfo?appKey=%s&svc=cc&post_id=%d&root_comment_id=%d&cursor=%s&limit=%d", appKey, postID, rootCommentID, cursor, limit)
		req, err := http.NewRequest("GET", testurl, nil)
		if err != nil {
			t.Error(err)
		}
		testR.ServeHTTP(w, req)

		t.Log(w.Body)
	})
	t.Run("get children comments", func(t *testing.T) {
		w := httptest.NewRecorder()
		appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
		if err != nil {
			t.Error(err)
		}
		var postID uint64 = 7468308864811663361
		var rootCommentID uint64 = 7468556624295100421
		cursor := "2000-09-21T00:00:00"
		limit := 2
		testurl := fmt.Sprintf("/api/v1/posts/comments/getinfo?appKey=%s&svc=cc&post_id=%d&root_comment_id=%d&cursor=%s&limit=%d", appKey, postID, rootCommentID, cursor, limit)
		req, err := http.NewRequest("GET", testurl, nil)
		if err != nil {
			t.Error(err)
		}
		testR.ServeHTTP(w, req)

		t.Log(w.Body)
	})
}

// Test for PUT for /api/v1/posts/comments/like
func TestLikeComment(t *testing.T) {

	t.Run("like", func(t *testing.T) {
		w := httptest.NewRecorder()

		appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
		if err != nil {
			t.Error(err)
		}

		var postID uint64 = 7468308864811663361

		var CommentID uint64 = 7468556624295100421

		userID := "mike"

		js := map[string]interface{}{
			"like": true,
		}

		jsdata, err := json.Marshal(js)
		if err != nil {
			t.Error(err)
		}

		testurl := fmt.Sprintf("/api/v1/posts/comments/like?appKey=%s&svc=cc&post_id=%d&comment_id=%d&user_id=%s", appKey, postID, CommentID, userID)

		req, err := http.NewRequest("PUT", testurl, bytes.NewReader(jsdata))

		if err != nil {
			t.Error(err)
		}

		testR.ServeHTTP(w, req)

		t.Log(w.Body)
	})
	t.Run("cancel like", func(t *testing.T) {
		w := httptest.NewRecorder()

		appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
		if err != nil {
			t.Error(err)
		}

		var postID uint64 = 7468308864811663361

		var CommentID uint64 = 7468556624295100421

		userID := "mike"

		js := map[string]interface{}{
			"like": false,
		}

		jsdata, err := json.Marshal(js)
		if err != nil {
			t.Error(err)
		}

		testurl := fmt.Sprintf("/api/v1/posts/comments/like?appKey=%s&svc=cc&post_id=%d&comment_id=%d&user_id=%s", appKey, postID, CommentID, userID)

		req, err := http.NewRequest("PUT", testurl, bytes.NewReader(jsdata))

		if err != nil {
			t.Error(err)
		}

		testR.ServeHTTP(w, req)

		t.Log(w.Body)
	})

}

// Test for GET for /api/v1/posts/comments/getlike
func TestGetCommentLike(t *testing.T) {
	w := httptest.NewRecorder()

	// 加密 appKey
	appKey, err := pkg.EncryptAES(fmt.Sprintf("%d", time.Now().Unix()), []byte("W7K8pJ3aQv2LcXgH"))
	if err != nil {
		t.Error(err)
	}

	// 构建请求的 URL
	testurl := fmt.Sprintf("/api/v1/posts/comments/getlike?appKey=%s&svc=cc", appKey)

	// 构造请求数据（改成 JSON 格式）
	payload := map[string]interface{}{
		"comment_ids": []uint64{7468556624295100421, 7468559811160834058},
	}

	// 序列化为 JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	// 创建请求对象
	req, err := http.NewRequest("GET", testurl, bytes.NewReader(jsonData))
	if err != nil {
		t.Error(err)
	}

	// 设置 Content-Type 为 application/json
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	testR.ServeHTTP(w, req)

	// 打印响应内容
	t.Log(w.Body)
}
