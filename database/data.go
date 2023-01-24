package database

import "go.mongodb.org/mongo-driver/bson"

type Words struct {
	Word        string `json:"Word" bson:"Word"`
	Description string `json:"Desc" bson:"Desc"`
}

func (w Words) NewBson() bson.M {
	return bson.M{
		"Word": w.Word,
		"Desc": w.Description,
	}
}

func (w Words) FillIn(buffer bson.M) Words {
	for k, v := range buffer {
		if str, ok := v.(string); ok && k == "Word" {
			w.Word = str
		}
		if str, ok := v.(string); ok && k == "Desc" {
			w.Description = str
		}
	}
	return w
}
