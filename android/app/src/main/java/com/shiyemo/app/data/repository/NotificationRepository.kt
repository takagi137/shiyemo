package com.shiyemo.app.data.repository

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.NotificationDto
import com.shiyemo.app.data.remote.NotificationApi
import javax.inject.Inject
import javax.inject.Singleton

/**
 * 通知数据仓库
 */
@Singleton
class NotificationRepository @Inject constructor(
    private val notificationApi: NotificationApi
) {
    /**
     * 获取通知列表
     */
    suspend fun getNotifications(
        page: Int = 1,
        pageSize: Int = 20
    ): ApiResponse<List<NotificationDto>> {
        return notificationApi.getNotifications(page = page, pageSize = pageSize)
    }

    /**
     * 标记单条通知为已读
     */
    suspend fun markAsRead(notificationId: Long): ApiResponse<Unit> {
        return notificationApi.markAsRead(notificationId)
    }

    /**
     * 标记所有通知为已读
     */
    suspend fun markAllAsRead(): ApiResponse<Unit> {
        return notificationApi.markAllAsRead()
    }
}
