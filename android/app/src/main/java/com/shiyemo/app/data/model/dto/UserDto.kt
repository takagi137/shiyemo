package com.shiyemo.app.data.model.dto

import com.google.gson.annotations.SerializedName

/**
 * 用户信息 DTO
 */
data class UserDto(
    @SerializedName("id")
    val id: Long = 0,

    @SerializedName("nickname")
    val nickname: String = "",

    @SerializedName("email")
    val email: String = "",

    @SerializedName("phone")
    val phone: String = "",

    @SerializedName("avatar")
    val avatar: String = "",

    @SerializedName("bio")
    val bio: String = "",

    @SerializedName("company")
    val company: String = "",

    @SerializedName("occupation")
    val occupation: String = "",

    @SerializedName("city")
    val city: String = "",

    @SerializedName("industry")
    val industry: String = "",

    @SerializedName("freedom_days")
    val freedomDays: Int = 0,

    @SerializedName("post_count")
    val postCount: Int = 0,

    @SerializedName("following_count")
    val followingCount: Int = 0,

    @SerializedName("follower_count")
    val followerCount: Int = 0,

    @SerializedName("is_following")
    val isFollowing: Boolean = false,

    @SerializedName("created_at")
    val createdAt: String = ""
)

/**
 * 登录/注册请求
 */
data class LoginEmailRequest(
    @SerializedName("email")
    val email: String,

    @SerializedName("password")
    val password: String
)

data class RegisterEmailRequest(
    @SerializedName("email")
    val email: String,

    @SerializedName("password")
    val password: String,

    @SerializedName("nickname")
    val nickname: String
)

data class LoginPhoneRequest(
    @SerializedName("phone")
    val phone: String,

    @SerializedName("code")
    val code: String
)

/**
 * 认证响应
 */
data class AuthResponse(
    @SerializedName("token")
    val token: String,

    @SerializedName("user")
    val user: UserDto
)
