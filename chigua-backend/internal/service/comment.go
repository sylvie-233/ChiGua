package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"chigua-backend/internal/sql"
	"errors"
	"time"
)

var ErrCommentNoPermission = errors.New("无权限删除此评论")

func CreateComment(comment model.CommentCreate, userID int64) (*model.Comment, error) {
	now := time.Now()
	newComment := model.Comment{
		ParentID:    comment.ParentID,
		ArticleID:   comment.ArticleID,
		ReplyUserID: comment.ReplyUserID,
		UserID:      userID,
		Content:     comment.Content,
		CreatedAt:   now,
	}

	err := database.DB.QueryRow(sql.CommentInsert, newComment.ParentID, newComment.ArticleID, newComment.ReplyUserID, newComment.UserID, newComment.Content, newComment.CreatedAt).Scan(&newComment.ID)
	if err != nil {
		return nil, err
	}

	return &newComment, nil
}

func GetCommentsByArticleID(articleID int64) ([]model.CommentResponse, error) {
	var comments []model.Comment
	err := database.DB.Select(&comments, sql.CommentSelectByArticle, articleID)
	if err != nil {
		return nil, err
	}

	response := make([]model.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		userResp, err := getUserResponse(comment.UserID)
		if err != nil {
			return nil, err
		}

		var replyUserResp *model.UserResponse
		if comment.ReplyUserID != 0 {
			reply, err := getUserResponse(comment.ReplyUserID)
			if err != nil {
				return nil, err
			}
			replyUserResp = reply
		}

		response = append(response, model.CommentResponse{
			ID:          comment.ID,
			ParentID:    comment.ParentID,
			ArticleID:   comment.ArticleID,
			ReplyUserID: comment.ReplyUserID,
			UserID:      comment.UserID,
			Content:     comment.Content,
			CreatedAt:   comment.CreatedAt,
			User:        userResp,
			ReplyUser:   replyUserResp,
		})
	}

	return response, nil
}

func DeleteComment(id int64, userID int64) error {
	var comment model.Comment
	err := database.DB.Get(&comment, sql.CommentSelectUserID, id)
	if err != nil {
		return err
	}
	if comment.UserID != userID {
		return ErrCommentNoPermission
	}

	// TODO: 删除评论应该把二级评论也删除掉
	_, err = database.DB.Exec(sql.CommentDelete, id)
	return err
}

