package controller

import (
	"com.w1nd/firstgin/common"
	"com.w1nd/firstgin/model"
	"com.w1nd/firstgin/response"
	"com.w1nd/firstgin/vo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

type IPostController interface {
	RestController
}

type PostController struct {
	DB *gorm.DB
}

func (p PostController) Create(c *gin.Context) {
	var requestPost vo.CreatePostRequest
	// 数据验证
	if err := c.ShouldBind(&requestPost); err != nil {
		log.Println(err.Error())
		response.Fail(c, "数据验证错误", nil)
		return
	}
	// 获取登陆用户user
	user, _ := c.Get("user")
	// 创建post
	post := model.Post{
		UserId:     user.(model.User).ID,
		CategoryId: requestPost.CategoryId,
		Title:      requestPost.Title,
		HeadImg:    requestPost.HeadImg,
		Content:    requestPost.Content,
	}

	if err := p.DB.Create(&post).Error; err != nil {
		panic(err)
		return
	}

	response.Success(c, "创建成功", nil)
}

func (p PostController) Update(c *gin.Context) {
	var requestPost vo.CreatePostRequest
	// 数据验证
	if err := c.ShouldBind(&requestPost); err != nil {
		log.Println(err.Error())
		response.Fail(c, "数据验证错误", nil)
		return
	}
	// 获取path中的id
	postId := c.Params.ByName("id")
	var post model.Post
	if err := p.DB.Where("id = ?", postId).First(&post).Error; err != nil {
		response.Fail(c, "文字不存在", nil)
		return
	}
	// 判断当前用户是否为文章的作者
	user, _ := c.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(c, "文字不属于你", nil)
		return
	}
	// 更新文字
	if err := p.DB.Model(&post).Updates(requestPost).Error; err != nil{
		response.Fail(c, "更新失败", nil)
		return
	}

	response.Success(c, "更新成功", gin.H{"post":post})
}

func (p PostController) Show(c *gin.Context) {
	// 获取path中的id
	postId := c.Params.ByName("id")
	var post model.Post
	if err := p.DB.Where("id = ?", postId).First(&post).Error; err != nil {
		response.Fail(c, "文字不存在", nil)
		return
	}
	response.Success(c, "查找成功", gin.H{"post":post})
}

func (p PostController) Delete(c *gin.Context) {
	// 获取path中的id
	postId := c.Params.ByName("id")
	var post model.Post
	if err := p.DB.Where("id = ?", postId).First(&post).Error; err != nil {
		response.Fail(c, "文字不存在", nil)
		return
	}

	// 判断当前用户是否为文章的作者
	user, _ := c.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(c, "文字不属于你", nil)
		return
	}

	// 删除
	p.DB.Delete(&post)

	response.Success(c, "删除成功", gin.H{"post": post})
}

func NewPostController() IPostController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	return PostController{DB: db}
}
