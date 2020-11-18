package meeting

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

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
	dur := meetingLog.End.Sub(meetingLog.Start)
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
			JSONTime{meetingLog.Start},
			JSONTime{meetingLog.End},
			meetingLog.Topic,
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
		JSONTime{meetingLog.Start},
		JSONTime{meetingLog.End},
		meetingLog.Topic,
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
	meetingLog = Meeting{}
	meetingLog.Topic = info.Topic
	meetingLog.Start = time.Now()
	// 重置会议记录
	getImage()
	c.JSON(200, map[string]string{
		"msg": "请求成功",
	})
	return
}

func getImage() {
	cli := http.Client{Timeout: time.Second*5 }
	resp, err := cli.Get(cameraURL)
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
	ioutil.WriteFile("./meeting.jpg", data, 0644)
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
		meetingLog.GetLog(),
	}
	c.JSON(200, logs)
	return
}

// 废弃
func PutRecords(c *gin.Context) {
	id := c.Param("id")
	log.Printf("meeting id: %s\n", id)
	meetingLog.End = time.Now()
	c.JSON(200, map[string]string{
		"msg": "请求成功",
	})
	return
}

type Speaker struct {
	Name string
	Time int
}

type Words struct {
	Word string
	Time int
}

type Meeting struct {
	Topic string
	Start time.Time
	End time.Time
	Log string
	Speaker []Speaker
	Words []Words
}

func (m *Meeting)GetLog() string {
	m.Log = ""
	for _, w := range m.Words {
		name := m.GetSpeaker(w.Time)
		m.Log = fmt.Sprintf("%s\n%s: %s", m.Log, name, w.Word)
	}
	return m.Log
}

func (m *Meeting)GetSpeaker(t int) string {
	var i int
	for index, s := range m.Speaker {
		if t >= s.Time {
			i = index - 1
			continue
		}
		i = index - 1
		break
	}
	if i < 0 {
		return "未知"
	}
	return m.Speaker[i].Name
}

// 会议记录
var meetingLog = Meeting{}

func PersistentLogToDisk() {
	data, err := json.Marshal(meetingLog)
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile("meeting.data", data, 0644)
	if err != nil {
		log.Println(err)
	}
}

func LoadLogFromDisk() {
	data, err := ioutil.ReadFile("meeting.data")
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(data, &meetingLog)
	if err != nil {
		log.Println(err)
	}
}

func PutSpeaker(c *gin.Context) {
	id := c.Param("id")
	log.Printf("meeting id: %s\n", id)

	var sp Speaker
	if err := c.BindJSON(&sp); err != nil || sp.Name == "" {
		c.JSON(400, map[string]string{
			"msg": "演讲人姓名不能为空",
		})
		return
	}
	meetingLog.Speaker = append(meetingLog.Speaker, sp)
	meetingLog.End = time.Now()
	c.JSON(200, map[string]string{
		"msg": "请求成功",
	})
	return
}

func PutWords(c *gin.Context) {
	id := c.Param("id")
	log.Printf("meeting id: %s\n", id)

	var words Words
	if err := c.BindJSON(&words); err != nil || words.Word == "" {
		c.JSON(400, map[string]string{
			"msg": "演讲内容不能为空",
		})
		return
	}
	meetingLog.Words = append(meetingLog.Words, words)
	meetingLog.End = time.Now()
	c.JSON(200, map[string]string{
		"msg": "请求成功",
	})
	return
}
