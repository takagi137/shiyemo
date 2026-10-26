package com.shiyemo.app.data.repository

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.PostDto
import com.shiyemo.app.data.remote.PostApi
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Feed 数据仓库
 */
@Singleton
class FeedRepository @Inject constructor(
    private val postApi: PostApi
) {
    /**
     * 获取帖子列表
     * @param type 帖子类型过滤（emotion/insight/case/action），null 表示全部
     * @param page 页码
     * @param pageSize 每页数量
     */
    suspend fun getPosts(
        type: String? = null,
        page: Int = 1,
        pageSize: Int = 20
    ): ApiResponse<List<PostDto>> {
        return postApi.getPosts(type = type, page = page, pageSize = pageSize)
    }

    /**
     * 获取关注动态 Feed
     */
    suspend fun getFollowingFeed(
        page: Int = 1,
        pageSize: Int = 20
    ): ApiResponse<List<PostDto>> {
        return postApi.getFollowingFeed(page = page, pageSize = pageSize)
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
}
