package pg

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func initdb() *gorm.DB {
	// PostgreSQL 数据库连接配置
	dsn := "host=localhost user=chen password=12345678 dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	// 连接到 PostgreSQL 数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接到数据库失败")
	}
	return db
}

func TestReadRepo_GetPostInfos(t *testing.T) {
	db := initdb()
	tt := &PgTable{}
	rr := &ReadRepo{db: db, tt: tt}
	t.Run("exist", func(t *testing.T) {
		res, err := rr.GetPostInfos(context.Background(), "testsvc", 1)
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})
	t.Run("don't exist", func(t *testing.T) {
		res, err := rr.GetPostInfos(context.Background(), "testsvc", 2)
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})

}

func TestReadRepo_GetPostLike(t *testing.T) {
	db := initdb()
	tt := &PgTable{}
	rr := &ReadRepo{db: db, tt: tt}
	t.Run("exist", func(t *testing.T) {
		res, err := rr.GetPostLike(context.Background(), "testsvc", 1)
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})
	t.Run("don't exist", func(t *testing.T) {
		res, err := rr.GetPostLike(context.Background(), "testsvc", 2)
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})

}

func TestReadRepo_GetCommentIds(t *testing.T) {
	db := initdb()
	tt := &PgTable{}
	rr := &ReadRepo{db: db, tt: tt}
	res, err := rr.GetCommentIds(context.Background(), "testsvc", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestReadRepo_GetCommentInfos(t *testing.T) {
	db := initdb()
	tt := &PgTable{}
	rr := &ReadRepo{db: db, tt: tt}
	res, err := rr.GetCommentInfos(context.Background(), "testsvc", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestReadRepo_GetCommentLike(t *testing.T) {
	db := initdb()
	tt := &PgTable{}
	rr := &ReadRepo{db: db, tt: tt}
	res, err := rr.GetCommentLike(context.Background(), "testsvc", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestReadRepo_GetPostInfosByExtra(t *testing.T) {
	db := initdb()
	tt := &PgTable{}
	rr := &ReadRepo{db: db, tt: tt}
	t.Run("have column", func(t *testing.T) {
		res, err := rr.GetPostInfosByExtra(context.Background(), "testsvc", "image", "hello")
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})
	t.Run("don't have column", func(t *testing.T) {
		res, err := rr.GetPostInfosByExtra(context.Background(), "testsvc", "haha", "hello")
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})

}
