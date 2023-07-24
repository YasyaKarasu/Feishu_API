package feishuapi

import "github.com/sirupsen/logrus"

type FileStatistics struct {
	// 文档历史访问人数，同一用户（user_id）多次访问按一次计算。
	Uv int `json:"uv"`
	// 文档历史访问次数，同一用户（user_id）多次访问按多次计算。（注：同一用户相邻两次访问间隔在半小时内视为一次访问）
	Pv int `json:"pv"`
	// 文档历史点赞总数，若对应的文档类型不支持点赞，返回 -1
	LikeCount int `json:"like_count"`
	// 时间戳（秒）
	Timestamp int `json:"timestamp"`
}

func (c AppClient) NewStatistics(data map[string]any) *FileStatistics {
	// data有三个key: file_token, file_type, statistics
	// 其中statistics是一个map[string]any
	statisticsMap := data["statistics"].(map[string]any)
	return &FileStatistics{
		Uv:        statisticsMap["uv"].(int),
		Pv:        statisticsMap["pv"].(int),
		LikeCount: statisticsMap["like_count"].(int),
		Timestamp: statisticsMap["timestamp"].(int),
	}
}

// StatisticsGetAllInfo Get the statistics of a file
func (c AppClient) StatisticsGetAllInfo(fileToken, fileType string) *FileStatistics {
	query := make(map[string]any)
	query["file_type"] = fileType

	responseMap := c.Request("get", "open-apis/drive/v1/files/"+fileToken+"/responseMap", query, nil, nil)

	if responseMap == nil {
		logrus.WithFields(logrus.Fields{
			"FileToken": fileToken,
			"FileType":  fileType,
		}).Warn("nil responseMap return")
		return nil
	}

	return c.NewStatistics(responseMap)
}
