package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

func init() {
	boltDB, err := bolt.Open("./dictionary.db", 0666, nil)
	if err != nil {
		log.Fatalf("init database failed %s", err.Error())
	}
	db = boltDB
}

func TodayNeedRemember() ([]MyDictionary, error) {
	now := time.Now()
	day1BeforeNow := now.Add(-1 * 24 * time.Hour)
	day7BeforeNow := now.Add(-1 * 24 * 7 * time.Hour)

	ret := make([]MyDictionary, 0)
	{
		dics, _ := listBucket(dataKey(day1BeforeNow))
		if dics != nil {
			ret = append(ret, dics...)
		}
	}
	{
		dics, _ := listBucket(dataKey(day7BeforeNow))
		if dics != nil {
			ret = append(ret, dics...)
		}
	}

	return ret, nil
}

func listBucket(bucketName string) ([]MyDictionary, error) {
	tx, err := db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return nil, nil
	}

	ret := make([]MyDictionary, 0)
	err = bucket.ForEach(func(k, v []byte) error {
		d := new(MyDictionary)
		if err := json.NewDecoder(bytes.NewBuffer(v)).Decode(d); err != nil {
			return err
		}
		ret = append(ret, *d)
		return nil
	})

	return ret, err

}

func dataKey(t time.Time) string {
	return fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
}

func Put(dic MyDictionary) error {
	now := time.Now()
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Commit()
	bucketName := dataKey(now)
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		bucket, err = tx.CreateBucket([]byte(bucketName))
		if err != nil {
			return err
		}
	}
	data, err := json.Marshal(dic)
	if err != nil {
		return err
	}
	return bucket.Put([]byte(dic.Word), data)
}
