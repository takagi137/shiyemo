package com.shiyemo.app.data.repository

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.NotificationDto
import com.shiyemo.app.data.model.dto.UserDto
import com.shiyemo.app.data.remote.NotificationApi
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.mockk
import kotlinx.coroutines.test.runTest
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

class NotificationRepositoryTest {

    private lateinit var notificationApi: NotificationApi
    private lateinit var repository: NotificationRepository

    @Before
    fun setup() {
        notificationApi = mockk()
        repository = NotificationRepository(notificationApi)
    }

    @Test
    fun `getNotifications 成功返回列表`() = runTest {
        // Given
        val notifications = listOf(
            NotificationDto(id = 1, type = "like", content = "赞了你的帖子", isRead = false),
            NotificationDto(id = 2, type = "comment", content = "评论了你的帖子", isRead = true),
            NotificationDto(id = 3, type = "follow", content = "关注了你", isRead = false)
        )
        coEvery { notificationApi.getNotifications(1, 20) } returns ApiResponse(
            code = 0, message = "ok", data = notifications
        )

        // When
        val response = repository.getNotifications()

        // Then
        assertEquals(0, response.code)
        assertEquals(3, response.data?.size)
        assertEquals("like", response.data?.get(0)?.type)
        assertEquals("comment", response.data?.get(1)?.type)
        assertEquals("follow", response.data?.get(2)?.type)
    }

    @Test
    fun `getNotifications 空列表`() = runTest {
        // Given
        coEvery { notificationApi.getNotifications(1, 20) } returns ApiResponse(
            code = 0, message = "ok", data = emptyList()
        )

        // When
        val response = repository.getNotifications()

        // Then
        assertEquals(0, response.code)
        assertTrue(response.data?.isEmpty() == true)
    }

    @Test
    fun `getNotifications 分页参数正确`() = runTest {
        // Given
        coEvery { notificationApi.getNotifications(3, 10) } returns ApiResponse(
            code = 0, message = "ok", data = emptyList()
        )

        // When
        repository.getNotifications(page = 3, pageSize = 10)

        // Then
        coVerify { notificationApi.getNotifications(3, 10) }
    }

    @Test
    fun `getNotifications 包含发送者信息`() = runTest {
        // Given
        val sender = UserDto(id = 10, nickname = "张三")
        coEvery { notificationApi.getNotifications(1, 20) } returns ApiResponse(
            code = 0, message = "ok", data = listOf(
                NotificationDto(id = 1, type = "like", fromUser = sender)
            )
        )

        // When
        val response = repository.getNotifications()

        // Then
        assertEquals("张三", response.data?.first()?.fromUser?.nickname)
        assertEquals(10L, response.data?.first()?.fromUser?.id)
    }

    @Test
    fun `markAsRead 成功标记`() = runTest {
        // Given
        coEvery { notificationApi.markAsRead(42) } returns ApiResponse(
            code = 0, message = "ok", data = Unit
        )

        // When
        val response = repository.markAsRead(42)

        // Then
        assertEquals(0, response.code)
        coVerify(exactly = 1) { notificationApi.markAsRead(42) }
    }

    @Test
    fun `markAsRead 服务端错误`() = runTest {
        // Given
        coEvery { notificationApi.markAsRead(999) } returns ApiResponse(
            code = 40001, message = "通知不存在", data = null
        )

        // When
        val response = repository.markAsRead(999)

        // Then
        assertEquals(40001, response.code)
        assertEquals("通知不存在", response.message)
    }

    @Test
    fun `markAllAsRead 成功标记全部`() = runTest {
        // Given
        coEvery { notificationApi.markAllAsRead() } returns ApiResponse(
            code = 0, message = "ok", data = Unit
        )

        // When
        val response = repository.markAllAsRead()

        // Then
        assertEquals(0, response.code)
        coVerify(exactly = 1) { notificationApi.markAllAsRead() }
    }

    @Test
    fun `markAllAsRead 网络异常`() = runTest {
        // Given
        coEvery { notificationApi.markAllAsRead() } throws RuntimeException("连接超时")

        // When & Then
        try {
            repository.markAllAsRead()
            assertTrue("应该抛出异常", false)
        } catch (e: RuntimeException) {
            assertEquals("连接超时", e.message)
        }
    }

    @Test
    fun `getNotifications 服务端错误`() = runTest {
        // Given
        coEvery { notificationApi.getNotifications(any(), any()) } returns ApiResponse(
            code = 50001, message = "服务内部错误", data = null
        )

        // When
        val response = repository.getNotifications()

        // Then
        assertEquals(50001, response.code)
        assertEquals(null, response.data)
    }
}
