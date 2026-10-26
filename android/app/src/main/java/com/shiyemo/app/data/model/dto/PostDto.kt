package com.shiyemo.app.data.model.dto

import com.google.gson.annotations.SerializedName

/**
 * 帖子 DTO
 */
data class PostDto(
    @SerializedName("id")
    val id: Long = 0,

    @SerializedName("user_id")
    val userId: Long = 0,

    @SerializedName("user")
    val user: UserDto? = null,

    @SerializedName("type")
    val type: String = "",

    @SerializedName("title")
    val title: String = "",

    @SerializedName("content")
    val content: String = "",

    @SerializedName("images")
    val images: List<String> = emptyList(),

    @SerializedName("tags")
    val tags: List<String> = emptyList(),

    @SerializedName("like_count")
    val likeCount: Int = 0,

    @SerializedName("comment_count")
    val commentCount: Int = 0,

    @SerializedName("bookmark_count")
    val bookmarkCount: Int = 0,

    @SerializedName("is_liked")
    val isLiked: Boolean = false,

    @SerializedName("is_bookmarked")
    val isBookmarked: Boolean = false,

    @SerializedName("is_anonymous")
    val isAnonymous: Boolean = false,

    @SerializedName("anonymous_name")
    val anonymousName: String = "",

    @SerializedName("is_top")
    val isTop: Boolean = false,

    @SerializedName("is_featured")
    val isFeatured: Boolean = false,

    @SerializedName("created_at")
    val createdAt: String = "",

    @SerializedName("updated_at")
    val updatedAt: String = ""
)

/**
 * 创建帖子请求
 */
data class CreatePostRequest(
    @SerializedName("type")
    val type: String,

    @SerializedName("title")
    val title: String,

    @SerializedName("content")
    val content: String,

    @SerializedName("images")
    val images: List<String> = emptyList(),

    @SerializedName("tags")
    val tags: List<String> = emptyList(),

    @SerializedName("is_anonymous")
    val isAnonymous: Boolean = false
)

/**
 * 评论 DTO
 */
data class CommentDto(
    @SerializedName("id")
    val id: Long = 0,

    @SerializedName("post_id")
    val postId: Long = 0,

    @SerializedName("user_id")
    val userId: Long = 0,

    @SerializedName("user")
    val user: UserDto? = null,

    @SerializedName("content")
    val content: String = "",

    @SerializedName("is_anonymous")
    val isAnonymous: Boolean = false,

    @SerializedName("anonymous_name")
    val anonymousName: String = "",

    @SerializedName("created_at")
    val createdAt: String = ""
)

/**
 * 搜索结果 DTO
 */
data class SearchResultDto(
    @SerializedName("id")
    val id: Long = 0,

    @SerializedName("type")
    val type: String = "", // "post" or "user"

    @SerializedName("title")
    val title: String = "",

    @SerializedName("content")
    val content: String = "",

    @SerializedName("nickname")
    val nickname: String = "",

    @SerializedName("avatar")
    val avatar: String = "",

    @SerializedName("created_at")
    val createdAt: String = ""
)

/**
 * 通知 DTO
 */
data class NotificationDto(
    @SerializedName("id")
    val id: Long = 0,

    @SerializedName("type")
    val type: String = "", // like, comment, follow, system

    @SerializedName("content")
    val content: String = "",

    @SerializedName("from_user")
    val fromUser: UserDto? = null,

    @SerializedName("post_id")
    val postId: Long = 0,

    @SerializedName("is_read")
    val isRead: Boolean = false,

    @SerializedName("created_at")
    val createdAt: String = ""
)
