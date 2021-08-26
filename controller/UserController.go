package controller

import (
	"com.w1nd/firstgin/common"
	"com.w1nd/firstgin/dto"
	"com.w1nd/firstgin/model"
	"com.w1nd/firstgin/response"
	"com.w1nd/firstgin/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	DB := common.DB
	// 获取参数
	// 使用map获取请求的参数
	//var resultMap = make(map[string]string)
	//json.NewDecoder(c.Request.Body).Decode(&resultMap)
	// 使用struct获取请求的参数
	//var requestUser = model.User{}
	//json.NewDecoder(c.Request.Body).Decode(&requestUser)
	// 使用gin带的bind
	var requestUser = model.User{}
	c.Bind(&requestUser)

	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	log.Println(telephone)
	// 数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码必须小于6位")
		return
	}

	if len(name) == 0 {
		name = util.RandomString(10)
	}
	log.Println(name, telephone, password)
	// 判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}
	// 创建用户
	hasedPassword, err  := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "加密失败")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)

	//	发送token
	token , err:= common.ReleaseToken(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code":500,"msg": "系统错误"})
		log.Printf("token generate error : %v", err)
		return
	}
	log.Println(token)
	//	返回结果
	response.Success(c, "注册成功", gin.H{"token": token})
}

func Login(c *gin.Context) {
	db := common.GetDB()
	// 使用gin带的bind
	var requestUser = model.User{}
	c.Bind(&requestUser)

	telephone := requestUser.Telephone
	password := requestUser.Password
	log.Println(telephone)
	// 数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码必须小于6位")
		return
	}
//	判断手机号是否存在
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
//	判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code":400, "msg": "密码错误"})
		return
	}
	//	发送token
	token , err:= common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code":500,"msg": "系统错误"})
		log.Printf("token generate error : %v", err)
		return
	}
	//	返回结果
	response.Success(c, "登陆成功", gin.H{"token": token})
}

func Info(c *gin.Context)  {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}