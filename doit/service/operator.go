package service

import (
	"Project/doit/form"
	"Project/doit/entity"
	"Project/doit/app"
	"Project/doit/code"
	"Project/doit/util"

	v "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/mediocregopher/radix.v2/redis"
	"io/ioutil"
	"fmt"
)

var Operator = &OperatorService{}

type OperatorService struct{}

type AreaNum struct {
	Area 	string `json:"area"`
	Count	int `json:"count"`
}

//获取文章排行
func (s *OperatorService) GetArticleTop()(res form.ArticleTopResponse,err error) {
	var art1 []entity.Article
	err = app.DB.Select().OrderBy("hot desc").Limit(10).All(&art1)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	var art2 []entity.Article
	err = app.DB.Select().OrderBy("read desc").Limit(10).All(&art2)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	for i,ar := range art1{
		res.Hot[i].ID = ar.ID
		res.Hot[i].Hot = ar.Hot
		res.Hot[i].Auth = ar.Auth
		res.Hot[i].Title = ar.Title
		res.Hot[i].Time = ar.CreateTime
	}
	for i,ar := range art2{
		res.Read[i].ID = ar.ID
		res.Read[i].Read = ar.Read
		res.Read[i].Auth = ar.Auth
		res.Read[i].Title = ar.Title
		res.Read[i].Time = ar.CreateTime
	}
	return
}

//获取用户地区分布统计
func (s *OperatorService) GetAreaStatic()(res form.AreaStatisticResponse,err error) {
	sess := app.DB.Select("area,count(*) as count").From(entity.TableUser).
		GroupBy("area").OrderBy("count desc").Limit(5)
	var counts []AreaNum
	err = sess.All(&counts)
	if err != nil {
		err = errors.Wrap(err, "fail to query article data.")
		return
	}
	var sum int = 0
	for _,cou := range counts{
		res.Area = append(res.Area,cou.Area)
		res.Array = append(res.Array,cou.Count)
		sum += cou.Count
	}
	var arts []entity.User
	err = app.DB.Select().All(&arts)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	res.Area = append(res.Area,"其他")
	res.Array = append(res.Array,len(arts)-sum)
	return
}

//获取性别时间段发文统计
func (s *OperatorService) GetGenderStatic()(res form.GenderStaticResponse,err error) {
	var genderS []entity.GenderStatistic
	err = app.DB.Select().OrderBy("times asc").All(&genderS)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	res.Female = []int{0,0,0,0,0,0}
	res.Male = []int{0,0,0,0,0,0}
	for i,gen := range genderS{
		res.Female[i] = gen.Female
		res.Male[i] = gen.Male
	}
	return
}

//获取文章类别统计信息
func (s *OperatorService) GetSortStatistic()(res form.SortStaticResponse,err error) {
	var sorts []entity.Sort
	err = app.DB.Select().Where(dbx.HashExp{"state":1}).All(&sorts)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	for _,sort := range sorts{
		res.Sorts = append(res.Sorts,sort.Name)
		res.Arry = append(res.Arry,sort.Sum)
	}
	return
}

//获取站点统计信息
func (s *OperatorService) GetStatistics()(res form.SiteStatisticResponse,err error) {
	var users []entity.User
	var arts []entity.Article
	err = app.DB.Select().All(&users)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	err = app.DB.Select().All(&arts)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	sess := app.DB.Select("count(*)").From(entity.TableArticle)
	st := time.Now().Format("2006-01-02")
	et := time.Now().Format("2006-01-02")
	fmt.Println("ccccccccc")
	fmt.Println(st)
	sess.AndWhere(dbx.NewExp("create_time>={:t1}",dbx.Params{"t1":st + " 00:00:00"})).
		AndWhere(dbx.NewExp("create_time<{:t2}",dbx.Params{"t2":et + " 23:59:59"}))
	//记录总数
	var n int
	err = sess.Row(&n)
	if err != nil {
		err = errors.Wrap(err, "fail to query arts.")
		return
	}
	fmt.Println("ccccccccc")
	fmt.Println(n)
	sess1 := app.DB.Select("count(*)").From(entity.TableUser)
	sess1.AndWhere(dbx.NewExp("create_time>={:t1}",dbx.Params{"t1":st + " 00:00:00"})).
		AndWhere(dbx.NewExp("create_time<{:t2}",dbx.Params{"t2":et + " 23:59:59"}))
	//记录总数
	var m int
	err = sess1.Row(&m)
	if err != nil {
		err = errors.Wrap(err, "fail to query arts.")
		return
	}
	var sum int = 0
	for _,ar := range arts{
		sum += ar.Read
	}

	var i int
	path := "./attachment"
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		fi,_ := ioutil.ReadDir(path+"/"+f.Name())
		i+=len(fi)
	}
	res.UserCount = len(users)
	res.ArtCount = len(arts)
	res.TodayArt = n
	res.TodayRegister = m
	res.ReadCount = sum
	res.AttachCount = i
	return
}

