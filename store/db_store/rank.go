package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/jinzhu/gorm"
)

type dbRank struct {
	db *gorm.DB
}

func (d *dbRank) RankDetailListByRankIds(rankIds []int64, wxId, categoryId int64) (ranks []*model.RankDetailAndWechat, err error) {
	ranks = make([]*model.RankDetailAndWechat, 40)
	q := d.db.Table("rank_details r").
		Select("r.*,w.*").
		Joins("left join wechats w on r.wx_id = w.id")
	if wxId != 0 {
		q = q.Where("r.wx_id = ? ", wxId)
	}
	if categoryId != 0 {
		q = q.Where("w.category_id = ? ", categoryId)
	}
	q = q.Where("r.rank_id in (?)) order by field(r.rank_id, ?", rankIds, rankIds)
	err = q.Find(&ranks).Error
	return
}

func (d *dbRank) RankLoad(rankId int64) (rank *model.Rank, err error) {
	rank = &model.Rank{}
	err = d.db.First(&rank, "id = ?", rankId).Error
	return
}

func (d *dbRank) RankDetail(rankId, categoryId int64, limit, offset int) (ranks []*model.RankDetailAndWechat, count int64, err error) {
	ranks = make([]*model.RankDetailAndWechat, 0, limit)
	q := d.db.Table("rank_details r").
		Select("r.*,w.*").
		Joins("left join wechats w on r.wx_id = w.id ").
		Where("r.rank_id = ?", rankId)

	if categoryId != 0 {
		q = q.Where(" w.category_id = ?", categoryId)
	} else {
		// cid == 0 && 不显示私有的分类下的公众号
		q = q.Joins("left join categories c on c.id = w.category_id").
			Where("c.is_private = ?", false)
	}
	q.Count(&count)
	q = q.Omit("created_at,updated_at").Order("r.wci desc")
	if limit != 0 {
		q = q.Offset(offset).Limit(limit)
	}
	err = q.Find(&ranks).Error
	return
}

func (d *dbRank) RankList(period string) (ranks []*model.Rank, err error) {
	ranks = make([]*model.Rank, 0, 5)
	err = d.db.Order("id desc").Limit(5).Find(&ranks, "period = ?", period).Error
	return
}

func (d *dbRank) RankDetailCreate(detail *model.RankDetail) error {
	var rd model.RankDetail
	err := d.db.First(&rd, "rank_id = ? and wx_id = ?", detail.RankId, detail.WxId).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = d.db.Create(&detail).Error
		}
	}
	return err
}

func (d *dbRank) RankCreate(rank *model.Rank) error {
	var r model.Rank
	err := d.db.Model(model.Rank{}).First(&r, "start_date = ? and end_date = ?", rank.StartDate, rank.EndDate).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = d.db.Create(&rank).Error
		}
	} else {
		rank.Id = r.Id
	}
	return err
}

func NewDBRank(db *gorm.DB) model.RankStore {
	return &dbRank{db: db}
}
