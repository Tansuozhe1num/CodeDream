package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Tansuozhe1num/codedream/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PostDTO 帖子数据传输对象
type PostDTO struct {
	ID         uint         `json:"id"`
	Title      string       `json:"title"`
	Content    string       `json:"content"`
	Author     string       `json:"author"`
	Type       string       `json:"type"`
	Votes      int          `json:"votes"`
	UserVote   *string      `json:"userVote"` // "up", "down", or null
	Bookmarked bool         `json:"bookmarked"`
	Time       string       `json:"time"`
	Tags       []string     `json:"tags"`
	Comments   []CommentDTO `json:"comments"`
}

// CommentDTO 评论数据传输对象
type CommentDTO struct {
	ID      uint   `json:"id"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

// ActiveUserDTO 活跃用户数据传输对象
type ActiveUserDTO struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Reputation int    `json:"reputation"`
}

// CreatePostRequest 创建帖子请求
type CreatePostRequest struct {
	Type    string   `json:"type" binding:"required,oneof=question solution discussion"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags"`
	Title   string   `json:"title"` // 可选，如果不提供则从content第一行提取
}

// AddCommentRequest 添加评论请求
type AddCommentRequest struct {
	PostID  uint   `json:"postId" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// VoteRequest 投票请求
type VoteRequest struct {
	PostID uint   `json:"postId" binding:"required"`
	Type   string `json:"type" binding:"required,oneof=up down"`
}

// BookmarkRequest 收藏请求
type BookmarkRequest struct {
	PostID uint `json:"postId" binding:"required"`
	State  bool `json:"state"`
}

// 获取当前用户ID（从认证中间件获取）
func getCurrentUserID(c *gin.Context, db *gorm.DB) uint {
	// 从认证中间件获取用户ID
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	// 如果没有认证，返回0（表示未登录）
	return 0
}

// formatTimeAgo 格式化时间为"X小时前"等格式
func formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "刚刚"
	}
	if diff < time.Hour {
		return strconv.Itoa(int(diff.Minutes())) + "分钟前"
	}
	if diff < 24*time.Hour {
		return strconv.Itoa(int(diff.Hours())) + "小时前"
	}
	if diff < 30*24*time.Hour {
		return strconv.Itoa(int(diff.Hours()/24)) + "天前"
	}
	return t.Format("2006-01-02")
}

// parseTags 解析JSON格式的标签
func parseTags(tagsJSON string) []string {
	if tagsJSON == "" {
		return []string{}
	}
	var tags []string
	json.Unmarshal([]byte(tagsJSON), &tags)
	return tags
}

// toPostDTO 转换为PostDTO
func toPostDTO(post model.Post, db *gorm.DB, currentUserID uint) PostDTO {
	// 解析标签
	tags := parseTags(post.Tags)

	// 获取评论
	var comments []model.Comment
	db.Where("post_id = ?", post.ID).
		Preload("User").
		Order("created_at ASC").
		Find(&comments)

	commentDTOs := make([]CommentDTO, len(comments))
	for i, c := range comments {
		commentDTOs[i] = CommentDTO{
			ID:      c.ID,
			Author:  c.User.Username,
			Content: c.Content,
			Time:    formatTimeAgo(c.CreatedAt),
		}
	}

	// 检查用户投票状态
	var userVote *string
	var vote model.Vote
	if err := db.Where("post_id = ? AND user_id = ?", post.ID, currentUserID).First(&vote).Error; err == nil {
		userVote = &vote.Type
	}

	// 检查是否收藏
	var bookmark model.Bookmark
	bookmarked := db.Where("post_id = ? AND user_id = ?", post.ID, currentUserID).First(&bookmark).Error == nil

	return PostDTO{
		ID:         post.ID,
		Title:      post.Title,
		Content:    post.Content,
		Author:     post.User.Username,
		Type:       post.Type,
		Votes:      post.VoteCount,
		UserVote:   userVote,
		Bookmarked: bookmarked,
		Time:       formatTimeAgo(post.CreatedAt),
		Tags:       tags,
		Comments:   commentDTOs,
	}
}

// HandleGetPosts 获取帖子列表
func HandleGetPosts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取查询参数
		postType := c.Query("type") // "question", "solution", "discussion", 或空（全部）
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("pageSize", "20")

		pageNum, _ := strconv.Atoi(page)
		pageSizeNum, _ := strconv.Atoi(pageSize)
		if pageNum < 1 {
			pageNum = 1
		}
		if pageSizeNum < 1 || pageSizeNum > 100 {
			pageSizeNum = 20
		}

		offset := (pageNum - 1) * pageSizeNum

		// 构建查询
		query := db.Model(&model.Post{}).
			Preload("User").
			Order("created_at DESC")

		if postType != "" && (postType == "question" || postType == "solution" || postType == "discussion") {
			query = query.Where("type = ?", postType)
		}

		// 查询帖子
		var posts []model.Post
		if err := query.Offset(offset).Limit(pageSizeNum).Find(&posts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "获取帖子列表失败: " + err.Error(),
			})
			return
		}

		// 转换为DTO
		currentUserID := getCurrentUserID(c, db)
		postDTOs := make([]PostDTO, len(posts))
		for i, post := range posts {
			postDTOs[i] = toPostDTO(post, db, currentUserID)
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    postDTOs,
		})
	}
}

// HandleCreatePost 创建新帖子
func HandleCreatePost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "请求参数错误: " + err.Error(),
			})
			return
		}

		currentUserID := getCurrentUserID(c, db)
		if currentUserID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "请先登录",
			})
			return
		}

		// 如果没有提供标题，从内容第一行提取
		title := req.Title
		if title == "" {
			// 从content第一行提取，最多50个字符
			contentRunes := []rune(req.Content)
			// 找到第一个换行符
			endIndex := len(contentRunes)
			for i, r := range contentRunes {
				if r == '\n' || r == '\r' {
					endIndex = i
					break
				}
			}
			// 提取第一行，最多50个字符
			if endIndex > 50 {
				title = string(contentRunes[:50]) + "..."
			} else {
				title = string(contentRunes[:endIndex])
			}
			if title == "" {
				title = "无标题"
			}
		}

		// 序列化标签
		tagsJSON, _ := json.Marshal(req.Tags)

		// 创建帖子
		post := model.Post{
			UserID:    currentUserID,
			Title:     title,
			Content:   req.Content,
			Type:      req.Type,
			Tags:      string(tagsJSON),
			VoteCount: 0,
			ViewCount: 0,
		}

		if err := db.Create(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "创建帖子失败: " + err.Error(),
			})
			return
		}

		// 加载关联数据
		db.Preload("User").First(&post, post.ID)

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "帖子创建成功",
			"postId":  post.ID,
			"data":    toPostDTO(post, db, currentUserID),
		})
	}
}

// HandleAddComment 添加评论
func HandleAddComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AddCommentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "请求参数错误: " + err.Error(),
			})
			return
		}

		// 检查帖子是否存在
		var post model.Post
		if err := db.First(&post, req.PostID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "帖子不存在",
			})
			return
		}

		currentUserID := getCurrentUserID(c, db)
		if currentUserID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "请先登录",
			})
			return
		}

		// 创建评论
		comment := model.Comment{
			PostID:  req.PostID,
			UserID:  currentUserID,
			Content: req.Content,
		}

		if err := db.Create(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "添加评论失败: " + err.Error(),
			})
			return
		}

		// 加载关联数据
		db.Preload("User").First(&comment, comment.ID)

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "评论添加成功",
			"data": CommentDTO{
				ID:      comment.ID,
				Author:  comment.User.Username,
				Content: comment.Content,
				Time:    formatTimeAgo(comment.CreatedAt),
			},
		})
	}
}

// HandlePostVote 帖子投票
func HandlePostVote(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VoteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "请求参数错误: " + err.Error(),
			})
			return
		}

		// 检查帖子是否存在
		var post model.Post
		if err := db.First(&post, req.PostID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "帖子不存在",
			})
			return
		}

		currentUserID := getCurrentUserID(c, db)
		if currentUserID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "请先登录",
			})
			return
		}

		// 查找现有投票
		var existingVote model.Vote
		err := db.Where("post_id = ? AND user_id = ?", req.PostID, currentUserID).First(&existingVote).Error

		if err == nil {
			// 如果投票类型相同，则取消投票
			if existingVote.Type == req.Type {
				// 删除投票
				db.Delete(&existingVote)
				// 更新帖子投票数
				if req.Type == "up" {
					post.VoteCount--
				} else {
					post.VoteCount++
				}
			} else {
				// 切换投票类型
				oldType := existingVote.Type
				existingVote.Type = req.Type
				db.Save(&existingVote)
				// 更新投票数：从down改为up，+2；从up改为down，-2
				if oldType == "down" && req.Type == "up" {
					post.VoteCount += 2
				} else if oldType == "up" && req.Type == "down" {
					post.VoteCount -= 2
				}
			}
		} else {
			// 创建新投票
			vote := model.Vote{
				PostID: req.PostID,
				UserID: currentUserID,
				Type:   req.Type,
			}
			db.Create(&vote)
			// 更新帖子投票数
			if req.Type == "up" {
				post.VoteCount++
			} else {
				post.VoteCount--
			}
		}

		// 保存帖子
		db.Save(&post)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "投票成功",
			"votes":   post.VoteCount,
		})
	}
}

// HandleBookmark 收藏帖子
func HandleBookmark(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BookmarkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "请求参数错误: " + err.Error(),
			})
			return
		}

		// 检查帖子是否存在
		var post model.Post
		if err := db.First(&post, req.PostID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "帖子不存在",
			})
			return
		}

		currentUserID := getCurrentUserID(c, db)
		if currentUserID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "请先登录",
			})
			return
		}

		// 查找现有收藏
		var bookmark model.Bookmark
		err := db.Where("post_id = ? AND user_id = ?", req.PostID, currentUserID).First(&bookmark).Error

		if req.State {
			// 收藏
			if err != nil {
				// 创建新收藏
				bookmark = model.Bookmark{
					PostID: req.PostID,
					UserID: currentUserID,
				}
				db.Create(&bookmark)
			}
		} else {
			// 取消收藏
			if err == nil {
				db.Delete(&bookmark)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "收藏状态更新成功",
		})
	}
}

// HandleHotPosts 获取热门帖子
func HandleHotPosts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := c.DefaultQuery("limit", "5")
		limitNum, _ := strconv.Atoi(limit)
		if limitNum < 1 || limitNum > 50 {
			limitNum = 5
		}

		var posts []model.Post
		if err := db.Model(&model.Post{}).
			Preload("User").
			Order("vote_count DESC, created_at DESC").
			Limit(limitNum).
			Find(&posts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "获取热门帖子失败: " + err.Error(),
			})
			return
		}

		// 转换为DTO（简化版，不需要完整信息）
		hotPosts := make([]map[string]interface{}, len(posts))
		for i, post := range posts {
			hotPosts[i] = map[string]interface{}{
				"id":    post.ID,
				"title": post.Title,
				"votes": post.VoteCount,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    hotPosts,
		})
	}
}

// HandleActiveUsers 获取活跃用户
func HandleActiveUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := c.DefaultQuery("limit", "5")
		limitNum, _ := strconv.Atoi(limit)
		if limitNum < 1 || limitNum > 50 {
			limitNum = 5
		}

		// 查询最近30天最活跃的用户（按发帖数+评论数排序）
		var users []model.User
		if err := db.Model(&model.User{}).
			Select("users.*, COUNT(DISTINCT posts.id) + COUNT(DISTINCT comments.id) as activity_count").
			Joins("LEFT JOIN posts ON posts.user_id = users.id AND posts.created_at > DATE_SUB(NOW(), INTERVAL 30 DAY)").
			Joins("LEFT JOIN comments ON comments.user_id = users.id AND comments.created_at > DATE_SUB(NOW(), INTERVAL 30 DAY)").
			Group("users.id").
			Order("activity_count DESC, reputation DESC").
			Limit(limitNum).
			Find(&users).Error; err != nil {
			// 如果上面的查询失败，使用简单的按声望排序
			db.Model(&model.User{}).
				Order("reputation DESC").
				Limit(limitNum).
				Find(&users)
		}

		// 转换为DTO
		activeUsers := make([]ActiveUserDTO, len(users))
		for i, user := range users {
			activeUsers[i] = ActiveUserDTO{
				ID:         user.ID,
				Name:       user.Username,
				Reputation: user.Reputation,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    activeUsers,
		})
	}
}
