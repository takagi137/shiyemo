package com.shiyemo.app.data.remote

import com.shiyemo.app.data.model.dto.ApiResponse
import com.shiyemo.app.data.model.dto.SearchResultDto
import retrofit2.http.GET
import retrofit2.http.Query

/**
 * 搜索相关 API
 */
interface SearchApi {

    /**
     * 搜索帖子和用户
     * @param query 搜索关键词
     * @param type 搜索类型过滤: posts / users, null 表示全部
     * @param page 页码
     * @param pageSize 每页数量
     */
    @GET("search")
    suspend fun search(
        @Query("q") query: String,
        @Query("type") type: String? = null,
        @Query("page") page: Int = 1,
        @Query("page_size") pageSize: Int = 20
    ): ApiResponse<List<SearchResultDto>>
}
