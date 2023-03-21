package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MongoOp is the MongoDB database
type MongoOp struct {
	URI        string
	Find       map[string]any
	DB         string
	Collection string
	Timeout    time.Duration
	scope      *Scope
	Result     any
}

// NewMongoOp is the constructor for MongoOp
func NewMongoOp(config map[string]any, scope *Scope) (*MongoOp, error) {
	config = SetDefault(config, "timeout", "10s")
	duration, err := time.ParseDuration(config["timeout"].(string))
	if err != nil {
		return nil, err
	}
	if err := PrototypeCheck(config, Proto{"URI": TYPE_STRING, "db": TYPE_STRING, "collection": TYPE_STRING,
		"find": TYPE_MAP, "timeout": TYPE_STRING}); err == nil {
		return &MongoOp{URI: config["URI"].(string), DB: config["db"].(string), Collection: config["collection"].(string),
			Find: config["find"].(map[string]any), Timeout: duration, scope: scope}, nil
	} else {
		return nil, err
	}
}

// Run will run the MongoDB query
func (o *MongoOp) Run(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	evalURI, err := o.scope.Render(ctx, o.URI)
	if err != nil {
		return err
	}
	// Connecting to MongoB
	if client, err := mongo.Connect(ctx, options.Client().ApplyURI(evalURI)); err == nil {
		defer func() {
			_ = client.Disconnect(ctx)
		}()
		evalDB, err := o.scope.Render(ctx, o.DB)
		if err != nil {
			return err
		}
		evalColl, err := o.scope.Render(ctx, o.Collection)
		if err != nil {
			return err
		}
		evalFind, err := o.scope.RenderMap(ctx, o.Find)
		if err != nil {
			return err
		}
		db := client.Database(evalDB)
		coll := db.Collection(evalColl)
		opts := options.FindOptions{}
		opts.SetMaxTime(o.Timeout)
		// Performing the Find
		if curr, err := coll.Find(ctx, evalFind, &opts); err == nil {
			defer func() {
				_ = curr.Close(ctx)
			}()
			data := make([]map[string]any, 0)
			if err := curr.Err(); err != nil {
				return err
			}
			err := curr.All(ctx, &data)
			if err != nil {
				return err
			}
			o.Result = data
		} else {
			// In case the Find fails
			return err
		}
	} else {
		// In case the connection fails
		return err
	}
	return nil
}

// GetResult returns a result
func (o *MongoOp) GetResult() any {
	return o.Result
}
