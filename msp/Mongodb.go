package msp

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

// MongoDB 数据库常用方法
type MongoDB struct {
	Client   *mongo.Client   //连接
	Ctx      context.Context //环境
	database string          //数据库名称
	lock     sync.Mutex      //数据锁
	IsLock   bool            //是否启用锁
}

// SetDB 初始化数据库 前置条件 需要设置数据库 url
func (c *MongoDB) SetDB(url string) error {
	c.Ctx = context.TODO()
	// Set the version of the Versioned API on the client.
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(c.Ctx, clientOptions)
	if err != nil {
		return errors.New("数据库连接失败")
	}
	c.Client = client
	return nil
}

func (c *MongoDB) SetDataBase(dataBaseName string) {
	c.database = dataBaseName
}

// Insert 插入单条数据
func (c *MongoDB) Insert(Collection string, document interface{}) error {
	if c.IsLock {
		c.lock.Lock()
		defer c.lock.Unlock()
	}

	_, err := c.Client.Database(c.database).Collection(Collection).InsertOne(c.Ctx, document)
	if err != nil {
		return errors.New("插入错误,该数据已存在")
	}
	return nil
}

// InsertMany 插入多条数据
func (c *MongoDB) InsertMany(collection string, document []interface{}) error {
	if c.IsLock {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	_, err := c.Client.Database(c.database).Collection(collection).InsertMany(c.Ctx, document)
	if err != nil {
		return errors.New("插入错误,该数据已存在")
	}
	return nil
}

// UpDateOne 更新数据单条 不需要$set
func (c *MongoDB) UpDateOne(collection string, find, update interface{}) error {
	if c.IsLock {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	res := c.Client.Database(c.database).Collection(collection).FindOneAndUpdate(c.Ctx, find, bson.M{
		"$set": update,
	})
	if res.Err() != nil {
		return errors.New("数据更新失败")
	}
	return nil
}

// UpdateMany 更新数据多条不需要$set
func (c *MongoDB) UpdateMany(collection string, find, update interface{}) error {
	if c.IsLock {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	_, err := c.Client.Database(c.database).Collection(collection).UpdateMany(c.Ctx, find, bson.M{
		"$set": update,
	})
	if err != nil {
		return errors.New("数据更新失败")
	}

	return nil
}

// DeleteOne 删除单条数据
func (c *MongoDB) DeleteOne(collection string, find interface{}) (int64, error) {
	if c.IsLock {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	one, err := c.Client.Database(c.database).Collection(collection).DeleteOne(c.Ctx, find)
	if err != nil {
		return 0, errors.New("没有可删除的数据")
	}
	return one.DeletedCount, nil
}

// DeleteMany 删除多条数据
func (c *MongoDB) DeleteMany(collection string, find interface{}) (int64, error) {
	if c.IsLock {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	one, err := c.Client.Database(c.database).Collection(collection).DeleteMany(c.Ctx, find)
	if err != nil {
		return 0, errors.New("没有可删除的数据")
	}
	return one.DeletedCount, nil
}

// FindOne 查询单个信息
func (c *MongoDB) FindOne(collection string, find interface{}, Data interface{}) error {
	if c.IsLock {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	one := c.Client.Database(c.database).Collection(collection).FindOne(c.Ctx, find)
	err := one.Decode(Data)
	if err != nil {
		return errors.New("数据不存在")
	}
	return nil
}

// FindMany 查询多个信息
func (c *MongoDB) FindMany(collection string, find interface{}, limit, skip int64, Data interface{}) error {
	if c.IsLock {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	opts := options.Find().SetLimit(limit).SetSkip(skip)
	one, err := c.Client.Database(c.database).Collection(collection).Find(c.Ctx, find, opts)
	err = one.All(c.Ctx, Data)
	if err != nil {
		return errors.New("数据不存在")
	}
	return nil
}

// FindManyAll 查询多个信息
func (c *MongoDB) FindManyAll(collection string, find interface{}, Data interface{}) error {
	if c.IsLock {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	one, err := c.Client.Database(c.database).Collection(collection).Find(c.Ctx, find)
	err = one.All(c.Ctx, Data)
	if err != nil {
		return errors.New("数据不存在")
	}
	return nil
}

// Count 计数
func (c *MongoDB) Count(collection string, find interface{}) (int64, error) {
	if c.IsLock {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	documents, err := c.Client.Database(c.database).Collection(collection).CountDocuments(c.Ctx, find)
	return documents, err
}
