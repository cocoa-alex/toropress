package models

import (
	"../utils"
	"database/sql"
	"fmt"
	//"github.com/coocood/goset"
	"github.com/coocood/qbs"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

const (
	dbName         = "./data/sqlite.db"
	dbUser         = "root"
	mysqlDriver    = "mymysql"
	mysqlDrvformat = "%v/%v/"
	pgDriver       = "postgres"
	pgDrvFormat    = "user=%v dbname=%v sslmode=disable"
	sqlite3Driver  = "sqlite3"
)

type User struct {
	Id            int64
	Email         string
	Password      string
	Nickname      string
	Realname      string
	Avatar        string
	Avatar_min    string
	Avatar_max    string
	Birth         time.Time
	Province      string
	City          string
	Company       string
	Address       string
	Postcode      string
	Mobile        string
	Website       string
	Sex           string
	Qq            string
	Msn           string
	Weibo         string
	Ctype         int64
	Role          int64
	Created       time.Time `qbs:"index"`
	Hotness       float64   `qbs:"index"`
	Hotup         int64
	Hotdown       int64
	Views         int64 `qbs:"index"`
	LastLoginTime time.Time
	LastLoginIp   string
	LoginCount    int64
}

//category,Pid:root
type Category struct {
	Id             int64
	Pid            int64
	Uid            int64
	Ctype          int64
	Title          string
	Content        string
	Attachment     string
	Created        time.Time `qbs:"index"`
	Hotness        float64   `qbs:"index"`
	Hotup          int64
	Hotdown        int64
	Views          int64 `qbs:"index"`
	Author         string
	NodeTime       time.Time
	NodeCount      int64
	NodeLastUserId int64
}

//node,Pid:category
type Node struct {
	Id              int64
	Pid             int64
	Uid             int64
	Ctype           int64
	Title           string
	Content         string
	Attachment      string
	Created         time.Time `qbs:"index"`
	Updated         time.Time `qbs:"index"`
	Hotness         float64   `qbs:"index"`
	Hotup           int64
	Hotdown         int64
	Views           int64 `qbs:"index"`
	Author          string
	TopicTime       time.Time
	TopicCount      int64
	TopicLastUserId int64
}

//topic,Pid:node
type Topic struct {
	Id              int64
	Cid             int64
	Nid             int64
	Uid             int64
	Ctype           int64
	Title           string
	Content         string
	Attachment      string
	Created         time.Time `qbs:"index"`
	Updated         time.Time `qbs:"index"`
	Hotness         float64   `qbs:"index"`
	Hotup           int64
	Hotdown         int64
	Views           int64 `qbs:"index"`
	Author          string
	ReplyTime       time.Time
	ReplyCount      int64
	ReplyLastUserId int64
}

//reply,Pid:topic
type Reply struct {
	Id         int64
	Uid        int64
	Pid        int64 //Topic id
	Ctype      int64
	Content    string
	Attachment string
	Created    time.Time `qbs:"index"`
	Hotness    float64   `qbs:"index"`
	Hotup      int64
	Hotdown    int64
	Views      int64 `qbs:"index"`
	Author     string
	Email      string
	Website    string
}

func ConnDb() (*qbs.Qbs, error) {
	db, err := sql.Open(sqlite3Driver, dbName)
	q := qbs.New(db, qbs.NewSqlite3())
	return q, err
}

func SetMg() (*qbs.Migration, error) {
	db, err := sql.Open(sqlite3Driver, dbName)
	mg := qbs.NewMigration(db, dbName, qbs.NewSqlite3())
	return mg, err
}

func CreateDb() bool {
	q, err := ConnDb()
	defer q.Db.Close()
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		mg, _ := SetMg()
		defer mg.Db.Close()

		mg.CreateTableIfNotExists(new(User))
		mg.CreateTableIfNotExists(new(Category))
		mg.CreateTableIfNotExists(new(Node))
		mg.CreateTableIfNotExists(new(Topic))
		mg.CreateTableIfNotExists(new(Reply))

		//用户等级划分：正数是普通用户，负数是管理员各种等级划分，为0则尚未注册
		if GetUserByRole(-1000).Role != -1000 {
			AddUser("root@insion.co", "root", utils.Encrypt_password("rootpass", nil), -1000)
			fmt.Println("Default User:root,Password:rootpass")

			if GetAllTopic(0, 0, "id") == nil {
				AddCategory("Hello Category", "This is Category!")
				AddNode("Hello Node!", "This is Node!", 1, 1)
				AddTopic("Hello World!", "This is Toropress!", 1, 1, 1)
			}
		}
		return true
	}

	return false

}

