package com.shiyemo.app.data.model.dto

import com.google.gson.annotations.SerializedName

/**
 * 通用 API 响应包装
 */
data class ApiResponse<T>(
    @SerializedName("code")
    val code: Int = 0,

    @SerializedName("message")
    val message: String = "ok",

    @SerializedName("data")
    val data: T? = null
) {
    val isSuccess: Boolean get() = code == 0
}
