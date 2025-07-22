package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []string{}})
}

func HandleSubmit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []string{}})
}

func HandleGetData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []string{}})
}

type Problem struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	Level   string `json:"level"`
	Content string `json:"content"`
}

// 获取每日题目
func HandleDailyProblems(c *gin.Context) {
	// 这里可以改为从数据库获取
	dailyProblems := []Problem{
		{
			ID:      1,
			Title:   "两数之和",
			Link:    "https://leetcode.cn/problems/two-sum/",
			Level:   "easy",
			Content: "给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出和为目标值 target 的那两个整数，并返回它们的数组下标。",
		},
		{
			ID:      2,
			Title:   "两数相加",
			Link:    "https://leetcode.cn/problems/add-two-numbers/",
			Level:   "medium",
			Content: "给你两个非空的链表，表示两个非负的整数。它们每位数字都是按照逆序的方式存储的，并且每个节点只能存储一位数字。",
		},
		{
			ID:      3,
			Title:   "寻找两个正序数组的中位数",
			Link:    "https://leetcode.cn/problems/median-of-two-sorted-arrays/",
			Level:   "hard",
			Content: "给定两个大小分别为 m 和 n 的正序（从小到大）数组 nums1 和 nums2。请你找出并返回这两个正序数组的中位数。",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dailyProblems,
	})
}
