package com.shiyemo.app.data.repository

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.SearchResultDto
import com.shiyemo.app.data.remote.SearchApi
import javax.inject.Inject
import javax.inject.Singleton

/**
 * 搜索数据仓库
 */
@Singleton
class SearchRepository @Inject constructor(
    private val searchApi: SearchApi
) {
    /**
     * 搜索帖子和用户
     * @param query 搜索关键词
     * @param type 搜索类型: posts / users, null 表示全部
     * @param page 页码
     * @param pageSize 每页数量
     */
    suspend fun search(
        query: String,
        type: String? = null,
        page: Int = 1,
        pageSize: Int = 20
    ): ApiResponse<List<SearchResultDto>> {
        return searchApi.search(query = query, type = type, page = page, pageSize = pageSize)
    }
}
