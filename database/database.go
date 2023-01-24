package database

import (
	"context"
	"errors"
	"fmt"
	"goquickstart/config"
	"goquickstart/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataBase struct {
	logger.Logger
	mongoClient       *mongo.Client
	currentCollection *mongo.Collection
}

func (db *DataBase) Init(dbConfig *config.SvrConfig, lgr logger.Logger) {
	db.Logger = lgr
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConfig.DateBaseLocation))
	if err != nil {
		println(err)
		return
	}
	db.mongoClient = client
	db.currentCollection = db.mongoClient.Database(dbConfig.DbName).Collection(dbConfig.DbContainer)
	db.WriteS("connect to container")
}

// Insert 插入一条bson数据
func (db *DataBase) Insert(data interface{}) error {
	docs, ok := data.(bson.M)
	docs["_id"] = docs["Word"]
	if !ok {
		db.Write("input data is not bson, cannot be inserted")
		return errors.New("not bson, canot be inserted")
	}
	_, err := db.currentCollection.InsertOne(context.Background(), docs)
	return err
}

// FindByWord 根据单词来查找
func (db DataBase) FindByWord(word string) (Words, bool) {
	sr := db.currentCollection.FindOne(context.Background(), bson.M{"Word": word})
	temp := bson.M{}
	if err := sr.Decode(&temp); err != nil {
		return Words{}, false
	}
	return Words{}.FillIn(temp), true
}

// DeleteWord 根据单词来删除数据
func (db *DataBase) DeleteWord(word string) (Words, error) {
	// opts := options.FindOneAndUpdate().SetUpsert(true)
	sr := db.currentCollection.FindOneAndDelete(context.Background(), bson.M{"Word": word})
	temp := bson.M{}
	if err := sr.Decode(&temp); err != nil {
		return Words{}, err
	}
	db.WriteS("delete word %s from databse successfully", word)
	return Words{}.FillIn(temp), nil
}

// UpdateWord 根据单词更新数据，如果找不到对应的单词，直接插入要更新的数据
func (db *DataBase) UpdateWord(word string, newword Words) (Words, error) {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	sr := db.currentCollection.FindOneAndUpdate(context.Background(), bson.M{"Word": word},
		bson.M{"$set": bson.M{"Word": newword.Word, "Desc": newword.Description}}, opts)

	temp := bson.M{}
	if err := sr.Decode(&temp); err != nil {
		fmt.Println(err)
		return Words{}, err
	}
	db.WriteS("update word %s successfully", word)
	return Words{}.FillIn(temp), nil
}

// FindAll 返回所有单词数据
func (db DataBase) FindAll() []Words {
	cursor, err := db.currentCollection.Find(context.Background(), bson.D{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer cursor.Close(context.Background())
	result := []Words{}

	for cursor.Next(context.Background()) {
		word := Words{}
		err = cursor.Decode(&word)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		result = append(result, word)
	}
	return result
}

// func bsonUnmarshal(elements []bson.RawElement, word *Words) {
// 	for _, element := range elements {
// 		if k := element.Key(); k == "Word" {
// 			word.Word = element.Value().String()
// 		} else if k == "Desc" {
// 			word.Description = element.Value().String()
// 		}
// 	}
// }
