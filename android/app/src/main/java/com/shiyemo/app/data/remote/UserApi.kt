package com.shiyemo.app.data.remote

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.PostDto
import com.shiyemo.app.data.model.dto.UserDto
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.PUT
import retrofit2.http.Path
import retrofit2.http.Query

/**
 * 用户相关 API
 */
interface UserApi {

    /**
     * 获取当前用户信息
     */
    @GET("users/me")
    suspend fun getCurrentUser(): ApiResponse<UserDto>

    /**
     * 获取指定用户信息
     */
    @GET("users/{id}")
    suspend fun getUser(@Path("id") userId: Long): ApiResponse<UserDto>

    /**
     * 更新用户信息
     */
    @PUT("users/me")
    suspend fun updateUser(@Body updates: Map<String, @JvmSuppressWildcards Any>): ApiResponse<UserDto>

    /**
     * 获取用户发布的帖子
     */
    @GET("users/{id}/posts")
    suspend fun getUserPosts(
        @Path("id") userId: Long,
        @Query("page") page: Int = 1,
        @Query("page_size") pageSize: Int = 20
    ): ApiResponse<List<PostDto>>

    /**
     * 获取收藏的帖子
     */
    @GET("users/me/bookmarks")
    suspend fun getBookmarks(
        @Query("page") page: Int = 1,
        @Query("page_size") pageSize: Int = 20
    ): ApiResponse<List<PostDto>>

    /**
     * 关注用户
     */
    @retrofit2.http.POST("users/{id}/follow")
    suspend fun followUser(@Path("id") userId: Long): ApiResponse<Unit>

    /**
     * 取消关注
     */
    @retrofit2.http.DELETE("users/{id}/follow")
    suspend fun unfollowUser(@Path("id") userId: Long): ApiResponse<Unit>

    /**
     * 获取关注列表
     */
    @GET("users/{id}/following")
    suspend fun getFollowing(
        @Path("id") userId: Long,
        @Query("page") page: Int = 1,
        @Query("page_size") pageSize: Int = 20
    ): ApiResponse<List<UserDto>>

    /**
     * 获取粉丝列表
     */
    @GET("users/{id}/followers")
    suspend fun getFollowers(
        @Path("id") userId: Long,
        @Query("page") page: Int = 1,
        @Query("page_size") pageSize: Int = 20
    ): ApiResponse<List<UserDto>>
}
