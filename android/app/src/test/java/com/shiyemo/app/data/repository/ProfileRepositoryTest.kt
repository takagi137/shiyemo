package com.shiyemo.app.data.repository

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.PostDto
import com.shiyemo.app.data.model.dto.UserDto
import com.shiyemo.app.data.remote.PostApi
import com.shiyemo.app.data.remote.UserApi
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.mockk
import kotlinx.coroutines.test.runTest
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

class ProfileRepositoryTest {

    private lateinit var userApi: UserApi
    private lateinit var postApi: PostApi
    private lateinit var repository: ProfileRepository

    @Before
    fun setup() {
        userApi = mockk()
        postApi = mockk()
        repository = ProfileRepository(userApi, postApi)
    }

    // ========== 关注/取关 ==========

    @Test
    fun `followUser 成功关注`() = runTest {
        // Given
        coEvery { userApi.followUser(10) } returns ApiResponse(
            code = 0, message = "ok", data = Unit
        )

        // When
        val response = repository.followUser(10)

        // Then
        assertEquals(0, response.code)
        coVerify(exactly = 1) { userApi.followUser(10) }
    }

    @Test
    fun `followUser 重复关注`() = runTest {
        // Given
        coEvery { userApi.followUser(10) } returns ApiResponse(
            code = 40001, message = "已经关注过该用户", data = null
        )

        // When
        val response = repository.followUser(10)

        // Then
        assertEquals(40001, response.code)
    }

    @Test
    fun `unfollowUser 成功取消关注`() = runTest {
        // Given
        coEvery { userApi.unfollowUser(10) } returns ApiResponse(
            code = 0, message = "ok", data = Unit
        )

        // When
        val response = repository.unfollowUser(10)

        // Then
        assertEquals(0, response.code)
        coVerify(exactly = 1) { userApi.unfollowUser(10) }
    }

    @Test
    fun `unfollowUser 未关注时取消`() = runTest {
        // Given
        coEvery { userApi.unfollowUser(10) } returns ApiResponse(
            code = 40001, message = "未关注该用户", data = null
        )

        // When
        val response = repository.unfollowUser(10)

        // Then
        assertEquals(40001, response.code)
    }

    // ========== 关注列表 ==========

    @Test
    fun `getFollowing 成功返回列表`() = runTest {
        // Given
        val users = listOf(
            UserDto(id = 2, nickname = "用户A"),
            UserDto(id = 3, nickname = "用户B")
        )
        coEvery { userApi.getFollowing(1, 1, 20) } returns ApiResponse(
            code = 0, message = "ok", data = users
        )

        // When
        val response = repository.getFollowing(userId = 1)

        // Then
        assertEquals(0, response.code)
        assertEquals(2, response.data?.size)
        assertEquals("用户A", response.data?.get(0)?.nickname)
    }

    @Test
    fun `getFollowing 空列表`() = runTest {
        // Given
        coEvery { userApi.getFollowing(1, 1, 20) } returns ApiResponse(
            code = 0, message = "ok", data = emptyList()
        )

        // When
        val response = repository.getFollowing(userId = 1)

        // Then
        assertEquals(0, response.code)
        assertTrue(response.data?.isEmpty() == true)
    }

    @Test
    fun `getFollowing 分页参数正确`() = runTest {
        // Given
        coEvery { userApi.getFollowing(1, 2, 10) } returns ApiResponse(
            code = 0, message = "ok", data = emptyList()
        )

        // When
        repository.getFollowing(userId = 1, page = 2, pageSize = 10)

        // Then
        coVerify { userApi.getFollowing(1, 2, 10) }
    }

    // ========== 粉丝列表 ==========

    @Test
    fun `getFollowers 成功返回列表`() = runTest {
        // Given
        val users = listOf(
            UserDto(id = 5, nickname = "粉丝A"),
            UserDto(id = 6, nickname = "粉丝B"),
            UserDto(id = 7, nickname = "粉丝C")
        )
        coEvery { userApi.getFollowers(1, 1, 20) } returns ApiResponse(
            code = 0, message = "ok", data = users
        )

        // When
        val response = repository.getFollowers(userId = 1)

        // Then
        assertEquals(0, response.code)
        assertEquals(3, response.data?.size)
    }

    @Test
    fun `getFollowers 空列表`() = runTest {
        // Given
        coEvery { userApi.getFollowers(1, 1, 20) } returns ApiResponse(
            code = 0, message = "ok", data = emptyList()
        )

        // When
        val response = repository.getFollowers(userId = 1)

        // Then
        assertEquals(0, response.code)
        assertTrue(response.data?.isEmpty() == true)
    }

    @Test
    fun `getFollowers 分页参数正确`() = runTest {
        // Given
        coEvery { userApi.getFollowers(1, 3, 5) } returns ApiResponse(
            code = 0, message = "ok", data = emptyList()
        )

        // When
        repository.getFollowers(userId = 1, page = 3, pageSize = 5)

        // Then
        coVerify { userApi.getFollowers(1, 3, 5) }
    }

    // ========== 用户信息 ==========

    @Test
    fun `getUser 成功返回用户信息`() = runTest {
        // Given
        coEvery { userApi.getUser(10) } returns ApiResponse(
            code = 0, message = "ok", data = UserDto(
                id = 10, nickname = "张三", bio = "失业中",
                followerCount = 5, isFollowing = true
            )
        )

        // When
        val response = repository.getUser(10)

        // Then
        assertEquals(0, response.code)
        assertEquals("张三", response.data?.nickname)
        assertEquals(true, response.data?.isFollowing)
        assertEquals(5, response.data?.followerCount)
    }

    @Test
    fun `getUser 用户不存在`() = runTest {
        // Given
        coEvery { userApi.getUser(999) } returns ApiResponse(
            code = 40201, message = "用户不存在", data = null
        )

        // When
        val response = repository.getUser(999)

        // Then
        assertEquals(40201, response.code)
        assertEquals(null, response.data)
    }

    // ========== 收藏列表 ==========

    @Test
    fun `getBookmarks 成功返回收藏`() = runTest {
        // Given
        coEvery { userApi.getBookmarks(1, 20) } returns ApiResponse(
            code = 0, message = "ok", data = listOf(
                PostDto(id = 1, title = "收藏的帖子", content = "内容")
            )
        )

        // When
        val response = repository.getBookmarks()

        // Then
        assertEquals(0, response.code)
        assertEquals(1, response.data?.size)
        assertEquals("收藏的帖子", response.data?.first()?.title)
    }

    // ========== 网络异常 ==========

    @Test
    fun `followUser 网络异常`() = runTest {
        // Given
        coEvery { userApi.followUser(any()) } throws RuntimeException("网络连接失败")

        // When & Then
        try {
            repository.followUser(1)
            assertTrue("应该抛出异常", false)
        } catch (e: RuntimeException) {
            assertEquals("网络连接失败", e.message)
        }
    }

    @Test
    fun `getFollowing 网络异常`() = runTest {
        // Given
        coEvery { userApi.getFollowing(any(), any(), any()) } throws RuntimeException("连接超时")

        // When & Then
        try {
            repository.getFollowing(userId = 1)
            assertTrue("应该抛出异常", false)
        } catch (e: RuntimeException) {
            assertEquals("连接超时", e.message)
        }
    }
}
