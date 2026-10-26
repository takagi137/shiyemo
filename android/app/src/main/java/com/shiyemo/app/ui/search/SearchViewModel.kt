package com.shiyemo.app.ui.search

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.shiyemo.app.data.model.dto.SearchResultDto
import com.shiyemo.app.data.repository.SearchRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Job
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

/**
 * 搜索页面 UI 状态
 */
data class SearchUiState(
    val query: String = "",
    val results: List<SearchResultDto> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null,
    val currentPage: Int = 1,
    val hasMore: Boolean = true,
    val isSearchPerformed: Boolean = false // 是否已执行过搜索
)

/**
 * 搜索 ViewModel
 */
@HiltViewModel
class SearchViewModel @Inject constructor(
    private val searchRepository: SearchRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(SearchUiState())
    val uiState: StateFlow<SearchUiState> = _uiState.asStateFlow()

    private var searchJob: Job? = null

    /**
     * 更新搜索关键词
     */
    fun updateQuery(query: String) {
        _uiState.value = _uiState.value.copy(query = query)
        // 自动搜索（防抖 500ms）
        if (query.isNotBlank()) {
            searchJob?.cancel()
            searchJob = viewModelScope.launch {
                delay(500)
                search()
            }
        }
    }

    /**
     * 执行搜索
     */
    fun search() {
        val query = _uiState.value.query.trim()
        if (query.isBlank()) return

        searchJob?.cancel()
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(
                isLoading = true,
                error = null,
                currentPage = 1,
                isSearchPerformed = true
            )

            try {
                val response = searchRepository.search(
                    query = query,
                    page = 1,
                    pageSize = PAGE_SIZE
                )

                if (response.code == 0) {
                    val results = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        results = results,
                        isLoading = false,
                        currentPage = 1,
                        hasMore = results.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = response.message
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "搜索失败"
                )
            }
        }
    }

    /**
     * 加载更多
     */
    fun loadMore() {
        val state = _uiState.value
        if (state.isLoading || !state.hasMore || state.query.isBlank()) return

        viewModelScope.launch {
            _uiState.value = state.copy(isLoading = true)

            try {
                val nextPage = state.currentPage + 1
                val response = searchRepository.search(
                    query = state.query,
                    page = nextPage,
                    pageSize = PAGE_SIZE
                )

                if (response.code == 0) {
                    val newResults = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        results = state.results + newResults,
                        isLoading = false,
                        currentPage = nextPage,
                        hasMore = newResults.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = response.message
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "加载更多失败"
                )
            }
        }
    }

    /**
     * 清除搜索
     */
    fun clearSearch() {
        searchJob?.cancel()
        _uiState.value = SearchUiState()
    }

    companion object {
        private const val PAGE_SIZE = 20
    }
}
