package com.shiyemo.app.data.remote

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.CommentDto
import com.shiyemo.app.data.model.dto.CreatePostRequest
import com.shiyemo.app.data.model.dto.PostDto
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.PUT
import retrofit2.http.Path
import retrofit2.http.Query

/**
 * 帖子相关 API
 */
interface PostApi {

    /**
     * 获取帖子列表
     */
    @GET("posts")
    suspend fun getPosts(
        @Query("type") type: String? = null,
        @Query("page") page: Int = 1,
        @Query("page_size") pageSize: Int = 20
    ): ApiResponse<List<PostDto>>

    /**
     * 获取帖子详情
     */
    @GET("posts/{id}")
    suspend fun getPostDetail(@Path("id") postId: Long): ApiResponse<PostDto>

    /**
     * 创建帖子
     */
    @POST("posts")
    suspend fun createPost(@Body request: CreatePostRequest): ApiResponse<PostDto>

    /**
     * 更新帖子
     */
    @PUT("posts/{id}")
    suspend fun updatePost(
        @Path("id") postId: Long,
        @Body request: CreatePostRequest
    ): ApiResponse<PostDto>

    /**
     * 删除帖子
     */
    @DELETE("posts/{id}")
    suspend fun deletePost(@Path("id") postId: Long): ApiResponse<Unit>

    /**
     * 点赞/取消点赞
     */
    @POST("posts/{id}/like")
    suspend fun toggleLike(@Path("id") postId: Long): ApiResponse<Unit>

    /**
     * 收藏/取消收藏
     */
    @POST("posts/{id}/bookmark")
    suspend fun toggleBookmark(@Path("id") postId: Long): ApiResponse<Unit>

    /**
     * 获取关注动态 Feed
     */
    @GET("posts/following")
    suspend fun getFollowingFeed(
        @Query("page") page: Int = 1,
        @Query("page_size") pageSize: Int = 20
    ): ApiResponse<List<PostDto>>

    /**
     * 获取帖子评论
     */
    @GET("posts/{id}/comments")
    suspend fun getComments(
        @Path("id") postId: Long,
        @Query("page") page: Int = 1,
        @Query("page_size") pageSize: Int = 20
    ): ApiResponse<List<CommentDto>>

    /**
     * 发表评论
     */
    @POST("posts/{id}/comments")
    suspend fun createComment(
        @Path("id") postId: Long,
        @Body request: Map<String, String>
    ): ApiResponse<CommentDto>

    /**
     * 删除评论
     */
    @DELETE("comments/{id}")
    suspend fun deleteComment(@Path("id") commentId: Long): ApiResponse<Unit>
}
