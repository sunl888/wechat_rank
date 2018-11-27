package model

type Article struct {
	Id          int64  `gorm:"primary_key"`  // ID
	WxId        string `gorm:"unique_index"` // 微信文章ID
	Top         int64                        // 文章位置
	Url         string                       // 微信地址
	Name        string                       // 微信昵称
	Type        string                       // 文章类型
	Title       string                       // 标题
	Author      string                       // 作者
	Picurl      string                       // 图片链接
	Digest      string                       // 描述
	WxName      string                       // 微信账号
	ReadCount   int64                        // 阅读数
	LikeCount   int64                        // 点赞数
	CreatedAt   string                       // 发布时间
	OriginalUrl string                       // 原始微信链接
}

type ArticleStore interface {
	ArticleSave([]*Article) error
}

type ArticleService interface {
	ArticleSave() error
}
