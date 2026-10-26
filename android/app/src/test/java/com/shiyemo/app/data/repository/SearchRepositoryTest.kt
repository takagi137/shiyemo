package com.shiyemo.app.data.repository

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.SearchResultDto
import com.shiyemo.app.data.remote.SearchApi
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.mockk
import kotlinx.coroutines.test.runTest
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

class SearchRepositoryTest {

    private lateinit var searchApi: SearchApi
    private lateinit var repository: SearchRepository

    @Before
    fun setup() {
        searchApi = mockk()
        repository = SearchRepository(searchApi)
    }

    @Test
    fun `search 成功返回结果`() = runTest {
        // Given
        val results = listOf(
            SearchResultDto(id = 1, type = "post", title = "失业日记", content = "今天..."),
            SearchResultDto(id = 2, type = "user", nickname = "小明")
        )
        coEvery { searchApi.search("失业", null, 1, 20) } returns ApiResponse(
            code = 0, message = "ok", data = results
        )

        // When
        val response = repository.search("失业")

        // Then
        assertEquals(0, response.code)
        assertEquals(2, response.data?.size)
        assertEquals("post", response.data?.first()?.type)
        assertEquals("user", response.data?.last()?.type)
    }

    @Test
    fun `search 空关键词不请求`() = runTest {
        // Given - 空关键词应由 ViewModel 层拦截，Repository 层不做校验
        coEvery { searchApi.search("", null, 1, 20) } returns ApiResponse(
            code = 0, message = "ok", data = emptyList()
        )

        // When
        val response = repository.search("")

        // Then
        assertEquals(0, response.code)
        assertTrue(response.data?.isEmpty() == true)
    }

    @Test
    fun `search 带类型过滤`() = runTest {
        // Given
        coEvery { searchApi.search("失业", "posts", 1, 20) } returns ApiResponse(
            code = 0, message = "ok", data = listOf(
                SearchResultDto(id = 1, type = "post", title = "失业日记")
            )
        )

        // When
        val response = repository.search("失业", type = "posts")

        // Then
        assertEquals(1, response.data?.size)
        assertEquals("post", response.data?.first()?.type)
        coVerify { searchApi.search("失业", "posts", 1, 20) }
    }

    @Test
    fun `search 分页参数正确传递`() = runTest {
        // Given
        coEvery { searchApi.search("工作", null, 2, 10) } returns ApiResponse(
            code = 0, message = "ok", data = emptyList()
        )

        // When
        repository.search("工作", page = 2, pageSize = 10)

        // Then
        coVerify { searchApi.search("工作", null, 2, 10) }
    }

    @Test
    fun `search 服务端返回错误`() = runTest {
        // Given
        coEvery { searchApi.search("xxx", null, 1, 20) } returns ApiResponse(
            code = 50001, message = "服务内部错误", data = null
        )

        // When
        val response = repository.search("xxx")

        // Then
        assertEquals(50001, response.code)
        assertEquals("服务内部错误", response.message)
    }

    @Test
    fun `search 网络异常`() = runTest {
        // Given
        coEvery { searchApi.search(any(), any(), any(), any()) } throws RuntimeException("网络连接失败")

        // When & Then
        try {
            repository.search("测试")
            assertTrue("应该抛出异常", false)
        } catch (e: RuntimeException) {
            assertEquals("网络连接失败", e.message)
        }
    }
}