type Count struct {
	Count	int `json:"count"`
	Date	int `json:"date"`
}

type Monthres struct {
	Arry []int `json:"arry"`
}

//获取每个月份文章发布数
func (s *OperatorService) GetMonthArticle(year int)(res Monthres,err error) {
	sess := app.DB.Select("count(*) AS count,MONTH(create_time) as date").From(entity.TableArticle).
		Where(dbx.NewExp("YEAR(create_time)={:ct}",dbx.Params{"ct":year})).
			GroupBy("MONTH(create_time)")
	var count []Count
	err = sess.All(&count)
	if err != nil {
		err = errors.Wrap(err, "fail to query article data.")
		return
	}
	res.Arry = []int{0,0,0,0,0,0,0,0,0,0,0,0}
	for _,cou := range count{
		res.Arry[cou.Date-1] = cou.Count
	}
	return
}

//管理员登陆
func (s *OperatorService) SignIn(request form.OperatorSignInRequest) (token string, operator entity.Operator, err error) {
	err = v.ValidateStruct(&request,
		v.Field(&request.Name, v.Required),
		v.Field(&request.Password, v.Required),
		v.Field(&request.CaptchaToken, v.Required),
		v.Field(&request.CaptchaCode, v.Required),
	)
	if err != nil {
		return
	}

	err = Captcha.Validate(request.CaptchaToken, request.CaptchaCode)
	if err != nil {
		return
	}
	err = app.DB.Select().Where(dbx.HashExp{"name":request.Name}).One(&operator)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusNotFound, code.CodeUserNotExist)
			return
		}
		err = errors.Wrap(err, "fail to find user")
		return
	}
	if operator.State != entity.OperatorStateEnabled {
		err = code.New(http.StatusBadRequest,code.CodeUserDisabled)
		return
	}

	err = util.ValidatePassword([]byte(request.Password), operator.PasswordHash)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			err = code.New(http.StatusBadRequest,code.CodeUserInvalidPassword)
			return
		}
		err = errors.WithStack(err)
		return
	}

	token = RandString(32)
	err = app.Redis.Cmd("SET", "go-blog:op:sessions:"+token, operator.ID, "EX", 3600).Err
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = s.UpdateSignInTimes(operator.ID)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//更新登录时间
func (s *OperatorService) UpdateSignInTimes(operatorID string) error {
	key := app.System + ":op:" + operatorID + ":sign-in-times"
	err := app.Redis.Cmd("lpush", key, time.Now().Format("2006-01-02 15:04:05")).Err
	if err != nil {
		return errors.WithStack(err)
	}
	err = app.Redis.Cmd("ltrim", key, 0, 2).Err
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

/*//查询用户总数
func (s *OperatorService) CountUsers(cond form.QueryUserRequest) (n int,err error) {
	sess := app.DB.Select("count(*)").From(entity.TableUser)
	if cond.ID != "" {
		sess.AndWhere(dbx.HashExp{"id": cond.ID})
	}
	if string(cond.Gender) != "" {
		sess.AndWhere(dbx.HashExp{"gender": cond.Gender})
	}
	if string(cond.State) != "" {
		sess.AndWhere(dbx.HashExp{"state": cond.State})
	}
	if cond.Oder != "" {
		if cond.Oder == "1"{
			sess.OrderBy("create_time desc")	//降序
		}else{
			sess.OrderBy("create_time asc")		//升序
		}
	}
	//记录总数
	err = sess.Row(&n)
	if err != nil {
		err = errors.Wrap(err, "fail to query devices.")
		return
	}
	return
}*/

//查询用户列表
func (s *OperatorService) QueryBlogUser(cond form.QueryUserRequest) (res []entity.User,err error) {
	sess := app.DB.Select("count(*)").From(entity.TableUser)
	if cond.ID != "" {
		sess.AndWhere(dbx.HashExp{"id": cond.ID})
	}
	if string(cond.Gender) != "" && cond.Gender!=0{
		sess.AndWhere(dbx.HashExp{"gender": cond.Gender})
	}
	if string(cond.State) != "" && cond.State!=0 {
		sess.AndWhere(dbx.HashExp{"state": cond.State})
	}
	if cond.Oder != "" {
		if cond.Oder == "1"{
			sess.OrderBy("create_time desc")	//降序
		}else{
			sess.OrderBy("create_time asc")		//升序
		}
	}
	err = sess.Select().All(&res)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//更改用户状态
func (s *OperatorService) ModifyUserStatus(req entity.ModifyUserStateRequest) (user entity.User,err error) {
	err = v.ValidateStruct(&req,
		v.Field(&req.ID, v.Required),
		v.Field(&req.State, v.Required),
	)
	if err != nil {
		return
	}
	err = app.DB.Select().Where(dbx.HashExp{"id": req.ID}).One(&user)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	user.State = req.State
	err = app.DB.Model(&user).Update("State")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//删除文章
func (s *OperatorService)DeleteArticle(articleID string) (err error) {
	var art entity.Article
	err = app.DB.Select().Where(dbx.HashExp{"id": articleID}).One(&art)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	err = app.DB.Model(&art).Delete()
	if err != nil{
		return err
	}
	//同步类别统计
	var sort entity.Sort
	err = app.DB.Select().Where(dbx.HashExp{"name": art.Sort}).One(&sort)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	sort.Sum -= 1
	err = app.DB.Model(&sort).Update("Sum")
	if err != nil{
		return
	}
	return
}

//删除文章分类
func (s *OperatorService)ModifyArticleSort(sortID string,state string) (sort entity.Sort,err error)  {
	err = app.DB.Select().Where(dbx.HashExp{"id": sortID}).One(&sort)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusBadRequest, code.CodeArticleNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}
	if state == "1"{
		sort.State = entity.SortStateAble
	}else{
		sort.State = entity.SortStateEnable
	}
	err = app.DB.Model(&sort).Update("State")
	if err != nil{
		return
	}
	return
}

//获取文章分类
func (s *OperatorService)GetArticlesSorts(admin string) (sorts []entity.Sort,err error) {
	if admin == "1"{
		err = app.DB.Select().All(&sorts)
		if err != nil {
			if util.IsDBNotFound(err) {
				err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
				return
			}
			err = errors.WithStack(err)
			return
		}
	}else {
		err = app.DB.Select().Where(dbx.HashExp{"state": entity.SortStateAble}).All(&sorts)
		if err != nil {
			if util.IsDBNotFound(err) {
				err = code.New(http.StatusBadRequest, code.CodeUserNotExist)
				return
			}
			err = errors.WithStack(err)
			return
		}
	}
	return
}

//获取文章列表（条件查询）
func (s *OperatorService) GetArticlesList(req form.QueryArticleRequest) (arts []form.QueryArticleResponse,err error) {
	var res []entity.Article
	sess := app.DB.Select("*").From(entity.TableArticle)
	if req.ID != "" {
		sess.AndWhere(dbx.HashExp{"id": req.ID})
	}
	if string(req.Sort) != "" {
		sess.AndWhere(dbx.HashExp{"sort": req.Sort})
	}
	err = sess.OrderBy("create_time desc").All(&res)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if res == nil {
		res = make([]entity.Article, 0)
	}
	var art form.QueryArticleResponse
	for _,re := range res{
		art.ID = re.ID
		art.Sort = re.Sort
		art.Title = re.Title
		art.Auth = re.Auth
		art.DatetimeAware = re.DatetimeAware
		arts = append(arts,art)
	}
	return
}

//创建文章分类
func (s *OperatorService) CreateArticleSort(req form.CreateArticleSortRequest)(sort entity.Sort,err error) {
	sort.Operator = req.Name
	sort.Name = req.Sort
	sort.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	sort.State = entity.SortStateAble
	err = app.DB.Transactional(func(tx *dbx.Tx) error {
		err = tx.Model(&sort).Insert()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if util.IsDBDuplicatedErr(err) {
			err = code.New(http.StatusConflict, code.CodeArticleExist)
			return
		}
		err = errors.Wrap(err, "fail to create article")
		return
	}
	return
}

//获取登录时间
func (s *OperatorService) GetSignInTimes(operatorID string) (times []string, err error) {
	key := app.System + ":op:" + operatorID + ":sign-in-times"
	times, err = app.Redis.Cmd("lrange", key, 0, 2).List()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//验证Token
func (s *OperatorService) CheckToken(token string) (operator entity.Operator, err error) {
	key := "go-blog:op:sessions:" + token
	ID, err := app.Redis.Cmd("GET", key).Str()
	if err != nil {
		if err == redis.ErrRespNil {
			app.Logger.Info().Msg("token expired.")
			err = code.New(http.StatusBadRequest,code.CodeTokenNotExist)
			return
		}
		err = errors.WithStack(err)
		return
	}

	err = app.DB.Select().Where(dbx.HashExp{"id":ID}).One(&operator)
	if err != nil {
		if util.IsDBNotFound(err) {
			err = code.New(http.StatusNotFound, code.CodeUserNotExist)
			return
		}
		err = errors.Wrap(err, "fail to find user")
		return
	}

	if operator.State != entity.OperatorStateEnabled {
		app.Logger.Info().Msg("operator status is no enabled.")
		err = code.New(http.StatusBadRequest,code.CodeStateInvalid)
		return
	}

	err = app.Redis.Cmd("EXPIRE", key, 3600).Err
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}
