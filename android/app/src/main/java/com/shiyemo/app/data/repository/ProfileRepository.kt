package com.shiyemo.app.data.repository

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.PostDto
import com.shiyemo.app.data.model.dto.UserDto
import com.shiyemo.app.data.remote.PostApi
import com.shiyemo.app.data.remote.UserApi
import javax.inject.Inject
import javax.inject.Singleton

/**
 * 个人中心数据仓库
 */
@Singleton
class ProfileRepository @Inject constructor(
    private val userApi: UserApi,
    private val postApi: PostApi
) {
    /**
     * 获取当前用户信息
     */
    suspend fun getCurrentUser(): ApiResponse<UserDto> {
        return userApi.getCurrentUser()
    }

    /**
     * 更新用户信息
     */
    suspend fun updateUser(updates: Map<String, Any>): ApiResponse<UserDto> {
        return userApi.updateUser(updates)
    }

    /**
     * 获取用户发布的帖子
     */
    suspend fun getUserPosts(
        userId: Long,
        page: Int = 1,
        pageSize: Int = 20
    ): ApiResponse<List<PostDto>> {
        return userApi.getUserPosts(userId, page = page, pageSize = pageSize)
    }

    /**
     * 获取用户收藏的帖子
     */
    suspend fun getBookmarks(
        page: Int = 1,
        pageSize: Int = 20
    ): ApiResponse<List<PostDto>> {
        return userApi.getBookmarks(page = page, pageSize = pageSize)
    }

    /**
     * 获取指定用户信息
     */
    suspend fun getUser(userId: Long): ApiResponse<UserDto> {
        return userApi.getUser(userId)
    }

    /**
     * 关注用户
     */
    suspend fun followUser(userId: Long): ApiResponse<Unit> {
        return userApi.followUser(userId)
    }

    /**
     * 取消关注
     */
    suspend fun unfollowUser(userId: Long): ApiResponse<Unit> {
        return userApi.unfollowUser(userId)
    }

    /**
     * 获取关注列表
     */
    suspend fun getFollowing(
        userId: Long,
        page: Int = 1,
        pageSize: Int = 20
    ): ApiResponse<List<UserDto>> {
        return userApi.getFollowing(userId, page = page, pageSize = pageSize)
    }

    /**
     * 获取粉丝列表
     */
    suspend fun getFollowers(
        userId: Long,
        page: Int = 1,
        pageSize: Int = 20
    ): ApiResponse<List<UserDto>> {
        return userApi.getFollowers(userId, page = page, pageSize = pageSize)
    }
}
