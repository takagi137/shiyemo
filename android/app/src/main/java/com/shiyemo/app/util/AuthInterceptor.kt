package com.shiyemo.app.util

import com.shiyemo.app.data.local.prefs.TokenManager
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.Response
import javax.inject.Inject
import javax.inject.Singleton

/**
 * JWT Token 拦截器
 * 自动在请求头中添加 Authorization: Bearer {token}
 */
@Singleton
class AuthInterceptor @Inject constructor(
    private val tokenManager: TokenManager
) : Interceptor {

    override fun intercept(chain: Interceptor.Chain): Response {
        val originalRequest = chain.request()

        // 从 DataStore 获取 token
        val token = runBlocking {
            tokenManager.tokenFlow.first()
        }

        // 如果有 token，添加到请求头
        val newRequest = if (!token.isNullOrEmpty()) {
            originalRequest.newBuilder()
                .header("Authorization", "Bearer $token")
                .build()
        } else {
            originalRequest
        }

        return chain.proceed(newRequest)
    }
}
