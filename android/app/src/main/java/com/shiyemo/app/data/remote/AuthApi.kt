package com.shiyemo.app.data.remote

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.AuthResponse
import com.shiyemo.app.data.model.dto.LoginEmailRequest
import com.shiyemo.app.data.model.dto.LoginPhoneRequest
import com.shiyemo.app.data.model.dto.RegisterEmailRequest
import retrofit2.http.Body
import retrofit2.http.POST

/**
 * 认证相关 API
 */
interface AuthApi {

    /**
     * 邮箱注册
     */
    @POST("auth/register/email")
    suspend fun registerEmail(@Body request: RegisterEmailRequest): ApiResponse<AuthResponse>

    /**
     * 邮箱登录
     */
    @POST("auth/login/email")
    suspend fun loginEmail(@Body request: LoginEmailRequest): ApiResponse<AuthResponse>

    /**
     * 手机号登录（验证码）
     */
    @POST("auth/login/phone")
    suspend fun loginPhone(@Body request: LoginPhoneRequest): ApiResponse<AuthResponse>

    /**
     * 发送验证码
     */
    @POST("auth/sms/send")
    suspend fun sendSmsCode(@Body request: Map<String, String>): ApiResponse<Unit>
}