func Counts() (categorys int, nodes int, topics int, menbers int) {
	q, _ := ConnDb()
	defer q.Db.Close()

	var categoryz []*Category
	if e := q.FindAll(&categoryz); e != nil {
		categorys = 0
		fmt.Println(e)
	} else {
		categorys = len(categoryz)
	}

	var nodez []*Node
	if e := q.FindAll(&nodez); e != nil {
		nodes = 0
		fmt.Println(e)
	} else {
		nodes = len(nodez)
	}

	var topicz []*Topic
	if e := q.FindAll(&topicz); e != nil {
		topics = 0
		fmt.Println(e)
	} else {
		topics = len(topicz)
	}

	var menberz []*User
	if e := q.FindAll(&menberz); e != nil {
		menbers = 0
		fmt.Println(e)
	} else {
		menbers = len(menberz)
	}

	return categorys, nodes, topics, menbers
}

func TopicCount() (today int, this_week int, this_month int) {
	q, _ := ConnDb()
	defer q.Db.Close()
	var topict, topicw, topicm []*Topic
	k := time.Now()

	//一天之前
	d, _ := time.ParseDuration("-24h")
	t := k.Add(d)
	e := q.Where("created>?", t).FindAll(&topict)
	if e != nil {
		today = 0
		fmt.Println(e)
	} else {
		today = len(topict)
	}

	//一周之前
	w := k.Add(d * 7)
	e = q.Where("created>?", w).FindAll(&topicw)
	if e != nil {
		this_week = 0
		fmt.Println(e)
	} else {
		this_week = len(topicw)
	}

	//一月之前
	m := k.Add(d * 30)
	e = q.Where("created>?", m).FindAll(&topicm)
	if e != nil {
		this_month = 0
		fmt.Println(e)
	} else {
		this_month = len(topicm)
	}

	return today, this_week, this_month
}

func AddUser(email string, nickname string, password string, role int) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	_, err := q.Save(&User{Email: email, Nickname: nickname, Password: password, Role: int64(role), Created: time.Now()})

	return err
}

func SaveUser(usr User) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	_, e := q.Save(&usr)
	return e
}

func UpdateUser(uid int, ur User) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	_, err := q.WhereEqual("id", int64(uid)).Update(&ur)
	return err
}

func GetUser(id int) (user User) {
	q, _ := ConnDb()
	defer q.Db.Close()
	q.Where("id=?", id).Find(&user)
	return user
}

func GetUserByRole(role int) (user User) {
	q, _ := ConnDb()
	defer q.Db.Close()
	q.Where("role=?", int64(role)).Find(&user)
	return user
}

func GetAllUserByRole(role int) (user []*User) {
	q, _ := ConnDb()
	defer q.Db.Close()
	q.Where("role=?", int64(role)).OrderByDesc("id").FindAll(&user)
	return user
}

func GetUserByNickname(nickname string) (user User) {
	q, _ := ConnDb()
	defer q.Db.Close()
	q.Where("nickname=?", nickname).Find(&user)
	return user
}

func AddCategory(title string, content string) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	_, err := q.Save(&Category{Title: title, Content: content, Created: time.Now()})

	return err
}

func SaveCategory(cat Category) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	_, err := q.Save(&cat)
	return err
}

func AddNode(title string, content string, cid int, uid int) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	if _, err := q.Save(&Node{Pid: int64(cid), Title: title, Content: content, Created: time.Now()}); err != nil {
		return err
	}

	type Category struct {
		NodeTime       time.Time
		NodeCount      int64
		NodeLastUserId int64
	}

	if _, err := q.WhereEqual("id", cid).Update(&Category{NodeTime: time.Now(), NodeCount: int64(len(GetAllNodeByCid(cid, 0, 0, "id"))), NodeLastUserId: int64(uid)}); err != nil {
		return err
	}
	/*
		ctr := GetCategory(cid)
		ctr.NodeTime = time.Now()
		ctr.NodeCount = int64(len(GetAllNodeByCid(cid, 0, 0, "id")))
		ctr.NodeLastUserId = int64(uid)
		if _, err := q.Save(&ctr); err != nil {
			return err
		}
	*/
	return nil
}

