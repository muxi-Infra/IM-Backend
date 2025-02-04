package pg

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	res, err := rr.GetPostCommentIds(context.Background(), "testsvc", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestReadRepo_GetCommentInfos(t *testing.T) {
	db := initdb()
	tt := &PgTable{}
	rr := &ReadRepo{db: db, tt: tt}
	res, err := rr.GetCommentInfosByID(context.Background(), "testsvc", 1)
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

func TestReadRepo_CheckPostExist(t *testing.T) {
	db := initdb()
	tt := &PgTable{}
	rr := &ReadRepo{db: db, tt: tt}
	t.Run("exist", func(t *testing.T) {
		// 测试存在的情况
		// 测试前添加id为1的post
		res := rr.CheckPostExist(context.Background(), "testsvc", 1)
		if !res {
			t.Error("should exist")
		}
	})
	t.Run("exist", func(t *testing.T) {
		// 测试存在的情况
		// 测试前删除id为2的post
		res := rr.CheckPostExist(context.Background(), "testsvc", 2)
		if res {
			t.Error("should not exist")
		}
	})
}

func TestReadRepo_GetChildCommentCnt(t *testing.T) {
	db := initdb()
	tt := &PgTable{}
	rr := &ReadRepo{db: db, tt: tt}
	t.Run("father_id == 0", func(t *testing.T) {
		res, err := rr.GetChildCommentCnt(context.Background(), "testsvc", 0)
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})

	t.Run("father_id !=0 ", func(t *testing.T) {
		res, err := rr.GetChildCommentCnt(context.Background(), "testsvc", 1)
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})
	t.Run("many father_id ", func(t *testing.T) {
		res, err := rr.GetChildCommentCnt(context.Background(), "testsvc", 0, 1)
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})

}

func TestReadRepo_GetUserIDByCommentID(t *testing.T) {
	db := initdb()
	tt := &PgTable{}
	rr := &ReadRepo{db: db, tt: tt}
	t.Run("test1", func(t *testing.T) {
		res, err := rr.GetUserIDByCommentID(context.Background(), "testsvc", 1)
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})
}

func TestReadRepo_CheckCommentExist(t *testing.T) {
	db := initdb()
	tt := &PgTable{}
	rr := &ReadRepo{db: db, tt: tt}
	t.Run("test1", func(t *testing.T) {
		res := rr.CheckCommentExist(context.Background(), "testsvc", 1, 2, 3, 4, 5, 666)
		t.Log(res)
	})
}

func TestReadRepo_GetChildCommentIDAfterCursor(t *testing.T) {
	db := initdb()
	tt := &PgTable{}
	rr := &ReadRepo{db: db, tt: tt}
	t.Run("test1", func(t *testing.T) {
		earlyTime := time.Date(1949, time.January, 1, 0, 0, 0, 0, time.UTC)
		res, err := rr.GetChildCommentIDAfterCursor(context.Background(), "testsvc", 1, earlyTime, 10)
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})
	t.Run("test2", func(t *testing.T) {
		earlyTime := time.Date(1949, time.January, 1, 0, 0, 0, 0, time.UTC)
		res, err := rr.GetChildCommentIDAfterCursor(context.Background(), "testsvc", 0, earlyTime, 10)
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})

	t.Run("test3", func(t *testing.T) {
		earlyTime := time.Date(2025, time.February, 4, 11, 55, 0, 0, time.UTC)
		res, err := rr.GetChildCommentIDAfterCursor(context.Background(), "testsvc", 0, earlyTime, 10)
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	})
}
