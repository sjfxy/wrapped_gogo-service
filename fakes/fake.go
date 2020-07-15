package fakes

import (
	"encoding/json"
	"github.com/cloudnativego/cfmgo"
	"gopkg.in/mgo.v2"
	"strconv"
)

var TargetCount int = 1

// FakeNewCollectionDialer -
//Fake将的我们注入了cintefvr 和内控处的操作一样的操作
// 返回了对应的func回调桉树注入对应的clo.conmf
// b,er := json.makr()
// if er
// reyirn
// s := &faler
// data:b
// fale
// data []byte
// mgo.clootoen
// Eorr
// 返回的colew包含即可
//
func FakeNewCollectionDialer(c interface{}) func(url, dbname, collectionname string) (col cfmgo.Collection, err error) {
	b, err := json.Marshal(c)
	if err != nil {
		panic("Unexpected Error:Unable to marshal fake data.")
	}
	return func(url, dbname, collectionname string) (col cfmgo.Collection, err error) {
		col = &FakeCollection{
			Data: b,
		}
		return
	}
}

type FakeCollection struct {
	mgo.Collection
	Data  []byte
	Error error
}

//Close
func (s *FakeCollection) Close() {

}
func (s *FakeCollection) Wake() {

}
func (s *FakeCollection) Find(params cfmgo.Params, result interface{}) (count int, err error) {
	count = TargetCount
	err = json.Unmarshal(s.Data, &result)
	return
}
func (s *FakeCollection) FindAndModify(selector interface{}, update interface{}, result interface{}) (info *mgo.ChangeInfo, err error) {
	return
}
func (s *FakeCollection) UpsertID(id interface{}, result interface{}) (changeInfo *mgo.ChangeInfo, err error) {
	var col []interface{}
	err = json.Unmarshal(s.Data, &col)
	if err != nil {
		return
	}
	col = append(col, result)
	b, err := json.Marshal(col)
	if err != nil {
		return
	}
	s.Data = b
	changeInfo = &mgo.ChangeInfo{
		Updated:    1,
		Removed:    0,
		UpsertedId: id,
	}
	return changeInfo, nil
}
func (s *FakeCollection) FindOne(id string, result interface{}) (err error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	var col []interface{}
	err = json.Unmarshal(s.Data, &col)
	if err != nil {
		return
	}
	b, err := json.Marshal(col[i])
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &result)
	return
	// fao.data
	// vatrlcol[[inyr
	// json.maurei
	//ikwem
	//mair
	// b.er colr[i[
	// er = json.un,web,reosu

}
