package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myVote/APP/model"
	"myVote/APP/tools"
	"net/http"
	"strconv"
	"time"
)

func Votes(c *gin.Context) {
	n := "游客"
	name, ok := model.GetSession(c)["name"]
	if ok && name != "" {
		n = name.(string)
	}
	c.HTML(200, "index.tmpl", gin.H{
		"name": n,
	})
	return
}

// GetVotes godoc
//
//	@Summary		获取投票列表
//	@Description	获取投票列表
//	@Tags			vote
//	@Accept			json
//	@Produce		json
//	@response		200,500	{object}	tools.HttpCode{data=[]model.Vote}
//	@Router			/votes [get]
func GetVotes(c *gin.Context) {
	ret := model.GetVotes()
	c.JSON(200, tools.HttpCode{
		Code: tools.OK,
		Data: ret,
	})
}

// GetVote godoc
//
//	@Summary		获取投票详情
//	@Description	根据ID获取某个用户的投票详情
//	@Tags			vote
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	false	"int valid"	minimum(1)
//	@response		200,500	{object}	tools.HttpCode{data=model.Vote}
//	@Router			/vote/:id [get]
func GetVote(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id <= 0 {
		c.JSON(404, tools.HttpCode{
			Code: tools.NotFound,
			Data: struct{}{},
		})
	}
	ret := model.GetVote(id)
	if ret.Id > 0 {
		fmt.Printf("vote:%+v\n", ret)
		c.JSON(200, tools.HttpCode{
			Code:    tools.OK,
			Message: "",
			Data:    ret,
		})
		return
	}
	c.JSON(404, tools.HttpCode{
		Code: tools.NotFound,
		Data: struct{}{},
	})
}

// AddVoteBasic godoc
//
//	@Summary		新增投票主体
//	@Description	新增加投票的主体
//	@Tags			vote
//	@Accept			json
//	@Produce		json
//	@Param			title		body		string	false	"投票主题"
//	@Param			start_time	body		int		false	"开始时间 时间戳"
//	@Param			type		body		int		false	"类型 0单选  1多选"
//	@Param			during		body		int		false	"持续时间 秒"
//	@response		200,400,500	{object}	tools.HttpCode
//	@Router			/vote/basic [post]
func AddVoteBasic(c *gin.Context) {
	vote := &model.Vote{}
	if err := c.ShouldBind(&vote); err != nil {
		c.JSON(200, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "数据解析失败",
			Data:    struct{}{},
		})
	}
	fmt.Printf("vote:%+v", vote)
	if uid, ok := c.Get("userId"); ok {
		vote.UserId = int(uid.(int64))
	}
	vote.CreatedTime = time.Now()
	err := model.CreateVote(vote)
	if err != nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "创建失败",
			Data:    struct{}{},
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code: tools.OK,
		Data: struct{}{},
	})
	return
}

// PutVoteBasic godoc
//
//	@Summary		更新投票主体
//	@Description	更新投票的主体
//	@Tags			vote
//	@Accept			json
//	@Produce		json
//	@Param			id			body		int		false	"投票主题ID"
//	@Param			title		body		string	false	"投票主题"
//	@Param			start_time	body		int		false	"开始时间 时间戳"
//	@Param			type		body		int		false	"类型 0单选  1多选"
//	@Param			during		body		int		false	"持续时间 秒"
//	@response		200,400,500	{object}	tools.HttpCode
//	@Router			/vote/basic [put]
func PutVoteBasic(c *gin.Context) {
	vote := &model.Vote{}
	if err := c.ShouldBind(&vote); err != nil || vote.UserId < 0 {
		c.JSON(200, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "数据解析失败,Id不能为空",
			Data:    struct{}{},
		})
	}
	fmt.Printf("vote:%+v", vote)
	if uid, ok := c.Get("userId"); ok {
		vote.UserId = int(uid.(int64))
	}
	vote.CreatedTime = time.Now()
	err := model.UpdateVote(vote)
	if err != nil {
		c.JSON(200, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "更新失败",
			Data:    struct{}{},
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code: tools.OK,
		Data: struct{}{},
	})
	return
}

// DeleteVoteBasic godoc
//
//	@Summary		删除投票主体
//	@Description	根据ID删除投票的主体
//	@Tags			vote
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int	false	"int valid"	minimum(1)
//	@response		200,404,500	{object}	tools.HttpCode{data=model.Vote}
//	@Router			/vote/:id [delete]
func DeleteVoteBasic(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id <= 0 {
		c.JSON(404, tools.HttpCode{
			Code: tools.NotFound,
			Data: struct{}{},
		})
		return
	}
	if err := model.DeleteVote(id); err == nil {
		c.JSON(200, tools.HttpCode{
			Code:    tools.OK,
			Message: "",
			Data:    struct{}{},
		})
		return
	}
	c.JSON(200, tools.HttpCode{
		Code:    tools.DoErr,
		Message: "删除失败",
		Data:    struct{}{},
	})
}

// DoVote godoc
//
//	@Summary		投票
//	@Description	根据用户投票信息执行投票操作
//	@Tags			vote
//	@Accept			json
//	@Produce		json
//	@Param			vote_id	formData	int		false	"投票主体ID"
//	@Param			opt[]	formData	[]int	false	"选项ID"
//	@response		200,500	{object}	tools.HttpCode{data=model.Vote}
//	@Router			/vote/do [post]
func DoVote(c *gin.Context) {
	voteIdStr, _ := c.GetPostForm("vote_id")
	opt, _ := c.GetPostFormArray("opt[]")
	id, _ := tools.Token.VerifyToken(c.GetHeader("Authorization"))
	fmt.Printf("vote:%s", voteIdStr)
	fmt.Printf("opt:%v", opt)
	fmt.Printf("id:%v", id.ID)
	optIds := make([]int64, 0)
	for _, v := range opt {
		tmp, _ := strconv.ParseInt(v, 10, 64)
		optIds = append(optIds, tmp)
	}
	voteId, _ := strconv.ParseInt(voteIdStr, 10, 64)
	userId := id.ID
	if model.DoVote(userId, voteId, optIds) {
		c.JSON(200, tools.HttpCode{
			Code: tools.OK,
			Data: VoteResult(voteId),
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.DoErr,
		Message: "投票失败！",
		Data:    struct{}{},
	})
	return
}
func VoteResult(id int64) *model.Vote {
	ret := model.GetVote(id)
	if ret.Id > 0 {
		for _, opt := range ret.VoteOpt {
			s, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(opt.Count)/float64(ret.Count)), 64)
			v, _ := strconv.ParseInt(fmt.Sprintf("%.f", s*100), 10, 64)
			opt.Percent = v
		}
	}
	return ret
}