// getUserResponse 根据用户ID获取UserResponse（不包含密码）
func getUserResponse(userID int64) (*model.UserResponse, error) {
	var user model.User
	err := database.DB.Get(&user, sql.UserSelectByID, userID)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

// getCommentResponseWithUser 将Comment转换为CommentResponse（使用UserResponse）
func getCommentResponseWithUser(comment model.Comment) (model.CommentResponse, error) {
	userResp, err := getUserResponse(comment.UserID)
	if err != nil {
		return model.CommentResponse{}, err
	}

	var replyUserResp *model.UserResponse
	if comment.ReplyUserID != 0 {
		reply, err := getUserResponse(comment.ReplyUserID)
		if err != nil {
			return model.CommentResponse{}, err
		}
		replyUserResp = reply
	}

	return model.CommentResponse{
		ID:          comment.ID,
		ParentID:    comment.ParentID,
		ArticleID:   comment.ArticleID,
		ReplyUserID: comment.ReplyUserID,
		UserID:      comment.UserID,
		Content:     comment.Content,
		CreatedAt:   comment.CreatedAt,
		User:        userResp,
		ReplyUser:   replyUserResp,
	}, nil
}

// GetAllComments 获取所有评论列表（Admin用）
func GetAllComments(page, pageSize int, keyword string) (*model.CommentListResponse, error) {
	var total int64
	var comments []model.Comment
	var err error

	offset := (page - 1) * pageSize

	if keyword != "" {
		err = database.DB.Get(&total, sql.CommentCountAllByContent, keyword)
		if err != nil {
			return nil, err
		}
		err = database.DB.Select(&comments, sql.CommentSelectAllByContent, keyword, pageSize, offset)
	} else {
		err = database.DB.Get(&total, sql.CommentCountAll)
		if err != nil {
			return nil, err
		}
		err = database.DB.Select(&comments, sql.CommentSelectAll, pageSize, offset)
	}

	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	response := model.CommentListResponse{
		PageResponse: model.PageResponse{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
		Items: make([]model.CommentWithChildren, 0, len(comments)),
	}

	for _, comment := range comments {
		commentResponse, err := getCommentResponseWithUser(comment)
		if err != nil {
			return nil, err
		}

		response.Items = append(response.Items, model.CommentWithChildren{
			ID:          commentResponse.ID,
			ParentID:    commentResponse.ParentID,
			ArticleID:   commentResponse.ArticleID,
			ReplyUserID: commentResponse.ReplyUserID,
			UserID:      commentResponse.UserID,
			Content:     commentResponse.Content,
			CreatedAt:   commentResponse.CreatedAt,
			User:        commentResponse.User,
			ReplyUser:   commentResponse.ReplyUser,
		})
	}

	return &response, nil
}

// AdminDeleteComment Admin删除评论（不检查权限）
func AdminDeleteComment(id int64) error {
	_, err := database.DB.Exec(sql.CommentDelete, id)
	return err
}

// CountComment 统计评论总数
func CountComment() (int64, error) {
	var count int64
	err := database.DB.Get(&count, sql.CommentCountAll)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetCommentsWithPagination 获取两级评论分页列表
// articleID: 文章ID
// page: 页码（从1开始）
// pageSize: 每页大小（一级评论数量）
// childPageSize: 每个一级评论显示的二级评论数量
func GetCommentsWithPagination(articleID int64, page, pageSize, childPageSize int) (*model.CommentListResponse, error) {
	// 查询一级评论
	offset := (page - 1) * pageSize
	var firstLevelComments []model.Comment
	err := database.DB.Select(&firstLevelComments, sql.CommentSelectFirstLevel, articleID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// 统计一级评论总数
	var totalFirstLevel int
	err = database.DB.Get(&totalFirstLevel, sql.CommentSelectFirstLevelCount, articleID)
	if err != nil {
		return nil, err
	}

	// 计算一级评论分页信息
	totalPages := int((int64(totalFirstLevel) + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	// 组装响应
	response := model.CommentListResponse{
		PageResponse: model.PageResponse{
			Total:      int64(totalFirstLevel),
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
		Items: make([]model.CommentWithChildren, 0, len(firstLevelComments)),
	}

	for _, comment := range firstLevelComments {
		// 获取一级评论详情
		commentResponse, err := getCommentResponseWithUser(comment)
		if err != nil {
			return nil, err
		}

		// 查询二级评论（前childPageSize条）
		var secondLevelComments []model.Comment
		err = database.DB.Select(&secondLevelComments, sql.CommentSelectSecondLevel, articleID, comment.ID, childPageSize, 0)
		if err != nil {
			return nil, err
		}

		// 统计二级评论总数
		var totalSecondLevel int
		err = database.DB.Get(&totalSecondLevel, sql.CommentSelectSecondLevelCount, articleID, comment.ID)
		if err != nil {
			return nil, err
		}

		// 计算二级评论分页信息
		childTotalPages := int((int64(totalSecondLevel) + int64(childPageSize) - 1) / int64(childPageSize))
		childHasNext := 1 < childTotalPages
		childHasPrev := false

		// 转换二级评论
		children := make([]model.CommentResponse, 0, len(secondLevelComments))
		for _, child := range secondLevelComments {
			childResponse, err := getCommentResponseWithUser(child)
			if err != nil {
				return nil, err
			}
			children = append(children, childResponse)
		}

		response.Items = append(response.Items, model.CommentWithChildren{
			ID:          commentResponse.ID,
			ParentID:    commentResponse.ParentID,
			ArticleID:   commentResponse.ArticleID,
			ReplyUserID: commentResponse.ReplyUserID,
			UserID:      commentResponse.UserID,
			Content:     commentResponse.Content,
			CreatedAt:   commentResponse.CreatedAt,
			User:        commentResponse.User,
			ReplyUser:   commentResponse.ReplyUser,
			Children: model.SecondLevelCommentList{
				PageResponse: model.PageResponse{
					Total:      int64(totalSecondLevel),
					Page:       1,
					PageSize:   childPageSize,
					TotalPages: childTotalPages,
					HasNext:    childHasNext,
					HasPrev:    childHasPrev,
				},
				Items: children,
			},
		})
	}

	return &response, nil
}

// GetMoreSecondLevelComments 获取更多二级评论
// articleID: 文章ID
// parentID: 父评论ID
// page: 页码（从1开始）
// pageSize: 每页大小
func GetMoreSecondLevelComments(articleID, parentID int64, page, pageSize int) (*model.SecondLevelCommentList, error) {
	offset := (page - 1) * pageSize
	var comments []model.Comment
	err := database.DB.Select(&comments, sql.CommentSelectSecondLevel, articleID, parentID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// 统计二级评论总数
	var total int
	err = database.DB.Get(&total, sql.CommentSelectSecondLevelCount, articleID, parentID)
	if err != nil {
		return nil, err
	}

	// 计算分页信息
	totalPages := int((int64(total) + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	// 转换响应
	items := make([]model.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		commentResponse, err := getCommentResponseWithUser(comment)
		if err != nil {
			return nil, err
		}
		items = append(items, commentResponse)
	}

	return &model.SecondLevelCommentList{
		PageResponse: model.PageResponse{
			Total:      int64(total),
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
		Items: items,
	}, nil
}
