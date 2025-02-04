package pg

import (
	"IM-Backend/dao"
	"IM-Backend/model/table"
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestWriteRepo_GetGormTx(t *testing.T) {
	db := &gorm.DB{}
	ctx := context.Background()
	vctx := context.WithValue(ctx, "123", db)
	wr := &WriteRepo{db: db}
	got := wr.GetGormTx(vctx)
	if got != db {
		t.Errorf("GetGormTx() = %v, want %v", got, db)
	}
}

func TestWriteRepo_Create(t *testing.T) {
	// PostgreSQL 数据库连接配置
	dsn := "host=localhost user=chen password=12345678 dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	// 连接到 PostgreSQL 数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接到数据库失败")
	}

	testTime := time.Now()
	pg := new(PgTable)
	wr := &WriteRepo{db: db, tt: pg}
	var tt dao.Table
	t.Run("create post_info", func(t *testing.T) {
		// 插入postinfo
		tt = &table.PostInfo{
			ID:      1,
			Content: "hello",
			Author:  "john",
			Extra: map[string]interface{}{
				"image": "hello",
			},
			CreatedAt: testTime,
			UpdatedAt: testTime,
		}
		err = wr.Create(context.Background(), "testsvc", tt)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("create post_like", func(t *testing.T) {
		// 插入comment
		tt = &table.PostLikeInfo{
			PostID:    1,
			UserID:    "hello",
			CreatedAt: testTime,
		}
		err = wr.Create(context.Background(), "testsvc", tt)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("create comment", func(t *testing.T) {
		tt = &table.PostCommentInfo{
			ID:       1,
			UserID:   "hello",
			FatherID: 0,
			TargetUserID: "test_user",
			PostID:   1,
			Content:  "hello world",
			Extra: map[string]interface{}{
				"image": "hello",
			},
			CreatedAt: testTime,
			UpdatedAt: testTime,
		}
		err = wr.Create(context.Background(), "testsvc", tt)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("create comment like", func(t *testing.T) {
		tt = &table.CommentLikeInfo{
			CommentID: 1,
			UserID:    "hello",
			CreatedAt: testTime,
		}
		err = wr.Create(context.Background(), "testsvc", tt)
		if err != nil {
			t.Error(err)
		}
	})

}
