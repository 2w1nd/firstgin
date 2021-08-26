package controller

import (
	"com.w1nd/firstgin/model"
	"com.w1nd/firstgin/repository"
	"com.w1nd/firstgin/response"
	"com.w1nd/firstgin/vo"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ICategoryController interface {
	RestController
}

type CategoryController struct {
	Repository repository.CategoryRepository
}

func NewCategoryController() ICategoryController {
	repository := repository.NewCategoryRepository()
	repository.DB.AutoMigrate(model.Category{})

	return CategoryController{Repository: repository}
}

func (c2 CategoryController) Create(c *gin.Context) {
	var requestCategory vo.CreateCategoryRequest
	if err := c.ShouldBind(&requestCategory); err != nil {
		response.Fail(c, "数据验证错误，分类名称必填", nil)
		return
	}

	category, err := c2.Repository.Create(requestCategory.Name)
	if err != nil {
		panic(err)
		return
	}

	response.Success(c, "", gin.H{"category": category})
}

func (c2 CategoryController) Update(c *gin.Context) {
	var requestCategory vo.CreateCategoryRequest
	if err := c.ShouldBind(&requestCategory); err != nil {
		response.Fail(c, "数据验证错误，分类名称必填", nil)
		return
	}
	// 获取path中的参数
	categoryId, _ := strconv.Atoi(c.Params.ByName("id"))

	updateCategory, err := c2.Repository.SelectById(categoryId)
	if err != nil {
		response.Fail(c, "分类不存在", nil)
		return
	}

	// 更新分类
	category, err := c2.Repository.Update(*updateCategory, requestCategory.Name)
	if err != nil {
		panic(err)
	}
	response.Success(c, "修改成功", gin.H{"category": category})
}

func (c2 CategoryController) Show(c *gin.Context) {
	// 获取path中的参数
	categoryId, _ := strconv.Atoi(c.Params.ByName("id"))

	category, err := c2.Repository.SelectById(categoryId)
	if err != nil {
		response.Fail(c, "分类不存在", nil)
		return
	}

	response.Success(c, "", gin.H{"category": category})
}

func (c2 CategoryController) Delete(c *gin.Context) {
	// 获取path中的参数
	categoryId, _ := strconv.Atoi(c.Params.ByName("id"))

	if err := c2.Repository.DeleteById(categoryId); err != nil {
		response.Fail(c, "删除失败，请重试", nil)
		return
	}
	response.Success(c, "删除成功", nil)
}

