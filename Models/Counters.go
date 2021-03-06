package Models

type Counters struct {
	TorrentCount         int32 `gorm:"column:torrent_count"`
	FileCount            int32 `gorm:"column:files_count"`
	TotalFilesSize       int64 `gorm:"column:total_file_size"`
	LastScrapedId        int32 `gorm:"column:last_scraped_id" json:"-"`
	LastExtractedVideoId int32 `gorm:"column:last_titled_video_id" json:"-"`
}

func (Counters) TableName() string {
	return "realtime_counters"
}
