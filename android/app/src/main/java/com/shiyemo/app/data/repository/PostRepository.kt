package com.shiyemo.app.data.repository

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.CommentDto
import com.shiyemo.app.data.model.dto.PostDto
import com.shiyemo.app.data.remote.PostApi
import javax.inject.Inject
import javax.inject.Singleton

/**
 * 帖子详情数据仓库
 */
@Singleton
class PostRepository @Inject constructor(
    private val postApi: PostApi
) {
    /**
     * 获取帖子详情
     */
    suspend fun getPostDetail(postId: Long): ApiResponse<PostDto> {
        return postApi.getPostDetail(postId)
    }

    /**
     * 获取帖子评论
     */
    suspend fun getComments(
        postId: Long,
        page: Int = 1,
        pageSize: Int = 20
    ): ApiResponse<List<CommentDto>> {
        return postApi.getComments(postId, page = page, pageSize = pageSize)
    }

    /**
     * 发表评论
     */
    suspend fun createComment(
        postId: Long,
        content: String
    ): ApiResponse<CommentDto> {
        return postApi.createComment(postId, mapOf("content" to content))
    }

    /**
     * 点赞/取消点赞
     */
    suspend fun toggleLike(postId: Long): ApiResponse<Unit> {
        return postApi.toggleLike(postId)
    }

    /**
     * 收藏/取消收藏
     */
    suspend fun toggleBookmark(postId: Long): ApiResponse<Unit> {
        return postApi.toggleBookmark(postId)
    }

    /**
     * 删除评论
     */
    suspend fun deleteComment(commentId: Long): ApiResponse<Unit> {
        return postApi.deleteComment(commentId)
    }
}