func AddTopic(title string, content string, cid int, nid int, uid int) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	if _, err := q.Save(&Topic{Cid: int64(cid), Nid: int64(nid), Title: title, Content: content, Created: time.Now()}); err != nil {
		return err
	}

	type Node struct {
		TopicTime       time.Time
		TopicCount      int64
		TopicLastUserId int64
	}

	if _, err := q.WhereEqual("id", nid).Update(&Node{TopicTime: time.Now(), TopicCount: int64(len(GetAllTopicByNid(nid, 0, 0, "id"))), TopicLastUserId: int64(uid)}); err != nil {
		return err
	}
	/*
		nd := GetNode(nid)
		nd.TopicTime = time.Now()
		nd.TopicCount = int64(len(GetAllTopicByNid(nid, 0, 0, "id")))
		nd.TopicLastUserId = int64(uid)
		if _, err := q.Save(&nd); err != nil {
			return err
		}
	*/
	return nil
}

func AddReply(tid int, uid int, content string, author string, email string, website string) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	if _, err := q.Save(&Reply{Pid: int64(tid), Uid: int64(uid), Content: content, Created: time.Now(), Author: author, Email: email, Website: website}); err != nil {
		return err
	}

	type Topic struct {
		ReplyTime       time.Time
		ReplyCount      int64
		ReplyLastUserId int64
	}

	if _, err := q.WhereEqual("id", tid).Update(&Topic{ReplyTime: time.Now(), ReplyCount: int64(len(GetReplyByPid(tid, 0, 0, "id"))), ReplyLastUserId: int64(uid)}); err != nil {
		return err
	}
	/*
		tp := GetTopic(tid)
		tp.ReplyCount = int64(len(GetReplyByPid(tid, 0, 0, "id")))
		tp.ReplyTime = time.Now()
		tp.ReplyLastUserId = int64(uid)
		if _, err := q.Save(&tp); err != nil {
			return err
		}
	*/
	return nil
}

func SaveNode(nd Node) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	_, err := q.Save(&nd)
	return err
}

func DelNodePlus(nid int) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	node := GetNode(nid)
	_, err := q.Delete(&node)

	for i, v := range GetAllTopicByNid(nid, 0, 0, "id") {
		if i > 0 {
			DelTopic(int(v.Id))
			for ii, vv := range GetReplyByPid(int(v.Id), 0, 0, "id") {
				if ii > 0 {
					DelReply(int(vv.Id))
				}
			}
		}
	}

	return err
}

func DelCategory(id int) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	category := GetCategory(id)
	_, err := q.Delete(&category)

	return err
}

func DelTopic(id int) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	topic := GetTopic(id)
	_, err := q.Delete(&topic)

	return err
}

func DelNode(nid int) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	node := GetNode(nid)
	_, err := q.Delete(&node)

	return err
}

func DelReply(tid int) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	reply := GetReply(tid)
	_, err := q.Delete(&reply)

	return err
}

func GetAllCategory() (allc []*Category) {
	q, _ := ConnDb()
	defer q.Db.Close()
	q.FindAll(&allc)
	return allc
}

func GetAllNode() (alln []*Node) {
	q, _ := ConnDb()
	defer q.Db.Close()
	//q.OrderByDesc("id").FindAll(&alln)
	q.OrderByDesc("created").FindAll(&alln)
	return alln
}

func GetAllTopic(offset int, limit int, path string) (allt []*Topic) {
	q, _ := ConnDb()
	defer q.Db.Close()
	q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("created").FindAll(&allt)
	return allt
}

func GetAllNodeByCid(cid int, offset int, limit int, path string) (alln []*Node) {
	//排序首先是热值优先，然后是时间优先。
	q, _ := ConnDb()
	defer q.Db.Close()
	if cid == 0 {
		q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("topic_count").OrderByDesc("created").FindAll(&alln)
	} else {
		q.WhereEqual("pid", cid).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("topic_count").OrderByDesc("created").FindAll(&alln)
	}
	return alln
}

func GetAllTopicByNid(nodeid int, offset int, limit int, path string) (allt []*Topic) {
	//排序首先是热值优先，然后是时间优先。
	q, _ := ConnDb()
	defer q.Db.Close()
	if nodeid == 0 {
		q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
	} else {
		q.WhereEqual("nid", nodeid).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
	}
	return allt
}

func SearchTopic(content string, offset int, limit int, path string) (allt []*Topic) {
	//排序首先是热值优先，然后是时间优先。
	if content != "" {
		q, _ := ConnDb()
		defer q.Db.Close()
		keyword := "%" + content + "%"
		condition := qbs.NewCondition("title like ?", keyword).Or("content like ?", keyword)
		q.Condition(condition).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
		//q.Where("title like ?", keyword).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("created").FindAll(&allt)
		return allt
	}
	return nil
}

func GetCategory(id int) (category Category) {
	q, _ := ConnDb()
	defer q.Db.Close()
	q.Where("id=?", id).Find(&category)
	return category
}

func GetNode(id int) (node Node) {
	q, _ := ConnDb()
	defer q.Db.Close()
	q.Where("id=?", id).Find(&node)
	return node
}

