package meeting

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var meetingTopic = "豆花甜咸之争"
var meetingContent = `第一行
第二行
第三行
第四行
`
var meetingStart JSONTime = JSONTime{time.Now()}
var meetingEnd JSONTime = JSONTime{time.Now()}
var cameraURL string

func init() {
	cameraURL = os.Getenv("CAMURL")
	if cameraURL == "" {
		cameraURL = "http://10.3.141.5:8080/?action=snapshot"
	}
}

type JSONTime struct {
	time.Time
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

type JSONDuration struct {
	time.Duration
}

func (t JSONDuration) MarshalJSON() ([]byte, error) {
	dura := fmt.Sprintf(`"%d 分钟"`, int(t.Duration/time.Minute))
	return []byte(dura), nil
}

func getDuration() JSONDuration {
	dur := meetingEnd.Time.Sub(meetingEnd.Time)
	intMinute := dur / time.Minute
	if intMinute < 1 {
		dur = time.Minute
	}
	return JSONDuration{dur}
}

type Info struct {
	Id int `json:""`
	// 会议开始时间
	Start JSONTime `json:""`
	//
	End JSONTime `json:""`
	// 会议主题
	Topic string `json:""`
	// 会议持续时间
	Duration JSONDuration `json:""`
	// 会议主持人
	Host string `json:""`
	// 参会人数
	Number int `json:""`
	// 参会人员
	Member []string `json:""`
}

func List(c *gin.Context) {
	var meetings = []Info{
		{
			1,
			meetingStart,
			meetingEnd,
			meetingTopic,
			getDuration(),
			"张三",
			4,
			[]string{"张三", "李四", "赵五", "钱六"},
		},
	}
	c.JSON(200, meetings)
	return
}

func Detail(c *gin.Context) {
	id := c.Param("id")
	log.Printf("meeting id: %s\n", id)
	var meetings = Info{
		1,
		meetingStart,
		meetingEnd,
		meetingTopic,
		getDuration(),
		"张三",
		4,
		[]string{"张三", "李四", "赵五", "钱六"},
	}
	c.JSON(200, meetings)
	return
}

func Add(c *gin.Context) {
	id := c.Param("id")
	log.Printf("meeting id: %s\n", id)

	var info Info
	if err := c.BindJSON(&info); err != nil || info.Topic == "" {
		c.JSON(400, map[string]string{
			"msg": "请求错误, topic不能为空",
		})
		return
	}
	meetingTopic = info.Topic
	meetingStart = JSONTime{time.Now()}
	getImage()
	c.JSON(200, map[string]string{
		"msg": "请求成功",
	})
	return
}

func getImage() {
	resp, err := http.Get(cameraURL)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	ioutil.WriteFile("/tmp/meeting.jpg", data, 0644)
	return
}

type Log struct {
	// 会议记录
	Text string `json:""`
}

func GetRecords(c *gin.Context) {
	id := c.Param("id")
	log.Printf("meeting id: %s\n", id)
	var logs = Log{
		meetingContent,
	}
	c.JSON(200, logs)
	return
}

func PutRecords(c *gin.Context) {
	id := c.Param("id")
	log.Printf("meeting id: %s\n", id)

	var log Log
	if err := c.BindJSON(&log); err != nil || log.Text == "" {
		c.JSON(400, map[string]string{
			"msg": "请求错误，会议内容不能为空",
		})
		return
	}
	meetingContent = log.Text
	meetingEnd = JSONTime{time.Now()}
	c.JSON(200, map[string]string{
		"msg": "请求成功",
	})
	return
}
