package com.shiyemo.app.data.remote

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.NotificationDto
import retrofit2.http.GET
import retrofit2.http.PUT
import retrofit2.http.Path
import retrofit2.http.Query

/**
 * 通知相关 API
 */
interface NotificationApi {

    /**
     * 获取通知列表
     * @param page 页码
     * @param pageSize 每页数量
     */
    @GET("notifications")
    suspend fun getNotifications(
        @Query("page") page: Int = 1,
        @Query("page_size") pageSize: Int = 20
    ): ApiResponse<List<NotificationDto>>

    /**
     * 标记单条通知为已读
     */
    @PUT("notifications/{id}/read")
    suspend fun markAsRead(@Path("id") notificationId: Long): ApiResponse<Unit>

    /**
     * 标记所有通知为已读
     */
    @PUT("notifications/read-all")
    suspend fun markAllAsRead(): ApiResponse<Unit>
}