func GetTopic(id int) (topic Topic) {
	q, _ := ConnDb()
	defer q.Db.Close()
	q.Where("id=?", id).Find(&topic)
	return topic
}

func SaveTopic(tp Topic) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	_, err := q.Save(&tp)
	return err
}

func UpdateCategory(cid int, cg Category) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	_, err := q.WhereEqual("id", int64(cid)).Update(&cg)
	return err
}

func UpdateNode(nid int, nd Node) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	_, err := q.WhereEqual("id", int64(nid)).Update(&nd)

	return err
}

func UpdateTopic(tid int, tp Topic) error {
	q, _ := ConnDb()
	defer q.Db.Close()
	_, err := q.WhereEqual("id", int64(tid)).Update(&tp)
	return err
}

func EditNode(nid int, cid int, uid int, title string, content string) error {
	nd := GetNode(nid)
	nd.Pid = int64(cid)
	nd.Title = title
	nd.Content = content
	nd.Updated = time.Now()
	if err := UpdateNode(nid, nd); err != nil {
		return err
	}

	q, _ := ConnDb()
	defer q.Db.Close()

	type Category struct {
		NodeTime       time.Time
		NodeCount      int64
		NodeLastUserId int64
	}

	if _, err := q.WhereEqual("id", cid).Update(&Category{NodeTime: time.Now(), NodeCount: int64(len(GetAllNodeByCid(cid, 0, 0, "id"))), NodeLastUserId: int64(uid)}); err != nil {
		return err
	}

	return nil
}

func EditTopic(tid int, nid int, cid int, uid int, title string, content string) error {
	tpc := GetTopic(tid)
	tpc.Cid = int64(cid)
	tpc.Nid = int64(nid)
	tpc.Title = title
	tpc.Content = content
	tpc.Updated = time.Now()

	if err := UpdateTopic(tid, tpc); err != nil {
		return err
	}

	q, _ := ConnDb()
	defer q.Db.Close()

	type Node struct {
		TopicTime       time.Time
		TopicCount      int64
		TopicLastUserId int64
	}

	if _, err := q.WhereEqual("id", nid).Update(&Node{TopicTime: tpc.Created, TopicCount: int64(len(GetAllTopicByNid(nid, 0, 0, "id"))), TopicLastUserId: int64(uid)}); err != nil {
		return err
	}

	return nil
}

func GetAllReply() (allr []*Reply) {
	q, _ := ConnDb()
	defer q.Db.Close()
	q.OrderByDesc("id").FindAll(&allr)
	return allr
}

func GetReply(id int) (reply Reply) {
	q, _ := ConnDb()
	defer q.Db.Close()
	q.Where("id=?", id).Find(&reply)
	return reply
}

func GetReplyByPid(tid int, offset int, limit int, path string) (allr []*Reply) {
	q, _ := ConnDb()
	defer q.Db.Close()
	if tid == 0 {
		q.Offset(offset).Limit(limit).OrderByDesc(path).FindAll(&allr)
	} else {
		//最热回复
		//q.Where("pid=?", tid).Offset(offset).Limit(limit).OrderByDesc("hotness").FindAll(&allr)
		q.WhereEqual("pid", tid).Offset(offset).Limit(limit).OrderByDesc(path).FindAll(&allr)
	}
	return allr
}

/*
func main() {

	ct()

	for i := 0; i < 100; i++ {
		AddCategory("我系标题", "我系内容啊~")
	}
	for i := 0; i < 100; i++ {
		AddUser("insion@lihuaer.com", "insion", "huhjj897857hggfgjhghsjg")
	}
	for i := 0; i < 100; i++ {
		AddNode("node title", "node content")
	}
	for i := 0; i < 100; i++ {
		AddTopic("topic title", "topic content")
	}
	for i := 0; i < 100; i++ {
		AddReply(int64(i), "a reply's content")
	}
	cc := GetAllCategory()
	for _, info := range cc {
		fmt.Println(info.Content)
	}

	c := GetCategory(1)
	fmt.Println(c.Content)

	n := GetNode(1)
	fmt.Println(n.Title)

	t := GetTopic(1)
	fmt.Println(t.Content)

	r := GetReply(1)
	fmt.Println(r.Content)

	for _, info := range GetAllCategory() {
		fmt.Println(info.Title)
	}
	for _, info := range GetAllNode() {
		fmt.Println(info.Content)
	}
	for _, info := range GetAllTopic() {
		fmt.Println(info.Title)
	}
	for _, info := range GetAllReply() {
		fmt.Println(info.Content)
	}
}
*/
