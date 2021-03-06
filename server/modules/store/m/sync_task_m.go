package m

import (
	"github.com/irisnet/irisplorer.io/server/modules/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	DocsNmSyncTask = "sync_task"
	PageSize       = 20
)

//同步信息
type SyncTask struct {
	ChainID string    `json:"chain_id" bson:"chain_id"`
	Height  int64     `json:"height" bson:"height"`
	Time    time.Time `json:"time" bson:"time"`
	Syncing bool      `json:"syncing" bson:"syncing"`
}

func (c SyncTask) Name() string {
	return DocsNmSyncTask
}

func (c SyncTask) PkKvPair() map[string]interface{} {
	return bson.M{"chain_id": c.ChainID}
}

func (c SyncTask) Index() mgo.Index {
	return mgo.Index{
		Key:        []string{"chain_id"}, // 索引字段， 默认升序,若需降序在字段前加-
		Unique:     true,                 // 唯一索引 同mysql唯一索引
		DropDups:   false,                // 索引重复替换旧文档,Unique为true时失效
		Background: true,                 // 后台创建索引
	}
}

func QuerySyncTask() (SyncTask, error) {
	result := SyncTask{}

	query := func(c *mgo.Collection) error {
		err := c.Find(bson.M{}).One(&result)
		return err
	}

	err := store.ExecCollection(DocsNmSyncTask, query)
	return result, err
}
