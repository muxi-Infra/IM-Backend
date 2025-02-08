package pg

import (
	"IM-Backend/dao"
	"IM-Backend/errcode"
	"IM-Backend/global"
	"context"
	"errors"
	"gorm.io/gorm"
)

type txKey struct{}

var (
	k = txKey{}
)

type WriteRepo struct {
	db *gorm.DB
	tt dao.TableTooler
}

func NewWriteRepo(db *gorm.DB, tt dao.TableTooler) *WriteRepo {
	return &WriteRepo{db: db, tt: tt}
}
func (t2 *WriteRepo) GetGormTx(ctx context.Context) (tx *gorm.DB) {
	//从上下文中读取*gorm.DB
	v, ok := ctx.Value(k).(*gorm.DB)
	if !ok {
		return t2.db
	}
	return v
}

func (t2 *WriteRepo) InTx(ctx context.Context, f func(ctx context.Context) error) error {

	return t2.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//将tx 存入ctx中
		ctx = context.WithValue(ctx, k, tx)
		return f(ctx)
	})
}

func (t2 *WriteRepo) Create(ctx context.Context, svc string, t dao.Table) error {
	db := t2.GetGormTx(ctx)
	if exist := t2.tt.CheckTableExist(db, t, svc); !exist {
		err := t2.tt.NewTable(db, t, svc)
		if err != nil {
			return errcode.ERRCreateTable.WrapError(err)
		}
	}
	tname := t.TableName(svc)
	res := db.Table(tname).Create(t)
	if res.Error != nil || res.RowsAffected == 0 {
		err := res.Error
		if res.RowsAffected == 0 {
			err = errors.New("create nothing")
		}
		global.Log.Errorf("create [%v] in svc:%v in db failed: %v", t, svc, err)
		return errcode.ERRCreateData.WrapError(res.Error)
	}
	return nil
}

func (t2 *WriteRepo) Update(ctx context.Context, svc string, t dao.Table) error {
	db := t2.GetGormTx(ctx)
	if exist := t2.tt.CheckTableExist(db, t, svc); !exist {
		return errcode.ERRNoTable
	}
	tname := t.TableName(svc)
	res := db.Table(tname).Updates(t)
	if res.Error != nil || res.RowsAffected == 0 {
		err := res.Error
		if res.RowsAffected == 0 {
			err = errors.New("update nothing")
		}
		global.Log.Errorf("update [%v] in svc:%v in db failed: %v", t, svc, err)
		return errcode.ERRUpdateData.WrapError(res.Error)
	}
	return nil
}

func (t2 *WriteRepo) Delete(ctx context.Context, svc string, t dao.Table, where map[string]interface{}) error {
	db := t2.GetGormTx(ctx)
	if exist := t2.tt.CheckTableExist(db, t, svc); !exist {
		return errcode.ERRNoTable
	}
	tname := t.TableName(svc)

	var res *gorm.DB
	if where == nil {
		res = db.Table(tname).Delete(t)
	} else {
		res = db.Table(tname).Where(where).Delete(t)
	}
	if res.Error != nil || res.RowsAffected == 0 {
		err := res.Error
		if res.RowsAffected == 0 {
			err = errors.New("delete nothing")
		}
		global.Log.Errorf("delete in svc:%v in db failed: %v", svc, err)
		return errcode.ERRDeleteData.WrapError(res.Error)
	}
	return nil
}
