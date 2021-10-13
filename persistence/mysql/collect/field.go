package collect

import (
	"github.com/dashengbuqi/spiderhub/middleware/mysql"
	"github.com/go-xorm/xorm"
	"time"
)

type AppFields struct {
	Id        int64  `bson:"_id"`
	Target    int    `bson:"target"`
	TargetId  int64  `bson:"target_id"`
	Content   string `bson:"content"`
	UpdatedAt int64  `bson:"updated_at"`
	CreatedAt int64  `bson:"created_at"`
}

type field struct {
	session *xorm.Engine
}

func NewAppField() *field {
	return &field{
		session: mysql.Engine[mysql.DATABASE_SPIDERHUB],
	}
}

//更新数据
func (this *field) Modify(target int, target_id int64, content []byte) error {
	var item AppFields

	has, err := this.session.Where("target=? AND target_id=?", target, target_id).Get(&item)
	if err != nil {
		return err
	}
	tm := time.Now().Unix()
	if has {
		item.Content = string(content)
		item.UpdatedAt = tm
		_, err = this.session.Where("id=?", item.Id).Cols("content", "updated_at").Update(item)
	} else {
		item.Target = target
		item.TargetId = target_id
		item.Content = string(content)
		item.UpdatedAt = tm
		item.CreatedAt = tm
		_, err = this.session.InsertOne(item)
	}
	return err
}
